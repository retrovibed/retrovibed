package ddisc

import (
	"context"
	"iter"

	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
)

// send all discovered media for the given partitions since the syncoffset.
func SyncDiscovered(q sqlx.Queryer, filter rendezvousf, syncoffset string) iterx.Seq[Discovered] {
	return &syncDiscovered{
		q:      q,
		filter: filter,
		offset: syncoffset,
	}
}

type syncDiscovered struct {
	q      sqlx.Queryer
	filter rendezvousf
	offset string
	err    error
}

func (t *syncDiscovered) Each(ctx context.Context) iter.Seq[Discovered] {
	return func(yield func(Discovered) bool) {
		s := sqlx.Scan(DiscoveredSinceSync(ctx, t.q, t.offset))

		for v := range s.Iter() {
			if t.filter.Filter(v.Infohash) {
				continue
			}

			if !yield(v) {
				return
			}
		}

		t.err = s.Err()
	}
}

func (t *syncDiscovered) Err() error {
	return t.err
}
