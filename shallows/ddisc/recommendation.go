package ddisc

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/localex"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
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
		Source:       md5x.String(library.RecommendationSourceRandom),
		ContentID:    random.KnownMediaID,
		KnownMediaID: random.KnownMediaID,
		Mimetype:     mimex.Category(random.Mimetype),
		Language:     localex.FirstDefined(random.AudioDefaultLocale, known.OriginalLanguage),
		Adult:        langx.FirstNonZero(random.Adult, known.Adult),
		TombstoneAt:  time.Now().Add(library.RecommendationTTL),
		Title:        known.Title,
		Overview:     known.Overview,
		Image:        stringsx.FirstNonBlank(known.PosterPath, known.BackdropPath),
		Popularity:   known.Popularity,
		Released:     known.Released,
	}).Scan(&rec); err != nil {
		return rec, err
	}

	return rec, nil
}

// RecommendationFromDiscovered maps a located Discovered row into a
// library.Recommendation keyed on d's own ddisc_media row. Pure mapping -
// does not touch the database; callers are responsible for persisting it
// via library.RecommendationInsertWithDefaults.
func RecommendationFromDiscovered(d Discovered) library.Recommendation {
	return library.Recommendation{
		Source:       md5x.String(library.RecommendationSourceDiscovered),
		ContentID:    d.ID,
		KnownMediaID: d.KnownMediaID,
		Mimetype:     mimex.Category(d.Mimetype),
		Language:     d.AudioDefaultLocale,
		Adult:        d.Adult,
		TombstoneAt:  time.Now().Add(library.RecommendationTTL),
		Title:        d.Title,
		Overview:     d.Description,
		Released:     d.ReleasedAt,
	}
}
