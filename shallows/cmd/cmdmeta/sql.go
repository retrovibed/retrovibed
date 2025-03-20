package cmdmeta

import (
	"context"
	"database/sql"
	"embed"
	"io/fs"

	"github.com/pressly/goose/v3"
	"github.com/retrovibed/retrovibed/internal/x/debugx"
	"github.com/retrovibed/retrovibed/internal/x/errorsx"
	"github.com/retrovibed/retrovibed/internal/x/goosex"
	"github.com/retrovibed/retrovibed/internal/x/userx"
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

	{
		mprov, err := goose.NewProvider("", db, errorsx.Must(fs.Sub(embedsqlite, ".migrations")), goose.WithStore(goosex.DuckdbStore{}))
		if err != nil {
			return nil, errorsx.Wrap(err, "unable to build migration provider")
		}

		if _, err := mprov.Up(ctx); err != nil {
			return nil, errorsx.Wrap(err, "unable to run migrations")
		}
	}

	return db, nil
}
