package ddisc_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/stretchr/testify/require"
)

func TestDiscoverRequestOptionAdult(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		r := ddisc.DiscoverRequest{}
		ddisc.DiscoverRequestOptionAdult(true)(&r)
		require.True(t, r.Adult)
	})

	t.Run("false", func(t *testing.T) {
		r := ddisc.DiscoverRequest{Adult: true}
		ddisc.DiscoverRequestOptionAdult(false)(&r)
		require.False(t, r.Adult)
	})
}
