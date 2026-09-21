package library_test

import (
	"database/sql"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestKnownScoreByID(t *testing.T) {
	t.Run("returns positive score for matching title", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		known.Title = "The Dark Knight"
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, known).Scan(&known))

		var relevance float64
		require.NoError(t, library.KnownScoreByID(ctx, db, known.UID, "The Dark Knight", 0.7).Scan(&relevance))
		require.Greater(t, relevance, 0.0)
	})

	t.Run("returns low score for non-matching title", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		known.Title = "The Dark Knight"
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, known).Scan(&known))

		var relevance float64
		require.NoError(t, library.KnownScoreByID(ctx, db, known.UID, "xyzzy unrelated query", 0.7).Scan(&relevance))
		require.Less(t, relevance, 0.5)
	})

	t.Run("returns no rows for unknown uid", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var relevance float64
		err := library.KnownScoreByID(ctx, db, "00000000-0000-0000-0000-000000000000", "anything", 0.7).Scan(&relevance)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})
}
