package cmdmeta_test

import (
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdmeta"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdtestx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/sshx"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestIdenRegister(t *testing.T) {
	genparser := cmdtestx.Genparser(cmdmeta.Identity{})

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

		require.NoError(t, cmdtestx.Execute(t, genparser(t), "command", "register", "--private-key-path", keypath, "--endpoint", srv.URL))

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

		require.NoError(t, cmdtestx.Execute(t, parser, "command", "register", "--private-key-path", keypath, "--endpoint", srv.URL))
		require.NoError(t, cmdtestx.Execute(t, parser, "command", "register", "--private-key-path", keypath, "--endpoint", srv.URL))

		require.Equal(t, 1, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM meta_profiles"))(t))
		require.Equal(t, 1, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM meta_sso_identity_ssh"))(t))
	})
}
