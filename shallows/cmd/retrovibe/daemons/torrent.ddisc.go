package daemons

import (
	"context"
	"database/sql"
	"encoding/hex"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/james-lawrence/torrent/storage"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/asynccompute"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/envx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/torrentx"
	"github.com/retrovibed/retrovibed/shallows/internal/userx"
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

		available := func(ctx context.Context, disc ddisc.Discovered) ([]dht.Peer, error) {
			dctx, done := context.WithTimeout(ctx, peerTimeout)
			defer done()
			return torrentx.Peers(dctx, dhts, int160.FromBytes(disc.Infohash))
		}

		peers, err := available(ctx, disc)
		if err != nil {
			return errorsx.Wrapf(err, "no peers available %s", disc.ID)
		}
		log.Println("located", len(peers), "initial peers")

		dctx, done := context.WithTimeout(ctx, infoTimeout)
		defer done()

		log.Println("identify initiated", disc.ID, hex.EncodeToString(disc.Infohash), disc.NextCheckAt)
		defer log.Println("identify completed", disc.ID, hex.EncodeToString(disc.Infohash), disc.NextCheckAt)

		if u := time.Until(disc.NextCheckAt); u > 0 {
			return errorsx.Wrapf(err, "premature attempt to identify media: %s - %v - %v", disc.ID, time.Until(disc.NextCheckAt), disc.NextCheckAt)
		}

		id := int160.FromBytes(disc.Infohash)
		metadata, err := torrent.New(metainfo.Hash(disc.Infohash), torrent.OptionStorage(ttstore))
		if err != nil {
			return errorsx.Wrapf(err, "unable to create metadata from infohash %s", disc.ID)
		}
		defer tclient.Stop(metadata)

		info, err := tclient.Info(
			dctx,
			metadata,
			torrent.TuneAnnounceUntilComplete,
			torrent.TuneNewConns,
		)
		if err != nil {
			return errorsx.Wrapf(err, "unable to initialize torrent for %s", disc.ID)
		}

		torrentx.FilePrint(
			info,
			torrentx.FileFirst(info, func(fi metainfo.FileInfo) bool { return strings.HasSuffix(fi.DisplayPath(info), ".nfo") }),
			torrentx.FileLargest(info),
		)

		_, off, length := torrentx.FileLargestRange(info)
		metadata, err = torrent.NewFromInfo(info, torrent.OptionStorage(ttstore))
		if err != nil {
			return errorsx.Wrapf(err, "unable to create metadata from info %s", id)
		}

		tt, _, err := tclient.Start(
			metadata,
			torrent.TuneResetBitmaps,
			torrent.TuneVerifyRange(off, min(length, 32*bytesx.KiB)),
			torrent.TuneAnnounceUntilComplete,
			torrent.TuneNewConns,
		)
		if err != nil {
			return errorsx.Wrapf(err, "unable to initialize torrent for %s", disc.ID)
		}

		log.Println("attempting to download", id, bytesx.Unit(length))
		r := torrent.DownloadRange(dctx, tt, off, min(length, 32*bytesx.KiB))
		go func() {
			torrentx.DownloadProgress(dctx, tt)
			errorsx.Log(errorsx.Wrapf(r.Close(), "%s failed to close reader", id))
			errorsx.Log(errorsx.Wrapf(tclient.Stop(metadata), "%s failed to stop torrent", id))
		}()

		mime, err := ddisc.Mimetype(r)
		if err != nil {
			return errorsx.Wrapf(err, "unable to initialize torrent for %s", disc.ID)
		}

		disc = langx.Clone(
			disc,
			ddisc.DiscoveredOptionMimetype(mime.String()),
			ddisc.DiscoveredOptionKnownMedia(uuid.Max.String()),
		)

		if err := ddisc.DiscoveredIndexed(dctx, db, disc).Scan(&disc); err != nil {
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
				ddisc.DiscoveredQueryNeedsCheck(),
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
