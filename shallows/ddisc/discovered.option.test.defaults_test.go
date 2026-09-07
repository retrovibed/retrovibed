package ddisc_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/userx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/localex"
	"github.com/stretchr/testify/require"
)

func TestDiscoveredOptionTestDefaults(t *testing.T) {
	d := ddisc.Discovered{}
	ddisc.DiscoveredOptionTestDefaults(&d)

	want := localex.FirstDefined(userx.LocaleLanguage())
	require.Equal(t, want, d.AudioDefaultLocale)
	require.Equal(t, want, d.SubtitlesDefaultLocale)
}
