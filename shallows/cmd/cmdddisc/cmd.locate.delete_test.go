package cmdddisc_test

import (
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdddisc"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdtestx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/stretchr/testify/require"
)

func TestLocateDelete(t *testing.T) {
	t.Run("removes an entry", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		var l ddisc.Locate
		require.NoError(t, ddisc.LocateInsertWithDefaults(ctx, q, ddisc.NewLocate("ubuntu", "video")).Scan(&l))

		routes := mux.NewRouter()
		ddiscapi.NewHTTPLocate(
			q,
			asyncx.NewWakeup(ctx),
			ddiscapi.HTTPLocateOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/l").Subrouter())
		srv := cmdtestx.NewTLSServer(t, q, routes)

		require.NoError(t, cmdtestx.Execute(t, cmdtestx.Genparser(cmdddisc.Commands{})(t), "command", "locate", "delete",
			"--private-key-path", keypath,
			"--endpoint", srv.URL,
			"--id", l.ID,
		))

		var target = sql.ErrNoRows
		require.ErrorAs(t, ddisc.LocateFindByID(ctx, q, l.ID).Scan(&l), &target)
	})
}
