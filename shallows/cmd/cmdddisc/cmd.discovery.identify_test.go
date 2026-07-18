package cmdddisc_test

import (
	"path/filepath"
	"testing"

	"github.com/gorilla/mux"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/autobind"
	"github.com/james-lawrence/torrent/bencode"
	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/james-lawrence/torrent/storage"
	"github.com/james-lawrence/torrent/torrenttest"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdtestx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/torrenttestx"
	"github.com/retrovibed/retrovibed/shallows/internal/torrentx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestDiscoveryIdentify(t *testing.T) {
	t.Run("fails fast when no peers are reachable for the infohash, leaving the entry in place", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		var uh tracking.UnknownHash
		require.NoError(t, testx.Fake(&uh, tracking.UnknownHashOptionTestDefaults))
		require.NoError(t, tracking.UnknownHashInsertWithDefaults(ctx, q, uh).Scan(&uh))

		routes := mux.NewRouter()
		ddiscapi.NewHTTPDiscovery(
			q,
			ddiscapi.HTTPDiscoveryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/ddisc/discovery").Subrouter())
		ddiscapi.NewHTTPMedia(
			q,
			ddiscapi.HTTPMediaOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/ddisc/media").Subrouter())
		srv := cmdtestx.NewTLSServer(t, q, routes)

		require.Error(t, cmdtestx.Execute(t, genparser(t), "discovery", "identify",
			"--private-key-path", keypath,
			"--insecure",
			"--library", srv.Listener.Addr().String(),
			"--id", uh.ID,
			"--peer-timeout", "1s",
			"--info-timeout", "1s",
		))

		var found tracking.UnknownHash
		require.NoError(t, tracking.UnknownHashDeleteByID(ctx, q, uh.ID).Scan(&found), "the entry should remain untouched after a failed identify attempt")
	})

	t.Run("identifies and persists a real torrent as media, then removes the source lead", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		keypath := filepath.Join(t.TempDir(), "id")
		q := sqltestx.Metadatabase(t)

		cmdtestx.Admin(t, ctx, q, keypath)

		target := torrenttestx.QuickDHT(t, dht.OptionBootstrapNodesNone)
		seeder := torrenttestx.QuickClientBinder(
			t,
			autobind.New(autobind.EnableDHT(target)),
		)
		defer seeder.Close()

		seedDir := t.TempDir()
		info, _, err := torrenttest.Random(seedDir, 32*bytesx.KiB)
		require.NoError(t, err)

		encoded, err := bencode.Marshal(info)
		require.NoError(t, err)
		hash := metainfo.NewHashFromBytes(encoded)
		id := int160.FromBytes(hash.Bytes())

		seedermd, err := torrent.NewFromInfo(info, torrent.OptionStorage(storage.NewFile(seedDir)))
		require.NoError(t, err)
		_, _, err = seeder.Start(seedermd, torrent.TuneAnnounceUntilComplete, torrent.TuneNewConns)
		require.NoError(t, err)

		var uh tracking.UnknownHash
		require.NoError(t, testx.Fake(&uh, tracking.UnknownHashOptionTestDefaults))
		uh.ID = torrentx.HashUID(&id)
		uh.Infohash = id.Bytes()
		require.NoError(t, tracking.UnknownHashInsertWithDefaults(ctx, q, uh).Scan(&uh))

		routes := mux.NewRouter()
		ddiscapi.NewHTTPDiscovery(
			q,
			ddiscapi.HTTPDiscoveryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/ddisc/discovery").Subrouter())
		ddiscapi.NewHTTPMedia(
			q,
			ddiscapi.HTTPMediaOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/ddisc/media").Subrouter())
		srv := cmdtestx.NewTLSServer(t, q, routes)

		seederAddrs := torrenttestx.ApprPorts(seeder)
		require.NotEmpty(t, seederAddrs)

		require.NoError(t, cmdtestx.Execute(t, genparser(t), "discovery", "identify",
			"--private-key-path", keypath,
			"--insecure",
			"--library", srv.Listener.Addr().String(),
			"--id", uh.ID,
			"--peer-timeout", "5s",
			"--info-timeout", "10s",
			"--dht-peers", target.DynamicAddrPort().String(),
			"--peer", seederAddrs[0].String(),
		))

		expected := ddisc.NewDiscovered(&id, "")
		var found ddisc.Discovered
		require.NoError(t, ddisc.DiscoveredFindByID(ctx, q, expected.ID).Scan(&found))
		require.NotEmpty(t, found.Mimetype)

		var missing tracking.UnknownHash
		require.Error(t, tracking.UnknownHashDeleteByID(ctx, q, uh.ID).Scan(&missing), "the source lead should be removed once identified")
	})
}
