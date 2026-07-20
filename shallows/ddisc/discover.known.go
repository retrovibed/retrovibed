package ddisc

import (
	"context"
	"iter"

	"github.com/Masterminds/squirrel"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/retrovibed/retrovibed/shallows/internal/slicesx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/squirrelx"
	"github.com/retrovibed/retrovibed/shallows/library"
)

// knownSimilarityCutoff is the minimum title-similarity score (see
// library.KnownQuerySimilarity) a library_known_media row must clear to be
// considered a catalog match.
const knownSimilarityCutoff = 0.6

// knownStrategyLimit caps how many catalog matches a single Discover call
// can surface.
const knownStrategyLimit = 10

// KnownStrategy searches the library_known_media catalog by fuzzy title
// match. Unlike LocalStrategy (which requires an exact known-media-id to
// already be resolved), it works from req.Title alone, so it can surface
// media we already know about even before any download source has been
// found for it. Every candidate it yields is catalog-only - see
// DiscoveredOptionCatalogOnly.
func KnownStrategy(q sqlx.Queryer) DiscoverStrategy {
	return knownStrategy{q: q}
}

type knownStrategy struct {
	q sqlx.Queryer
}

func (t knownStrategy) Discover(ctx context.Context, req DiscoverRequest) iterx.Seq[Discovered] {
	return &knownSeq{q: t.q, req: req}
}

type knownSeq struct {
	q   sqlx.Queryer
	req DiscoverRequest
	err error
}

func (t *knownSeq) Each(ctx context.Context) iter.Seq[Discovered] {
	return func(yield func(Discovered) bool) {
		if t.req.Title == "" {
			return
		}

		// req.Mimetypes carries discovery-specific mimetypes (see Category);
		// library_known_media.mimetype stores the coarse category, so
		// Generalize each one back before filtering.
		mimetypes := slicesx.MapTransform(Generalize, t.req.Mimetypes...)

		qq := library.KnownSearchBuilder().Where(squirrel.And{
			library.KnownQueryExplicit(t.req.Adult),
			squirrelx.In("library_known_media.mimetype", mimetypes...),
			library.KnownQuerySimilarity(t.req.Title, knownSimilarityCutoff),
		}).Limit(knownStrategyLimit)

		s := sqlx.Scan(library.KnownSearch(ctx, t.q, qq))
		for known := range s.Iter() {
			d := NewDiscoveredFromKnown(
				int160.FromHashedBytes([]byte(known.UID)),
				known,
				DiscoveredOptionCatalogOnly,
				DiscoveredOptionURI(catalogURI(known.UID)),
			)

			if !yield(d) {
				return
			}
		}

		t.err = s.Err()
	}
}

func (t *knownSeq) Err() error {
	return t.err
}
