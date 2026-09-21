package cmdmedia

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/Masterminds/squirrel"
	"github.com/retrovibed/retrovibed/retroapi/asynccompute"
	"github.com/retrovibed/retrovibed/retroapi/fsx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/duckdbx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/lucenex"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
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

	type ScoredKnown struct {
		library.Known
		Relevance float64
	}

	// report carries the parsing stages of a query alongside its result so the printer can log them together.
	type report struct {
		Input     string
		Lucene    string
		Cleaned   string
		Title     string
		Subtitle  string
		Collation string
		Released  string
		Query     string
		Subquery  string
		Result    ScoredKnown
	}

	var count, matched uint

	// single worker so results are reported serially, matched is only touched by this worker until shutdown.
	printer := asynccompute.New(func(ctx context.Context, rpt report) error {
		var (
			buf bytes.Buffer
			c   = fsx.NewWriteErrCompact()
			tw  = tabwriter.NewWriter(&buf, 1, 0, 2, ' ', 0)
		)

		if stringsx.Present(rpt.Result.Title) {
			c.Compact(fmt.Fprintf(tw, "\n--------------------- result found ---------------------\n"))
			c.Compact(fmt.Fprintf(tw, "input:\t%s\n", rpt.Input))
			c.Compact(fmt.Fprintf(tw, "parsed title:\t'%s'\n", rpt.Title))
			c.Compact(fmt.Fprintf(tw, "parsed episode:\t'%s'\n", rpt.Subtitle))
			c.Compact(fmt.Fprintf(tw, "parsed collation:\t'%s'\n", rpt.Collation))
			c.Compact(fmt.Fprintf(tw, "parsed released:\t'%s'\n", rpt.Released))
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
			c.Compact(fmt.Fprintf(tw, "lucene:\t'%s'\n", rpt.Lucene))
			c.Compact(fmt.Fprintf(tw, "parsed title:\t'%s'\n", rpt.Title))
			c.Compact(fmt.Fprintf(tw, "parsed episode:\t'%s'\n", rpt.Subtitle))
			c.Compact(fmt.Fprintf(tw, "parsed collation:\t'%s'\n", rpt.Collation))
			c.Compact(fmt.Fprintf(tw, "parsed released:\t'%s'\n", rpt.Released))
			c.Compact(fmt.Fprintf(tw, "stripped:\t'%s' | '%s'\n", rpt.Query, rpt.Subquery))
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
		var (
			cleaned, query, subquery, release, episode string
		)

		if cleaned, err = cleaner.Clean(ctx, rec.Query); err != nil {
			log.Println("unable to clean query", err)
			query = rec.Query
		} else {
			query, subquery, release, episode = library.ParseReleaseEpisode(cleaned)
		}

		collation := library.KnownStringCollationEpisode(episode)
		rpt := report{
			Input:     rec.Query,
			Cleaned:   cleaned,
			Title:     query,
			Subtitle:  subquery,
			Collation: episode,
			Released:  release,
			Lucene:    lucenex.Clean(query),
			Query:     library.StripHallucinations(rec.Query, query),
			Subquery:  library.StripHallucinations(rec.Query, subquery),
			Result:    ScoredKnown{Relevance: t.MinRelevance},
		}

		{
			q := library.KnownSearchBuilder().Where(squirrel.And{
				library.KnownQueryExplicit(t.Explicit),
				lucenex.Query(duckdbx.NewLucene(), rpt.Lucene, lucenex.WithDefaultField("auto_description")),
			}).
				OrderByClause(library.KnownOrderCollationNearest(collation)).
				OrderByClause(library.KnownOrderReleasedNearest(library.KnownStringRelease(release))).
				OrderByClause(library.KnownOrderTitleSimilarity(rpt.Query, 0.7)).
				OrderByClause(library.KnownOrderSubtitleSimilarity(rpt.Subquery, 0.7)).
				Limit(1028)

			scanner := sqlx.Scan(library.KnownSearch(ctx, db, q))

			for v := range scanner.Iter() {
				var cur = ScoredKnown{Known: v}

				if err := library.KnownScoreByID(ctx, db, v.UID, rpt.Query, t.Cutoff).Scan(&cur.Relevance); err != nil {
					log.Println("unable to score", v.UID, err)
					continue
				}

				if cur.Relevance > rpt.Result.Relevance {
					rpt.Result = cur
				}
			}

			if err := scanner.Err(); err != nil {
				return err
			}
		}

		if stringsx.Present(rpt.Result.Title) {
			return printer.Run(ctx, rpt)
		}

		{
			terms := strings.ReplaceAll(stringsx.CompactWhitespace(rpt.Lucene), " ", " OR ")
			q := library.KnownSearchBuilder().Where(squirrel.And{
				library.KnownQueryExplicit(t.Explicit),
				lucenex.Query(duckdbx.NewLucene(), terms, lucenex.WithDefaultField("title")),
			}).
				OrderByClause(library.KnownOrderCollationNearest(collation)).
				OrderByClause(library.KnownOrderReleasedNearest(library.KnownStringRelease(release))).
				OrderByClause(library.KnownOrderTitleSimilarity(rpt.Query, 0.7)).
				OrderByClause(library.KnownOrderSubtitleSimilarity(rpt.Subquery, 0.7)).
				Limit(1028)

			scanner := sqlx.Scan(library.KnownSearch(ctx, db, q))

			for v := range scanner.Iter() {
				var cur = ScoredKnown{Known: v}

				if err := library.KnownScoreByID(ctx, db, v.UID, rpt.Query, t.Cutoff).Scan(&cur.Relevance); err != nil {
					log.Println("unable to score", v.UID, err)
					continue
				}

				if cur.Relevance > rpt.Result.Relevance {
					rpt.Result = cur
				}
			}

			if err := scanner.Err(); err != nil {
				return err
			}
		}

		return printer.Run(ctx, rpt)
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
