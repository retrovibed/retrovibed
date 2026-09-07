package ddisc_test

import (
	"testing"
	"time"

	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/stretchr/testify/require"
)

func TestDiscoveredOptionAutoReleased(t *testing.T) {
	now := time.Now()
	epoch := time.Unix(0, 0)

	tcs := []struct {
		name     string
		input    time.Time
		expected time.Time
	}{
		{name: "zero defaults to neg-infinity", input: time.Time{}, expected: timex.NegInf()},
		{name: "infinity is preserved", input: timex.Inf(), expected: timex.Inf()},
		{name: "neg-infinity is preserved", input: timex.NegInf(), expected: timex.NegInf()},
		{name: "epoch is preserved", input: epoch, expected: epoch},
		{name: "now is preserved", input: now, expected: now},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			d := &ddisc.Discovered{ReleasedAt: tc.input}
			ddisc.DiscoveredOptionAutoReleased(d)

			require.True(t, d.ReleasedAt.Equal(tc.expected))
		})
	}
}
