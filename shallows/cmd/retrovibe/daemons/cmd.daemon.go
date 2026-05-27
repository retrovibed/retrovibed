package daemons

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang-jwt/jwt/v5"
	"github.com/justinas/alice"
	"golang.org/x/crypto/ssh"
	"golang.zx2c4.com/wireguard/tun/netstack"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/retroapi/blockcache"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/netmonx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/dnscache"
	"github.com/retrovibed/retrovibed/shallows/downloads"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/contextx"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/envx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/netx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/internal/tlsx"
	"github.com/retrovibed/retrovibed/shallows/internal/userx"
	"github.com/retrovibed/retrovibed/shallows/internal/wireguardx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/retrovibed/retrovibed/shallows/meta/identityssh"
	"github.com/retrovibed/retrovibed/shallows/metaapi"

	_ "github.com/duckdb/duckdb-go/v2"

	"github.com/james-lawrence/torrent/storage"

	"github.com/gorilla/mux"
	"github.com/logrusorgru/aurora"
)

func DefaultDialer(wgnet *netstack.Net, cache dnscache.Resolver) netx.Dialer {
	if wgnet != nil {
		return dnscache.NewDialer(cache, wgnet)
	}

	return dnscache.NewDialer(cache, &net.Dialer{})
}

func DefaultResolver(d netx.Dialer) *net.Resolver {
	return &net.Resolver{
		PreferGo: true,
		Dial:     d.DialContext,
	}
}

type Command struct {
	Alpha               bool             `flag:"" name:"alpha" help:"enable alpha functionality" default:"false" negatable:"" hidden:"true"`
	AutoMDNS            bool             `flag:"" name:"auto-mdns" help:"enable the multicast dns service" env:"${env_auto_mdns}" default:"true" negatable:""`
	AutoBootstrap       bool             `flag:"" name:"auto-bootstrap" help:"bootstrap from a predefined set of peers" env:"${env_auto_bootstrap}" default:"true" negatable:""`
	AutoDiscovery       bool             `flag:"" name:"auto-discovery" help:"enable automatic discovery of content from peers" env:"${env_auto_discovery}" default:"true" negatable:""`
	AutoIdentifyMedia   bool             `flag:"" name:"auto-identify-media" help:"enable automatically identifying media" env:"${env_auto_identify_media}" default:"true" negatable:""`
	AutoLocateMedia     bool             `flag:"" name:"auto-locate-media" help:"enable automatically locating media from distributed index" env:"${env_auto_locate_media}" default:"true" negatable:""`
	AutoArchive         bool             `flag:"" name:"auto-archive" help:"enable automatic archiving of eligible media" env:"${env_auto_archive}" negatable:""`
	AutoReclaim         bool             `flag:"" name:"auto-reclaim" help:"EXPERIMENTAL: enable automatic reclaiming of disk space of archived media" negatable:"" env:"${env_auto_reclaim}"`
	AutoRecommendations bool             `flag:"" name:"auto-recommendations" help:"enable automatic daily recommendations" default:"true" negatable:""`
	AutoSocks5          bool             `flag:"" name:"auto-socks5" help:"enable the socks5 proxy service" default:"true" negatable:""`
	Socks5              cmdopts.Listener `flag:"" name:"socks5-address" help:"enable socks5 proxy, requires a vpn to be configured" default:"tcp://:9999"`
	DHTLogging          bool             `flag:"" name:"dht-logging" help:"enable debug logging for the dht" default:"false" negatable:"" hidden:"true"`
	TorrentResume       bool             `flag:"" name:"torrent-resume" help:"enable announcing and resuming torrents" default:"true" negatable:""`
	TorrentFirewalled   bool             `flag:"" name:"torrent-firewalled" help:"restrict torrent connections to private networks" env:"${env_torrent_private}"`
	TorrentLogging      bool             `flag:"" name:"torrent-logging" help:"enable torrent logging" default:"false" negatable:"" env:"${env_torrent_logging}"`
	TorrentDebug        bool             `flag:"" name:"torrent-debug" help:"enable torrent debug logging" default:"false" negatable:"" env:"${env_torrent_debug}"`
	TorrentPort         uint16           `flag:"" name:"torrent-port" help:"port to use for torrenting" env:"${env_torrent_port}" default:"10000"`
	TorrentPublicIP4    string           `flag:"" name:"torrent-ipv4" help:"public ipv4 address of the torrent" env:"${env_torrent_ipv4}"`
	TorrentPublicIP6    string           `flag:"" name:"torrent-ipv6" help:"public ipv6 address of the torrent" env:"${env_torrent_ipv6}"`
	TorrentMaxRequests  uint32           `flag:"" name:"torrent-max-outstanding" help:"maximum piece requests to allow" default:"1024"`
	TorrentFolderWatch  []string         `flag:"" name:"torrent-watch" help:"monitor the provided directories for torrent files to automatically download" env:"${env_torrent_directory_watch}" default:"${vars_downloads_directory}"`
	DiscoveryWorkloads  uint64           `flag:"" name:"discovery-workloads" help:"maximum number of infohashes to concurrently process while indexing" default:"1"`
	DiscoveryRatio      uint8            `flag:"" name:"discovery-ratio" help:"percentage of infohashes to index, range from 0-100. 0 = off, 100 = attempt to index every infohash, 1-99 percentage of the partition to index" default:"1"`
	DiscoveryPartitions uint8            `flag:"" name:"discovery-partition" help:"number of partitions to split the infohash space into, adjustments to this value are not recommended as it'll seperate you from identifying synchronization peers" default:"128"`
	DiscoverySeed       string           `flag:"" name:"discovery-seed" help:"seed to generate partition spaces, adjustments to this value are not recommended as it'll seperate you from identifying synchronization peers" default:"retrovibed-ddisc"`
	HTTP                cmdopts.Listener `flag:"" name:"http-address" help:"address to serve daemon api from" default:"tcp://:9998" env:"${env_daemon_socket}"`
	SelfSignedHosts     []string         `flag:"" name:"self-signed-hosts" help:"comma seperated list of hosts to add to the sign signed certificate" env:"${env_self_signed_hosts}"`
}

func (t Command) torrentsettings() *TorrentSettings {
	return AutoTorrentSettings(&TorrentSettings{
		Seed:            envx.Boolean(true, env.TorrentAllowSeeding),
		Pex:             envx.Boolean(true, env.TorrentPEX),
		Log:             envx.Boolean(t.TorrentLogging, env.TorrentLogging),
		Debug:           envx.Boolean(t.TorrentDebug, env.TorrentDebug),
		MaximumRequests: uint64(t.TorrentMaxRequests),
		Ip4:             t.TorrentPublicIP4,
		Ip6:             t.TorrentPublicIP6,
		Resumable:       t.TorrentResume,
		Firewalled:      t.TorrentFirewalled,
		Port:            uint32(t.TorrentPort),
		AutoBootstrap:   t.AutoBootstrap,
		AutoLocateMedia: t.AutoLocateMedia,
	})
}

func (t Command) discoverysettings() *DiscoverySettings {
	return &DiscoverySettings{
		Enabled:    t.AutoDiscovery,
		Ratio:      uint32(t.DiscoveryRatio),
		Partitions: uint32(t.DiscoveryPartitions),
		Workloads:  uint32(t.DiscoveryWorkloads),
		Seed:       t.DiscoverySeed,
	}
}

func (t Command) Run(gctx *cmdopts.Global, sshid *cmdopts.SSHID, tlscfg *cmdopts.TLSConfig) (err error) {
	var (
		db             *sql.DB
		id             ssh.Signer
		_socks5        net.Listener
		deepjwt        = httpx.NewFixedStatusClient(http.StatusMethodNotAllowed)
		mediameta      = asyncx.NewWakeup(gctx.Context)
		archival       = asyncx.NewWakeup(gctx.Context)
		publishing     = asyncx.NewWakeup(gctx.Context)
		vpncfgpath     = userx.DefaultConfigDir(userx.DefaultRelRoot(), "vpn.cfg")
		storagecfgpath = userx.DefaultConfigDir(userx.DefaultRelRoot(), "storage.cfg")
		mc             = library.NewQueryerCleanerAuto()
	)

	gctx.Cleanup.Add(1)
	defer gctx.Cleanup.Done()

	// envx.Debug(os.Environ()...)

	if id, err = sshid.Signer(); err != nil {
		return err
	}

	sshjwt := jwtx.NewSSHSigner()
	jwt.RegisterSigningMethod(sshjwt.Alg(), func() jwt.SigningMethod { return sshjwt })
	jwtx.RegisterAlgorithms(sshjwt, jwt.SigningMethodHS512)

	if c, err := authn.AutoJWTClient(gctx.Context); err == nil {
		deepjwt = c
	} else {
		// we allow creation to fail the application should function even without the api.
		// just warn that the api is unavailable.
		errorsx.Log(errorsx.Wrap(err, "failed to create oauth2, api will fail"))
	}

	httpbind, err := t.HTTP.Socket()
	if err != nil {
		return errorsx.Wrap(err, "unable to setup http socket")
	}

	go func() {
		<-gctx.Context.Done()
		errorsx.Log(errorsx.Wrap(httpbind.Close(), "http server shutdown failed"))
	}()

	if db, err = cmdopts.DatabaseMeta(gctx.Context); err != nil {
		return err
	}
	defer db.Close()

	if err = identityssh.InitializeAdmin(gctx.Context, db, id.PublicKey()); err != nil {
		return errorsx.Wrap(err, "unable to import ssh identity")
	}

	errorsx.Log(errorsx.Wrap(PrepareDefaultFeeds(gctx.Context, db), "unable to initialize default rss feeds"))

	rootstore := fsx.DirVirtual(env.RootStorageDir())
	mediastore := fsx.DirVirtual(env.MediaDir())
	tvfs := fsx.DirVirtual(env.TorrentDir())

	if err := fsx.MkDirs(0700, rootstore.Path(), mediastore.Path(), tvfs.Path(), wireguardx.ConfigDirectory()); err != nil {
		return err
	}

	var tstore storage.ClientImpl = blockcache.NewTorrentFromVirtualFS(tvfs)

	if t.AutoArchive && deepjwt != http.DefaultClient {
		log.Println("automatic archival is enabled")
		errorsx.Log(AutoArchival(gctx.Context, db, mediastore, archival, t.AutoArchive))
		errorsx.Log(AutoPublishing(gctx.Context, db, deepjwt, mediastore, tvfs, publishing))
		errorsx.Log(AutoFeedSync(gctx.Context, db, deepjwt, publishing))
		errorsx.Log(SubscriptionSync(gctx.Context, db, deepjwt))
		tstore = library.NewTorrentStorageFromHTTP(deepjwt, db, tstore)
	} else {
		log.Println("automatic archival is disabled")
	}

	errorsx.Log(AutoReclaim(gctx.Context, db, mediastore, asyncx.NewWakeup(gctx.Context), t.AutoReclaim))

	if t.AutoSocks5 {
		if _socks5, err = t.Socks5.Socket(); err != nil {
			return errorsx.Wrap(err, "unable to setup socks5 socket")
		}
	}

	torrenting := newTorrenting(db, id, rootstore, mediastore, tvfs, mc, tstore, _socks5)

	if err = torrenting.Reload(gctx.Context, t.torrentsettings(), t.discoverysettings()); err != nil {
		return errorsx.Wrap(err, "failed to reload torrent")
	}

	if err = torrenting.Watch(
		gctx.Context,
		vpncfgpath,
		storagecfgpath,
	); err != nil {
		return err
	}

	if netmon := netmonx.Global(); netmon != nil {
		go func() {
			for delta := range netmon.Each(gctx.Context) {
				log.Println("network delta", spew.Sdump(delta))
				torrenting.Broadcast()
			}

			errorsx.Log(errorsx.Wrap(netmon.Err(), "netmon failed"))
		}()
	} else {
		log.Println("network monitor unavailable, network change detection disabled")
	}

	asyncx.Background(gctx.Context, mediameta, func(ctx context.Context) error {
		return errorsx.Wrap(MediaMetadataImport(ctx, db, tvfs, tstore), "media metadata import failed")
	})
	asyncx.Background(gctx.Context, mediameta, func(ctx context.Context) error {
		return errorsx.Wrap(NeuralImport(ctx, db, userx.DefaultCacheDirectory(userx.DefaultRelRoot()), tvfs, tstore), "media metadata import failed")
	})

	go func() {
		errorsx.Log(errorsx.Wrap(asyncx.WatchDirectories(gctx.Context, mediameta, asyncx.FileCreated, mediastore.Path()), "media metadata file watch failed"))
	}()
	mediameta.Broadcast() // trigger an initial pass incase we shutdown midway through.

	go timex.NowAndEveryVoid(gctx.Context, 24*time.Hour, func(ctx context.Context) {
		errorsx.Log(library.NewTombstonedCleanup(ctx, mediastore, db))
	})

	if len(t.TorrentFolderWatch) > 0 {
		dwatcher, err := downloads.NewDirectoryWatcher(gctx.Context, tlsx.MustClone(tlscfg.Config(), tlsx.OptionInsecureSkipVerify), db)
		if err != nil {
			return errorsx.Wrap(err, "unable to setup directory monitoring for torrents")
		}

		for _, dir := range t.TorrentFolderWatch {
			if err = dwatcher.Add(dir); err != nil {
				log.Println(aurora.Yellow("WARNING"), errorsx.Wrapf(err, "unable to watch directory: %s", dir))
				continue
			}
		}
	} else {
		log.Println("download folder monitoring disabled - no watch folders defined")
	}

	go PrintStatistics(gctx.Context, db)

	// block for first refresh
	errorsx.Log(cmdopts.Checkpoint(gctx.Context, db))
	go timex.Every(10*time.Minute, func() {
		errorsx.Log(cmdopts.Checkpoint(gctx.Context, db))
	})

	if t.AutoIdentifyMedia {
		go timex.NowAndEvery(gctx.Context, 15*time.Minute, func(ctx context.Context) error {
			errorsx.Log(IdentifyTorrentMedia(gctx.Context, db))
			return nil
		})
	} else {
		log.Println("auto identify media is disabled, to enable add --auto-identify-media flag, this is an experimental feature.")
	}

	if t.AutoRecommendations {
		errorsx.Log(RecommendationsBackground(gctx.Context, db))
	} else {
		log.Println("auto recommendations is disabled")
	}

	httpmux := mux.NewRouter()
	httpmux.NotFoundHandler = httpx.NotFound(alice.New())
	httpmux.Use(
		httpx.RouteInvoked,
		httpx.Chaos(
			envx.Float64(0.0, env.ChaosRate),
			httpx.ChaosStatusCodes(http.StatusBadGateway),
			httpx.ChaosRateLimited(time.Second),
		),
	)

	httpmux.HandleFunc(
		"/healthz",
		httpx.Healthz(
			cmdopts.MachineID(),
			envx.Float64(0.0, env.HTTPHealthzProbability),
			envx.Int(http.StatusOK, env.HTTPHealthzCode),
		),
	).Methods(http.MethodGet)

	oauth2mux := httpmux.PathPrefix("/oauth2").Subrouter()
	metaapi.NewSSHOauth2(db).Bind(oauth2mux.PathPrefix("/ssh").Subrouter())
	metamux := httpmux.PathPrefix("/meta").Subrouter()

	metaapi.NewHTTP(db).Bind(httpmux.PathPrefix("/sso").Subrouter())
	metaapi.NewHTTPWireguard(wireguardx.ConfigDirectory(), db).Bind(httpmux.PathPrefix("/wireguard").Subrouter())
	metaapi.NewHTTPUsermanagement(db).Bind(metamux.PathPrefix("/u12t").Subrouter())
	metaapi.NewHTTPDaemons(db).Bind(metamux.PathPrefix("/d").Subrouter())
	metaapi.NewHTTPAuthz(db).Bind(metamux.PathPrefix("/authz").Subrouter())
	media.NewHTTPLibrary(
		db,
		archival,
		mediastore,
		deepjwt,
		media.HTTPLibraryOptionTorrentStorage(tvfs),
		media.HTTPLibraryOptionQueryCleaner(mc),
	).Bind(httpmux.PathPrefix("/m").Subrouter())
	media.NewHTTPDiscovered(
		db,
		torrenting._tclient,
		tstore,
		media.HTTPDiscoveredOptionRootStorage(rootstore),
		media.HTTPDiscoveredOptionQueryCleaner(mc),
	).Bind(httpmux.PathPrefix("/d").Subrouter())
	media.NewHTTPRecommendations(db).Bind(httpmux.PathPrefix("/r").Subrouter())
	media.NewHTTPRecent(db).Bind(httpmux.PathPrefix("/w").Subrouter())
	ddiscapi.NewHTTPPeerManagement(db).Bind(httpmux.PathPrefix("/ddisc").Subrouter())
	media.NewHTTPRSSFeed(db).Bind(httpmux.PathPrefix("/rss").Subrouter())
	media.NewHTTPKnown(db).Bind(httpmux.PathPrefix("/k").Subrouter())
	media.NewHTTPLocate(db).Bind(httpmux.PathPrefix("/l").Subrouter())

	metaapi.NewHTTPFileConfig(vpncfgpath).Bind(httpmux.PathPrefix("/s/wireguard").Subrouter())
	metaapi.NewHTTPFileConfig(torrenting.cfgpath).Bind(httpmux.PathPrefix("/s/torrents").Subrouter())
	metaapi.NewHTTPFileConfig(storagecfgpath).Bind(httpmux.PathPrefix("/s/storage").Subrouter())

	community.NewHTTP(
		db,
		envx.Toggle(community.HTTPOptionNoop, community.HTTPOptionHTTPClient(deepjwt), t.AutoArchive),
		community.HTTPOptionPublishing(publishing),
		community.HTTPOptionMediaStorage(mediastore),
		community.HTTPOptionTorrentStorage(tvfs),
	).Bind(httpmux.PathPrefix("/c").Subrouter())

	community.NewHTTPYouTube(db, deepjwt).Bind(httpmux.PathPrefix("/integrations/youtube").Subrouter())

	tlspem := envx.String(userx.DefaultCacheDirectory(userx.DefaultRelRoot(), "tls.pem"), env.DaemonTLSPEM)
	if err = tlsx.SelfSignedLocalHostTLS(tlspem, tlsx.X509OptionHosts(t.SelfSignedHosts...)); err != nil {
		return err
	}

	_ = httpmux.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		if uri, err := route.URLPath(); err == nil {
			log.Println("Route", errorsx.Zero(route.GetPathTemplate()), errorsx.Zero(route.GetMethods()), uri.String())
		}

		return nil
	})

	if t.AutoMDNS {
		if err := MulticastService(gctx.Context, httpbind); err != nil {
			return errorsx.Wrap(err, "unable to setup multicast service")
		}
	} else {
		log.Println("mdns service is disabled")
	}

	log.Println("https listening on:", httpbind.Addr().String(), tlspem)
	if cause := http.ServeTLS(httpbind, httpmux, tlspem, tlspem); netx.IgnoreConnectionClosed(cause) != nil {
		return errorsx.Wrap(cause, "http server stopped")
	}

	// report any async failures.
	return contextx.IgnoreCancelled(errorsx.Compact(context.Cause(gctx.Context), gctx.Context.Err()))
}
