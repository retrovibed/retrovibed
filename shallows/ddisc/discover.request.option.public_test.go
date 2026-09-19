package ddisc_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/stretchr/testify/require"
)

func TestDiscoverRequestOptionPublic(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		r := ddisc.DiscoverRequest{}
		ddisc.DiscoverRequestOptionPublic(true)(&r)
		require.True(t, r.Public)
	})

	t.Run("false", func(t *testing.T) {
		r := ddisc.DiscoverRequest{Public: true}
		ddisc.DiscoverRequestOptionPublic(false)(&r)
		require.False(t, r.Public)
	})
}
