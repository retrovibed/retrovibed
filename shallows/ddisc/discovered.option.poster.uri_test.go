package ddisc_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/stretchr/testify/require"
)

func TestDiscoveredOptionPosterURI(t *testing.T) {
	d := ddisc.Discovered{}
	ddisc.DiscoveredOptionPosterURI("http://example.com/poster.jpg")(&d)
	require.Equal(t, "http://example.com/poster.jpg", d.PosterURI)
}
