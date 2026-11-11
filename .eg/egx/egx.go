package egx

import (
	"context"

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
