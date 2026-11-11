package egx

import (
	"context"
	"time"

	"github.com/egdaemon/eg/runtime/wasi/eg"
)

// Fallback tries each op in order and returns nil on the first that succeeds.
// Returns the last error if all ops fail.
func Fallback(ops ...eg.OpFn) eg.OpFn {
	return func(ctx context.Context, o eg.Op) error {
		var err error
		for _, op := range ops {
			if err = op(ctx, o); err == nil {
				return nil
			}
		}
		return err
	}
}

// RetryUntilSuccess retries the operation until it succeeds or the context
// is cancelled, waiting delay between attempts.
func RetryUntilSuccess(delay time.Duration, op eg.OpFn) eg.OpFn {
	return func(ctx context.Context, o eg.Op) error {
		for {
			err := op(ctx, o)
			if err == nil {
				return nil
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}
	}
}
