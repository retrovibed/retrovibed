package blockcache_test

import (
	"crypto/md5"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/autobind"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/james-lawrence/torrent/storage"
	"github.com/james-lawrence/torrent/torrenttest"
	"github.com/retrovibed/retrovibed/retroapi/blockcache"
	"github.com/retrovibed/retrovibed/retroapi/internal/bytesx"
	"github.com/retrovibed/retrovibed/retroapi/internal/fsx"
	"github.com/retrovibed/retrovibed/retroapi/internal/md5x"
	"github.com/retrovibed/retrovibed/retroapi/internal/torrenttestx"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestingConfig(t testing.TB, mdstore torrent.MetadataStore, store storage.ClientImpl, options ...torrent.ClientConfigOption) *torrent.ClientConfig {
	return torrent.NewDefaultClientConfig(
		mdstore,
		store,
		torrent.ClientConfigCacheDirectory(t.TempDir()),
		torrent.ClientConfigDebugLogger(log.New(os.Stderr, "[debug] ", log.Flags())),
		torrent.ClientConfigCompose(options...),
	)
}

func testClientTransfer(t *testing.T, seedercfg, leechercfg torrent.ClientConfigOption) {
	var (
		actual   = md5.New()
		expected = md5.New()
	)

	ctx := t.Context()
	seedir := t.TempDir()

	mi, err := torrenttest.RandomMulti(seedir, 3, 16*bytesx.MiB, 64*bytesx.MiB)
	require.NoError(t, err)

	// Create seeder and a Torrent.
	cfg := TestingConfig(
		t,
		torrent.NewMetadataCache(seedir),
		blockcache.NewTorrentFromVirtualFS(fsx.DirVirtual(seedir)),
		torrent.ClientConfigSeed(true),
		seedercfg,
	)

	seeder, err := autobind.NewLoopback(
		autobind.EnableDHT(torrenttestx.QuickDHT(t)),
	).Bind(torrent.NewClient(cfg))
	require.NoError(t, err)

	md, err := torrent.NewFromInfo(mi, torrent.OptionStorage(storage.NewFile(filepath.Join(seedir))))
	require.NoError(t, err)

	seederTorrent, _, err := seeder.Start(md)
	require.NoError(t, err)
	// Run a Stats right after Closing the Client. This will trigger the Stats
	// panic in #214 caused by RemoteAddr on Closed uTP sockets.
	defer seederTorrent.Stats()
	defer seeder.Close()

	require.NoError(t, torrent.Verify(ctx, seederTorrent))
	n, err := torrent.DownloadInto(ctx, expected, seederTorrent, torrent.TuneSeeding)
	require.NoError(t, err)
	require.Equal(t, mi.TotalLength(), n)

	leechdir := t.TempDir()
	cfg = TestingConfig(
		t,
		torrent.NewMetadataCache(leechdir),
		blockcache.NewTorrentFromVirtualFS(fsx.DirVirtual(leechdir)),
		torrent.ClientConfigSeed(false),
		leechercfg,
	)

	leecher, err := autobind.NewLoopback(
		autobind.EnableDHT(torrenttestx.QuickDHT(t)),
	).Bind(torrent.NewClient(cfg))
	require.NoError(t, err)
	defer leecher.Close()

	leecherTorrent, added, err := leecher.MaybeStart(
		torrent.NewFromInfo(
			mi,
		),
	)
	require.NoError(t, err)
	assert.True(t, added)

	// Now do some things with leecher and seeder.
	require.NoError(t, leecherTorrent.Tune(torrent.TuneClientPeer(seeder)))

	// The Torrent should not be interested in obtaining peers, so the one we
	// just added should be the only one.
	require.False(t, leecherTorrent.Stats().Seeding)

	// begin downloading
	_, err = torrent.DownloadInto(ctx, actual, leecherTorrent)
	require.NoError(t, err)

	// fsx.PrintFS(os.DirFS(leechdir))
	// log.Println("WAAAT", md5x.FormatUUID(expected), md5x.FormatUUID(actual))

	// seederStats := seederTorrent.Stats()
	// assert.GreaterOrEqual(t, mi.Length, seederStats.BytesWrittenData.Int64())

	// leecherStats := leecherTorrent.Stats()
	// assert.GreaterOrEqual(t, mi.Length, leecherStats.BytesReadData.Int64())
	require.Equal(t, md5x.FormatUUID(expected), md5x.FormatUUID(actual))
}

func TestClientTransferDefault(t *testing.T) {
	testClientTransfer(t, torrent.ClientConfigNoop, torrent.ClientConfigNoop)
}

func TestTorrentFileSystem(t *testing.T) {
	t.Run("WalkDir yields all torrent files", func(t *testing.T) {
		dir := t.TempDir()
		info, err := torrenttest.RandomMulti(dir, 5, bytesx.KiB, bytesx.MiB)
		require.NoError(t, err)

		bcache, err := blockcache.NewDirectoryCache(dir)
		require.NoError(t, err)
		fsys := blockcache.TorrentFilesystem(bcache, info)

		count := 0
		err = fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			for d := range metainfo.Files(info) {
				if path == d.Path {
					count++
					return nil
				}
			}

			return nil
		})
		require.NoError(t, err)
		require.Equal(t, 5, count)
	})

	t.Run("Stat reports correct file size", func(t *testing.T) {
		dir := t.TempDir()
		info, err := torrenttest.RandomMulti(dir, 5, bytesx.KiB, bytesx.MiB)
		require.NoError(t, err)

		bcache, err := blockcache.NewDirectoryCache(dir)
		require.NoError(t, err)
		fsys := blockcache.TorrentFilesystem(bcache, info)

		for fn := range metainfo.Files(info) {
			si, serr := fsys.Stat(fn.Path)
			require.NoError(t, serr)
			require.Equal(t, int64(fn.Length), si.Size())
		}
	})

	t.Run("Open returns file with correct torrent byte offset", func(t *testing.T) {
		dir := t.TempDir()
		info, err := torrenttest.RandomMulti(dir, 5, bytesx.KiB, bytesx.MiB)
		require.NoError(t, err)

		bcache, err := blockcache.NewDirectoryCache(dir)
		require.NoError(t, err)
		fsys := blockcache.TorrentFilesystem(bcache, info)

		for fn := range metainfo.Files(info) {
			f, ferr := fsys.Open(fn.Path)
			require.NoError(t, ferr)
			defer f.Close()
			bf, ok := f.(*blockcache.File)
			require.True(t, ok, "Open returned unexpected type for %s", fn.Path)
			require.Equal(t, fn.Offset, bf.Offset)
		}
	})

	t.Run("ReadDir with n>0 advances through entries across calls", func(t *testing.T) {
		dir := t.TempDir()
		info, err := torrenttest.RandomMulti(dir, 4, bytesx.KiB, bytesx.MiB)
		require.NoError(t, err)

		bcache, err := blockcache.NewDirectoryCache(dir)
		require.NoError(t, err)
		fsys := blockcache.TorrentFilesystem(bcache, info)

		root, err := fsys.Open(".")
		require.NoError(t, err)
		rdf := root.(fs.ReadDirFile)

		batch1, err := rdf.ReadDir(2)
		require.NoError(t, err)
		require.Len(t, batch1, 2)

		batch2, err := rdf.ReadDir(2)
		require.NoError(t, err)
		require.Len(t, batch2, 2)

		names1 := map[string]bool{batch1[0].Name(): true, batch1[1].Name(): true}
		for _, e := range batch2 {
			require.False(t, names1[e.Name()], "entry %s appeared in both batches", e.Name())
		}

		_, err = rdf.ReadDir(1)
		require.ErrorIs(t, err, io.EOF)
	})

	t.Run("Open returns independent directory instances", func(t *testing.T) {
		dir := t.TempDir()
		root := filepath.Join(dir, "torrent")
		require.NoError(t, os.MkdirAll(filepath.Join(root, "subdir"), 0700))
		require.NoError(t, os.WriteFile(filepath.Join(root, "subdir", "file.txt"), []byte("x"), 0600))

		info, err := metainfo.NewFromPath(root)
		require.NoError(t, err)

		bcache, err := blockcache.NewDirectoryCache(dir)
		require.NoError(t, err)
		fsys := blockcache.TorrentFilesystem(bcache, info)

		d1, err := fsys.Open("subdir")
		require.NoError(t, err)
		d2, err := fsys.Open("subdir")
		require.NoError(t, err)

		require.False(t, d1 == d2, "Open returned the same directory instance twice")
	})

	t.Run("WalkDir does not return duplicates for files in subdirectories", func(t *testing.T) {
		dir := t.TempDir()
		root := filepath.Join(dir, "torrent")
		require.NoError(t, os.MkdirAll(filepath.Join(root, "subdir"), 0700))
		require.NoError(t, os.WriteFile(filepath.Join(root, "subdir", "file1.txt"), []byte("aaa"), 0600))
		require.NoError(t, os.WriteFile(filepath.Join(root, "subdir", "file2.txt"), []byte("bbb"), 0600))

		info, err := metainfo.NewFromPath(root)
		require.NoError(t, err)

		bcache, err := blockcache.NewDirectoryCache(dir)
		require.NoError(t, err)
		fsys := blockcache.TorrentFilesystem(bcache, info)

		seen := make(map[string]int)
		err = fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				seen[path]++
			}
			return nil
		})
		require.NoError(t, err)
		require.Len(t, seen, 2)
		for path, count := range seen {
			require.Equal(t, 1, count, "file %s returned %d times", path, count)
		}
	})

	t.Run("Open and Stat return ErrNotExist for missing paths", func(t *testing.T) {
		dir := t.TempDir()
		info, err := torrenttest.RandomMulti(dir, 2, bytesx.KiB, bytesx.MiB)
		require.NoError(t, err)

		bcache, err := blockcache.NewDirectoryCache(dir)
		require.NoError(t, err)
		fsys := blockcache.TorrentFilesystem(bcache, info)

		_, err = fsys.Open("does/not/exist.txt")
		require.ErrorIs(t, err, fs.ErrNotExist)

		_, err = fsys.Stat("does/not/exist.txt")
		require.ErrorIs(t, err, fs.ErrNotExist)
	})
}
