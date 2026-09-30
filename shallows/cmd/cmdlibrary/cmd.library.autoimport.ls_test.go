package cmdlibrary_test

import (
	"bytes"
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdlibrary"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdtestx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/stretchr/testify/require"
)

func TestAutoimportLs(t *testing.T) {
	t.Run("lists monitored directories", func(t *testing.T) {
		var (
			d library.AutoimportDirectory
		)

		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		require.NoError(t, testx.Fake(&d, library.AutoimportDirectoryOptionTestDefaults, library.AutoimportDirectoryOptionPath(t.TempDir())))
		require.NoError(t, library.AutoimportDirectoryInsertWithDefaults(ctx, q, d).Scan(&d))

		routes := mux.NewRouter()
		media.NewHTTPAutoimport(
			q,
			asyncx.NewWakeup(t.Context()),
			media.HTTPAutoimportOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/autoimport").Subrouter())
		srv := cmdtestx.NewTLSServer(t, q, routes)

		var buf bytes.Buffer
		genparser := cmdtestx.Genparser(cmdlibrary.Commands{}, kong.Writers(&buf, nil), kong.Vars{"vars_cores": strconv.Itoa(runtime.GOMAXPROCS(0)), "vars_media_directory": t.TempDir(), "vars_user_configuration_directory": t.TempDir()})
		require.NoError(t, cmdtestx.Execute(t, genparser(t), "command", "autoimport", "ls", "--private-key-path", keypath, "--endpoint", srv.URL))

		require.Contains(t, buf.String(), fmt.Sprintf("id='%s' path='%s' mode=copy debounce=1h0m0s", d.ID, d.Path))
	})
}
