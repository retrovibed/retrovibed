package library_test

import (
	"math"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestKnownOrderCollationNearest(t *testing.T) {
	t.Run("MaxUint32 targets collation 0", func(t *testing.T) {
		expected, eargs, err := library.KnownOrderCollationNearest(0).ToSql()
		require.NoError(t, err)

		sql, args, err := library.KnownOrderCollationNearest(math.MaxUint32).ToSql()
		require.NoError(t, err)
		require.Equal(t, expected, sql)
		require.Equal(t, eargs, args)
		require.Equal(t, []interface{}{uint32(0)}, args)
	})

	t.Run("MaxUint32 orders collation 0 first", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		for _, c := range []uint32{
			library.KnownCollationEpisode(3, 1),
			0,
			library.KnownCollationEpisode(1, 2),
		} {
			var known library.Known
			require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionCollation(c)))
			require.NoError(t, library.KnownInsertWithDefaults(ctx, db, known).Scan(&known))
		}

		var found []library.Known
		b := library.KnownSearchBuilder().OrderByClause(library.KnownOrderCollationNearest(math.MaxUint32))
		require.NoError(t, sqlx.ScanInto(library.KnownSearch(ctx, db, b), &found))

		actual := make([]uint32, 0, len(found))
		for _, k := range found {
			actual = append(actual, k.Collation)
		}

		require.Equal(t, []uint32{
			0,
			library.KnownCollationEpisode(1, 2),
			library.KnownCollationEpisode(3, 1),
		}, actual)
	})

	t.Run("orders by distance from the target collation", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		for _, c := range []uint32{
			library.KnownCollationEpisode(3, 1),
			library.KnownCollationEpisode(1, 9),
			library.KnownCollationEpisode(1, 2),
			library.KnownCollationEpisode(1, 4),
			0,
		} {
			var known library.Known
			require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionCollation(c)))
			require.NoError(t, library.KnownInsertWithDefaults(ctx, db, known).Scan(&known))
		}

		var found []library.Known
		b := library.KnownSearchBuilder().OrderByClause(library.KnownOrderCollationNearest(library.KnownCollationEpisode(1, 5)))
		require.NoError(t, sqlx.ScanInto(library.KnownSearch(ctx, db, b), &found))

		actual := make([]uint32, 0, len(found))
		for _, k := range found {
			actual = append(actual, k.Collation)
		}

		require.Equal(t, []uint32{
			library.KnownCollationEpisode(1, 4),
			library.KnownCollationEpisode(1, 2),
			library.KnownCollationEpisode(1, 9),
			0,
			library.KnownCollationEpisode(3, 1),
		}, actual)
	})

	t.Run("orders by distance within each parent group when combined with parent uid nil last", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		parent := uuid.Must(uuid.NewV4()).String()
		for _, r := range []struct {
			parent    string
			collation uint32
		}{
			{uuid.Nil.String(), 0},
			{parent, library.KnownCollationEpisode(3, 1)},
			{uuid.Nil.String(), library.KnownCollationEpisode(1, 5)},
			{parent, library.KnownCollationEpisode(1, 9)},
			{parent, library.KnownCollationEpisode(1, 4)},
		} {
			var known library.Known
			require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionParentUID(r.parent), library.KnownOptionCollation(r.collation)))
			require.NoError(t, library.KnownInsertWithDefaults(ctx, db, known).Scan(&known))
		}

		var found []library.Known
		b := library.KnownSearchBuilder().
			OrderByClause(library.KnownOrderParentUIDNilLast()).
			OrderByClause(library.KnownOrderCollationNearest(library.KnownCollationEpisode(1, 5)))
		require.NoError(t, sqlx.ScanInto(library.KnownSearch(ctx, db, b), &found))

		type row struct {
			ParentUID string
			Collation uint32
		}

		actual := make([]row, 0, len(found))
		for _, k := range found {
			actual = append(actual, row{ParentUID: k.ParentUID, Collation: k.Collation})
		}

		// the exact collation match has a nil parent, so it still sorts after
		// every row that has a parent, however far away those rows are.
		require.Equal(t, []row{
			{parent, library.KnownCollationEpisode(1, 4)},
			{parent, library.KnownCollationEpisode(1, 9)},
			{parent, library.KnownCollationEpisode(3, 1)},
			{uuid.Nil.String(), library.KnownCollationEpisode(1, 5)},
			{uuid.Nil.String(), 0},
		}, actual)
	})

	t.Run("orders collation 0 first within each parent group when there is no target collation", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		parent := uuid.Must(uuid.NewV4()).String()
		for _, r := range []struct {
			parent    string
			collation uint32
		}{
			{uuid.Nil.String(), library.KnownCollationEpisode(2, 1)},
			{parent, library.KnownCollationEpisode(1, 2)},
			{uuid.Nil.String(), 0},
			{parent, 0},
		} {
			var known library.Known
			require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionParentUID(r.parent), library.KnownOptionCollation(r.collation)))
			require.NoError(t, library.KnownInsertWithDefaults(ctx, db, known).Scan(&known))
		}

		var found []library.Known
		b := library.KnownSearchBuilder().
			OrderByClause(library.KnownOrderParentUIDNilLast()).
			OrderByClause(library.KnownOrderCollationNearest(math.MaxUint32)).
			OrderBy("title DESC")
		require.NoError(t, sqlx.ScanInto(library.KnownSearch(ctx, db, b), &found))

		type row struct {
			ParentUID string
			Collation uint32
		}

		actual := make([]row, 0, len(found))
		for _, k := range found {
			actual = append(actual, row{ParentUID: k.ParentUID, Collation: k.Collation})
		}

		require.Equal(t, []row{
			{parent, 0},
			{parent, library.KnownCollationEpisode(1, 2)},
			{uuid.Nil.String(), 0},
			{uuid.Nil.String(), library.KnownCollationEpisode(2, 1)},
		}, actual)
	})
}
