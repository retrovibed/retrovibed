package ddisc

import (
	"context"
	"iter"

	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
)

// LocalStrategy searches the local ddisc_media table for an exact
// known_media_id match.
func LocalStrategy(q sqlx.Queryer) DiscoverStrategy {
	return localStrategy{q: q}
}

type localStrategy struct {
	q sqlx.Queryer
}

func (t localStrategy) Discover(ctx context.Context, req DiscoverRequest) iterx.Seq[Discovered] {
	return &localSeq{q: t.q, kid: req.KnownMediaID}
}

type localSeq struct {
	q   sqlx.Queryer
	kid string
	err error
}

func (t *localSeq) Each(ctx context.Context) iter.Seq[Discovered] {
	return func(yield func(Discovered) bool) {
		qq := DiscoveredSearchBuilder().Where(DiscoveredQueryKnownMediaID(t.kid))
		s := sqlx.Scan(DiscoveredSearch(ctx, t.q, qq))

		for d := range s.Iter() {
			if !yield(d) {
				return
			}
		}

		t.err = s.Err()
	}
}

func (t *localSeq) Err() error {
	return t.err
}
