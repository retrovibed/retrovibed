package daemons

import (
	"context"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

// PublishDiscoveredMedia scans library media that has been matched to a known media entity and
// belongs to a torrent, and publishes a ddisc.Discovered record for it so it can be
// announced/synced to peers, even when this node downloaded the content on its own behalf
// (i.e. outside of DHT-based discovery). Torrents flagged private are still published so
// this node can use them locally, but are marked private so they're never handed to peers
// (see ddisc.DiscoveredOptionPrivate and the private-excluding sync/search queries).
func PublishDiscoveredMedia(ctx context.Context, db sqlx.Queryer) error {
	q := library.MetadataSearchBuilder().Where(
		squirrel.And{
			library.MetadataQueryHasKnownMedia(true),
			library.MetadataQueryHasTorrent(true),
			library.MetadataQueryNotTombstoned(),
			library.MetadataQueryDirectory(false),
			library.MetadataQueryUpdatedBetween(timex.NewRangeDuration(48 * time.Hour)),
		},
	)

	iter := sqlx.Scan(library.MetadataSearch(ctx, db, q))

	log.Println("publishing discovered media initiated")
	defer log.Println("publishing discovered media completed")

	for lmd := range iter.Iter() {
		if err := publishDiscoveredMediaOne(ctx, db, lmd); err != nil {
			log.Println("unable to publish discovered media", lmd.ID, err)
		}
	}

	return errorsx.Wrap(iter.Err(), "failed to scan library metadata for ddisc publishing")
}

func publishDiscoveredMediaOne(ctx context.Context, db sqlx.Queryer, lmd library.Metadata) error {
	var tmd tracking.Metadata
	if err := tracking.MetadataFindByID(ctx, db, lmd.TorrentID).Scan(&tmd); err != nil {
		return errorsx.Wrap(err, "unable to find torrent metadata")
	}

	var known library.Known
	if err := library.KnownFindByID(ctx, db, lmd.KnownMediaID).Scan(&known); err != nil {
		return errorsx.Wrap(err, "unable to find known media")
	}

	id := int160.FromBytes(tmd.Infohash)

	// skip files already published, keyed deterministically on (infohash, known_media_id).
	candidate := ddisc.NewDiscoveredFromKnown(
		id, known,
		ddisc.DiscoveredOptionPrivate(tmd.Private),
		ddisc.DiscoveredOptionAutoMagnet,
		ddisc.DiscoveredOptionAcquisitionState(ddisc.AcquisitionStateCompleted),
	)
	var existing ddisc.Discovered
	if err := ddisc.DiscoveredFindByID(ctx, db, candidate.ID).Scan(&existing); err == nil {
		return nil
	}

	if err := ddisc.DiscoveredInsertWithDefaults(ctx, db, candidate).Scan(&candidate); err != nil {
		return errorsx.Wrap(err, "unable to insert discovered record")
	}

	log.Println("published discovered media", candidate.ID, tmd.ID, "->", known.UID, known.Title)

	return nil
}
