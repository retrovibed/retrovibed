package cmdddisc_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdtestx"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestDiscoveryLs(t *testing.T) {
	t.Run("lists all entries", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		var uh tracking.UnknownHash
		require.NoError(t, testx.Fake(&uh, tracking.UnknownHashOptionTestDefaults))
		require.NoError(t, tracking.UnknownHashInsertWithDefaults(ctx, q, uh).Scan(&uh))

		routes := mux.NewRouter()
		ddiscapi.NewHTTPDiscovery(
			q,
			ddiscapi.HTTPDiscoveryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/ddisc/discovery").Subrouter())
		srv := cmdtestx.NewTLSServer(t, q, routes)

		require.NoError(t, cmdtestx.Execute(t, genparser(t), "discovery", "ls",
			"--private-key-path", keypath,
			"--insecure",
			"--library", srv.Listener.Addr().String(),
		))
	})

	t.Run("filters to entries needing a check", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		var due tracking.UnknownHash
		require.NoError(t, testx.Fake(&due, tracking.UnknownHashOptionTestDefaults))
		due.NextCheck = time.Now().Add(-time.Minute)
		require.NoError(t, tracking.UnknownHashInsertWithDefaults(ctx, q, due).Scan(&due))

		var future tracking.UnknownHash
		require.NoError(t, testx.Fake(&future, tracking.UnknownHashOptionTestDefaults))
		future.NextCheck = time.Now().Add(time.Hour)
		require.NoError(t, tracking.UnknownHashInsertWithDefaults(ctx, q, future).Scan(&future))

		routes := mux.NewRouter()
		ddiscapi.NewHTTPDiscovery(
			q,
			ddiscapi.HTTPDiscoveryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/ddisc/discovery").Subrouter())
		srv := cmdtestx.NewTLSServer(t, q, routes)

		require.NoError(t, cmdtestx.Execute(t, genparser(t), "discovery", "ls",
			"--private-key-path", keypath,
			"--insecure",
			"--library", srv.Listener.Addr().String(),
			"--needs-check",
		))
	})

	t.Run("filters to a specific id", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		var wanted tracking.UnknownHash
		require.NoError(t, testx.Fake(&wanted, tracking.UnknownHashOptionTestDefaults))
		require.NoError(t, tracking.UnknownHashInsertWithDefaults(ctx, q, wanted).Scan(&wanted))

		var other tracking.UnknownHash
		require.NoError(t, testx.Fake(&other, tracking.UnknownHashOptionTestDefaults))
		require.NoError(t, tracking.UnknownHashInsertWithDefaults(ctx, q, other).Scan(&other))

		routes := mux.NewRouter()
		ddiscapi.NewHTTPDiscovery(
			q,
			ddiscapi.HTTPDiscoveryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/ddisc/discovery").Subrouter())
		srv := cmdtestx.NewTLSServer(t, q, routes)

		require.NoError(t, cmdtestx.Execute(t, genparser(t), "discovery", "ls",
			"--private-key-path", keypath,
			"--insecure",
			"--library", srv.Listener.Addr().String(),
			"--id", wanted.ID,
		))
	})
}
