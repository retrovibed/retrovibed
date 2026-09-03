package ddisc_test

import (
	"testing"

	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/stretchr/testify/require"
)

func TestDiscoveredOptionFromTorrentInfo(t *testing.T) {
	t.Run("copies name, length, and private", func(t *testing.T) {
		priv := true
		info := &metainfo.Info{Name: "the.title", Length: 1024, Private: &priv}

		d := ddisc.Discovered{}
		ddisc.DiscoveredOptionFromTorrentInfo(info)(&d)

		require.Equal(t, "the.title", d.Title)
		require.EqualValues(t, 1024, d.Bytes)
		require.True(t, d.Private)
	})

	t.Run("defaults private to false when unset", func(t *testing.T) {
		info := &metainfo.Info{Name: "the.title", Length: 1024}

		d := ddisc.Discovered{}
		ddisc.DiscoveredOptionFromTorrentInfo(info)(&d)

		require.False(t, d.Private)
	})
}
