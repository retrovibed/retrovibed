package cmdopts

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/duckdb/duckdb-go/v2"

	"github.com/retrovibed/retrovibed/retroapi/userx"
	"github.com/retrovibed/retrovibed/shallows/internal/debugx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
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
// AttachCache and DatabaseCache.
func DatabaseCustom(ctx context.Context, metapath, cachepath string) (db *sql.DB, err error) {
	if db, err = DatabaseMetaCustom(ctx, metapath); err != nil {
		return nil, err
	}

	cachedb, err := migrateCache(ctx, cachepath)
	if err != nil {
		return nil, err
	}

	if err = cachedb.Close(); err != nil {
		return nil, err
	}

	if err = AttachCache(ctx, db, "cache", cachepath); err != nil {
		return nil, err
	}

	return db, nil
}

func DatabaseMetaCustom(ctx context.Context, path string) (db *sql.DB, err error) {
	return InitializeDatabase(ctx, path, errorsx.Must(fs.Sub(embedmigrations, ".migrations")))
}

// InitializeDatabase opens path, creating its parent directory if necessary, and
// applies migrations. Shared by DatabaseCustom (meta.db) and DatabaseCache
// (cache.db) so the two databases are opened and initialized identically,
// differing only in which migration set they apply.
func InitializeDatabase(ctx context.Context, path string, migrations fs.FS) (db *sql.DB, err error) {
	log.Println("database path", path)

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, err
	}

	if db, err = sql.Open("duckdb", path); err != nil {
		return nil, errorsx.Wrap(err, "unable to open db")
	}

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

	return db, goosex.InitializeDatabase(ctx, db, migrations)
}

// .migrations.cache holds the schema for cache.db: data that is rebuildable
// from external sources (tmdb/tvdb/musicbrainz/deeppool imports, ddisc
// discovery) such as library_known_media, and is therefore excluded from
// meta.db's backup/restore flow.
//
//go:embed .migrations.cache/*.sql
var embedmigrationscache embed.FS

// DatabaseCache opens and migrates path (see .migrations.cache) and attaches it
// under the "cache" catalog of an otherwise empty in-memory database, with no
// dependency on a meta database. Cache-resident tables (e.g.
// library_known_media) are referenced fully qualified as cache.<table> (or
// "cache"."<table>" for INSERT/UPDATE/DELETE targets) throughout the
// genieql-generated code. A database's catalog is named after its file, so
// attaching from a separate instance is what makes those queries resolve for
// any path rather than only one named cache.db.
func DatabaseCache(ctx context.Context, path string) (db *sql.DB, err error) {
	return DatabaseCustom(ctx, "", path)
}

// migrateCache opens and migrates path (see .migrations.cache) as its own
// standalone database, whose catalog is named after its file.
func migrateCache(ctx context.Context, path string) (db *sql.DB, err error) {
	return InitializeDatabase(ctx, path, errorsx.Must(fs.Sub(embedmigrationscache, ".migrations.cache")))
}

// AttachCache attaches cachepath (see DatabaseCache) onto db under the
// "cache" catalog. ATTACH is instance-wide — shared by every connection
// duckdb's pool opens against the same file, not just the one that issued it
// — so this only needs to run once per db; IF NOT EXISTS keeps a repeat call
// a no-op instead of an "already attached" error.
func AttachCache(ctx context.Context, db *sql.DB, name, path string) (err error) {
	attach := fmt.Sprintf("ATTACH IF NOT EXISTS '%s' AS %s;", path, name)
	if _, err = db.ExecContext(ctx, attach); err != nil {
		return errorsx.Wrap(err, "unable to attach cache database")
	}

	return nil
}

func Checkpoint(ctx context.Context, db *sql.DB) (err error) {
	log.Println("------------------------------------------------ database checkpoint initiated ------------------------------------------------")
	defer log.Println("------------------------------------------------ database checkpoint completed ------------------------------------------------")

	if _, err := db.ExecContext(ctx, "FORCE CHECKPOINT;"); err != nil {
		return errorsx.Wrap(err, "failed to checkpoint database")
	}

	return nil
}
