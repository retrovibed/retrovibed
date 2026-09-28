package cmdddisc_test

import (
	"path/filepath"
	"testing"
	"time"

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

func TestLocateRetry(t *testing.T) {
	t.Run("clears the cooldown and resets attempts", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		var l ddisc.Locate
		require.NoError(t, ddisc.LocateInsertWithDefaults(ctx, q, ddisc.NewLocate("ubuntu", "video")).Scan(&l))
		require.NoError(t, ddisc.LocateCooldown(ctx, q, l.ID, time.Now().Add(time.Hour)).Scan(&l))

		routes := mux.NewRouter()
		ddiscapi.NewHTTPLocate(
			q,
			asyncx.NewWakeup(ctx),
			ddiscapi.HTTPLocateOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/l").Subrouter())
		srv := cmdtestx.NewTLSServer(t, q, routes)

		require.NoError(t, cmdtestx.Execute(t, cmdtestx.Genparser(cmdddisc.Commands{})(t), "command", "locate", "retry",
			"--private-key-path", keypath,
			"--endpoint", srv.URL,
			"--id", l.ID,
			"--reset-attempts",
		))

		require.NoError(t, ddisc.LocateFindByID(ctx, q, l.ID).Scan(&l))
		require.EqualValues(t, 0, l.Attempts)
		require.False(t, l.NextCheckAt.After(time.Now()))
	})
}
