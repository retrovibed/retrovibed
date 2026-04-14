package timex_test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/stretchr/testify/require"
)

func TestJSONSafeDecodeNowShouldRemainUnchanged(t *testing.T) {
	type foo struct {
		Timestamp time.Time
		Bar       struct {
			Timestamp time.Time
		}
	}

	ts := time.Now()
	tmp := timex.JSONSafeDecode(&foo{Timestamp: ts, Bar: struct{ Timestamp time.Time }{Timestamp: ts}})
	require.Equal(t, tmp.Timestamp, ts)
	require.Equal(t, tmp.Bar.Timestamp, ts)
}

func TestJSONSafeDecodeInfShouldBeAdjusted(t *testing.T) {
	type foo struct {
		Timestamp time.Time
		Bar       struct {
			Timestamp time.Time
		}
	}

	tmp := timex.JSONSafeDecode(&foo{Timestamp: timex.RFC3339Inf(), Bar: struct{ Timestamp time.Time }{Timestamp: timex.RFC3339Inf()}})
	log.Println(tmp.Timestamp, timex.Inf())
	require.Equal(t, tmp.Timestamp, timex.RFC3339NanoDecode(timex.Inf()))
	require.NotEqual(t, tmp.Timestamp, timex.RFC3339Inf())

	require.Equal(t, tmp.Bar.Timestamp, timex.RFC3339NanoDecode(timex.Inf()))
	require.NotEqual(t, tmp.Bar.Timestamp, timex.RFC3339Inf())
}

func TestJSONSafeEncodeNowShouldRemainUnchanged(t *testing.T) {
	type foo struct {
		Timestamp time.Time
		Bar       struct {
			Timestamp time.Time
		}
	}

	ts := time.Now()
	tmp := timex.JSONSafeEncode(&foo{Timestamp: ts, Bar: struct{ Timestamp time.Time }{Timestamp: ts}})
	require.Equal(t, tmp.Timestamp, ts)
	require.Equal(t, tmp.Bar.Timestamp, ts)
}

func TestJSONSafeEncodeInfShouldBeAdjusted(t *testing.T) {
	type foo struct {
		Timestamp time.Time
		Bar       struct {
			Timestamp time.Time
		}
	}

	tmp := timex.JSONSafeEncode(&foo{Timestamp: timex.Inf(), Bar: struct{ Timestamp time.Time }{Timestamp: timex.Inf()}})
	log.Println(tmp.Timestamp, timex.Inf())
	require.Equal(t, tmp.Timestamp, timex.RFC3339Inf())
	require.NotEqual(t, tmp.Timestamp, timex.Inf())

	require.Equal(t, tmp.Bar.Timestamp, timex.RFC3339Inf())
	require.NotEqual(t, tmp.Bar.Timestamp, timex.Inf())
}

func TestMax(t *testing.T) {
	expected := time.Now().Add(time.Hour)
	require.Equal(t, expected, timex.Max(time.UnixMicro(0), time.Now(), time.Now().Add(time.Minute), expected))
	require.Equal(t, timex.NegInf(), timex.Max(timex.NegInf()))
	require.Equal(t, timex.Inf(), timex.Max(timex.NegInf(), timex.Inf()))
}

func TestMin(t *testing.T) {
	require.Equal(t, time.UnixMicro(0), timex.Min(time.UnixMicro(0), time.Now(), time.Now().Add(time.Minute), time.Now().Add(time.Hour)))
	require.Equal(t, timex.Inf(), timex.Min(timex.Inf()))
	require.Equal(t, timex.NegInf(), timex.Min(timex.NegInf(), timex.Inf()))
	require.Equal(t, time.UnixMicro(0), timex.Min(time.UnixMicro(0), timex.Inf()))
}

func TestRFC3339NanoEncode(t *testing.T) {
	t.Run("time after max limit", func(t *testing.T) {
		maxTime := timex.RFC3339Inf()
		pastMax := maxTime.Add(time.Millisecond)
		encoded := timex.RFC3339NanoEncode(pastMax)
		require.WithinDuration(t, maxTime, encoded, 0)
	})

	t.Run("time before min limit", func(t *testing.T) {
		minTime := timex.RFC3339NegInf()
		beforeMin := minTime.Add(-time.Millisecond)
		encoded := timex.RFC3339NanoEncode(beforeMin)
		require.WithinDuration(t, minTime, encoded, 0)
	})

	t.Run("time within limits", func(t *testing.T) {
		ts := time.Date(2025, time.October, 25, 12, 30, 0, 500000000, time.UTC)
		encoded := timex.RFC3339NanoEncode(ts)
		require.WithinDuration(t, ts, encoded, 0)
	})

	t.Run("zero time", func(t *testing.T) {
		zero := time.Time{}
		encoded := timex.RFC3339NanoEncode(zero)
		require.Equal(t, zero, encoded)
		require.WithinDuration(t, zero, encoded, 0)
	})

	t.Run("time at max limit", func(t *testing.T) {
		maxTime := timex.RFC3339Inf()
		require.WithinDuration(t, maxTime, timex.RFC3339NanoEncode(maxTime), 0)
		require.WithinDuration(t, maxTime, timex.RFC3339NanoEncode(maxTime.Add(time.Nanosecond)), 0)
	})

	t.Run("time at min limit", func(t *testing.T) {
		minTime := timex.RFC3339NegInf()
		require.WithinDuration(t, minTime, timex.RFC3339NanoEncode(minTime), 0)
		require.WithinDuration(t, minTime, timex.RFC3339NanoEncode(minTime.Add(-time.Nanosecond)), 0)
	})
}

func TestRFC3339NanoDecode(t *testing.T) {
	t.Run("time after max limit", func(t *testing.T) {
		decoded := timex.RFC3339NanoDecode(timex.RFC3339Inf().Add(time.Nanosecond))
		require.WithinDuration(t, timex.Inf(), decoded, 0)
	})

	t.Run("time before min limit", func(t *testing.T) {
		require.WithinDuration(t, timex.NegInf(), timex.RFC3339NanoDecode(timex.RFC3339NegInf().Add(-time.Nanosecond)), 0)
	})

	t.Run("time within limits", func(t *testing.T) {
		now := time.Date(2025, time.October, 25, 12, 30, 0, 500000000, time.UTC)
		require.WithinDuration(t, now, timex.RFC3339NanoDecode(now), 0)
	})

	t.Run("time at max limit", func(t *testing.T) {
		require.WithinDuration(t, timex.Inf(), timex.RFC3339NanoDecode(timex.Inf()), 0)
	})

	t.Run("time at min limit", func(t *testing.T) {
		require.WithinDuration(t, timex.NegInf(), timex.RFC3339NanoDecode(timex.RFC3339NegInf()), 0)
	})
}

func TestRFC3339NanoEncodeDecodeRoundtrip(t *testing.T) {
	tests := []time.Time{
		{},
		time.Date(2025, time.October, 25, 12, 30, 0, 123456789, time.UTC),
		time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
		timex.RFC3339Inf().Add(-time.Nanosecond),
		timex.RFC3339NegInf().Add(time.Nanosecond),
		timex.Inf(),
		timex.NegInf(),
	}

	for i, ts := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			encoded := timex.RFC3339NanoEncode(ts)
			decoded := timex.RFC3339NanoDecode(encoded)
			require.WithinDuration(t, ts, decoded, 0)
		})
	}
}

func TestStartOfDay(t *testing.T) {
	t.Run("midday", func(t *testing.T) {
		ts := time.Date(2025, time.December, 5, 14, 30, 45, 123456789, time.UTC)
		result := timex.StartOfDay(ts)
		expected := time.Date(2025, time.December, 5, 0, 0, 0, 0, time.UTC)
		require.Equal(t, expected, result)
	})

	t.Run("already at start of day", func(t *testing.T) {
		ts := time.Date(2025, time.December, 5, 0, 0, 0, 0, time.UTC)
		result := timex.StartOfDay(ts)
		require.Equal(t, ts, result)
	})

	t.Run("end of day", func(t *testing.T) {
		ts := time.Date(2025, time.December, 5, 23, 59, 59, 999999999, time.UTC)
		result := timex.StartOfDay(ts)
		expected := time.Date(2025, time.December, 5, 0, 0, 0, 0, time.UTC)
		require.Equal(t, expected, result)
	})
}

func TestEndOfDay(t *testing.T) {
	t.Run("midday", func(t *testing.T) {
		ts := time.Date(2025, time.December, 5, 14, 30, 45, 123456789, time.UTC)
		result := timex.EndOfDay(ts)
		expected := time.Date(2025, time.December, 5, 23, 59, 59, 999999999, time.UTC)
		require.Equal(t, expected, result)
	})

	t.Run("start of day", func(t *testing.T) {
		ts := time.Date(2025, time.December, 5, 0, 0, 0, 0, time.UTC)
		result := timex.EndOfDay(ts)
		expected := time.Date(2025, time.December, 5, 23, 59, 59, 999999999, time.UTC)
		require.Equal(t, expected, result)
	})

	t.Run("already at end of day", func(t *testing.T) {
		ts := time.Date(2025, time.December, 5, 23, 59, 59, 999999999, time.UTC)
		result := timex.EndOfDay(ts)
		require.Equal(t, ts, result)
	})
}

func TestStartOfWeek(t *testing.T) {
	t.Run("wednesday returns monday", func(t *testing.T) {
		ts := time.Date(2025, time.December, 3, 14, 30, 45, 0, time.UTC)
		result := timex.StartOfWeek(ts)
		expected := time.Date(2025, time.December, 1, 0, 0, 0, 0, time.UTC)
		require.Equal(t, expected, result)
		require.Equal(t, time.Monday, result.Weekday())
	})

	t.Run("sunday returns previous monday", func(t *testing.T) {
		ts := time.Date(2025, time.December, 7, 14, 30, 45, 0, time.UTC)
		result := timex.StartOfWeek(ts)
		expected := time.Date(2025, time.December, 1, 0, 0, 0, 0, time.UTC)
		require.Equal(t, expected, result)
		require.Equal(t, time.Monday, result.Weekday())
	})

	t.Run("monday returns same monday", func(t *testing.T) {
		ts := time.Date(2025, time.December, 1, 14, 30, 45, 0, time.UTC)
		result := timex.StartOfWeek(ts)
		expected := time.Date(2025, time.December, 1, 0, 0, 0, 0, time.UTC)
		require.Equal(t, expected, result)
		require.Equal(t, time.Monday, result.Weekday())
	})

	t.Run("saturday returns monday", func(t *testing.T) {
		ts := time.Date(2025, time.December, 6, 23, 59, 59, 0, time.UTC)
		result := timex.StartOfWeek(ts)
		expected := time.Date(2025, time.December, 1, 0, 0, 0, 0, time.UTC)
		require.Equal(t, expected, result)
	})
}

func TestStartOfMonth(t *testing.T) {
	t.Run("mid month", func(t *testing.T) {
		ts := time.Date(2025, time.December, 15, 14, 30, 45, 123456789, time.UTC)
		result := timex.StartOfMonth(ts)
		expected := time.Date(2025, time.December, 1, 0, 0, 0, 0, time.UTC)
		require.Equal(t, expected, result)
	})

	t.Run("first of month", func(t *testing.T) {
		ts := time.Date(2025, time.December, 1, 0, 0, 0, 0, time.UTC)
		result := timex.StartOfMonth(ts)
		require.Equal(t, ts, result)
	})

	t.Run("last of month", func(t *testing.T) {
		ts := time.Date(2025, time.December, 31, 23, 59, 59, 999999999, time.UTC)
		result := timex.StartOfMonth(ts)
		expected := time.Date(2025, time.December, 1, 0, 0, 0, 0, time.UTC)
		require.Equal(t, expected, result)
	})

	t.Run("february", func(t *testing.T) {
		ts := time.Date(2025, time.February, 28, 14, 30, 45, 0, time.UTC)
		result := timex.StartOfMonth(ts)
		expected := time.Date(2025, time.February, 1, 0, 0, 0, 0, time.UTC)
		require.Equal(t, expected, result)
	})
}

func TestRange(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		r := timex.Range{}
		require.True(t, r.Start.IsZero())
		require.True(t, r.End.IsZero())
	})

	t.Run("with values", func(t *testing.T) {
		start := time.Date(2025, time.December, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(2025, time.December, 31, 23, 59, 59, 0, time.UTC)
		r := timex.Range{Start: start, End: end}
		require.Equal(t, start, r.Start)
		require.Equal(t, end, r.End)
	})

	t.Run("duration calculation", func(t *testing.T) {
		start := time.Date(2025, time.December, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(2025, time.December, 8, 0, 0, 0, 0, time.UTC)
		r := timex.Range{Start: start, End: end}
		require.Equal(t, 7*24*time.Hour, r.End.Sub(r.Start))
	})

	t.Run("month range", func(t *testing.T) {
		ts := time.Date(2025, time.December, 15, 14, 30, 0, 0, time.UTC)
		r := timex.Range{
			Start: timex.StartOfMonth(ts),
			End:   timex.EndOfDay(ts),
		}
		require.Equal(t, time.Date(2025, time.December, 1, 0, 0, 0, 0, time.UTC), r.Start)
		require.Equal(t, time.Date(2025, time.December, 15, 23, 59, 59, 999999999, time.UTC), r.End)
	})
}
