package library_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
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
