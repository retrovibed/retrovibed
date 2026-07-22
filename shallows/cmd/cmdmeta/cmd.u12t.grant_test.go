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
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGrant(t *testing.T) {
	genparser := cmdtestx.Genparser(cmdmeta.Usermanagement{})

	newServer := func(t *testing.T, q sqlx.Queryer) *httptest.Server {
		t.Helper()
		routes := mux.NewRouter()
		metaapi.NewHTTPUsermanagement(
			q,
			metaapi.HTTPUsermanagementOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/meta/u12t").Subrouter())
		return cmdtestx.NewTLSServer(t, q, routes)
	}

	t.Run("grants library read by default", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		var p meta.Profile
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))

		parser := genparser(t)
		srv := newServer(t, q)

		require.NoError(t, cmdtestx.Execute(t, parser, "command", "grant", "--private-key-path", keypath, "--endpoint", srv.URL, p.ID))

		var authz meta.Authz
		require.NoError(t, meta.AuthzFindByProfileID(ctx, q, sqlx.NewNullString(p.ID)).Scan(&authz))
		assert.True(t, authz.LibraryRead)
		assert.False(t, authz.LibraryModify)
		assert.False(t, authz.BillingRead)
		assert.False(t, authz.BillingModify)
		assert.False(t, authz.CommunityModify)
		assert.False(t, authz.Usermanagement)
		require.Equal(t, 1, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM meta_profiles WHERE disabled_pending_approval_at > NOW() AND id = '"+p.ID+"'"))(t))
	})

	t.Run("grants only the specified permissions", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		var p meta.Profile
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))

		parser := genparser(t)
		srv := newServer(t, q)

		require.NoError(t, cmdtestx.Execute(t, parser, "command", "grant", "--private-key-path", keypath, "--endpoint", srv.URL, "--no-library-read", "--library-modify", "--usermanagement", p.ID))

		var authz meta.Authz
		require.NoError(t, meta.AuthzFindByProfileID(ctx, q, sqlx.NewNullString(p.ID)).Scan(&authz))
		assert.False(t, authz.LibraryRead)
		assert.True(t, authz.LibraryModify)
		assert.False(t, authz.BillingRead)
		assert.False(t, authz.BillingModify)
		assert.False(t, authz.CommunityModify)
		assert.True(t, authz.Usermanagement)
	})

	t.Run("grants all permissions", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		var p meta.Profile
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))

		parser := genparser(t)
		srv := newServer(t, q)

		require.NoError(t, cmdtestx.Execute(t, parser, "command", "grant", "--private-key-path", keypath, "--endpoint", srv.URL, "--library-read", "--library-modify", "--billing-read", "--billing-modify", "--community-modify", "--usermanagement", p.ID))

		var authz meta.Authz
		require.NoError(t, meta.AuthzFindByProfileID(ctx, q, sqlx.NewNullString(p.ID)).Scan(&authz))
		assert.True(t, authz.LibraryRead)
		assert.True(t, authz.LibraryModify)
		assert.True(t, authz.BillingRead)
		assert.True(t, authz.BillingModify)
		assert.True(t, authz.CommunityModify)
		assert.True(t, authz.Usermanagement)
	})

	t.Run("idempotent - granting twice does not error", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		var p meta.Profile
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))

		parser := genparser(t)
		srv := newServer(t, q)

		require.NoError(t, cmdtestx.Execute(t, parser, "command", "grant", "--private-key-path", keypath, "--endpoint", srv.URL, p.ID))
		require.NoError(t, cmdtestx.Execute(t, parser, "command", "grant", "--private-key-path", keypath, "--endpoint", srv.URL, p.ID))

		var authz meta.Authz
		require.NoError(t, meta.AuthzFindByProfileID(ctx, q, sqlx.NewNullString(p.ID)).Scan(&authz))
		assert.True(t, authz.LibraryRead)
	})
}
