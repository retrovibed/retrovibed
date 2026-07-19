package daemons

import (
	"context"
	"log"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/searchplugin"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/backoffx"
	"github.com/retrovibed/retrovibed/shallows/internal/contextx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

func SearchQueueBackgroundRun(ctx context.Context, q sqlx.Queryer, importer tracking.URIImport, plugins searchplugin.T) error {
	// SearchQueueBackgroundRun drains ddisc_search_queue: for each pending
	// known-media-id, ask the loaded search plugins for candidates, resolve
	// each candidate's real infohash (without importing/downloading it - see
	// tracking.URIImport.Resolve) and persist whatever they find, or push the
	// entry's cooldown out if nothing turned up. maxAge bounds how long a
	// known-media-id stays queued for search-plugin discovery before it's
	// given up on and purged.
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
		seq := ddisc.Discover(sctx, ddisc.DefaultPolicy(), ddisc.DiscoverRequestFromKnown(known), ddisc.PluginStrategy(q, plugins))

		found := false
		for d := range seq.Each(sctx) {
			found = true

			if resolved, rerr := importer.Resolve(sctx, d.URI); rerr != nil {
				errorsx.Log(errorsx.Wrap(rerr, "unable to resolve discovered candidate"))
				continue
			} else {
				d.Infohash = resolved.Infohash
			}

			if err := ddisc.DiscoveredInsertWithDefaults(sctx, q, d).Scan(&d); err != nil {
				errorsx.Log(errorsx.Wrap(err, "unable to persist discovered candidate"))
				continue
			}
		}
		err := seq.Err()
		cancel()

		// we don't care *what* error occurs here (if any) — cool down on
		// any failure to find a candidate, same as a clean not-found.
		if !found || err != nil {
			if err != nil {
				log.Println("search queue failed:", entry.KnownMediaID, known.Title, err)
			} else {
				log.Println("search queue: no candidates found", entry.KnownMediaID, known.Title)
			}
			errorsx.Log(ddisc.SearchQueueCooldown(ctx, q, entry.KnownMediaID).Scan(&entry))
			continue
		}

		errorsx.Log(ddisc.SearchQueueResolve(ctx, q, entry.KnownMediaID).Scan(&entry))
	}

	return s.Err()
}

// SearchQueueBackground drains the queue, then polls for new entries on an
// exponential backoff that maxes out at an hour.
func SearchQueueBackground(ctx context.Context, q sqlx.Queryer, importer tracking.URIImport, plugins searchplugin.T) error {
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
			return SearchQueueBackgroundRun(ctx, q, importer, plugins)
		}))
	})

	return nil
}
