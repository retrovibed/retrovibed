package cmdlibrary

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/internal/timex"

	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
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

		dstDB := sqltestx.Metadatabase(t)
		c := newImportDirectoryServer(t, dstDB)
		require.NoError(t, cmd.run(ctx, c, "https://localhost", exportToBuffer()))
		require.NoError(t, cmd.run(ctx, c, "https://localhost", exportToBuffer()))
		require.Equal(t, 1, testx.Must(sqlx.Count(ctx, dstDB, "SELECT COUNT(*) FROM library_metadata"))(t))
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
