package cmdmedia

import (
	"bytes"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
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

		require.NoError(t, knownquery{Query: "inception"}.run(ctx, db))
	})

	t.Run("finds record matching description", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		known.Title = "Inception"
		known.Overview = "A mind-bending thriller about dreams within dreams"

		var buf bytes.Buffer
		require.NoError(t, jsonl.NewEncoder(&buf).Encode(known))
		require.NoError(t, knownimport{}.run(ctx, db, &buf))

		require.NoError(t, knownquery{Query: "Inception"}.run(ctx, db))
	})

	t.Run("falls back to title search when no fts match", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		known.Title = "The Dark Knight"

		var buf bytes.Buffer
		require.NoError(t, jsonl.NewEncoder(&buf).Encode(known))
		require.NoError(t, knownimport{}.run(ctx, db, &buf))

		require.NoError(t, knownquery{Query: "Dark Knight"}.run(ctx, db))
	})
}
