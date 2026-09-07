package ddisc_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/stretchr/testify/require"
)

func TestDiscoveredOptionContentMime(t *testing.T) {
	d := ddisc.Discovered{}
	ddisc.DiscoveredOptionContentMime(mimex.HTTP)(&d)
	require.Equal(t, mimex.HTTP, d.Contentmime)
}
