package cmdmeta_test

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdmeta"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdtestx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestDeviceEdit(t *testing.T) {
	t.Run("marks a device as default", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		var first, second meta.Daemon
		require.NoError(t, testx.Fake(&first, meta.DaemonOptionTestDefaults, meta.DaemonOptionMaybeID, func(d *meta.Daemon) { d.Default = true }))
		require.NoError(t, meta.DaemonInsertWithDefaults(ctx, q, first).Scan(&first))
		require.NoError(t, testx.Fake(&second, meta.DaemonOptionTestDefaults, meta.DaemonOptionMaybeID))
		require.NoError(t, meta.DaemonInsertWithDefaults(ctx, q, second).Scan(&second))

		routes := mux.NewRouter()
		metaapi.NewHTTPDaemons(
			q,
			metaapi.HTTPDaemonsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/meta/d").Subrouter())
		srv := cmdtestx.NewTLSServer(t, q, routes)

		var buf bytes.Buffer
		genparser := cmdtestx.Genparser(cmdmeta.Device{}, kong.Writers(&buf, nil))
		require.NoError(t, cmdtestx.Execute(t, genparser(t), "command", "edit", second.ID, "--default", "--private-key-path", keypath, "--endpoint", srv.URL))

		var got meta.Daemon
		require.NoError(t, meta.DaemonFindDefault(ctx, q).Scan(&got))
		require.Equal(t, second.ID, got.ID)
		require.EqualValues(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM meta_daemons WHERE \"default\""))
	})

	t.Run("marks a device as the download target", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		var d meta.Daemon
		require.NoError(t, testx.Fake(&d, meta.DaemonOptionTestDefaults, meta.DaemonOptionMaybeID))
		require.NoError(t, meta.DaemonInsertWithDefaults(ctx, q, d).Scan(&d))

		routes := mux.NewRouter()
		metaapi.NewHTTPDaemons(
			q,
			metaapi.HTTPDaemonsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/meta/d").Subrouter())
		srv := cmdtestx.NewTLSServer(t, q, routes)

		var buf bytes.Buffer
		genparser := cmdtestx.Genparser(cmdmeta.Device{}, kong.Writers(&buf, nil))
		require.NoError(t, cmdtestx.Execute(t, genparser(t), "command", "edit", d.ID, "--download", "--private-key-path", keypath, "--endpoint", srv.URL))

		var got meta.Daemon
		require.NoError(t, meta.DaemonFindByDownload(ctx, q).Scan(&got))
		require.Equal(t, d.ID, got.ID)
		require.Equal(t, d.Hostname, got.Hostname)
		require.Equal(t, d.Description, got.Description)
	})

	t.Run("updates description while preserving other fields", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		var d meta.Daemon
		require.NoError(t, testx.Fake(&d, meta.DaemonOptionTestDefaults, meta.DaemonOptionMaybeID, meta.DaemonOptionDownload))
		require.NoError(t, meta.DaemonInsertWithDefaults(ctx, q, d).Scan(&d))

		routes := mux.NewRouter()
		metaapi.NewHTTPDaemons(
			q,
			metaapi.HTTPDaemonsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/meta/d").Subrouter())
		srv := cmdtestx.NewTLSServer(t, q, routes)

		var buf bytes.Buffer
		genparser := cmdtestx.Genparser(cmdmeta.Device{}, kong.Writers(&buf, nil))
		require.NoError(t, cmdtestx.Execute(t, genparser(t), "command", "edit", d.ID, "--description", "updated", "--private-key-path", keypath, "--endpoint", srv.URL))

		var got meta.Daemon
		require.NoError(t, meta.DaemonFindByDownload(ctx, q).Scan(&got))
		require.Equal(t, d.ID, got.ID)
		require.Equal(t, d.Hostname, got.Hostname)
		require.Equal(t, "updated", got.Description)
		require.True(t, got.Downloads)
	})

	t.Run("applies default and download together", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		var d meta.Daemon
		require.NoError(t, testx.Fake(&d, meta.DaemonOptionTestDefaults, meta.DaemonOptionMaybeID))
		require.NoError(t, meta.DaemonInsertWithDefaults(ctx, q, d).Scan(&d))

		routes := mux.NewRouter()
		metaapi.NewHTTPDaemons(
			q,
			metaapi.HTTPDaemonsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/meta/d").Subrouter())
		srv := cmdtestx.NewTLSServer(t, q, routes)

		var buf bytes.Buffer
		genparser := cmdtestx.Genparser(cmdmeta.Device{}, kong.Writers(&buf, nil))
		require.NoError(t, cmdtestx.Execute(t, genparser(t), "command", "edit", d.ID, "--default", "--download", "--private-key-path", keypath, "--endpoint", srv.URL))

		var gotDefault, gotDownload meta.Daemon
		require.NoError(t, meta.DaemonFindDefault(ctx, q).Scan(&gotDefault))
		require.Equal(t, d.ID, gotDefault.ID)
		require.NoError(t, meta.DaemonFindByDownload(ctx, q).Scan(&gotDownload))
		require.Equal(t, d.ID, gotDownload.ID)
	})
}
