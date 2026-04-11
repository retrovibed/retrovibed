package asyncx_test

import (
	"context"
	"runtime"
	"testing"

	"github.com/retrovibed/retrovibed/internal/asyncx"
	"github.com/stretchr/testify/require"
)

func TestWakeup(t *testing.T) {
	t.Run("wakeup when signalled", func(t *testing.T) {
		var (
			c  uint
			cc = make(chan struct{}, 1)
		)

		wakeup := asyncx.NewWakeup(t.Context())
		cc <- struct{}{}

		for range 6 {
			select {
			case <-cc:
				runtime.Gosched() // forces scheduling
				wakeup.Broadcast()
			case <-wakeup.C:
				c += 1
				cc <- struct{}{}
			case <-t.Context().Done():
				require.NoError(t, t.Context().Err())
			}
		}

		require.NoError(t, wakeup.Close())
		require.EqualValues(t, 3, c)
	})
}

func TestRun(t *testing.T) {
	t.Run("run to completion", func(t *testing.T) {
		var (
			c uint
		)
		ctx, done := context.WithCancelCause(t.Context())
		wakeup := asyncx.NewWakeup(ctx)
		go func() {
			for {
				if ctx.Err() != nil {
					return
				}
				wakeup.Broadcast()
			}
		}()
		require.NoError(t, asyncx.Run(ctx, wakeup, func(ctx context.Context) error {
			if c < 3 {
				c += 1
				return nil
			}

			done(nil)
			return nil
		}))

		require.NoError(t, wakeup.Close())
		require.EqualValues(t, c, uint(3))
	})
}
