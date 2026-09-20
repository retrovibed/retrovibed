package tracking_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/autobind"
	"github.com/james-lawrence/torrent/storage"
	"github.com/james-lawrence/torrent/torrenttest"
	"github.com/james-lawrence/torrent/torrenttestx"
	"github.com/retrovibed/retrovibed/retroapi/blockcache"
	"github.com/retrovibed/retrovibed/retroapi/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResume(t *testing.T) {
	t.Run("runs the download to completion when the torrent is newly added", func(t *testing.T) {
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
		defer seeder.Close()

		smd, err := torrent.NewFromInfo(mi, torrent.OptionStorage(storage.NewFile(filepath.Join(seedir))))
		require.NoError(t, err)
		seederTorrent, _, err := seeder.Start(smd)
		require.NoError(t, err)
		require.NoError(t, torrent.Verify(ctx, seederTorrent))

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
			new(smd.ID),
			tracking.MetadataOptionFromInfo(mi),
			tracking.MetadataOptionAutoDescription,
		)
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		_, added, err := tracking.Resume(ctx, q, root, library.QueryCleanerNoop(), leecher, blockcache.NewTorrentFromVirtualFS(fsx.DirVirtual(leechdir)), lmd, asyncx.NewWakeup(ctx), torrent.TuneClientPeer(seeder), torrent.TuneNewConns)
		require.NoError(t, err)
		require.True(t, added)

		var actual tracking.Metadata
		require.Eventually(t, func() bool {
			require.NoError(t, tracking.MetadataFindByID(ctx, q, lmd.ID).Scan(&actual))
			return actual.CompletedAt.Before(time.Now())
		}, 15*time.Second, 100*time.Millisecond)
	})

	t.Run("does not run the download again when the torrent is already running", func(t *testing.T) {
		ctx := t.Context()
		q := sqltestx.Metadatabase(t)

		seedir := t.TempDir()
		mi, err := torrenttest.RandomMulti(seedir, 5, 16*bytesx.KiB, 64*bytesx.KiB)
		require.NoError(t, err)

		root := fsx.DirVirtual(t.TempDir())
		leechdir := root.Path("torrent")
		require.NoError(t, fsx.MkDirs(0700, leechdir, root.Path("media")))

		leecher := torrenttestx.Client(
			t,
			autobind.NewLoopback(
				autobind.EnableDHT(torrenttestx.QuickDHT(t)),
			),
			torrent.NewMetadataCache(leechdir),
			blockcache.NewTorrentFromVirtualFS(fsx.DirVirtual(leechdir)),
		)
		defer leecher.Close()

		imd, err := torrent.NewFromInfo(mi)
		require.NoError(t, err)

		lmd := tracking.NewMetadata(
			new(imd.ID),
			tracking.MetadataOptionFromInfo(mi),
			tracking.MetadataOptionAutoDescription,
		)
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		tstore := blockcache.NewTorrentFromVirtualFS(fsx.DirVirtual(leechdir))
		metadata, err := tracking.NewResumableMetadata(root, tstore, lmd)
		require.NoError(t, err)

		_, added, err := leecher.Start(metadata)
		require.NoError(t, err)
		require.True(t, added)

		_, added, err = tracking.Resume(ctx, q, root, library.QueryCleanerNoop(), leecher, tstore, lmd, asyncx.NewWakeup(ctx))
		require.NoError(t, err)
		require.False(t, added)
	})

	t.Run("records the bytes downloaded from peers when the download is not a resume", func(t *testing.T) {
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
		defer seeder.Close()

		smd, err := torrent.NewFromInfo(mi, torrent.OptionStorage(storage.NewFile(filepath.Join(seedir))))
		require.NoError(t, err)
		seederTorrent, _, err := seeder.Start(smd)
		require.NoError(t, err)
		require.NoError(t, torrent.Verify(ctx, seederTorrent))

		root := fsx.DirVirtual(t.TempDir())
		leechdir := root.Path("torrent")
		require.NoError(t, fsx.MkDirs(0700, leechdir, root.Path("media")))

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
			new(smd.ID),
			tracking.MetadataOptionFromInfo(mi),
			tracking.MetadataOptionAutoDescription,
		)
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		_, added, err := tracking.Resume(ctx, q, root, library.QueryCleanerNoop(), leecher, blockcache.NewTorrentFromVirtualFS(fsx.DirVirtual(leechdir)), lmd, asyncx.NewWakeup(ctx), torrent.TuneClientPeer(seeder), torrent.TuneNewConns)
		require.NoError(t, err)
		require.True(t, added)

		var actual tracking.Metadata
		require.Eventually(t, func() bool {
			require.NoError(t, tracking.MetadataFindByID(ctx, q, lmd.ID).Scan(&actual))
			return actual.CompletedAt.Before(time.Now())
		}, 15*time.Second, 100*time.Millisecond)

		require.EqualValues(t, mi.TotalLength(), actual.Downloaded)
		require.EqualValues(t, mi.TotalLength(), actual.Available)
	})

	t.Run("refreshes downloaded from the verified on-disk bytes and preserves uploaded when resuming", func(t *testing.T) {
		const previouslyUploaded = 7 * bytesx.KiB

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
		defer seeder.Close()

		smd, err := torrent.NewFromInfo(mi, torrent.OptionStorage(storage.NewFile(filepath.Join(seedir))))
		require.NoError(t, err)
		seederTorrent, _, err := seeder.Start(smd)
		require.NoError(t, err)
		require.NoError(t, torrent.Verify(ctx, seederTorrent))

		root := fsx.DirVirtual(t.TempDir())
		leechdir := root.Path("torrent")
		require.NoError(t, fsx.MkDirs(0700, leechdir, root.Path("media")))

		// first session: download the torrent from the seeder so the data is on disk.
		first := torrenttestx.Client(
			t,
			autobind.NewLoopback(
				autobind.EnableDHT(torrenttestx.QuickDHT(t)),
			),
			torrent.NewMetadataCache(leechdir),
			blockcache.NewTorrentFromVirtualFS(fsx.DirVirtual(leechdir)),
		)
		defer first.Close()

		lmd := tracking.NewMetadata(
			new(smd.ID),
			tracking.MetadataOptionFromInfo(mi),
			tracking.MetadataOptionAutoDescription,
		)
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		_, added, err := tracking.Resume(ctx, q, root, library.QueryCleanerNoop(), first, blockcache.NewTorrentFromVirtualFS(fsx.DirVirtual(leechdir)), lmd, asyncx.NewWakeup(ctx), torrent.TuneClientPeer(seeder), torrent.TuneNewConns)
		require.NoError(t, err)
		require.True(t, added)

		require.Eventually(t, func() bool {
			require.NoError(t, tracking.MetadataFindByID(ctx, q, lmd.ID).Scan(&lmd))
			return lmd.CompletedAt.Before(time.Now())
		}, 15*time.Second, 100*time.Millisecond)

		// restart: stop the first session, and put the row back to what an interrupted
		// session would have left behind: a stale downloaded value and an accumulated uploaded.
		require.NoError(t, first.Close())

		staleDownloaded := uint64(mi.TotalLength()) / 2
		require.NoError(t, tracking.MetadataResetByID(ctx, q, lmd.ID).Scan(&lmd))
		require.NoError(t, tracking.MetadataProgressByID(ctx, q, lmd.ID, 0, lmd.Bytes, staleDownloaded, lmd.Bytes).Scan(&lmd))
		require.NoError(t, tracking.MetadataUploadedByID(ctx, q, lmd.Infohash, previouslyUploaded).Scan(&lmd))
		require.Equal(t, staleDownloaded, lmd.Downloaded)
		require.EqualValues(t, previouslyUploaded, lmd.Uploaded)
		require.True(t, lmd.CompletedAt.After(time.Now()), "the restarted torrent must not be marked completed")

		// second session: same on disk data, no peers involved.
		second := torrenttestx.Client(
			t,
			autobind.NewLoopback(
				autobind.EnableDHT(torrenttestx.QuickDHT(t)),
			),
			torrent.NewMetadataCache(leechdir),
			blockcache.NewTorrentFromVirtualFS(fsx.DirVirtual(leechdir)),
		)
		defer second.Close()

		_, added, err = tracking.Resume(ctx, q, root, library.QueryCleanerNoop(), second, blockcache.NewTorrentFromVirtualFS(fsx.DirVirtual(leechdir)), lmd, asyncx.NewWakeup(ctx))
		require.NoError(t, err)
		require.True(t, added)

		var actual tracking.Metadata
		require.Eventually(t, func() bool {
			require.NoError(t, tracking.MetadataFindByID(ctx, q, lmd.ID).Scan(&actual))
			return actual.CompletedAt.Before(time.Now())
		}, 15*time.Second, 100*time.Millisecond)

		assert.EqualValues(t, mi.TotalLength(), actual.Available, "all bytes are on disk")
		assert.EqualValues(t, mi.TotalLength(), actual.Downloaded, "downloaded is the verified bytes on disk, the stale value from the interrupted session must be refreshed")
		assert.EqualValues(t, previouslyUploaded, actual.Uploaded, "uploaded is a running total across sessions and must not be replaced by the new session's count")
	})
}
