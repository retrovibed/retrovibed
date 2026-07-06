package daemons

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"math"
	"net"
	"net/netip"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/fsnotify/fsnotify"
	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/bep0051"
	"github.com/james-lawrence/torrent/connections"
	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/dht/krpc"
	"github.com/james-lawrence/torrent/storage"
	"github.com/retrovibed/retrovibed/retroapi/netmonx"
	retronetx "github.com/retrovibed/retrovibed/retroapi/netx"
	"github.com/retrovibed/retrovibed/retroapi/userx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/ddisc/ddisctorrent"
	"github.com/retrovibed/retrovibed/shallows/dnscache"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/contextx"
	"github.com/retrovibed/retrovibed/shallows/internal/cryptox"
	"github.com/retrovibed/retrovibed/shallows/internal/dhtx"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/envx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/internal/torrentx"
	"github.com/retrovibed/retrovibed/shallows/internal/wireguardx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"golang.org/x/crypto/ssh"
	"golang.org/x/time/rate"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun/netstack"
	"google.golang.org/protobuf/proto"
)

func AutoTorrentSettings(defaults *TorrentSettings, options ...func(*TorrentSettings)) *TorrentSettings {
	// highly conservative settings for vpn providers like protonvpn.
	return new(langx.Clone(TorrentSettings{
		Seed:            defaults.Seed,
		Pex:             defaults.Pex,
		Log:             defaults.Log,
		Debug:           defaults.Debug,
		MaximumRequests: defaults.MaximumRequests,
		Ip4:             defaults.Ip4,
		Ip6:             defaults.Ip6,
		Resumable:       defaults.Resumable,
		Firewalled:      defaults.Firewalled,
		Port:            defaults.Port,
		AutoBootstrap:   defaults.AutoBootstrap,
		AutoLocateMedia: defaults.AutoLocateMedia,
		Peers:           &Peers{Min: 16, Max: 64},
		Upload:          &Limit{Rate: 128 * bytesx.MiB, Burst: 128 * bytesx.MiB},
		Download:        &Limit{Rate: 128 * bytesx.MiB, Burst: 128 * bytesx.MiB},
		Outbound:        &Limit{Rate: 4, Burst: 1},
		Inbound:         &Limit{Rate: 4, Burst: 1},
	}, options...))
}

func newTorrenting(db *sql.DB, id ssh.Signer, root, media, tvfs fsx.Virtual, mc library.QueryCleaner, tstore storage.ClientImpl, socks5 net.Listener) _torrenting {
	return _torrenting{
		cond:             sync.NewCond(&sync.Mutex{}),
		cfgpath:          userx.DefaultConfigDir(userx.DefaultRelRoot(), "torrent.cfg"),
		discoverycfgpath: userx.DefaultConfigDir(userx.DefaultRelRoot(), "discovery.cfg"),
		ddiscidpath:      userx.DefaultConfigDir(userx.DefaultRelRoot(), "ddisc.id"),
		peercachepath:    userx.DefaultCacheDirectory(userx.DefaultRelRoot(), "torrent.peers"),
		samplecachepath:  userx.DefaultCacheDirectory(userx.DefaultRelRoot(), "torrent.bep51.samples"),
		machineid:        cmdopts.MachineID(),
		wgconfigdir:      wireguardx.ConfigDirectory(),
		wglatest:         wireguardx.Latest(),
		db:               db,
		id:               id,
		rootstore:        root,
		mediastore:       media,
		tvfs:             tvfs,
		tstore:           tstore,
		socks5:           socks5,
		mc:               mc,
		_tclient:         &atomic.Pointer[torrent.Client]{},
		_dnscache:        dnscache.AutoProxyResolver(),
		_wgdev:           &atomic.Pointer[device.Device]{},
		_dhts:            &atomic.Pointer[dht.Server]{},
		_discovery:       &atomic.Pointer[ddisc.Snapshot]{},
	}
}

type _torrenting struct {
	cfgpath          string
	discoverycfgpath string
	ddiscidpath      string
	peercachepath    string
	samplecachepath  string
	machineid        string
	wgconfigdir      string
	wglatest         string
	db               *sql.DB
	id               ssh.Signer
	rootstore        fsx.Virtual
	mediastore       fsx.Virtual
	tvfs             fsx.Virtual
	mc               library.QueryCleaner
	tstore           storage.ClientImpl
	socks5           net.Listener
	_tclient         *atomic.Pointer[torrent.Client]
	_dnscache        *dnscache.ProxyPtr
	_wgdev           *atomic.Pointer[device.Device]
	_dhts            *atomic.Pointer[dht.Server]
	_discovery       *atomic.Pointer[ddisc.Snapshot]
	cond             *sync.Cond
}

func (t _torrenting) loadcfg(path string, v proto.Message) error {
	encoded, err := fsx.AutoCached(path, func() ([]byte, error) {
		return json.Marshal(v)
	})
	if err != nil {
		return err
	}

	var (
		d = proto.Clone(v)
	)

	if err = json.Unmarshal(encoded, d); err != nil {
		return err
	}

	proto.Merge(v, d)

	return nil
}

func (t *_torrenting) WireguardSnapshot() (wireguardx.Statistics, error) {
	if dev := t._wgdev.Load(); dev != nil {
		return wireguardx.Snapshot(dev)
	}
	return wireguardx.Statistics{}, nil
}

func (t *_torrenting) DHTSnapshot() (dht.ServerStats, error) {
	if d := t._dhts.Load(); d != nil {
		return d.Stats(), nil
	}
	return dht.ServerStats{}, nil
}

func (t *_torrenting) DiscoverySnapshot() (ddisc.Snapshot, error) {
	if d := t._discovery.Load(); d != nil {
		return *d, nil
	}
	return ddisc.Snapshot{}, nil
}

func (t *_torrenting) Reload(ctx context.Context, cfg *TorrentSettings, disc *DiscoverySettings) error {
	go func() {
		limiter := rate.NewLimiter(rate.Every(5*time.Second), 1)
		for {
			var (
				mcfg  = proto.CloneOf(cfg)
				mdisc = proto.CloneOf(disc)
			)

			log.Println("torrent settings initiated", spew.Sdump(mcfg))
			if err := t.loadcfg(t.cfgpath, mcfg); err != nil {
				errorsx.Log(errorsx.Wrap(err, "failed to read torrent config"))
				continue
			}
			log.Println("torrent settings completed", spew.Sdump(mcfg))

			log.Println("discovery settings initiated", spew.Sdump(mdisc))
			if err := t.loadcfg(t.discoverycfgpath, mdisc); err != nil {
				errorsx.Log(errorsx.Wrap(err, "failed to read discovery config"))
				continue
			}
			log.Println("discovery settings completed", spew.Sdump(mdisc))

			if metered := netmonx.Metered(); metered {
				log.Println("applying metered settings to configuration")
				mcfg.Inbound.Rate = 0         // block inbound connections.
				mcfg.Upload.Rate = bytesx.KiB // dramatically slow down distribution to peers
				mcfg.Seed = false             // dont seed
				mcfg.Resumable = false        // dont attempt to resume downloads
				mcfg.AutoLocateMedia = false  // dont attempt to index the swarm
				mdisc.Enabled = false         // dont attempt to discover content from the swarm
			}

			_ctx, _done := context.WithCancelCause(ctx)
			log.Println("torrent settings", spew.Sdump(mcfg))
			log.Println("discovery settings", spew.Sdump(mdisc))
			asyncfailure := func(cause error) {
				defer t.cond.Broadcast()
				errorsx.Log(errorsx.Wrap(cause, "async failure"))
				_done(contextx.IgnoreCancelled(cause))
			}

			if err := t.Init(_ctx, asyncfailure, mcfg, mdisc); err != nil {
				_done(err)
				errorsx.Log(errorsx.Wrap(err, "reloading torrent failed"))
				if err := limiter.Wait(ctx); err != nil {
					return
				}
				continue
			}

			log.Println("torrent client initialized, waiting for reload signal")
			t.cond.L.Lock()
			t.cond.Wait()
			log.Println("reload signal received")
			_done(nil)
			t.cond.L.Unlock()

			select {
			case <-ctx.Done():
				errorsx.Log(context.Cause(ctx))
				return
			default:
			}

			log.Println("reloading torrent client")
		}
	}()

	return nil
}

func (t *_torrenting) Broadcast() {
	t.cond.Broadcast()
}

func (t *_torrenting) Watch(ctx context.Context, paths ...string) error {
	if err := fsx.Touch(0600, paths...); err != nil {
		return err
	}

	if err := fsx.Touch(0600, t.wglatest); err != nil {
		return err
	}

	if err := t.loadcfg(t.cfgpath, &TorrentSettings{}); err != nil {
		return err
	}

	if err := t.loadcfg(t.discoverycfgpath, &DiscoverySettings{}); err != nil {
		return err
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	addpath := func(path string) {
		if err = w.Add(path); err != nil {
			errorsx.Log(errorsx.Wrapf(err, "unable to watch %s", path))
			return
		}
	}

	addpath(t.cfgpath)
	addpath(t.discoverycfgpath)
	addpath(t.wgconfigdir)

	for _, path := range paths {
		addpath(path)
	}

	go func() {
		defer log.Println("torrent file watch done")
		defer w.Close()
		for {
			select {
			case evt := <-w.Events:
				log.Println("resetting torrent client due to filesystem event", evt.Op, evt.Name)
				t.cond.Broadcast()
			case err := <-w.Errors:
				log.Println("watch error", err)
			case <-ctx.Done():
				log.Println("context completed", ctx.Err())
				return
			}
		}
	}()

	return nil
}

func (t *_torrenting) Init(dctx context.Context, asyncfailure context.CancelCauseFunc, cfg *TorrentSettings, disc *DiscoverySettings) (err error) {
	const (
		dhtk        = 128
		dhtminpeers = dhtk / 4
	)

	var (
		tnetwork     torrent.Binder
		wgnet        *netstack.Net
		firewall     torrent.ClientConfigOption
		wgcfg        meta.Wireguard
		dhts         *dht.Server
		limit                   = retronetx.NewConnUnlimited()
		peerid                  = int160.New(ssh.FingerprintSHA256(t.id.PublicKey()))
		torrentpeers            = t.peercachepath
		bootstrap    dht.Option = dht.OptionNoop
		info                    = torrent.ClientConfigInfoLogger(torrent.LogDiscard())
		debug                   = torrent.ClientConfigDebugLogger(torrent.LogDiscard())
		dhtdebug                = dht.OptionLogger(torrent.LogDiscard())
	)

	ddisckey, err := fsx.AutoCached(t.ddiscidpath, func() (_ []byte, _ error) {
		return md5x.Digest(t.machineid).Sum(nil), nil
	})
	if err != nil {
		return errorsx.Wrap(err, "unable to generate ddisc key")
	}

	partitions := ddisc.Partitions(uint16(disc.Partitions), cryptox.NewChaCha8(disc.Seed))
	localpartition := partitions.Max(ddisckey)
	t._discovery.Store(&ddisc.Snapshot{
		Enabled:        disc.Enabled,
		Ratio:          disc.Ratio,
		Partitions:     disc.Partitions,
		Workloads:      disc.Workloads,
		LocalPartition: localpartition.String(),
	})

	log.Println("dht peer id", peerid.String())
	log.Println("ddisc partitions digest", disc.Partitions, disc.Seed, "->", ddisc.PartitionsDigest(partitions).String())
	log.Println("ddisc partition", localpartition)

	if c := t._tclient.Load(); c != nil {
		errorsx.Log(errorsx.Wrap(c.Close(), "failed to close previous client"))
	}

	log.Println("------------------------------------------------------- initiated torrent initialization -------------------------------------------------------")
	defer log.Println("------------------------------------------------------- completed torrent initialization -------------------------------------------------------")

	if cfg.Log {
		info = torrent.ClientConfigInfoLogger(log.New(os.Stderr, "[torrent] ", log.Flags()))
	}

	if cfg.Debug {
		debug = torrent.ClientConfigDebugLogger(log.New(os.Stderr, "[torrent - debug] ", log.Flags()))
	}

	if envx.Boolean(false, env.DHTDebug) {
		dhtdebug = dht.OptionLogger(log.New(os.Stderr, "[dht] ", log.Flags()))
	}

	if fsx.IsRegularFile(torrentpeers) {
		log.Println("file cache bootstrap enabled", torrentpeers)
		bootstrap = dht.OptionBootstrapPeerFile(torrentpeers)
	}

	if cfg.AutoBootstrap {
		log.Println("public bootstrap enabled")
		bootstrap = langx.Compose(bootstrap, dht.OptionBootstrapGlobal)
	}

	if err = meta.WireguardCurrent(dctx, t.db).Scan(&wgcfg); errorsx.Ignore(err, sql.ErrNoRows) != nil {
		return errorsx.Wrap(err, "failed to read wireguard config")
	}

	tmux := dht.DefaultMuxer().
		Method(bep0051.Query, bep0051.NewEndpoint(NewSampler(t.db, time.Duration(bep0051.TTLMax)*time.Second, t.samplecachepath))).
		Method(ddisctorrent.MethodMeta, ddisctorrent.NewMeta(localpartition)).
		Method(ddisctorrent.MethodSearch, ddisctorrent.NewSearch(t.db)).
		Method(ddisctorrent.MethodSync, ddisctorrent.NewSync(t.db)).
		Method(ddisctorrent.MethodDisc, ddisctorrent.NewDiscovered(t.db, partitions)).
		Method(ddisctorrent.MethodMedia, ddisctorrent.NewMediaRecorder(t.db))

	if path := fsx.DirVirtual(t.wgconfigdir).Path(uuid.FromStringOrNil(wgcfg.ID).String()); fsx.Exists(path) {
		wcfg, err := wireguardx.Parse(path)
		if err != nil {
			return errorsx.Wrapf(err, "unable to parse wireguard configuration: %s", path)
		}

		log.Println("loaded wireguard configuration", path)

		var wgdev *device.Device
		if wgnet, wgdev, err = torrentx.WireguardSocket(dctx, wcfg); err != nil {
			return errorsx.Wrap(err, "unable to setup wireguard tunnel")
		}
		t._wgdev.Store(wgdev)

		log.Println("------------------------------------", cfg.Port, wgcfg.Port, "------------------------------------")
		t._dnscache.Store(
			dnscache.New(
				wireguardx.HostLookupAdapter(wgnet),
				dnscache.CacheOptionRateLimit(wgcfg.DNSRateLimit),
			),
		)

		dhts, err = dht.NewServer(
			dhtk,
			dht.OptionNodeID(peerid),
			dht.OptionMuxer(tmux),
			dht.OptionHostResolver(t._dnscache),
			torrentx.AutomaticIP(wcfg, wgnet, wgcfg.Port),
			dhtdebug,
			bootstrap,
		)
		if err != nil {
			return errorsx.Wrap(err, "unable to setup dht server")
		}
		t._dhts.Store(dhts)

		limit = retronetx.NewConnLimited(langx.FirstNonZero(wgcfg.MaximumConnections, math.MaxUint64))
		tnetwork, err = torrentx.SetupTorrentBinder(
			wgnet,
			uint16(cfg.Port),
			limit,
			torrent.BinderOptionDHT(dhts),
		)
		if err != nil {
			return errorsx.Wrap(err, "unable to setup torrent binder")
		}
	} else {
		t._dnscache.Store(dnscache.New(net.DefaultResolver))

		log.Println("no wireguard configuration found at", path, wgnet == nil)

		dhts, err = dht.NewServer(
			dhtk,
			dht.OptionNodeID(peerid),
			dht.OptionMuxer(tmux),
			dht.OptionHostResolver(t._dnscache),
			dht.OptionOnQuery(func(source netip.AddrPort, query *krpc.Msg) (propagate bool) {
				const samplerate = 0.01
				ctx, done := context.WithTimeout(context.Background(), time.Second)
				defer done()
				errorsx.Log(errorsx.Wrap(tracking.SamplePeer(ctx, t.db, samplerate, query.A.ID.Int160(), source), "unable to sample peer"))
				return true
			}),
			dht.OptionUPnP,
			dhtdebug,
			bootstrap,
		)
		if err != nil {
			return errorsx.Wrap(err, "unable to setup dht server")
		}
		t._dhts.Store(dhts)

		log.Println("------------------------------------", cfg.Port, wgcfg.Port, "------------------------------------")
		if tnetwork, err = torrentx.Autosocket(dhts, uint16(cfg.Port), limit); err != nil {
			return errorsx.Wrap(err, "unable to setup torrent socket")
		}
	}

	// debugx.OnSignal(dctx, dhtx.Statistics(dhts), syscall.SIGUSR1)
	// go dhtx.BackgroundStatistics(dctx, time.Minute, dhts)
	// go timex.NowAndEveryVoid(dctx, 5*time.Second, func(_ context.Context) {
	// 	retronetx.ConnLimitStatistics(os.Stderr, limit)
	// })
	go dhtx.RecordBootstrapNodes(dctx, time.Minute, dhtminpeers, dhts, torrentpeers)
	go dhts.TableMaintainer(dctx)

	firewall = torrent.ClientConfigFirewall(connections.NewFirewall(
		connections.BanInvalidPort{},
	))
	if cfg.Firewalled {
		log.Println("disabling public networks for torrent")
		firewall = torrent.ClientConfigFirewall(connections.NewFirewall(
			connections.Private{},
			connections.BanInvalidPort{},
		))
	}

	log.Printf("USING STORAGE %T - %s\n", t.tstore, t.tvfs.Path())

	torconfig := torrent.NewDefaultClientConfig(
		torrent.NewMetadataCache(t.tvfs.Path()),
		t.tstore,
		torrent.ClientConfigCacheDirectory(t.tvfs.Path()),
		torrent.ClientConfigPEX(cfg.Pex),
		torrent.ClientConfigSeed(cfg.Seed),
		torrent.ClientConfigDialer(DefaultDialer(wgnet, t._dnscache)),
		torrent.ClientConfigDialTimeouts(time.Second, 4*time.Second),
		torrent.ClientConfigHandshakeTimeout(30*time.Second),
		torrent.ClientConfigDialPoolSize(128*runtime.NumCPU()),
		torrent.ClientConfigDialRateLimit(rate.NewLimiter(rate.Limit(cfg.Outbound.Rate), int(cfg.Outbound.Burst))),
		torrent.ClientConfigAcceptLimit(rate.NewLimiter(rate.Limit(cfg.Inbound.Rate), int(cfg.Inbound.Burst))),
		torrent.ClientConfigMaxOutstandingRequests(int(cfg.MaximumRequests)),
		torrent.ClientConfigPeerLimits(cfg.Peers.Min, cfg.Peers.Max),
		torrent.ClientConfigUploadLimit(rate.NewLimiter(rate.Limit(cfg.Upload.Rate), int(cfg.Upload.Burst))),
		torrent.ClientConfigDownloadLimit(rate.NewLimiter(rate.Limit(cfg.Download.Rate), int(cfg.Download.Burst))),
		torrent.ClientConfigHTTPUserAgent("retrovibed/0.0"),
		torrent.ClientConfigExtension(ddisctorrent.ExtensionName),
		torrent.ClientConfigConnectionClosed(func(ih int160.T, stats torrent.ConnStats, remaining int) {
			if stats.BytesWrittenData.Uint64() == 0 {
				return
			}

			var md tracking.Metadata
			ictx, done := context.WithTimeout(context.Background(), 3*time.Second)
			defer done()
			if err := tracking.MetadataUploadedByID(ictx, t.db, ih.Bytes(), stats.BytesWrittenData.Uint64()).Scan(&md); errorsx.Ignore(err, sql.ErrNoRows) != nil {
				log.Println(errorsx.Wrapf(err, "%s: unable to record uploaded metrics", ih.String()))
				return
			}
		}),
		info,
		debug,
		firewall,
	)

	tclient, err := tnetwork.Bind(torrent.NewClient(torconfig))
	if err != nil {
		return errorsx.Wrap(err, "unable to setup torrent client")
	}

	// TODO: AutoLocateMedia should be located within distributed indexing.
	if cfg.AutoLocateMedia {
		go timex.NowAndEvery(dctx, 15*time.Minute, func(ctx context.Context) error {
			errorsx.Log(LocateMedia(dctx, t.db, tclient, disc))
			return nil
		})
	} else {
		log.Println("auto locate media is disabled, to enable add --auto-locate-media flag.")
	}

	if disc.Enabled && disc.Ratio > 0 {
		go dhtx.WaitForMinimumNodes(dctx, 32, dhts, func() {
			go func() {
				// finds infohashes from the dht
				if err := AutoDiscovery(dctx, t.db, dhts, t.tstore); err != nil {
					asyncfailure(errorsx.Wrap(err, "autodiscovery from peers failed"))
					return
				}
			}()

			go func() {
				// extracts media metadata we've been able to download the metadata for.
				if err := DiscoverMedia(dctx, t.db, dhts, tclient); err != nil {
					asyncfailure(errorsx.Wrap(err, "index discovered media failed"))
					return
				}
			}()

			go func() {
				if err := ddisctorrent.Announce(dctx, tclient, dhts, localpartition); err != nil {
					asyncfailure(errorsx.Wrap(err, "failed to announce partition(s)"))
					return
				}
			}()

			go func() {
				log.Println("auto retrieval of unknown infohashes initiated")
				defer log.Println("auto retrieval of unknown infohashes completed")

				var (
					block ddisc.Filter = ddisc.FilterNone
				)

				if disc.Ratio < 100 {
					n := partitions.Max(peerid.Bytes())
					block = ddisc.FilterRatio(cryptox.NewChaCha8(n[:]), uint8(disc.Ratio)).Filter
				}

				if err := DiscoverDHTMetadata(dctx, uint64(disc.Workloads), t.db, tclient, block); err != nil {
					asyncfailure(errorsx.Wrap(err, "infohash sampling failed"))
					return
				}
			}()
		})
		log.Println("autodiscovery is enabled. this is an experimental feature.")
	} else if disc.Ratio == 0 {
		log.Println("autodiscovery is disabled, to enable remove --discovery-ratio=0 flag, this is an experimental feature.")
	} else {
		log.Println("autodiscovery is disabled, due to no dht servers. this is an experimental feature.")
	}

	go func() {
		if err := DiscoverFromRSSFeeds(dctx, t.db, t.rootstore, t.mc, tclient, t.tstore); errorsx.Ignore(err, context.Canceled) != nil {
			asyncfailure(errorsx.Wrap(err, "autodiscovery of RSS feeds failed"))
			return
		}
	}()

	go timex.NowAndEveryVoid(dctx, envx.Duration(24*time.Hour, env.TorrentVerifyFrequency), func(ctx context.Context) {
		VerifyTorrents(dctx, t.db, t.rootstore, t.mc, tclient, t.tstore)
	})

	if cfg.Resumable {
		go AnnounceSeeded(dctx, t.db, dhts, t.rootstore, tclient, t.tstore)
		go ResumeDownloads(dctx, t.db, t.rootstore, t.mc, tclient, t.tstore)
	} else {
		log.Println("announce/resume disabled")
	}

	go torrentx.ClearIdleTorrents(dctx, time.Hour, tclient)

	if t.socks5 == nil || wgnet == nil {
		log.Printf("socks5 service is disabled - enabled(%t) vpn(%t)\n", t.socks5 != nil, wgnet != nil)
	} else {
		if err = t.socks5.Close(); err != nil {
			log.Printf("socks5 socket closed, reopening - %T - %v\n", err, err)
		}

		if t.socks5, err = net.Listen(t.socks5.Addr().Network(), t.socks5.Addr().String()); err != nil {
			return errorsx.Wrap(err, "unable to serve socks5 - listen")
		}

		if err := Socks5(dctx, wgnet, WireguardResolver{Resolver: t._dnscache}, t.socks5); err != nil {
			return errorsx.Wrap(err, "unable to serve socks5 - serve")
		}
	}

	t._tclient.Store(tclient)

	return nil
}
