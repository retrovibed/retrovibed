package daemons_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/james-lawrence/torrent/storage"
	"github.com/retrovibed/retrovibed/retroapi/blockcache"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/cmd/retrovibe/daemons"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestPublishPluginTorrentImport(t *testing.T) {
	wasmModule := func(t *testing.T, n int) []byte {
		t.Helper()
		return append([]byte{0x00, 0x61, 0x73, 0x6D}, bytes.Repeat([]byte{0x42}, n)...)
	}

	notWasm := func(t *testing.T, n int) []byte {
		t.Helper()
		return bytes.Repeat([]byte{0x24}, n)
	}

	setup := func(t *testing.T, description string, content []byte, options ...func(*tracking.Metadata)) (sqlx.Queryer, fsx.Virtual, storage.ClientImpl, string, tracking.Metadata) {
		t.Helper()

		q := sqltestx.Metadatabase(t)

		seedir := t.TempDir()
		tvfs := fsx.DirVirtual(seedir)
		tstore := blockcache.NewTorrentFromVirtualFS(tvfs)
		mdcache := torrent.NewMetadataCache(seedir)
		plugindir := t.TempDir()

		contentdir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(contentdir, "plugin.wasm"), content, 0600))

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
			new(md.ID),
			append([]func(*tracking.Metadata){
				tracking.MetadataOptionFromInfo(info),
				tracking.MetadataOptionMimetype(mimex.RetrovibedPublishModule),
				tracking.MetadataOptionDescription(description),
			}, options...)...,
		)
		require.NoError(t, tracking.MetadataInsertWithDefaults(t.Context(), q, lmd).Scan(&lmd))

		return q, tvfs, tstore, plugindir, lmd
	}

	t.Run("imports a valid wasm module and marks it imported", func(t *testing.T) {
		content := wasmModule(t, 4096)
		q, tvfs, tstore, plugindir, lmd := setup(t, "lemmy", content, tracking.MetadataOptionCompleted)

		require.NoError(t, daemons.PublishPluginTorrentImport(t.Context(), q, plugindir, tvfs, tstore))

		written, err := os.ReadFile(filepath.Join(plugindir, "lemmy.wasm"))
		require.NoError(t, err)
		require.Equal(t, content, written)

		var updated tracking.Metadata
		require.NoError(t, tracking.MetadataFindByID(t.Context(), q, lmd.ID).Scan(&updated))
		require.NotEqual(t, timex.Inf(), updated.ImportedAt)
	})

	t.Run("sanitizes a name that would escape the plugin directory", func(t *testing.T) {
		content := wasmModule(t, 4096)
		q, tvfs, tstore, plugindir, _ := setup(t, "../../escape", content, tracking.MetadataOptionCompleted)

		require.NoError(t, daemons.PublishPluginTorrentImport(t.Context(), q, plugindir, tvfs, tstore))

		_, err := os.Stat(filepath.Join(filepath.Dir(filepath.Dir(plugindir)), "escape.wasm"))
		require.ErrorIs(t, err, os.ErrNotExist)

		written, err := os.ReadFile(filepath.Join(plugindir, "escape.wasm"))
		require.NoError(t, err)
		require.Equal(t, content, written)
	})

	t.Run("tombstones a name with nothing usable left in it", func(t *testing.T) {
		content := wasmModule(t, 4096)
		q, tvfs, tstore, plugindir, lmd := setup(t, "../..", content, tracking.MetadataOptionCompleted)

		require.NoError(t, daemons.PublishPluginTorrentImport(t.Context(), q, plugindir, tvfs, tstore))

		entries, err := os.ReadDir(plugindir)
		require.NoError(t, err)
		require.Empty(t, entries)

		var updated tracking.Metadata
		require.NoError(t, tracking.MetadataFindByID(t.Context(), q, lmd.ID).Scan(&updated))
		require.NotEqual(t, timex.Inf(), updated.TombstonedAt)
		require.Equal(t, timex.Inf(), updated.ImportedAt)
	})

	t.Run("tombstones and does not copy a payload that isn't a valid wasm module", func(t *testing.T) {
		content := notWasm(t, 4096)
		q, tvfs, tstore, plugindir, lmd := setup(t, "lemmy", content, tracking.MetadataOptionCompleted)

		require.NoError(t, daemons.PublishPluginTorrentImport(t.Context(), q, plugindir, tvfs, tstore))

		_, err := os.Stat(filepath.Join(plugindir, "lemmy.wasm"))
		require.ErrorIs(t, err, os.ErrNotExist)

		var updated tracking.Metadata
		require.NoError(t, tracking.MetadataFindByID(t.Context(), q, lmd.ID).Scan(&updated))
		require.NotEqual(t, timex.Inf(), updated.TombstonedAt)
		require.Equal(t, timex.Inf(), updated.ImportedAt)
	})

	t.Run("skips a plugin torrent not yet completed", func(t *testing.T) {
		content := wasmModule(t, 4096)
		q, tvfs, tstore, plugindir, _ := setup(t, "lemmy", content)

		require.NoError(t, daemons.PublishPluginTorrentImport(t.Context(), q, plugindir, tvfs, tstore))

		_, err := os.Stat(filepath.Join(plugindir, "lemmy.wasm"))
		require.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("skips a plugin torrent already imported", func(t *testing.T) {
		content := wasmModule(t, 4096)
		q, tvfs, tstore, plugindir, lmd := setup(t, "lemmy", content, tracking.MetadataOptionCompleted)
		require.NoError(t, tracking.MetadataImportedByID(t.Context(), q, lmd.ID).Scan(&lmd))

		require.NoError(t, daemons.PublishPluginTorrentImport(t.Context(), q, plugindir, tvfs, tstore))

		_, err := os.Stat(filepath.Join(plugindir, "lemmy.wasm"))
		require.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("skips a plugin torrent already tombstoned", func(t *testing.T) {
		content := notWasm(t, 4096)
		q, tvfs, tstore, plugindir, lmd := setup(t, "lemmy", content, tracking.MetadataOptionCompleted)
		require.NoError(t, tracking.MetadataTombstoneByID(t.Context(), q, lmd.ID).Scan(&lmd))

		require.NoError(t, daemons.PublishPluginTorrentImport(t.Context(), q, plugindir, tvfs, tstore))

		require.Equal(t, 0, errorsx.Zero(sqlx.Count(t.Context(), q, "SELECT COUNT(*) FROM torrents_metadata WHERE imported_at < NOW()")))
	})

	t.Run("no matching plugin metadata returns no error", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)

		seedir := t.TempDir()
		tvfs := fsx.DirVirtual(seedir)
		tstore := blockcache.NewTorrentFromVirtualFS(tvfs)
		plugindir := t.TempDir()

		require.NoError(t, daemons.PublishPluginTorrentImport(t.Context(), q, plugindir, tvfs, tstore))
	})
}
