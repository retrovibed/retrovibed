package cmdddisc_test

import (
	"database/sql"
	"encoding/hex"
	"path/filepath"
	"testing"

	"github.com/james-lawrence/torrent/dht/krpc"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdddisc"
	cmdtestx "github.com/retrovibed/retrovibed/shallows/cmd/cmdtestx/v2"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestPeerDelete(t *testing.T) {
	t.Run("deletes a peer", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		peer := tracking.NewPeerFromInfo(krpc.RandomNodeInfo(16), tracking.PeerOptionTestDefaults)
		require.NoError(t, tracking.PeerInsertWithDefaults(ctx, q, peer).Scan(&peer))

		routes := bindPeerManagement(t, q)
		srv := cmdtestx.NewTLSServer(t, q, routes)

		require.NoError(t, cmdtestx.New(t, &cmdddisc.Commands{}).Options(cmdtestx.Standard(t)).Run("peers", "delete",
			"--private-key-path", keypath,
			"--endpoint", srv.URL,
			"--peer", hex.EncodeToString(peer.Peer)))

		var target = sql.ErrNoRows
		require.ErrorAs(t, tracking.PeerFindByID(ctx, q, peer.ID).Scan(&peer), &target)
	})
}
