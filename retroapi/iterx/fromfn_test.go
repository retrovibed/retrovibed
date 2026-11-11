package iterx_test

import (
	"context"
	"errors"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/stretchr/testify/require"
)

func TestFromFn(t *testing.T) {
	t.Run("yields values and reports the returned error", func(t *testing.T) {
		boom := errors.New("boom")
		seq := iterx.FromFn(func(ctx context.Context, yield func(int) bool) error {
			for i := 1; i <= 3; i++ {
				if !yield(i) {
					return nil
				}
			}

			return boom
		})

		var got []int
		for v := range seq.Each(context.Background()) {
			got = append(got, v)
		}

		require.Equal(t, []int{1, 2, 3}, got)
		require.ErrorIs(t, seq.Err(), boom)
	})

	t.Run("discards the error when the consumer stops early", func(t *testing.T) {
		seq := iterx.FromFn(func(ctx context.Context, yield func(int) bool) error {
			for i := 1; i <= 3; i++ {
				if !yield(i) {
					return context.Canceled
				}
			}

			return nil
		})

		var got []int
		for v := range iterx.Take(seq, 2).Each(context.Background()) {
			got = append(got, v)
		}

		require.Equal(t, []int{1, 2}, got)
		require.NoError(t, seq.Err())
	})
}
