package daemons

import (
	"context"
	"errors"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/ddisc/ddisctorrent"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

// LocateMedia drains pending library_locate rows: for each, it runs (if dhts
// and partitions are non-nil) the DHT partition-swarm peer query, then (if
// plugins is non-nil) external search plugins to (re)populate ddisc_media,
// then ranks every candidate already in ddisc_media for that known-media-id
// with policy and downloads the best one.
func LocateMedia(ctx context.Context, db sqlx.Queryer, c *torrent.Client, disc *DiscoverySettings, dhts *dht.Server, partitions *ddisc.Partition, plugins searchPlugins, policy ddisc.Policy) error {
	if !disc.LocateP2P {
		return nil
	}

	q := library.KnownSearchBuilder().InnerJoin("library_locate ON library_locate.known_media_id = library_known_media.uid").Where(squirrel.And{
		library.LocateQueryPending(),
		library.KnownQueryExplicit(false),
	})

	download := func(r library.Known) (err error) {
		strategies := []ddisc.DiscoverStrategy{}
		if dhts != nil && partitions != nil {
			strategies = append(strategies, ddisctorrent.NewPartitionStrategy(dhts, partitions))
		}
		if plugins != nil {
			strategies = append(strategies, ddisc.PluginStrategy(db, plugins))
		}

		seq := ddisc.Discover(ctx, ddisc.DiscoverRequestFromKnown(r), strategies...)
		for range seq.Each(ctx) {
			// draining for persistence side effects only; selection happens below
		}
		if err := seq.Err(); err != nil {
			return err
		}

		d, err := ddisc.RankAndSelect(ctx, db, policy, r.UID)
		if errors.Is(err, ddisc.ErrNoCandidate) {
			return nil
		} else if err != nil {
			return err
		}

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
