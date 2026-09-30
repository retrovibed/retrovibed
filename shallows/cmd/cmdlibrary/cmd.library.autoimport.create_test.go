package cmdlibrary_test

import (
	"bytes"
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/alecthomas/kong"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdlibrary"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdtestx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/stretchr/testify/require"
)

func TestAutoimportCreate(t *testing.T) {
	t.Run("creates a monitored directory", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)
		root := t.TempDir()

		cmdtestx.Admin(t, ctx, q, keypath)

		routes := mux.NewRouter()
		media.NewHTTPAutoimport(
			q,
			asyncx.NewWakeup(t.Context()),
			media.HTTPAutoimportOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/autoimport").Subrouter())
		srv := cmdtestx.NewTLSServer(t, q, routes)

		var buf bytes.Buffer
		genparser := cmdtestx.Genparser(cmdlibrary.Commands{}, kong.Writers(&buf, nil), kong.Vars{"vars_cores": strconv.Itoa(runtime.GOMAXPROCS(0)), "vars_media_directory": t.TempDir(), "vars_user_configuration_directory": t.TempDir()})
		require.NoError(t, cmdtestx.Execute(t, genparser(t), "command", "autoimport", "create", root, "--debounce", "1m", "--mode", "move", "--description", "inbox", "--private-key-path", keypath, "--endpoint", srv.URL))

		d, err := sqlx.ScanOne(library.AutoimportDirectorySearch(ctx, q, library.AutoimportDirectorySearchBuilder()))
		require.NoError(t, err)
		require.Equal(t, root, d.Path)
		require.Equal(t, "inbox", d.Description)
		require.Equal(t, time.Minute, d.Debounce)
		require.Equal(t, library.AutoimportModeMove, d.Mode)

		require.Contains(t, buf.String(), fmt.Sprintf("id='%s' path='%s' mode=move debounce=1m0s", d.ID, root))
	})
}
