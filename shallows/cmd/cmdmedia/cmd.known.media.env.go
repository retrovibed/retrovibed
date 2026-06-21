package cmdmedia

import (
	"context"
	"database/sql"
	"io"
	"os"
	"time"

	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/envx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
)

type knownenv struct {
	Database string `flag:"" name:"database" help:"database to read" default:"${vars_user_configuration_directory}/meta.db"`
}

func (t knownenv) Run(gctx *cmdopts.Global) (err error) {
	var db *sql.DB

	if db, err = cmdopts.DatabaseCustom(gctx.Context, t.Database); err != nil {
		return err
	}
	defer db.Close()

	return t.run(gctx.Context, db, os.Stdout)
}

func (t knownenv) run(ctx context.Context, db *sql.DB, w io.Writer) (err error) {
	begin, err := sqlx.Timestamp(ctx, db, `SELECT COALESCE(MAX(released), '-infinity'::TIMESTAMP) AS start FROM library_known_media WHERE released < NOW()`)
	if err != nil {
		return err
	}

	return envx.Build().Var(
		"RETROVIBED_ARCHIVE_START_DATE",
		begin.Format(time.DateOnly),
	).CopyTo(w)
}
