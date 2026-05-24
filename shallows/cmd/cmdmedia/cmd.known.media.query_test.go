package cmdmedia

import (
	"bytes"
	"strings"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestKnownQueryRun(t *testing.T) {
	t.Run("empty database returns no error", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		require.NoError(t, knownquery{}.run(ctx, strings.NewReader("{\"query\":\"inception\"}\n"), db, library.QueryCleanerNoop()))
	})

	t.Run("finds record matching description", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		known.Title = "Inception"
		known.Overview = "A mind-bending thriller about dreams within dreams"

		require.NoError(t, knownimport{}.run(ctx, db, jsonlBuffer(t, known)))
		require.NoError(t, knownquery{}.run(ctx, strings.NewReader("{\"query\":\"Inception\"}\n"), db, library.QueryCleanerNoop()))
	})

	t.Run("falls back to title search when no fts match", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		known.Title = "The Dark Knight"

		require.NoError(t, knownimport{}.run(ctx, db, jsonlBuffer(t, known)))
		require.NoError(t, knownquery{}.run(ctx, strings.NewReader("{\"query\":\"Dark Knight\"}\n"), db, library.QueryCleanerNoop()))
	})

	t.Run("query with embedded lucene keywords does not error", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		require.NoError(t, knownquery{}.run(ctx, strings.NewReader("{\"query\":\"How to School 101 Brilliant Ideas to Keep\"}\n"), db, library.QueryCleanerNoop()))
	})

	t.Run("fails when stdin is empty", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		require.Error(t, knownquery{}.run(ctx, bytes.NewReader(nil), db, library.QueryCleanerNoop()))
	})
}
