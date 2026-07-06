package ddisc

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/localex"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
)

func Recommendation(ctx context.Context, q sqlx.Queryer, mimetype string, lang string, adult bool) (rec library.Recommendation, err error) {
	var (
		random Discovered
		known  library.Known
	)

	scanner := DiscoveredSearch(ctx, q, DiscoveredSearchBuilder().Where(squirrel.And{
		DiscoveredQueryKnown(),
		DiscoveredQueryCategory(mimetype),
		DiscoveredQueryExplicit(adult),
		DiscoveredQueryLanguage(lang),
	}).OrderBy("random()").Limit(1))

	if random, err = sqlx.ScanOne(scanner); err != nil {
		return rec, err
	}

	if err = library.KnownFindByID(ctx, q, random.KnownMediaID).Scan(&known); err != nil {
		return rec, err
	}

	if err = library.RecommendationInsertWithDefaults(ctx, q, library.Recommendation{
		ID:           uuid.Must(uuid.NewV7()).String(),
		Source:       md5x.String(library.RecommendationSourceRandom),
		KnownMediaID: random.KnownMediaID,
		Mimetype:     mimex.Category(random.Mimetype),
		Language:     localex.FirstDefined(random.AudioDefaultLocale, known.OriginalLanguage),
		Adult:        langx.FirstNonZero(random.Adult, known.Adult),
		TombstoneAt:  time.Now().Add(library.RecommendationTTL),
	}).Scan(&rec); err != nil {
		return rec, err
	}

	return rec, nil
}
