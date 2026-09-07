package media_test

import (
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/retroapi/userx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/pqueuetestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/stretchr/testify/require"
)

func TestRecommendationsBackgroundRun(t *testing.T) {
	t.Run("generates a video recommendation when none exist", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		w := pqueuetestx.NewQueue()

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionMimetype(mimex.Video)))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		id := int160.Random()
		d := ddisc.NewDiscovered(&id, ddisc.DiscoveredOptionKnownMedia(known.UID), ddisc.DiscoveredOptionMimetype("video/mp4"), ddisc.DiscoveredOptionTestDefaults, ddisc.DiscoveredOptionAutoMagnet)
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))

		require.NoError(t, media.RecommendationsBackgroundRun(ctx, q, w))

		require.Equal(t, 2, w.Len())
		var videoProfileID string
		var videoLanguage string
		for _, e := range w.Snapshot() {
			var req media.RecommendationRefreshRequest
			require.NoError(t, jsonx.Unmarshal(e, &req))
			if req.Mimetype == mimex.Video {
				videoProfileID, videoLanguage = req.ProfileId, req.Language
			}
		}
		require.Equal(t, uuid.Nil.String(), videoProfileID)
		require.Equal(t, userx.LocaleLanguage(), videoLanguage)
	})

	t.Run("generates a audio recommendation when none exist", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		w := pqueuetestx.NewQueue()

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionMimetype(mimex.Audio)))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		id := int160.Random()
		d := ddisc.NewDiscovered(&id, ddisc.DiscoveredOptionKnownMedia(known.UID), ddisc.DiscoveredOptionMimetype("audio/mpeg"), ddisc.DiscoveredOptionTestDefaults, ddisc.DiscoveredOptionAutoMagnet)
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))

		require.NoError(t, media.RecommendationsBackgroundRun(ctx, q, w))

		require.Equal(t, 2, w.Len())
		var audioProfileID string
		var audioLanguage string
		for _, e := range w.Snapshot() {
			var req media.RecommendationRefreshRequest
			require.NoError(t, jsonx.Unmarshal(e, &req))
			if req.Mimetype == mimex.Audio {
				audioProfileID, audioLanguage = req.ProfileId, req.Language
			}
		}
		require.Equal(t, uuid.Nil.String(), audioProfileID)
		require.Equal(t, userx.LocaleLanguage(), audioLanguage)
	})

	t.Run("skips generation when last recommendation is within 24 hours", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		w := pqueuetestx.NewQueue()

		// RecommendationsBackgroundRun only enqueues now, it never writes to
		// library_recommendations itself, so the cooldown it checks via
		// RecommendationLastGeneratedAt has to be seeded directly here.
		var rec library.Recommendation
		require.NoError(t, testx.Fake(&rec, library.RecommendationOptionTestDefaults))
		require.NoError(t, library.RecommendationInsertWithDefaults(ctx, q, rec).Scan(&rec))

		require.NoError(t, media.RecommendationsBackgroundRun(ctx, q, w))
		require.Zero(t, w.Len())
	})

	t.Run("no known media returns no error", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		w := pqueuetestx.NewQueue()

		require.NoError(t, media.RecommendationsBackgroundRun(ctx, q, w))
		require.Equal(t, 2, w.Len())
	})

	t.Run("queues both audio and video requests in a single run", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		w := pqueuetestx.NewQueue()

		require.NoError(t, media.RecommendationsBackgroundRun(ctx, q, w))

		require.Equal(t, 2, w.Len())
		mimetypes := map[string]bool{}
		for _, e := range w.Snapshot() {
			var req media.RecommendationRefreshRequest
			require.NoError(t, jsonx.Unmarshal(e, &req))
			mimetypes[req.Mimetype] = true
		}
		require.Equal(t, map[string]bool{mimex.Audio: true, mimex.Video: true}, mimetypes)
	})
}
