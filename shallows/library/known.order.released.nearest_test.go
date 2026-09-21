package library_test

import (
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestKnownOrderReleasedNearest(t *testing.T) {
	t.Run("the zero time is a noop", func(t *testing.T) {
		sql, args, err := library.KnownOrderReleasedNearest(time.Time{}).ToSql()
		require.NoError(t, err)
		require.Empty(t, sql)
		require.Empty(t, args)
	})

	t.Run("orders by distance from the target release", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		for _, r := range []time.Time{
			library.KnownStringRelease("2020-01-01"),
			library.KnownStringRelease("2018-07-30"),
			library.KnownStringRelease("2018-07-20"),
			library.KnownStringRelease("2018-08-15"),
			library.KnownStringRelease("1999-01-01"),
		} {
			var known library.Known
			require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionReleased(r)))
			require.NoError(t, library.KnownInsertWithDefaults(ctx, db, known).Scan(&known))
		}

		var found []library.Known
		b := library.KnownSearchBuilder().OrderByClause(library.KnownOrderReleasedNearest(library.KnownStringRelease("2018-07-27")))
		require.NoError(t, sqlx.ScanInto(library.KnownSearch(ctx, db, b), &found))

		actual := make([]string, 0, len(found))
		for _, k := range found {
			actual = append(actual, k.Released.UTC().Format(time.DateOnly))
		}

		require.Equal(t, []string{"2018-07-30", "2018-07-20", "2018-08-15", "2020-01-01", "1999-01-01"}, actual)
	})

	t.Run("orders rows with an infinite release last", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var unknown library.Known
		require.NoError(t, testx.Fake(&unknown, library.KnownOptionTestDefaults, library.KnownOptionReleased(library.KnownStringRelease("2018-07-27"))))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, unknown).Scan(&unknown))
		_, err := db.ExecContext(ctx, "UPDATE cache.library_known_media SET released = 'infinity' WHERE uid = ?", unknown.UID)
		require.NoError(t, err)

		for _, r := range []string{"1999-01-01", "2018-07-20"} {
			var known library.Known
			require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionReleased(library.KnownStringRelease(r))))
			require.NoError(t, library.KnownInsertWithDefaults(ctx, db, known).Scan(&known))
		}

		var found []library.Known
		b := library.KnownSearchBuilder().OrderByClause(library.KnownOrderReleasedNearest(library.KnownStringRelease("2018-07-27")))
		require.NoError(t, sqlx.ScanInto(library.KnownSearch(ctx, db, b), &found))

		require.Len(t, found, 3)
		require.Equal(t, "2018-07-20", found[0].Released.UTC().Format(time.DateOnly))
		require.Equal(t, "1999-01-01", found[1].Released.UTC().Format(time.DateOnly))
		require.Equal(t, unknown.UID, found[2].UID)
	})

	t.Run("orders by release after parent uid nil last and nearest collation", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		parent := uuid.Must(uuid.NewV4()).String()
		for _, r := range []struct {
			parent    string
			collation uint32
			released  string
		}{
			{uuid.Nil.String(), 0, "2018-07-27"},
			{parent, library.KnownCollationEpisode(1, 1), "2018-07-27"},
			{parent, library.KnownCollationEpisode(1, 1), "2018-08-27"},
			{parent, library.KnownCollationEpisode(1, 1), "2018-07-26"},
			{parent, library.KnownCollationEpisode(1, 2), "2018-07-27"},
		} {
			var known library.Known
			require.NoError(t, testx.Fake(
				&known,
				library.KnownOptionTestDefaults,
				library.KnownOptionParentUID(r.parent),
				library.KnownOptionCollation(r.collation),
				library.KnownOptionReleased(library.KnownStringRelease(r.released)),
			))
			require.NoError(t, library.KnownInsertWithDefaults(ctx, db, known).Scan(&known))
		}

		var found []library.Known
		b := library.KnownSearchBuilder().
			OrderByClause(library.KnownOrderParentUIDNilLast()).
			OrderByClause(library.KnownOrderCollationNearest(library.KnownCollationEpisode(1, 1))).
			OrderByClause(library.KnownOrderReleasedNearest(library.KnownStringRelease("2018-07-27"))).
			OrderBy("title DESC")
		require.NoError(t, sqlx.ScanInto(library.KnownSearch(ctx, db, b), &found))

		type row struct {
			ParentUID string
			Collation uint32
			Released  string
		}

		actual := make([]row, 0, len(found))
		for _, k := range found {
			actual = append(actual, row{ParentUID: k.ParentUID, Collation: k.Collation, Released: k.Released.UTC().Format(time.DateOnly)})
		}

		require.Equal(t, []row{
			{parent, library.KnownCollationEpisode(1, 1), "2018-07-27"},
			{parent, library.KnownCollationEpisode(1, 1), "2018-07-26"},
			{parent, library.KnownCollationEpisode(1, 1), "2018-08-27"},
			{parent, library.KnownCollationEpisode(1, 2), "2018-07-27"},
			{uuid.Nil.String(), 0, "2018-07-27"},
		}, actual)
	})

	t.Run("a zero target composes as a noop between other orderings", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, known).Scan(&known))

		var found []library.Known
		b := library.KnownSearchBuilder().
			OrderByClause(library.KnownOrderParentUIDNilLast()).
			OrderByClause(library.KnownOrderReleasedNearest(time.Time{})).
			OrderBy("title DESC")
		require.NoError(t, sqlx.ScanInto(library.KnownSearch(ctx, db, b), &found))
		require.Len(t, found, 1)
	})
}
