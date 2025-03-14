package daemons

import (
	"context"
	"iter"
	"log"
	"net/http"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/james-lawrence/deeppool/internal/x/backoffx"
	"github.com/james-lawrence/deeppool/internal/x/contextx"
	"github.com/james-lawrence/deeppool/internal/x/errorsx"
	"github.com/james-lawrence/deeppool/internal/x/httpx"
	"github.com/james-lawrence/deeppool/internal/x/sqlx"
	"github.com/james-lawrence/deeppool/rss"
	"github.com/james-lawrence/deeppool/tracking"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/james-lawrence/torrent/storage"
	"golang.org/x/time/rate"
)

func DiscoverFromRSSFeeds(ctx context.Context, q sqlx.Queryer, tclient *torrent.Client, tstore storage.ClientImpl) (err error) {
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

				if err = tracking.MetadataInsertWithDefaults(ctx, q, tracking.NewMetadata(&md.InfoHash, tracking.MetadataOptionFromInfo(mi), tracking.MetadataOptionDescription(item.Title), autodownload)).Scan(&meta); err != nil {
					log.Println("unable to record torrent metadata", feed.ID, err)
					continue
				}

				// log.Println("recorded", feed.ID, meta.ID, meta.Description)
			}

			if err = tracking.RSSCooldownByID(fctx, q, feed.ID, channel.TTL).Scan(&feed); err != nil {
				log.Println("unable to mark rss feed for cooldown", err)
				continue
			}

			// begin any torrent provided by this feed
			ResumeDownloads(ctx, q, tclient, tstore)
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
