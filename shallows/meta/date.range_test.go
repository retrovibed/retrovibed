package meta_test

import (
	"testing"
	"time"

	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/stretchr/testify/require"
)

func TestNewDateRange(t *testing.T) {
	t.Run("normal range", func(t *testing.T) {
		start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(2025, 6, 15, 12, 30, 0, 0, time.UTC)

		dr := meta.NewDateRange(timex.Range{Start: start, End: end})

		require.Equal(t, start.Format(time.RFC3339Nano), dr.Oldest)
		require.Equal(t, end.Format(time.RFC3339Nano), dr.Newest)
	})

	t.Run("everything range clamps to RFC3339 sentinels", func(t *testing.T) {
		dr := meta.NewDateRange(timex.NewRangeEverything())

		require.Equal(t, timex.RFC3339NegInf().UTC().Format(time.RFC3339Nano), dr.Oldest)
		require.Equal(t, timex.RFC3339Inf().UTC().Format(time.RFC3339Nano), dr.Newest)
	})
}
