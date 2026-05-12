package cmdmedia

import (
	"testing"

	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestKnownDetectRun(t *testing.T) {
	cmd := knowndetect{}

	t.Run("finds a matching record by title", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		known.Title = "The Grand Budapest Hotel"

		require.NoError(t, knownimport{}.run(ctx, db, jsonlBuffer(t, known)))

		result, err := cmd.run(ctx, db, "The Grand Budapest Hotel")
		require.NoError(t, err)
		require.Equal(t, known.UID, result.UID)
	})

	t.Run("returns Unknown when no match found", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		result, err := cmd.run(ctx, db, "something completely unknown")
		require.NoError(t, err)
		require.Equal(t, library.Unknown().UID, result.UID)
	})
}
