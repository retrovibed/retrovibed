package cmdddisc_test

import (
	"context"
	"database/sql"
	"path/filepath"
	"sync"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/gorilla/mux"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdddisc"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdtestx"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func genparser(t *testing.T, options ...kong.Option) *kong.Kong {
	var cli struct {
		cmdopts.Global
		cmdopts.TLSConfig
		cmdopts.SSHID
		cmdddisc.Commands
	}
	cli.Context, cli.Shutdown = context.WithCancel(context.Background())
	cli.Cleanup = &sync.WaitGroup{}

	return kong.Must(
		&cli,
		append(options,
			kong.Bind(&cli.TLSConfig),
			kong.Bind(&cli.Global),
			kong.Bind(&cli.SSHID),
			kong.Vars{
				"vars_private_key":                  env.PrivateKeyPath(),
				"vars_user_configuration_directory": t.TempDir(),
			},
			kong.NamedMapper("durationinf", kong.MapperFunc(cmdopts.ParseDurationInf)),
			kong.NamedMapper("envvar", kong.MapperFunc(cmdopts.ParseEnviron)),
		)...,
	)
}

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

		require.NoError(t, cmdtestx.Execute(t, genparser(t), "peers", "create",
			"--private-key-path", keypath,
			"--insecure",
			"--library", srv.Listener.Addr().String(),
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
