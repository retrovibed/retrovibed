package cmdmedia

import (
	"context"
	"database/sql"
	"log"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdmeta"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/duckdbx"
	"github.com/retrovibed/retrovibed/shallows/internal/lucenex"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/library"
)

type knownquery struct {
	Database string `flag:"" name:"database" help:"database to read" default:"${vars_user_configuration_directory}/meta.db"`
	Explicit bool   `flag:"" name:"explicit" help:"include explicit content in results" default:"false"`
	Query    string `arg:"" name:"query" help:"query to search for" required:"true"`
}

func (t knownquery) Run(gctx *cmdopts.Global) (err error) {
	var db *sql.DB
	if db, err = cmdmeta.DatabaseCustom(gctx.Context, t.Database); err != nil {
		return err
	}
	defer db.Close()
	return t.run(gctx.Context, db)
}

func (t knownquery) run(ctx context.Context, db *sql.DB) (err error) {
	type ScoredKnown struct {
		library.Known
		Relevance float64
	}

	result := ScoredKnown{Relevance: -1.0}

	{
		q := library.KnownSearchBuilder().Where(squirrel.And{
			library.KnownQueryExplicit(t.Explicit),
			lucenex.Query(duckdbx.NewLucene(), t.Query, lucenex.WithDefaultField("auto_description")),
		}).OrderBy("title DESC").Limit(1028)

		scanner := sqlx.Scan(library.KnownSearch(ctx, db, q))

		for v := range scanner.Iter() {
			var cur = ScoredKnown{Known: v}

			if err := library.KnownScoreByID(ctx, db, v.UID, t.Query).Scan(&cur.Relevance); err != nil {
				log.Println("unable to score", v.UID, err)
				continue
			}

			if cur.Relevance > result.Relevance {
				result = cur
				log.Println(cur.Relevance, cur.UID, cur.Title)
			}
		}

		if err := scanner.Err(); err != nil {
			return err
		}
	}

	if result.Relevance > 0 {
		log.Println("result", result.Relevance, result.UID, result.Title)
		return nil
	}

	{
		terms := strings.ReplaceAll(stringsx.CompactWhitespace(t.Query), " ", " OR ")
		q := library.KnownSearchBuilder().Where(squirrel.And{
			library.KnownQueryExplicit(t.Explicit),
			lucenex.Query(duckdbx.NewLucene(), terms, lucenex.WithDefaultField("title")),
		})

		scanner := sqlx.Scan(library.KnownSearch(ctx, db, q))

		for v := range scanner.Iter() {
			var cur = ScoredKnown{Known: v}

			if err := library.KnownScoreByID(ctx, db, v.UID, t.Query).Scan(&cur.Relevance); err != nil {
				log.Println("unable to score", v.UID, err)
				continue
			}

			if cur.Relevance > result.Relevance {
				result = cur
				log.Println(cur.Relevance, cur.UID, cur.Title)
			}
		}

		if err := scanner.Err(); err != nil {
			return err
		}
	}

	log.Println("result", result.Relevance, result.UID, result.Title)
	return nil
}
