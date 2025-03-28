package krpc

import "github.com/anacrolix/missinggo/slices"

type (
	CompactIPv4NodeInfo []NodeInfo
)

func (CompactIPv4NodeInfo) ElemSize() int {
	return 26
}

func (me CompactIPv4NodeInfo) MarshalBinary() ([]byte, error) {
	return marshalBinarySlice(slices.Map(func(ni NodeInfo) NodeInfo {
		ni.Addr = NewNodeAddrFromIPPort(ni.Addr.IP().To4(), int(ni.Addr.Port()))
		return ni
	}, me).(CompactIPv4NodeInfo))
}

func (me CompactIPv4NodeInfo) MarshalBencode() ([]byte, error) {
	return bencodeBytesResult(me.MarshalBinary())
}

func (me *CompactIPv4NodeInfo) UnmarshalBinary(b []byte) error {
	return unmarshalBinarySlice(me, b)
}

func (me *CompactIPv4NodeInfo) UnmarshalBencode(b []byte) error {
	return unmarshalBencodedBinary(me, b)
}
