package cmdmedia

import (
	"database/sql"
	"log"
	"os"
	"sync/atomic"

	"github.com/retrovibed/retrovibed/shallows/cmd/cmdmeta"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
)

type duckdbexport struct {
	Database string `flag:"" name:"database" help:"database to read" default:"${vars_user_configuration_directory}/meta.db"`
}

func (t duckdbexport) Run(gctx *cmdopts.Global) (err error) {
	var (
		db       *sql.DB
		progress uint64
	)

	if db, err = cmdmeta.DatabaseCustom(gctx.Context, t.Database); err != nil {
		return err
	}
	defer db.Close()

	d := jsonl.NewEncoder(os.Stdout)

	q := sqlx.Scan(library.KnownSearch(gctx.Context, db, library.KnownSearchBuilder()))

	for v := range q.Iter() {
		if err := d.Encode(v); err != nil {
			return err
		}

		if atomic.AddUint64(&progress, 1)%100 == 0 {
			log.Println("exported", progress, "records")
			log.Println("current", v.ID, v.UID, v.Title)
		}
	}

	return q.Err()
}
