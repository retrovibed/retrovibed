package ddisc_test

import (
	"testing"

	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/stretchr/testify/require"
)

func TestDiscoveredOptionAutoMagnet(t *testing.T) {
	t.Run("synthesizes a magnet uri from the infohash when uri is empty", func(t *testing.T) {
		infohash := []byte{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
			11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
		}
		d := ddisc.Discovered{Infohash: infohash}
		ddisc.DiscoveredOptionAutoMagnet(&d)

		require.Equal(t, metainfo.NewMagnetFromInfohash(infohash).String(), d.URI)
		require.Equal(t, mimex.Bittorrent, d.Contentmime)
	})

	t.Run("leaves an already set uri untouched", func(t *testing.T) {
		d := ddisc.Discovered{URI: "https://tracker.example/download/1.torrent"}
		ddisc.DiscoveredOptionAutoMagnet(&d)

		require.Equal(t, "https://tracker.example/download/1.torrent", d.URI)
		require.Empty(t, d.Contentmime)
	})
}
