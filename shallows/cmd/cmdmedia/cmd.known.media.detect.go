package cmdmedia

import (
	"bytes"
	"context"
	"database/sql"
	"io"
	"log"
	"os"

	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
)

type knowndetect struct {
	Database string `flag:"" name:"database" help:"cache database to read" default:"${vars_user_cache_directory}/cache.db"`
}

func (t knowndetect) Run(gctx *cmdopts.Global) (err error) {
	var db *sql.DB

	if db, err = cmdopts.DatabaseCache(gctx.Context, t.Database); err != nil {
		return err
	}
	defer db.Close()

	cleaner := library.NewQueryerCleanerAuto()

	var in io.Reader = bytes.NewReader(nil)
	if cmdopts.Readable(os.Stdin) {
		in = os.Stdin
	}

	return t.run(gctx.Context, in, db, cleaner)
}

func (t knowndetect) run(ctx context.Context, in io.Reader, db sqlx.Queryer, cleaner library.QueryCleaner) error {
	type input struct {
		Query    string `json:"query"`
		Mimetype string `json:"mimetype"`
	}

	var count int
	identifier := library.NewKnownIdentifier(db, cleaner)
	seq := jsonl.Iter[input](jsonl.NewDecoder(in))
	for rec := range seq.Each(ctx) {
		count++

		result, err := identifier.Identify(ctx, mimex.Category(rec.Mimetype), rec.Query)
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
