package daemons

import (
	"context"
	"errors"
	"iter"
	"log"
	"net/http"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/james-lawrence/torrent/storage"
	"github.com/retrovibed/retrovibed/internal/backoffx"
	"github.com/retrovibed/retrovibed/internal/contextx"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/fsx"
	"github.com/retrovibed/retrovibed/internal/httpx"
	"github.com/retrovibed/retrovibed/internal/md5x"
	"github.com/retrovibed/retrovibed/internal/sqlx"
	"github.com/retrovibed/retrovibed/internal/stringsx"
	"github.com/retrovibed/retrovibed/rss"
	"github.com/retrovibed/retrovibed/tracking"
	"golang.org/x/time/rate"
)

func PrepareDefaultFeeds(ctx context.Context, q sqlx.Queryer) error {
	feedcreate := func(description, url string) (err error) {
		feed := tracking.RSS{
			ID:           md5x.FormatUUID(md5x.Digest(url)),
			Description:  description,
			URL:          url,
			Contributing: true,
		}

		if err = tracking.RSSInsertWithDefaults(ctx, q, feed).Scan(&feed); err != nil {
			return errorsx.Wrapf(err, "feed creation failed: %s - %s", description, url)
		}

		return nil
	}

	return errors.Join(
		feedcreate("Arch Linux", "https://archlinux.org/feeds/releases/"),
	)
}

// retrieve torrents from rss feeds.
func DiscoverFromRSSFeeds(ctx context.Context, q sqlx.Queryer, rootstore fsx.Virtual, tclient *torrent.Client, tstore storage.ClientImpl) (err error) {
	queryfeeds := func(ctx context.Context, done context.CancelCauseFunc) iter.Seq[tracking.RSS] {
		return func(yield func(tracking.RSS) bool) {
			query := tracking.RSSSearchBuilder().Where(
				squirrel.And{
					tracking.RSSQueryNeedsCheck(),
				},
			).Limit(128)

			scanner := tracking.RSSSearch(ctx, q, query)
			defer scanner.Close()

			for scanner.Next() {
				var (
					p tracking.RSS
				)

				if err := scanner.Scan(&p); err != nil {
					done(err)
					return
				}

				if !yield(p) {
					return
				}
			}

			if err := scanner.Err(); err != nil {
				done(err)
				return
			}

			done(nil)
		}
	}

	bs := backoffx.New(
		backoffx.Exponential(time.Minute),
		backoffx.Maximum(15*time.Minute),
	)

	for attempts := 0; true; attempts++ {
		if c := errorsx.Zero(sqlx.Count(ctx, q, "SELECT COUNT (*) FROM torrents_feed_rss WHERE next_check < NOW()")); c == 0 {
			time.Sleep(bs.Backoff(attempts))
			continue
		} else {
			attempts = -1
		}

		c := httpx.BindRetryTransport(http.DefaultClient, http.StatusTooManyRequests, http.StatusBadGateway)
		l := rate.NewLimiter(rate.Every(3*time.Second), 1)

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
				if err = tracking.RSSCooldownByID(fctx, q, feed.ID, 10).Scan(&feed); err != nil {
					log.Println("unable to mark rss feed for cooldown", err)
				}
				continue
			}

			channel, items, err := rss.Parse(ctx, resp.Body)
			if err != nil {
				log.Println("unable to parse feed", feed.ID, err)
				continue
			}

			autodownload := tracking.MetadataOptionNoop
			if feed.Autodownload {
				autodownload = tracking.MetadataOptionInitiate
			}

			for _, item := range items {
				var (
					meta tracking.Metadata
				)

				if err = l.Wait(ctx); err != nil {
					log.Println("rate limit failure", err)
					continue
				}

				req, err := http.NewRequestWithContext(ctx, http.MethodGet, item.Link, nil)
				if err != nil {
					log.Println("unable to build torrent request", feed.ID, err)
					continue
				}

				resp, err := httpx.AsError(http.DefaultClient.Do(req))
				if err != nil {
					log.Println("unable to retrieve feed", feed.ID, err)
					continue
				}

				mi, err := metainfo.NewFromReader(resp.Body)
				if err != nil {
					log.Println("unable to read metainfo from response", feed.ID, err)
					continue
				}
				md, err := torrent.NewFromInfo(*mi)
				if err != nil {
					log.Println("unable to read metainfo from response", feed.ID, err)
					continue
				}

				if err = tracking.MetadataInsertWithDefaults(ctx, q, tracking.NewMetadata(&md.ID, tracking.MetadataOptionFromInfo(mi), tracking.MetadataOptionDescription(item.Title), autodownload)).Scan(&meta); err != nil {
					log.Println("unable to record torrent metadata", feed.ID, err)
					continue
				}

				// log.Println("recorded", feed.ID, meta.ID, meta.Description)
			}

			if updated := stringsx.FirstNonBlank(feed.Description, channel.Title); updated != feed.Description {
				feed.Description = updated
				if cause := tracking.RSSInsertWithDefaults(fctx, q, feed).Scan(&feed); cause != nil {
					log.Println("failed to update rss feed", cause)
					continue
				}
			} else {
				if err = tracking.RSSCooldownByID(fctx, q, feed.ID, channel.TTL).Scan(&feed); err != nil {
					log.Println("unable to mark rss feed for cooldown", err)
					continue
				}
			}

			// begin any torrent provided by this feed
			ResumeDownloads(ctx, q, rootstore, tclient, tstore)
		}

		if err := fctx.Err(); contextx.IgnoreCancelled(err) != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}

	return nil
}
