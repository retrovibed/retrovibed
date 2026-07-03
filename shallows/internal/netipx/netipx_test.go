package netipx_test

import (
	"net/netip"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/internal/netipx"
	"github.com/stretchr/testify/require"
)

func TestAddrPortFromStrings(t *testing.T) {
	t.Run("keeps valid addresses", func(t *testing.T) {
		result := netipx.AddrPortFromStrings("127.0.0.1:80", "10.0.0.1:443")
		require.Equal(t, []netip.AddrPort{
			netip.MustParseAddrPort("127.0.0.1:80"),
			netip.MustParseAddrPort("10.0.0.1:443"),
		}, result)
	})

	t.Run("keeps valid ipv6 address", func(t *testing.T) {
		result := netipx.AddrPortFromStrings("[::1]:8080")
		require.Equal(t, []netip.AddrPort{
			netip.MustParseAddrPort("[::1]:8080"),
		}, result)
	})

	t.Run("skips invalid strings", func(t *testing.T) {
		result := netipx.AddrPortFromStrings("not-an-addr", "127.0.0.1:80")
		require.Equal(t, []netip.AddrPort{
			netip.MustParseAddrPort("127.0.0.1:80"),
		}, result)
	})

	t.Run("all invalid returns empty", func(t *testing.T) {
		result := netipx.AddrPortFromStrings("bad", "worse")
		require.Empty(t, result)
	})

	t.Run("no input returns empty", func(t *testing.T) {
		result := netipx.AddrPortFromStrings()
		require.Empty(t, result)
	})
}
