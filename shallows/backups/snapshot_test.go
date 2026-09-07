package backups_test

import (
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/backups"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/stretchr/testify/require"
)

func TestSnapshot(t *testing.T) {
	t.Run("a snapshot restores to the same contents under the same key", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		src := sqltestx.Metadatabase(t)
		_, err := src.ExecContext(ctx, "CREATE TABLE fixture AS SELECT range AS id FROM range(1000)")
		require.NoError(t, err)

		key, err := backups.Key("seed", []byte("identity"))
		require.NoError(t, err)

		encrypted := filepath.Join(t.TempDir(), "backup.db")
		require.NoError(t, backups.Snapshot(ctx, src, encrypted, key))

		dst, err := sql.Open("duckdb", filepath.Join(t.TempDir(), "restored.db"))
		require.NoError(t, err)
		defer dst.Close()

		require.NoError(t, backups.Restore(ctx, dst, encrypted, key))

		n, err := sqlx.Count(ctx, dst, "SELECT COUNT(*) FROM fixture")
		require.NoError(t, err)
		require.Equal(t, 1000, n)

		// hnsw indexes cannot be written to an encrypted database; the owning package
		// recreates them on start, so the restored copy has none and the source keeps its own.
		n, err = sqlx.Count(ctx, dst, "SELECT COUNT(*) FROM duckdb_indexes() WHERE sql ILIKE '%USING HNSW%'")
		require.NoError(t, err)
		require.Equal(t, 0, n)
		n, err = sqlx.Count(ctx, src, "SELECT COUNT(*) FROM duckdb_indexes() WHERE sql ILIKE '%USING HNSW%'")
		require.NoError(t, err)
		require.Equal(t, 1, n)
		require.False(t, fsx.Exists(encrypted+".staging"), "the plain staging copy is removed")
	})

	t.Run("the wrong key cannot open a snapshot", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		src := sqltestx.Metadatabase(t)

		key, err := backups.Key("seed", []byte("identity"))
		require.NoError(t, err)
		wrong, err := backups.Key("seed", []byte("another identity"))
		require.NoError(t, err)

		encrypted := filepath.Join(t.TempDir(), "backup.db")
		require.NoError(t, backups.Snapshot(ctx, src, encrypted, key))

		dst, err := sql.Open("duckdb", filepath.Join(t.TempDir(), "restored.db"))
		require.NoError(t, err)
		defer dst.Close()

		require.Error(t, backups.Restore(ctx, dst, encrypted, wrong))
	})

	t.Run("a snapshot leaves the source usable", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		src := sqltestx.Metadatabase(t)

		key, err := backups.Key("seed", []byte("identity"))
		require.NoError(t, err)

		require.NoError(t, backups.Snapshot(ctx, src, filepath.Join(t.TempDir(), "backup.db"), key))

		// the backup catalog is detached on the way out, so a second snapshot attaches cleanly.
		require.NoError(t, backups.Snapshot(ctx, src, filepath.Join(t.TempDir(), "again.db"), key))
	})
}
