package library_test

import (
	"database/sql"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecommendationSourceMD5(t *testing.T) {
	t.Run("fixed md5 sums", func(t *testing.T) {
		// any changes to these will require downstream changes in the ui.
		assert.Equal(t, "7ddf32e1-7a6a-c5ce-04a8-ecbf782ca509", md5x.String(library.RecommendationSourceRandom), "mismatched md5 for source random")
		assert.Equal(t, "538416cf-3bc5-9332-670a-f4cae9485ebe", md5x.String(library.RecommendationSourceDiscovered), "mismatched md5 for source discovered")
		assert.Equal(t, "e15ee067-d8b5-a64b-ffd9-617121fa925b", md5x.String(library.RecommendationSourceGenerative), "mismatched md5 for source generative")
		assert.Equal(t, "ab1c952c-a77a-0bc1-3e23-68588d7a0c6e", md5x.String(library.RecommendationSourceSearchPlugin), "mismatched md5 for source search plugin")
	})
}

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
