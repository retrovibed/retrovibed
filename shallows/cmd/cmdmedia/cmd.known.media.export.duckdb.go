package cmdmedia

import (
	"database/sql"
	"log"
	"os"
	"sync/atomic"

	"github.com/Masterminds/squirrel"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
)

type duckdbexport struct {
	Database string   `flag:"" name:"database" help:"database to read" default:"${vars_user_configuration_directory}/meta.db"`
	Offset   string   `flag:"" name:"offset" help:"uid offset to start export from" default:"00000000-0000-0000-0000-000000000000"`
	ID       []string `flag:"" name:"id" help:"restrict export to the specified known id(s)" optional:""`
	Language string   `flag:"" name:"language" help:"restrict export to the specified language" optional:""`
	Mimetype string   `flag:"" name:"mimetype" help:"restrict export to the specified mimetype" optional:""`
	Source   []string `flag:"" name:"source" help:"restrict export to the specified source(s)" optional:""`
	Explicit bool     `flag:"" name:"explicit" help:"include explicit content in export" default:"false"`
	Limit    uint64   `flag:"" name:"limit" help:"maximum number of records to export (0 = unlimited)" default:"0"`
}

func (t duckdbexport) Run(gctx *cmdopts.Global) (err error) {
	var (
		db       *sql.DB
		progress uint64
	)

	if db, err = cmdopts.DatabaseCustom(gctx.Context, t.Database); err != nil {
		return err
	}
	defer db.Close()

	d := jsonl.NewEncoder(os.Stdout)

	b := library.KnownSearchBuilder().OrderBy("uid ASC").Where(squirrel.And{
		library.KnownQueryUIDGreaterThan(t.Offset),
		library.KnownQueryUID(t.ID...),
		library.KnownQueryLanguage(t.Language),
		library.KnownQueryMimetype(t.Mimetype),
		library.KnownQuerySource(t.Source...),
		library.KnownQueryExplicit(t.Explicit),
	}).Limit(t.Limit)

	q := sqlx.Scan(library.KnownSearch(gctx.Context, db, b))

	for v := range q.Iter() {
		if err := d.Encode(v); err != nil {
			return err
		}

		if atomic.AddUint64(&progress, 1)%8192 == 0 {
			log.Println("exported", progress, "records")
			log.Println("current", v.ID, v.UID, v.Title)
		}
	}

	return q.Err()
}
