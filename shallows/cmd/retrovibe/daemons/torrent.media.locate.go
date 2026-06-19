package daemons

import (
	"context"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/duckdbx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/lucenex"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

func LocateTorrentMedia(ctx context.Context, db sqlx.Queryer, c *torrent.Client) error {
	q := library.KnownSearchBuilder().InnerJoin("library_locate ON library_locate.known_media_id = library_known_media.uid").Where(squirrel.And{
		library.LocateQueryPending(),
		library.KnownQueryExplicit(false),
	})

	download := func(r library.Known) (err error) {
		qq := ddisc.DiscoveredSearchBuilder().Where(squirrel.And{
			lucenex.Query(duckdbx.NewLucene(), r.Title, lucenex.WithDefaultField("title")),
		})
		ss := sqlx.Scan(ddisc.DiscoveredSearch(ctx, db, qq))

		for d := range ss.Iter() {
			var (
				l library.Locate
			)
			metadata, err := torrent.New(metainfo.Hash(d.Infohash))
			if err != nil {
				return errorsx.Wrapf(err, "unable to create torrent from infohash %s", d.ID)
			}

			info, err := c.Info(
				ctx,
				metadata,
				torrent.TuneAnnounceUntilComplete,
			)

			if err != nil {
				return errorsx.Wrapf(err, "unable to retrieve metadata from infohash %s", d.ID)
			}

			lmd := tracking.NewMetadata(
				new(int160.FromBytes(d.Infohash)),
				tracking.MetadataOptionFromInfo(info),
				tracking.MetadataOptionAutoDescription,
				tracking.MetadataOptionAutoHidden,
			)

			if err = tracking.MetadataInsertWithDefaults(ctx, db, lmd).Scan(&lmd); err != nil {
				return errorsx.Wrapf(err, "unable to record metadata for download from infohash %s", d.ID)
			}

			if err = tracking.MetadataAutoDownloadByID(ctx, db, lmd.ID).Scan(&lmd); err != nil {
				return errorsx.Wrapf(err, "unable to mark torrent for download from infohash %s", d.ID)
			}

			log.Println("marked for download", lmd.ID, lmd.Description)
			if err = library.LocateMarkTorrent(ctx, db, r.ID, lmd.ID).Scan(&l); err != nil {
				return errorsx.Wrapf(err, "unable to mark locate torrent %s", r.ID)
			}
			return nil
		}

		return ss.Err()
	}

	s := sqlx.Scan(library.KnownSearch(ctx, db, q))

	for r := range s.Iter() {
		log.Println("locating initiated", r.ID, r.Title)
		if err := download(r); err != nil {
			errorsx.Log(err)
			continue
		}
		log.Println("locating completed", r.ID, r.Title)
	}

	return s.Err()
}
