package duckdbx_test

import (
	"database/sql"
	"testing"

	_ "github.com/duckdb/duckdb-go/v2"
	"github.com/retrovibed/retrovibed/shallows/internal/duckdbx"
	"github.com/stretchr/testify/require"
)

func TestSearchToSql(t *testing.T) {
	t.Run("query", func(t *testing.T) {
		q, args, err := duckdbx.Search("arch linux", "a", "b").ToSql()
		require.NoError(t, err)
		require.Equal(t, `SELECT true FROM (SELECT UNNEST($1::VARCHAR[]) AS val) AS t WHERE ("val" ILIKE '%' || $2 || '%') AND ("val" ILIKE '%' || $3 || '%')`, q)
		require.Equal(t, []any{[]string{"a", "b"}, "arch", "linux"}, args)
	})

	t.Run("blank query", func(t *testing.T) {
		q, args, err := duckdbx.Search("", "a", "b").ToSql()
		require.NoError(t, err)
		require.Equal(t, `SELECT true FROM (SELECT UNNEST($1::VARCHAR[]) AS val) AS t`, q)
		require.Equal(t, []any{[]string{"a", "b"}}, args)
	})
}

func TestSearch(t *testing.T) {
	db, err := sql.Open("duckdb", "")
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	values := []string{
		"golang has built-in concurrency support",
		"apache lucene provides full text indexing",
		"bleve is a golang library for text search",
	}

	t.Run("matches boolean query", func(t *testing.T) {
		rows, err := duckdbx.Search(`golang AND ("text search" OR concurrency)`, values...).RunWith(db).Query()
		require.NoError(t, err)
		defer rows.Close()

		count := 0
		for rows.Next() {
			var matched bool
			require.NoError(t, rows.Scan(&matched))
			require.True(t, matched)
			count++
		}
		require.NoError(t, rows.Err())
		require.Equal(t, 2, count)
	})

	t.Run("blank query matches everything", func(t *testing.T) {
		rows, err := duckdbx.Search("", values...).RunWith(db).Query()
		require.NoError(t, err)
		defer rows.Close()

		count := 0
		for rows.Next() {
			count++
		}
		require.NoError(t, rows.Err())
		require.Equal(t, len(values), count)
	})

	t.Run("no matches", func(t *testing.T) {
		rows, err := duckdbx.Search("nonexistentterm12345", values...).RunWith(db).Query()
		require.NoError(t, err)
		defer rows.Close()

		require.False(t, rows.Next())
		require.NoError(t, rows.Err())
	})
}
