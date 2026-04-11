package library_test

import (
	"database/sql"
	"testing"

	"github.com/retrovibed/retrovibed/internal/sqltestx"
	"github.com/retrovibed/retrovibed/internal/testx"
	"github.com/retrovibed/retrovibed/internal/timex"
	"github.com/retrovibed/retrovibed/library"
	"github.com/stretchr/testify/require"
)

func TestRecommendationLastGeneratedAt(t *testing.T) {
	t.Run("returns neg infinity when no recommendations exist", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		db := sqltestx.Metadatabase(t)

		ts, err := library.RecommendationLastGeneratedAt(ctx, db, library.RecommendationSourceRandom)
		require.NoError(t, err)
		require.Equal(t, timex.NegInf(), ts)
	})

	t.Run("returns timestamp of most recent recommendation", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		db := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, known).Scan(&known))

		rec, err := library.RecommendationFromRandomKnown(ctx, db)
		require.NoError(t, err)

		ts, err := library.RecommendationLastGeneratedAt(ctx, db, library.RecommendationSourceRandom)
		require.NoError(t, err)
		require.Equal(t, rec.UpdatedAt, ts)
	})

	t.Run("does not return results from a different source", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		db := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, known).Scan(&known))

		_, err := library.RecommendationFromRandomKnown(ctx, db)
		require.NoError(t, err)

		ts, err := library.RecommendationLastGeneratedAt(ctx, db, library.RecommendationSourceGenerative)
		require.NoError(t, err)
		require.Equal(t, timex.NegInf(), ts)
	})
}

func TestRecommendationFromRandomKnown(t *testing.T) {
	t.Run("creates recommendation from random known media", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		db := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, known).Scan(&known))

		rec, err := library.RecommendationFromRandomKnown(ctx, db)
		require.NoError(t, err)
		require.Equal(t, known.UID, rec.KnownMediaID)
		require.NotEmpty(t, rec.ID)
	})

	t.Run("sets source to random", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		db := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, known).Scan(&known))

		rec, err := library.RecommendationFromRandomKnown(ctx, db)
		require.NoError(t, err)
		require.NotEmpty(t, rec.Source)
	})

	t.Run("increments recommendations counter on repeat", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		db := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, known).Scan(&known))

		_, err := library.RecommendationFromRandomKnown(ctx, db)
		require.NoError(t, err)

		rec, err := library.RecommendationFromRandomKnown(ctx, db)
		require.NoError(t, err)
		require.EqualValues(t, 1, rec.Recommendations)
	})

	t.Run("empty database returns error", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		db := sqltestx.Metadatabase(t)

		_, err := library.RecommendationFromRandomKnown(ctx, db)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("excludes adult content", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		db := sqltestx.Metadatabase(t)

		var adult library.Known
		require.NoError(t, testx.Fake(&adult, library.KnownOptionTestDefaults))
		adult.Adult = true
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, adult).Scan(&adult))

		_, err := library.RecommendationFromRandomKnown(ctx, db)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})
}
