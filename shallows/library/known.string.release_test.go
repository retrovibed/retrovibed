package library_test

import (
	"testing"
	"time"

	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestKnownStringRelease(t *testing.T) {
	t.Run("parses a date", func(t *testing.T) {
		require.Equal(t, time.Date(2018, time.July, 27, 0, 0, 0, 0, time.UTC), library.KnownStringRelease("2018-07-27"))
	})

	t.Run("is always in UTC regardless of the local timezone", func(t *testing.T) {
		local := time.Local
		time.Local = time.FixedZone("test", -8*60*60)
		defer func() { time.Local = local }()

		result := library.KnownStringRelease("2018-07-27")
		require.Equal(t, time.UTC, result.Location())
		require.Equal(t, time.Date(2018, time.July, 27, 0, 0, 0, 0, time.UTC).Unix(), result.Unix())
	})

	t.Run("ignores surrounding whitespace", func(t *testing.T) {
		require.Equal(t, time.Date(2008, time.January, 20, 0, 0, 0, 0, time.UTC), library.KnownStringRelease(" 2008-01-20\n"))
	})

	t.Run("returns the zero time for blank input", func(t *testing.T) {
		require.True(t, library.KnownStringRelease("").IsZero())
		require.True(t, library.KnownStringRelease("   ").IsZero())
	})

	t.Run("returns the zero time for a malformed date", func(t *testing.T) {
		require.True(t, library.KnownStringRelease("2018").IsZero())
		require.True(t, library.KnownStringRelease("2018-07").IsZero())
		require.True(t, library.KnownStringRelease("07-27-2018").IsZero())
		require.True(t, library.KnownStringRelease("garbage").IsZero())
	})

	t.Run("returns the zero time for an impossible date", func(t *testing.T) {
		require.True(t, library.KnownStringRelease("2018-13-40").IsZero())
	})
}
