package ddisc_test

import (
	"context"
	"errors"
	"iter"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/stretchr/testify/require"
)

type fakeDiscoverStrategy struct {
	results []ddisc.Discovered
	err     error
	called  *bool
}

func (t fakeDiscoverStrategy) Discover(ctx context.Context, req ddisc.DiscoverRequest) iterx.Seq[ddisc.Discovered] {
	if t.called != nil {
		*t.called = true
	}
	return fakeDiscoverSeq{results: t.results, err: t.err}
}

type fakeDiscoverSeq struct {
	results []ddisc.Discovered
	err     error
}

func (t fakeDiscoverSeq) Each(ctx context.Context) iter.Seq[ddisc.Discovered] {
	return iterx.From(t.results...)
}

func (t fakeDiscoverSeq) Err() error { return t.err }

// newDiscoverHit builds a fully valid, unpersisted Discovered candidate -
// Discover's central seq now persists and ranks every candidate it yields,
// so fakeDiscoverStrategy results must satisfy ddisc_media's constraints
// rather than being bare ID-only structs.
func newDiscoverHit(kid string) ddisc.Discovered {
	id := int160.Random()
	return ddisc.NewDiscovered(&id,
		ddisc.DiscoveredOptionKnownMedia(kid),
		ddisc.DiscoveredOptionMimetype(mimex.Binary),
		ddisc.DiscoveredOptionTitle("Some.Movie.2024.1080p.BluRay.x264"),
		ddisc.DiscoveredOptionAutoMagnet,
	)
}

func TestDiscoverStopsAtFirstHit(t *testing.T) {
	q := sqltestx.Metadatabase(t)
	kid := uuid.Must(uuid.NewV4()).String()
	hit := newDiscoverHit(kid)
	called2 := false

	seq := ddisc.Discover(
		context.Background(),
		q,
		ddisc.DefaultPolicy(),
		ddisc.DiscoverRequest{KnownMediaID: kid},
		fakeDiscoverStrategy{results: []ddisc.Discovered{hit}},
		fakeDiscoverStrategy{called: &called2},
	)

	var got []ddisc.Discovered
	for d := range seq.Each(context.Background()) {
		got = append(got, d)
	}
	require.NoError(t, seq.Err())
	require.Len(t, got, 1)
	require.Equal(t, hit.ID, got[0].ID)
	require.False(t, called2, "second strategy should not run once the first finds something")
	require.Equal(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM ddisc_media WHERE id = ?", got[0].ID))
}

func TestDiscoverFallsThroughOnMiss(t *testing.T) {
	q := sqltestx.Metadatabase(t)
	kid := uuid.Must(uuid.NewV4()).String()
	hit := newDiscoverHit(kid)

	seq := ddisc.Discover(
		context.Background(),
		q,
		ddisc.DefaultPolicy(),
		ddisc.DiscoverRequest{KnownMediaID: kid},
		fakeDiscoverStrategy{},
		fakeDiscoverStrategy{results: []ddisc.Discovered{hit}},
	)

	var got []ddisc.Discovered
	for d := range seq.Each(context.Background()) {
		got = append(got, d)
	}
	require.NoError(t, seq.Err())
	require.Len(t, got, 1)
	require.Equal(t, hit.ID, got[0].ID)
}

func TestDiscoverPropagatesGenuineErrorAndStopsChain(t *testing.T) {
	q := sqltestx.Metadatabase(t)
	cause := errors.New("boom")
	called2 := false

	seq := ddisc.Discover(
		context.Background(),
		q,
		ddisc.DefaultPolicy(),
		ddisc.DiscoverRequest{KnownMediaID: uuid.Must(uuid.NewV4()).String()},
		fakeDiscoverStrategy{err: cause},
		fakeDiscoverStrategy{called: &called2},
	)

	for range seq.Each(context.Background()) {
	}
	require.ErrorIs(t, seq.Err(), cause)
	require.False(t, called2, "chain should stop on a genuine error, not fall through")
}

func TestDiscoverEmptyWhenAllStrategiesMiss(t *testing.T) {
	q := sqltestx.Metadatabase(t)

	seq := ddisc.Discover(
		context.Background(),
		q,
		ddisc.DefaultPolicy(),
		ddisc.DiscoverRequest{KnownMediaID: uuid.Must(uuid.NewV4()).String()},
		fakeDiscoverStrategy{},
		fakeDiscoverStrategy{},
	)

	var count int
	for range seq.Each(context.Background()) {
		count++
	}
	require.NoError(t, seq.Err())
	require.Equal(t, 0, count)
}

func TestDiscoverRanksYieldedCandidates(t *testing.T) {
	q := sqltestx.Metadatabase(t)
	kid := uuid.Must(uuid.NewV4()).String()
	hit := newDiscoverHit(kid)

	seq := ddisc.Discover(
		context.Background(),
		q,
		ddisc.DefaultPolicy(),
		ddisc.DiscoverRequest{KnownMediaID: kid},
		fakeDiscoverStrategy{results: []ddisc.Discovered{hit}},
	)

	var got []ddisc.Discovered
	for d := range seq.Each(context.Background()) {
		got = append(got, d)
	}
	require.NoError(t, seq.Err())
	require.Len(t, got, 1)
	require.NotEqual(t, uint16(0), got[0].PolicyRank)
	require.Less(t, got[0].PolicyRank, uint16(65535), "a clean title should rank better than the unranked sentinel")

	var persisted ddisc.Discovered
	require.NoError(t, ddisc.DiscoveredFindByID(context.Background(), q, got[0].ID).Scan(&persisted))
	require.Equal(t, got[0].PolicyRank, persisted.PolicyRank)
}
