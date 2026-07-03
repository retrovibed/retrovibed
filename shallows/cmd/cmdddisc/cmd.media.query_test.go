package cmdddisc_test

import (
	"bytes"
	"net/netip"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/autobind"
	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/ddisc/ddisctorrent"
	"github.com/retrovibed/retrovibed/shallows/internal/slicesx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/torrenttestx"
	"github.com/stretchr/testify/require"
)

func TestMediaQuery(t *testing.T) {
	t.Run("finds media known to a peer", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		id := uuid.Must(uuid.NewV7()).String()

		pdb := sqltestx.Metadatabase(t)
		defer pdb.Close()

		rid := int160.Random()
		record := ddisc.NewDiscovered(&rid, ddisc.DiscoveredOptionKnownMedia(id))
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, pdb, record).Scan(&record))

		ptm := dht.DefaultMuxer().
			Method(ddisctorrent.MethodSearch, ddisctorrent.NewSearch(pdb))
		pdht, err := dht.NewServer(
			32,
			dht.OptionMuxer(ptm),
		)
		require.NoError(t, err)

		tpeer := torrenttestx.QuickClientBinder(
			t,
			autobind.New(autobind.EnableDHT(pdht)),
			torrent.ClientConfigExtension(ddisctorrent.ExtensionName),
		)
		defer tpeer.Close()

		var buf bytes.Buffer

		peers := slicesx.MapTransform(func(n netip.AddrPort) string { return n.String() }, torrenttestx.ApprPorts(tpeer)...)

		kctx, err := genparser(t, kong.Writers(nil, &buf)).Parse(
			append([]string{
				"media",
				"query",
				"--no-bootstrap",
				"--timeout", "2s",
				"--wait", "3s",
				id,
			},
				peers...,
			),
		)

		require.NoError(t, err)
		require.NoError(t, kctx.Run())

		require.Contains(t, buf.String(), id)
	})
}
