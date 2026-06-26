package cmdmedia

import (
	"context"
	"database/sql"
	"io"
	"log"
	"os"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/asynccompute"
	"github.com/retrovibed/retrovibed/shallows/internal/backoffx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/slicesx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/library"
)

type knownimport struct {
	Database string `flag:"" name:"database" help:"database to read" default:"${vars_user_configuration_directory}/meta.db"`
	Batch    int    `flag:"" name:"batch" help:"number of records to insert per batch" default:"8192"`
	Backlog  uint16 `flag:"" name:"backlog" help:"number of batches to allowed to queue up" default:"128"`
	Workers  uint16 `flag:"" name:"workers" help:"number of async database workers to run" default:"1"`
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
	type batch []library.Known
	inserts := asynccompute.New(func(ctx context.Context, chunk batch) (err error) {
		bs := backoffx.New(
			backoffx.Constant(time.Second),
			backoffx.JitterRandom(200*time.Millisecond),
		)
		return backoffx.Attempt(ctx, bs, func(ctx context.Context) error {
			ts := time.Now()
			s := library.NewKnownBatchInsertWithDefaults(ctx, db, chunk...)

			if err := sqlx.Discard(sqlx.Scan(s)); err != nil {
				return errorsx.Wrap(err, "failed to insert batch")
			}

			log.Println("imported", time.Since(ts), len(chunk), "records")
			return nil
		})
	}, asynccompute.Backlog[batch](t.Backlog), asynccompute.Workers[batch](t.Workers))

	d := jsonl.Iter[library.Known](jsonl.NewDecoder(r))

	for chunk := range iterx.Chunk(d.Each(ctx), t.Batch) {
		owned := slicesx.Map(func(v library.Known) library.Known {
			v.AutoDescription = stringsx.Join("\n", v.Title, v.OriginalTitle, v.Overview)
			return v
		}, append(batch{}, chunk...)...)

		if err = inserts.Run(ctx, owned); err != nil {
			return errorsx.Compact(err, asynccompute.Shutdown(ctx, inserts))
		}
	}

	if err = errorsx.Compact(d.Err(), asynccompute.Shutdown(ctx, inserts)); err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, "CHECKPOINT"); err != nil {
		return errorsx.Wrap(err, "failed to checkpoint database")
	}

	return nil
}
