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
	"github.com/retrovibed/retrovibed/shallows/internal/squirrelx"
	"github.com/retrovibed/retrovibed/shallows/library"
)

type knownenv struct {
	Database      string   `flag:"" name:"database" help:"database to read" default:"${vars_user_configuration_directory}/meta.db"`
	ExcludeSource []string `flag:"" name:"exclude-source" help:"exclude the specified source(s) from the calculation" optional:""`
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
	query, args, err := squirrelx.PSQL.
		Select("COALESCE(MAX(released), '1700-01-01'::TIMESTAMP) AS start").
		From("library_known_media").
		Where("released < NOW()").
		Where(library.KnownQueryExcludeSource(t.ExcludeSource...)).
		ToSql()
	if err != nil {
		return err
	}

	begin, err := sqlx.Timestamp(ctx, db, query, args...)
	if err != nil {
		return err
	}

	return envx.Build().Var(
		"RETROVIBED_ARCHIVE_START_DATE",
		begin.Format(time.DateOnly),
	).CopyTo(w)
}
