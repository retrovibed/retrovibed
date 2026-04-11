package cmdmedia

import (
	"database/sql"
	"log"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/retrovibed/retrovibed/cmd/cmdmeta"
	"github.com/retrovibed/retrovibed/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/internal/duckdbx"
	"github.com/retrovibed/retrovibed/internal/lucenex"
	"github.com/retrovibed/retrovibed/internal/sqlx"
	"github.com/retrovibed/retrovibed/internal/stringsx"
	"github.com/retrovibed/retrovibed/library"
)

type knownquery struct {
	Query string `arg:"" name:"query" help:"query to search for" required:"true"`
}

func (t knownquery) Run(gctx *cmdopts.Global) (err error) {
	type ScoredKnown struct {
		library.Known
		Relevance float64
	}
	var (
		db     *sql.DB
		result = ScoredKnown{
			Relevance: -1.0,
		}
	)

	if db, err = cmdmeta.Database(gctx.Context); err != nil {
		return err
	}
	defer db.Close()

	{

		q := library.KnownSearchBuilder().Where(squirrel.And{
			library.KnownQueryExplicit(false),
			lucenex.Query(duckdbx.NewLucene(), t.Query, lucenex.WithDefaultField("auto_description")),
		}).OrderBy("title DESC").Limit(1028)

		scanner := sqlx.Scan(library.KnownSearch(gctx.Context, db, q))

		for v := range scanner.Iter() {
			var cur = ScoredKnown{Known: v}

			if err := library.KnownScoreByID(gctx.Context, db, v.UID, t.Query).Scan(&cur.Relevance); err != nil {
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
			lucenex.Query(duckdbx.NewLucene(), terms, lucenex.WithDefaultField("title")),
		})

		scanner := sqlx.Scan(library.KnownSearch(gctx.Context, db, q))

		for v := range scanner.Iter() {
			var cur = ScoredKnown{Known: v}

			if err := library.KnownScoreByID(gctx.Context, db, v.UID, t.Query).Scan(&cur.Relevance); err != nil {
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
