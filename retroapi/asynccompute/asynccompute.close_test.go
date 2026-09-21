package asynccompute_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/asynccompute"
	"github.com/stretchr/testify/require"
)

func TestPoolClose(t *testing.T) {
	t.Run("repeated close does not panic and returns nil", func(t *testing.T) {
		p := asynccompute.New(func(ctx context.Context, w int) error { return nil })

		require.NoError(t, p.Close())
		require.NotPanics(t, func() {
			require.NoError(t, p.Close())
			require.NoError(t, p.Close())
		})
	})

	t.Run("repeated close waits for and does not rerun completed work", func(t *testing.T) {
		var processed atomic.Int64
		p := asynccompute.New(func(ctx context.Context, w int) error {
			processed.Add(1)
			return nil
		})

		for i := 0; i < 10; i++ {
			require.NoError(t, p.Run(t.Context(), i))
		}

		require.NoError(t, p.Close())
		require.EqualValues(t, 10, processed.Load())
		require.NoError(t, p.Close())
		require.EqualValues(t, 10, processed.Load())
	})

	t.Run("repeated close returns the same failure", func(t *testing.T) {
		expected := errors.New("workload failed")
		p := asynccompute.New(func(ctx context.Context, w int) error { return expected })

		require.NoError(t, p.Run(t.Context(), 1))

		require.ErrorIs(t, p.Close(), expected)
		require.ErrorIs(t, p.Close(), expected)
	})

	t.Run("concurrent close does not panic", func(t *testing.T) {
		p := asynccompute.New(func(ctx context.Context, w int) error { return nil })

		require.NoError(t, p.Run(t.Context(), 1))

		var wg sync.WaitGroup
		require.NotPanics(t, func() {
			for i := 0; i < 8; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					require.NoError(t, p.Close())
				}()
			}
			wg.Wait()
		})
	})
}
