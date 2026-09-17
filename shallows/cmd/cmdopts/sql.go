package cmdopts

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync/atomic"

	duckdb "github.com/duckdb/duckdb-go/v2"

	"github.com/retrovibed/retrovibed/retroapi/userx"
	"github.com/retrovibed/retrovibed/shallows/internal/debugx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/goosex"
)

//go:embed .migrations/*.sql
var embedmigrations embed.FS

func DatabaseMeta(ctx context.Context) (db *sql.DB, err error) {
	return DatabaseCustom(ctx, userx.DefaultConfigDir(userx.DefaultRelRoot(), "meta.db"))
}

func DatabaseCustom(ctx context.Context, path string) (db *sql.DB, err error) {
	log.Println("database path", path)

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, err
	}

	if db, err = sql.Open("duckdb", path); err != nil {
		return nil, errorsx.Wrap(err, "unable to open db")
	}
	defer func() {
		if err == nil {
			return
		}
		debugx.Println("closing database due to error during initialization", err)
		errorsx.Log(db.Close())
	}()

	return db, InitializeDatabase(ctx, db)
}

func InitializeDatabase(ctx context.Context, db *sql.DB) (err error) {
	return goosex.InitializeDatabase(ctx, db, errorsx.Must(fs.Sub(embedmigrations, ".migrations")))
}

// .migrations.cache holds the schema for cache.db: data that is rebuildable
// from external sources (tmdb/tvdb/musicbrainz/deeppool imports, ddisc
// discovery) such as library_known_media, and is therefore excluded from
// meta.db's backup/restore flow.
//
//go:embed .migrations.cache/*.sql
var embedmigrationscache embed.FS

// DatabaseCacheDatabase opens path and attaches its sibling cache.db (see
// .migrations.cache) under the "cache" catalog, so genieql-generated queries
// against cache-resident tables (e.g. library_known_media) resolve
// unqualified against whichever database this connection has open.
func DatabaseCacheDatabase(ctx context.Context, path string) (db *sql.DB, err error) {
	if path != "" {
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			return nil, err
		}
	}

	return DatabaseCacheDatabaseCustom(ctx, path, filepath.Join(filepath.Dir(path), "cache.db"))
}

// DatabaseCacheDatabaseCustom opens metapath (a duckdb DSN; "" for in-memory,
// used by tests) and attaches cachepath (see .migrations.cache) under the
// "cache" catalog. ATTACH and SET search_path are re-run via a per-connection
// init hook on every physical connection duckdb's pool opens against
// metapath, rather than once on a single shared connection, so unqualified
// references to cache-resident tables (e.g. library_known_media) resolve
// regardless of which pooled connection serves a given query.
//
// The attach is deliberately withheld (via the attachready gate) until after
// meta's own migrations have run: cache.db carries its own goose_db_version
// bookkeeping table (see migrateCacheDatabase), and attaching cache while
// meta's identically-named, not-yet-created goose_db_version is being
// established makes references to it ambiguous across the two catalogs,
// which DuckDB rejects ("a single transaction can only write to a single
// attached database"). Once meta's migrations are known-applied, that
// ambiguity can no longer arise.
func DatabaseCacheDatabaseCustom(ctx context.Context, metapath, cachepath string) (db *sql.DB, err error) {
	log.Println("database path", metapath, "cache database path", cachepath)

	if err = migrateCacheDatabase(ctx, cachepath); err != nil {
		return nil, err
	}

	var attachready atomic.Bool
	// ATTACH is instance-wide (shared by every connection duckdb pools against
	// the same file), so IF NOT EXISTS keeps re-running it on later
	// connections a no-op instead of a "already attached" error; search_path
	// is session-scoped and must be (re)set on every connection regardless.
	attach := fmt.Sprintf("ATTACH IF NOT EXISTS '%s' AS cache; SET search_path = 'main,cache';", cachepath)
	connector, err := duckdb.NewConnector(metapath, func(execer driver.ExecerContext) error {
		if !attachready.Load() {
			return nil
		}
		_, err := execer.ExecContext(context.Background(), attach, nil)
		return err
	})
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to create duckdb connector")
	}

	db = sql.OpenDB(connector)
	if err = InitializeDatabase(ctx, db); err != nil {
		debugx.Println("closing database due to error during initialization", err)
		errorsx.Log(db.Close())
		return nil, err
	}

	attachready.Store(true)
	// force any connection(s) opened during migration (which ran with the
	// gate closed, so never attached cache) out of the idle pool, so every
	// connection serving a query from here on is freshly initialized with
	// cache attached.
	db.SetMaxIdleConns(0)
	db.SetMaxIdleConns(2)

	return db, nil
}

// migrateCacheDatabase applies cache.db's own migrations, scoped to its own
// goose version table, via a standalone connection that is closed before
// returning so the file's exclusive lock is free for the caller to ATTACH it.
func migrateCacheDatabase(ctx context.Context, path string) (err error) {
	log.Println("cache database path", path)

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	db, err := sql.Open("duckdb", path)
	if err != nil {
		return errorsx.Wrap(err, "unable to open cache db")
	}
	defer func() { errorsx.Log(db.Close()) }()

	return goosex.InitializeDatabase(ctx, db, errorsx.Must(fs.Sub(embedmigrationscache, ".migrations.cache")))
}

func Checkpoint(ctx context.Context, db *sql.DB) (err error) {
	log.Println("------------------------------------------------ database checkpoint initiated ------------------------------------------------")
	defer log.Println("------------------------------------------------ database checkpoint completed ------------------------------------------------")

	if _, err := db.ExecContext(ctx, "FORCE CHECKPOINT;"); err != nil {
		return errorsx.Wrap(err, "failed to checkpoint database")
	}

	return nil
}
