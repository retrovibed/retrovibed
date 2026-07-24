package iterx

import (
	"context"
	"iter"
)

// Take yields at most n elements of s, then stops pulling from s. Take(s, 0)
// yields no elements; callers that want "0 means unlimited" semantics should
// only wrap with Take when n > 0.
func Take[T any](s Seq[T], n uint64) Seq[T] {
	return &takeSeq[T]{inner: s, n: n}
}

type takeSeq[T any] struct {
	inner Seq[T]
	n     uint64
}

func (t *takeSeq[T]) Each(ctx context.Context) iter.Seq[T] {
	return func(yield func(T) bool) {
		i := uint64(0)
		for v := range t.inner.Each(ctx) {
			if i >= t.n {
				return
			}
			if !yield(v) {
				return
			}

			i++
		}
	}
}

func (t *takeSeq[T]) Err() error {
	return t.inner.Err()
}
