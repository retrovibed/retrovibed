//go:build !linux

package audiox

import (
	"context"
	"errors"
	"iter"

	"github.com/retrovibed/retrovibed/retroapi/iterx"
)

// Supported reports whether this build can talk to PulseAudio.
const Supported = false

type unsupportedSinkSeq struct{}

func listSinks() iterx.Seq[Sink] {
	return unsupportedSinkSeq{}
}

func (unsupportedSinkSeq) Each(_ context.Context) iter.Seq[Sink] {
	return func(yield func(Sink) bool) {}
}

func (unsupportedSinkSeq) Err() error {
	return errors.ErrUnsupported
}

func current() (Sink, error) {
	return Sink{}, errors.ErrUnsupported
}

func activate(id string) error {
	return errors.ErrUnsupported
}
