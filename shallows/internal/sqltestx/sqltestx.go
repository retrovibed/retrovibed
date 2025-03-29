package sqltestx

import (
	"database/sql"
	"testing"

	_ "github.com/marcboeker/go-duckdb/v2"
	"github.com/retrovibed/retrovibed/cmd/cmdmeta"
	"github.com/retrovibed/retrovibed/internal/testx"
	"github.com/stretchr/testify/require"
)

func Metadatabase(t testing.TB) *sql.DB {
	ctx, done := testx.Context(t)
	defer done()

	db, err := sql.Open("duckdb", "")
	require.NoError(t, err)
	require.NoError(t, cmdmeta.InitializeDatabase(ctx, db))
	return db
}
