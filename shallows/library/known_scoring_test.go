package library_test

import (
	"database/sql"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
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

func TestKnownBestMatch(t *testing.T) {
	t.Run("returns best matching entry for exact title", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		known.Title = "Inception"
		known.Adult = false
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, known).Scan(&known))

		var result library.Known
		require.NoError(t, library.KnownBestMatch(ctx, db, mimex.Application, "Inception", 0.7).Scan(&result))
		require.Equal(t, known.UID, result.UID)
	})

	t.Run("returns best match among multiple candidates", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		for _, title := range []string{"Interstellar", "Inception", "The Prestige"} {
			var known library.Known
			require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
			known.Title = title
			known.Adult = false
			require.NoError(t, library.KnownInsertWithDefaults(ctx, db, known).Scan(&known))
		}

		var result library.Known
		require.NoError(t, library.KnownBestMatch(ctx, db, mimex.Application, "Inception", 0.7).Scan(&result))
		require.Equal(t, "Inception", result.Title)
	})

	t.Run("returns no rows when nothing matches above cutoff", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		known.Title = "Inception"
		known.Adult = false
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, known).Scan(&known))

		var result library.Known
		err := library.KnownBestMatch(ctx, db, mimex.Application, "xyzzy completely unrelated", 0.99).Scan(&result)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("excludes adult content", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		known.Title = "Inception"
		known.Adult = true
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, known).Scan(&known))

		var result library.Known
		err := library.KnownBestMatch(ctx, db, mimex.Application, "Inception", 0.7).Scan(&result)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("blank mime matches any mimetype", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var video library.Known
		require.NoError(t, testx.Fake(&video, library.KnownOptionTestDefaults, library.KnownOptionMimetype(mimex.Video)))
		video.Title = "Inception"
		video.Adult = false
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, video).Scan(&video))

		var result library.Known
		require.NoError(t, library.KnownBestMatch(ctx, db, "", "Inception", 0.7).Scan(&result))
		require.Equal(t, video.UID, result.UID)
	})

	t.Run("mime filters out non-matching mimetype", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var video library.Known
		require.NoError(t, testx.Fake(&video, library.KnownOptionTestDefaults, library.KnownOptionMimetype(mimex.Video)))
		video.Title = "Inception"
		video.Adult = false
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, video).Scan(&video))

		var result library.Known
		err := library.KnownBestMatch(ctx, db, mimex.Audio, "Inception", 0.7).Scan(&result)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})
}
