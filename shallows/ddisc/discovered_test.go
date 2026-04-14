package ddisc_test

import (
	"testing"

	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/torrenttest"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/stretchr/testify/require"
)

func TestDiscovered(t *testing.T) {
	t.Run("InsertWithDefaults", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		tmpdir := t.TempDir()

		q := sqltestx.Metadatabase(t)

		id := int160.Random()
		info, _, err := torrenttest.Random(tmpdir, 128*bytesx.KiB)
		require.NoError(t, err)

		d := ddisc.NewDiscovered(
			&id,
			ddisc.DiscoveredOptionIndex(true),
			ddisc.DiscoveredOptionMimetype(mimex.Binary),
			ddisc.DiscoveredOptionFromTorrentInfo(info),
		)
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))
	})
}
