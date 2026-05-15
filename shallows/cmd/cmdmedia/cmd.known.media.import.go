package cmdmedia

import (
	"context"
	"database/sql"
	"io"
	"log"
	"os"
	"time"

	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/iterx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/slicesx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/library"
)

type knownimport struct {
	Database string `flag:"" name:"database" help:"database to read" default:"${vars_user_configuration_directory}/meta.db"`
	Batch    int    `flag:"" name:"batch" help:"number of records to insert per batch" default:"8192"`
}

func (t knownimport) Run(gctx *cmdopts.Global) (err error) {
	var db *sql.DB

	if db, err = cmdopts.DatabaseCustom(gctx.Context, t.Database); err != nil {
		return err
	}
	defer db.Close()

	return t.run(gctx.Context, db, os.Stdin)
}

func (t knownimport) run(ctx context.Context, db *sql.DB, r io.Reader) (err error) {
	d := jsonl.Iter[library.Known](jsonl.NewDecoder(r))
	for chunk := range iterx.Chunk(d.Each(ctx), t.Batch) {
		chunk = slicesx.Map(func(v library.Known) library.Known {
			v.AutoDescription = stringsx.Join("\n", v.Title, v.OriginalTitle, v.Overview)
			return v
		}, chunk...)

		ts := time.Now()
		s := library.NewKnownBatchInsertWithDefaults(ctx, db, chunk...)
		for s.Next() {
			var v library.Known
			if err = s.Scan(&v); err != nil {
				return errorsx.Wrap(err, "failed to scan inserted record")
			}
		}

		if err = s.Err(); err != nil {
			return errorsx.Wrap(err, "failed to insert batch")
		}

		if err = s.Close(); err != nil {
			return errorsx.Wrap(err, "failed to close batch")
		}

		log.Println("imported", time.Since(ts), len(chunk), "records")
	}

	if err = d.Err(); err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, "CHECKPOINT"); err != nil {
		return errorsx.Wrap(err, "failed to checkpoint database")
	}

	return nil
}
