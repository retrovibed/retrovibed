package tracking_test

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/autobind"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/james-lawrence/torrent/storage"
	"github.com/james-lawrence/torrent/torrenttest"
	"github.com/retrovibed/retrovibed/retroapi/blockcache"
	"github.com/retrovibed/retrovibed/retroapi/cryptox"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/james-lawrence/torrent/torrenttestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// provided a mapping of files to description, and returns the media and torrent directory.
func downloadTree(t *testing.T, q sqlx.Queryer, m []string, options ...metainfo.Option) (tracking.Metadata, string, string) {
	t.Helper()
	var (
		actual   = md5.New()
		expected = md5.New()
	)

	ctx := t.Context()

	seedir := t.TempDir()

	mi, err := torrenttest.Tree(
		seedir, cryptox.NewChaCha8(t.Name()), 64*bytesx.KiB, 128*bytesx.KiB,
		m,
		options...,
	)
	require.NoError(t, err)

	seeder := torrenttestx.Client(
		t,
		autobind.NewLoopback(
			autobind.EnableDHT(torrenttestx.QuickDHT(t)),
		),
		torrent.NewMetadataCache(seedir),
		blockcache.NewTorrentFromVirtualFS(fsx.DirVirtual(seedir)),
	)

	md, err := torrent.NewFromInfo(mi, torrent.OptionStorage(storage.NewFile(filepath.Join(seedir))))
	require.NoError(t, err)

	seederTorrent, _, err := seeder.Start(md)
	require.NoError(t, err)
	defer seeder.Close()

	require.NoError(t, torrent.Verify(ctx, seederTorrent))
	n, err := torrent.DownloadInto(ctx, expected, seederTorrent, torrent.TuneSeeding)
	require.NoError(t, err)
	require.Equal(t, mi.TotalLength(), n)

	root := fsx.DirVirtual(t.TempDir())

	leechdir := root.Path("torrent")
	mediadir := root.Path("media")
	require.NoError(t, fsx.MkDirs(0700, leechdir, mediadir))

	leecher := torrenttestx.Client(
		t,
		autobind.NewLoopback(
			autobind.EnableDHT(torrenttestx.QuickDHT(t)),
		),
		torrent.NewMetadataCache(leechdir),
		blockcache.NewTorrentFromVirtualFS(fsx.DirVirtual(leechdir)),
	)
	defer leecher.Close()

	lmd := tracking.NewMetadata(
		new(md.ID),
		tracking.MetadataOptionFromInfo(mi),
		tracking.MetadataOptionAutoDescription,
	)

	require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

	ltor, added, err := leecher.MaybeStart(
		torrent.NewFromInfo(mi),
	)
	require.NoError(t, err)
	assert.True(t, added)

	require.NoError(t, ltor.Tune(torrent.TuneClientPeer(seeder)))

	require.NoError(t, tracking.DownloadInto(t.Context(), q, root, library.QueryCleanerNoop(), &lmd, ltor, actual, asyncx.NewWakeup(t.Context())))
	require.Equal(t, md5x.FormatUUID(expected), md5x.FormatUUID(actual))

	return lmd, mediadir, leechdir
}

func TestDownloadInto(t *testing.T) {
	t.Run("properly download a multitorrent", func(t *testing.T) {
		var (
			actual   = md5.New()
			expected = md5.New()
		)

		ctx := t.Context()
		q := sqltestx.Metadatabase(t)

		seedir := t.TempDir()

		mi, err := torrenttest.RandomMulti(seedir, 5, 16*bytesx.KiB, 64*bytesx.KiB)
		require.NoError(t, err)

		seeder := torrenttestx.Client(
			t,
			autobind.NewLoopback(
				autobind.EnableDHT(torrenttestx.QuickDHT(t)),
			),
			torrent.NewMetadataCache(seedir),
			blockcache.NewTorrentFromVirtualFS(fsx.DirVirtual(seedir)),
		)

		md, err := torrent.NewFromInfo(mi, torrent.OptionStorage(storage.NewFile(filepath.Join(seedir))))
		require.NoError(t, err)

		seederTorrent, _, err := seeder.Start(md)
		require.NoError(t, err)
		defer seeder.Close()

		require.NoError(t, torrent.Verify(ctx, seederTorrent))
		n, err := torrent.DownloadInto(ctx, expected, seederTorrent, torrent.TuneSeeding)
		require.NoError(t, err)
		require.Equal(t, mi.TotalLength(), n)

		root := fsx.DirVirtual(t.TempDir())

		leechdir := root.Path("torrent")
		mediadir := root.Path("media")
		require.NoError(t, fsx.MkDirs(0700, leechdir, mediadir))

		leecher := torrenttestx.Client(
			t,
			autobind.NewLoopback(
				autobind.EnableDHT(torrenttestx.QuickDHT(t)),
			),
			torrent.NewMetadataCache(leechdir),
			blockcache.NewTorrentFromVirtualFS(fsx.DirVirtual(leechdir)),
		)
		defer leecher.Close()

		lmd := tracking.NewMetadata(
			new(md.ID),
			tracking.MetadataOptionFromInfo(mi),
			tracking.MetadataOptionAutoDescription,
		)

		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		ltor, added, err := leecher.MaybeStart(
			torrent.NewFromInfo(
				mi,
			),
		)
		require.NoError(t, err)
		assert.True(t, added)

		require.NoError(t, ltor.Tune(torrent.TuneClientPeer(seeder)))

		require.NoError(t, tracking.DownloadInto(t.Context(), q, root, library.QueryCleanerNoop(), &lmd, ltor, actual, asyncx.NewWakeup(t.Context())))

		require.Equal(t, md5x.FormatUUID(expected), md5x.FormatUUID(actual))

		w0 := fsx.WalkDir(os.DirFS(leechdir))
		require.EqualValues(t, 3, testx.Seq2Count(w0.Walk()))
		require.NoError(t, w0.Err())

		w1 := fsx.WalkDir(os.DirFS(mediadir))
		require.EqualValues(t, 6, testx.Seq2Count(w1.Walk()))
		require.NoError(t, w1.Err())

		var libMDs []library.Metadata
		require.NoError(t, sqlx.ScanInto(library.MetadataSearch(t.Context(), q, library.MetadataSearchBuilder().Where(
			library.MetadataQueryByTorrentID(lmd.ID),
		)), &libMDs))
		require.Len(t, libMDs, 5)
		for _, m := range libMDs {
			assert.Equal(t, lmd.ID, m.TorrentID)
			assert.Greater(t, m.Bytes, uint64(0))
			assert.NotEmpty(t, m.Description)
		}
	})

	t.Run("properly track bytes", func(t *testing.T) {
		var (
			actual   = md5.New()
			expected = md5.New()
		)

		ctx := t.Context()
		q := sqltestx.Metadatabase(t)

		seedir := t.TempDir()

		mi, err := torrenttest.RandomMulti(seedir, 5, 16*bytesx.KiB, 64*bytesx.KiB)
		require.NoError(t, err)

		seeder := torrenttestx.Client(
			t,
			autobind.NewLoopback(
				autobind.EnableDHT(torrenttestx.QuickDHT(t)),
			),
			torrent.NewMetadataCache(seedir),
			blockcache.NewTorrentFromVirtualFS(fsx.DirVirtual(seedir)),
		)

		md, err := torrent.NewFromInfo(mi, torrent.OptionStorage(storage.NewFile(filepath.Join(seedir))))
		require.NoError(t, err)

		seederTorrent, _, err := seeder.Start(md)
		require.NoError(t, err)
		defer seeder.Close()

		require.NoError(t, torrent.Verify(ctx, seederTorrent))
		n, err := torrent.DownloadInto(ctx, expected, seederTorrent, torrent.TuneSeeding)
		require.NoError(t, err)
		require.Equal(t, mi.TotalLength(), n)

		root := fsx.DirVirtual(t.TempDir())

		leechdir := root.Path("torrent")
		mediadir := root.Path("media")
		require.NoError(t, fsx.MkDirs(0700, leechdir, mediadir))

		leecher := torrenttestx.Client(
			t,
			autobind.NewLoopback(
				autobind.EnableDHT(torrenttestx.QuickDHT(t)),
			),
			torrent.NewMetadataCache(leechdir),
			blockcache.NewTorrentFromVirtualFS(fsx.DirVirtual(leechdir)),
		)
		defer leecher.Close()

		lmd := tracking.NewMetadata(
			new(md.ID),
			tracking.MetadataOptionFromInfo(mi),
			tracking.MetadataOptionAutoDescription,
		)

		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		ltor, added, err := leecher.MaybeStart(
			torrent.NewFromInfo(
				mi,
			),
		)
		require.NoError(t, err)
		assert.True(t, added)

		require.NoError(t, ltor.Tune(torrent.TuneClientPeer(seeder)))

		require.NoError(t, tracking.DownloadInto(t.Context(), q, root, library.QueryCleanerNoop(), &lmd, ltor, actual, asyncx.NewWakeup(t.Context())))

		require.Equal(t, md5x.FormatUUID(expected), md5x.FormatUUID(actual))

		w0 := fsx.WalkDir(os.DirFS(leechdir))
		require.EqualValues(t, 3, testx.Seq2Count(w0.Walk()))
		require.NoError(t, w0.Err())

		w1 := fsx.WalkDir(os.DirFS(mediadir))
		require.EqualValues(t, 6, testx.Seq2Count(w1.Walk()))
		require.NoError(t, w1.Err())

		require.NoError(t, tracking.MetadataFindByID(t.Context(), q, lmd.ID).Scan(&lmd))
		assert.EqualValues(t, mi.TotalLength(), lmd.Bytes)
		assert.WithinDuration(t, time.Now(), lmd.CompletedAt, time.Second)

		// a fresh leecher had nothing before this download, so every byte of
		// the completed torrent came from peers: available and downloaded
		// should both equal the full size.
		assert.EqualValues(t, mi.TotalLength(), lmd.Available)
		assert.EqualValues(t, mi.TotalLength(), lmd.Downloaded)

		var libMDs []library.Metadata
		require.NoError(t, sqlx.ScanInto(library.MetadataSearch(t.Context(), q, library.MetadataSearchBuilder().Where(
			library.MetadataQueryByTorrentID(lmd.ID),
		)), &libMDs))
		require.Len(t, libMDs, 5)
		for _, m := range libMDs {
			assert.Equal(t, lmd.ID, m.TorrentID)
			assert.Greater(t, m.Bytes, uint64(0))
			assert.NotEmpty(t, m.Description)
		}
	})

	t.Run("bluray torrent treated as single file", func(t *testing.T) {
		var (
			actual   = md5.New()
			expected = md5.New()
		)

		ctx := t.Context()
		q := sqltestx.Metadatabase(t)

		seedir := t.TempDir()

		mi, err := torrenttest.Tree(seedir, rand.Reader, 16*bytesx.KiB, 64*bytesx.KiB, []string{
			"BDMV/index.bdmv",
			"BDMV/MovieObject.bdmv",
			"BDMV/STREAM/00001.m2ts",
			"BDMV/STREAM/00002.m2ts",
			"CERTIFICATE/id.bdmv",
		})
		require.NoError(t, err)

		seeder := torrenttestx.Client(
			t,
			autobind.NewLoopback(
				autobind.EnableDHT(torrenttestx.QuickDHT(t)),
			),
			torrent.NewMetadataCache(seedir),
			blockcache.NewTorrentFromVirtualFS(fsx.DirVirtual(seedir)),
		)

		md, err := torrent.NewFromInfo(mi, torrent.OptionStorage(storage.NewFile(filepath.Join(seedir))))
		require.NoError(t, err)

		seederTorrent, _, err := seeder.Start(md)
		require.NoError(t, err)
		defer seeder.Close()

		require.NoError(t, torrent.Verify(ctx, seederTorrent))
		n, err := torrent.DownloadInto(ctx, expected, seederTorrent, torrent.TuneSeeding)
		require.NoError(t, err)
		require.Equal(t, mi.TotalLength(), n)

		root := fsx.DirVirtual(t.TempDir())

		leechdir := root.Path("torrent")
		mediadir := root.Path("media")
		require.NoError(t, fsx.MkDirs(0700, leechdir, mediadir))

		leecher := torrenttestx.Client(
			t,
			autobind.NewLoopback(
				autobind.EnableDHT(torrenttestx.QuickDHT(t)),
			),
			torrent.NewMetadataCache(leechdir),
			blockcache.NewTorrentFromVirtualFS(fsx.DirVirtual(leechdir)),
		)
		defer leecher.Close()

		lmd := tracking.NewMetadata(
			new(md.ID),
			tracking.MetadataOptionFromInfo(mi),
			tracking.MetadataOptionAutoDescription,
		)

		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		ltor, added, err := leecher.MaybeStart(
			torrent.NewFromInfo(
				mi,
			),
		)
		require.NoError(t, err)
		assert.True(t, added)

		require.NoError(t, ltor.Tune(torrent.TuneClientPeer(seeder)))

		require.NoError(t, tracking.DownloadInto(t.Context(), q, root, library.QueryCleanerNoop(), &lmd, ltor, actual, asyncx.NewWakeup(t.Context())))

		require.Equal(t, md5x.FormatUUID(expected), md5x.FormatUUID(actual))

		// bluray: entire torrent is treated as a single file, so only 1 symlink is created in the media dir
		w1 := fsx.WalkDir(os.DirFS(mediadir))
		require.EqualValues(t, 2, testx.Seq2Count(w1.Walk()))
		require.NoError(t, w1.Err())

		var libMDs []library.Metadata
		require.NoError(t, sqlx.ScanInto(library.MetadataSearch(t.Context(), q, library.MetadataSearchBuilder().Where(
			library.MetadataQueryByTorrentID(lmd.ID),
		)), &libMDs))
		require.Len(t, libMDs, 1)
		assert.Equal(t, lmd.ID, libMDs[0].TorrentID)
		assert.EqualValues(t, mi.TotalLength(), libMDs[0].Bytes)
		assert.NotEmpty(t, libMDs[0].Description)
	})

	t.Run("dvd torrent treated as single file", func(t *testing.T) {
		var (
			actual   = md5.New()
			expected = md5.New()
		)

		ctx := t.Context()
		q := sqltestx.Metadatabase(t)

		seedir := t.TempDir()

		mi, err := torrenttest.Tree(seedir, rand.Reader, 16*bytesx.KiB, 64*bytesx.KiB, []string{
			"VIDEO_TS/VIDEO_TS.IFO",
			"VIDEO_TS/VIDEO_TS.VOB",
			"VIDEO_TS/VIDEO_TS.BUP",
			"VIDEO_TS/VTS_01_0.IFO",
			"VIDEO_TS/VTS_01_1.VOB",
		})
		require.NoError(t, err)

		seeder := torrenttestx.Client(
			t,
			autobind.NewLoopback(
				autobind.EnableDHT(torrenttestx.QuickDHT(t)),
			),
			torrent.NewMetadataCache(seedir),
			blockcache.NewTorrentFromVirtualFS(fsx.DirVirtual(seedir)),
		)

		md, err := torrent.NewFromInfo(mi, torrent.OptionStorage(storage.NewFile(filepath.Join(seedir))))
		require.NoError(t, err)

		seederTorrent, _, err := seeder.Start(md)
		require.NoError(t, err)
		defer seeder.Close()

		require.NoError(t, torrent.Verify(ctx, seederTorrent))
		n, err := torrent.DownloadInto(ctx, expected, seederTorrent, torrent.TuneSeeding)
		require.NoError(t, err)
		require.Equal(t, mi.TotalLength(), n)

		root := fsx.DirVirtual(t.TempDir())

		leechdir := root.Path("torrent")
		mediadir := root.Path("media")
		require.NoError(t, fsx.MkDirs(0700, leechdir, mediadir))

		leecher := torrenttestx.Client(
			t,
			autobind.NewLoopback(
				autobind.EnableDHT(torrenttestx.QuickDHT(t)),
			),
			torrent.NewMetadataCache(leechdir),
			blockcache.NewTorrentFromVirtualFS(fsx.DirVirtual(leechdir)),
		)
		defer leecher.Close()

		lmd := tracking.NewMetadata(
			new(md.ID),
			tracking.MetadataOptionFromInfo(mi),
			tracking.MetadataOptionAutoDescription,
		)

		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		ltor, added, err := leecher.MaybeStart(
			torrent.NewFromInfo(
				mi,
			),
		)
		require.NoError(t, err)
		assert.True(t, added)

		require.NoError(t, ltor.Tune(torrent.TuneClientPeer(seeder)))

		require.NoError(t, tracking.DownloadInto(t.Context(), q, root, library.QueryCleanerNoop(), &lmd, ltor, actual, asyncx.NewWakeup(t.Context())))

		require.Equal(t, md5x.FormatUUID(expected), md5x.FormatUUID(actual))

		// dvd: entire torrent is treated as a single file, so only 1 symlink is created in the media dir
		w1 := fsx.WalkDir(os.DirFS(mediadir))
		require.EqualValues(t, 2, testx.Seq2Count(w1.Walk()))
		require.NoError(t, w1.Err())

		var libMDs []library.Metadata
		require.NoError(t, sqlx.ScanInto(library.MetadataSearch(t.Context(), q, library.MetadataSearchBuilder().Where(
			library.MetadataQueryByTorrentID(lmd.ID),
		)), &libMDs))
		require.Len(t, libMDs, 1)
		assert.Equal(t, lmd.ID, libMDs[0].TorrentID)
		assert.EqualValues(t, mi.TotalLength(), libMDs[0].Bytes)
		assert.NotEmpty(t, libMDs[0].Description)
	})

	t.Run("torrent tree description from tracking metadata propagates to library entries", func(t *testing.T) {
		var (
			actual   = md5.New()
			expected = md5.New()
		)

		ctx := t.Context()
		q := sqltestx.Metadatabase(t)

		seedir := t.TempDir()

		mi, err := torrenttest.Tree(seedir, rand.Reader, 16*bytesx.KiB, 64*bytesx.KiB, []string{
			"file1.mkv",
			"file2.mkv",
			"file3.mkv",
		})
		require.NoError(t, err)

		seeder := torrenttestx.Client(
			t,
			autobind.NewLoopback(
				autobind.EnableDHT(torrenttestx.QuickDHT(t)),
			),
			torrent.NewMetadataCache(seedir),
			blockcache.NewTorrentFromVirtualFS(fsx.DirVirtual(seedir)),
		)

		md, err := torrent.NewFromInfo(mi, torrent.OptionStorage(storage.NewFile(filepath.Join(seedir))))
		require.NoError(t, err)

		seederTorrent, _, err := seeder.Start(md)
		require.NoError(t, err)
		defer seeder.Close()

		require.NoError(t, torrent.Verify(ctx, seederTorrent))
		n, err := torrent.DownloadInto(ctx, expected, seederTorrent, torrent.TuneSeeding)
		require.NoError(t, err)
		require.Equal(t, mi.TotalLength(), n)

		root := fsx.DirVirtual(t.TempDir())

		leechdir := root.Path("torrent")
		mediadir := root.Path("media")
		require.NoError(t, fsx.MkDirs(0700, leechdir, mediadir))

		leecher := torrenttestx.Client(
			t,
			autobind.NewLoopback(
				autobind.EnableDHT(torrenttestx.QuickDHT(t)),
			),
			torrent.NewMetadataCache(leechdir),
			blockcache.NewTorrentFromVirtualFS(fsx.DirVirtual(leechdir)),
		)
		defer leecher.Close()

		const rawDesc = "My.Test.Movie.2024"
		lmd := tracking.NewMetadata(
			new(md.ID),
			tracking.MetadataOptionFromInfo(mi),
			tracking.MetadataOptionDescription(rawDesc),
			tracking.MetadataOptionAutoDescription,
		)

		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		ltor, added, err := leecher.MaybeStart(
			torrent.NewFromInfo(
				mi,
			),
		)
		require.NoError(t, err)
		assert.True(t, added)

		require.NoError(t, ltor.Tune(torrent.TuneClientPeer(seeder)))

		require.NoError(t, tracking.DownloadInto(t.Context(), q, root, library.QueryCleanerNoop(), &lmd, ltor, actual, asyncx.NewWakeup(t.Context())))

		require.Equal(t, md5x.FormatUUID(expected), md5x.FormatUUID(actual))

		scanner := sqlx.Scan(library.MetadataSearch(t.Context(), q, library.MetadataSearchBuilder().Where(
			library.MetadataQueryByTorrentID(lmd.ID),
		)))

		libMDs := slices.Collect(scanner.Iter())
		require.NoError(t, scanner.Err())
		require.Len(t, libMDs, 3)

		// With QueryCleanerNoop the cleaner returns the input unchanged.
		// as a result we'll see `{Torrent Name} {File Name}` as the description.
		wantDesc := []string{
			fmt.Sprintf("%s %s", rawDesc, "file1.mkv"),
			fmt.Sprintf("%s %s", rawDesc, "file2.mkv"),
			fmt.Sprintf("%s %s", rawDesc, "file3.mkv"),
		}

		for _, m := range libMDs {
			assert.Equal(t, lmd.ID, m.TorrentID)
			assert.Contains(t, wantDesc, m.Description, "library metadata should match one of the normalized description")
		}
	})

	t.Run("single file torrent creates one media entry named from tracking metadata", func(t *testing.T) {
		var (
			actual   = md5.New()
			expected = md5.New()
		)

		ctx := t.Context()
		q := sqltestx.Metadatabase(t)

		seedir := t.TempDir()

		mi, _, err := torrenttest.Random(seedir, 64*bytesx.KiB)
		require.NoError(t, err)

		seeder := torrenttestx.Client(
			t,
			autobind.NewLoopback(
				autobind.EnableDHT(torrenttestx.QuickDHT(t)),
			),
			torrent.NewMetadataCache(seedir),
			blockcache.NewTorrentFromVirtualFS(fsx.DirVirtual(seedir)),
		)

		md, err := torrent.NewFromInfo(mi, torrent.OptionStorage(storage.NewFile(filepath.Join(seedir))))
		require.NoError(t, err)

		seederTorrent, _, err := seeder.Start(md)
		require.NoError(t, err)
		defer seeder.Close()

		require.NoError(t, torrent.Verify(ctx, seederTorrent))
		n, err := torrent.DownloadInto(ctx, expected, seederTorrent, torrent.TuneSeeding)
		require.NoError(t, err)
		require.Equal(t, mi.TotalLength(), n)

		root := fsx.DirVirtual(t.TempDir())

		leechdir := root.Path("torrent")
		mediadir := root.Path("media")
		require.NoError(t, fsx.MkDirs(0700, leechdir, mediadir))

		leecher := torrenttestx.Client(
			t,
			autobind.NewLoopback(
				autobind.EnableDHT(torrenttestx.QuickDHT(t)),
			),
			torrent.NewMetadataCache(leechdir),
			blockcache.NewTorrentFromVirtualFS(fsx.DirVirtual(leechdir)),
		)
		defer leecher.Close()

		lmd := tracking.NewMetadata(
			new(md.ID),
			tracking.MetadataOptionFromInfo(mi),
			tracking.MetadataOptionAutoDescription,
		)

		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		ltor, added, err := leecher.MaybeStart(
			torrent.NewFromInfo(mi),
		)
		require.NoError(t, err)
		assert.True(t, added)

		require.NoError(t, ltor.Tune(torrent.TuneClientPeer(seeder)))

		require.NoError(t, tracking.DownloadInto(t.Context(), q, root, library.QueryCleanerNoop(), &lmd, ltor, actual, asyncx.NewWakeup(t.Context())))

		require.Equal(t, md5x.FormatUUID(expected), md5x.FormatUUID(actual))

		// single file torrent: exactly one symlink in the media dir (dir + symlink = 2 walk entries)
		w1 := fsx.WalkDir(os.DirFS(mediadir))
		require.EqualValues(t, 2, testx.Seq2Count(w1.Walk()))
		require.NoError(t, w1.Err())

		libMD, err := sqlx.ScanOne(library.MetadataSearch(t.Context(), q, library.MetadataSearchBuilder().Where(
			library.MetadataQueryByTorrentID(lmd.ID),
		)))
		require.NoError(t, err)

		assert.Equal(t, lmd.ID, libMD.TorrentID)
		assert.EqualValues(t, mi.TotalLength(), libMD.Bytes)
		// QueryCleanerNoop returns md.Description unchanged; so the library entry carries
		// the tracking metadata's description verbatim.
		assert.Equal(t, lmd.Description, libMD.Description)
	})

	t.Run("properly extract descriptions for library content", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)

		m := map[string]string{
			"Star Trek The Next Generation (1987) S02E01 (1080p NF WEB-DL DDP5.1 AV1) - Vialle.mkv": "Star Trek The Next Generation 1987 S02E01 Vialle.mkv",
			"Star Trek The Next Generation (1987) S02E02 (1080p NF WEB-DL DDP5.1 AV1) - Vialle.mkv": "Star Trek The Next Generation 1987 S02E02 Vialle.mkv",
			"Star Trek The Next Generation (1987) S02E03 (1080p NF WEB-DL DDP5.1 AV1) - Vialle.mkv": "Star Trek The Next Generation 1987 S02E03 Vialle.mkv",
			"Star Trek The Next Generation (1987) S02E04 (1080p NF WEB-DL DDP5.1 AV1) - Vialle.mkv": "Star Trek The Next Generation 1987 S02E04 Vialle.mkv",
			"Star Trek The Next Generation (1987) S02E05 (1080p NF WEB-DL DDP5.1 AV1) - Vialle.mkv": "Star Trek The Next Generation 1987 S02E05 Vialle.mkv",
			"Star Trek The Next Generation (1987) S02E06 (1080p NF WEB-DL DDP5.1 AV1) - Vialle.mkv": "Star Trek The Next Generation 1987 S02E06 Vialle.mkv",
			"Star Trek The Next Generation (1987) S02E07 (1080p NF WEB-DL DDP5.1 AV1) - Vialle.mkv": "Star Trek The Next Generation 1987 S02E07 Vialle.mkv",
			"Star Trek The Next Generation (1987) S02E08 (1080p NF WEB-DL DDP5.1 AV1) - Vialle.mkv": "Star Trek The Next Generation 1987 S02E08 Vialle.mkv",
			"Star Trek The Next Generation (1987) S02E09 (1080p NF WEB-DL DDP5.1 AV1) - Vialle.mkv": "Star Trek The Next Generation 1987 S02E09 Vialle.mkv",
			"Star Trek The Next Generation (1987) S02E10 (1080p NF WEB-DL DDP5.1 AV1) - Vialle.mkv": "Star Trek The Next Generation 1987 S02E10 Vialle.mkv",
			"Star Trek The Next Generation (1987) S02E11 (1080p NF WEB-DL DDP5.1 AV1) - Vialle.mkv": "Star Trek The Next Generation 1987 S02E11 Vialle.mkv",
			"Star Trek The Next Generation (1987) S02E12 (1080p NF WEB-DL DDP5.1 AV1) - Vialle.mkv": "Star Trek The Next Generation 1987 S02E12 Vialle.mkv",
			"Star Trek The Next Generation (1987) S02E13 (1080p NF WEB-DL DDP5.1 AV1) - Vialle.mkv": "Star Trek The Next Generation 1987 S02E13 Vialle.mkv",
			"Star Trek The Next Generation (1987) S02E14 (1080p NF WEB-DL DDP5.1 AV1) - Vialle.mkv": "Star Trek The Next Generation 1987 S02E14 Vialle.mkv",
			"Star Trek The Next Generation (1987) S02E15 (1080p NF WEB-DL DDP5.1 AV1) - Vialle.mkv": "Star Trek The Next Generation 1987 S02E15 Vialle.mkv",
			"Star Trek The Next Generation (1987) S02E16 (1080p NF WEB-DL DDP5.1 AV1) - Vialle.mkv": "Star Trek The Next Generation 1987 S02E16 Vialle.mkv",
			"Star Trek The Next Generation (1987) S02E17 (1080p NF WEB-DL DDP5.1 AV1) - Vialle.mkv": "Star Trek The Next Generation 1987 S02E17 Vialle.mkv",
			"Star Trek The Next Generation (1987) S02E18 (1080p NF WEB-DL DDP5.1 AV1) - Vialle.mkv": "Star Trek The Next Generation 1987 S02E18 Vialle.mkv",
			"Star Trek The Next Generation (1987) S02E19 (1080p NF WEB-DL DDP5.1 AV1) - Vialle.mkv": "Star Trek The Next Generation 1987 S02E19 Vialle.mkv",
			"Star Trek The Next Generation (1987) S02E20 (1080p NF WEB-DL DDP5.1 AV1) - Vialle.mkv": "Star Trek The Next Generation 1987 S02E20 Vialle.mkv",
			"Star Trek The Next Generation (1987) S02E21 (1080p NF WEB-DL DDP5.1 AV1) - Vialle.mkv": "Star Trek The Next Generation 1987 S02E21 Vialle.mkv",
			"Star Trek The Next Generation (1987) S02E22 (1080p NF WEB-DL DDP5.1 AV1) - Vialle.mkv": "Star Trek The Next Generation 1987 S02E22 Vialle.mkv",
		}

		lmd, mediadir, _ := downloadTree(t, q, slices.Collect(maps.Keys(m)), metainfo.OptionDisplayName("Star Trek The Next Generation (1987) S02 (1080p NF WEB-DL DDP5.1 AV1) - Vialle"))

		// single file torrent: exactly one symlink in the media dir (dir + symlinks = 23 walk entries)
		w1 := fsx.WalkDir(os.DirFS(mediadir))
		require.EqualValues(t, 23, testx.Seq2Count(w1.Walk()))
		require.NoError(t, w1.Err())

		require.EqualValues(
			t,
			len(m),
			sqltestx.Count(t, q, "SELECT COUNT(*) FROM library_metadata WHERE torrent_id = ?", lmd.ID),
		)

		s := sqlx.Scan(library.MetadataSearch(t.Context(), q, library.MetadataSearchBuilder().Where(
			library.MetadataQueryByTorrentID(lmd.ID),
		)))

		for lmd := range s.Iter() {
			_, ok := m[lmd.Description]
			require.True(t, ok, "library description missing: %s", lmd.Description)
		}

		require.NoError(t, s.Err())
	})

	t.Run("InfoFromPath and FileInfoFromOffset resolve each file's real path after a download", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)

		names := []string{"file1.mkv", "file2.mkv", "file3.mkv"}
		lmd, _, leechdir := downloadTree(t, q, names)

		p := fsx.DirVirtual(leechdir).Path(int160.FromBytes(lmd.Infohash).String())

		info, err := tracking.InfoFromPath(p)
		require.NoError(t, err)
		require.EqualValues(t, lmd.Bytes, info.TotalLength())

		scanner := sqlx.Scan(library.MetadataSearch(t.Context(), q, library.MetadataSearchBuilder().Where(
			library.MetadataQueryByTorrentID(lmd.ID),
		)))
		libMDs := slices.Collect(scanner.Iter())
		require.NoError(t, scanner.Err())
		require.Len(t, libMDs, 3)

		resolved := map[string]bool{}
		for _, m := range libMDs {
			fi, err := tracking.FileInfoFromOffset(p, m.DiskOffset)
			require.NoError(t, err)
			require.Contains(t, names, fi.Path)
			require.EqualValues(t, m.Bytes, fi.Length)
			resolved[fi.Path] = true
		}
		require.Len(t, resolved, 3, "each library.Metadata row should resolve to a distinct file")
	})
}
