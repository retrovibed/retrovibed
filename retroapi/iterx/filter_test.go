package iterx_test

import (
	"context"
	"errors"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/stretchr/testify/require"
)

func TestFilterKeepsOnlyMatchingValues(t *testing.T) {
	seq := iterx.Filter[int](fakeSeq{values: []int{1, 2, 3, 4}}, func(ctx context.Context, v int) (bool, error) {
		return v%2 == 0, nil
	})

	var got []int
	for v := range seq.Each(context.Background()) {
		got = append(got, v)
	}

	require.Equal(t, []int{2, 4}, got)
	require.NoError(t, seq.Err())
}

func TestFilterPredicateErrorAbortsAndSurfaces(t *testing.T) {
	cause := errors.New("boom")
	seq := iterx.Filter[int](fakeSeq{values: []int{1, 2, 3}}, func(ctx context.Context, v int) (bool, error) {
		if v == 2 {
			return false, cause
		}
		return true, nil
	})

	var got []int
	for v := range seq.Each(context.Background()) {
		got = append(got, v)
	}

	require.Equal(t, []int{1}, got)
	require.ErrorIs(t, seq.Err(), cause)
}

func TestFilterPropagatesInnerSequenceError(t *testing.T) {
	cause := errors.New("inner boom")
	seq := iterx.Filter[int](fakeSeq{values: []int{1, 2}, err: cause}, func(ctx context.Context, v int) (bool, error) {
		return true, nil
	})

	for range seq.Each(context.Background()) {
	}

	require.ErrorIs(t, seq.Err(), cause)
}
