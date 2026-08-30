package ddisc_test

import (
	"context"
	"database/sql"
	"iter"
	"testing"

	"github.com/james-lawrence/torrent/dht/int160"
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

func TestRecommendation(t *testing.T) {
	t.Run("creates recommendation from random discovered media", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		id := int160.Random()
		d := ddisc.NewDiscovered(&id, ddisc.DiscoveredOptionKnownMedia(known.UID), ddisc.DiscoveredOptionAutoMagnet)
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))

		rec, err := ddisc.Recommendation(ctx, q, mimex.Application, "", false)
		require.NoError(t, err)
		require.Equal(t, d.ID, rec.ContentID)
		require.NotEmpty(t, rec.ID)
	})

	t.Run("sets source to random", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		id := int160.Random()
		d := ddisc.NewDiscovered(&id, ddisc.DiscoveredOptionKnownMedia(known.UID), ddisc.DiscoveredOptionAutoMagnet)
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))

		rec, err := ddisc.Recommendation(ctx, q, mimex.Application, "", false)
		require.NoError(t, err)
		require.NotEmpty(t, rec.Source)
	})

	t.Run("empty database returns error", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		_, err := ddisc.Recommendation(ctx, q, mimex.Application, "", false)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("excludes adult content", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		id := int160.Random()
		d := ddisc.NewDiscovered(&id, ddisc.DiscoveredOptionKnownMedia(known.UID), ddisc.DiscoveredOptionAutoMagnet)
		d.Adult = true
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))

		_, err := ddisc.Recommendation(ctx, q, mimex.Application, "", false)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("returns adult content when requested", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		id := int160.Random()
		d := ddisc.NewDiscovered(&id, ddisc.DiscoveredOptionKnownMedia(known.UID), ddisc.DiscoveredOptionAutoMagnet)
		d.Adult = true
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))

		rec, err := ddisc.Recommendation(ctx, q, mimex.Application, "", true)
		require.NoError(t, err)
		require.Equal(t, d.ID, rec.ContentID)
		require.True(t, rec.Adult)
	})

	t.Run("filters by mimetype", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var video, audio library.Known
		require.NoError(t, testx.Fake(&video, library.KnownOptionTestDefaults))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, video).Scan(&video))
		require.NoError(t, testx.Fake(&audio, library.KnownOptionTestDefaults))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, audio).Scan(&audio))

		vid := int160.Random()
		dv := ddisc.NewDiscovered(&vid, ddisc.DiscoveredOptionKnownMedia(video.UID), ddisc.DiscoveredOptionMimetype("video/mp4"), ddisc.DiscoveredOptionAutoMagnet)
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, dv).Scan(&dv))

		aid := int160.Random()
		da := ddisc.NewDiscovered(&aid, ddisc.DiscoveredOptionKnownMedia(audio.UID), ddisc.DiscoveredOptionMimetype("audio/mpeg"), ddisc.DiscoveredOptionAutoMagnet)
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, da).Scan(&da))

		rec, err := ddisc.Recommendation(ctx, q, mimex.Video, "", false)
		require.NoError(t, err)
		require.Equal(t, dv.ID, rec.ContentID)
		require.Equal(t, mimex.Video, rec.Mimetype)
	})

	t.Run("increments recommendations counter on repeat", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		id := int160.Random()
		d := ddisc.NewDiscovered(&id, ddisc.DiscoveredOptionKnownMedia(known.UID), ddisc.DiscoveredOptionAutoMagnet)
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))

		_, err := ddisc.Recommendation(ctx, q, mimex.Application, "", false)
		require.NoError(t, err)

		rec, err := ddisc.Recommendation(ctx, q, mimex.Application, "", false)
		require.NoError(t, err)
		require.EqualValues(t, 1, rec.Recommendations)
	})

	t.Run("prefers discovered locale over known original language", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		known.OriginalLanguage = "fr"
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		id := int160.Random()
		d := ddisc.NewDiscovered(&id, ddisc.DiscoveredOptionKnownMedia(known.UID), ddisc.DiscoveredOptionAutoMagnet)
		d.AudioDefaultLocale = "en"
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))

		rec, err := ddisc.Recommendation(ctx, q, mimex.Application, "", false)
		require.NoError(t, err)
		require.Equal(t, "en", rec.Language)
	})

	t.Run("falls back to known original language when discovered locale is undetermined", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		known.OriginalLanguage = "fr"
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		id := int160.Random()
		d := ddisc.NewDiscovered(&id, ddisc.DiscoveredOptionKnownMedia(known.UID), ddisc.DiscoveredOptionAutoMagnet)
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))

		rec, err := ddisc.Recommendation(ctx, q, mimex.Application, "", false)
		require.NoError(t, err)
		require.Equal(t, "fr", rec.Language)
	})
}

func TestRecommendationFromDiscovered(t *testing.T) {
	id := int160.Random()
	d := ddisc.NewDiscovered(&id, ddisc.DiscoveredOptionMimetype("video/mp4"), ddisc.DiscoveredOptionAutoMagnet)
	d.AudioDefaultLocale = "en"
	d.Adult = true

	rec := ddisc.RecommendationFromDiscovered(d)
	require.Equal(t, d.ID, rec.ContentID)
	require.Equal(t, md5x.String(library.RecommendationSourceDiscovered), rec.Source)
	require.Equal(t, mimex.Category(d.Mimetype), rec.Mimetype)
	require.Equal(t, "en", rec.Language)
	require.True(t, rec.Adult)
}

func TestRecommendationsFromPlugins(t *testing.T) {
	t.Run("persists a recommendation per yielded import", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		plugins := fakeRecommendPlugins{results: []*ddiscapi.Import{
			{Uri: "magnet:?xt=urn:btih:4444444444444444444444444444444444444444&dn=one", Mimetype: mimex.Video, Title: "one", PosterPath: "http://example.com/one.jpg"},
			{Uri: "magnet:?xt=urn:btih:5555555555555555555555555555555555555555&dn=two", Mimetype: mimex.Video, Title: "two"},
		}}

		require.NoError(t, ddisc.RecommendationsFromPlugins(ctx, q, plugins, mimex.Video, 5, "", false))

		count, err := sqlx.Count(ctx, q, "SELECT COUNT(*) FROM library_recommendations")
		require.NoError(t, err)
		require.Equal(t, 2, count)
	})

	t.Run("tags source and mimetype from the import, not the discovered default", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		plugins := fakeRecommendPlugins{results: []*ddiscapi.Import{
			{Uri: "magnet:?xt=urn:btih:6666666666666666666666666666666666666666&dn=one", Mimetype: mimex.Video, Title: "one", PosterPath: "http://example.com/one.jpg"},
		}}

		require.NoError(t, ddisc.RecommendationsFromPlugins(ctx, q, plugins, mimex.Video, 5, "", false))

		rec, err := sqlx.ScanOne(library.RecommendationSearch(ctx, q, library.RecommendationSearchBuilder()))
		require.NoError(t, err)
		require.Equal(t, md5x.String(library.RecommendationSourceSearchPlugin), rec.Source)
		require.Equal(t, mimex.Video, rec.Mimetype)
		require.Equal(t, "http://example.com/one.jpg", rec.Image)
	})

	t.Run("dedups repeat recommendations of the same content", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		plugins := fakeRecommendPlugins{results: []*ddiscapi.Import{
			{Uri: "magnet:?xt=urn:btih:7777777777777777777777777777777777777777&dn=one", Mimetype: mimex.Video, Title: "one"},
		}}

		require.NoError(t, ddisc.RecommendationsFromPlugins(ctx, q, plugins, mimex.Video, 5, "", false))
		require.NoError(t, ddisc.RecommendationsFromPlugins(ctx, q, plugins, mimex.Video, 5, "", false))

		rec, err := sqlx.ScanOne(library.RecommendationSearch(ctx, q, library.RecommendationSearchBuilder()))
		require.NoError(t, err)
		require.EqualValues(t, 1, rec.Recommendations)
	})

	t.Run("no plugins configured is a silent no-op", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, ddisc.RecommendationsFromPlugins(ctx, q, searchplugin.Unimplemented{}, mimex.Video, 5, "", false))

		count, err := sqlx.Count(ctx, q, "SELECT COUNT(*) FROM library_recommendations")
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})
}
