package blockcache_test

import (
	"crypto/md5"
	"io"
	"math/rand/v2"
	"testing"

	"github.com/retrovibed/retrovibed/blockcache"
	"github.com/retrovibed/retrovibed/internal/bytesx"
	"github.com/retrovibed/retrovibed/internal/cryptox"
	"github.com/retrovibed/retrovibed/internal/md5x"
	"github.com/stretchr/testify/require"
)

func blockcachecycletest(t *testing.T, dcache *blockcache.DirCache, nsize int64, src io.Reader) {
	var (
		expected = md5.New()
		actual   = md5.New()
	)

	n, err := io.Copy(io.NewOffsetWriter(dcache, 0), io.LimitReader(io.TeeReader(src, expected), nsize))
	require.NoError(t, err)
	require.Equal(t, nsize, n)

	n, err = io.Copy(actual, io.NewSectionReader(dcache, 0, nsize))
	require.NoError(t, err)
	require.Equal(t, nsize, n)
	require.Equal(t, md5x.FormatUUID(expected), md5x.FormatUUID(actual))
}

func TestDirCache(t *testing.T) {
	t.Run("single byte chunks", func(t *testing.T) {
		tmpdir := t.TempDir()
		dcache, err := blockcache.NewDirectoryCache(tmpdir, blockcache.OptionDirCacheBlockLength(1))
		require.NoError(t, err)

		blockcachecycletest(t, dcache, 5, cryptox.NewChaCha8(t.Name()))
	})

	t.Run("block length", func(t *testing.T) {
		tmpdir := t.TempDir()
		dcache, err := blockcache.NewDirectoryCache(tmpdir)
		require.NoError(t, err)

		blockcachecycletest(t, dcache, dcache.BlockLength, cryptox.NewChaCha8(t.Name()))
	})

	t.Run("block length + 1", func(t *testing.T) {
		tmpdir := t.TempDir()
		dcache, err := blockcache.NewDirectoryCache(tmpdir)
		require.NoError(t, err)

		blockcachecycletest(t, dcache, dcache.BlockLength+1, cryptox.NewChaCha8(t.Name()))
	})
}

func FuzzDirCache(f *testing.F) {
	f.Fuzz(func(t *testing.T, seed string, n uint8) {
		tmpdir := t.TempDir()
		dcache, err := blockcache.NewDirectoryCache(tmpdir)
		require.NoError(t, err)
		blockcachecycletest(t, dcache, rand.Int64N(int64(n+1)*bytesx.MiB), cryptox.NewChaCha8(seed))
	})
}
