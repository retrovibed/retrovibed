//go:build linux

package audiox

import (
	"context"
	"iter"

	"github.com/jfreymuth/pulse"
	"github.com/jfreymuth/pulse/proto"
	"github.com/retrovibed/retrovibed/retroapi/iterx"
)

// Supported reports whether this build can talk to PulseAudio.
const Supported = true

type pulseSinkSeq struct {
	err error
}

func listSinks() iterx.Seq[Sink] {
	return &pulseSinkSeq{}
}

func (s *pulseSinkSeq) Each(_ context.Context) iter.Seq[Sink] {
	return func(yield func(Sink) bool) {
		c, err := pulse.NewClient()
		if err != nil {
			s.err = err
			return
		}
		defer c.Close()

		sinks, err := c.ListSinks()
		if err != nil {
			s.err = err
			return
		}

		for _, sk := range sinks {
			if !yield(Sink{ID: sk.ID(), Name: sk.Name()}) {
				return
			}
		}
	}
}

func (s *pulseSinkSeq) Err() error {
	return s.err
}

func current() (Sink, error) {
	c, err := pulse.NewClient()
	if err != nil {
		return Sink{}, err
	}
	defer c.Close()

	s, err := c.DefaultSink()
	if err != nil {
		return Sink{}, err
	}

	return Sink{ID: s.ID(), Name: s.Name()}, nil
}

func activate(id string) error {
	c, err := pulse.NewClient()
	if err != nil {
		return err
	}
	defer c.Close()

	return c.RawRequest(&proto.SetDefaultSink{SinkName: id}, nil)
}
