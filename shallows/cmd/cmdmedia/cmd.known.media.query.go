package cmdmedia

import (
	"bytes"
	"context"
	"database/sql"
	"io"
	"log"
	"os"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/retrovibed/retrovibed/retroapi/unsafepretty"
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

func (t knownquery) run(ctx context.Context, in io.Reader, db *sql.DB, cleaner library.QueryCleaner) error {
	type input struct {
		Query string `json:"query"`
	}

	type ScoredKnown struct {
		library.Known
		Relevance float64
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

		log.Println("clean result", unsafepretty.Print(query, unsafepretty.OptionDisplaySpaces()))
		query, subquery, release, episode := library.ParseReleaseEpisode(query)
		log.Printf("parsed title '%s'\n", query)
		log.Printf("parsed episode '%s'\n", subquery)
		log.Printf("parsed collation '%s'\n", episode)
		log.Printf("parsed released '%s'\n", release)
		query = library.StripHallucinations(rec.Query, query)
		subquery = library.StripHallucinations(rec.Query, subquery)
		log.Printf("stripped '%s' | '%s'\n", query, subquery)
		query = lucenex.Clean(query)
		log.Println("query cleaned", query, release, episode)
		result := ScoredKnown{Relevance: t.MinRelevance}

		collation := library.KnownStringCollationEpisode(episode)

		{
			q := library.KnownSearchBuilder().Where(squirrel.And{
				library.KnownQueryExplicit(t.Explicit),
				lucenex.Query(duckdbx.NewLucene(), query, lucenex.WithDefaultField("auto_description")),
			}).
				OrderByClause(library.KnownOrderCollationNearest(collation)).
				OrderByClause(library.KnownOrderReleasedNearest(library.KnownStringRelease(release))).
				OrderByClause(library.KnownOrderTitleSimilarity(query, 0.7)).
				OrderByClause(library.KnownOrderSubtitleSimilarity(subquery, 0.7)).
				Limit(1028)

			scanner := sqlx.Scan(library.KnownSearch(ctx, db, q))

			for v := range scanner.Iter() {
				var cur = ScoredKnown{Known: v}

				if err := library.KnownScoreByID(ctx, db, v.UID, query, t.Cutoff).Scan(&cur.Relevance); err != nil {
					log.Println("unable to score", v.UID, err)
					continue
				}

				if cur.Relevance > result.Relevance {
					result = cur
					log.Println(cur.Relevance, cur.UID, cur.Title, library.KnownCollationString(cur.Collation), cur.Subtitle, cur.Released)
				}
			}

			if err := scanner.Err(); err != nil {
				return err
			}
		}

		if stringsx.Present(result.Title) {
			log.Println("result", result.Relevance, result.UID, result.Title, library.KnownCollationString(result.Collation), result.Subtitle, result.Released, result.Mimetype)
			continue
		}

		{
			terms := strings.ReplaceAll(stringsx.CompactWhitespace(query), " ", " OR ")
			log.Println("DERP DERP", terms)
			q := library.KnownSearchBuilder().Where(squirrel.And{
				library.KnownQueryExplicit(t.Explicit),
				lucenex.Query(duckdbx.NewLucene(), terms, lucenex.WithDefaultField("title")),
			}).
				OrderByClause(library.KnownOrderCollationNearest(collation)).
				OrderByClause(library.KnownOrderReleasedNearest(library.KnownStringRelease(release))).
				OrderByClause(library.KnownOrderTitleSimilarity(query, 0.7)).
				OrderByClause(library.KnownOrderSubtitleSimilarity(subquery, 0.7)).
				Limit(1028)

			scanner := sqlx.Scan(library.KnownSearch(ctx, sqlx.Debug(db), q))

			for v := range scanner.Iter() {
				var cur = ScoredKnown{Known: v}

				if err := library.KnownScoreByID(ctx, db, v.UID, query, t.Cutoff).Scan(&cur.Relevance); err != nil {
					log.Println("unable to score", v.UID, err)
					continue
				}

				if cur.Relevance > result.Relevance {
					result = cur
					log.Println(cur.Relevance, cur.UID, cur.Title, library.KnownCollationString(cur.Collation), cur.Subtitle, cur.Released)
				}
			}

			if err := scanner.Err(); err != nil {
				return err
			}
		}

		if stringsx.Present(result.Title) {
			log.Println("result", result.Relevance, result.UID, result.Title, library.KnownCollationString(result.Collation), result.Subtitle, result.Released, result.Mimetype)
		} else {
			log.Println("no result found")
		}
	}

	if err := seq.Err(); err != nil {
		return errorsx.Wrap(err, "failed to decode input")
	}

	if count == 0 {
		return errorsx.New("query is required: stdin was empty")
	}

	return nil
}
