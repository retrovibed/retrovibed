package ddisc_test

import (
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestKnownMediaDetector(t *testing.T) {
	t.Run("should stamp a known-media-id onto an unresolved candidate", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionMimetype(mimex.Video)))
		known.Title = "The Grand Budapest Hotel"
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		id := int160.Random()
		cand := ddisc.NewDiscovered(&id,
			ddisc.DiscoveredOptionTitle(known.Title),
			ddisc.DiscoveredOptionMimetype(mimex.Video),
		)
		require.Equal(t, uuid.Nil.String(), cand.KnownMediaID, "candidate must start unresolved for the detector to have anything to do")

		seq := ddisc.KnownMediaDetector(q, library.QueryCleanerNoop())(fakeDiscoverSeq{results: []ddisc.Discovered{cand}})

		var got []ddisc.Discovered
		for d := range seq.Each(ctx) {
			got = append(got, d)
		}
		require.NoError(t, seq.Err())
		require.Len(t, got, 1)
		require.Equal(t, known.UID, got[0].KnownMediaID)
	})

	t.Run("should leave an already resolved candidate alone", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		kid := uuid.Must(uuid.NewV4()).String()
		id := int160.Random()
		cand := ddisc.NewDiscovered(&id,
			ddisc.DiscoveredOptionKnownMedia(kid),
			ddisc.DiscoveredOptionTitle("some completely unrelated title"),
			ddisc.DiscoveredOptionMimetype(mimex.Video),
		)

		seq := ddisc.KnownMediaDetector(q, library.QueryCleanerNoop())(fakeDiscoverSeq{results: []ddisc.Discovered{cand}})

		var got []ddisc.Discovered
		for d := range seq.Each(ctx) {
			got = append(got, d)
		}
		require.NoError(t, seq.Err())
		require.Len(t, got, 1)
		require.Equal(t, kid, got[0].KnownMediaID, "a candidate that already carries a known-media-id must not be overwritten")
	})
}
