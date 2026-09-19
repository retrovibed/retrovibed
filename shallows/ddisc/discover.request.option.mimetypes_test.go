package ddisc_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/stretchr/testify/require"
)

func TestDiscoverRequestOptionMimetypes(t *testing.T) {
	t.Run("sets mimetypes", func(t *testing.T) {
		r := ddisc.DiscoverRequest{}
		ddisc.DiscoverRequestOptionMimetypes(mimex.RetrovibedDiscoveryMovies, mimex.RetrovibedDiscoveryTV)(&r)
		require.Equal(t, []string{mimex.RetrovibedDiscoveryMovies, mimex.RetrovibedDiscoveryTV}, r.Mimetypes)
	})

	t.Run("replaces existing mimetypes", func(t *testing.T) {
		r := ddisc.DiscoverRequest{Mimetypes: []string{mimex.RetrovibedDiscoveryMusic}}
		ddisc.DiscoverRequestOptionMimetypes(mimex.RetrovibedDiscoveryMovies)(&r)
		require.Equal(t, []string{mimex.RetrovibedDiscoveryMovies}, r.Mimetypes)
	})
}
