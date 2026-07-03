package cmdddisc_test

import (
	"encoding/hex"
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
	"github.com/retrovibed/retrovibed/shallows/internal/torrentx"
	"github.com/stretchr/testify/require"
)

func TestMediaCreate(t *testing.T) {
	t.Run("creates a media record", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		id := int160.Random()
		routes := mux.NewRouter()
		ddiscapi.NewHTTPMedia(
			q,
			ddiscapi.HTTPMediaOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/ddisc/media").Subrouter())
		srv := cmdtestx.NewTLSServer(t, q, routes)

		require.NoError(t, cmdtestx.Execute(t, genparser(t), "media", "create",
			"--private-key-path", keypath,
			"--insecure",
			"--library", srv.Listener.Addr().String(),
			"--infohash", hex.EncodeToString(id.Bytes()),
			"--title", "derp",
			"--description", "a description",
			"--mimetype", "video/mp4",
		))

		var stored ddisc.Discovered
		require.NoError(t, ddisc.DiscoveredFindByID(ctx, q, torrentx.HashUID(&id)).Scan(&stored))
		require.Equal(t, "derp", stored.Title)
		require.Equal(t, "a description", stored.Description)
		require.Equal(t, "video/mp4", stored.Mimetype)
	})
}
