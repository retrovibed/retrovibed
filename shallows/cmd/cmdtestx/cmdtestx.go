package cmdtestx

import (
	"context"
	"database/sql"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/sshx"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/meta/identityssh"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func Genparser[T any](cmd T, opts ...kong.Option) func(t *testing.T) *kong.Kong {
	return func(t *testing.T) *kong.Kong {
		t.Helper()
		var cli struct {
			cmdopts.Global
			cmdopts.TLSConfig
			cmdopts.SSHID
			cmdopts.Endpoint
			Command T `cmd:""`
		}

		cli.Context, cli.Shutdown = context.WithCancel(context.Background())
		cli.Cleanup = &sync.WaitGroup{}
		return kong.Must(
			&cli,
			append(opts,
				kong.Bind(&cli.TLSConfig),
				kong.Bind(&cli.Global),
				kong.Bind(&cli.SSHID),
				kong.Bind(&cli.Endpoint),
				kong.Vars{
					"vars_private_key":                  env.PrivateKeyPath(),
					"vars_user_configuration_directory": t.TempDir(),
					"env_http_endpoint":                 env.Endpoint,
				},
				kong.NamedMapper("durationinf", kong.MapperFunc(cmdopts.ParseDurationInf)),
				kong.NamedMapper("envvar", kong.MapperFunc(cmdopts.ParseEnviron)),
			)...,
		)
	}
}

func Execute(t *testing.T, parser *kong.Kong, cmd ...string) error {
	t.Helper()

	kctx, err := parser.Parse(cmd)
	require.NoError(t, err)
	return kctx.Run()
}

func Admin(t *testing.T, ctx context.Context, q *sql.DB, keypath string) (admin meta.Profile, authz meta.Authz, iden identityssh.Identity) {
	t.Helper()

	signer, err := sshx.AutoCached(sshx.NewKeyGen(), keypath)
	require.NoError(t, err)

	require.NoError(t, testx.Fake(&admin, meta.ProfileOptionTestDefaults))
	require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, admin).Scan(&admin))
	require.NoError(t, meta.ProfileEnable(ctx, q, admin.ID).Scan(&admin))
	require.NoError(t, testx.Fake(&authz, meta.AuthzOptionProfileID(admin.ID), meta.AuthzOptionAdmin))
	require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, authz).Scan(&authz))
	require.NoError(t, identityssh.IdentityInsertWithDefaults(ctx, q, identityssh.Identity{
		ID:        sshx.FingerprintMD5(signer.PublicKey()),
		ProfileID: admin.ID,
		PublicKey: sshx.EncodeBase64PublicKey(signer.PublicKey()),
	}).Scan(&iden))

	return admin, authz, iden
}

func NewTLSServer(t *testing.T, q sqlx.Queryer, routes *mux.Router) *httptest.Server {
	t.Helper()

	metamux := routes.PathPrefix("/meta").Subrouter()
	metaapi.NewHTTPAuthz(
		q,
		metaapi.HTTPAuthzOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(metamux.PathPrefix("/authz").Subrouter())

	oauth2mux := routes.PathPrefix("/oauth2").Subrouter()
	metaapi.NewSSHOauth2(
		q,
		metaapi.SSHOauth2OptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(oauth2mux.PathPrefix("/ssh").Subrouter())

	srv := httptest.NewTLSServer(routes)
	t.Cleanup(srv.Close)
	return srv
}
