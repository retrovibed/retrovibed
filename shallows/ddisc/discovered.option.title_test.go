package ddisc_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/stretchr/testify/require"
)

func TestDiscoveredOptionTitle(t *testing.T) {
	d := ddisc.Discovered{}
	ddisc.DiscoveredOptionTitle("the title")(&d)
	require.Equal(t, "the title", d.Title)
}
