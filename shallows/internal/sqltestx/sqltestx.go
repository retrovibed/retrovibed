package sqltestx

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/duckdb/duckdb-go/v2"
	"github.com/retrovibed/retrovibed/cmd/cmdmeta"
	"github.com/retrovibed/retrovibed/internal/sqlx"
	"github.com/stretchr/testify/require"
)

func Metadatabase(t testing.TB) *sql.DB {
	db, err := sql.Open("duckdb", "")
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })
	require.NoError(t, cmdmeta.InitializeDatabase(t.Context(), db))
	return db
}

func Count(t testing.TB, q sqlx.Queryer, query string, args ...any) (count int) {
	require.NoError(t, sqlx.NewValueRowScanner[int](q.QueryRowContext(t.Context(), query, args...)).Scan(&count))
	return count
}

func Timestamp(t testing.TB, q sqlx.Queryer, query string) (ts time.Time) {
	require.NoError(t, sqlx.NewValueRowScanner[time.Time](q.QueryRowContext(t.Context(), query)).Scan(&ts))
	return ts
}

func String(t testing.TB, q sqlx.Queryer, query string) (v string) {
	require.NoError(t, sqlx.NewValueRowScanner[string](q.QueryRowContext(t.Context(), query)).Scan(&v))
	return v
}

func Bool(t testing.TB, q sqlx.Queryer, query string) (v bool) {
	require.NoError(t, sqlx.NewValueRowScanner[bool](q.QueryRowContext(t.Context(), query)).Scan(&v))
	return v
}
