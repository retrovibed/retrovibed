package media

import (
	"context"
	"errors"

	"github.com/retrovibed/retrovibed/retroapi/searchplugin"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
)

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
		d := ddisc.NewDiscoveredFromImport(
			imp,
			ddisc.DiscoveredOptionMimetype(imp.Mimetype, mimetype),
		)

		if err := ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d); err != nil {
			return err
		}

		var rec library.Recommendation
		if err := library.RecommendationInsertWithDefaults(ctx, q, ddisc.RecommendationFromDiscovered(d, library.RecommendationOptionSourceSearchPlugin)).Scan(&rec); err != nil {
			return err
		}
	}

	return errorsx.Ignore(seq.Err(), errors.ErrUnsupported)
}
