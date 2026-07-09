package daemons

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/retroapi/ddiscapi"
	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/backoffx"
	"github.com/retrovibed/retrovibed/shallows/internal/contextx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
)

// searchPlugins is the narrow interface SearchQueueBackgroundRun needs from
// *retroapi/searchplugin.Registry, so tests can inject a fake without a real
// compiled wasm plugin.
type searchPlugins interface {
	Search(ctx context.Context, category, query string) iterx.Seq[*ddiscapi.Import]
}

func SearchQueueBackgroundRun(ctx context.Context, q sqlx.Queryer, plugins searchPlugins) error {
	// SearchQueueBackgroundRun drains ddisc_search_queue: for each pending
	// known-media-id, ask the loaded search plugins for candidates and persist
	// whatever they find, or push the entry's cooldown out if nothing turned up.
	// maxAge bounds how long a known-media-id stays queued for
	// search-plugin discovery before it's given up on and purged.
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
		seq := iterx.NotFound(plugins.Search(sctx, mimex.Category(known.Mimetype), known.Title))
		for r := range seq.Each(sctx) {
			m, err := metainfo.ParseMagnetURI(r.Magnet)
			if err != nil {
				errorsx.Log(err)
				continue
			}
			d := ddisc.NewDiscoveredFromKnown(
				int160.FromBytes(m.InfoHash.Bytes()),
				known,
				ddisc.DiscoveredOptionMimetype(langx.FirstNonZero(r.Mimetype, known.Mimetype)),
			)
			errorsx.Log(ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))
		}
		cancel()

		// TODO: I dont think we care *what* error occurs here.
		// we should cool down on any error.
		if err := seq.Err(); errors.Is(err, iterx.ErrNotFound) {
			log.Println("search queue: no candidates found", entry.KnownMediaID, known.Title)
			errorsx.Log(ddisc.SearchQueueCooldown(ctx, q, entry.KnownMediaID).Scan(&entry))
			continue
		} else if err != nil {
			log.Println("search queue failed:", entry.KnownMediaID, known.Title, err)
			errorsx.Log(ddisc.SearchQueueCooldown(ctx, q, entry.KnownMediaID).Scan(&entry))
			continue
		}

		errorsx.Log(ddisc.SearchQueueResolve(ctx, q, entry.KnownMediaID).Scan(&entry))
	}

	return s.Err()
}

// SearchQueueBackground drains the queue, then polls for new entries on an
// exponential backoff that maxes out at an hour.
func SearchQueueBackground(ctx context.Context, q sqlx.Queryer, plugins searchPlugins) error {
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
			return SearchQueueBackgroundRun(ctx, q, plugins)
		}))
	})

	return nil
}
