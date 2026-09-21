package library_test

import (
	"math"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestKnownStringCollationEpisode(t *testing.T) {
	t.Run("parses season and episode", func(t *testing.T) {
		require.Equal(t, library.KnownCollationEpisode(1, 2), library.KnownStringCollationEpisode("S01E02"))
		require.Equal(t, library.KnownCollationEpisode(12, 345), library.KnownStringCollationEpisode("S12E345"))
	})

	t.Run("packs season into the high 16 bits and episode into the low 16 bits", func(t *testing.T) {
		require.Equal(t, uint32(3)<<16|uint32(7), library.KnownStringCollationEpisode("S03E07"))
	})

	t.Run("season zero maps to the specials season", func(t *testing.T) {
		result := library.KnownStringCollationEpisode("S00E03")
		require.Equal(t, library.KnownCollationEpisode(0, 3), result)
		require.Equal(t, uint32(library.KnownCollationSpecialsSeason)<<16|uint32(3), result)
	})

	t.Run("returns MaxUint32 when the S prefix is missing", func(t *testing.T) {
		require.Equal(t, uint32(math.MaxUint32), library.KnownStringCollationEpisode("1E2"))
		require.Equal(t, uint32(math.MaxUint32), library.KnownStringCollationEpisode("01E02"))
	})

	t.Run("returns MaxUint32 for blank input", func(t *testing.T) {
		require.Equal(t, uint32(math.MaxUint32), library.KnownStringCollationEpisode(""))
		require.Equal(t, uint32(math.MaxUint32), library.KnownStringCollationEpisode("   "))
	})

	t.Run("returns MaxUint32 when there is no episode separator", func(t *testing.T) {
		require.Equal(t, uint32(math.MaxUint32), library.KnownStringCollationEpisode("S01"))
		require.Equal(t, uint32(math.MaxUint32), library.KnownStringCollationEpisode("0102"))
	})

	t.Run("returns MaxUint32 for a non numeric season", func(t *testing.T) {
		require.Equal(t, uint32(math.MaxUint32), library.KnownStringCollationEpisode("SxxE02"))
		require.Equal(t, uint32(math.MaxUint32), library.KnownStringCollationEpisode("E02"))
	})

	t.Run("returns MaxUint32 for a non numeric episode", func(t *testing.T) {
		require.Equal(t, uint32(math.MaxUint32), library.KnownStringCollationEpisode("S01Exx"))
		require.Equal(t, uint32(math.MaxUint32), library.KnownStringCollationEpisode("S01E"))
	})
}
