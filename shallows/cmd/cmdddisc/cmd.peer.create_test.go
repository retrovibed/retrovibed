package cmdddisc_test

import (
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/gorilla/mux"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdddisc"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdtestx"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func bindPeerManagement(t *testing.T, q *sql.DB) *mux.Router {
	t.Helper()

	routes := mux.NewRouter()
	ddiscapi.NewHTTPPeerManagement(
		q,
		ddiscapi.HTTPPeerManagementOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(routes.PathPrefix("/ddisc").Subrouter())

	return routes
}

func TestPeerCreate(t *testing.T) {
	t.Run("creates a peer", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		routes := bindPeerManagement(t, q)
		srv := cmdtestx.NewTLSServer(t, q, routes)

		require.NoError(t, cmdtestx.Execute(t, cmdtestx.Genparser(cmdddisc.Commands{})(t), "command", "peers", "create",
			"--private-key-path", keypath,
			"--endpoint", srv.URL,
			"--name", "derp",
			"--peer", "34363564353033612d643263352d363338332d30",
			"--partition", "033292b1-98c2-5e96-38a4-956548a40b55",
		))

		id, err := int160.FromHexEncodedString("34363564353033612d643263352d363338332d30")
		require.NoError(t, err)

		var stored tracking.Peer
		require.NoError(t, tracking.PeerFindByID(ctx, q, tracking.PeerUID(id)).Scan(&stored))
		require.Equal(t, "derp", stored.Description)
		require.Equal(t, "033292b1-98c2-5e96-38a4-956548a40b55", stored.DdiscPartition)
	})
}
