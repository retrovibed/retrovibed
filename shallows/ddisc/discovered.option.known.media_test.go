package ddisc_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/stretchr/testify/require"
)

func TestDiscoveredOptionKnownMedia(t *testing.T) {
	d := ddisc.Discovered{}
	ddisc.DiscoveredOptionKnownMedia("known-media-id")(&d)
	require.Equal(t, "known-media-id", d.KnownMediaID)
}
