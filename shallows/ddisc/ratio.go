package ddisc

import (
	"io"

	"github.com/retrovibed/retrovibed/internal/rendezvous"
)

type Filter func(k []byte) bool

// attempt to index everything.
func FilterNone(k []byte) bool { return false }

// only accept a percentage of the infohashes based on the provided key.
// ratio is an integer 0-100.
func FilterRatio(src io.Reader, ratio uint8) rendezvousf {
	nodes := make([]rendezvous.Paired[bool], 0, 100)
	for range ratio {
		nodes = append(nodes, rendezvous.Seeded(src, false))
	}

	for range 100 - ratio {
		nodes = append(nodes, rendezvous.Seeded(src, true))
	}

	return rendezvousf(nodes)
}

type rendezvousf []rendezvous.Paired[bool]

func (t rendezvousf) Filter(k []byte) bool {
	return rendezvous.Max(k, func(b rendezvous.Paired[bool]) []byte {
		return b.Bi.Bytes()
	}, ([]rendezvous.Paired[bool])(t)...).N
}
