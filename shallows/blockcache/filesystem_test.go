package blockcache_test

import (
	"crypto/rand"
	"io/fs"
	"testing"

	"github.com/james-lawrence/torrent/torrenttest"
	"github.com/retrovibed/retrovibed/shallows/blockcache"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTorrentSingleFile(t *testing.T) {
	t.Run("exposes the torrent as a single file regardless of file count", func(t *testing.T) {
		seedir := t.TempDir()
		mi, err := torrenttest.Tree(seedir, rand.Reader, 16*bytesx.KiB, 64*bytesx.KiB, []string{
			"BDMV/index.bdmv",
			"BDMV/MovieObject.bdmv",
			"BDMV/STREAM/00001.m2ts",
			"BDMV/STREAM/00002.m2ts",
			"CERTIFICATE/id.bdmv",
		})
		require.NoError(t, err)

		dcache, err := blockcache.NewDirectoryCache(t.TempDir())
		require.NoError(t, err)

		vfs := blockcache.TorrentSingleFile(dcache, mi)

		var entries []fs.DirEntry
		require.NoError(t, fs.WalkDir(vfs, ".", func(path string, d fs.DirEntry, err error) error {
			require.NoError(t, err)
			if !d.IsDir() {
				entries = append(entries, d)
			}
			return nil
		}))

		require.Len(t, entries, 1)
		info, err := entries[0].Info()
		require.NoError(t, err)
		assert.Equal(t, mi.Name, info.Name())
		assert.EqualValues(t, mi.TotalLength(), info.Size())
	})

	t.Run("stat finds the single file by name", func(t *testing.T) {
		seedir := t.TempDir()
		mi, err := torrenttest.Tree(seedir, rand.Reader, 16*bytesx.KiB, 64*bytesx.KiB, []string{
			"BDMV/index.bdmv",
			"BDMV/STREAM/00001.m2ts",
		})
		require.NoError(t, err)

		dcache, err := blockcache.NewDirectoryCache(t.TempDir())
		require.NoError(t, err)

		vfs := blockcache.TorrentSingleFile(dcache, mi)

		info, err := vfs.Stat(mi.Name)
		require.NoError(t, err)
		assert.Equal(t, mi.Name, info.Name())
		assert.EqualValues(t, mi.TotalLength(), info.Size())
		assert.False(t, info.IsDir())
	})

	t.Run("single file size equals total torrent length", func(t *testing.T) {
		seedir := t.TempDir()
		mi, err := torrenttest.RandomMulti(seedir, 5, 16*bytesx.KiB, 64*bytesx.KiB)
		require.NoError(t, err)

		dcache, err := blockcache.NewDirectoryCache(t.TempDir())
		require.NoError(t, err)

		vfs := blockcache.TorrentSingleFile(dcache, mi)

		info, err := vfs.Stat(mi.Name)
		require.NoError(t, err)
		assert.EqualValues(t, mi.TotalLength(), info.Size())
	})
}
