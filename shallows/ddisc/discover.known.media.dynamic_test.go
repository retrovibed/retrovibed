package ddisc_test

import (
	"context"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

// TestDiscoverKnownMediaDynamic exercises DiscoverOptionKnownMediaDynamic
// end to end through Discover itself, rather than any one strategy - every
// strategy's yielded candidates flow through the same central stage (see
// KnownMediaDynamic), so it only needs to be proven once here rather than
// duplicated per-strategy.
func TestDiscoverKnownMediaDynamic(t *testing.T) {
	t.Run("mints a catalog entry and stamps the candidate when unresolved with sufficient content", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)

		id := int160.Random()
		hit := ddisc.NewDiscovered(&id,
			ddisc.DiscoveredOptionMimetype(mimex.Video),
			ddisc.DiscoveredOptionTitle("Ubuntu Documentary"),
			ddisc.DiscoveredOptionDescription("a documentary about ubuntu"),
			ddisc.DiscoveredOptionPosterURI("https://video.example/large.jpg"),
		)
		hit.Source = "retrovibed.discovery.peertube"

		seq := ddisc.Discover(
			context.Background(),
			ddisc.DefaultPolicy(),
			ddisc.DiscoverRequest{Query: "ubuntu"},
			[]ddisc.DiscoverOption{ddisc.DiscoverOptionKnownMediaDynamic(ddisc.KnownMediaDynamic(q))},
			fakeDiscoverStrategy{results: []ddisc.Discovered{hit}},
		)

		var got []ddisc.Discovered
		for d := range seq.Each(context.Background()) {
			got = append(got, d)
		}
		require.NoError(t, seq.Err())
		require.Len(t, got, 1)
		require.NotEqual(t, uuid.Nil.String(), got[0].KnownMediaID, "candidate must be stamped with the newly minted known-media id")

		var known library.Known
		require.NoError(t, library.KnownFindByID(context.Background(), q, got[0].KnownMediaID).Scan(&known))
		require.Equal(t, "Ubuntu Documentary", known.Title)
		require.Equal(t, "a documentary about ubuntu", known.Overview)
		require.Equal(t, "https://video.example/large.jpg", known.PosterPath)
		require.Equal(t, "retrovibed.discovery.peertube", known.Source)
	})

	t.Run("skips a candidate that already has a resolved known-media id", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)
		kid := uuid.Must(uuid.NewV4()).String()

		id := int160.Random()
		hit := ddisc.NewDiscovered(&id,
			ddisc.DiscoveredOptionKnownMedia(kid),
			ddisc.DiscoveredOptionMimetype(mimex.Video),
			ddisc.DiscoveredOptionTitle("Ubuntu Documentary"),
			ddisc.DiscoveredOptionPosterURI("https://video.example/large.jpg"),
		)

		seq := ddisc.Discover(
			context.Background(),
			ddisc.DefaultPolicy(),
			ddisc.DiscoverRequest{KnownMediaID: kid},
			[]ddisc.DiscoverOption{ddisc.DiscoverOptionKnownMediaDynamic(ddisc.KnownMediaDynamic(q))},
			fakeDiscoverStrategy{results: []ddisc.Discovered{hit}},
		)

		var got []ddisc.Discovered
		for d := range seq.Each(context.Background()) {
			got = append(got, d)
		}
		require.NoError(t, seq.Err())
		require.Len(t, got, 1)
		require.Equal(t, kid, got[0].KnownMediaID, "an already-resolved candidate must not be re-stamped")
		require.Equal(t, 0, sqltestx.Count(t, q, "SELECT COUNT(*) FROM library_known_media"), "must not mint a duplicate catalog entry for an already-resolved candidate")
	})

	t.Run("skips a candidate missing a title", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)

		id := int160.Random()
		hit := ddisc.NewDiscovered(&id,
			ddisc.DiscoveredOptionMimetype(mimex.Video),
			ddisc.DiscoveredOptionPosterURI("https://video.example/large.jpg"),
		)

		seq := ddisc.Discover(
			context.Background(),
			ddisc.DefaultPolicy(),
			ddisc.DiscoverRequest{Query: "ubuntu"},
			[]ddisc.DiscoverOption{ddisc.DiscoverOptionKnownMediaDynamic(ddisc.KnownMediaDynamic(q))},
			fakeDiscoverStrategy{results: []ddisc.Discovered{hit}},
		)

		var got []ddisc.Discovered
		for d := range seq.Each(context.Background()) {
			got = append(got, d)
		}
		require.NoError(t, seq.Err())
		require.Len(t, got, 1)
		require.Equal(t, uuid.Nil.String(), got[0].KnownMediaID)
		require.Equal(t, 0, sqltestx.Count(t, q, "SELECT COUNT(*) FROM library_known_media"), "a candidate with no title isn't worth cataloging")
	})

	t.Run("skips a candidate missing a poster uri", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)

		id := int160.Random()
		hit := ddisc.NewDiscovered(&id,
			ddisc.DiscoveredOptionMimetype(mimex.Video),
			ddisc.DiscoveredOptionTitle("Ubuntu Documentary"),
		)

		seq := ddisc.Discover(
			context.Background(),
			ddisc.DefaultPolicy(),
			ddisc.DiscoverRequest{Query: "ubuntu"},
			[]ddisc.DiscoverOption{ddisc.DiscoverOptionKnownMediaDynamic(ddisc.KnownMediaDynamic(q))},
			fakeDiscoverStrategy{results: []ddisc.Discovered{hit}},
		)

		var got []ddisc.Discovered
		for d := range seq.Each(context.Background()) {
			got = append(got, d)
		}
		require.NoError(t, seq.Err())
		require.Len(t, got, 1)
		require.Equal(t, uuid.Nil.String(), got[0].KnownMediaID)
		require.Equal(t, 0, sqltestx.Count(t, q, "SELECT COUNT(*) FROM library_known_media"), "a candidate with nothing to show in the UI isn't worth cataloging")
	})

	t.Run("does nothing when the option is omitted", func(t *testing.T) {
		kid := uuid.Must(uuid.NewV4()).String()
		hit := newDiscoverHit(kid)

		seq := ddisc.Discover(
			context.Background(),
			ddisc.DefaultPolicy(),
			ddisc.DiscoverRequest{KnownMediaID: kid},
			nil,
			fakeDiscoverStrategy{results: []ddisc.Discovered{hit}},
		)

		var count int
		for range seq.Each(context.Background()) {
			count++
		}
		require.NoError(t, seq.Err())
		require.Equal(t, 1, count)
	})
}
