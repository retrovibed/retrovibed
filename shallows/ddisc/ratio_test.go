package ddisc_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/cryptox"
	"github.com/stretchr/testify/require"
)

func TestFilterNone(t *testing.T) {
	t.Run("returns false", func(t *testing.T) {
		require.False(t, ddisc.FilterNone(nil))
		require.False(t, ddisc.FilterNone([]byte("test")))
		require.False(t, ddisc.FilterNone([]byte("another test")))
	})
}

func TestFilterRatio(t *testing.T) {
	t.Run("creates a filter with the correct ratio", func(t *testing.T) {
		filter := ddisc.FilterRatio(cryptox.NewChaCha8(t.Name()), 50)
		require.Equal(t, 100, len(filter))

		var trueCount, falseCount int
		for _, node := range filter {
			if node.N {
				trueCount++
			} else {
				falseCount++
			}
		}

		require.Equal(t, 50, trueCount)
		require.Equal(t, 50, falseCount)
	})

	t.Run("creates a filter with a zero ratio", func(t *testing.T) {
		filter := ddisc.FilterRatio(cryptox.NewChaCha8(t.Name()), 0)
		require.Equal(t, 100, len(filter))

		var filtered int
		for _, node := range filter {
			if node.N {
				filtered++
			}
		}

		require.Equal(t, 100, filtered)
	})

	t.Run("creates a filter with a full ratio", func(t *testing.T) {
		filter := ddisc.FilterRatio(cryptox.NewChaCha8(t.Name()), 100)
		require.Equal(t, 100, len(filter))

		var filtered int
		for _, node := range filter {
			if node.N {
				filtered++
			}
		}

		require.Equal(t, 0, filtered)
	})
}
