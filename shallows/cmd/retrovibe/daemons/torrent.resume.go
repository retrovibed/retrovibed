package daemons

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/james-lawrence/torrent/storage"
	"github.com/james-lawrence/torrent/tracker"
	"github.com/retrovibed/retrovibed/retroapi/asynccompute"
	"github.com/retrovibed/retrovibed/retroapi/backoffx"
	"github.com/retrovibed/retrovibed/retroapi/env"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/internal/torrentx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

func ResumeDownloads(ctx context.Context, db sqlx.Queryer, rootstore fsx.Virtual, mc library.QueryCleaner, tclient *torrent.Client, tstore storage.ClientImpl, pub *asyncx.Wakeup) {
	q := tracking.MetadataSearchBuilder().Where(
		squirrel.And{
			tracking.MetadataQueryInitiated(),
			tracking.MetadataQueryIncomplete(),
			tracking.MetadataQueryNotPaused(),
		},
	)

	iter := sqlx.Scan(tracking.MetadataSearch(ctx, db, q))

	for md := range iter.Iter() {
		id := metainfo.Hash(md.Infohash)
		infopath := rootstore.Path(env.TorrentDirName, fmt.Sprintf("%s.torrent", id))
		log.Println("resuming", md.ID, md.Description, md.Private, infopath)

		metadata, err := torrent.New(
			metainfo.Hash(md.Infohash),
			torrent.OptionStorage(tstore),
			torrentx.OptionTracker(md.Tracker),
			torrentx.OptionInfoFromFile(infopath),
			torrent.OptionPublicTrackers(md.Private, tracking.PublicTrackers()...),
			torrent.OptionDisplayName(md.Description),
		)
		if err != nil {
			log.Println(errorsx.Wrapf(err, "unable to create metadata from %s - %s - %s", md.ID, md.Description, infopath))
			return
		}

		t, added, err := tclient.Start(metadata)
		if err != nil {
			log.Println(errorsx.Wrapf(err, "unable to start download %s - %s - %s", md.ID, md.Description, infopath))
			continue
		}

		if !added {
			log.Printf("torrent already running %s - %s\n", md.ID, infopath)
			continue
		}

		go func(infopath string, md tracking.Metadata, dl torrent.Torrent) {
			var (
				mhash = md5.New()
			)

			errorsx.Log(errorsx.Wrap(tracking.DownloadInto(ctx, db, rootstore, mc, &md, dl, mhash, pub), "resume failed"))
		}(infopath, md, t)

		log.Println("resumed", md.ID, hex.EncodeToString(md.Infohash), md.Description)
	}

	errorsx.Log(errorsx.Wrap(iter.Err(), "failed to resume all downloads"))
}

func VerifyTorrents(ctx context.Context, db sqlx.Queryer, rootstore fsx.Virtual, mc library.QueryCleaner, tclient *torrent.Client, tstore storage.ClientImpl, pub *asyncx.Wakeup) {
	q := tracking.MetadataSearchBuilder().Where(
		squirrel.And{
			tracking.MetadataQueryNeedsVerification(),
		},
	)
	log.Println("torrent verification initiated")
	defer log.Println("torrent verification completed")

	iter := sqlx.Scan(tracking.MetadataSearch(ctx, db, q))

	for md := range iter.Iter() {
		id := metainfo.Hash(md.Infohash)
		infopath := rootstore.Path(env.TorrentDirName, fmt.Sprintf("%s.torrent", id))
		log.Println("verifying", md.ID, md.Description, md.Private, infopath)

		metadata, err := torrent.New(
			metainfo.Hash(md.Infohash),
			torrent.OptionStorage(tstore),
			torrentx.OptionTracker(md.Tracker),
			torrentx.OptionInfoFromFile(infopath),
			torrent.OptionPublicTrackers(md.Private, tracking.PublicTrackers()...),
			torrent.OptionDisplayName(md.Description),
		)
		if err != nil {
			log.Println(errorsx.Wrapf(err, "unable to create metadata from %s - %s", md.ID, infopath))
			return
		}

		log.Println("verification initiated", md.ID, int160.FromBytes(md.Infohash), md.Description)
		t, _, err := tclient.Start(metadata, torrent.TuneVerifySample(32))
		if err != nil {
			log.Println(errorsx.Wrapf(err, "unable to start download %s - %s", md.ID, infopath))
			continue
		}

		if t.Info() == nil {
			log.Println("verification ignored", md.ID, int160.FromBytes(md.Infohash), md.Description)
			// ignoring for now
			continue
		}

		if err = sqlx.Discard(sqlx.Scan(library.MetadataDeleteByTorrentID(ctx, db, md.ID))); err != nil {
			log.Println(errorsx.Wrapf(err, "unable reset media for torrent %s - %s", md.ID, infopath))
			continue
		}

		if err = tracking.DownloadInto(ctx, db, rootstore, mc, &md, t, md5.New(), pub); err != nil {
			log.Println(errorsx.Wrapf(err, "unable reimport torrent %s - %s", md.ID, infopath))
			continue
		}

		log.Println("verification completed", md.ID, md.Description, t.BytesCompleted(), "/", t.Info().TotalLength())

		stats := t.Stats()
		if err := tracking.MetadataVerifiedByID(ctx, db, md.ID, 0, (&stats.BytesValidated).Uint64(), uint64(t.BytesCompleted())).Scan(&md); err != nil {
			log.Println(errorsx.Wrapf(err, "unable to update bytes completed during verification %s - %s", md.ID, infopath))
			continue
		}
	}

	errorsx.Log(errorsx.Wrap(iter.Err(), "failed to resume all downloads"))
}

func AnnounceSeeded(ctx context.Context, q sqlx.Queryer, dhts *dht.Server, rootstore fsx.Virtual, tclient *torrent.Client, tstore storage.ClientImpl) {
	const defaultDelay = 40 * time.Minute

	defer log.Println("announce seeded completed")

	announcer := torrentx.AnnouncerFromClient(tclient)

	announceDHT := func(cctx context.Context, md tracking.Metadata) time.Duration {
		if md.Private {
			return 0
		}

		id := int160.FromBytes(md.Infohash)
		log.Println("announcing dht", md.ID, md.Description, id)
		errorsx.Log(errorsx.Wrap(torrent.DHTAnnounceOnce(cctx, dhts, id), "announce dht failed"))

		return defaultDelay
	}

	announceTracker := func(cctx context.Context, md tracking.Metadata) time.Duration {
		if stringsx.Blank(md.Tracker) {
			return 0
		}

		req := tracker.NewAccounceRequest(
			dhts.ID(dhts.DynamicAddrPort()),
			dhts.DynamicAddrPort().Port(),
			int160.FromBytes(md.Infohash),
			tracker.AnnounceOptionKey,
			tracker.AnnounceOptionDownloaded(int64(md.Downloaded)),
			tracker.AnnounceOptionUploaded(int64(md.Uploaded)),
			tracker.AnnounceOptionSeeding,
		)

		log.Println("announcing tracker", md.ID, int160.FromBytes(md.Infohash).String())
		announced, err := announcer.ForTracker(md.Tracker).Do(cctx, req)
		if errors.Is(err, tracker.ErrMissingInfoHash) {
			errorsx.Log(errorsx.Wrapf(tracking.MetadataDisableAnnounced(cctx, q, md.ID).Scan(&md), "unable to disable announcements for torrent: %s", md.ID))
			return defaultDelay
		} else if err != nil {
			log.Println("failed to announce seeded torrent", md.ID, int160.FromBytes(md.Infohash).String(), err)
			return defaultDelay
		}

		return timex.DurationMax(time.Duration(announced.Interval)*time.Second, 10*time.Minute)
	}

	pool := asynccompute.New(func(ctx context.Context, i tracking.Metadata) error {
		nextts := time.Now().Add(timex.DurationMax(announceDHT(ctx, i), announceTracker(ctx, i))).Add(backoffx.Random(5 * time.Minute))
		if err := tracking.MetadataAnnounced(ctx, q, i.ID, nextts).Scan(&i); err != nil {
			return errorsx.Wrap(err, "failed to record announcement")
		}

		log.Println("announced", i.ID, int160.FromBytesOrZero(i.Infohash), i.Description, "next", i.NextAnnounceAt)
		return nil
	})
	defer func() {
		errorsx.Log(asynccompute.Shutdown(ctx, pool))
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			log.Println("running announcements")
		}

		query := tracking.MetadataSearchBuilder().Where(
			squirrel.And{
				tracking.MetadataQueryNotPaused(),
				tracking.MetadataQuerySeeding(),
				tracking.MetadataQueryNeedsAnnounce(),
			},
		)

		s := sqlx.Scan(tracking.MetadataSearch(ctx, q, query))
		for i := range s.Iter() {
			if err := pool.Run(ctx, i); err != nil {
				log.Println(err)
				continue
			}
		}

		if err := s.Err(); err != nil {
			errorsx.Log(errorsx.Wrap(err, "failed to retrieve torrents to announce"))
			continue
		}

		// we finished announcing, determine when the next announcement is due.
		{
			var (
				err    error
				nextts time.Time
				delay  = backoffx.Random(16 * time.Second)
			)

			query, args, err := tracking.MetadataSearchBuilder().RemoveColumns().Columns("next_announce_at").Where(
				squirrel.And{
					tracking.MetadataQueryNotPaused(),
					tracking.MetadataQuerySeeding(),
					tracking.MetadataQueryAnnounceable(),
				},
			).OrderBy("next_announce_at ASC").Limit(1).ToSql()
			if err != nil {
				panic(errorsx.Wrap(err, "failed to generate announce query"))
			}

			// log.Println("derp", errorsx.Zero(sqlx.Timestamp(ctx, q, "UPDATE torrents_metadata SET next_announce_at = NOW() + INTERVAL '5s' WHERE id = '8ff31fd8-3950-1ee9-0bef-186e20d8e812' RETURNING next_announce_at")))

			if nextts, err = sqlx.Timestamp(ctx, q, query, args...); err != nil {
				log.Printf("unable to determine next timestamp - %v", err)
				nextts = time.Now().Add(10 * time.Minute)
			}

			nextts = nextts.Add(delay)

			log.Printf("announce - next will be %v %v\n", time.Until(nextts), nextts)
			time.Sleep(time.Until(nextts))
		}
	}
}
