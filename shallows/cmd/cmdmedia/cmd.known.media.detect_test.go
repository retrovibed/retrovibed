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

func TestKnownDetectRun(t *testing.T) {
	t.Run("finds a matching record by title", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		known.Title = "The Grand Budapest Hotel"

		require.NoError(t, knownimport{}.run(ctx, db, jsonlBuffer(t, known)))

		require.NoError(t, knowndetect{}.run(ctx, strings.NewReader("{\"query\":\"The Grand Budapest Hotel\"}\n"), db, library.NoopQueryCleaner{}))
	})

	t.Run("returns no error when no match found", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		require.NoError(t, knowndetect{}.run(ctx, strings.NewReader("{\"query\":\"something completely unknown\"}\n"), db, library.NoopQueryCleaner{}))
	})

	t.Run("query with embedded lucene keywords does not error", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		require.NoError(t, knowndetect{}.run(ctx, strings.NewReader("{\"query\":\"How to School 101 Brilliant Ideas to Keep\"}\n"), db, library.NoopQueryCleaner{}))
	})

	t.Run("fails when stdin is empty", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		require.Error(t, knowndetect{}.run(ctx, bytes.NewReader(nil), db, library.NoopQueryCleaner{}))
	})
}
