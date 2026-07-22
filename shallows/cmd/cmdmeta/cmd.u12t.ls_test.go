package cmdmeta_test

import (
	"path/filepath"
	"testing"

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

func TestLs(t *testing.T) {
	genparser := cmdtestx.Genparser(cmdmeta.Usermanagement{})

	t.Run("lists all profiles when no filters specified", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		var p meta.Profile
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))

		routes := mux.NewRouter()
		metaapi.NewHTTPUsermanagement(
			q,
			metaapi.HTTPUsermanagementOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/meta/u12t").Subrouter())
		srv := cmdtestx.NewTLSServer(t, q, routes)

		require.NoError(t, cmdtestx.Execute(t, genparser(t), "command", "ls", "--private-key-path", keypath, "--endpoint", srv.URL))
	})

	t.Run("filters to pending profiles", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		var p meta.Profile
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))

		routes := mux.NewRouter()
		metaapi.NewHTTPUsermanagement(
			q,
			metaapi.HTTPUsermanagementOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/meta/u12t").Subrouter())
		srv := cmdtestx.NewTLSServer(t, q, routes)

		require.NoError(t, cmdtestx.Execute(t, genparser(t), "command", "ls", "--private-key-path", keypath, "--pending", "--endpoint", srv.URL))
	})

	t.Run("filters to enabled profiles", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		var p meta.Profile
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, meta.ProfileEnable(ctx, q, p.ID).Scan(&p))

		routes := mux.NewRouter()
		metaapi.NewHTTPUsermanagement(
			q,
			metaapi.HTTPUsermanagementOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/meta/u12t").Subrouter())
		srv := cmdtestx.NewTLSServer(t, q, routes)

		require.NoError(t, cmdtestx.Execute(t, genparser(t), "command", "ls", "--private-key-path", keypath, "--enabled", "--endpoint", srv.URL))
	})

	t.Run("filters to disabled profiles", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		var p meta.Profile
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, meta.ProfileDisableByID(ctx, q, p.ID).Scan(&p))

		routes := mux.NewRouter()
		metaapi.NewHTTPUsermanagement(
			q,
			metaapi.HTTPUsermanagementOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/meta/u12t").Subrouter())
		srv := cmdtestx.NewTLSServer(t, q, routes)

		require.NoError(t, cmdtestx.Execute(t, genparser(t), "command", "ls", "--private-key-path", keypath, "--disabled", "--endpoint", srv.URL))
	})

	t.Run("combines status filters additively", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		var pending meta.Profile
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, pending).Scan(&pending))

		var enabled meta.Profile
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, enabled).Scan(&enabled))
		require.NoError(t, meta.ProfileEnable(ctx, q, enabled.ID).Scan(&enabled))

		routes := mux.NewRouter()
		metaapi.NewHTTPUsermanagement(
			q,
			metaapi.HTTPUsermanagementOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/meta/u12t").Subrouter())
		srv := cmdtestx.NewTLSServer(t, q, routes)

		require.NoError(t, cmdtestx.Execute(t, genparser(t), "command", "ls", "--private-key-path", keypath, "--pending", "--enabled", "--endpoint", srv.URL))
	})

	t.Run("text search with --query", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		var p meta.Profile
		p.Display = "searchable testuser"
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))

		routes := mux.NewRouter()
		metaapi.NewHTTPUsermanagement(
			q,
			metaapi.HTTPUsermanagementOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/meta/u12t").Subrouter())
		srv := cmdtestx.NewTLSServer(t, q, routes)

		require.NoError(t, cmdtestx.Execute(t, genparser(t), "command", "ls", "--private-key-path", keypath, "--query", "searchable", "--endpoint", srv.URL))
	})
}
