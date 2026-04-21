package cmdmedia

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestKnownImportRun(t *testing.T) {
	cmd := knownimport{}

	jsonl := func(t *testing.T, records ...library.Known) *bytes.Buffer {
		t.Helper()
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		for _, r := range records {
			require.NoError(t, enc.Encode(r))
		}
		return &buf
	}

	t.Run("inserts records from JSONL", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var a, b library.Known
		require.NoError(t, testx.Fake(&a, library.KnownOptionTestDefaults))
		require.NoError(t, testx.Fake(&b, library.KnownOptionTestDefaults))

		require.NoError(t, cmd.run(ctx, db, jsonl(t, a, b)))
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

		require.NoError(t, cmd.run(ctx, db, jsonl(t, known)))
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

		require.NoError(t, cmd.run(ctx, db, jsonl(t, records...)))
		require.Equal(t, 150, sqltestx.Count(t, db, `SELECT COUNT(*) FROM library_known_media`))
	})

	t.Run("upserts on conflict", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))

		require.NoError(t, cmd.run(ctx, db, jsonl(t, known)))
		require.NoError(t, cmd.run(ctx, db, jsonl(t, known)))
		require.Equal(t, 1, sqltestx.Count(t, db, `SELECT COUNT(*) FROM library_known_media`))
		require.Equal(t, 1, sqltestx.Count(t, db, `SELECT duplicates FROM library_known_media WHERE uid = '`+known.UID+`'`))
	})
}
