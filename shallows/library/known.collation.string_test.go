package library_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestKnownCollationString(t *testing.T) {
	t.Run("blank for the standalone collation", func(t *testing.T) {
		require.Equal(t, "", library.KnownCollationString(0))
	})

	t.Run("formats season and episode", func(t *testing.T) {
		require.Equal(t, "S01E02", library.KnownCollationString(library.KnownCollationEpisode(1, 2)))
		require.Equal(t, "S03E07", library.KnownCollationString(uint32(3)<<16|uint32(7)))
	})

	t.Run("does not truncate wide values", func(t *testing.T) {
		require.Equal(t, "S12E345", library.KnownCollationString(library.KnownCollationEpisode(12, 345)))
	})

	t.Run("specials season renders as season zero", func(t *testing.T) {
		require.Equal(t, "S00E03", library.KnownCollationString(library.KnownCollationEpisode(0, 3)))
	})

	t.Run("episode zero is not blank", func(t *testing.T) {
		require.Equal(t, "S01E00", library.KnownCollationString(library.KnownCollationEpisode(1, 0)))
	})

	t.Run("round trips through KnownStringCollationEpisode", func(t *testing.T) {
		for _, s := range []string{"S01E02", "S12E345", "S00E03"} {
			require.Equal(t, s, library.KnownCollationString(library.KnownStringCollationEpisode(s)))
		}
	})
}
