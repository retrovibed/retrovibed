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

func TestUpdate(t *testing.T) {
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

	t.Run("sets the display name", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		var p meta.Profile
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))

		parser := genparser(t)
		srv := newServer(t, q)

		require.NoError(t, cmdtestx.Execute(t, parser, "command", "update", "--private-key-path", keypath, "--endpoint", srv.URL, "--display", "New Name", p.ID))

		var updated meta.Profile
		require.NoError(t, meta.ProfileFindByID(ctx, q, p.ID).Scan(&updated))
		assert.Equal(t, "New Name", updated.Display)
		assert.True(t, updated.DisabledAt.Equal(p.DisabledAt))
		assert.True(t, updated.DisabledManuallyAt.Equal(p.DisabledManuallyAt))
		assert.True(t, updated.DisabledPendingApprovalAt.Equal(p.DisabledPendingApprovalAt))
	})

	t.Run("no flags is a no-op that does not error", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		var p meta.Profile
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))

		parser := genparser(t)
		srv := newServer(t, q)

		require.NoError(t, cmdtestx.Execute(t, parser, "command", "update", "--private-key-path", keypath, "--endpoint", srv.URL, p.ID))

		var updated meta.Profile
		require.NoError(t, meta.ProfileFindByID(ctx, q, p.ID).Scan(&updated))
		assert.Equal(t, p.Display, updated.Display)
	})
}
