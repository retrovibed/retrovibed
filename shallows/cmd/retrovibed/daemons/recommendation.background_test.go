package daemons_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/cmd/retrovibed/daemons"
	"github.com/retrovibed/retrovibed/internal/sqltestx"
	"github.com/retrovibed/retrovibed/internal/sqlx"
	"github.com/retrovibed/retrovibed/internal/testx"
	"github.com/retrovibed/retrovibed/library"
	"github.com/stretchr/testify/require"
)

func TestRecommendationsBackgroundRun(t *testing.T) {
	t.Run("generates a recommendation when none exist", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		require.NoError(t, daemons.RecommendationsBackgroundRun(ctx, q))

		count, err := sqlx.Count(ctx, q, "SELECT COUNT(*) FROM library_recommendations")
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})

	t.Run("skips generation when last recommendation is within 24 hours", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		// generate the first recommendation
		require.NoError(t, daemons.RecommendationsBackgroundRun(ctx, q))

		// second run should not generate another
		require.NoError(t, daemons.RecommendationsBackgroundRun(ctx, q))

		count, err := sqlx.Count(ctx, q, "SELECT COUNT(*) FROM library_recommendations")
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})

	t.Run("no known media returns no error", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, daemons.RecommendationsBackgroundRun(ctx, q))
	})
}
