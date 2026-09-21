package library

import (
	"context"
	"log"
	"runtime/trace"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/retrovibed/retrovibed/shallows/internal/duckdbx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/lucenex"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
)

type KnownScored struct {
	Known
	Relevance float64
}

// KnownIdentifierOptionMimetype restricts identification to the given mimetype category.
func KnownIdentifierOptionMimetype(v string) func(*KnownIdentifier) {
	return func(t *KnownIdentifier) {
		t.Mimetype = v
	}
}

func NewKnownIdentifier(q sqlx.Queryer, c QueryCleaner, options ...func(*KnownIdentifier)) *KnownIdentifier {
	return new(langx.Clone(KnownIdentifier{
		q:            q,
		cleaner:      c,
		Cutoff:       0.7,
		Threshold:    0.7,
		MinRelevance: 0.85,
		Limit:        8,
	}, options...))
}

type KnownIdentifier struct {
	q            sqlx.Queryer
	cleaner      QueryCleaner
	Cutoff       float32
	Threshold    float32
	MinRelevance float64
	Limit        uint
	Explicit     bool
	// Mimetype restricts candidates to a mimetype category, blank does not restrict.
	Mimetype string
}

func (t KnownIdentifier) Identify(ctx context.Context, i string) (res KnownScored, err error) {
	var (
		cleaned, query, subquery, release, episode string
		task                                       *trace.Task
	)

	ctx, task = trace.NewTask(ctx, "known.identify")
	defer task.End()

	trace.Logf(ctx, "input", "%q", i)

	// a miss is the unknown media.
	res = KnownScored{Known: Unknown()}

	if cleaned, err = t.cleaner.Clean(ctx, i); err != nil {
		log.Println("unable to clean query", err)
		query = i
	} else {
		query, subquery, release, episode = ParseReleaseEpisode(cleaned)
	}

	collation := KnownStringCollationEpisode(episode)
	squery := StripHallucinations(i, query)
	ssubquery := StripHallucinations(i, subquery)

	terms := strings.ReplaceAll(stringsx.CompactWhitespace(lucenex.Clean(query)), " ", " OR ")

	trace.Logf(ctx, "cleaned", "%q", cleaned)
	trace.Logf(ctx, "parsed", "title: %q subtitle: %q episode: %q release: %q", query, subquery, episode, release)
	trace.Logf(ctx, "lucene", "%q", terms)
	trace.Logf(ctx, "stripped", "%q | %q", squery, ssubquery)

	// nothing to search with (i.e. the cleaner detected a messy input), an empty lucene query would
	// match the entire catalog.
	if stringsx.Blank(terms) {
		return res, nil
	}

	q := KnownSearchBuilder().Where(squirrel.And{
		KnownQueryExplicit(t.Explicit),
		KnownQueryMimetype(t.Mimetype),
		lucenex.Query(duckdbx.NewLucene(), terms, lucenex.WithDefaultField("auto_description")),
	}).
		OrderByClause(KnownOrderCollationNearest(collation)).
		OrderByClause(KnownOrderReleasedNearest(KnownStringRelease(release))).
		OrderByClause(KnownOrderTitleSimilarity(squery, t.Threshold)).
		OrderByClause(KnownOrderSubtitleSimilarity(ssubquery, t.Threshold)).
		Limit(uint64(t.Limit))
	scanner := sqlx.Scan(KnownSearch(ctx, t.q, q))

	for v := range scanner.Iter() {
		cur := KnownScored{Known: v}
		if err := KnownScoreByID(ctx, t.q, v.UID, squery, t.Cutoff).Scan(&cur.Relevance); err != nil {
			log.Println("unable to score", v.UID, err)
			continue
		}

		// only candidates exceeding the minimum relevance can be identified.
		if cur.Relevance > max(res.Relevance, t.MinRelevance) {
			res = cur
		}
	}

	if err := scanner.Err(); err != nil {
		return res, err
	}

	trace.Logf(ctx, "result", "title: %q relevance: %v", res.Title, res.Relevance)

	return res, nil
}
