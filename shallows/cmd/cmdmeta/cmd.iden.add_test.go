package cmdmeta_test

import (
	"path/filepath"
	"testing"

	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdmeta"
	cmdtestx "github.com/retrovibed/retrovibed/shallows/cmd/cmdtestx/v2"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/sshx"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIdenAdd(t *testing.T) {
	harness := cmdtestx.New(t, &cmdmeta.Identity{}).Options(cmdtestx.Standard(t))

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

		require.NoError(t, harness.Run("add", "--private-key-path", keypath, "--endpoint", srv.URL, string(pub)))

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

	t.Run("display defaults to the ssh key comment when username not provided", func(t *testing.T) {
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
		pub = sshx.Comment(pub, "alice@example.com")

		require.NoError(t, harness.Run("add", "--private-key-path", keypath, "--endpoint", srv.URL, string(pub)))

		display, err := sqlx.String(ctx, q, "SELECT display FROM meta_profiles WHERE id = (SELECT profile_id FROM meta_sso_identity_ssh ORDER BY created_at DESC LIMIT 1)")
		require.NoError(t, err)
		require.Equal(t, "alice@example.com", display)
	})

	t.Run("display uses provided username over the ssh key comment", func(t *testing.T) {
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
		pub = sshx.Comment(pub, "alice@example.com")

		require.NoError(t, harness.Run("add", "--private-key-path", keypath, "--endpoint", srv.URL, "--username", "alice", string(pub)))

		display, err := sqlx.String(ctx, q, "SELECT display FROM meta_profiles WHERE id = (SELECT profile_id FROM meta_sso_identity_ssh ORDER BY created_at DESC LIMIT 1)")
		require.NoError(t, err)
		require.Equal(t, "alice", display)
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

		require.NoError(t, harness.Run("add", "--private-key-path", keypath, "--endpoint", srv.URL, string(pub)))
		require.NoError(t, harness.Run("add", "--private-key-path", keypath, "--endpoint", srv.URL, string(pub)))

		require.Equal(t, 2, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM meta_profiles"))(t))
		require.Equal(t, 1, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM meta_sso_identity_ssh WHERE profile_id != (SELECT profile_id FROM meta_sso_identity_ssh LIMIT 1)"))(t))
	})
}
