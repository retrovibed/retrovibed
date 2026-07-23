package ddisc

import (
	"context"
	"iter"

	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/retrovibed/retrovibed/retroapi/uuidx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/lucenex"
	"github.com/retrovibed/retrovibed/shallows/internal/slicesx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/squirrelx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
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
		if t.req.Query == "" {
			return
		}

		// req.Mimetypes carries discovery-specific mimetypes (see Category);
		// library_known_media.mimetype stores the coarse category, so
		// Generalize each one back before filtering.
		mimetypes := slicesx.MapTransform(Generalize, t.req.Mimetypes...)

		qq := library.KnownSearchBuilder().Where(squirrel.And{
			library.KnownQueryExplicit(t.req.Adult),
			squirrelx.In("library_known_media.mimetype", mimetypes...),
			library.KnownQuerySimilarity(t.req.Query, knownSimilarityCutoff),
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

// KnownMediaDetector builds a DiscoverOptionDetectMedia transform: for every
// Discovered candidate that doesn't already carry a known-media-id, it
// cleans the candidate's own Title through mc and lucenex.Clean, strips any
// trailing release/episode tokens via library.ParseReleaseEpisode, and looks
// up the remaining title via library.DetectKnownMedia - stamping a match
// onto the candidate before it's yielded. Unlike KnownStrategy (which
// matches the raw, uncleaned request query and yields its own catalog-only
// candidates), this enriches candidates already produced by other
// strategies (e.g. plugin/peertube hits) that don't know their own catalog
// match.
func KnownMediaDetector(q sqlx.Queryer, mc library.QueryCleaner) func(iterx.Seq[Discovered]) iterx.Seq[Discovered] {
	return func(s iterx.Seq[Discovered]) iterx.Seq[Discovered] {
		return &knownMediaDetectSeq{inner: s, q: q, mc: mc}
	}
}

type knownMediaDetectSeq struct {
	inner iterx.Seq[Discovered]
	q     sqlx.Queryer
	mc    library.QueryCleaner
	err   error
}

func (t *knownMediaDetectSeq) Each(ctx context.Context) iter.Seq[Discovered] {
	return func(yield func(Discovered) bool) {
		for d := range t.inner.Each(ctx) {
			if stringsx.Blank(d.Title) {
				if !yield(d) {
					return
				}
				continue
			}

			cleaned, err := t.mc.Clean(ctx, d.Title)
			if err != nil {
				t.err = errorsx.Wrapf(err, "unable to clean title: %s", d.Title)
				return
			}

			cleaned = lucenex.Clean(cleaned)
			title, _, _ := library.ParseReleaseEpisode(cleaned)

			known, err := library.DetectKnownMedia(ctx, t.q, d.Mimetype, title)
			if err != nil {
				t.err = errorsx.Wrap(err, "unable to detect known media")
				return
			}

			if !uuidx.IsMinMax(uuid.FromStringOrNil(known.UID)) {
				d.KnownMediaID = known.UID
			}

			if !yield(d) {
				return
			}
		}

		t.err = langx.FirstNonNil(t.err, t.inner.Err())
	}
}

func (t *knownMediaDetectSeq) Err() error {
	return t.err
}
