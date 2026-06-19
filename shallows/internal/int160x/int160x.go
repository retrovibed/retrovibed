package int160x

import (
	"io"
	"net/netip"

	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/shallows/internal/cryptox"
	"github.com/retrovibed/retrovibed/shallows/internal/rendezvous"
)

type Ranger interface {
	Generate() int160.T
}

func NewRangeRandom() *RangerRandom {
	return &RangerRandom{}
}

type RangerRandom struct{}

func (t RangerRandom) Generate() int160.T {
	return int160.Random()
}

func NewRangeDynamic(s *dht.Server, n uint16) *RangerDynamic {
	stable := int160.StableSuffix(s.ID(s.DynamicAddrPort()))
	src := cryptox.NewChaCha8(stable.String())
	return &RangerDynamic{
		stable: stable,
		ranges: rendezvous.SeededAddrN(n, src),
		src:    src,
	}
}

type RangerDynamic struct {
	stable int160.T
	ranges []rendezvous.Paired[netip.Addr]
	src    io.Reader
}

func (t RangerDynamic) Generate() int160.T {
	var (
		buf [16]byte
	)

	_, _ = io.ReadFull(t.src, buf[:]) // ignore the error not the end of the world if we get one every now and then here.

	sampleip := rendezvous.Max(buf[:], func(n rendezvous.Paired[netip.Addr]) []byte { return n.Bi.Bytes() }, t.ranges...)
	target := t.stable.Secure(sampleip.N)

	return target
}

func NewRangeFixed(v int160.T) *RangerFixed {
	return &RangerFixed{v: v}
}

type RangerFixed struct {
	v int160.T
}

func (t RangerFixed) Generate() int160.T {
	return t.v
}
