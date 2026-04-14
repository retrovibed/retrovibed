package cmdmedia

import (
	"database/sql"
	"os"

	"github.com/retrovibed/retrovibed/shallows/cmd/cmdmeta"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
)

type duckdbimport struct{}

func (t duckdbimport) Run(gctx *cmdopts.Global) (err error) {
	var (
		db *sql.DB
	)

	if db, err = cmdmeta.Database(gctx.Context); err != nil {
		return err
	}
	defer db.Close()

	d := jsonl.NewEncoder(os.Stdout)

	q := sqlx.Scan(library.KnownSearch(gctx.Context, db, library.KnownSearchBuilder()))

	for v := range q.Iter() {
		if err := d.Encode(v); err != nil {
			return err
		}
	}

	return q.Err()
}
