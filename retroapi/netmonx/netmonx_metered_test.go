package netmonx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsMeteredInterface(t *testing.T) {
	cases := []struct {
		name    string
		iface   string
		metered bool
	}{
		{name: "wwan prefix", iface: "wwan0", metered: true},
		{name: "wwan with suffix", iface: "wwan0.1", metered: true},
		{name: "rmnet prefix", iface: "rmnet0", metered: true},
		{name: "rmnet android prefixed", iface: "v4-rmnet1", metered: true},
		{name: "rmnet android bare", iface: "rmnet_data0", metered: true},
		{name: "pdp_ip prefix", iface: "pdp_ip0", metered: true},
		{name: "ccmni prefix", iface: "ccmni0", metered: true},
		{name: "eth not metered", iface: "eth0", metered: false},
		{name: "wifi not metered", iface: "wlan0", metered: false},
		{name: "loopback not metered", iface: "lo", metered: false},
		{name: "empty not metered", iface: "", metered: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.metered, isMeteredInterface(tc.iface))
		})
	}
}
