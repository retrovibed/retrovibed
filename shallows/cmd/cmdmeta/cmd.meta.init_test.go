package cmdmeta

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"sync"
	"testing"

	_ "github.com/duckdb/duckdb-go/v2"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/stretchr/testify/require"
)

func metaInitGlobal(ctx context.Context, cancel context.CancelFunc) *cmdopts.Global {
	return &cmdopts.Global{
		Context:  ctx,
		Shutdown: cancel,
		Cleanup:  &sync.WaitGroup{},
	}
}

func metaInitOpen(t *testing.T, path string) *sql.DB {
	db, err := sql.Open("duckdb", path)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	return db
}

func TestMetaInit(t *testing.T) {
	t.Run("creates and migrates both databases in missing directories", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		dir := t.TempDir()
		cmd := MetaInit{
			Database:      filepath.Join(dir, "a", "b", "meta.db"),
			CacheDatabase: filepath.Join(dir, "c", "d", "cache.db"),
		}

		require.NoError(t, cmd.Run(metaInitGlobal(ctx, done)))

		require.FileExists(t, cmd.Database)
		require.FileExists(t, cmd.CacheDatabase)

		meta := metaInitOpen(t, cmd.Database)
		require.Equal(t, 0, testx.Must(sqlx.Count(ctx, meta, "SELECT COUNT(*) FROM meta_profiles"))(t))

		cache := metaInitOpen(t, cmd.CacheDatabase)
		require.Equal(t, 0, testx.Must(sqlx.Count(ctx, cache, "SELECT COUNT(*) FROM library_known_media"))(t))
	})

	t.Run("does not attach the cache tables into the meta database file", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		dir := t.TempDir()
		cmd := MetaInit{
			Database:      filepath.Join(dir, "meta.db"),
			CacheDatabase: filepath.Join(dir, "cache.db"),
		}

		require.NoError(t, cmd.Run(metaInitGlobal(ctx, done)))

		meta := metaInitOpen(t, cmd.Database)
		_, err := sqlx.Count(ctx, meta, "SELECT COUNT(*) FROM library_known_media")
		require.Error(t, err)
	})

	t.Run("idempotent - rerunning preserves existing data", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		dir := t.TempDir()
		cmd := MetaInit{
			Database:      filepath.Join(dir, "meta.db"),
			CacheDatabase: filepath.Join(dir, "cache.db"),
		}

		require.NoError(t, cmd.Run(metaInitGlobal(ctx, done)))

		func() {
			db, err := sql.Open("duckdb", cmd.Database)
			require.NoError(t, err)
			defer db.Close()
			_, err = db.ExecContext(ctx, "CREATE TABLE sentinel (id INTEGER); INSERT INTO sentinel VALUES (1);")
			require.NoError(t, err)
		}()

		require.NoError(t, cmd.Run(metaInitGlobal(ctx, done)))

		meta := metaInitOpen(t, cmd.Database)
		require.Equal(t, 1, testx.Must(sqlx.Count(ctx, meta, "SELECT COUNT(*) FROM sentinel"))(t))
	})

	t.Run("fails when the database directory cannot be created", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		dir := t.TempDir()
		blocker := filepath.Join(dir, "blocker")
		require.NoError(t, os.WriteFile(blocker, nil, 0600))

		cmd := MetaInit{
			Database:      filepath.Join(blocker, "meta.db"),
			CacheDatabase: filepath.Join(dir, "cache.db"),
		}

		err := cmd.Run(metaInitGlobal(ctx, done))
		require.Error(t, err)
		require.Contains(t, err.Error(), "unable to initialize database")
		require.NoFileExists(t, cmd.CacheDatabase)
	})

	t.Run("fails when the cache directory cannot be created", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		dir := t.TempDir()
		blocker := filepath.Join(dir, "blocker")
		require.NoError(t, os.WriteFile(blocker, nil, 0600))

		cmd := MetaInit{
			Database:      filepath.Join(dir, "meta.db"),
			CacheDatabase: filepath.Join(blocker, "cache.db"),
		}

		err := cmd.Run(metaInitGlobal(ctx, done))
		require.Error(t, err)
		require.Contains(t, err.Error(), "unable to initialize database")
	})
}
