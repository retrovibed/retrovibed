package ddisc

import (
	"context"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/searchplugin"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
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
		DiscoveredQueryNotTombstoned(),
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
		RecommendationFromDiscovered(
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

// RecommendationFromDiscovered maps a located Discovered row into a
// library.Recommendation keyed on d's own ddisc_media row. Pure mapping -
// does not touch the database; callers are responsible for persisting it
// via library.RecommendationInsertWithDefaults.
func RecommendationFromDiscovered(d Discovered, options ...library.RecommendationOption) library.Recommendation {
	return langx.Clone(library.Recommendation{
		Source:       md5x.String(library.RecommendationSourceDiscovered),
		ContentID:    d.ID,
		KnownMediaID: d.KnownMediaID,
		Mimetype:     mimex.Category(d.Mimetype),
		Language:     d.AudioDefaultLocale,
		Adult:        d.Adult,
		TombstoneAt:  time.Now().Add(library.RecommendationTTL),
		Title:        d.Title,
		Overview:     d.Description,
		Image:        d.PosterURI,
		Released:     d.ReleasedAt,
	}, options...)
}

// RecommendationsFromPlugins asks every loaded search plugin for up to limit
// recommended candidates for mimetype (via searchplugin.R.Recommend, which is
// query-less by design), restricted to lang/adult the same way
// ddisc.Recommendation restricts its own random pick, and persists each as a
// library.Recommendation, deduplicated by content id through
// RecommendationInsertWithDefaults's existing ON CONFLICT. A single plugin
// failing is already logged and skipped inside Registry.Recommend, and no
// registry being wired at all (searchplugin.Unimplemented's
// errors.ErrUnsupported) is treated the same as zero plugins responding -
// that's a searchplugin-layer detail callers of this pipeline shouldn't have
// to know about.
func RecommendationsFromPlugins(ctx context.Context, q sqlx.Queryer, plugins searchplugin.R, mimetype string, limit uint, lang string, adult bool) error {
	seq := plugins.Recommend(ctx, []string{mimetype}, limit, lang, adult, false)

	for imp := range seq.Each(ctx) {
		d := NewDiscoveredFromImport(imp, DiscoveredOptionMimetype(imp.Mimetype))

		var rec library.Recommendation
		if err := library.RecommendationInsertWithDefaults(ctx, q, RecommendationFromDiscovered(d, library.RecommendationOptionSourceSearchPlugin)).Scan(&rec); err != nil {
			return err
		}
	}

	return errorsx.Ignore(seq.Err(), errors.ErrUnsupported)
}
