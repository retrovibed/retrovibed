package netipx

import (
	"net/netip"

	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

func AddrPortFromStrings(in ...string) (res []netip.AddrPort) {
	for _, ap := range in {
		v := errorsx.Zero(netip.ParseAddrPort(ap))
		if !v.IsValid() {
			continue
		}

		res = append(res, v)
	}
	return res
}

func AddrFromSlice(ip []byte) netip.Addr {
	addr, _ := netip.AddrFromSlice(ip)
	return addr
}

func IPv4Loopback() netip.Addr {
	return netip.AddrFrom4([4]byte{127, 0, 0, 1})
}
