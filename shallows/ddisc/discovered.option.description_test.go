package ddisc_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/stretchr/testify/require"
)

func TestDiscoveredOptionDescription(t *testing.T) {
	d := ddisc.Discovered{}
	ddisc.DiscoveredOptionDescription("the description")(&d)
	require.Equal(t, "the description", d.Description)
}
