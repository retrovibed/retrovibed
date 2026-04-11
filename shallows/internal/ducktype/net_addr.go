package ducktype

import (
	"database/sql/driver"
	"fmt"
	"math/big"
	"net/netip"

	"github.com/davecgh/go-spew/spew"
)

// NullNetAddr represents a netip.Addr that may be null.
// The V field holds the netip.Addr value, and Valid indicates its validity.
type NullNetAddr struct {
	V     netip.Addr
	Valid bool
}

// Scan implements the sql.Scanner interface.
// It supports scanning a value from a database driver, including NULL,
// and can handle database types that represent an IP address as a string or a byte slice,
// including values from `net.IP`.
func (n *NullNetAddr) Scan(src any) error {
	if src == nil {
		n.V, n.Valid = netip.Addr{}, false
		return nil
	}
	n.Valid = true
	switch v := src.(type) {
	case []byte:
		addr, ok := netip.AddrFromSlice(v)
		if !ok {
			return fmt.Errorf("NullNetAddr: cannot scan []byte %q into netip.Addr", v)
		}
		n.V = addr
		return nil
	case string:
		addr, err := netip.ParseAddr(v)
		if err != nil {
			return fmt.Errorf("NullNetAddr: failed to parse string %q as netip.Addr: %s", err, v)
		}
		n.V = addr
		return nil
	case map[string]any:
		switch _addr := v["address"].(type) {
		case string:
			addr, err := netip.ParseAddr(_addr)
			if err != nil {
				return fmt.Errorf("NullNetAddr: failed to parse string %q as netip.Addr: %s", err, v)
			}
			n.V = addr
			return nil
		case *big.Int:
			const (
				IPv6BitLen = 64
				IPv4BitLen = 32
			)
			var (
				decoded []byte
				ip4buf  [4]byte
				ip6buf  [16]byte
			)

			if _addr.BitLen() <= IPv4BitLen {
				decoded = ip4buf[:]
			} else {
				decoded = ip6buf[:]
			}

			addr, ok := netip.AddrFromSlice(_addr.FillBytes(decoded))
			if !ok {
				return fmt.Errorf("NullNetAddr: failed to convert big.Int as netip.Addr: %v - %v - %v", _addr, _addr.BitLen(), spew.Sdump(v))
			}
			n.V = addr
			return nil
		default:
			return fmt.Errorf("NullNetAddr: address returned is an unknown type: %T - %v", _addr, spew.Sdump(v))
		}
	default:
		n.Valid = false
		return fmt.Errorf("NullNetAddr: cannot scan type %T into NullNetIP %v", src, src)
	}
}

// Value implements the driver.Valuer interface.
// It returns nil if the value is not valid, otherwise it returns a byte slice
// representation of the netip.Addr for database storage.
func (n NullNetAddr) Value() (driver.Value, error) {
	if !n.Valid || !n.V.IsValid() {
		return nil, nil
	}

	return n.V.String(), nil
}
