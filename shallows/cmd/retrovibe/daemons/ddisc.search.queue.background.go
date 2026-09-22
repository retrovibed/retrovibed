package daemons

import (
	"context"
	"log"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/backoffx"
	"github.com/retrovibed/retrovibed/retroapi/searchplugin"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/contextx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/library"
)

func SearchQueueBackgroundRun(ctx context.Context, q sqlx.Queryer, plugins searchplugin.T, peertube ddisc.DiscoverStrategy, mc library.QueryCleaner) error {
	// SearchQueueBackgroundRun drains ddisc_search_queue: for each pending
	// known-media-id, ask the external search strategies (wasm plugins,
	// PeerTube/SepiaSearch) for candidates and persist whatever they find, or
	// push the entry's cooldown out if nothing turned up. Candidates are
	// persisted as-is (placeholder infohash, Private defaulted true - see
	// ddisc.NewDiscoveredFromImport) with no network fetch of the candidate's
	// .torrent: resolving the real infohash and BEP 27 privacy only happens
	// once a candidate is actually selected for download (see
	// ddisc.DownloadDiscovered), not for every candidate the plugins/peertube
	// turn up on every drain. maxAge bounds how long a known-media-id stays
	// queued for external discovery before it's given up on and purged.
	const maxAge = 30 * 24 * time.Hour
	errorsx.Log(sqlx.Discard(sqlx.Scan(ddisc.SearchQueuePurge(ctx, q, maxAge))))

	s := sqlx.Scan(ddisc.SearchQueuePending(ctx, q))
	for entry := range s.Iter() {
		var known library.Known
		if err := library.KnownFindByID(ctx, q, entry.KnownMediaID).Scan(&known); err != nil {
			errorsx.Log(err)
			continue
		}

		sctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		req := ddisc.DiscoverRequestFromKnown(known)
		options := []ddisc.DiscoverOption{
			ddisc.DiscoverOptionFilter(ddisc.NewTitleFilter(q, req).Match),
			ddisc.DiscoverOptionDetectMedia(ddisc.KnownMediaDetector(q, mc)),
			ddisc.DiscoverOptionKnownMediaDynamic(ddisc.KnownMediaDynamic(q)),
		}
		seq := ddisc.Discover(sctx, ddisc.DefaultPolicy(), req, options, ddisc.ExternalStrategies(plugins, peertube)...)

		found := false
		for d := range seq.Each(sctx) {
			// a candidate the policy has already rejected (cam, ts, etc.)
			// isn't worth persisting - see ddisc.Policy.Rank.
			if stringsx.Present(d.PolicyRejection) {
				continue
			}

			if err := ddisc.DiscoveredInsertWithDefaults(sctx, q, d).Scan(&d); err != nil {
				errorsx.Log(errorsx.Wrap(err, "unable to persist discovered candidate"))
				continue
			}

			found = true
		}
		err := seq.Err()
		cancel()

		if err != nil {
			log.Println("search queue failed:", entry.KnownMediaID, known.Title, err)
		} else {
			log.Println("search queue: no candidates found", entry.KnownMediaID, known.Title)
		}

		// we don't care *what* error occurs here (if any) — cool down on
		// any failure to find a candidate, same as a clean not-found.
		if !found || err != nil {
			errorsx.Log(ddisc.SearchQueueCooldown(ctx, q, entry.KnownMediaID).Scan(&entry))
			continue
		}

		errorsx.Log(ddisc.SearchQueueResolve(ctx, q, entry.KnownMediaID).Scan(&entry))
	}

	return s.Err()
}

// SearchQueueBackground drains the queue, then polls for new entries on an
// exponential backoff that maxes out at an hour.
func SearchQueueBackground(ctx context.Context, q sqlx.Queryer, plugins searchplugin.T, peertube ddisc.DiscoverStrategy, mc library.QueryCleaner) error {
	wakeup := asyncx.NewWakeup(ctx)
	defer wakeup.Broadcast() // kick off an initial drain
	s := backoffx.New(
		backoffx.Exponential(time.Second),
		backoffx.Maximum(time.Hour),
		backoffx.Jitter(0.1),
	)

	go asyncx.Periodic(ctx, wakeup, s, "ddisc search queue drain")
	contextx.Run(ctx, func() {
		errorsx.Log(asyncx.Run(ctx, wakeup, func(ctx context.Context) error {
			return SearchQueueBackgroundRun(ctx, q, plugins, peertube, mc)
		}))
	})

	return nil
}
