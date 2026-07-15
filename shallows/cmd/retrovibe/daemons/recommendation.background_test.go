package daemons_test

import (
	"testing"

	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/retrovibe/daemons"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestRecommendationsBackgroundRun(t *testing.T) {
	t.Run("generates a video recommendation when none exist", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionMimetype(mimex.Video)))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		id := int160.Random()
		d := ddisc.NewDiscovered(&id, ddisc.DiscoveredOptionKnownMedia(known.UID), ddisc.DiscoveredOptionMimetype("video/mp4"), ddisc.DiscoveredOptionTestDefaults)
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))

		require.NoError(t, daemons.RecommendationsBackgroundRun(ctx, q))

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
		d := ddisc.NewDiscovered(&id, ddisc.DiscoveredOptionKnownMedia(known.UID), ddisc.DiscoveredOptionMimetype("audio/mpeg"), ddisc.DiscoveredOptionTestDefaults)
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))

		require.NoError(t, daemons.RecommendationsBackgroundRun(ctx, q))

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
		d := ddisc.NewDiscovered(&id, ddisc.DiscoveredOptionKnownMedia(known.UID), ddisc.DiscoveredOptionMimetype("video/mp4"), ddisc.DiscoveredOptionTestDefaults)
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))

		// generate the first recommendation
		require.NoError(t, daemons.RecommendationsBackgroundRun(ctx, q))

		// second run should not generate another
		require.NoError(t, daemons.RecommendationsBackgroundRun(ctx, q))

		count, err := sqlx.Count(ctx, q, "SELECT COUNT(*) FROM library_recommendations")
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})

	t.Run("no known media returns no error", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, daemons.RecommendationsBackgroundRun(ctx, q))
	})
}
