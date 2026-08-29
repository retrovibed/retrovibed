package backoffx_test

import (
	"testing"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/internal/backoffx"
	"github.com/stretchr/testify/require"
)

func TestCycle(t *testing.T) {
	t.Run("more attempts than delays", func(t *testing.T) {
		s := backoffx.Cycle(1*time.Second, 2*time.Second, 3*time.Second)
		expected := []time.Duration{1 * time.Second, 2 * time.Second, 3 * time.Second, 1 * time.Second, 2 * time.Second}
		for i, want := range expected {
			require.Equal(t, want, s.Backoff(i), "attempt %d", i)
		}
	})
}
