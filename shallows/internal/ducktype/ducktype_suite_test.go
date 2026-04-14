package ducktype_test

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/duckdb/duckdb-go/v2"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	testx.Logging()
	os.Exit(m.Run())
}

func newDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("duckdb", "")
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })
	return db
}
