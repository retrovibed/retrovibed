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
	"github.com/retrovibed/retrovibed/shallows/cmd/retrovibed/daemons"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/tarx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestMediaMetadataImport(t *testing.T) {
	knownArchive := func(t *testing.T, records ...library.Known) []byte {
		t.Helper()

		dir := t.TempDir()
		f, err := os.Create(filepath.Join(dir, "media.jsonl"))
		require.NoError(t, err)

		enc := jsonl.NewEncoder(f)
		for _, r := range records {
			require.NoError(t, enc.Encode(r))
		}
		require.NoError(t, f.Close())

		var buf bytes.Buffer
		require.NoError(t, tarx.Pack(&buf, dir))
		return buf.Bytes()
	}

	t.Run("imports library.Known records from a completed metadata archive", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)

		seedir := t.TempDir()
		tvfs := fsx.DirVirtual(seedir)
		tstore := blockcache.NewTorrentFromVirtualFS(tvfs)
		mdcache := torrent.NewMetadataCache(seedir)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))

		archiveBytes := knownArchive(t, known)

		contentdir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(contentdir, "archive.jsonl.tar.gz"), archiveBytes, 0600))

		info, err := metainfo.NewFromPath(contentdir)
		require.NoError(t, err)

		md, err := torrent.NewFromInfo(info, torrent.OptionStorage(storage.NewFile(contentdir)))
		require.NoError(t, err)
		require.NoError(t, mdcache.Write(md))

		storedir := storage.InfoHashPathMaker(seedir, md.ID, info, nil)
		cache, err := blockcache.NewDirectoryCache(storedir)
		require.NoError(t, err)
		_, err = cache.WriteAt(archiveBytes, 0)
		require.NoError(t, err)

		lmd := tracking.NewMetadata(
			langx.Autoptr(md.ID),
			tracking.MetadataOptionFromInfo(info),
			tracking.MetadataOptionCompleted,
			tracking.MetadataOptionMimetype(mimex.RetrovibedMediaArchive),
		)
		require.NoError(t, tracking.MetadataInsertWithDefaults(t.Context(), q, lmd).Scan(&lmd))

		require.NoError(t, daemons.MediaMetadataImport(t.Context(), q, tvfs, tstore))

		require.Equal(t, 1, errorsx.Zero(sqlx.Count(t.Context(), q, "SELECT COUNT(*) FROM library_known_media")))
		require.Equal(t, 1, errorsx.Zero(sqlx.Count(t.Context(), q, "SELECT COUNT(*) FROM torrents_metadata WHERE imported_at < NOW()")))
	})

	t.Run("imports multiple library.Known records from a single archive", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)

		seedir := t.TempDir()
		tvfs := fsx.DirVirtual(seedir)
		tstore := blockcache.NewTorrentFromVirtualFS(tvfs)
		mdcache := torrent.NewMetadataCache(seedir)

		var known1, known2 library.Known
		require.NoError(t, testx.Fake(&known1, library.KnownOptionTestDefaults))
		require.NoError(t, testx.Fake(&known2, library.KnownOptionTestDefaults))

		archiveBytes := knownArchive(t, known1, known2)

		contentdir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(contentdir, "archive.jsonl.tar.gz"), archiveBytes, 0600))

		info, err := metainfo.NewFromPath(contentdir)
		require.NoError(t, err)

		md, err := torrent.NewFromInfo(info, torrent.OptionStorage(storage.NewFile(contentdir)))
		require.NoError(t, err)
		require.NoError(t, mdcache.Write(md))

		storedir := storage.InfoHashPathMaker(seedir, md.ID, info, nil)
		cache, err := blockcache.NewDirectoryCache(storedir)
		require.NoError(t, err)
		_, err = cache.WriteAt(archiveBytes, 0)
		require.NoError(t, err)

		lmd := tracking.NewMetadata(
			langx.Autoptr(md.ID),
			tracking.MetadataOptionFromInfo(info),
			tracking.MetadataOptionCompleted,
			tracking.MetadataOptionMimetype(mimex.RetrovibedMediaArchive),
		)
		require.NoError(t, tracking.MetadataInsertWithDefaults(t.Context(), q, lmd).Scan(&lmd))

		require.NoError(t, daemons.MediaMetadataImport(t.Context(), q, tvfs, tstore))

		require.Equal(t, 2, errorsx.Zero(sqlx.Count(t.Context(), q, "SELECT COUNT(*) FROM library_known_media")))
	})

	t.Run("skips archives not yet completed", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)

		seedir := t.TempDir()
		tvfs := fsx.DirVirtual(seedir)
		tstore := blockcache.NewTorrentFromVirtualFS(tvfs)

		lmd := tracking.NewMetadata(
			langx.Autoptr(int160.Random()),
			tracking.MetadataOptionMimetype(mimex.RetrovibedMediaArchive),
			// no MetadataOptionCompleted — completed_at remains at infinity
		)
		require.NoError(t, tracking.MetadataInsertWithDefaults(t.Context(), q, lmd).Scan(&lmd))

		require.NoError(t, daemons.MediaMetadataImport(t.Context(), q, tvfs, tstore))

		require.Equal(t, 0, errorsx.Zero(sqlx.Count(t.Context(), q, "SELECT COUNT(*) FROM library_known_media")))
	})

	t.Run("skips archives already imported", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)

		seedir := t.TempDir()
		tvfs := fsx.DirVirtual(seedir)
		tstore := blockcache.NewTorrentFromVirtualFS(tvfs)

		lmd := tracking.NewMetadata(
			langx.Autoptr(int160.Random()),
			tracking.MetadataOptionMimetype(mimex.RetrovibedMediaArchive),
			tracking.MetadataOptionCompleted,
		)
		require.NoError(t, tracking.MetadataInsertWithDefaults(t.Context(), q, lmd).Scan(&lmd))
		require.NoError(t, tracking.MetadataImportedByID(t.Context(), q, lmd.ID).Scan(&lmd))

		require.NoError(t, daemons.MediaMetadataImport(t.Context(), q, tvfs, tstore))

		require.Equal(t, 0, errorsx.Zero(sqlx.Count(t.Context(), q, "SELECT COUNT(*) FROM library_known_media")))
	})

	t.Run("no archive metadata returns no error", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)

		seedir := t.TempDir()
		tvfs := fsx.DirVirtual(seedir)
		tstore := blockcache.NewTorrentFromVirtualFS(tvfs)

		require.NoError(t, daemons.MediaMetadataImport(t.Context(), q, tvfs, tstore))
	})
}
