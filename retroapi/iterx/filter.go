package iterx

import (
	"context"
	"iter"
)

// Filter yields only the elements of s for which pred returns true. A
// non-nil error from pred aborts iteration immediately (same as an error
// from s itself) and is surfaced via Err().
func Filter[T any](s Seq[T], pred func(context.Context, T) (bool, error)) Seq[T] {
	return &filterSeq[T]{inner: s, pred: pred}
}

type filterSeq[T any] struct {
	inner Seq[T]
	pred  func(context.Context, T) (bool, error)
	err   error
}

func (t *filterSeq[T]) Each(ctx context.Context) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range t.inner.Each(ctx) {
			ok, err := t.pred(ctx, v)
			if err != nil {
				t.err = err
				return
			}
			if !ok {
				continue
			}
			if !yield(v) {
				return
			}
		}
	}
}

func (t *filterSeq[T]) Err() error {
	if t.err != nil {
		return t.err
	}
	return t.inner.Err()
}
