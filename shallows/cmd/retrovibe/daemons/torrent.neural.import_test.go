package daemons_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/james-lawrence/torrent/storage"
	"github.com/retrovibed/retrovibed/retroapi/blockcache"
	"github.com/retrovibed/retrovibed/shallows/cmd/retrovibe/daemons"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestNeuralImport(t *testing.T) {
	neuralWeights := func(t *testing.T, n int) []byte {
		t.Helper()
		return bytes.Repeat([]byte{0x42}, n)
	}

	t.Run("imports neural weights from a completed neural torrent", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)

		seedir := t.TempDir()
		tvfs := fsx.DirVirtual(seedir)
		tstore := blockcache.NewTorrentFromVirtualFS(tvfs)
		mdcache := torrent.NewMetadataCache(seedir)
		neuraldir := t.TempDir()

		content := neuralWeights(t, 4096)

		contentdir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(contentdir, "weights.bin"), content, 0600))

		info, err := metainfo.NewFromPath(contentdir)
		require.NoError(t, err)

		md, err := torrent.NewFromInfo(info, torrent.OptionStorage(storage.NewFile(contentdir)))
		require.NoError(t, err)
		require.NoError(t, mdcache.Write(md))

		storedir := storage.InfoHashPathMaker(seedir, md.ID, info, nil)
		cache, err := blockcache.NewDirectoryCache(storedir)
		require.NoError(t, err)
		_, err = cache.WriteAt(content, 0)
		require.NoError(t, err)

		lmd := tracking.NewMetadata(
			langx.Autoptr(md.ID),
			tracking.MetadataOptionFromInfo(info),
			tracking.MetadataOptionMimetype(mimex.RetrovibedNeural),
			tracking.MetadataOptionCompleted,
			tracking.MetadataOptionDescription("weights.bin"),
		)
		require.NoError(t, tracking.MetadataInsertWithDefaults(t.Context(), q, lmd).Scan(&lmd))

		require.NoError(t, daemons.NeuralImport(t.Context(), q, neuraldir, tvfs, tstore))

		written, err := os.ReadFile(filepath.Join(neuraldir, "weights.bin"))
		require.NoError(t, err)
		require.Equal(t, content, written)
		require.Equal(t, 1, errorsx.Zero(sqlx.Count(t.Context(), q, "SELECT COUNT(*) FROM torrents_metadata WHERE imported_at < NOW()")))
	})

	t.Run("imports multiple neural weights", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)

		seedir := t.TempDir()
		tvfs := fsx.DirVirtual(seedir)
		tstore := blockcache.NewTorrentFromVirtualFS(tvfs)
		mdcache := torrent.NewMetadataCache(seedir)
		neuraldir := t.TempDir()

		setupNeural := func(t *testing.T, filename string, content []byte) {
			t.Helper()

			contentdir := t.TempDir()
			require.NoError(t, os.WriteFile(filepath.Join(contentdir, filename), content, 0600))

			info, err := metainfo.NewFromPath(contentdir)
			require.NoError(t, err)

			md, err := torrent.NewFromInfo(info, torrent.OptionStorage(storage.NewFile(contentdir)))
			require.NoError(t, err)
			require.NoError(t, mdcache.Write(md))

			storedir := storage.InfoHashPathMaker(seedir, md.ID, info, nil)
			cache, err := blockcache.NewDirectoryCache(storedir)
			require.NoError(t, err)
			_, err = cache.WriteAt(content, 0)
			require.NoError(t, err)

			lmd := tracking.NewMetadata(
				langx.Autoptr(md.ID),
				tracking.MetadataOptionFromInfo(info),
				tracking.MetadataOptionMimetype(mimex.RetrovibedNeural),
				tracking.MetadataOptionCompleted,
				tracking.MetadataOptionDescription(filename),
			)
			require.NoError(t, tracking.MetadataInsertWithDefaults(t.Context(), q, lmd).Scan(&lmd))
		}

		content1 := neuralWeights(t, 4096)
		content2 := bytes.Repeat([]byte{0x24}, 8192)
		setupNeural(t, "weights1.bin", content1)
		setupNeural(t, "weights2.bin", content2)

		require.NoError(t, daemons.NeuralImport(t.Context(), q, neuraldir, tvfs, tstore))

		written1, err := os.ReadFile(filepath.Join(neuraldir, "weights1.bin"))
		require.NoError(t, err)
		require.Equal(t, content1, written1)

		written2, err := os.ReadFile(filepath.Join(neuraldir, "weights2.bin"))
		require.NoError(t, err)
		require.Equal(t, content2, written2)

		require.Equal(t, 2, errorsx.Zero(sqlx.Count(t.Context(), q, "SELECT COUNT(*) FROM torrents_metadata WHERE imported_at < NOW()")))
	})

	t.Run("skips neural not yet completed", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)

		seedir := t.TempDir()
		tvfs := fsx.DirVirtual(seedir)
		tstore := blockcache.NewTorrentFromVirtualFS(tvfs)
		neuraldir := t.TempDir()

		lmd := tracking.NewMetadata(
			langx.Autoptr(int160.Random()),
			tracking.MetadataOptionMimetype(mimex.RetrovibedNeural),
			tracking.MetadataOptionDescription("weights.bin"),
			// no MetadataOptionCompleted — completed_at remains at infinity
		)
		require.NoError(t, tracking.MetadataInsertWithDefaults(t.Context(), q, lmd).Scan(&lmd))

		require.NoError(t, daemons.NeuralImport(t.Context(), q, neuraldir, tvfs, tstore))

		_, err := os.Stat(filepath.Join(neuraldir, "weights.bin"))
		require.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("skips neural already imported", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)

		seedir := t.TempDir()
		tvfs := fsx.DirVirtual(seedir)
		tstore := blockcache.NewTorrentFromVirtualFS(tvfs)
		neuraldir := t.TempDir()

		lmd := tracking.NewMetadata(
			langx.Autoptr(int160.Random()),
			tracking.MetadataOptionMimetype(mimex.RetrovibedNeural),
			tracking.MetadataOptionCompleted,
			tracking.MetadataOptionDescription("weights.bin"),
		)
		require.NoError(t, tracking.MetadataInsertWithDefaults(t.Context(), q, lmd).Scan(&lmd))
		require.NoError(t, tracking.MetadataImportedByID(t.Context(), q, lmd.ID).Scan(&lmd))

		require.NoError(t, daemons.NeuralImport(t.Context(), q, neuraldir, tvfs, tstore))

		_, err := os.Stat(filepath.Join(neuraldir, "weights.bin"))
		require.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("no neural metadata returns no error", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)

		seedir := t.TempDir()
		tvfs := fsx.DirVirtual(seedir)
		tstore := blockcache.NewTorrentFromVirtualFS(tvfs)
		neuraldir := t.TempDir()

		require.NoError(t, daemons.NeuralImport(t.Context(), q, neuraldir, tvfs, tstore))
	})
}
