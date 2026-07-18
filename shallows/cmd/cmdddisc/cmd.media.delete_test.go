package cmdddisc_test

import (
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/gorilla/mux"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdtestx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/stretchr/testify/require"
)

func TestMediaDelete(t *testing.T) {
	t.Run("deletes a media record", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		id := int160.Random()
		d := ddisc.NewDiscovered(&id, "")
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))

		routes := mux.NewRouter()
		ddiscapi.NewHTTPMedia(
			q,
			ddiscapi.HTTPMediaOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/ddisc/media").Subrouter())
		srv := cmdtestx.NewTLSServer(t, q, routes)

		require.NoError(t, cmdtestx.Execute(t, genparser(t), "media", "delete",
			"--private-key-path", keypath,
			"--insecure",
			"--library", srv.Listener.Addr().String(),
			"--id", d.ID,
		))

		var target = sql.ErrNoRows
		require.ErrorAs(t, ddisc.DiscoveredFindByID(ctx, q, d.ID).Scan(&d), &target)
	})
}
