package daemons_test

import (
	"context"
	"iter"
	"testing"

	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/ddiscapi"
	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/searchplugin"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/retrovibe/daemons"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

// fakeRecommendPlugins is a searchplugin.R that yields a fixed set of
// results instead of running a real wasm plugin - mirrors the
// fakeDiscoverStrategy/fakeDiscoverSeq pair in shallows/ddisc/discover_test.go.
type fakeRecommendPlugins struct {
	results []*ddiscapi.Import
}

func (t fakeRecommendPlugins) Recommend(ctx context.Context, mimetypes []string, limit uint, lang string, adult, public bool) iterx.Seq[*ddiscapi.Import] {
	return fakeRecommendSeq(t)
}

type fakeRecommendSeq struct {
	results []*ddiscapi.Import
}

func (t fakeRecommendSeq) Each(ctx context.Context) iter.Seq[*ddiscapi.Import] {
	return iterx.From(t.results...)
}

func (t fakeRecommendSeq) Err() error { return nil }

func TestRecommendationsBackgroundRun(t *testing.T) {
	t.Run("generates a video recommendation when none exist", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionMimetype(mimex.Video)))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		id := int160.Random()
		d := ddisc.NewDiscovered(&id, ddisc.DiscoveredOptionKnownMedia(known.UID), ddisc.DiscoveredOptionMimetype("video/mp4"), ddisc.DiscoveredOptionTestDefaults, ddisc.DiscoveredOptionAutoMagnet)
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))

		require.NoError(t, daemons.RecommendationsBackgroundRun(ctx, q, searchplugin.Unimplemented{}))

		count, err := sqlx.Count(ctx, q, "SELECT COUNT(*) FROM library_recommendations")
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})

	t.Run("generates a audio recommendation when none exist", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionMimetype(mimex.Audio)))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		id := int160.Random()
		d := ddisc.NewDiscovered(&id, ddisc.DiscoveredOptionKnownMedia(known.UID), ddisc.DiscoveredOptionMimetype("audio/mpeg"), ddisc.DiscoveredOptionTestDefaults, ddisc.DiscoveredOptionAutoMagnet)
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))

		require.NoError(t, daemons.RecommendationsBackgroundRun(ctx, q, searchplugin.Unimplemented{}))

		count, err := sqlx.Count(ctx, q, "SELECT COUNT(*) FROM library_recommendations")
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})

	t.Run("skips generation when last recommendation is within 24 hours", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionMimetype(mimex.Video)))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		id := int160.Random()
		d := ddisc.NewDiscovered(&id, ddisc.DiscoveredOptionKnownMedia(known.UID), ddisc.DiscoveredOptionMimetype("video/mp4"), ddisc.DiscoveredOptionTestDefaults, ddisc.DiscoveredOptionAutoMagnet)
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))

		// generate the first recommendation
		require.NoError(t, daemons.RecommendationsBackgroundRun(ctx, q, searchplugin.Unimplemented{}))

		// second run should not generate another
		require.NoError(t, daemons.RecommendationsBackgroundRun(ctx, q, searchplugin.Unimplemented{}))

		count, err := sqlx.Count(ctx, q, "SELECT COUNT(*) FROM library_recommendations")
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})

	t.Run("no known media returns no error", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, daemons.RecommendationsBackgroundRun(ctx, q, searchplugin.Unimplemented{}))
	})

	t.Run("persists recommendations from search plugins", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		plugins := fakeRecommendPlugins{results: []*ddiscapi.Import{
			{Uri: "magnet:?xt=urn:btih:3333333333333333333333333333333333333333&dn=plugin-rec", Mimetype: mimex.Video, Title: "plugin recommendation"},
		}}

		require.NoError(t, daemons.RecommendationsBackgroundRun(ctx, q, plugins))

		count, err := sqlx.Count(ctx, q, "SELECT COUNT(*) FROM library_recommendations")
		require.NoError(t, err)
		require.Equal(t, 1, count)

		rec, err := sqlx.ScanOne(library.RecommendationSearch(ctx, q, library.RecommendationSearchBuilder()))
		require.NoError(t, err)
		require.Equal(t, library.RecommendationSourceSearchPlugin, library.RecommendationSourceString(rec.Source))
		require.Equal(t, mimex.Video, rec.Mimetype)
	})
}
