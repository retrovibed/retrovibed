package library_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/internal/mimex"
	"github.com/retrovibed/retrovibed/internal/sqltestx"
	"github.com/retrovibed/retrovibed/internal/testx"
	"github.com/retrovibed/retrovibed/library"
	"github.com/stretchr/testify/require"
)

func TestKnownInsertWithDefaults(t *testing.T) {
	t.Run("inserts known media with mimetype", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		known.Mimetype = mimex.Video
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, known).Scan(&known))
		require.Equal(t, mimex.Video, known.Mimetype)
	})

	t.Run("upsert preserves updated mimetype", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		known.Mimetype = mimex.Video
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, known).Scan(&known))
		require.Equal(t, mimex.Video, known.Mimetype)

		known.Mimetype = mimex.Audio
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, known).Scan(&known))
		require.Equal(t, mimex.Audio, known.Mimetype)
	})
}
