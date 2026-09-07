package ddisc_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/stretchr/testify/require"
)

func TestDiscoveredOptionHealth(t *testing.T) {
	d := ddisc.Discovered{}
	ddisc.DiscoveredOptionHealth(42)(&d)
	require.EqualValues(t, 42, d.Health)
}
