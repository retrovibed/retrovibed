package daemons_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/retrovibe/daemons"
	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/userx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestRecommendationsBackgroundRun(t *testing.T) {
	t.Run("generates a video recommendation when none exist", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionMimetype(mimex.Video)))
		known.OriginalLanguage = userx.LocaleLanguage()
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		require.NoError(t, daemons.RecommendationsBackgroundRun(ctx, q))

		count, err := sqlx.Count(ctx, q, "SELECT COUNT(*) FROM library_recommendations")
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})

	t.Run("generates a audio recommendation when none exist", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionMimetype(mimex.Audio)))
		known.OriginalLanguage = userx.LocaleLanguage()
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
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionMimetype(mimex.Video)))
		known.OriginalLanguage = userx.LocaleLanguage()
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
