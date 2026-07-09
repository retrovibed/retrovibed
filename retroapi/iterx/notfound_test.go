package iterx_test

import (
	"context"
	"errors"
	"iter"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/stretchr/testify/require"
)

type fakeSeq struct {
	values []int
	err    error
}

func (t fakeSeq) Each(ctx context.Context) iter.Seq[int] {
	return func(yield func(int) bool) {
		for _, v := range t.values {
			if !yield(v) {
				return
			}
		}
	}
}

func (t fakeSeq) Err() error { return t.err }

func TestNotFoundEmptySequenceYieldsErrNotFound(t *testing.T) {
	seq := iterx.NotFound[int](fakeSeq{})

	var got []int
	for v := range seq.Each(context.Background()) {
		got = append(got, v)
	}

	require.Empty(t, got)
	require.ErrorIs(t, seq.Err(), iterx.ErrNotFound)
}

func TestNotFoundRealErrorPassesThroughUnchanged(t *testing.T) {
	cause := errors.New("boom")
	seq := iterx.NotFound[int](fakeSeq{err: cause})

	for range seq.Each(context.Background()) {
	}

	require.ErrorIs(t, seq.Err(), cause)
	require.NotErrorIs(t, seq.Err(), iterx.ErrNotFound)
}

func TestNotFoundYieldsAndReturnsNilWhenResultsExist(t *testing.T) {
	seq := iterx.NotFound[int](fakeSeq{values: []int{1, 2, 3}})

	var got []int
	for v := range seq.Each(context.Background()) {
		got = append(got, v)
	}

	require.Equal(t, []int{1, 2, 3}, got)
	require.NoError(t, seq.Err())
}
