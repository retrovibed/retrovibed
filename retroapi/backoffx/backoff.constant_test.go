package backoffx_test

import (
	"testing"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/backoffx"
	"github.com/stretchr/testify/require"
)

func TestConstant(t *testing.T) {
	t.Run("should remain constant", func(t *testing.T) {
		s := backoffx.Constant(1 * time.Second)
		for i := 0; i < 5; i++ {
			require.Equal(t, 1*time.Second, s.Backoff(i), "attempt %d", i)
		}
	})
}
