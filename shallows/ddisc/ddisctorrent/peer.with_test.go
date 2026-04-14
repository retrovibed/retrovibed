package ddisctorrent_test

import (
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/dht/krpc"
	"github.com/retrovibed/retrovibed/shallows/ddisc/ddisctorrent"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/netx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/torrenttestx"
	"github.com/stretchr/testify/require"
)

func TestPeerWith(t *testing.T) {
	t.Run("should receive results of the meta protocol", func(t *testing.T) {
		pdb := sqltestx.Metadatabase(t)
		defer pdb.Close()

		ptm := dht.DefaultMuxer().
			Method(ddisctorrent.MethodMeta, ddisctorrent.NewMeta(uuid.Max))
		pdht, err := dht.NewServer(
			32,
			dht.OptionMuxer(ptm),
		)
		require.NoError(t, err)
		tpeer := torrenttestx.QuickClientWithDHT(
			t,
			pdht,
			torrent.ClientConfigExtension(ddisctorrent.ExtensionName),
		)
		defer tpeer.Close()

		cdb := sqltestx.Metadatabase(t)
		defer cdb.Close()

		ctm := dht.DefaultMuxer().
			Method(ddisctorrent.MethodMeta, ddisctorrent.NewMeta(uuid.Max))
		cdht, err := dht.NewServer(
			32,
			dht.OptionMuxer(ctm),
		)
		require.NoError(t, err)
		tclient := torrenttestx.QuickClientWithDHT(
			t,
			cdht,
			torrent.ClientConfigExtension(ddisctorrent.ExtensionName),
		)
		defer tclient.Close()

		info, err := ddisctorrent.PeerWith(
			t.Context(),
			cdht,
			uuid.Nil,
			krpc.NewInfo(
				pdht.ID().AsByteArray(),
				krpc.NewNodeAddrFromAddrPort(
					langx.Autoderef(netx.AddrPort(pdht.Addr())),
				),
			),
		)
		require.NoError(t, err)
		require.Equal(t, uuid.Max.String(), info.Partition)
		require.Equal(t, pdht.ID(), info.Peer.Int160())
		require.Equal(t, ddisctorrent.ExtensionName, info.Version)
	})

	t.Run("peer doesnt support ddisc meta protocol", func(t *testing.T) {
		pdb := sqltestx.Metadatabase(t)
		defer pdb.Close()

		ptm := dht.DefaultMuxer()
		pdht, err := dht.NewServer(
			32,
			dht.OptionMuxer(ptm),
		)
		require.NoError(t, err)
		tpeer := torrenttestx.QuickClientWithDHT(
			t,
			pdht,
			torrent.ClientConfigExtension(ddisctorrent.ExtensionName),
		)
		defer tpeer.Close()

		cdb := sqltestx.Metadatabase(t)
		defer cdb.Close()

		ctm := dht.DefaultMuxer().
			Method(ddisctorrent.MethodMeta, ddisctorrent.NewMeta(uuid.Max))
		cdht, err := dht.NewServer(
			32,
			dht.OptionMuxer(ctm),
		)
		require.NoError(t, err)
		tclient := torrenttestx.QuickClientWithDHT(
			t,
			cdht,
			torrent.ClientConfigExtension(ddisctorrent.ExtensionName),
		)
		defer tclient.Close()

		info, err := ddisctorrent.PeerWith(
			t.Context(),
			cdht,
			uuid.Nil,
			krpc.NewInfo(
				pdht.ID().AsByteArray(),
				krpc.NewNodeAddrFromAddrPort(
					langx.Autoderef(netx.AddrPort(pdht.Addr())),
				),
			),
		)
		var (
			expected = new(krpc.Error)
		)
		require.ErrorAs(t, err, &expected)
		require.Equal(t, krpc.ErrorCodeMethodUnknown, expected.Code)
		require.Equal(t, "", info.Partition)
	})
}
