package cmdmedia

import (
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
	var (
		db   *sql.DB
		v    library.Known
		derr error
	)

	if db, err = cmdmeta.DatabaseCustom(gctx.Context, t.Database); err != nil {
		return err
	}
	defer db.Close()

	d := jsonl.NewDecoder(os.Stdin)

	for derr = d.Decode(&v); derr == nil; derr = d.Decode(&v) {
		v.AutoDescription = stringsx.Join("\n", v.Title, v.OriginalTitle, v.Overview)
		if err = library.KnownInsertWithDefaults(gctx.Context, db, v).Scan(&v); err != nil {
			return err
		}
	}

	if _, err := db.ExecContext(gctx.Context, "CHECKPOINT"); err != nil {
		return errorsx.Wrap(err, "failed to checkpoint database")
	}

	return errorsx.Ignore(derr, io.EOF)
}
