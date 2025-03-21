package netipx

import "net/netip"

func AddrFromSlice(ip []byte) netip.Addr {
	addr, _ := netip.AddrFromSlice(ip)
	return addr
}
