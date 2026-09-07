package ddisc_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/stretchr/testify/require"
)

func TestDiscoveredOptionSource(t *testing.T) {
	d := ddisc.Discovered{}
	ddisc.DiscoveredOptionSource("the source")(&d)
	require.Equal(t, "the source", d.Source)
}
