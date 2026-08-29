package backoffx_test

import (
	"testing"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/backoffx"
	"github.com/stretchr/testify/require"
)

func TestJitterRandom(t *testing.T) {
	t.Run("stays within [x, x+d)", func(t *testing.T) {
		const (
			x = 100 * time.Millisecond
			d = 20 * time.Millisecond
		)
		s := backoffx.New(backoffx.Constant(x), backoffx.JitterRandom(d))

		for range 10000 {
			got := s.Backoff(0)
			require.GreaterOrEqual(t, got, x)
			require.Less(t, got, x+d)
		}
	})

	t.Run("reaches values near the base x", func(t *testing.T) {
		const (
			x = 100 * time.Millisecond
			d = 20 * time.Millisecond
		)
		s := backoffx.New(backoffx.Constant(x), backoffx.JitterRandom(d))

		min := x + d
		for range 10000 {
			if got := s.Backoff(0); got < min {
				min = got
			}
		}

		require.Less(t, min, x+d/10)
	})

	t.Run("reaches values near the upper bound x+d", func(t *testing.T) {
		const (
			x = 100 * time.Millisecond
			d = 20 * time.Millisecond
		)
		s := backoffx.New(backoffx.Constant(x), backoffx.JitterRandom(d))

		max := x
		for range 10000 {
			if got := s.Backoff(0); got > max {
				max = got
			}
		}

		require.Greater(t, max, x+d-d/10)
	})
}
