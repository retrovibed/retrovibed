package daemons

import (
	"context"
	"errors"
	"log"
	"runtime/trace"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/dht"
	"github.com/retrovibed/retrovibed/retroapi/backoffx"
	"github.com/retrovibed/retrovibed/retroapi/searchplugin"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/ddisc/ddisctorrent"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

// Locate runs (if dhts and partitions are non-nil, and loc has a resolved
// known-media-id) the DHT partition-swarm peer query, (if plugins is
// non-nil) external wasm search plugins, (if peertube is non-nil) the
// in-process PeerTube/SepiaSearch strategy, and (if loc has a resolved
// known-media-id) a scan of whatever's already recorded in ddisc_media -
// every pass, regardless of whether an earlier strategy already produced
// something - and picks the best-ranked candidate across all of them via
// ddisc.Select/Compare. Returns ddisc.ErrNoCandidate if nothing ranked yet -
// this is a normal "nothing found this pass" outcome, not a failure. Never
// downloads, and never persists a candidate itself - see
// DiscoveredDownload, which persists only the winner, once import has
// resolved its real infohash.
//
// The DHT partition strategy never yields synchronously - any peer response
// lands in ddisc_media asynchronously via the already-registered MethodMedia
// responder (see ddisctorrent.NewPartitionStrategy). The local ddisc_media
// scan is what notices those (and any previously-imported winner) on a
// later pass.
func Locate(ctx context.Context, db sqlx.Queryer, disc *DiscoverySettings, dhts *dht.Server, partitions *ddisc.Partition, plugins searchplugin.T, peertube ddisc.DiscoverStrategy, policy ddisc.Policy, mc library.QueryCleaner, loc ddisc.Locate) (ddisc.Discovered, error) {
	trace.Logf(ctx, "input", "known_media_id: %q query: %q", loc.KnownMediaID, loc.Query)

	strategies := []ddisc.DiscoverStrategy{}
	// the DHT partition strategy, the local fallback strategy, and the
	// partition's peer-side responder all key strictly off known_media_id
	// equality, which is meaningless for an unresolved (Nil) known_media_id:
	// every free-text locate would converge on the same partition/local scan
	// and match unrelated Nil-tagged data from completely different searches.
	if dhts != nil && partitions != nil && loc.KnownMediaID != uuid.Nil.String() {
		strategies = append(strategies, ddisctorrent.NewPartitionStrategy(dhts, partitions))
	}
	strategies = append(strategies, ddisc.SyncStrategies(db, plugins, peertube, loc.KnownMediaID)...)

	req := ddisc.DiscoverRequestFromLocate(loc)

	options := []ddisc.DiscoverOption{
		ddisc.DiscoverOptionFilter(ddisc.NewTitleFilter(db, req).Match),
		ddisc.DiscoverOptionDetectMedia(ddisc.KnownMediaDetector(db, mc)),
		ddisc.DiscoverOptionKnownMediaDynamic(ddisc.KnownMediaDynamic(db)),
	}
	seq := ddisc.Discover(ctx, policy, req, options, strategies...)
	best, err := ddisc.Select(func(yield func(ddisc.Discovered) bool) {
		for v := range seq.Each(ctx) {
			log.Println("located", v.Title, v.PolicyRank, v.PolicyRejection)
			trace.Logf(ctx, "candidate", "title: %q rank: %v rejection: %q", v.Title, v.PolicyRank, v.PolicyRejection)
			if !yield(v) {
				return
			}
		}
	})
	if serr := seq.Err(); serr != nil {
		trace.Logf(ctx, "error", "%v", serr)
		return ddisc.Discovered{}, serr
	}

	return best, err
}

// DiscoveredDownload resolves d's torrent metadata via importer, records it
// for download, marks it for auto-download, and stamps loc as located. Full
// info-dict resolution (if not already available from d.URI) is left to the
// normal download/resume machinery (see ResumeDownloads) once initiated_at
// is set - same as every other ingestion path (RSS, etc.). This is also the
// only ddisc_media write the whole locate ever makes: non-winning candidates
// are never persisted.
func DiscoveredDownload(ctx context.Context, db sqlx.Queryer, importer tracking.URIImport, loc ddisc.Locate, d ddisc.Discovered) (err error) {
	var l ddisc.Locate

	acquisition := ddisc.AcquisitionStateAvailable
	if loc.Autodownload {
		acquisition = ddisc.AcquisitionStateDownloading
	}

	_, lmd, err := ddisc.DownloadDiscovered(ctx, db, importer, d, acquisition)
	if err != nil {
		return err
	}

	if err = ddisc.LocateLocated(ctx, db, loc.ID, lmd.ID).Scan(&l); err != nil {
		return errorsx.Wrapf(err, "unable to mark locate torrent %s", loc.ID)
	}

	return nil
}

// LocateMedia drains pending ddisc_locate rows, locating and downloading
// the best candidate for each.
func LocateMedia(ctx context.Context, db sqlx.Queryer, importer tracking.URIImport, disc *DiscoverySettings, dhts *dht.Server, partitions *ddisc.Partition, plugins searchplugin.T, peertube ddisc.DiscoverStrategy, policy ddisc.Policy, mc library.QueryCleaner) error {
	var task *trace.Task
	ctx, task = trace.NewTask(ctx, "ddisc.locate_media")
	defer task.End()

	log.Println("locate media initiated")
	defer log.Println("locate media completed")

	// p2p consent is enforced client-side only (see console's ensureP2P);
	// this daemon no longer gates on disc.LocateP2P.
	locateCooldown := backoffx.New(
		backoffx.Multiple(24*time.Hour),
		backoffx.Maximum(7*24*time.Hour),
		backoffx.JitterRandom(15*time.Minute),
	)
	q := ddisc.LocateSearchBuilder().Where(ddisc.LocateQueryPending())
	s := sqlx.Scan(ddisc.LocateSearch(ctx, db, q))

	for loc := range s.Iter() {
		log.Println("locating initiated", loc.ID, loc.Query)
		trace.Logf(ctx, "entry", "known_media_id: %q query: %q", loc.KnownMediaID, loc.Query)

		nextCheckAt := time.Now().Add(locateCooldown.Backoff(int(loc.Attempts)))
		if err := ddisc.LocateCooldown(ctx, db, loc.ID, nextCheckAt).Scan(&loc); err != nil {
			errorsx.Log(err)
			continue
		}

		d, err := Locate(ctx, db, disc, dhts, partitions, plugins, peertube, policy, mc, loc)
		if errors.Is(err, ddisc.ErrNoCandidate) {
			trace.Logf(ctx, "outcome", "known_media_id: %q status: no_candidate", loc.KnownMediaID)
			continue
		} else if err != nil {
			errorsx.Log(err)
			trace.Logf(ctx, "outcome", "known_media_id: %q status: error err: %v", loc.KnownMediaID, err)
			continue
		}

		if err := DiscoveredDownload(ctx, db, importer, loc, d); err != nil {
			errorsx.Log(err)
			trace.Logf(ctx, "outcome", "known_media_id: %q status: download_error err: %v", loc.KnownMediaID, err)
			continue
		}

		if err := ddisc.LocateCompleted(ctx, db, loc.ID).Scan(&loc); err != nil {
			errorsx.Log(err)
			trace.Logf(ctx, "outcome", "known_media_id: %q status: complete_error err: %v", loc.KnownMediaID, err)
			continue
		}

		trace.Logf(ctx, "outcome", "known_media_id: %q status: completed", loc.KnownMediaID)
		log.Println("locating completed", loc.ID, loc.Query)
	}

	return s.Err()
}
