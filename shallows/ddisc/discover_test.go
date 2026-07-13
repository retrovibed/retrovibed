package ddisc_test

import (
	"context"
	"errors"
	"iter"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/iterx"
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

func TestDiscoverStopsAtFirstHit(t *testing.T) {
	hit := ddisc.Discovered{ID: "found"}
	called2 := false

	seq := ddisc.Discover(
		context.Background(),
		ddisc.DiscoverRequest{KnownMediaID: "x"},
		fakeDiscoverStrategy{results: []ddisc.Discovered{hit}},
		fakeDiscoverStrategy{called: &called2},
	)

	var got []ddisc.Discovered
	for d := range seq.Each(context.Background()) {
		got = append(got, d)
	}
	require.NoError(t, seq.Err())
	require.Equal(t, []ddisc.Discovered{hit}, got)
	require.False(t, called2, "second strategy should not run once the first finds something")
}

func TestDiscoverFallsThroughOnMiss(t *testing.T) {
	hit := ddisc.Discovered{ID: "found"}

	seq := ddisc.Discover(
		context.Background(),
		ddisc.DiscoverRequest{KnownMediaID: "x"},
		fakeDiscoverStrategy{},
		fakeDiscoverStrategy{results: []ddisc.Discovered{hit}},
	)

	var got []ddisc.Discovered
	for d := range seq.Each(context.Background()) {
		got = append(got, d)
	}
	require.NoError(t, seq.Err())
	require.Equal(t, []ddisc.Discovered{hit}, got)
}

func TestDiscoverPropagatesGenuineErrorAndStopsChain(t *testing.T) {
	cause := errors.New("boom")
	called2 := false

	seq := ddisc.Discover(
		context.Background(),
		ddisc.DiscoverRequest{KnownMediaID: "x"},
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
		ddisc.DiscoverRequest{KnownMediaID: "x"},
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
