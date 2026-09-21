package library_test

import (
	"math"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestKnownQueryCollation(t *testing.T) {
	t.Run("MaxUint32 is a noop", func(t *testing.T) {
		sql, args, err := library.KnownQueryCollation(math.MaxUint32).ToSql()
		require.NoError(t, err)
		require.Empty(t, sql)
		require.Empty(t, args)
	})

	t.Run("generates an equality predicate", func(t *testing.T) {
		sql, args, err := library.KnownQueryCollation(library.KnownCollationEpisode(1, 2)).ToSql()
		require.NoError(t, err)
		require.Equal(t, "cache.library_known_media.\"collation\" = ?", sql)
		require.Equal(t, []interface{}{library.KnownCollationEpisode(1, 2)}, args)
	})

	t.Run("zero is a filter, not a noop", func(t *testing.T) {
		sql, args, err := library.KnownQueryCollation(0).ToSql()
		require.NoError(t, err)
		require.NotEmpty(t, sql)
		require.Equal(t, []interface{}{uint32(0)}, args)
	})

	t.Run("matches only the requested episode", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		collations := []uint32{
			0,
			library.KnownCollationEpisode(0, 1),
			library.KnownCollationEpisode(1, 1),
			library.KnownCollationEpisode(1, 2),
			library.KnownCollationEpisode(2, 1),
		}

		for _, c := range collations {
			var known library.Known
			require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionCollation(c)))
			require.NoError(t, library.KnownInsertWithDefaults(ctx, db, known).Scan(&known))
		}

		for _, c := range collations {
			var found []library.Known
			b := library.KnownSearchBuilder().Where(library.KnownQueryCollation(c))
			require.NoError(t, sqlx.ScanInto(library.KnownSearch(ctx, db, b), &found))
			require.Len(t, found, 1)
			require.Equal(t, c, found[0].Collation)
		}
	})

	t.Run("resolves a string episode to the matching row", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		for _, c := range []uint32{library.KnownCollationEpisode(1, 2), library.KnownCollationEpisode(0, 2)} {
			var known library.Known
			require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionCollation(c)))
			require.NoError(t, library.KnownInsertWithDefaults(ctx, db, known).Scan(&known))
		}

		var found []library.Known
		b := library.KnownSearchBuilder().Where(library.KnownQueryCollation(library.KnownStringCollationEpisode("S00E02")))
		require.NoError(t, sqlx.ScanInto(library.KnownSearch(ctx, db, b), &found))
		require.Len(t, found, 1)
		require.Equal(t, library.KnownCollationEpisode(0, 2), found[0].Collation)
	})

	t.Run("an unparseable string episode does not filter", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		for _, c := range []uint32{0, library.KnownCollationEpisode(1, 1), library.KnownCollationEpisode(1, 2)} {
			var known library.Known
			require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionCollation(c)))
			require.NoError(t, library.KnownInsertWithDefaults(ctx, db, known).Scan(&known))
		}

		var found []library.Known
		b := library.KnownSearchBuilder().Where(squirrel.And{
			library.KnownQueryCollation(library.KnownStringCollationEpisode("garbage")),
			library.KnownQueryNotTombstoned(),
		})
		require.NoError(t, sqlx.ScanInto(library.KnownSearch(ctx, db, b), &found))
		require.Len(t, found, 3)
	})
}
