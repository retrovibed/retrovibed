package cmdmeta_test

import (
	"context"
	"path/filepath"
	"sync"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/cmd/cmdmeta"
	"github.com/retrovibed/retrovibed/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/cmd/cmdtestx"
	"github.com/retrovibed/retrovibed/httpauthtest"
	"github.com/retrovibed/retrovibed/internal/env"
	"github.com/retrovibed/retrovibed/internal/sqltestx"
	"github.com/retrovibed/retrovibed/internal/sqlx"
	"github.com/retrovibed/retrovibed/internal/sshx"
	"github.com/retrovibed/retrovibed/internal/testx"
	"github.com/retrovibed/retrovibed/meta"
	"github.com/retrovibed/retrovibed/metaapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIdenAdd(t *testing.T) {
	var cli struct {
		cmdopts.Global
		cmdopts.TLSConfig
		cmdopts.SSHID
		Identity cmdmeta.Identity `cmd:""`
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

	t.Run("creates profile with admin permissions", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		routes := mux.NewRouter()
		metaapi.NewHTTPUsermanagement(
			q,
			metaapi.HTTPUsermanagementOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/meta/u12t").Subrouter())
		srv := cmdtestx.NewTLSServer(t, q, routes)

		_, pub, err := sshx.UnsafeNewKeyGen().Generate()
		require.NoError(t, err)

		require.NoError(t, cmdtestx.Execute(t, parser, "identity", "add", "--private-key-path", keypath, "--endpoint", srv.Listener.Addr().String(), string(pub)))

		require.Equal(t, 2, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM meta_profiles"))(t))
		require.Equal(t, 1, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM meta_sso_identity_ssh WHERE profile_id != (SELECT profile_id FROM meta_sso_identity_ssh LIMIT 1)"))(t))

		pid, err := sqlx.String(ctx, q, "SELECT profile_id::text FROM meta_sso_identity_ssh ORDER BY created_at DESC LIMIT 1")
		require.NoError(t, err)

		var authz meta.Authz
		require.NoError(t, meta.AuthzFindByProfileID(ctx, q, sqlx.NewNullString(pid)).Scan(&authz))
		require.Equal(t, pid, authz.ProfileID)
		assert.False(t, authz.Usermanagement)
		assert.False(t, authz.LibraryModify)
		assert.True(t, authz.LibraryRead)
	})

	t.Run("idempotent - same key does not create duplicates", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		routes := mux.NewRouter()
		metaapi.NewHTTPUsermanagement(
			q,
			metaapi.HTTPUsermanagementOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/meta/u12t").Subrouter())
		srv := cmdtestx.NewTLSServer(t, q, routes)

		_, pub, err := sshx.UnsafeNewKeyGen().Generate()
		require.NoError(t, err)

		require.NoError(t, cmdtestx.Execute(t, parser, "identity", "add", "--private-key-path", keypath, "--endpoint", srv.Listener.Addr().String(), string(pub)))
		require.NoError(t, cmdtestx.Execute(t, parser, "identity", "add", "--private-key-path", keypath, "--endpoint", srv.Listener.Addr().String(), string(pub)))

		require.Equal(t, 2, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM meta_profiles"))(t))
		require.Equal(t, 1, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM meta_sso_identity_ssh WHERE profile_id != (SELECT profile_id FROM meta_sso_identity_ssh LIMIT 1)"))(t))
	})
}
