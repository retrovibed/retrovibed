package ddisc

import (
	"context"
	"iter"

	"github.com/retrovibed/retrovibed/internal/iterx"
	"github.com/retrovibed/retrovibed/internal/sqlx"
)

// search for known media
func FindMedia(q sqlx.Queryer, kid string) iterx.Seq[Discovered] {
	return &knownMedia{
		ID: kid,
		q:  q,
	}
}

type knownMedia struct {
	q   sqlx.Queryer
	ID  string
	err error
}

func (t *knownMedia) Each(ctx context.Context) iter.Seq[Discovered] {
	return func(yield func(Discovered) bool) {
		s := sqlx.Scan(DiscoveredByKnownID(ctx, t.q, t.ID))

		for v := range s.Iter() {
			if !yield(v) {
				return
			}
		}

		t.err = s.Err()
	}
}

func (t *knownMedia) Err() error {
	return t.err
}
