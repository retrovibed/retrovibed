package ddisc_test

import (
	"database/sql"
	"testing"

	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestRecommendation(t *testing.T) {
	t.Run("creates recommendation from random discovered media", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		id := int160.Random()
		d := ddisc.NewDiscovered(&id, "", ddisc.DiscoveredOptionKnownMedia(known.UID))
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))

		rec, err := ddisc.Recommendation(ctx, q, mimex.Application, "", false)
		require.NoError(t, err)
		require.Equal(t, known.UID, rec.ContentID)
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
		d := ddisc.NewDiscovered(&id, "", ddisc.DiscoveredOptionKnownMedia(known.UID))
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
		d := ddisc.NewDiscovered(&id, "", ddisc.DiscoveredOptionKnownMedia(known.UID))
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
		d := ddisc.NewDiscovered(&id, "", ddisc.DiscoveredOptionKnownMedia(known.UID))
		d.Adult = true
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))

		rec, err := ddisc.Recommendation(ctx, q, mimex.Application, "", true)
		require.NoError(t, err)
		require.Equal(t, known.UID, rec.ContentID)
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
		dv := ddisc.NewDiscovered(&vid, "", ddisc.DiscoveredOptionKnownMedia(video.UID), ddisc.DiscoveredOptionMimetype("video/mp4"))
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, dv).Scan(&dv))

		aid := int160.Random()
		da := ddisc.NewDiscovered(&aid, "", ddisc.DiscoveredOptionKnownMedia(audio.UID), ddisc.DiscoveredOptionMimetype("audio/mpeg"))
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, da).Scan(&da))

		rec, err := ddisc.Recommendation(ctx, q, mimex.Video, "", false)
		require.NoError(t, err)
		require.Equal(t, video.UID, rec.ContentID)
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
		d := ddisc.NewDiscovered(&id, "", ddisc.DiscoveredOptionKnownMedia(known.UID))
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
		d := ddisc.NewDiscovered(&id, "", ddisc.DiscoveredOptionKnownMedia(known.UID))
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
		d := ddisc.NewDiscovered(&id, "", ddisc.DiscoveredOptionKnownMedia(known.UID))
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))

		rec, err := ddisc.Recommendation(ctx, q, mimex.Application, "", false)
		require.NoError(t, err)
		require.Equal(t, "fr", rec.Language)
	})
}

func TestRecommendationFromDiscovered(t *testing.T) {
	id := int160.Random()
	d := ddisc.NewDiscovered(&id, "", ddisc.DiscoveredOptionMimetype("video/mp4"))
	d.AudioDefaultLocale = "en"
	d.Adult = true

	rec := ddisc.RecommendationFromDiscovered(d)
	require.Equal(t, d.ID, rec.ContentID)
	require.Equal(t, md5x.String(library.RecommendationSourceDiscovered), rec.Source)
	require.Equal(t, mimex.Category(d.Mimetype), rec.Mimetype)
	require.Equal(t, "en", rec.Language)
	require.True(t, rec.Adult)
}
