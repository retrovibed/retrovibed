package library_test

import (
	"database/sql"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
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

		rec, err := library.RecommendationFromRandomKnown(ctx, db, "", "", false)
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

		_, err := library.RecommendationFromRandomKnown(ctx, db, "", "", false)
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

		rec, err := library.RecommendationFromRandomKnown(ctx, db, "", "", false)
		require.NoError(t, err)
		require.Equal(t, known.UID, rec.ContentID)
		require.NotEmpty(t, rec.ID)
	})

	t.Run("sets source to random", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		db := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, known).Scan(&known))

		rec, err := library.RecommendationFromRandomKnown(ctx, db, "", "", false)
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

		_, err := library.RecommendationFromRandomKnown(ctx, db, "", "", false)
		require.NoError(t, err)

		rec, err := library.RecommendationFromRandomKnown(ctx, db, "", "", false)
		require.NoError(t, err)
		require.EqualValues(t, 1, rec.Recommendations)
	})

	t.Run("empty database returns error", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		db := sqltestx.Metadatabase(t)

		_, err := library.RecommendationFromRandomKnown(ctx, db, "", "", false)
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

		_, err := library.RecommendationFromRandomKnown(ctx, db, "", "", false)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("filters by mimetype", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		db := sqltestx.Metadatabase(t)

		var video, audio library.Known
		require.NoError(t, testx.Fake(&video, library.KnownOptionTestDefaults, library.KnownOptionMimetype(mimex.Video)))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, video).Scan(&video))
		require.NoError(t, testx.Fake(&audio, library.KnownOptionTestDefaults, library.KnownOptionMimetype(mimex.Audio)))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, audio).Scan(&audio))

		rec, err := library.RecommendationFromRandomKnown(ctx, db, mimex.Audio, "", false)
		require.NoError(t, err)
		require.Equal(t, mimex.Audio, rec.Mimetype)
	})

	t.Run("copies mimetype from known media", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		db := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionMimetype(mimex.Video)))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, known).Scan(&known))

		rec, err := library.RecommendationFromRandomKnown(ctx, db, mimex.Video, "", false)
		require.NoError(t, err)
		require.Equal(t, mimex.Video, rec.Mimetype)
	})
}
