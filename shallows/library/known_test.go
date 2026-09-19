package library_test

import (
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestKnownOptionAutoParentUID(t *testing.T) {
	t.Run("defaults blank parent uid to the nil uuid", func(t *testing.T) {
		var known library.Known
		library.KnownOptionAutoParentUID(&known)
		require.Equal(t, uuid.Nil.String(), known.ParentUID)
	})

	t.Run("preserves an existing parent uid", func(t *testing.T) {
		parent := uuid.Must(uuid.NewV4()).String()
		known := library.Known{ParentUID: parent}
		library.KnownOptionAutoParentUID(&known)
		require.Equal(t, parent, known.ParentUID)
	})
}

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

func TestKnownInsertWithDefaultsTOFU(t *testing.T) {
	t.Run("inserts known media", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		known.Mimetype = mimex.Video
		require.NoError(t, library.KnownInsertWithDefaultsTOFU(ctx, db, known).Scan(&known))
		require.Equal(t, mimex.Video, known.Mimetype)
	})

	t.Run("noop when record already exists", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var original library.Known
		require.NoError(t, testx.Fake(&original, library.KnownOptionTestDefaults))
		original.Mimetype = mimex.Video
		require.NoError(t, library.KnownInsertWithDefaultsTOFU(ctx, db, original).Scan(&original))

		conflicting := original
		conflicting.Title = original.Title + " changed"
		conflicting.Mimetype = mimex.Audio

		var scanned library.Known
		require.NoError(t, library.KnownInsertWithDefaultsTOFU(ctx, db, conflicting).Scan(&scanned))
		require.Equal(t, original.Title, scanned.Title)
		require.Equal(t, original.Mimetype, scanned.Mimetype)

		var found library.Known
		require.NoError(t, library.KnownFindByID(ctx, db, original.UID).Scan(&found))
		require.Equal(t, original.Title, found.Title)
		require.Equal(t, original.Mimetype, found.Mimetype)
	})
}
