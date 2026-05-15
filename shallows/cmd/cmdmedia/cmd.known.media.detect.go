package cmdmedia

import (
	"bytes"
	"context"
	"database/sql"
	"io"
	"log"
	"os"

	"github.com/retrovibed/retrovibed/retroapi/deeppool"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/lucenex"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
)

type knowndetect struct {
	Database string `flag:"" name:"database" help:"database to read" default:"${vars_user_configuration_directory}/meta.db"`
	Clean    bool   `flag:"" name:"clean" help:"clean query text via deeppool before searching" default:"false"`
}

func (t knowndetect) Run(gctx *cmdopts.Global, dpc cmdopts.DeeppoolClient) (err error) {
	var db *sql.DB

	if db, err = cmdopts.DatabaseCustom(gctx.Context, t.Database); err != nil {
		return err
	}
	defer db.Close()

	cleaner := library.QueryCleaner(library.NoopQueryCleaner{})
	if t.Clean {
		httpc, err := dpc.HTTPClient(gctx.Context)
		if err != nil {
			return errorsx.Wrap(err, "unable to create deeppool client")
		}
		cleaner = deeppool.NewMediaID(httpc)
	}

	var in io.Reader = bytes.NewReader(nil)
	if cmdopts.Readable(os.Stdin) {
		in = os.Stdin
	}

	return t.run(gctx.Context, in, db, cleaner)
}

func (t knowndetect) run(ctx context.Context, in io.Reader, db sqlx.Queryer, cleaner library.QueryCleaner) error {
	type input struct {
		Query string `json:"query"`
	}

	var count int
	seq := jsonl.Iter[input](jsonl.NewDecoder(in))
	for rec := range seq.Each(ctx) {
		count++

		query, err := cleaner.Clean(ctx, rec.Query)
		if err != nil {
			log.Println("unable to clean query", err)
			query = rec.Query
		}
		query = lucenex.Clean(query)

		if rec.Query != query {
			log.Println("query cleaned", rec.Query, "->", query)
		}

		result, err := library.DetectKnownMedia(ctx, db, query)
		if err != nil {
			return err
		}

		log.Println("result", result.UID, result.Title)
	}

	if err := seq.Err(); err != nil {
		return errorsx.Wrap(err, "failed to decode input")
	}

	if count == 0 {
		return errorsx.New("query is required: stdin was empty")
	}

	return nil
}
