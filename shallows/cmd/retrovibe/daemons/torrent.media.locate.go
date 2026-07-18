package daemons

import (
	"context"
	"errors"
	"log"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/dht"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/ddisc/ddisctorrent"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

// Locate runs (if dhts and partitions are non-nil, and loc has a resolved
// known-media-id) the DHT partition-swarm peer query, then (if plugins is
// non-nil) external search plugins to (re)populate ddisc_media, then ranks
// every candidate already in ddisc_media for loc's known-media-id whose
// title matches loc.Query with policy and returns the best one. Returns
// ddisc.ErrNoCandidate if nothing ranks yet - this is a normal "nothing
// found this pass" outcome, not a failure. Never downloads.
func Locate(ctx context.Context, db sqlx.Queryer, disc *DiscoverySettings, dhts *dht.Server, partitions *ddisc.Partition, plugins searchPlugins, policy ddisc.Policy, loc ddisc.Locate) (ddisc.Discovered, error) {
	strategies := []ddisc.DiscoverStrategy{}
	// the DHT partition strategy and its peer-side responder both key
	// strictly off known_media_id equality, which is meaningless for an
	// unresolved (Nil) known_media_id: every free-text locate would converge
	// on the same partition and query peers' unrelated Nil-tagged data.
	if dhts != nil && partitions != nil && loc.KnownMediaID != uuid.Nil.String() {
		strategies = append(strategies, ddisctorrent.NewPartitionStrategy(dhts, partitions))
	}
	if plugins != nil {
		strategies = append(strategies, ddisc.PluginStrategy(db, plugins))
	}

	req := ddisc.DiscoverRequest{
		KnownMediaID: loc.KnownMediaID,
		Title:        loc.Query,
		Mimetypes:    ddisc.Category(loc.Mimetype),
		Adult:        loc.Adult,
	}

	seq := ddisc.Discover(ctx, db, policy, req, strategies...)
	for v := range seq.Each(ctx) {
		log.Println("located", v.Title, v.PolicyRank, v.PolicyRejection)
		// draining for persistence side effects only; selection happens below
	}
	if err := seq.Err(); err != nil {
		return ddisc.Discovered{}, err
	}

	return ddisc.RankAndSelect(ctx, db, policy, loc)
}

// DiscoveredDownload resolves d's torrent metadata via importer, records it
// for download, marks it for auto-download, and stamps loc as located. Full
// info-dict resolution (if not already available from d.URI) is left to the
// normal download/resume machinery (see ResumeDownloads) once initiated_at
// is set - same as every other ingestion path (RSS, etc.).
func DiscoveredDownload(ctx context.Context, db sqlx.Queryer, importer tracking.URIImport, loc ddisc.Locate, d ddisc.Discovered) (err error) {
	var (
		l ddisc.Locate
	)

	lmd, err := importer.Import(
		ctx,
		d.URI,
		tracking.MetadataOptionKnownMediaID(d.KnownMediaID),
		tracking.MetadataOptionAutoDescription,
		tracking.MetadataOptionAutoHidden,
	)
	if err != nil {
		return errorsx.Wrapf(err, "unable to import uri for download %s", d.ID)
	}

	if loc.Autodownload {
		if err = tracking.MetadataAutoDownloadByID(ctx, db, lmd.ID).Scan(&lmd); err != nil {
			return errorsx.Wrapf(err, "unable to mark torrent for download from uri %s", d.ID)
		}
		log.Println("marked for download", lmd.ID, lmd.Description)
	} else {
		var rec library.Recommendation
		if err = library.RecommendationInsertWithDefaults(ctx, db, ddisc.RecommendationFromDiscovered(d)).Scan(&rec); err != nil {
			return errorsx.Wrap(err, "unable to record recommendation for located media")
		}
		log.Println("recommendation created", rec.ID, rec.Title)
	}

	if err = ddisc.LocateLocated(ctx, db, loc.ID, lmd.ID).Scan(&l); err != nil {
		return errorsx.Wrapf(err, "unable to mark locate torrent %s", loc.ID)
	}

	return nil
}

// LocateMedia drains pending ddisc_locate rows, locating and downloading
// the best candidate for each.
func LocateMedia(ctx context.Context, db sqlx.Queryer, importer tracking.URIImport, disc *DiscoverySettings, dhts *dht.Server, partitions *ddisc.Partition, plugins searchPlugins, policy ddisc.Policy) error {
	log.Println("locate media initiated")
	defer log.Println("locate media completed")
	if !disc.LocateP2P {
		return nil
	}

	q := ddisc.LocateSearchBuilder().Where(ddisc.LocateQueryPending())
	s := sqlx.Scan(ddisc.LocateSearch(ctx, db, q))

	for loc := range s.Iter() {
		log.Println("locating initiated", loc.ID, loc.Query)

		d, err := Locate(ctx, db, disc, dhts, partitions, plugins, policy, loc)
		if errors.Is(err, ddisc.ErrNoCandidate) {
			continue
		} else if err != nil {
			errorsx.Log(err)
			continue
		}

		if err := DiscoveredDownload(ctx, db, importer, loc, d); err != nil {
			errorsx.Log(err)
			continue
		}

		if err := ddisc.LocateCompleted(ctx, db, loc.ID).Scan(&loc); err != nil {
			errorsx.Log(err)
			continue
		}

		log.Println("locating completed", loc.ID, loc.Query)
	}

	return s.Err()
}
