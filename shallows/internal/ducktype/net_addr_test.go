package ducktype_test

import (
	"math/big"
	"net/netip"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/internal/ducktype"
	"github.com/stretchr/testify/require"
)

func TestNullNetAddr(t *testing.T) {
	t.Run("Scan", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			var n ducktype.NullNetAddr
			require.NoError(t, n.Scan(nil))
			require.False(t, n.Valid)
			require.Equal(t, netip.Addr{}, n.V)
		})

		t.Run("string IPv4", func(t *testing.T) {
			var n ducktype.NullNetAddr
			require.NoError(t, n.Scan("192.168.1.1"))
			require.True(t, n.Valid)
			require.Equal(t, netip.MustParseAddr("192.168.1.1"), n.V)
		})

		t.Run("string IPv6", func(t *testing.T) {
			var n ducktype.NullNetAddr
			require.NoError(t, n.Scan("::1"))
			require.True(t, n.Valid)
			require.Equal(t, netip.MustParseAddr("::1"), n.V)
		})

		t.Run("string IPv4 broadcast", func(t *testing.T) {
			var n ducktype.NullNetAddr
			require.NoError(t, n.Scan("255.255.255.255"))
			require.True(t, n.Valid)
			require.Equal(t, netip.MustParseAddr("255.255.255.255"), n.V)
		})

		t.Run("string IPv4 unspecified", func(t *testing.T) {
			var n ducktype.NullNetAddr
			require.NoError(t, n.Scan("0.0.0.0"))
			require.True(t, n.Valid)
			require.Equal(t, netip.MustParseAddr("0.0.0.0"), n.V)
		})

		t.Run("string IPv6 unspecified", func(t *testing.T) {
			var n ducktype.NullNetAddr
			require.NoError(t, n.Scan("::"))
			require.True(t, n.Valid)
			require.Equal(t, netip.MustParseAddr("::"), n.V)
		})

		t.Run("string IPv4 loopback", func(t *testing.T) {
			var n ducktype.NullNetAddr
			require.NoError(t, n.Scan("127.0.0.1"))
			require.True(t, n.Valid)
			require.Equal(t, netip.MustParseAddr("127.0.0.1"), n.V)
		})

		t.Run("string IPv4 multicast", func(t *testing.T) {
			var n ducktype.NullNetAddr
			require.NoError(t, n.Scan("224.0.0.1"))
			require.True(t, n.Valid)
			require.Equal(t, netip.MustParseAddr("224.0.0.1"), n.V)
		})

		t.Run("string IPv6 multicast", func(t *testing.T) {
			var n ducktype.NullNetAddr
			require.NoError(t, n.Scan("ff02::1"))
			require.True(t, n.Valid)
			require.Equal(t, netip.MustParseAddr("ff02::1"), n.V)
		})

		t.Run("string IPv4-mapped IPv6", func(t *testing.T) {
			var n ducktype.NullNetAddr
			require.NoError(t, n.Scan("::ffff:192.0.2.1"))
			require.True(t, n.Valid)
			require.Equal(t, netip.MustParseAddr("::ffff:192.0.2.1"), n.V)
		})

		t.Run("string invalid", func(t *testing.T) {
			var n ducktype.NullNetAddr
			require.Error(t, n.Scan("not-an-ip"))
		})

		t.Run("bytes IPv4", func(t *testing.T) {
			var n ducktype.NullNetAddr
			require.NoError(t, n.Scan([]byte{10, 0, 0, 1}))
			require.True(t, n.Valid)
			addr, ok := netip.AddrFromSlice([]byte{10, 0, 0, 1})
			require.True(t, ok)
			require.Equal(t, addr, n.V)
		})

		t.Run("map with string address", func(t *testing.T) {
			var n ducktype.NullNetAddr
			require.NoError(t, n.Scan(map[string]any{"address": "10.0.0.2"}))
			require.True(t, n.Valid)
			require.Equal(t, netip.MustParseAddr("10.0.0.2"), n.V)
		})

		t.Run("map with big.Int IPv4", func(t *testing.T) {
			var n ducktype.NullNetAddr
			b := new(big.Int).SetBytes([]byte{192, 168, 0, 1})
			require.NoError(t, n.Scan(map[string]any{"address": b}))
			require.True(t, n.Valid)
			addr, ok := netip.AddrFromSlice([]byte{192, 168, 0, 1})
			require.True(t, ok)
			require.Equal(t, addr, n.V)
		})

		t.Run("map with big.Int IPv6", func(t *testing.T) {
			var n ducktype.NullNetAddr
			b := new(big.Int).SetBytes([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1})
			require.NoError(t, n.Scan(map[string]any{"address": b}))
			require.True(t, n.Valid)
		})

		t.Run("map with unknown address type", func(t *testing.T) {
			var n ducktype.NullNetAddr
			require.Error(t, n.Scan(map[string]any{"address": 12345}))
		})

		t.Run("unsupported type", func(t *testing.T) {
			var n ducktype.NullNetAddr
			require.Error(t, n.Scan(12345))
		})
	})

	t.Run("Value", func(t *testing.T) {
		t.Run("invalid", func(t *testing.T) {
			n := ducktype.NullNetAddr{Valid: false}
			v, err := n.Value()
			require.NoError(t, err)
			require.Nil(t, v)
		})

		t.Run("valid", func(t *testing.T) {
			n := ducktype.NullNetAddr{V: netip.MustParseAddr("1.2.3.4"), Valid: true}
			v, err := n.Value()
			require.NoError(t, err)
			require.Equal(t, "1.2.3.4", v)
		})
	})

	t.Run("select from duckdb", func(t *testing.T) {
		t.Run("string", func(t *testing.T) {
			db := newDB(t)
			var n ducktype.NullNetAddr
			row := db.QueryRowContext(t.Context(), "SELECT '192.168.1.100'::VARCHAR")
			require.NoError(t, row.Scan(&n))
			require.True(t, n.Valid)
			require.Equal(t, netip.MustParseAddr("192.168.1.100"), n.V)
		})

		t.Run("null", func(t *testing.T) {
			db := newDB(t)
			var n ducktype.NullNetAddr
			row := db.QueryRowContext(t.Context(), "SELECT NULL::VARCHAR")
			require.NoError(t, row.Scan(&n))
			require.False(t, n.Valid)
		})
	})
}
