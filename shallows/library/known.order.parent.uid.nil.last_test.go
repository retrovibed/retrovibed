package library_test

import (
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestKnownOrderParentUIDNilLast(t *testing.T) {
	t.Run("orders nil parent rows after rows with a parent", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		parent := uuid.Must(uuid.NewV4()).String()
		for _, p := range []string{uuid.Nil.String(), parent, uuid.Nil.String(), parent} {
			var known library.Known
			require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionParentUID(p)))
			require.NoError(t, library.KnownInsertWithDefaults(ctx, db, known).Scan(&known))
		}

		var found []library.Known
		b := library.KnownSearchBuilder().OrderByClause(library.KnownOrderParentUIDNilLast())
		require.NoError(t, sqlx.ScanInto(library.KnownSearch(ctx, db, b), &found))

		require.Len(t, found, 4)
		require.Equal(t, parent, found[0].ParentUID)
		require.Equal(t, parent, found[1].ParentUID)
		require.Equal(t, uuid.Nil.String(), found[2].ParentUID)
		require.Equal(t, uuid.Nil.String(), found[3].ParentUID)
	})
}
