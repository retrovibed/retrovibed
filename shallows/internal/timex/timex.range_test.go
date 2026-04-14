package timex_test

import (
	"testing"
	"time"

	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/stretchr/testify/require"
)

func TestNewRangeDuration(t *testing.T) {
	t.Run("end is after start", func(t *testing.T) {
		r := timex.NewRangeDuration(time.Hour)
		require.True(t, r.End.After(r.Start))
	})

	t.Run("span matches duration", func(t *testing.T) {
		d := 2 * time.Hour
		r := timex.NewRangeDuration(d)
		require.InDelta(t, d.Seconds(), r.End.Sub(r.Start).Seconds(), 1)
	})

	t.Run("end is approximately now", func(t *testing.T) {
		before := time.Now()
		r := timex.NewRangeDuration(time.Minute)
		require.False(t, r.End.Before(before))
		require.False(t, r.End.After(time.Now().Add(time.Second)))
	})
}

func TestNewRangeISO8601(t *testing.T) {
	t.Run("valid UTC timestamps", func(t *testing.T) {
		r, err := timex.NewRangeISO8601("2025-12-01T00:00:00Z", "2025-12-31T23:59:59Z")
		require.NoError(t, err)
		require.Equal(t, time.Date(2025, time.December, 1, 0, 0, 0, 0, time.UTC), r.Start)
		require.Equal(t, time.Date(2025, time.December, 31, 23, 59, 59, 0, time.UTC), r.End)
	})

	t.Run("valid timestamps with nanoseconds", func(t *testing.T) {
		r, err := timex.NewRangeISO8601("2025-12-01T00:00:00.123456789Z", "2025-12-31T23:59:59.999999999Z")
		require.NoError(t, err)
		require.Equal(t, time.Date(2025, time.December, 1, 0, 0, 0, 123456789, time.UTC), r.Start)
		require.Equal(t, time.Date(2025, time.December, 31, 23, 59, 59, 999999999, time.UTC), r.End)
	})

	t.Run("valid timestamps with timezone offset", func(t *testing.T) {
		r, err := timex.NewRangeISO8601("2025-12-01T00:00:00+05:00", "2025-12-31T23:59:59-08:00")
		require.NoError(t, err)
		loc := time.FixedZone("", 5*60*60)
		require.Equal(t, time.Date(2025, time.December, 1, 0, 0, 0, 0, loc), r.Start)
		require.Equal(t, time.Date(2025, time.December, 31, 23, 59, 59, 0, time.FixedZone("", -8*60*60)), r.End)
	})

	t.Run("invalid begin string", func(t *testing.T) {
		_, err := timex.NewRangeISO8601("not-a-timestamp", "2025-12-31T23:59:59Z")
		require.Error(t, err)
	})

	t.Run("invalid end string", func(t *testing.T) {
		_, err := timex.NewRangeISO8601("2025-12-01T00:00:00Z", "not-a-timestamp")
		require.Error(t, err)
	})

	t.Run("both invalid strings", func(t *testing.T) {
		_, err := timex.NewRangeISO8601("bad", "also-bad")
		require.Error(t, err)
	})

	t.Run("empty strings", func(t *testing.T) {
		_, err := timex.NewRangeISO8601("", "")
		require.Error(t, err)
	})

	t.Run("same start and end", func(t *testing.T) {
		ts := "2025-06-15T12:00:00Z"
		r, err := timex.NewRangeISO8601(ts, ts)
		require.NoError(t, err)
		require.True(t, r.Start.Equal(r.End))
	})

	t.Run("end before start is allowed", func(t *testing.T) {
		r, err := timex.NewRangeISO8601("2025-12-31T00:00:00Z", "2025-01-01T00:00:00Z")
		require.NoError(t, err)
		require.True(t, r.End.Before(r.Start))
	})

	t.Run("RFC3339Inf decodes to Inf", func(t *testing.T) {
		s := timex.RFC3339Inf().Format(time.RFC3339Nano)
		r, err := timex.NewRangeISO8601("2025-01-01T00:00:00Z", s)
		require.NoError(t, err)
		require.True(t, r.End.Equal(timex.Inf()), "expected Inf(), got %v", r.End)
	})

	t.Run("RFC3339NegInf decodes to NegInf", func(t *testing.T) {
		s := timex.RFC3339NegInf().Format(time.RFC3339Nano)
		r, err := timex.NewRangeISO8601(s, "2025-01-01T00:00:00Z")
		require.NoError(t, err)
		require.True(t, r.Start.Equal(timex.NegInf()), "expected NegInf(), got %v", r.Start)
	})

	t.Run("both inf and neginf", func(t *testing.T) {
		begin := timex.RFC3339NegInf().Format(time.RFC3339Nano)
		end := timex.RFC3339Inf().Format(time.RFC3339Nano)
		r, err := timex.NewRangeISO8601(begin, end)
		require.NoError(t, err)
		require.True(t, r.Start.Equal(timex.NegInf()), "expected NegInf(), got %v", r.Start)
		require.True(t, r.End.Equal(timex.Inf()), "expected Inf(), got %v", r.End)
	})
}
