package cmdddisc_test

import (
	"database/sql"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/searchplugin"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdddisc"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdtestx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestDiscoveryDelete(t *testing.T) {
	t.Run("removes an entry", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		var uh tracking.UnknownHash
		require.NoError(t, testx.Fake(&uh, tracking.UnknownHashOptionTestDefaults))
		require.NoError(t, tracking.UnknownHashInsertWithDefaults(ctx, q, uh).Scan(&uh))

		routes := mux.NewRouter()
		ddiscapi.NewHTTPDiscovery(
			q,
			searchplugin.Unimplemented{},
			ddisc.UnimplementedStrategy{},
			tracking.NewURIImport(q, http.DefaultClient, fsx.DirVirtual(t.TempDir())),
			ddiscapi.HTTPDiscoveryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/ddisc/discovery").Subrouter())
		srv := cmdtestx.NewTLSServer(t, q, routes)

		require.NoError(t, cmdtestx.Execute(
			t,
			cmdtestx.Genparser(cmdddisc.Commands{})(t), "command",
			"discovery", "delete",
			"--private-key-path", keypath,
			"--endpoint", srv.URL,
			"--id", uh.ID,
		))

		var target = sql.ErrNoRows
		require.ErrorAs(t, tracking.UnknownHashDeleteByID(ctx, q, uh.ID).Scan(&uh), &target)
	})
}
