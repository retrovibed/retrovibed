package sqltestx

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/duckdb/duckdb-go/v2"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/stretchr/testify/require"
)

func Metadatabase(t testing.TB) *sql.DB {
	db, err := sql.Open("duckdb", "")
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })
	require.NoError(t, cmdopts.InitializeDatabase(t.Context(), db))
	return db
}

func Count(t testing.TB, q sqlx.Queryer, query string, args ...any) (count int) {
	ctx := t.Context()
	err := sqlx.NewValueRowScanner[int](q.QueryRowContext(ctx, query, args...)).Scan(&count)
	// t's context is only canceled once the test starts tearing down. A canceled ctx
	// here means this call is an orphaned require.Eventually/Never polling goroutine
	// (testify doesn't join them at the deadline) still in flight after the subtest
	// already returned — require.NoError would panic with "Fail in goroutine after
	// <test> has completed", so treat it as condition-not-met instead.
	if err != nil && ctx.Err() != nil {
		return count
	}
	require.NoError(t, err)
	return count
}

func Timestamp(t testing.TB, q sqlx.Queryer, query string) (ts time.Time) {
	ctx := t.Context()
	err := sqlx.NewValueRowScanner[time.Time](q.QueryRowContext(ctx, query)).Scan(&ts)
	// t's context is only canceled once the test starts tearing down. A canceled ctx
	// here means this call is an orphaned require.Eventually/Never polling goroutine
	// (testify doesn't join them at the deadline) still in flight after the subtest
	// already returned — require.NoError would panic with "Fail in goroutine after
	// <test> has completed", so treat it as condition-not-met instead.
	if err != nil && ctx.Err() != nil {
		return ts
	}
	require.NoError(t, err)
	return ts
}

func String(t testing.TB, q sqlx.Queryer, query string) (v string) {
	ctx := t.Context()
	err := sqlx.NewValueRowScanner[string](q.QueryRowContext(ctx, query)).Scan(&v)
	// t's context is only canceled once the test starts tearing down. A canceled ctx
	// here means this call is an orphaned require.Eventually/Never polling goroutine
	// (testify doesn't join them at the deadline) still in flight after the subtest
	// already returned — require.NoError would panic with "Fail in goroutine after
	// <test> has completed", so treat it as condition-not-met instead.
	if err != nil && ctx.Err() != nil {
		return v
	}
	require.NoError(t, err)
	return v
}

func Bool(t testing.TB, q sqlx.Queryer, query string) (v bool) {
	ctx := t.Context()
	err := sqlx.NewValueRowScanner[bool](q.QueryRowContext(ctx, query)).Scan(&v)
	// t's context is only canceled once the test starts tearing down. A canceled ctx
	// here means this call is an orphaned require.Eventually/Never polling goroutine
	// (testify doesn't join them at the deadline) still in flight after the subtest
	// already returned — require.NoError would panic with "Fail in goroutine after
	// <test> has completed", so treat it as condition-not-met instead.
	if err != nil && ctx.Err() != nil {
		return v
	}
	require.NoError(t, err)
	return v
}
