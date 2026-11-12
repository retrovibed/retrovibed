package ddisc_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/uuidx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/stretchr/testify/require"
)

func TestDiscoveredOptionIndex(t *testing.T) {
	t.Run("true sets known media id to the max uuid", func(t *testing.T) {
		d := ddisc.Discovered{KnownMediaID: uuidx.Nil().String()}
		ddisc.DiscoveredOptionIndex(true)(&d)
		require.Equal(t, uuidx.Max().String(), d.KnownMediaID)
	})

	t.Run("false leaves known media id unchanged", func(t *testing.T) {
		d := ddisc.Discovered{KnownMediaID: "some-id"}
		ddisc.DiscoveredOptionIndex(false)(&d)
		require.Equal(t, "some-id", d.KnownMediaID)
	})
}
