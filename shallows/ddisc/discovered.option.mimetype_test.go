package ddisc_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/stretchr/testify/require"
)

func TestDiscoveredOptionMimetype(t *testing.T) {
	tcs := []struct {
		name             string
		input            []string
		expectedMimetype string
		expectedCategory string
	}{
		{name: "uses the first non-zero value", input: []string{"", "video/mp4", "audio/mpeg"}, expectedMimetype: "video/mp4", expectedCategory: mimex.Video},
		{name: "falls back to binary when all values are zero", input: []string{"", ""}, expectedMimetype: mimex.Binary, expectedCategory: mimex.Application},
		{name: "falls back to binary when no values are given", input: nil, expectedMimetype: mimex.Binary, expectedCategory: mimex.Application},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			d := ddisc.Discovered{}
			ddisc.DiscoveredOptionMimetype(tc.input...)(&d)

			require.Equal(t, tc.expectedMimetype, d.Mimetype)
			require.Equal(t, tc.expectedCategory, d.Category)
		})
	}
}
