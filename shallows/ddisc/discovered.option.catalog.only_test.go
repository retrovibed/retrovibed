package ddisc_test

import (
	"math"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/stretchr/testify/require"
)

func TestDiscoveredOptionCatalogOnly(t *testing.T) {
	d := ddisc.Discovered{}
	ddisc.DiscoveredOptionCatalogOnly(&d)

	require.Equal(t, ddisc.PolicyRejectionCatalogOnly, d.PolicyRejection)
	require.EqualValues(t, math.MaxUint16, d.PolicyRank)
}
