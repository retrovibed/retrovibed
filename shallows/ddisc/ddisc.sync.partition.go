package ddisc

import (
	"context"
	"iter"

	"github.com/retrovibed/retrovibed/internal/iterx"
	"github.com/retrovibed/retrovibed/internal/sqlx"
)

// send all discovered media matching the partition
func SyncPartition(q sqlx.Queryer, partition, syncoffset string) iterx.Seq[Discovered] {
	return &syncPartition{
		q:         q,
		partition: partition,
		offset:    syncoffset,
	}
}

type syncPartition struct {
	q         sqlx.Queryer
	partition string
	offset    string
	err       error
}

func (t *syncPartition) Each(ctx context.Context) iter.Seq[Discovered] {
	return func(yield func(Discovered) bool) {
		s := sqlx.Scan(DiscoveredPartitionSync(ctx, t.q, t.offset, t.partition))

		for v := range s.Iter() {
			if !yield(v) {
				return
			}
		}

		t.err = s.Err()
	}
}

func (t *syncPartition) Err() error {
	return t.err
}
