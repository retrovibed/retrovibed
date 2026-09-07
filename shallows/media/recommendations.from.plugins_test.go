package media_test

import (
	"context"
	"iter"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/ddiscapi"
	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/searchplugin"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/stretchr/testify/require"
)

// fakeRecommendPlugins is a searchplugin.R that yields a fixed set of
// results instead of running a real wasm plugin - mirrors the
// fakeDiscoverStrategy/fakeDiscoverSeq pair in discover_test.go.
type fakeRecommendPlugins struct {
	results []*ddiscapi.Import
	err     error
}

func (t fakeRecommendPlugins) Recommend(ctx context.Context, mimetypes []string, limit uint, lang string, adult, public bool) iterx.Seq[*ddiscapi.Import] {
	return fakeRecommendSeq(t)
}

type fakeRecommendSeq struct {
	results []*ddiscapi.Import
	err     error
}

func (t fakeRecommendSeq) Each(ctx context.Context) iter.Seq[*ddiscapi.Import] {
	return iterx.From(t.results...)
}

func (t fakeRecommendSeq) Err() error { return t.err }

func TestRecommendationsFromPlugins(t *testing.T) {
	t.Run("persists a recommendation per yielded import", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		plugins := fakeRecommendPlugins{results: []*ddiscapi.Import{
			{Uri: "magnet:?xt=urn:btih:4444444444444444444444444444444444444444&dn=one", Mimetype: mimex.Video, Title: "one", PosterPath: "http://example.com/one.jpg"},
			{Uri: "magnet:?xt=urn:btih:5555555555555555555555555555555555555555&dn=two", Mimetype: mimex.Video, Title: "two"},
		}}

		require.NoError(t, media.RecommendationsFromPlugins(ctx, q, plugins, mimex.Video, 5, "", false))

		count, err := sqlx.Count(ctx, q, "SELECT COUNT(*) FROM library_recommendations")
		require.NoError(t, err)
		require.Equal(t, 2, count)
	})

	t.Run("sets source md5 to search plugin", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		plugins := fakeRecommendPlugins{results: []*ddiscapi.Import{
			{Uri: "magnet:?xt=urn:btih:9999999999999999999999999999999999999999&dn=one", Mimetype: mimex.Video, Title: "one"},
		}}

		require.NoError(t, media.RecommendationsFromPlugins(ctx, q, plugins, mimex.Video, 5, "", false))

		rec, err := sqlx.ScanOne(library.RecommendationSearch(ctx, q, library.RecommendationSearchBuilder()))
		require.NoError(t, err)
		require.Equal(t, md5x.String(library.RecommendationSourceSearchPlugin), rec.Source)
	})

	t.Run("tags source and mimetype from the import, not the discovered default", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		plugins := fakeRecommendPlugins{results: []*ddiscapi.Import{
			{Uri: "magnet:?xt=urn:btih:6666666666666666666666666666666666666666&dn=one", Mimetype: mimex.Video, Title: "one", PosterPath: "http://example.com/one.jpg"},
		}}

		require.NoError(t, media.RecommendationsFromPlugins(ctx, q, plugins, mimex.Video, 5, "", false))

		rec, err := sqlx.ScanOne(library.RecommendationSearch(ctx, q, library.RecommendationSearchBuilder()))
		require.NoError(t, err)
		require.Equal(t, md5x.String(library.RecommendationSourceSearchPlugin), rec.Source)
		require.Equal(t, mimex.Video, rec.Mimetype)
		require.Equal(t, "http://example.com/one.jpg", rec.Image)
	})

	t.Run("falls back to the requested mimetype when the import omits one", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		plugins := fakeRecommendPlugins{results: []*ddiscapi.Import{
			{Uri: "magnet:?xt=urn:btih:8888888888888888888888888888888888888888&dn=one", Title: "one"},
		}}

		require.NoError(t, media.RecommendationsFromPlugins(ctx, q, plugins, mimex.Video, 5, "", false))

		rec, err := sqlx.ScanOne(library.RecommendationSearch(ctx, q, library.RecommendationSearchBuilder()))
		require.NoError(t, err)
		require.Equal(t, mimex.Video, rec.Mimetype)
	})

	t.Run("dedups repeat recommendations of the same content", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		plugins := fakeRecommendPlugins{results: []*ddiscapi.Import{
			{Uri: "magnet:?xt=urn:btih:7777777777777777777777777777777777777777&dn=one", Mimetype: mimex.Video, Title: "one"},
		}}

		require.NoError(t, media.RecommendationsFromPlugins(ctx, q, plugins, mimex.Video, 5, "", false))
		require.NoError(t, media.RecommendationsFromPlugins(ctx, q, plugins, mimex.Video, 5, "", false))

		rec, err := sqlx.ScanOne(library.RecommendationSearch(ctx, q, library.RecommendationSearchBuilder()))
		require.NoError(t, err)
		require.EqualValues(t, 1, rec.Recommendations)
	})

	t.Run("discovered record is created and reachable from the recommendation content id", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		discoveredCount, err := sqlx.Count(ctx, q, "SELECT COUNT(*) FROM ddisc_media")
		require.NoError(t, err)
		require.Zero(t, discoveredCount)

		recommendationCount, err := sqlx.Count(ctx, q, "SELECT COUNT(*) FROM library_recommendations")
		require.NoError(t, err)
		require.Zero(t, recommendationCount)

		plugins := fakeRecommendPlugins{results: []*ddiscapi.Import{
			{Uri: "magnet:?xt=urn:btih:3333333333333333333333333333333333333333&dn=one", Mimetype: mimex.Video, Title: "one"},
		}}

		require.NoError(t, media.RecommendationsFromPlugins(ctx, q, plugins, mimex.Video, 5, "", false))

		rec, err := sqlx.ScanOne(library.RecommendationSearch(ctx, q, library.RecommendationSearchBuilder()))
		require.NoError(t, err)

		var d ddisc.Discovered
		require.NoError(t, ddisc.DiscoveredFindByID(ctx, q, rec.ContentID).Scan(&d))
		require.Equal(t, rec.ContentID, d.ID)
	})

	t.Run("no plugins configured is a silent no-op", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, media.RecommendationsFromPlugins(ctx, q, searchplugin.Unimplemented{}, mimex.Video, 5, "", false))

		count, err := sqlx.Count(ctx, q, "SELECT COUNT(*) FROM library_recommendations")
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})
}
