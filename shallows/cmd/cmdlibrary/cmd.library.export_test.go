package cmdlibrary

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/blockcache"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/cryptox"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func writeBlockData(t *testing.T, vfs fsx.Virtual, md library.Metadata) {
	t.Helper()
	dcache, err := blockcache.NewDirectoryCache(vfs.Path(md.ID))
	require.NoError(t, err)
	src := io.LimitReader(cryptox.NewChaCha8(md.ID), int64(md.Bytes))
	_, err = io.Copy(io.NewOffsetWriter(dcache, int64(md.DiskOffset)), src)
	require.NoError(t, err)
}

func decodeExportStream(t *testing.T, r io.Reader) []exportTrailer {
	t.Helper()
	dec := json.NewDecoder(r)
	var trailers []exportTrailer
	for {
		var hdr exportHeader
		err := dec.Decode(&hdr)
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		for i := uint64(0); i < hdr.Chunks; i++ {
			var chunk exportChunk
			require.NoError(t, dec.Decode(&chunk))
		}
		var trailer exportTrailer
		require.NoError(t, dec.Decode(&trailer))
		trailers = append(trailers, trailer)
	}
	return trailers
}

func TestExportJSONLRun(t *testing.T) {
	cmd := exportJSONL{}

	t.Run("exports all non-tombstoned records", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)
		vfs := fsx.DirVirtual(t.TempDir())

		var a, b library.Metadata
		require.NoError(t, testx.Fake(&a, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
		require.NoError(t, testx.Fake(&b, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, a).Scan(&a))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, b).Scan(&b))
		writeBlockData(t, vfs, a)
		writeBlockData(t, vfs, b)

		var buf bytes.Buffer
		require.NoError(t, cmd.run(ctx, db, vfs, &buf))

		trailers := decodeExportStream(t, &buf)
		require.Len(t, trailers, 2)
	})

	t.Run("skips tombstoned records", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)
		vfs := fsx.DirVirtual(t.TempDir())

		var a, b library.Metadata
		require.NoError(t, testx.Fake(&a, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
		require.NoError(t, testx.Fake(&b, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, a).Scan(&a))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, b).Scan(&b))
		writeBlockData(t, vfs, a)
		_, err := db.ExecContext(ctx, `UPDATE library_metadata SET tombstoned_at = NOW() WHERE id = '`+b.ID+`'`)
		require.NoError(t, err)

		var buf bytes.Buffer
		require.NoError(t, cmd.run(ctx, db, vfs, &buf))

		trailers := decodeExportStream(t, &buf)
		require.Len(t, trailers, 1)
		require.Equal(t, a.ID, trailers[0].Metadata.ID)
	})

	t.Run("filter known-media true", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)
		vfs := fsx.DirVirtual(t.TempDir())

		realKnownID := uuid.Must(uuid.NewV4()).String()
		var withKnown, withoutKnown library.Metadata
		require.NoError(t, testx.Fake(&withKnown, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID, library.MetadataOptionKnownMediaID(realKnownID)))
		require.NoError(t, testx.Fake(&withoutKnown, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, withKnown).Scan(&withKnown))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, withoutKnown).Scan(&withoutKnown))
		writeBlockData(t, vfs, withKnown)
		writeBlockData(t, vfs, withoutKnown)

		hasKnown := true
		c := exportJSONL{KnownMedia: &hasKnown}
		var buf bytes.Buffer
		require.NoError(t, c.run(ctx, db, vfs, &buf))

		trailers := decodeExportStream(t, &buf)
		require.Len(t, trailers, 1)
		require.Equal(t, withKnown.ID, trailers[0].Metadata.ID)
	})

	t.Run("filter known-media false", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)
		vfs := fsx.DirVirtual(t.TempDir())

		realKnownID := uuid.Must(uuid.NewV4()).String()
		var withKnown, withoutKnown library.Metadata
		require.NoError(t, testx.Fake(&withKnown, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID, library.MetadataOptionKnownMediaID(realKnownID)))
		require.NoError(t, testx.Fake(&withoutKnown, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, withKnown).Scan(&withKnown))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, withoutKnown).Scan(&withoutKnown))
		writeBlockData(t, vfs, withKnown)
		writeBlockData(t, vfs, withoutKnown)

		hasKnown := false
		c := exportJSONL{KnownMedia: &hasKnown}
		var buf bytes.Buffer
		require.NoError(t, c.run(ctx, db, vfs, &buf))

		trailers := decodeExportStream(t, &buf)
		require.Len(t, trailers, 1)
		require.Equal(t, withoutKnown.ID, trailers[0].Metadata.ID)
	})

	t.Run("filter by id", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)
		vfs := fsx.DirVirtual(t.TempDir())

		var a, b library.Metadata
		require.NoError(t, testx.Fake(&a, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
		require.NoError(t, testx.Fake(&b, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, a).Scan(&a))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, b).Scan(&b))
		writeBlockData(t, vfs, a)
		writeBlockData(t, vfs, b)

		c := exportJSONL{ID: []string{a.ID}}
		var buf bytes.Buffer
		require.NoError(t, c.run(ctx, db, vfs, &buf))

		trailers := decodeExportStream(t, &buf)
		require.Len(t, trailers, 1)
		require.Equal(t, a.ID, trailers[0].Metadata.ID)
	})

	t.Run("filter torrent true", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)
		vfs := fsx.DirVirtual(t.TempDir())

		torrentID := uuid.Must(uuid.NewV4()).String()
		var withTorrent, withoutTorrent library.Metadata
		require.NoError(t, testx.Fake(&withTorrent, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID, library.MetadataOptionTorrentID(torrentID)))
		require.NoError(t, testx.Fake(&withoutTorrent, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, withTorrent).Scan(&withTorrent))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, withoutTorrent).Scan(&withoutTorrent))
		writeBlockData(t, vfs, withTorrent)
		writeBlockData(t, vfs, withoutTorrent)

		hasTorrent := true
		c := exportJSONL{Torrent: &hasTorrent}
		var buf bytes.Buffer
		require.NoError(t, c.run(ctx, db, vfs, &buf))

		trailers := decodeExportStream(t, &buf)
		require.Len(t, trailers, 1)
		require.Equal(t, withTorrent.ID, trailers[0].Metadata.ID)
	})

	t.Run("errors when blocks unavailable", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)
		vfs := fsx.DirVirtual(t.TempDir())

		var md library.Metadata
		require.NoError(t, testx.Fake(&md, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, md).Scan(&md))
		// intentionally no writeBlockData — blocks don't exist

		var buf bytes.Buffer
		require.Error(t, cmd.run(ctx, db, vfs, &buf))
	})

	t.Run("md5 is correct in trailer", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)
		vfs := fsx.DirVirtual(t.TempDir())

		var md library.Metadata
		require.NoError(t, testx.Fake(&md, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, md).Scan(&md))
		writeBlockData(t, vfs, md)

		var buf bytes.Buffer
		require.NoError(t, cmd.run(ctx, db, vfs, &buf))

		trailers := decodeExportStream(t, &buf)
		require.Len(t, trailers, 1)
		require.NotEmpty(t, trailers[0].MD5)
	})

	t.Run("skips folders", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)
		vfs := fsx.DirVirtual(t.TempDir())

		var md, dir library.Metadata
		require.NoError(t, testx.Fake(&md, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, md).Scan(&md))
		writeBlockData(t, vfs, md)

		require.NoError(t, testx.Fake(&dir, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID, library.MetadataOptionMimetype(mimex.Directory), library.MetadataOptionBytes(0)))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, dir).Scan(&dir))

		var buf bytes.Buffer
		require.NoError(t, cmd.run(ctx, db, vfs, &buf))

		trailers := decodeExportStream(t, &buf)
		require.Len(t, trailers, 1)
		require.Equal(t, md.ID, trailers[0].Metadata.ID)

		// export opens a block cache per row, and that call creates the directory it is
		// given, so a folder row leaves an empty tree behind on every export.
		require.NoDirExists(t, vfs.Path(dir.ID))
	})
}
