package ddisc_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/stretchr/testify/require"
)

func TestDiscoveredOptionPrivate(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		d := ddisc.Discovered{}
		ddisc.DiscoveredOptionPrivate(true)(&d)
		require.True(t, d.Private)
	})

	t.Run("false", func(t *testing.T) {
		d := ddisc.Discovered{Private: true}
		ddisc.DiscoveredOptionPrivate(false)(&d)
		require.False(t, d.Private)
	})
}
