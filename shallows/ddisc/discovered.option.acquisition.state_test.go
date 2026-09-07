package ddisc_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/stretchr/testify/require"
)

func TestDiscoveredOptionAcquisitionState(t *testing.T) {
	d := ddisc.Discovered{}
	ddisc.DiscoveredOptionAcquisitionState(ddisc.AcquisitionStateDownloading)(&d)
	require.EqualValues(t, ddisc.AcquisitionStateDownloading, d.AcquisitionState)
}
