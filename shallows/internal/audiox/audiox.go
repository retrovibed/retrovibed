// Package audiox enumerates and selects the active system audio output
// device (PulseAudio "sink" on linux). Platform-specific behavior lives in
// audiox_linux.go / audiox_other.go, split by build tag.
package audiox

import (
	"github.com/retrovibed/retrovibed/retroapi/iterx"
)

// Sink is an audio output device.
type Sink struct {
	ID   string
	Name string
}

// ListSinks returns all available output devices.
func ListSinks() iterx.Seq[Sink] {
	return listSinks()
}

// Current returns the active output device.
func Current() (Sink, error) {
	return current()
}

// Activate makes the sink identified by id the active output device.
func Activate(id string) error {
	return activate(id)
}
