package ddisc

import (
	"io"

	"github.com/retrovibed/retrovibed/retroapi/uuidx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/rendezvous"
	"github.com/retrovibed/retrovibed/shallows/internal/slicesx"
)

func Partitions(n uint16, s io.Reader) *Partition {
	return &Partition{
		parts: rendezvous.SeededN(n, s),
	}
}

type Partition struct {
	parts []rendezvous.Paired[uuidx.UUID]
}

func (t Partition) Max(k []byte) uuidx.UUID {
	return rendezvous.Max(
		k,
		func(v rendezvous.Paired[uuidx.UUID]) []byte { return v.N.Bytes() },
		t.parts...,
	).N
}

func PartitionsDigest(p *Partition) uuidx.UUID {
	return uuidx.FromBytesOrNil(md5x.Digest(slicesx.MapTransform(func(v rendezvous.Paired[uuidx.UUID]) []byte { return v.N[:] })...).Sum(nil))
}
