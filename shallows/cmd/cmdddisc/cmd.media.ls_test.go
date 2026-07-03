package cmdddisc_test

import (
	"path/filepath"
	"testing"

	"github.com/gofrs/uuid/v5"
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

func TestMediaLs(t *testing.T) {
	t.Run("lists all media", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		id := int160.Random()
		d := ddisc.NewDiscovered(&id)
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))

		routes := mux.NewRouter()
		ddiscapi.NewHTTPMedia(
			q,
			ddiscapi.HTTPMediaOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/ddisc/media").Subrouter())
		srv := cmdtestx.NewTLSServer(t, q, routes)

		require.NoError(t, cmdtestx.Execute(t, genparser(t), "media", "ls",
			"--private-key-path", keypath,
			"--insecure",
			"--library", srv.Listener.Addr().String(),
		))
	})

	t.Run("filters by known media id", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		knownMediaID := uuid.Must(uuid.NewV7()).String()

		matchID := int160.Random()
		match := ddisc.NewDiscovered(&matchID, ddisc.DiscoveredOptionKnownMedia(knownMediaID))
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, match).Scan(&match))

		otherID := int160.Random()
		other := ddisc.NewDiscovered(&otherID)
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, other).Scan(&other))

		routes := mux.NewRouter()
		ddiscapi.NewHTTPMedia(
			q,
			ddiscapi.HTTPMediaOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/ddisc/media").Subrouter())
		srv := cmdtestx.NewTLSServer(t, q, routes)

		require.NoError(t, cmdtestx.Execute(t, genparser(t), "media", "ls",
			"--private-key-path", keypath,
			"--insecure",
			"--library", srv.Listener.Addr().String(),
			"--known-media-id", knownMediaID,
		))
	})
}
