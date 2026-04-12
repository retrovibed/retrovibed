package cmdmeta

import (
	"context"
	"database/sql"
	"embed"
	"io/fs"
	"log"
	"os"
	"strings"
	"sync"

	_ "github.com/duckdb/duckdb-go/v2"

	"github.com/pressly/goose/v3"
	"github.com/retrovibed/retrovibed/internal/debugx"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/goosex"
	"github.com/retrovibed/retrovibed/internal/sqlx"
	"github.com/retrovibed/retrovibed/internal/userx"
	"github.com/retrovibed/retrovibed/meta"
)

//go:embed .migrations/*.sql
var embedsqlite embed.FS

func Database(ctx context.Context) (db *sql.DB, err error) {
	log.Println("database path", userx.DefaultConfigDir(userx.DefaultRelRoot(), "meta.db"))

	if err := os.MkdirAll(userx.DefaultConfigDir(userx.DefaultRelRoot()), 0700); err != nil {
		return nil, err
	}

	if db, err = sql.Open("duckdb", userx.DefaultConfigDir(userx.DefaultRelRoot(), "meta.db")); err != nil {
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

var m sync.Mutex

func InitializeDatabase(ctx context.Context, db *sql.DB) (err error) {
	m.Lock()
	defer m.Unlock()

	if _, err := db.ExecContext(ctx, "CHECKPOINT;"); err != nil {
		return errorsx.Wrap(err, "failed to checkpoint database")
	}

	if version, err := sqlx.String(ctx, db, "SELECT library_version FROM pragma_version()"); err != nil {
		return err
	} else {
		log.Println("detected duckdb", version)
	}

	if _, err := db.ExecContext(ctx, inetSQL); err != nil {
		return errorsx.Wrap(err, "failed to load inet extension")
	}

	mprov, err := goose.NewProvider(
		"",
		db,
		errorsx.Must(fs.Sub(embedsqlite, ".migrations")),
		goose.WithStore(goosex.DuckdbStore{}),
		goose.WithAllowOutofOrder(true),
		// goose.WithVerbose(true),
	)
	if err != nil {
		return errorsx.Wrap(err, "unable to build migration provider")
	}

	if _, err := mprov.Up(ctx); err != nil {
		return errorsx.Wrap(err, "unable to run migrations")
	}

	if _, err := db.ExecContext(ctx, "CHECKPOINT;"); err != nil {
		return errorsx.Wrap(err, "failed to checkpoint database")
	}

	return nil
}

func Checkpoint(ctx context.Context, db *sql.DB) (err error) {
	log.Println("------------------------------------------------ database checkpoint initiated ------------------------------------------------")
	defer log.Println("------------------------------------------------ database checkpoint completed ------------------------------------------------")

	if _, err := db.ExecContext(ctx, "CHECKPOINT;"); err != nil {
		return errorsx.Wrap(err, "failed to checkpoint database")
	}

	return nil
}

func Hostnames(ctx context.Context, q sqlx.Queryer) (res []string, _ error) {
	s := sqlx.Scan(meta.DaemonSearch(ctx, q, meta.DaemonSearchBuilder().Limit(128)))
	for d := range s.Iter() {
		before, _, _ := strings.Cut(d.Hostname, ":")
		res = append(res, before)
	}

	return res, s.Err()
}
