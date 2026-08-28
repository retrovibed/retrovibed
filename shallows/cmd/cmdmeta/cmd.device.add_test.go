package cmdmeta_test

import (
	"bytes"
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdmeta"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdtestx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestDeviceAdd(t *testing.T) {
	t.Run("adds a reachable device", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		routes := mux.NewRouter()
		routes.Path("/healthz").Methods(http.MethodGet).HandlerFunc(httpx.Healthz("test", 0, http.StatusOK))
		metaapi.NewHTTPDaemons(
			q,
			metaapi.HTTPDaemonsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/meta/d").Subrouter())
		srv := cmdtestx.NewTLSServer(t, q, routes)
		hostname := strings.TrimPrefix(srv.URL, "https://")

		var buf bytes.Buffer
		genparser := cmdtestx.Genparser(cmdmeta.Device{}, kong.Writers(&buf, nil))
		require.NoError(t, cmdtestx.Execute(t, genparser(t), "command", "add", hostname, "--description", "self", "--private-key-path", keypath, "--endpoint", srv.URL))

		var got meta.Daemon
		require.NoError(t, meta.DaemonFindByLatestUpdated(ctx, q).Scan(&got))
		require.Equal(t, hostname, got.Hostname)
		require.Equal(t, "self", got.Description)
	})

	t.Run("fails without --force when unreachable", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		routes := mux.NewRouter()
		metaapi.NewHTTPDaemons(
			q,
			metaapi.HTTPDaemonsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/meta/d").Subrouter())
		srv := cmdtestx.NewTLSServer(t, q, routes)

		var buf bytes.Buffer
		genparser := cmdtestx.Genparser(cmdmeta.Device{}, kong.Writers(&buf, nil))
		err := cmdtestx.Execute(t, genparser(t), "command", "add", "unreachable.invalid:9999", "--private-key-path", keypath, "--endpoint", srv.URL)
		require.Error(t, err)
		require.EqualValues(t, 0, sqltestx.Count(t, q, "SELECT COUNT(*) FROM meta_daemons"))
	})

	t.Run("adds anyway with --force when unreachable", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		routes := mux.NewRouter()
		metaapi.NewHTTPDaemons(
			q,
			metaapi.HTTPDaemonsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/meta/d").Subrouter())
		srv := cmdtestx.NewTLSServer(t, q, routes)

		var buf bytes.Buffer
		genparser := cmdtestx.Genparser(cmdmeta.Device{}, kong.Writers(&buf, nil))
		require.NoError(t, cmdtestx.Execute(t, genparser(t), "command", "add", "unreachable.invalid:9999", "--force", "--private-key-path", keypath, "--endpoint", srv.URL))

		require.EqualValues(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM meta_daemons"))
	})

	t.Run("round-trips default and download flags", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		routes := mux.NewRouter()
		routes.Path("/healthz").Methods(http.MethodGet).HandlerFunc(httpx.Healthz("test", 0, http.StatusOK))
		metaapi.NewHTTPDaemons(
			q,
			metaapi.HTTPDaemonsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/meta/d").Subrouter())
		srv := cmdtestx.NewTLSServer(t, q, routes)
		hostname := strings.TrimPrefix(srv.URL, "https://")

		var buf bytes.Buffer
		genparser := cmdtestx.Genparser(cmdmeta.Device{}, kong.Writers(&buf, nil))
		require.NoError(t, cmdtestx.Execute(t, genparser(t), "command", "add", hostname, "--default", "--download", "--private-key-path", keypath, "--endpoint", srv.URL))

		var got meta.Daemon
		require.NoError(t, meta.DaemonFindDefault(ctx, q).Scan(&got))
		require.Equal(t, hostname, got.Hostname)
		require.True(t, got.Downloads)
	})
}
