package cmdopts

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"time"

	duckdb "github.com/duckdb/duckdb-go/v2"
	"github.com/egdaemon/duckproxy/duckshare"

	"github.com/retrovibed/retrovibed/retroapi/debugx"
	"github.com/retrovibed/retrovibed/retroapi/errorsx"
	"github.com/retrovibed/retrovibed/retroapi/userx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/goosex"
)

//go:embed .migrations/*.sql
var embedmigrations embed.FS

func DatabaseMeta(ctx context.Context) (db *sql.DB, err error) {
	return DatabaseMetaCustom(ctx, userx.DefaultConfigDir(userx.DefaultRelRoot(), "meta.db"))
}

// Database opens meta.db, opens and migrates its sibling cache.db, and attaches
// the latter onto the former under the "cache" catalog, returning the single
// connection callers of library.Known* (which query cache.<table>) need.
func Database(ctx context.Context) (db *sql.DB, err error) {
	return DatabaseCustom(
		ctx,
		userx.DefaultConfigDir(userx.DefaultRelRoot(), "meta.db"),
		userx.DefaultCacheDirectory(userx.DefaultRelRoot(), "cache.db"),
	)
}

// DatabaseCustom opens metapath, opens and migrates cachepath, and attaches
// cachepath onto the metapath connection under the "cache" catalog. See
// DatabaseCache.
//
// cache.db is never opened directly: duckdb only allows a single process to have
// a file open, so the meta.db connector is shared between processes (see
// shared) and the cache is attached to, and migrated through, that one
// instance.
func DatabaseCustom(ctx context.Context, metapath, cachepath string) (db *sql.DB, err error) {
	// attaching doesn't create the directory of the cache.
	if err = fsx.MkDirs(0700, filepath.Dir(metapath), filepath.Dir(cachepath)); err != nil {
		return nil, err
	}

	// the cache is attached by the connector, not once after opening, because the
	// attachment is lost when the process owning the shared file exits and
	// another takes over, see duckshare.WithInit.
	connector, err := shared(metapath, fmt.Sprintf("ATTACH IF NOT EXISTS '%s' AS cache;", cachepath))
	if err != nil {
		return nil, err
	}

	if db, err = initialize(ctx, connector, errorsx.Must(fs.Sub(embedmigrations, ".migrations"))); err != nil {
		return nil, err
	}

	defer func() {
		if err == nil {
			return
		}
		debugx.Println("closing database due to error during initialization", err)
		errorsx.Log(db.Close())
	}()

	// the version table is qualified because goose checks for pending migrations
	// on a connection whose default catalog is meta. USE makes the unqualified
	// statements of the migrations, and the migration lock, apply to the cache.
	err = goosex.InitializeDatabase(
		ctx,
		db,
		errorsx.Must(fs.Sub(embedmigrationscache, ".migrations.cache")),
		goosex.WithTableName("cache.goose_db_version"),
		goosex.WithPreStatements("USE cache;"),
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func DatabaseMetaCustom(ctx context.Context, path string) (db *sql.DB, err error) {
	return InitializeDatabase(ctx, path, errorsx.Must(fs.Sub(embedmigrations, ".migrations")))
}

// shared returns a connector for the duckdb file at path which lets concurrent
// processes use it, the first to connect owns the file and serves it to the
// others (see duckshare). init statements run on every connection the owner
// opens to the file, they must be safe to repeat. an empty path is a private
// in-memory database.
func shared(path string, init ...string) (driver.Connector, error) {
	log.Printf("database path: '%s'\n", path)

	// state held by the duckdb instance, restored whenever a process takes the
	// shared file over.
	init = append([]string{
		extensionSQL,
		"SET GLOBAL hnsw_enable_experimental_persistence = true;",
		// debugging.
		// "CALL enable_logging('QueryLog', storage_path = '/tmp/retrovibed.duckdb.log', storage_buffer_size = 0);",
	}, init...)

	if path == "" {
		return duckdb.NewConnector("", func(execer driver.ExecerContext) error {
			for _, stmt := range init {
				if _, err := execer.ExecContext(context.Background(), stmt, nil); err != nil {
					return errorsx.Wrap(err, "unable to initialize the connection")
				}
			}

			return nil
		})
	}

	if err := fsx.MkDirs(0700, filepath.Dir(path)); err != nil {
		return nil, err
	}

	return duckshare.New(userx.DefaultRuntimeDirectory(userx.DefaultRelRoot()), path, duckshare.WithInit(init...)), nil
}

// InitializeDatabase opens path, creating its parent directory if necessary, and
// applies migrations. Shared by DatabaseMetaCustom (meta.db) and tests so
// databases are opened and initialized identically, differing only in which
// migration set they apply.
func InitializeDatabase(ctx context.Context, path string, migrations fs.FS, opts ...goosex.Option) (db *sql.DB, err error) {
	connector, err := shared(path)
	if err != nil {
		return nil, err
	}

	return initialize(ctx, connector, migrations, opts...)
}

// initialize opens connector and applies migrations.
func initialize(ctx context.Context, connector driver.Connector, migrations fs.FS, opts ...goosex.Option) (db *sql.DB, err error) {
	db = sql.OpenDB(connector)

	// duckdb-go caches returned connections as idle rather than actually closing
	// them, which leaves their underlying DuckDB session (and whatever MVCC/index
	// state it's still pinning) alive indefinitely. disable idle caching so every
	// connection is genuinely closed once returned, and bound how long any single
	// connection may be reused as defense in depth against the same staleness.
	db.SetMaxIdleConns(0)
	db.SetConnMaxLifetime(time.Minute)

	defer func() {
		if err == nil {
			return
		}
		debugx.Println("closing database due to error during initialization", err)
		errorsx.Log(db.Close())
	}()

	if err = Checkpoint(ctx, db); err != nil {
		return db, err
	}

	if err = goosex.InitializeDatabase(ctx, db, migrations, opts...); err != nil {
		return db, err
	}

	return db, Checkpoint(ctx, db)
}

// .migrations.cache holds the schema for cache.db: data that is rebuildable
// from external sources (tmdb/tvdb/musicbrainz/deeppool imports, ddisc
// discovery) such as library_known_media, and is therefore excluded from
// meta.db's backup/restore flow.
//
//go:embed .migrations.cache/*.sql
var embedmigrationscache embed.FS

// DatabaseCache opens and migrates path (see .migrations.cache) and attaches it
// under the "cache" catalog of the default meta.db. The cache is only ever
// opened through meta.db, so concurrent commands (and the daemon) share the one
// instance rather than contending for the cache file. Cache-resident tables
// (e.g. library_known_media) are referenced fully qualified as cache.<table>
// (or "cache"."<table>" for INSERT/UPDATE/DELETE targets) throughout the
// genieql-generated code. A database's catalog is named after its file, so
// attaching from a separate instance is what makes those queries resolve for
// any path rather than only one named cache.db.
func DatabaseCache(ctx context.Context, path string) (db *sql.DB, err error) {
	return DatabaseCustom(ctx, userx.DefaultConfigDir(userx.DefaultRelRoot(), "meta.db"), path)
}

func Checkpoint(ctx context.Context, db *sql.DB) (err error) {
	log.Println("------------------------------------------------ database checkpoint initiated ------------------------------------------------")
	defer log.Println("------------------------------------------------ database checkpoint completed ------------------------------------------------")

	if _, err := db.ExecContext(ctx, "FORCE CHECKPOINT;"); err != nil {
		return errorsx.Wrap(err, "failed to checkpoint database")
	}

	return nil
}
