package cmdmeta

import (
	"context"
	"database/sql"
	"embed"
	"io/fs"
	"strings"

	_ "github.com/marcboeker/go-duckdb/v2"

	"github.com/pressly/goose/v3"
	"github.com/retrovibed/retrovibed/internal/debugx"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/goosex"
	"github.com/retrovibed/retrovibed/internal/slicesx"
	"github.com/retrovibed/retrovibed/internal/sqlx"
	"github.com/retrovibed/retrovibed/internal/sqlxx"
	"github.com/retrovibed/retrovibed/internal/userx"
	"github.com/retrovibed/retrovibed/meta"
)

//go:embed .migrations/*.sql
var embedsqlite embed.FS

func Database(ctx context.Context) (db *sql.DB, err error) {
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

func InitializeDatabase(ctx context.Context, db *sql.DB) (err error) {
	mprov, err := goose.NewProvider("", db, errorsx.Must(fs.Sub(embedsqlite, ".migrations")), goose.WithStore(goosex.DuckdbStore{}))
	if err != nil {
		return errorsx.Wrap(err, "unable to build migration provider")
	}

	if _, err := mprov.Up(ctx); err != nil {
		return errorsx.Wrap(err, "unable to run migrations")
	}

	return nil
}

func Hostnames(ctx context.Context, q sqlx.Queryer) ([]string, error) {
	var (
		results []meta.Daemon
	)

	if err := sqlxx.ScanInto(meta.DaemonSearch(ctx, q, meta.DaemonSearchBuilder().Limit(128)), &results); err != nil {
		return nil, errorsx.Wrap(err, "unable to retrieve hostnames")
	}

	return slicesx.MapTransform(func(d meta.Daemon) string {
		before, _, _ := strings.Cut(d.Hostname, ":")
		return before
	}, results...), nil
}
