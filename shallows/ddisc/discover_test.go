package ddisc_test

import (
	"context"
	"database/sql"
	"errors"
	"iter"
	"testing"

	_ "github.com/duckdb/duckdb-go/v2"
	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
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
// Discover ranks every candidate it yields, so fakeDiscoverStrategy results
// must satisfy the same shape a real strategy would produce rather than
// being bare ID-only structs.
func newDiscoverHit(kid string) ddisc.Discovered {
	id := int160.Random()
	return ddisc.NewDiscovered(&id,
		ddisc.DiscoveredOptionKnownMedia(kid),
		ddisc.DiscoveredOptionMimetype(mimex.Binary),
		ddisc.DiscoveredOptionTitle("Some.Movie.2024.1080p.BluRay.x264"),
		ddisc.DiscoveredOptionAutoMagnet,
	)
}

func TestDiscoverTriesEveryStrategy(t *testing.T) {
	kid := uuid.Must(uuid.NewV4()).String()
	hit1 := newDiscoverHit(kid)
	hit2 := newDiscoverHit(kid)
	called2 := false

	seq := ddisc.Discover(
		context.Background(),
		ddisc.DefaultPolicy(),
		ddisc.DiscoverRequest{KnownMediaID: kid},
		nil,
		fakeDiscoverStrategy{results: []ddisc.Discovered{hit1}},
		fakeDiscoverStrategy{results: []ddisc.Discovered{hit2}, called: &called2},
	)

	var got []ddisc.Discovered
	for d := range seq.Each(context.Background()) {
		got = append(got, d)
	}
	require.NoError(t, seq.Err())
	require.True(t, called2, "second strategy should still run even though the first already found something")
	require.Len(t, got, 2)
	require.Equal(t, hit1.ID, got[0].ID)
	require.Equal(t, hit2.ID, got[1].ID)
}

func TestDiscoverFallsThroughOnMiss(t *testing.T) {
	kid := uuid.Must(uuid.NewV4()).String()
	hit := newDiscoverHit(kid)

	seq := ddisc.Discover(
		context.Background(),
		ddisc.DefaultPolicy(),
		ddisc.DiscoverRequest{KnownMediaID: kid},
		nil,
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
	cause := errors.New("boom")
	called2 := false

	seq := ddisc.Discover(
		context.Background(),
		ddisc.DefaultPolicy(),
		ddisc.DiscoverRequest{KnownMediaID: uuid.Must(uuid.NewV4()).String()},
		nil,
		fakeDiscoverStrategy{err: cause},
		fakeDiscoverStrategy{called: &called2},
	)

	for range seq.Each(context.Background()) {
	}
	require.ErrorIs(t, seq.Err(), cause)
	require.False(t, called2, "chain should stop on a genuine error, not fall through")
}

func TestDiscoverEmptyWhenAllStrategiesMiss(t *testing.T) {
	seq := ddisc.Discover(
		context.Background(),
		ddisc.DefaultPolicy(),
		ddisc.DiscoverRequest{KnownMediaID: uuid.Must(uuid.NewV4()).String()},
		nil,
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

func TestTitleFilterCleansHostileFreeText(t *testing.T) {
	db, err := sql.Open("duckdb", "")
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	id := int160.Random()
	candidate := ddisc.NewDiscovered(&id, ddisc.DiscoveredOptionTitle("Arch Linux ISO"))

	// a trailing ":" is a dangling field:value token in lucene grammar -
	// free-typed UI search text hits this whenever a query is submitted
	// mid-word, so TitleFilter must sanitize it rather than handing it to
	// the parser raw (previously surfaced as a DuckDB "Parser Error: syntax
	// error at end of input").
	filter := ddisc.NewTitleFilter(db, ddisc.DiscoverRequest{Query: "arch linux:"})
	matched, err := filter.Match(context.Background(), candidate)
	require.NoError(t, err)
	require.True(t, matched)
}

func TestDiscoverRanksYieldedCandidates(t *testing.T) {
	kid := uuid.Must(uuid.NewV4()).String()
	hit := newDiscoverHit(kid)

	seq := ddisc.Discover(
		context.Background(),
		ddisc.DefaultPolicy(),
		ddisc.DiscoverRequest{KnownMediaID: kid},
		nil,
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
}
