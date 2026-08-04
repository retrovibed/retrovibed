package cmdddisc_test

import (
	"fmt"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/gorilla/mux"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/searchplugin"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdddisc"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdtestx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/torrentx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestDiscoveryCreate(t *testing.T) {
	t.Run("creates a discovery entry", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		infohash := int160.Random()

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
			"discovery", "create",
			"--private-key-path", keypath,
			"--endpoint", srv.URL,
			"--magnet", fmt.Sprintf("magnet:?xt=urn:btih:%s", infohash.String()),
		))

		expectedID := torrentx.HashUID(&infohash)

		found := sqlx.Scan(tracking.UnknownSearch(ctx, q, tracking.UnknownSearchBuilder().Where(tracking.UnknownHashQueryByIDs(expectedID))))
		rows := 0
		for range found.Iter() {
			rows++
		}
		require.NoError(t, found.Err())
		require.Equal(t, 1, rows)
	})
}
