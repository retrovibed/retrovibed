package tracking_test

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/autobind"
	"github.com/james-lawrence/torrent/storage"
	"github.com/james-lawrence/torrent/torrenttest"
	"github.com/retrovibed/retrovibed/retroapi/blockcache"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/torrenttestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
			langx.Autoptr(md.ID),
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

		require.NoError(t, tracking.DownloadInto(t.Context(), q, root, library.QueryCleanerNoop(), &lmd, ltor, actual))

		require.Equal(t, md5x.FormatUUID(expected), md5x.FormatUUID(actual))

		w0 := fsx.Walk(os.DirFS(leechdir))
		require.EqualValues(t, 3, testx.Seq2Count(w0.Walk()))
		require.NoError(t, w0.Err())

		w1 := fsx.Walk(os.DirFS(mediadir))
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
			langx.Autoptr(md.ID),
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

		require.NoError(t, tracking.DownloadInto(t.Context(), q, root, library.QueryCleanerNoop(), &lmd, ltor, actual))

		require.Equal(t, md5x.FormatUUID(expected), md5x.FormatUUID(actual))

		w0 := fsx.Walk(os.DirFS(leechdir))
		require.EqualValues(t, 3, testx.Seq2Count(w0.Walk()))
		require.NoError(t, w0.Err())

		w1 := fsx.Walk(os.DirFS(mediadir))
		require.EqualValues(t, 6, testx.Seq2Count(w1.Walk()))
		require.NoError(t, w1.Err())

		require.NoError(t, tracking.MetadataFindByID(t.Context(), q, lmd.ID).Scan(&lmd))
		assert.EqualValues(t, mi.TotalLength(), lmd.Bytes)
		assert.WithinDuration(t, time.Now(), lmd.CompletedAt, time.Second)

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
			langx.Autoptr(md.ID),
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

		require.NoError(t, tracking.DownloadInto(t.Context(), q, root, library.QueryCleanerNoop(), &lmd, ltor, actual))

		require.Equal(t, md5x.FormatUUID(expected), md5x.FormatUUID(actual))

		// bluray: entire torrent is treated as a single file, so only 1 symlink is created in the media dir
		w1 := fsx.Walk(os.DirFS(mediadir))
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
			langx.Autoptr(md.ID),
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

		require.NoError(t, tracking.DownloadInto(t.Context(), q, root, library.QueryCleanerNoop(), &lmd, ltor, actual))

		require.Equal(t, md5x.FormatUUID(expected), md5x.FormatUUID(actual))

		// dvd: entire torrent is treated as a single file, so only 1 symlink is created in the media dir
		w1 := fsx.Walk(os.DirFS(mediadir))
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
			langx.Autoptr(md.ID),
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

		require.NoError(t, tracking.DownloadInto(t.Context(), q, root, library.QueryCleanerNoop(), &lmd, ltor, actual))

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
			langx.Autoptr(md.ID),
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

		require.NoError(t, tracking.DownloadInto(t.Context(), q, root, library.QueryCleanerNoop(), &lmd, ltor, actual))

		require.Equal(t, md5x.FormatUUID(expected), md5x.FormatUUID(actual))

		// single file torrent: exactly one symlink in the media dir (dir + symlink = 2 walk entries)
		w1 := fsx.Walk(os.DirFS(mediadir))
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
}
