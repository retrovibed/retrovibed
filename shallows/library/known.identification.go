package library

import (
	"context"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/retrovibed/retrovibed/shallows/internal/duckdbx"
	"github.com/retrovibed/retrovibed/shallows/internal/lucenex"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
)

type KnownScored struct {
	Known
	Relevance float64
}

type KnownIdentifier struct {
	q         sqlx.Queryer
	Cutoff    float32
	Threshold float32
	Limit     uint
	Explicit  bool
	cleaner   QueryCleaner
}

func (t KnownIdentifier) Identify(ctx context.Context, i string) (res KnownScored, err error) {
	var (
		cleaned, query, subquery, release, episode string
	)

	if cleaned, err = t.cleaner.Clean(ctx, i); err != nil {
		log.Println("unable to clean query", err)
		query = i
	} else {
		query, subquery, release, episode = ParseReleaseEpisode(cleaned)
	}

	collation := KnownStringCollationEpisode(episode)
	squery := StripHallucinations(i, query)
	ssubquery := StripHallucinations(i, subquery)

	q := KnownSearchBuilder().Where(squirrel.And{
		KnownQueryExplicit(t.Explicit),
		lucenex.Query(duckdbx.NewLucene(), lucenex.Clean(query), lucenex.WithDefaultField("auto_description")),
	}).
		OrderByClause(KnownOrderCollationNearest(collation)).
		OrderByClause(KnownOrderReleasedNearest(KnownStringRelease(release))).
		OrderByClause(KnownOrderTitleSimilarity(squery, t.Threshold)).
		OrderByClause(KnownOrderSubtitleSimilarity(ssubquery, t.Threshold)).
		Limit(1028)

	scanner := sqlx.Scan(KnownSearch(ctx, t.q, q))

	for v := range scanner.Iter() {
		cur := KnownScored{Known: v}
		if err := KnownScoreByID(ctx, t.q, v.UID, squery, t.Cutoff).Scan(&cur.Relevance); err != nil {
			log.Println("unable to score", v.UID, err)
			continue
		}

		if cur.Relevance > res.Relevance {
			res = cur
		}
	}

	if err := scanner.Err(); err != nil {
		return res, err
	}

	return res, nil
}
