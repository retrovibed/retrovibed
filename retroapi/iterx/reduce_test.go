package iterx_test

import (
	"context"
	"errors"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/stretchr/testify/require"
)

func TestReduce(t *testing.T) {
	t.Run("folds all values into the accumulator", func(t *testing.T) {
		got, err := iterx.Reduce[int, int](context.Background(), fakeSeq{values: []int{1, 2, 3, 4}}, 10, func(ctx context.Context, acc int, v int) (int, error) {
			return acc + v, nil
		})

		require.NoError(t, err)
		require.Equal(t, 20, got)
	})

	t.Run("accumulator type may differ from element type", func(t *testing.T) {
		got, err := iterx.Reduce[int, []string](context.Background(), fakeSeq{values: []int{1, 2}}, nil, func(ctx context.Context, acc []string, v int) ([]string, error) {
			return append(acc, string(rune('a'+v)), "-"), nil
		})

		require.NoError(t, err)
		require.Equal(t, []string{"b", "-", "c", "-"}, got)
	})

	t.Run("empty sequence returns init", func(t *testing.T) {
		got, err := iterx.Reduce[int, int](context.Background(), fakeSeq{}, 7, func(ctx context.Context, acc int, v int) (int, error) {
			return acc + v, nil
		})

		require.NoError(t, err)
		require.Equal(t, 7, got)
	})

	t.Run("fn error aborts and returns accumulator so far", func(t *testing.T) {
		cause := errors.New("boom")
		calls := 0
		got, err := iterx.Reduce[int, int](context.Background(), fakeSeq{values: []int{1, 2, 3}}, 0, func(ctx context.Context, acc int, v int) (int, error) {
			calls++
			if v == 2 {
				return 0, cause
			}
			return acc + v, nil
		})

		require.ErrorIs(t, err, cause)
		require.Equal(t, 1, got)
		require.Equal(t, 2, calls)
	})

	t.Run("surfaces inner sequence error", func(t *testing.T) {
		cause := errors.New("inner boom")
		_, err := iterx.Reduce[int, int](context.Background(), fakeSeq{values: []int{1, 2}, err: cause}, 0, func(ctx context.Context, acc int, v int) (int, error) {
			return acc + v, nil
		})

		require.ErrorIs(t, err, cause)
	})
}
