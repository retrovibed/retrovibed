package main

import (
	"database/sql"
	"embed"
	"io/fs"

	"github.com/pressly/goose/v3"

	"github.com/james-lawrence/deeppool/cmd/cmdopts"
	"github.com/james-lawrence/deeppool/downloads"
	"github.com/james-lawrence/deeppool/internal/x/errorsx"
	"github.com/james-lawrence/deeppool/internal/x/goosex"
	"github.com/james-lawrence/deeppool/internal/x/userx"

	_ "github.com/marcboeker/go-duckdb"
)

//go:embed .migrations/*.sql
var embedsqlite embed.FS

type cmdDaemon struct{}

func (t cmdDaemon) Run(ctx *cmdopts.Global) (err error) {
	{
		db, err := sql.Open("duckdb", "dpool.db")
		if err != nil {
			return errorsx.Wrap(err, "unable to open db")
		}
		defer db.Close()

		mprov, err := goose.NewProvider("", db, errorsx.Must(fs.Sub(embedsqlite, ".migrations")), goose.WithStore(goosex.DuckdbStore{}))
		if err != nil {
			return errorsx.Wrap(err, "unable to build migration provider")
		}

		if _, err := mprov.Up(ctx.Context); err != nil {
			return errorsx.Wrap(err, "unable to run migrations")
		}
	}

	dwatcher, err := downloads.NewDirectoryWatcher()
	if err != nil {
		return errorsx.Wrap(err, "unable to setup directory monitoring for torrents")
	}

	if err = dwatcher.Add(userx.DefaultDownloadDirectory()); err != nil {
		return errorsx.Wrap(err, "unable to add the download directory to be watched")
	}

	<-ctx.Context.Done()
	return nil
}
