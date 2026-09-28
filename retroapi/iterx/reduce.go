package iterx

import (
	"context"
)

// Reduce folds every element of s into a single accumulator, starting from
// init. A non-nil error from fn aborts iteration immediately and is returned
// along with the accumulator built so far. Otherwise the error (if any) of s
// itself is returned once iteration completes.
func Reduce[T, A any](ctx context.Context, s Seq[T], init A, fn func(context.Context, A, T) (A, error)) (A, error) {
	acc := init
	for v := range s.Each(ctx) {
		next, err := fn(ctx, acc, v)
		if err != nil {
			return acc, err
		}
		acc = next
	}

	return acc, s.Err()
}
