package backups_test

import (
	"database/sql"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/backups"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/stretchr/testify/require"
)

func TestSnapshot(t *testing.T) {
	// encrypts a snapshot the way the upload does, into a file the way the download does.
	encrypt := func(t *testing.T, key backups.Key, src string) string {
		f, err := os.Open(src)
		require.NoError(t, err)
		defer f.Close()

		r, err := key.Encrypt(f)
		require.NoError(t, err)

		dst := filepath.Join(t.TempDir(), "backup.enc")
		out, err := os.Create(dst)
		require.NoError(t, err)
		_, err = io.Copy(out, r)
		require.NoError(t, err)
		require.NoError(t, out.Close())
		return dst
	}

	t.Run("a snapshot restores to the same contents under the same key", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		src := sqltestx.Metadatabase(t)
		_, err := src.ExecContext(ctx, "CREATE TABLE fixture AS SELECT range AS id FROM range(1000)")
		require.NoError(t, err)

		key, err := backups.NewKey("seed", []byte("identity"))
		require.NoError(t, err)

		snapshot := filepath.Join(t.TempDir(), "backup.db")
		require.NoError(t, backups.Snapshot(ctx, src, snapshot))

		restored := filepath.Join(t.TempDir(), "restored.db")
		require.NoError(t, backups.Restore(key, encrypt(t, key, snapshot), restored))

		dst, err := sql.Open("duckdb", restored)
		require.NoError(t, err)
		defer dst.Close()

		n, err := sqlx.Count(ctx, dst, "SELECT COUNT(*) FROM fixture")
		require.NoError(t, err)
		require.Equal(t, 1000, n)

		// the copy is a plain duckdb file, so every index survives, hnsw included.
		n, err = sqlx.Count(ctx, dst, "SELECT COUNT(*) FROM duckdb_indexes() WHERE sql ILIKE '%USING HNSW%'")
		require.NoError(t, err)
		require.Equal(t, 1, n)
	})

	t.Run("the wrong key does not restore a database", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		src := sqltestx.Metadatabase(t)

		key, err := backups.NewKey("seed", []byte("identity"))
		require.NoError(t, err)
		wrong, err := backups.NewKey("seed", []byte("another identity"))
		require.NoError(t, err)

		snapshot := filepath.Join(t.TempDir(), "backup.db")
		require.NoError(t, backups.Snapshot(ctx, src, snapshot))

		restored := filepath.Join(t.TempDir(), "restored.db")
		require.NoError(t, backups.Restore(wrong, encrypt(t, key, snapshot), restored))

		// the driver refuses the file on open; the ping covers a driver that defers it.
		opened := func() error {
			dst, err := sql.Open("duckdb", restored)
			if err != nil {
				return err
			}
			defer dst.Close()
			return dst.PingContext(ctx)
		}
		require.Error(t, opened())
	})

	t.Run("a snapshot leaves the source usable", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		src := sqltestx.Metadatabase(t)

		require.NoError(t, backups.Snapshot(ctx, src, filepath.Join(t.TempDir(), "backup.db")))

		// the backup catalog is detached on the way out, so a second snapshot attaches cleanly.
		require.NoError(t, backups.Snapshot(ctx, src, filepath.Join(t.TempDir(), "again.db")))
	})
}
