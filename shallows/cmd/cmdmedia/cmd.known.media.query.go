package cmdmedia

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"text/tabwriter"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/asynccompute"
	"github.com/retrovibed/retrovibed/retroapi/fsx"
	"github.com/retrovibed/retrovibed/retroapi/uuidx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/library"
)

type knownquery struct {
	Database     string  `flag:"" name:"database" help:"cache database to read" default:"${vars_user_cache_directory}/cache.db"`
	Explicit     bool    `flag:"" name:"explicit" help:"include explicit content in results" default:"false"`
	Cutoff       float32 `flag:"" name:"cutoff" help:"similarity cutoff for scoring" default:"0.7"`
	MinRelevance float64 `flag:"" name:"min-relevance" help:"minimum relevance score" default:"0.85"`
	Backlog      uint16  `flag:"" name:"backlog" help:"number of queries allowed to queue up" default:"128"`
	Workers      uint16  `flag:"" name:"workers" help:"number of concurrent query workers to run" default:"8"`
}

func (t knownquery) Run(gctx *cmdopts.Global) (err error) {
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

func (t knownquery) run(ctx context.Context, in io.Reader, db *sql.DB, cleaner library.QueryCleaner) (err error) {
	type input struct {
		Query string `json:"query"`
	}

	type report struct {
		Input  string
		Result library.KnownScored
	}

	var count, matched uint

	identifier := library.NewKnownIdentifier(db, cleaner)
	identifier.Cutoff = t.Cutoff
	identifier.MinRelevance = t.MinRelevance
	identifier.Explicit = t.Explicit

	// single worker so results are reported serially, matched is only touched by this worker until shutdown.
	printer := asynccompute.New(func(ctx context.Context, rpt report) error {
		var (
			buf bytes.Buffer
			c   = fsx.NewWriteErrCompact()
			tw  = tabwriter.NewWriter(&buf, 1, 0, 2, ' ', 0)
		)

		if !uuidx.IsMinMax(uuid.FromStringOrNil(rpt.Result.UID)) {
			c.Compact(fmt.Fprintf(tw, "\n--------------------- result found ---------------------\n"))
			c.Compact(fmt.Fprintf(tw, "input:\t%s\n", rpt.Input))
			c.Compact(fmt.Fprintf(tw, "result relevance:\t%v\n", rpt.Result.Relevance))
			c.Compact(fmt.Fprintf(tw, "result uid:\t%s\n", rpt.Result.UID))
			c.Compact(fmt.Fprintf(tw, "result title:\t%s\n", rpt.Result.Title))
			c.Compact(fmt.Fprintf(tw, "result collation:\t%s\n", library.KnownCollationString(rpt.Result.Collation)))
			c.Compact(fmt.Fprintf(tw, "result subtitle:\t%s\n", rpt.Result.Subtitle))
			c.Compact(fmt.Fprintf(tw, "result released:\t%v\n", rpt.Result.Released))
			c.Compact(fmt.Fprintf(tw, "result mimetype:\t%s\n", rpt.Result.Mimetype))
			c.Compact(fmt.Fprintf(tw, "--------------------------------------------------------\n"))
			matched += 1
		} else {
			c.Compact(fmt.Fprintf(tw, "\n-------------------- no result found -------------------\n"))
			c.Compact(fmt.Fprintf(tw, "input:\t'%s'\n", rpt.Input))
			c.Compact(fmt.Fprintf(tw, "--------------------------------------------------------\n"))
		}

		if err := errorsx.Compact(tw.Flush(), c.Err()); err != nil {
			return err
		}

		// single write so a report is never interleaved with other log output.
		log.Print(buf.String())

		return nil
	}, asynccompute.Backlog[report](t.Backlog), asynccompute.Workers[report](1))

	queries := asynccompute.New(func(ctx context.Context, rec input) (err error) {
		var res library.KnownScored

		if res, err = identifier.Identify(ctx, rec.Query); err != nil {
			return err
		}

		return printer.Run(ctx, report{Input: rec.Query, Result: res})
	}, asynccompute.Backlog[input](t.Backlog), asynccompute.Workers[input](t.Workers))

	defer func() {
		err = errorsx.Compact(err, asynccompute.Shutdown(ctx, queries, printer))
	}()

	seq := jsonl.Iter[input](jsonl.NewDecoder(in))
	for rec := range seq.Each(ctx) {
		count++

		if err := queries.Run(ctx, rec); err != nil {
			return err
		}
	}

	// drain order: queries feed the printer.
	if err := errorsx.Compact(errorsx.Wrap(seq.Err(), "failed to decode input"), asynccompute.Shutdown(ctx, queries, printer)); err != nil {
		return err
	}

	if count == 0 {
		return errorsx.New("query is required: stdin was empty")
	}

	log.Println("completed - matched", matched, "/", count)

	return nil
}
