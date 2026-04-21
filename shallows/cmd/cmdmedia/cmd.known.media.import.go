package cmdmedia

import (
	"context"
	"database/sql"
	"io"
	"os"

	"github.com/retrovibed/retrovibed/shallows/cmd/cmdmeta"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/library"
)

type knownimport struct {
	Database string `flag:"" name:"database" help:"database to read" default:"${vars_user_configuration_directory}/meta.db"`
}

func (t knownimport) Run(gctx *cmdopts.Global) (err error) {
	var db *sql.DB

	if db, err = cmdmeta.DatabaseCustom(gctx.Context, t.Database); err != nil {
		return err
	}
	defer db.Close()

	return t.run(gctx.Context, db, os.Stdin)
}

func (t knownimport) run(ctx context.Context, db *sql.DB, r io.Reader) (err error) {
	d := jsonl.Iter[library.Known](jsonl.NewDecoder(r))
	for v := range d.Each(ctx) {
		v.AutoDescription = stringsx.Join("\n", v.Title, v.OriginalTitle, v.Overview)
		if err = library.KnownInsertWithDefaults(ctx, db, v).Scan(&v); err != nil {
			return err
		}
	}

	if err = d.Err(); err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, "CHECKPOINT"); err != nil {
		return errorsx.Wrap(err, "failed to checkpoint database")
	}

	return nil
}
