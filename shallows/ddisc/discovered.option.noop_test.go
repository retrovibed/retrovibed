package ddisc_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/stretchr/testify/require"
)

func TestDiscoveredOptionNoop(t *testing.T) {
	before := ddisc.Discovered{Title: "unchanged", Mimetype: "video/mp4"}
	after := before

	ddisc.DiscoveredOptionNoop(&after)

	require.Equal(t, before, after)
}
