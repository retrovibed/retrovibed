package ddisc

import (
	"io"

	"github.com/gofrs/uuid/v5"
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
	parts []rendezvous.Paired[uuid.UUID]
}

func (t Partition) Max(k []byte) uuid.UUID {
	return rendezvous.Max(
		k,
		func(v rendezvous.Paired[uuid.UUID]) []byte { return v.N.Bytes() },
		t.parts...,
	).N
}

func PartitionsDigest(p *Partition) uuid.UUID {
	return uuid.FromBytesOrNil(md5x.Digest(slicesx.MapTransform(func(v rendezvous.Paired[uuid.UUID]) []byte { return v.N[:] })...).Sum(nil))
}
