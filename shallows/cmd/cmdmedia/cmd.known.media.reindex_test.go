package cmdmedia

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
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/tarx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestKnownReindex(t *testing.T) {
	t.Run("truncates the cache and reimports previously imported archives", func(t *testing.T) {
		db := sqltestx.Metadatabase(t)

		seedir := t.TempDir()
		tvfs := fsx.DirVirtual(seedir)
		tstore := blockcache.NewTorrentFromVirtualFS(tvfs)
		mdcache := torrent.NewMetadataCache(seedir)

		var known, stale library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		require.NoError(t, testx.Fake(&stale, library.KnownOptionTestDefaults))
		require.NoError(t, sqlx.Discard(sqlx.Scan(library.NewKnownBatchInsertWithDefaults(t.Context(), db, stale))))

		archivedir := t.TempDir()
		f, err := os.Create(filepath.Join(archivedir, "media.jsonl"))
		require.NoError(t, err)
		require.NoError(t, jsonl.NewEncoder(f).Encode(known))
		require.NoError(t, f.Close())

		var archive bytes.Buffer
		require.NoError(t, tarx.Pack(&archive, archivedir))

		contentdir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(contentdir, "archive.jsonl.tar.gz"), archive.Bytes(), 0600))

		info, err := metainfo.NewFromPath(contentdir)
		require.NoError(t, err)

		md, err := torrent.NewFromInfo(info, torrent.OptionStorage(storage.NewFile(contentdir)))
		require.NoError(t, err)
		require.NoError(t, mdcache.Write(md))

		cache, err := blockcache.NewDirectoryCache(storage.InfoHashPathMaker(seedir, md.ID, info, nil))
		require.NoError(t, err)
		_, err = cache.WriteAt(archive.Bytes(), 0)
		require.NoError(t, err)

		lmd := tracking.NewMetadata(
			new(md.ID),
			tracking.MetadataOptionFromInfo(info),
			tracking.MetadataOptionCompleted,
			tracking.MetadataOptionMimetype(mimex.RetrovibedMediaArchive),
		)
		require.NoError(t, tracking.MetadataInsertWithDefaults(t.Context(), db, lmd).Scan(&lmd))
		require.NoError(t, tracking.MetadataImportedByID(t.Context(), db, lmd.ID).Scan(&lmd))

		require.NoError(t, knownreindex{}.run(t.Context(), db, tvfs, tstore))

		require.Equal(t, 1, sqltestx.Count(t, db, `SELECT COUNT(*) FROM cache.library_known_media`))
		require.Equal(t, known.UID, sqltestx.String(t, db, `SELECT uid::VARCHAR FROM cache.library_known_media`))
		require.Equal(t, 1, sqltestx.Count(t, db, `SELECT COUNT(*) FROM torrents_metadata WHERE imported_at < 'infinity'`))
	})
}
