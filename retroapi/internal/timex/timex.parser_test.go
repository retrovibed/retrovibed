package timex_test

import (
	"slices"
	"testing"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/internal/timex"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	t.Run("parses the values in order", func(t *testing.T) {
		a := time.Date(2026, 1, 2, 3, 4, 5, 6, time.UTC)
		b := time.Date(2027, 1, 2, 3, 4, 5, 6, time.UTC)

		parsed := slices.Collect(timex.Parser(time.RFC3339Nano, a.Format(time.RFC3339Nano), b.Format(time.RFC3339Nano)))
		require.Len(t, parsed, 2)
		require.True(t, a.Equal(parsed[0]))
		require.True(t, b.Equal(parsed[1]))
	})

	t.Run("decodes the json safe infinities", func(t *testing.T) {
		parsed := slices.Collect(timex.Parser(
			time.RFC3339Nano,
			timex.RFC3339Inf().Format(time.RFC3339Nano),
			timex.RFC3339NegInf().Format(time.RFC3339Nano),
		))
		require.Len(t, parsed, 2)
		require.True(t, parsed[0].Equal(timex.Inf()))
		require.True(t, parsed[1].Equal(timex.NegInf()))
	})

	t.Run("skips blank and unparseable values", func(t *testing.T) {
		a := time.Date(2026, 1, 2, 3, 4, 5, 6, time.UTC)

		parsed := slices.Collect(timex.Parser(time.RFC3339Nano, "", "garbage", a.Format(time.RFC3339Nano)))
		require.Len(t, parsed, 1)
		require.True(t, a.Equal(parsed[0]))
	})

	t.Run("stops when the consumer stops", func(t *testing.T) {
		a := time.Date(2026, 1, 2, 3, 4, 5, 6, time.UTC).Format(time.RFC3339Nano)

		var seen int
		for range timex.Parser(time.RFC3339Nano, a, a, a) {
			seen++
			break
		}
		require.Equal(t, 1, seen)
	})
}
