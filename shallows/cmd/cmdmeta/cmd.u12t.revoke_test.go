package cmdmeta_test

import (
	"path/filepath"
	"testing"

	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdmeta"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdtestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/stretchr/testify/require"
)

func TestPendingRevoke(t *testing.T) {
	genparser := cmdtestx.Genparser(cmdmeta.Usermanagement{})
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

		require.NoError(t, cmdtestx.Execute(t, genparser(t), "command", "revoke", "--private-key-path", keypath, "--endpoint", srv.URL, p.ID))

		require.Equal(t, 0, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM authz_meta WHERE profile_id = '"+p.ID+"'"))(t))
	})
}
