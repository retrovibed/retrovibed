package cmdmeta_test

import (
	"bytes"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdmeta"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdtestx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestDeviceLs(t *testing.T) {
	t.Run("lists known devices", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		var d meta.Daemon
		require.NoError(t, testx.Fake(&d, meta.DaemonOptionTestDefaults, meta.DaemonOptionMaybeID))
		require.NoError(t, meta.DaemonInsertWithDefaults(ctx, q, d).Scan(&d))

		routes := mux.NewRouter()
		metaapi.NewHTTPDaemons(
			q,
			metaapi.HTTPDaemonsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/meta/d").Subrouter())
		srv := cmdtestx.NewTLSServer(t, q, routes)

		var buf bytes.Buffer
		genparser := cmdtestx.Genparser(cmdmeta.Device{}, kong.Writers(&buf, nil))
		require.NoError(t, cmdtestx.Execute(t, genparser(t), "command", "ls", "--private-key-path", keypath, "--endpoint", srv.URL))

		require.Contains(t, buf.String(), fmt.Sprintf("id='%s' hostname='%s' description='%s'", d.ID, d.Hostname, d.Description))
	})
}
