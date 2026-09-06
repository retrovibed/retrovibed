package daemons

import (
	"context"
	"iter"
	"log"
	"net/http"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/storage"
	"github.com/retrovibed/retrovibed/retroapi/backoffx"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/userx"
	"github.com/retrovibed/retrovibed/retroapi/uuidx"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/contextx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/slicesx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/rss"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

func PrepareDefaultFeeds(ctx context.Context, q sqlx.Queryer) error {
	var (
		feeds []tracking.RSS
	)

	log.Println("syncing default rss feeds initialized")
	defer log.Println("syncing default rss feeds completed")

	encoded, err := fsx.AutoCached(userx.DefaultConfigDir(userx.DefaultRelRoot(), "default.feeds.json"), func() (_ []byte, _ error) {
		return jsonx.Marshal([]tracking.RSS{
			{
				Description:  "Arch Linux - iso",
				URL:          "https://archlinux.org/feeds/releases/",
				Contributing: true,
			},

			{
				Description:  "Retrovibed - test data",
				URL:          "https://vibed.community.retrovibe.space",
				Contributing: true,
			},
		})

	})
	if err != nil {
		return err
	}

	if err = jsonx.Unmarshal(encoded, &feeds); err != nil {
		return err
	}

	for _, feed := range feeds {
		feed = langx.Clone(feed, tracking.RSSOptionDefaultFeeds(feed), tracking.RSSOptionDefaultEncryptionSeed)
		if err = tracking.RSSInsertDefaultFeed(ctx, q, feed).Scan(&feed); err != nil {
			return errorsx.Wrapf(err, "feed creation failed: %s - %s", feed.Description, feed.URL)
		}
	}

	return nil
}

func DiscoverFromRSSFeedsOnce(
	ctx context.Context,
	q sqlx.Queryer,
	c *http.Client,
	rootstore fsx.Virtual,
	mc library.QueryCleaner,
	tclient *torrent.Client,
	tstore storage.ClientImpl,
	pub *asyncx.Wakeup,
) (err error) {
	const defaultttl = 1440 // 1 day in minutes
	queryfeeds := func(ctx context.Context, done context.CancelCauseFunc) iter.Seq[tracking.RSS] {
		return func(yield func(tracking.RSS) bool) {
			query := tracking.RSSSearchBuilder().Where(
				squirrel.And{
					tracking.RSSQueryNeedsCheck(),
				},
			).Limit(128)

			qiter := sqlx.Scan(tracking.RSSSearch(ctx, q, query))

			for p := range qiter.Iter() {
				if !yield(p) {
					break
				}
			}

			done(qiter.Err())
		}
	}

	fctx, fdone := context.WithCancelCause(ctx)
	for feed := range queryfeeds(fctx, fdone) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, feed.URL, nil)
		if err != nil {
			log.Println("unable to build feed request", feed.ID, err)
			continue
		}

		resp, err := httpx.AsError(c.Do(req))
		if err != nil {
			log.Println("unable to retrieve feed", feed.ID, err)
			if err = tracking.RSSCooldownByID(fctx, q, feed.ID, defaultttl, feed.Digest, feed.LastBuiltAt).Scan(&feed); err != nil {
				log.Println("unable to mark rss feed for cooldown", err)
			}
			continue
		}

		digest, channel, items, err := rss.Parse(ctx, resp.Body)
		if err != nil {
			log.Println("unable to parse feed", feed.ID, err)
			continue
		}
		md5digest := md5x.FormatUUID(digest)

		if v := channel.LastBuildDate.Timestamp(time.Now()); md5digest == feed.Digest {
			log.Println("torrent rss feed has not updated since last check", feed.ID, feed.Description, channel.TTL, md5digest, "==", feed.Digest)
			if err = tracking.RSSCooldownByID(fctx, q, feed.ID, langx.FirstNonZero(channel.TTL, defaultttl), md5digest, v).Scan(&feed); err != nil {
				log.Println("unable to mark rss feed for cooldown", err)
			}
			continue
		} else {
			log.Println("torrent rss feed changes detected", feed.ID, feed.Description, "fetching", len(items), "torrents", md5digest, "!=", feed.Digest)
		}

		importer := tracking.NewURIImport(q, c, rootstore)
		encryptionseed := uuidx.FirstNonNil(
			uuid.FromStringOrNil(feed.EncryptionSeed),
			uuid.FromStringOrNil(channel.Retrovibed.Entropy),
		).Bytes()

		for _, item := range items {
			uri := slicesx.FirstOrDefault(rss.ItemToEnclosure(item, mimex.Bittorrent), rss.FindBittorrentEnclosures(item)...)
			if uri.URL == "" {
				continue
			}

			meta, err := importer.Import(
				fctx,
				uri.URL,
				tracking.MetadataOptionMimetype(stringsx.FirstNonBlank(channel.Retrovibed.Mimetype, uri.Mimetype)),
				tracking.MetadataOptionDescription(item.Title),
				tracking.MetadataOptionKnownMediaID(uuid.Max.String()),
				tracking.MetadataOptionAutoEntropySeed(encryptionseed),
				tracking.MetadataOptionAutoArchive(feed.Autoarchive),
				tracking.MetadataOptionAutoBytes(uri.Length),
			)
			if err != nil {
				log.Println("unable to import uri:", feed.ID, uri.URL, err)
				continue
			}

			if feed.Autodownload {
				log.Println("marking torrent to be automatically downloaded", meta.Description, feed.Autodownload)
				if err := tracking.MetadataAutoDownloadByID(ctx, q, meta.ID).Scan(&meta); err != nil {
					log.Println("unable to mark torrent for automatic download:", feed.ID, err)
					continue
				}
			}
		}

		if updated := stringsx.FirstNonBlank(feed.Description, channel.Title); updated != feed.Description {
			feed.Description = updated
			feed.Digest = md5digest

			if cause := tracking.RSSInsertWithDefaults(fctx, q, feed).Scan(&feed); cause != nil {
				log.Println("failed to update rss feed", cause)
				continue
			}
		} else {
			if err = tracking.RSSCooldownByID(fctx, q, feed.ID, langx.FirstNonZero(channel.TTL, defaultttl), md5digest, channel.LastBuildDate.Timestamp(time.Now())).Scan(&feed); err != nil {
				log.Println("unable to mark rss feed for cooldown", err)
				continue
			}
		}

		log.Println("starting any downloads", feed.Description)
		// begin any torrent provided by this feed
		ResumeDownloads(ctx, q, rootstore, mc, tclient, tstore, pub)
	}

	if err := fctx.Err(); contextx.IgnoreCancelled(err) != nil {
		return err
	}

	return nil
}

// retrieve torrents from rss feeds.
func DiscoverFromRSSFeeds(ctx context.Context, q sqlx.Queryer, c *http.Client, rootstore fsx.Virtual, mc library.QueryCleaner, tclient *torrent.Client, tstore storage.ClientImpl, pub *asyncx.Wakeup) (err error) {
	bs := backoffx.New(
		backoffx.Exponential(time.Minute),
		backoffx.Maximum(15*time.Minute),
	)

	for attempts := 0; true; attempts++ {
		if err := context.Cause(ctx); err != nil {
			return err
		}

		if c := errorsx.Zero(sqlx.Count(ctx, q, "SELECT COUNT (*) FROM torrents_feed_rss WHERE next_check < NOW()")); c == 0 {
			time.Sleep(bs.Backoff(attempts))
			continue
		} else {
			attempts = -1
		}

		if err := DiscoverFromRSSFeedsOnce(ctx, q, c, rootstore, mc, tclient, tstore, pub); err != nil {
			log.Println("failed to discover torrents", err)
			continue
		}
	}

	return nil
}
