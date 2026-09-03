package ddisc_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/stretchr/testify/require"
)

func TestDiscoveredOptionInfoHash(t *testing.T) {
	d := ddisc.Discovered{}
	want := []byte{1, 2, 3, 4}
	ddisc.DiscoveredOptionInfoHash(want)(&d)
	require.Equal(t, want, d.Infohash)
}
