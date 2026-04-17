package netx

import (
	"context"
	"errors"
	"log"
	"net"
	"net/netip"
	"strconv"
	"strings"
)

// Dialer missing interface from the net package.
type Dialer interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

func DefaultIfNil(d0, d1 Dialer) Dialer {
	if d0 != nil {
		return d0
	}

	return d1
}

func DefaultIfZero(fallback net.IP, v net.IP) net.IP {
	if v != nil {
		return v
	}

	return fallback
}

// IPString returns ip.String() or "" if ip is nil.
func IPString(ip net.IP) string {
	if ip == nil {
		return ""
	}
	return ip.String()
}

// FirstNonZeroIP returns the first non-nil net.IP from the provided candidates,
// or nil if all are nil. A nil net.IP result is safe to call .String() on.
func FirstNonZeroIP(candidates ...net.IP) net.IP {
	for _, ip := range candidates {
		if ip != nil {
			return ip
		}
	}
	return nil
}

// HostIP ...
func HostIP(host string) net.IP {
	ip, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		log.Println("failed to resolve ip for", host, "falling back to 127.0.0.1:", err)
		return net.ParseIP("127.0.0.1")
	}

	return ip.IP
}

func Port(s string) (p uint16, err error) {
	var (
		sport string
		port  uint64
	)

	if _, sport, err = net.SplitHostPort(s); err != nil {
		return 0, err
	}

	if port, err = strconv.ParseUint(sport, 10, 16); err != nil {
		return 0, err
	}

	return uint16(port), nil
}

func IP(s string) net.IP {
	var (
		err  error
		host string
	)

	if host, _, err = net.SplitHostPort(s); err != nil {
		log.Println("unable to parse host", host)
		return nil
	}

	return HostIP(host)
}

func AddrPort(a net.Addr) *netip.AddrPort {
	switch v := a.(type) {
	case *net.TCPAddr:
		ip, _ := netip.AddrFromSlice(v.IP)
		tmp := netip.AddrPortFrom(ip, uint16(v.Port))
		return &tmp
	case *net.UDPAddr:
		ip, _ := netip.AddrFromSlice(v.IP)
		tmp := netip.AddrPortFrom(ip, uint16(v.Port))
		return &tmp
	default:
		log.Printf("unknown address type: %T\n", a)
		return nil
	}
}

func IgnoreConnectionClosed(err error) error {
	var (
		c = &net.OpError{}
	)

	if !errors.As(err, &c) {
		return err
	}

	if !strings.HasSuffix(c.Error(), "use of closed network connection") {
		return err
	}

	return nil
}
