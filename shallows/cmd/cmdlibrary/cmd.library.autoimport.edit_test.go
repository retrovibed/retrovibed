package cmdlibrary_test

import (
	"bytes"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/alecthomas/kong"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdlibrary"
	cmdtestx "github.com/retrovibed/retrovibed/shallows/cmd/cmdtestx/v2"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/stretchr/testify/require"
)

func TestAutoimportEdit(t *testing.T) {
	t.Run("only overridden fields change", func(t *testing.T) {
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
		harness := cmdtestx.New(t, &cmdlibrary.Commands{}).Options(cmdtestx.Standard(t), kong.Writers(&buf, nil), kong.Vars{"vars_cores": strconv.Itoa(runtime.GOMAXPROCS(0)), "vars_media_directory": t.TempDir(), "vars_user_configuration_directory": t.TempDir()})
		require.NoError(t, harness.Run("autoimport", "edit", d.ID, "--mode", "move", "--private-key-path", keypath, "--endpoint", srv.URL))

		var updated library.AutoimportDirectory
		require.NoError(t, library.AutoimportDirectoryFindByID(ctx, q, d.ID).Scan(&updated))
		require.Equal(t, library.AutoimportModeMove, updated.Mode)
		require.Equal(t, d.Description, updated.Description)
		require.Equal(t, time.Hour, updated.Debounce)
		require.Equal(t, d.LibraryDirectoryID, updated.LibraryDirectoryID)

		require.Contains(t, buf.String(), "mode=move")
	})
}
