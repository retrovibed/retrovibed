package daemons

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/storage"
	"github.com/retrovibed/retrovibed/retroapi/userx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/asynccompute"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/envx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"golang.org/x/time/rate"
)

type discoverMediaOptions struct {
	frequency   time.Duration
	peerTimeout time.Duration
	infoTimeout time.Duration
}

type DiscoverMediaOption func(*discoverMediaOptions)

func DiscoverMediaOptionFrequency(d time.Duration) DiscoverMediaOption {
	return func(o *discoverMediaOptions) {
		o.frequency = d
	}
}

func DiscoverMediaOptionPeerTimeout(d time.Duration) DiscoverMediaOption {
	return func(o *discoverMediaOptions) {
		o.peerTimeout = d
	}
}

func DiscoverMediaOptionInfoTimeout(d time.Duration) DiscoverMediaOption {
	return func(o *discoverMediaOptions) {
		o.infoTimeout = d
	}
}

// newIdentifyOne builds the worker that locates, downloads, and identifies a single
// discovered torrent. peerTimeout bounds the DHT peer lookup; infoTimeout bounds the
// wait for torrent metadata/info once peers are being searched for.
func newIdentifyOne(db sqlx.Queryer, dhts *dht.Server, tclient *torrent.Client, ttstore storage.ClientImpl, peerTimeout, infoTimeout time.Duration) func(ctx context.Context, disc ddisc.Discovered) error {
	return func(ctx context.Context, disc ddisc.Discovered) (err error) {
		defer func() {
			if err == nil {
				return
			}

			log.Println("discover media failed", err)
			errorsx.Log(
				errorsx.Wrap(
					ddisc.DiscoveredCooldown(ctx, db, disc).Scan(&disc),
					"failed to mark discovered media for cooldown",
				),
			)
		}()

		if u := time.Until(disc.NextCheckAt); u > 0 {
			return errorsx.Wrapf(err, "premature attempt to identify media: %s - %v - %v", disc.ID, time.Until(disc.NextCheckAt), disc.NextCheckAt)
		}

		dctx, done := context.WithTimeout(ctx, infoTimeout)
		defer done()

		result, err := ddisc.IdentifyOne(ctx, dhts, tclient, ttstore, peerTimeout, infoTimeout, disc)
		if err != nil {
			return err
		}

		if err := ddisc.DiscoveredIndexed(dctx, db, result).Scan(&disc); err != nil {
			return errorsx.Wrapf(err, "unable to mark as indexed: %s", disc.ID)
		}

		log.Println("successfully downloaded media", disc.Description)

		return nil
	}
}

func DiscoverMedia(ctx context.Context, db sqlx.Queryer, dhts *dht.Server, tclient *torrent.Client, options ...DiscoverMediaOption) error {
	opts := langx.Clone(discoverMediaOptions{
		frequency:   envx.Duration(time.Hour, env.DDiscFrequency),
		peerTimeout: time.Minute,
		infoTimeout: 10 * time.Minute,
	}, options...)

	ttstore := storage.NewFile(userx.DefaultCacheDirectory("torrentddisc"))

	l := rate.NewLimiter(rate.Every(opts.frequency), 1)
	identifyone := asynccompute.New(
		newIdentifyOne(db, dhts, tclient, ttstore, opts.peerTimeout, opts.infoTimeout),
		asynccompute.Workers[ddisc.Discovered](envx.Uint[uint16](1, env.DDiscBackgroundWorkers)),
	)

	for err := l.Wait(ctx); err == nil; err = l.Wait(ctx) {
		q := ddisc.DiscoveredSearchBuilder().Where(
			squirrel.And{
				ddisc.DiscoveredQueryNextCheck(timex.NewRangeWithin(0)),
			},
		).OrderBy("attempts ASC, created_at DESC")

		s := sqlx.Scan(ddisc.DiscoveredSearch(ctx, db, q))

		for disc := range s.Iter() {
			if err := identifyone.Run(ctx, disc); err != nil {
				log.Println(err)
				continue
			}
		}

		if err := s.Err(); err == nil {
			errorsx.Log(errorsx.Wrap(os.RemoveAll(userx.DefaultCacheDirectory("torrentddisc")), "failed to reset torrent discovery cache"))
		} else if errorsx.Ignore(err, sql.ErrNoRows) == nil {
			time.Sleep(time.Minute)
		} else {
			errorsx.Log(errorsx.Wrap(s.Err(), "failed to retrieved indexable data"))
		}
	}

	return asynccompute.Shutdown(ctx, identifyone)
}
