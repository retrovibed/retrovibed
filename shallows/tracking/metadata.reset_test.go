package tracking_test

import (
	"crypto/rand"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/autobind"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/storage"
	"github.com/james-lawrence/torrent/torrenttest"
	"github.com/retrovibed/retrovibed/retroapi/blockcache"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/james-lawrence/torrent/torrenttestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReset(t *testing.T) {
	t.Run("resets tracking metadata progress fields after a completed download", func(t *testing.T) {
		ctx := t.Context()
		q := sqltestx.Metadatabase(t)

		seedir := t.TempDir()
		mi, err := torrenttest.Tree(seedir, rand.Reader, 16*bytesx.KiB, 64*bytesx.KiB, []string{
			"file1.mkv", "file2.mkv", "file3.mkv",
		})
		require.NoError(t, err)

		seeder := torrenttestx.Client(
			t,
			autobind.NewLoopback(autobind.EnableDHT(torrenttestx.QuickDHT(t))),
			torrent.NewMetadataCache(seedir),
			blockcache.NewTorrentFromVirtualFS(fsx.DirVirtual(seedir)),
		)
		md, err := torrent.NewFromInfo(mi, torrent.OptionStorage(storage.NewFile(filepath.Join(seedir))))
		require.NoError(t, err)
		seederTorrent, _, err := seeder.Start(md)
		require.NoError(t, err)
		defer seeder.Close()
		require.NoError(t, torrent.Verify(ctx, seederTorrent))
		_, err = torrent.DownloadInto(ctx, io.Discard, seederTorrent, torrent.TuneSeeding)
		require.NoError(t, err)

		root := fsx.DirVirtual(t.TempDir())
		require.NoError(t, fsx.MkDirs(0700, root.Path("torrent"), root.Path("media")))

		leecher := torrenttestx.Client(
			t,
			autobind.NewLoopback(autobind.EnableDHT(torrenttestx.QuickDHT(t))),
			torrent.NewMetadataCache(root.Path("torrent")),
			blockcache.NewTorrentFromVirtualFS(fsx.DirVirtual(root.Path("torrent"))),
		)
		defer leecher.Close()

		lmd := tracking.NewMetadata(
			new(md.ID),
			tracking.MetadataOptionFromInfo(mi),
			tracking.MetadataOptionAutoDescription,
		)
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		ltor, added, err := leecher.MaybeStart(torrent.NewFromInfo(mi))
		require.NoError(t, err)
		assert.True(t, added)
		require.NoError(t, ltor.Tune(torrent.TuneClientPeer(seeder)))

		require.NoError(t, tracking.DownloadInto(t.Context(), q, root, library.QueryCleanerNoop(), &lmd, ltor, io.Discard, asyncx.NewWakeup(t.Context())))

		require.NoError(t, tracking.Reset(ctx, q, root, &lmd))

		assert.EqualValues(t, 0, lmd.Downloaded)
		assert.EqualValues(t, 0, lmd.Available)
		assert.EqualValues(t, 0, lmd.Uploaded)
		assert.False(t, lmd.Seeding)
		assert.EqualValues(t, 0, lmd.Peers)
		assert.Equal(t, timex.Inf(), lmd.InitiatedAt)
		assert.Equal(t, timex.Inf(), lmd.CompletedAt)
	})

	t.Run("tombstones library metadata entries after a completed download", func(t *testing.T) {
		ctx := t.Context()
		q := sqltestx.Metadatabase(t)

		seedir := t.TempDir()
		mi, err := torrenttest.Tree(seedir, rand.Reader, 16*bytesx.KiB, 64*bytesx.KiB, []string{
			"file1.mkv", "file2.mkv", "file3.mkv",
		})
		require.NoError(t, err)

		seeder := torrenttestx.Client(
			t,
			autobind.NewLoopback(autobind.EnableDHT(torrenttestx.QuickDHT(t))),
			torrent.NewMetadataCache(seedir),
			blockcache.NewTorrentFromVirtualFS(fsx.DirVirtual(seedir)),
		)
		md, err := torrent.NewFromInfo(mi, torrent.OptionStorage(storage.NewFile(filepath.Join(seedir))))
		require.NoError(t, err)
		seederTorrent, _, err := seeder.Start(md)
		require.NoError(t, err)
		defer seeder.Close()
		require.NoError(t, torrent.Verify(ctx, seederTorrent))
		_, err = torrent.DownloadInto(ctx, io.Discard, seederTorrent, torrent.TuneSeeding)
		require.NoError(t, err)

		root := fsx.DirVirtual(t.TempDir())
		require.NoError(t, fsx.MkDirs(0700, root.Path("torrent"), root.Path("media")))

		leecher := torrenttestx.Client(
			t,
			autobind.NewLoopback(autobind.EnableDHT(torrenttestx.QuickDHT(t))),
			torrent.NewMetadataCache(root.Path("torrent")),
			blockcache.NewTorrentFromVirtualFS(fsx.DirVirtual(root.Path("torrent"))),
		)
		defer leecher.Close()

		lmd := tracking.NewMetadata(
			new(md.ID),
			tracking.MetadataOptionFromInfo(mi),
			tracking.MetadataOptionAutoDescription,
		)
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		ltor, added, err := leecher.MaybeStart(torrent.NewFromInfo(mi))
		require.NoError(t, err)
		assert.True(t, added)
		require.NoError(t, ltor.Tune(torrent.TuneClientPeer(seeder)))

		require.NoError(t, tracking.DownloadInto(t.Context(), q, root, library.QueryCleanerNoop(), &lmd, ltor, io.Discard, asyncx.NewWakeup(t.Context())))

		require.NoError(t, tracking.Reset(ctx, q, root, &lmd))

		// all 3 library entries should be tombstoned (not deleted)
		assert.Equal(t, 3, sqltestx.Count(t, q,
			"SELECT COUNT(*) FROM library_metadata WHERE tombstoned_at < 'infinity' AND torrent_id = $1",
			lmd.ID,
		))
		assert.Equal(t, 3, sqltestx.Count(t, q,
			"SELECT COUNT(*) FROM library_metadata WHERE torrent_id = $1",
			lmd.ID,
		))
	})

	t.Run("removes media directories for all library entries after a completed download", func(t *testing.T) {
		ctx := t.Context()
		q := sqltestx.Metadatabase(t)

		seedir := t.TempDir()
		mi, err := torrenttest.Tree(seedir, rand.Reader, 16*bytesx.KiB, 64*bytesx.KiB, []string{
			"file1.mkv", "file2.mkv", "file3.mkv",
		})
		require.NoError(t, err)

		seeder := torrenttestx.Client(
			t,
			autobind.NewLoopback(autobind.EnableDHT(torrenttestx.QuickDHT(t))),
			torrent.NewMetadataCache(seedir),
			blockcache.NewTorrentFromVirtualFS(fsx.DirVirtual(seedir)),
		)
		md, err := torrent.NewFromInfo(mi, torrent.OptionStorage(storage.NewFile(filepath.Join(seedir))))
		require.NoError(t, err)
		seederTorrent, _, err := seeder.Start(md)
		require.NoError(t, err)
		defer seeder.Close()
		require.NoError(t, torrent.Verify(ctx, seederTorrent))
		_, err = torrent.DownloadInto(ctx, io.Discard, seederTorrent, torrent.TuneSeeding)
		require.NoError(t, err)

		root := fsx.DirVirtual(t.TempDir())
		require.NoError(t, fsx.MkDirs(0700, root.Path("torrent"), root.Path("media")))

		leecher := torrenttestx.Client(
			t,
			autobind.NewLoopback(autobind.EnableDHT(torrenttestx.QuickDHT(t))),
			torrent.NewMetadataCache(root.Path("torrent")),
			blockcache.NewTorrentFromVirtualFS(fsx.DirVirtual(root.Path("torrent"))),
		)
		defer leecher.Close()

		lmd := tracking.NewMetadata(
			new(md.ID),
			tracking.MetadataOptionFromInfo(mi),
			tracking.MetadataOptionAutoDescription,
		)
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		ltor, added, err := leecher.MaybeStart(torrent.NewFromInfo(mi))
		require.NoError(t, err)
		assert.True(t, added)
		require.NoError(t, ltor.Tune(torrent.TuneClientPeer(seeder)))

		require.NoError(t, tracking.DownloadInto(t.Context(), q, root, library.QueryCleanerNoop(), &lmd, ltor, io.Discard, asyncx.NewWakeup(t.Context())))

		var libMDs []library.Metadata
		require.NoError(t, sqlx.ScanInto(library.MetadataSearch(ctx, q,
			library.MetadataSearchBuilder().Where(library.MetadataQueryByTorrentID(lmd.ID)),
		), &libMDs))
		require.Len(t, libMDs, 3)

		mediavfs := fsx.DirVirtual(root.Path("media"))

		require.NoError(t, tracking.Reset(ctx, q, root, &lmd))

		for _, m := range libMDs {
			_, err := os.Stat(mediavfs.Path(m.ID))
			assert.True(t, os.IsNotExist(err), "media dir %s should be removed", m.ID)
		}
	})

	t.Run("removes torrent files matching infohash pattern", func(t *testing.T) {
		ctx := t.Context()
		q := sqltestx.Metadatabase(t)

		vfs := fsx.DirVirtual(t.TempDir())
		require.NoError(t, fsx.MkDirs(0700, vfs.Path("media"), vfs.Path("torrent")))

		lmd := tracking.NewMetadata(
			new(int160.Random()),
			tracking.MetadataOptionDescription("test torrent"),
		)
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		torrentvfs := fsx.DirVirtual(vfs.Path("torrent"))
		infohashStr := int160.FromBytes(lmd.Infohash).String()
		torrentFile := torrentvfs.Path(infohashStr + tracking.TorrentSuffix)
		torrentDir := torrentvfs.Path(infohashStr)
		require.NoError(t, os.WriteFile(torrentFile, []byte("torrent data"), 0600))
		require.NoError(t, os.MkdirAll(torrentDir, 0700))

		require.NoError(t, tracking.Reset(ctx, q, vfs, &lmd))

		_, err := os.Stat(torrentFile)
		assert.True(t, os.IsNotExist(err), "torrent file should be removed")
		_, err = os.Stat(torrentDir)
		assert.True(t, os.IsNotExist(err), "torrent subdir should be removed")
	})
}
