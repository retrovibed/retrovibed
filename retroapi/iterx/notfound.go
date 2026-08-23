package iterx

import (
	"context"
	"iter"

	"github.com/retrovibed/retrovibed/retroapi/errorsx"
)

// ErrNotFound is returned by a NotFound-wrapped Seq's Err() when Each never
// yielded a value and the wrapped sequence otherwise completed without error.
const ErrNotFound = errorsx.String("iterx: not found")

// NotFound wraps s so Err() returns ErrNotFound if Each never yielded a value
// and the underlying sequence otherwise completed without error.
func NotFound[T any](s Seq[T]) Seq[T] {
	return &notFoundSeq[T]{inner: s}
}

type notFoundSeq[T any] struct {
	inner Seq[T]
	found bool
}

func (t *notFoundSeq[T]) Each(ctx context.Context) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range t.inner.Each(ctx) {
			t.found = true
			if !yield(v) {
				return
			}
		}
	}
}

func (t *notFoundSeq[T]) Err() error {
	if err := t.inner.Err(); err != nil {
		return err
	}
	if !t.found {
		return ErrNotFound
	}
	return nil
}
