package ddisc_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/stretchr/testify/require"
)

func TestDiscoveredOptionURI(t *testing.T) {
	d := ddisc.Discovered{}
	ddisc.DiscoveredOptionURI("https://tracker.example/download/1.torrent")(&d)
	require.Equal(t, "https://tracker.example/download/1.torrent", d.URI)
}
