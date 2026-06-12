package cmdmeta_test

import (
	"context"
	"net/http/httptest"
	"path/filepath"
	"sync"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdmeta"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdtestx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/sshx"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestIdenRegister(t *testing.T) {
	genparser := func(t *testing.T) *kong.Kong {
		var cli struct {
			cmdopts.Global
			cmdopts.TLSConfig
			cmdopts.SSHID
			Identity cmdmeta.Identity `cmd:""`
		}

		cli.Context, cli.Shutdown = context.WithCancel(context.Background())
		cli.Cleanup = &sync.WaitGroup{}

		return kong.Must(
			&cli,
			kong.Bind(&cli.TLSConfig),
			kong.Bind(&cli.Global),
			kong.Bind(&cli.SSHID),
			kong.Vars{
				"vars_private_key":                  env.PrivateKeyPath(),
				"vars_user_configuration_directory": t.TempDir(),
			},
		)
	}

	newServer := func(t *testing.T, q sqlx.Queryer) *httptest.Server {
		t.Helper()
		routes := mux.NewRouter()
		metaapi.NewHTTP(
			q,
			metaapi.HTTPOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/sso").Subrouter())
		return cmdtestx.NewTLSServer(t, q, routes)
	}

	t.Run("registers a new identity, creating a pending profile", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		_, err := sshx.AutoCached(sshx.NewKeyGen(), keypath)
		require.NoError(t, err)

		q := sqltestx.Metadatabase(t)
		srv := newServer(t, q)

		require.NoError(t, cmdtestx.Execute(t, genparser(t), "identity", "register", "--private-key-path", keypath, "--insecure", "--endpoint", srv.Listener.Addr().String()))

		require.Equal(t, 1, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM meta_profiles"))(t))
		require.Equal(t, 1, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM meta_sso_identity_ssh"))(t))

		pid, err := sqlx.String(ctx, q, "SELECT profile_id::text FROM meta_sso_identity_ssh")
		require.NoError(t, err)
		require.NotEmpty(t, pid)

		var p meta.Profile
		require.NoError(t, meta.ProfileFindByID(ctx, q, pid).Scan(&p))
	})

	t.Run("idempotent - registering twice does not error or duplicate", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		_, err := sshx.AutoCached(sshx.NewKeyGen(), keypath)
		require.NoError(t, err)

		q := sqltestx.Metadatabase(t)
		srv := newServer(t, q)
		parser := genparser(t)

		require.NoError(t, cmdtestx.Execute(t, parser, "identity", "register", "--private-key-path", keypath, "--insecure", "--endpoint", srv.Listener.Addr().String()))
		require.NoError(t, cmdtestx.Execute(t, parser, "identity", "register", "--private-key-path", keypath, "--insecure", "--endpoint", srv.Listener.Addr().String()))

		require.Equal(t, 1, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM meta_profiles"))(t))
		require.Equal(t, 1, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM meta_sso_identity_ssh"))(t))
	})
}
