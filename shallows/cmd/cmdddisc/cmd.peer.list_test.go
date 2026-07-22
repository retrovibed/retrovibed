package cmdddisc_test

import (
	"path/filepath"
	"testing"

	"github.com/james-lawrence/torrent/dht/krpc"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdddisc"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdtestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestPeerList(t *testing.T) {
	t.Run("lists all peers", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		peer := tracking.NewPeerFromInfo(krpc.RandomNodeInfo(16), tracking.PeerOptionTestDefaults)
		require.NoError(t, tracking.PeerInsertWithDefaults(ctx, q, peer).Scan(&peer))

		routes := bindPeerManagement(t, q)
		srv := cmdtestx.NewTLSServer(t, q, routes)

		require.NoError(t, cmdtestx.Execute(t, cmdtestx.Genparser(cmdddisc.Commands{})(t), "command", "peers", "list",
			"--private-key-path", keypath,
			"--insecure",
			"--endpoint", srv.URL,
		))
	})

	t.Run("filters by query", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		peer := tracking.NewPeerFromInfo(krpc.RandomNodeInfo(16), tracking.PeerOptionTestDefaults, tracking.PeerOptionDescription("searchable peer"))
		require.NoError(t, tracking.PeerInsertWithDefaults(ctx, q, peer).Scan(&peer))

		routes := bindPeerManagement(t, q)
		srv := cmdtestx.NewTLSServer(t, q, routes)

		require.NoError(t, cmdtestx.Execute(t, cmdtestx.Genparser(cmdddisc.Commands{})(t), "command", "peers", "list",
			"--private-key-path", keypath,
			"--insecure",
			"--endpoint", srv.URL,
			"--query", "searchable",
		))
	})
}
