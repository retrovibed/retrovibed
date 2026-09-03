package ddisc_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/stretchr/testify/require"
)

func TestDiscoveredOptionPartition(t *testing.T) {
	d := ddisc.Discovered{}
	ddisc.DiscoveredOptionPartition("the-partition")(&d)
	require.Equal(t, "the-partition", d.Partition)
}
