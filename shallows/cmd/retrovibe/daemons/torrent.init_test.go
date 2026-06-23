package daemons

import (
	"context"
	"database/sql"
	"io"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/autobind"
	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/james-lawrence/torrent/storage"
	"github.com/james-lawrence/torrent/torrenttest"
	"github.com/retrovibed/retrovibed/retroapi/blockcache"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/dnscache"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sshx"
	"github.com/retrovibed/retrovibed/shallows/internal/torrenttestx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
	"golang.zx2c4.com/wireguard/device"
)

func newTestTorrenting(t *testing.T, q *sql.DB) _torrenting {
	t.Helper()
	signer, err := sshx.SignerFromGenerator(sshx.UnsafeNewKeyGen())
	require.NoError(t, err)

	rootstore := fsx.DirVirtual(t.TempDir())
	require.NoError(t, fsx.MkDirs(0700, rootstore.Path("torrent")))
	vfs := fsx.DirVirtual(t.TempDir())
	tstore := blockcache.NewTorrentFromVirtualFS(vfs)

	return _torrenting{
		cond:             sync.NewCond(&sync.Mutex{}),
		cfgpath:          rootstore.Path("torrent.cfg"),
		discoverycfgpath: rootstore.Path("discovery.cfg"),
		ddiscidpath:      rootstore.Path("ddisc.id"),
		peercachepath:    rootstore.Path("torrent.peers"),
		machineid:        "test-machine-id",
		wgconfigdir:      t.TempDir(),
		wglatest:         rootstore.Path("wireguard.latest"),
		db:               q,
		id:               signer,
		rootstore:        rootstore,
		mediastore:       vfs,
		tvfs:             vfs,
		tstore:           tstore,
		_tclient:         &atomic.Pointer[torrent.Client]{},
		_dnscache:        dnscache.AutoProxyResolver(),
		_wgdev:           &atomic.Pointer[device.Device]{},
		_dhts:            &atomic.Pointer[dht.Server]{},
		_discovery:       &atomic.Pointer[ddisc.Snapshot]{},
	}
}

func TestInit(t *testing.T) {
	t.Run("resumes in-progress torrent after repeated Init calls", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)

		signer, err := sshx.SignerFromGenerator(sshx.UnsafeNewKeyGen())
		require.NoError(t, err)

		rootstore := fsx.DirVirtual(t.TempDir())
		require.NoError(t, fsx.MkDirs(0700, rootstore.Path("torrent")))
		vfs := fsx.DirVirtual(t.TempDir())
		tstore := blockcache.NewTorrentFromVirtualFS(vfs)

		seedir := t.TempDir()
		mcache := torrent.NewMetadataCache(seedir)
		info, _, err := torrenttest.Random(seedir, 128*bytesx.KiB)
		require.NoError(t, err)
		torrentMD, err := torrent.NewFromInfo(info, torrent.OptionStorage(storage.NewFile(seedir)))
		require.NoError(t, err)
		require.NoError(t, mcache.Write(torrentMD))

		seeder := torrenttestx.Client(t, autobind.NewLoopback(autobind.EnableDHT(torrenttestx.QuickDHT(t))), mcache, storage.NewFile(seedir))
		defer seeder.Close()
		_, _, err = seeder.Start(torrentMD)
		require.NoError(t, err)

		infoBytes := testx.Must(metainfo.Encode(info))(t)
		infohash := int160.New(infoBytes)

		// insert as initiated and incomplete
		md := tracking.NewMetadata(
			&infohash,
			tracking.MetadataOptionInitiate,
			tracking.MetadataOptionBytes(info.TotalLength()),
			tracking.MetadataOptionDownloaded(0),
		)
		require.NoError(t, tracking.MetadataInsertWithDefaults(t.Context(), q, md).Scan(&md))

		tr := _torrenting{
			cond:          sync.NewCond(&sync.Mutex{}),
			cfgpath:       rootstore.Path("torrent.cfg"),
			ddiscidpath:   rootstore.Path("ddisc.id"),
			peercachepath: rootstore.Path("torrent.peers"),
			machineid:     "test-machine-id",
			wgconfigdir:   t.TempDir(),
			wglatest:      rootstore.Path("wireguard.latest"),
			db:            q,
			id:            signer,
			rootstore:     rootstore,
			mediastore:    vfs,
			tvfs:          vfs,
			tstore:        tstore,
			_tclient:      &atomic.Pointer[torrent.Client]{},
			_dnscache:     dnscache.AutoProxyResolver(),
			_wgdev:        &atomic.Pointer[device.Device]{},
			_dhts:         &atomic.Pointer[dht.Server]{},
			_discovery:    &atomic.Pointer[ddisc.Snapshot]{},
		}

		cfg := AutoTorrentSettings(&TorrentSettings{
			Resumable: true,
		})
		disc := &DiscoverySettings{}

		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()

		noop := func(e error) {
			require.NoError(t, err)
		}

		startAndRead := func(c *torrent.Client) {
			dl, _, err := c.Start(torrentMD, torrent.TuneClientPeer(seeder), torrent.TuneNewConns)
			require.NoError(t, err)
			r := torrent.NewReader(dl)
			defer r.Close()
			buf := make([]byte, 1)
			_, err = io.ReadFull(r, buf)
			require.NoError(t, err)
		}

		// first Init - can start and download
		require.NoError(t, tr.Init(ctx, noop, cfg, disc))
		startAndRead(tr._tclient.Load())

		// second Init (simulates reload) - new client, torrent resumes
		require.NoError(t, tr.Init(ctx, noop, cfg, disc))
		startAndRead(tr._tclient.Load())
	})

	t.Run("creates a valid client when Resumable is false", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)
		tr := newTestTorrenting(t, q)

		cfg := AutoTorrentSettings(&TorrentSettings{Resumable: false})
		disc := &DiscoverySettings{}

		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()

		noop := func(e error) { require.NoError(t, e) }

		require.NoError(t, tr.Init(ctx, noop, cfg, disc))
		require.NotNil(t, tr._tclient.Load())
	})

	t.Run("re-init replaces the previous client", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)
		tr := newTestTorrenting(t, q)

		cfg := AutoTorrentSettings(&TorrentSettings{})
		disc := &DiscoverySettings{}

		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()

		noop := func(e error) { require.NoError(t, e) }

		require.NoError(t, tr.Init(ctx, noop, cfg, disc))
		first := tr._tclient.Load()
		require.NotNil(t, first)

		require.NoError(t, tr.Init(ctx, noop, cfg, disc))
		second := tr._tclient.Load()
		require.NotNil(t, second)
		require.NotSame(t, first, second)
	})

	t.Run("creates a valid client when Firewalled is true", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)
		tr := newTestTorrenting(t, q)

		cfg := AutoTorrentSettings(&TorrentSettings{Firewalled: true})
		disc := &DiscoverySettings{}

		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()

		noop := func(e error) { require.NoError(t, e) }

		require.NoError(t, tr.Init(ctx, noop, cfg, disc))
		require.NotNil(t, tr._tclient.Load())
	})
}
