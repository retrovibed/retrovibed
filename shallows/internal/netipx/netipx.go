package netipx

import "net/netip"

func AddrFromSlice(ip []byte) netip.Addr {
	addr, _ := netip.AddrFromSlice(ip)
	return addr
}

func IPv4Loopback() netip.Addr {
	return netip.AddrFrom4([4]byte{127, 0, 0, 1})
}
