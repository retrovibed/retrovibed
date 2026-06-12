package cmdlibrary

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/internal/timex"

	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/cryptox"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestImportJSONLRun(t *testing.T) {
	cmd := importJSONL{}

	t.Run("handles empty input", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		q := sqltestx.Metadatabase(t)
		c := newImportDirectoryServer(t, q)

		require.NoError(t, cmd.run(ctx, c, "https://localhost", &bytes.Buffer{}))
		require.Equal(t, 0, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM library_metadata"))(t))
	})

	t.Run("round-trips export", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		srcDB := sqltestx.Metadatabase(t)
		srcVFS := fsx.DirVirtual(t.TempDir())

		var a, b library.Metadata
		require.NoError(t, testx.Fake(&a, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
		require.NoError(t, testx.Fake(&b, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, srcDB, a).Scan(&a))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, srcDB, b).Scan(&b))
		writeBlockData(t, srcVFS, a)
		writeBlockData(t, srcVFS, b)

		var buf bytes.Buffer
		exp := exportJSONL{}
		require.NoError(t, exp.run(ctx, srcDB, srcVFS, &buf))

		dstDB := sqltestx.Metadatabase(t)
		c := newImportDirectoryServer(t, dstDB)
		require.NoError(t, cmd.run(ctx, c, "https://localhost", &buf))

		require.Equal(t, 2, testx.Must(sqlx.Count(ctx, dstDB, "SELECT COUNT(*) FROM library_metadata"))(t))
	})

	t.Run("round-trips export with multi-chunk data", func(t *testing.T) {
		const chunksize = 16 * bytesx.MiB
		const blocklength = 32 * bytesx.MiB

		sizes := map[string]uint64{
			"exactly one chunk":            chunksize,
			"one byte under one chunk":     chunksize - 1,
			"exactly two chunks":           2 * chunksize,
			"exactly three chunks":         3 * chunksize,
			"one byte over a block":        blocklength + 1,
			"multiple chunks across block": 40 * bytesx.MiB,
		}

		for name, size := range sizes {
			t.Run(name, func(t *testing.T) {
				ctx, done := testx.Context(t)
				defer done()
				srcDB := sqltestx.Metadatabase(t)
				srcVFS := fsx.DirVirtual(t.TempDir())

				var md library.Metadata
				require.NoError(t, testx.Fake(&md, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID, library.MetadataOptionBytes(size)))
				require.NoError(t, library.MetadataInsertWithDefaults(ctx, srcDB, md).Scan(&md))
				writeBlockData(t, srcVFS, md)

				expected := testx.IOMD5(io.LimitReader(cryptox.NewChaCha8(md.ID), int64(md.Bytes)))

				var buf bytes.Buffer
				exp := exportJSONL{}
				require.NoError(t, exp.run(ctx, srcDB, srcVFS, &buf))

				dstDB := sqltestx.Metadatabase(t)
				c := newImportDirectoryServer(t, dstDB)
				require.NoError(t, cmd.run(ctx, c, "https://localhost", &buf))

				require.Equal(t, 1, testx.Must(sqlx.Count(ctx, dstDB, "SELECT COUNT(*) FROM library_metadata"))(t))

				var dstMD library.Metadata
				require.NoError(t, library.MetadataFindByID(ctx, dstDB, expected).Scan(&dstMD))
				require.Equal(t, md.Bytes, dstMD.Bytes)
			})
		}
	})

	t.Run("upserts on conflict", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		srcDB := sqltestx.Metadatabase(t)
		srcVFS := fsx.DirVirtual(t.TempDir())

		var md library.Metadata
		require.NoError(t, testx.Fake(&md, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, srcDB, md).Scan(&md))
		writeBlockData(t, srcVFS, md)

		exportToBuffer := func() *bytes.Buffer {
			var buf bytes.Buffer
			exp := exportJSONL{}
			require.NoError(t, exp.run(ctx, srcDB, srcVFS, &buf))
			return &buf
		}

		expectedID := testx.IOMD5(io.LimitReader(cryptox.NewChaCha8(md.ID), int64(md.Bytes)))

		dstDB := sqltestx.Metadatabase(t)
		c := newImportDirectoryServer(t, dstDB)
		require.NoError(t, cmd.run(ctx, c, "https://localhost", exportToBuffer()))

		var first library.Metadata
		require.NoError(t, library.MetadataFindByID(ctx, dstDB, expectedID).Scan(&first))

		require.NoError(t, cmd.run(ctx, c, "https://localhost", exportToBuffer()))
		require.Equal(t, 1, testx.Must(sqlx.Count(ctx, dstDB, "SELECT COUNT(*) FROM library_metadata"))(t))

		var second library.Metadata
		require.NoError(t, library.MetadataFindByID(ctx, dstDB, expectedID).Scan(&second))

		// fields synced from the source record on every import.
		require.Equal(t, md.Description, second.Description)
		require.Equal(t, md.KnownMediaID, second.KnownMediaID)
		require.Equal(t, md.ArchiveID, second.ArchiveID)
		require.Equal(t, md.Mimetype, second.Mimetype)
		require.Equal(t, md.Bytes, second.Bytes)

		// re-importing the same export is idempotent: the second import
		// leaves these fields unchanged from the first.
		require.Equal(t, first.Description, second.Description)
		require.Equal(t, first.AutoDescription, second.AutoDescription)
		require.Equal(t, first.KnownMediaID, second.KnownMediaID)
		require.Equal(t, first.ArchiveID, second.ArchiveID)
		require.Equal(t, first.Bytes, second.Bytes)
		require.Equal(t, first.Mimetype, second.Mimetype)
		require.Equal(t, first.DiskOffset, second.DiskOffset)
		require.Equal(t, first.TorrentID, second.TorrentID)
	})

	t.Run("rejects bad MD5", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)

		var md library.Metadata
		require.NoError(t, testx.Fake(&md, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
		md.Bytes = 0
		timex.JSONSafeEncodeOption(&md)
		require.NoError(t, enc.Encode(exportHeader{Chunks: 0}))
		require.NoError(t, enc.Encode(exportTrailer{Metadata: md, MD5: "nottherighthash00000000000000000"}))

		dstDB := sqltestx.Metadatabase(t)
		c := newImportDirectoryServer(t, dstDB)
		err := cmd.run(ctx, c, "https://localhost", &buf)
		require.Error(t, err)
		require.Contains(t, err.Error(), "MD5 mismatch")
	})
}
