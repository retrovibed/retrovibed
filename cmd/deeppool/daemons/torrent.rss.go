package daemons

import (
	"context"
	"iter"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/james-lawrence/deeppool/internal/x/backoffx"
	"github.com/james-lawrence/deeppool/internal/x/errorsx"
	"github.com/james-lawrence/deeppool/internal/x/sqlx"
	"github.com/james-lawrence/deeppool/tracking"
)

func DiscoverFromRSSFeeds(ctx context.Context, q sqlx.Queryer) (err error) {
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

		fctx, fdone := context.WithCancelCause(ctx)
		for feed := range queryfeeds(fctx, fdone) {
			if err = tracking.RSSCooldown(fctx, q, feed).Scan(&feed); err != nil {
				log.Println("unable to mark rss feed for cooldown", err)
				continue
			}
		}

		if err := fctx.Err(); err != nil {
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
