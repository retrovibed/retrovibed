package iterx_test

import (
	"context"
	"errors"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/stretchr/testify/require"
)

func TestTakeStopsAfterNElements(t *testing.T) {
	seq := iterx.Take[int](fakeSeq{values: []int{1, 2, 3, 4, 5}}, 3)

	var got []int
	for v := range seq.Each(context.Background()) {
		got = append(got, v)
	}

	require.Equal(t, []int{1, 2, 3}, got)
	require.NoError(t, seq.Err())
}

func TestTakeZeroYieldsNothing(t *testing.T) {
	seq := iterx.Take[int](fakeSeq{values: []int{1, 2, 3}}, 0)

	var got []int
	for v := range seq.Each(context.Background()) {
		got = append(got, v)
	}

	require.Empty(t, got)
}

func TestTakeNGreaterThanLengthYieldsAll(t *testing.T) {
	seq := iterx.Take[int](fakeSeq{values: []int{1, 2}}, 5)

	var got []int
	for v := range seq.Each(context.Background()) {
		got = append(got, v)
	}

	require.Equal(t, []int{1, 2}, got)
}

func TestTakeSurfacesInnerError(t *testing.T) {
	cause := errors.New("boom")
	seq := iterx.Take[int](fakeSeq{values: []int{1, 2}, err: cause}, 5)

	for range seq.Each(context.Background()) {
	}

	require.ErrorIs(t, seq.Err(), cause)
}
