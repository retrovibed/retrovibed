package krpc

import (
	"net/netip"

	"github.com/anacrolix/missinggo/slices"
)

type CompactIPv4NodeAddrs []NodeAddr

func (CompactIPv4NodeAddrs) ElemSize() int { return 6 }

func (me CompactIPv4NodeAddrs) MarshalBinary() ([]byte, error) {
	return marshalBinarySlice(slices.Map(func(addr NodeAddr) NodeAddr {
		addr.AddrPort = netip.AddrPortFrom(netip.AddrFrom4(addr.Addr().As4()), addr.Port())
		return addr
	}, me).(CompactIPv4NodeAddrs))
}

func (me CompactIPv4NodeAddrs) MarshalBencode() ([]byte, error) {
	return bencodeBytesResult(me.MarshalBinary())
}

func (me *CompactIPv4NodeAddrs) UnmarshalBinary(b []byte) error {
	return unmarshalBinarySlice(me, b)
}

func (me *CompactIPv4NodeAddrs) UnmarshalBencode(b []byte) error {
	return unmarshalBencodedBinary(me, b)
}

func (me CompactIPv4NodeAddrs) NodeAddrs() []NodeAddr {
	return me
}

func (me CompactIPv4NodeAddrs) Index(x NodeAddr) int {
	return addrIndex(me, x)
}
