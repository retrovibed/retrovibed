package cmdmedia

import (
	"context"
	"database/sql"
	"log"

	"github.com/retrovibed/retrovibed/shallows/cmd/cmdmeta"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
)

type knowndetect struct {
	Database string `flag:"" name:"database" help:"database to read" default:"${vars_user_configuration_directory}/meta.db"`
	Query    string `arg:"" name:"query" help:"query to detect" required:"true"`
}

func (t knowndetect) Run(gctx *cmdopts.Global) (err error) {
	var db *sql.DB

	if db, err = cmdmeta.DatabaseCustom(gctx.Context, t.Database); err != nil {
		return err
	}
	defer db.Close()

	result, err := t.run(gctx.Context, db, t.Query)
	if err != nil {
		return err
	}

	log.Println("result", result.UID, result.Title)
	return nil
}

func (t knowndetect) run(ctx context.Context, db sqlx.Queryer, query string) (library.Known, error) {
	return library.DetectKnownMedia(ctx, db, query)
}
