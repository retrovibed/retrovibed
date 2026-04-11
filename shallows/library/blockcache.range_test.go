package library

import (
	"log"
	"math/rand/v2"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/internal/bytesx"
	"github.com/retrovibed/retrovibed/internal/cryptox"
	"github.com/stretchr/testify/require"
)

func TestCalculateBlockRange(t *testing.T) {
	t.Run("first block", func(t *testing.T) {
		doffset, dlength := calculateBlockRange(32, 0, 64)
		require.Equal(t, uint64(0), doffset)
		require.Equal(t, uint64(32), dlength)
	})

	t.Run("second block", func(t *testing.T) {
		doffset, dlength := calculateBlockRange(32, 32, 64)
		require.Equal(t, uint64(32), doffset)
		require.Equal(t, uint64(32), dlength)
	})

	t.Run("partial block at start", func(t *testing.T) {
		doffset, dlength := calculateBlockRange(32, 0, 33)
		require.Equal(t, uint64(0), doffset)
		require.Equal(t, uint64(32), dlength)
	})

	t.Run("partial block at end", func(t *testing.T) {
		doffset, dlength := calculateBlockRange(32, 33, 33)
		require.EqualValues(t, 32, int(doffset))
		require.EqualValues(t, 1, int(dlength))
	})

	t.Run("random lengths - every offset", func(t *testing.T) {
		const blength = 64
		prng := rand.New(cryptox.NewChaCha8(uuid.Must(uuid.NewV4()).Bytes()))
		length := prng.Uint64N(128*bytesx.KiB) + bytesx.KiB
		for offset := range length {
			eoffset := (offset / blength) * blength
			doffset, dlength := calculateBlockRange(blength, offset, length)

			require.EqualValues(t, eoffset, doffset, "blength(%d) offset(%d) length(%d)", blength, offset, length)

			elength := uint64(blength)
			if eoffset+blength > length {
				elength = length - eoffset
			}
			log.Printf("blength(%d) offset(%d) length(%d)\n", blength, offset, length)
			require.EqualValues(t, elength, int(dlength), "blength(%d) offset(%d) length(%d)", blength, offset, length)
		}
	})

	t.Run("random lengths - actual blocks", func(t *testing.T) {
		const blength = 64
		prng := rand.New(cryptox.NewChaCha8(uuid.Must(uuid.NewV4()).Bytes()))
		length := prng.Uint64N(128*bytesx.KiB) + bytesx.KiB
		for offset := uint64(0); offset < length; {
			eoffset := (offset / blength) * blength
			doffset, dlength := calculateBlockRange(blength, offset, length)

			require.EqualValues(t, eoffset, doffset, "blength(%d) offset(%d) length(%d)", blength, offset, length)
			require.EqualValues(t, eoffset%blength, 0, "offset must always be a multiple of 64 for chacha")

			elength := uint64(blength)
			if eoffset+blength > length {
				elength = length - eoffset
			}

			require.EqualValues(t, elength, int(dlength), "blength(%d) offset(%d) length(%d)", blength, offset, length)
			offset += dlength
		}
	})
}
