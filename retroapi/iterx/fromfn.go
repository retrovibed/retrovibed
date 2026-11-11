package iterx

import (
	"context"
	"iter"
)

// FromFn adapts fn into a Seq. fn pushes values to yield and stops when yield
// returns false; the error it returns is reported by Err. When the consumer
// stops early (yield returned false) fn's error is discarded, since it is
// expected teardown noise such as the cancellation used to stop in-flight work.
func FromFn[T any](fn func(ctx context.Context, yield func(T) bool) error) Seq[T] {
	return &fnSeq[T]{fn: fn}
}

type fnSeq[T any] struct {
	fn  func(ctx context.Context, yield func(T) bool) error
	err error
}

func (t *fnSeq[T]) Each(ctx context.Context) iter.Seq[T] {
	return func(yield func(T) bool) {
		stopped := false
		t.err = t.fn(ctx, func(v T) bool {
			if !yield(v) {
				stopped = true
				return false
			}

			return true
		})

		if stopped {
			t.err = nil
		}
	}
}

func (t *fnSeq[T]) Err() error {
	return t.err
}
