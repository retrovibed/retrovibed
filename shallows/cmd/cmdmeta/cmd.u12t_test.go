package cmdmeta_test

import (
	"context"
	"path/filepath"
	"sync"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdmeta"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdtestx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLs(t *testing.T) {
	var cli struct {
		cmdopts.Global
		cmdopts.TLSConfig
		cmdopts.SSHID
		Usermanagement cmdmeta.Usermanagement `cmd:""`
	}

	cli.Context, cli.Shutdown = context.WithCancel(context.Background())
	cli.Cleanup = &sync.WaitGroup{}
	parser := kong.Must(
		&cli,
		kong.Bind(&cli.TLSConfig),
		kong.Bind(&cli.Global),
		kong.Bind(&cli.SSHID),
		kong.Vars{
			"vars_private_key":                  env.PrivateKeyPath(),
			"vars_user_configuration_directory": t.TempDir(),
		},
	)

	t.Run("lists pending profiles", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		var pending meta.Profile
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, pending).Scan(&pending))

		routes := mux.NewRouter()
		metaapi.NewHTTPUsermanagement(
			q,
			metaapi.HTTPUsermanagementOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/meta/u12t").Subrouter())
		srv := cmdtestx.NewTLSServer(t, q, routes)

		require.NoError(t, cmdtestx.Execute(t, parser, "usermanagement", "ls", "--private-key-path", keypath, "--insecure", "--endpoint", srv.Listener.Addr().String()))
	})

	t.Run("lists pending profiles with --pending flag", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		var pending meta.Profile
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, pending).Scan(&pending))

		routes := mux.NewRouter()
		metaapi.NewHTTPUsermanagement(
			q,
			metaapi.HTTPUsermanagementOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/meta/u12t").Subrouter())
		srv := cmdtestx.NewTLSServer(t, q, routes)

		require.NoError(t, cmdtestx.Execute(t, parser, "usermanagement", "ls", "--private-key-path", keypath, "--pending", "--insecure", "--endpoint", srv.Listener.Addr().String()))
	})
}

func TestGrant(t *testing.T) {
	var cli struct {
		cmdopts.Global
		cmdopts.TLSConfig
		cmdopts.SSHID
		Usermanagement cmdmeta.Usermanagement `cmd:""`
	}

	cli.Context, cli.Shutdown = context.WithCancel(context.Background())
	cli.Cleanup = &sync.WaitGroup{}

	parser := kong.Must(
		&cli,
		kong.Bind(&cli.TLSConfig),
		kong.Bind(&cli.Global),
		kong.Bind(&cli.SSHID),
		kong.Vars{
			"vars_private_key":                  env.PrivateKeyPath(),
			"vars_user_configuration_directory": t.TempDir(),
		},
	)

	t.Run("enables profile and grants library read", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		var pending meta.Profile
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, pending).Scan(&pending))

		routes := mux.NewRouter()
		metaapi.NewHTTPUsermanagement(
			q,
			metaapi.HTTPUsermanagementOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/meta/u12t").Subrouter())
		srv := cmdtestx.NewTLSServer(t, q, routes)

		require.NoError(t, cmdtestx.Execute(t, parser, "usermanagement", "grant", "--private-key-path", keypath, "--insecure", "--endpoint", srv.Listener.Addr().String(), pending.ID))

		var authz meta.Authz
		require.NoError(t, meta.AuthzFindByProfileID(ctx, q, sqlx.NewNullString(pending.ID)).Scan(&authz))
		assert.True(t, authz.LibraryRead)
		assert.False(t, authz.Usermanagement)
		assert.False(t, authz.LibraryModify)

		require.Equal(t, 1, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM meta_profiles WHERE disabled_pending_approval_at > NOW() AND id = '"+pending.ID+"'"))(t))
	})

	t.Run("idempotent - granting twice does not error", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		var pending meta.Profile
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, pending).Scan(&pending))

		routes := mux.NewRouter()
		metaapi.NewHTTPUsermanagement(
			q,
			metaapi.HTTPUsermanagementOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/meta/u12t").Subrouter())
		srv := cmdtestx.NewTLSServer(t, q, routes)

		require.NoError(t, cmdtestx.Execute(t, parser, "usermanagement", "grant", "--private-key-path", keypath, "--insecure", "--endpoint", srv.Listener.Addr().String(), pending.ID))
		require.NoError(t, cmdtestx.Execute(t, parser, "usermanagement", "grant", "--private-key-path", keypath, "--insecure", "--endpoint", srv.Listener.Addr().String(), pending.ID))

		var authz meta.Authz
		require.NoError(t, meta.AuthzFindByProfileID(ctx, q, sqlx.NewNullString(pending.ID)).Scan(&authz))
		assert.True(t, authz.LibraryRead)
	})
}

func TestPendingRevoke(t *testing.T) {
	var cli struct {
		cmdopts.Global
		cmdopts.TLSConfig
		cmdopts.SSHID
		Usermanagement cmdmeta.Usermanagement `cmd:""`
	}

	cli.Context, cli.Shutdown = context.WithCancel(context.Background())
	cli.Cleanup = &sync.WaitGroup{}

	parser := kong.Must(
		&cli,
		kong.Bind(&cli.TLSConfig),
		kong.Bind(&cli.Global),
		kong.Bind(&cli.SSHID),
		kong.Vars{
			"vars_private_key":                  env.PrivateKeyPath(),
			"vars_user_configuration_directory": t.TempDir(),
		},
	)

	t.Run("removes authz for a profile", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		var p meta.Profile
		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, meta.ProfileEnable(ctx, q, p.ID).Scan(&p))

		authz := meta.Authz{ProfileID: p.ID, LibraryRead: true}
		require.NoError(t, meta.AuthzUpsertWithDefaults(ctx, q, authz).Scan(&authz))

		routes := mux.NewRouter()
		srv := cmdtestx.NewTLSServer(t, q, routes)

		require.NoError(t, cmdtestx.Execute(t, parser, "usermanagement", "revoke", "--private-key-path", keypath, "--insecure", "--endpoint", srv.Listener.Addr().String(), p.ID))

		require.Equal(t, 0, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM authz_meta WHERE profile_id = '"+p.ID+"'"))(t))
	})
}
