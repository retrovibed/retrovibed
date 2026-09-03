package media

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
)

func Recommendation(ctx context.Context, q sqlx.Queryer, mimetype string, lang string, adult bool) (rec library.Recommendation, err error) {
	var (
		random ddisc.Discovered
		known  library.Known
	)

	scanner := ddisc.DiscoveredSearch(ctx, q, ddisc.DiscoveredSearchBuilder().Where(squirrel.And{
		ddisc.DiscoveredQueryKnown(),
		ddisc.DiscoveredQueryCategory(mimetype),
		ddisc.DiscoveredQueryExplicit(adult),
		ddisc.DiscoveredQueryLanguage(lang),
		ddisc.DiscoveredQueryNotTombstoned(),
	}).OrderBy("random()").Limit(1))

	if random, err = sqlx.ScanOne(scanner); err != nil {
		return rec, err
	}

	if err = library.KnownFindByID(ctx, q, random.KnownMediaID).Scan(&known); err != nil {
		return rec, err
	}

	err = library.RecommendationInsertWithDefaults(
		ctx,
		q,
		ddisc.RecommendationFromDiscovered(
			random,
			library.RecommendationOptionAutoKnown(known),
			library.RecommendationOptionSourceRandom,
		),
	).Scan(&rec)

	if err != nil {
		return rec, err
	}

	return rec, nil
}
