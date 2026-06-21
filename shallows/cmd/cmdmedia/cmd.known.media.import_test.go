package cmdmedia

import (
	"bytes"
	"strings"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/slicesx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestKnownImportRun(t *testing.T) {
	cmd := knownimport{}

	t.Run("inserts records from JSONL", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var a, b library.Known
		require.NoError(t, testx.Fake(&a, library.KnownOptionTestDefaults))
		require.NoError(t, testx.Fake(&b, library.KnownOptionTestDefaults))

		var input bytes.Buffer
		require.NoError(t, jsonl.NewEncoder(&input).Encode(a, b))
		require.NoError(t, cmd.run(ctx, db, &input))
		require.Equal(t, 2, sqltestx.Count(t, db, `SELECT COUNT(*) FROM library_known_media`))
	})

	t.Run("sets AutoDescription from title, original title, and overview", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		known.Title = "My Title"
		known.OriginalTitle = "Original Title"
		known.Overview = "A brief overview"

		var input bytes.Buffer
		require.NoError(t, jsonl.NewEncoder(&input).Encode(known))
		require.NoError(t, cmd.run(ctx, db, &input))
		expected := strings.Join([]string{"My Title", "Original Title", "A brief overview"}, "\n")
		require.Equal(t, expected, sqltestx.String(t, db, `SELECT auto_description FROM library_known_media LIMIT 1`))
	})

	t.Run("handles empty input", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		require.NoError(t, cmd.run(ctx, db, &bytes.Buffer{}))
		require.Equal(t, 0, sqltestx.Count(t, db, `SELECT COUNT(*) FROM library_known_media`))
	})

	t.Run("inserts more than one batch", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		records := make([]library.Known, 150)
		for i := range records {
			require.NoError(t, testx.Fake(&records[i], library.KnownOptionTestDefaults, library.KnownOptionRandomID))
		}

		var input bytes.Buffer
		items := slicesx.MapTransform(func(v library.Known) any { return v }, records...)
		require.NoError(t, jsonl.NewEncoder(&input).Encode(items...))
		require.NoError(t, cmd.run(ctx, db, &input))
		require.Equal(t, 150, sqltestx.Count(t, db, `SELECT COUNT(*) FROM library_known_media`))
	})

	t.Run("upserts on conflict", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))

		var first, second bytes.Buffer
		require.NoError(t, jsonl.NewEncoder(&first).Encode(known))
		require.NoError(t, jsonl.NewEncoder(&second).Encode(known))
		require.NoError(t, cmd.run(ctx, db, &first))
		require.NoError(t, cmd.run(ctx, db, &second))
		require.Equal(t, 1, sqltestx.Count(t, db, `SELECT COUNT(*) FROM library_known_media`))
		require.Equal(t, 1, sqltestx.Count(t, db, `SELECT duplicates FROM library_known_media WHERE uid = '`+known.UID+`'`))
	})
}
