package netx_test

import (
	"net"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/internal/netx"
	"github.com/stretchr/testify/require"
)

func TestIPString(t *testing.T) {
	t.Run("nil returns empty string", func(t *testing.T) {
		require.Equal(t, "", netx.IPString(nil))
	})

	t.Run("non-nil returns address", func(t *testing.T) {
		require.Equal(t, "127.0.0.1", netx.IPString(net.ParseIP("127.0.0.1")))
	})
}

func TestFirstNonZeroIP(t *testing.T) {
	loopback := net.ParseIP("127.0.0.1")
	public := net.ParseIP("10.0.0.1")

	t.Run("returns first non-nil", func(t *testing.T) {
		result := netx.FirstNonZeroIP(nil, loopback, public)
		require.Equal(t, loopback, result)
	})

	t.Run("skips leading nils", func(t *testing.T) {
		result := netx.FirstNonZeroIP(nil, nil, public)
		require.Equal(t, public, result)
	})

	t.Run("all nil returns empty string via IPString", func(t *testing.T) {
		result := netx.FirstNonZeroIP(nil, nil)
		require.Nil(t, result)
		require.Equal(t, "", netx.IPString(result))
	})

	t.Run("no args returns empty string via IPString", func(t *testing.T) {
		result := netx.FirstNonZeroIP()
		require.Nil(t, result)
		require.Equal(t, "", netx.IPString(result))
	})

	t.Run("single non-nil", func(t *testing.T) {
		result := netx.FirstNonZeroIP(loopback)
		require.Equal(t, loopback, result)
	})
}
