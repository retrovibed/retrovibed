package daemons

import (
	"context"
	"database/sql"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/justinas/alice"
	"golang.org/x/crypto/ssh"

	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/dht/krpc"
	"github.com/james-lawrence/torrent/storage"
	"github.com/retrovibed/retrovibed/cmd/cmdmeta"
	"github.com/retrovibed/retrovibed/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/downloads"
	"github.com/retrovibed/retrovibed/internal/contextx"
	"github.com/retrovibed/retrovibed/internal/dhtx"
	"github.com/retrovibed/retrovibed/internal/env"
	"github.com/retrovibed/retrovibed/internal/envx"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/fsx"
	"github.com/retrovibed/retrovibed/internal/httpx"
	"github.com/retrovibed/retrovibed/internal/slicesx"
	"github.com/retrovibed/retrovibed/internal/timex"
	"github.com/retrovibed/retrovibed/internal/tlsx"
	"github.com/retrovibed/retrovibed/internal/torrentx"
	"github.com/retrovibed/retrovibed/internal/userx"
	"github.com/retrovibed/retrovibed/media"

	_ "github.com/marcboeker/go-duckdb"

	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/bep0051"

	"github.com/gorilla/mux"
)

type Command struct {
	DisableMDNS   bool             `flag:"" name:"no-mdns" help:"disable the multicast dns service" default:"false" env:"${env_mdns_enabled}"`
	AutoBootstrap bool             `flag:"" name:"auto-bootstrap" help:"bootstrap from a predefined set of peers" default:"false" env:"${env_auto_bootstrap}"`
	AutoDiscovery bool             `flag:"" name:"auto-discovery" help:"enable autodiscovery of content from peers" default:"false" env:"${env_auto_discovery}"`
	HTTP          cmdopts.Listener `flag:"" name:"http-address" help:"address to serve daemon api from" default:"tcp://:9998"`
}

func (t Command) Run(gctx *cmdopts.Global, id *cmdopts.SSHID) (err error) {
	var (
		db           *sql.DB
		torrentpeers                            = userx.DefaultCacheDirectory(userx.DefaultRelRoot(), "torrent.peers")
		peerid                                  = krpc.IdFromString(ssh.FingerprintSHA256(id.PublicKey()))
		bootstrap    torrent.ClientConfigOption = torrent.ClientConfigNoop
	)

	// envx.Debug(os.Environ()...)

	dctx, done := context.WithCancelCause(gctx.Context)
	asyncfailure := func(cause error) {
		done(contextx.IgnoreCancelled(cause))
	}
	defer asyncfailure(nil)

	httpbind, err := t.HTTP.Socket()
	if err != nil {
		return errorsx.Wrap(err, "unable to setup http socket")
	}
	go func() {
		<-dctx.Done()
		httpbind.Close()
	}()

	if db, err = cmdmeta.Database(dctx); err != nil {
		return err
	}
	defer db.Close()

	go func() {
		errorsx.Log(errorsx.Wrap(PrepareDefaultFeeds(dctx, db), "unable to initialize default rss feeds"))
	}()

	tnetwork, err := torrentx.Autosocket(0)
	if err != nil {
		return errorsx.Wrap(err, "unable to setup torrent socket")
	}

	if fsx.IsRegularFile(torrentpeers) {
		bootstrap = torrent.ClientConfigBootstrapPeerFile(torrentpeers)
	}

	if t.AutoBootstrap {
		bootstrap = torrent.ClientConfigBootstrapGlobal
	}

	tm := dht.DefaultMuxer().
		Method(bep0051.Query, bep0051.NewEndpoint(bep0051.EmptySampler{}))
	tclient, err := tnetwork.Bind(
		torrent.NewClient(
			torrent.NewDefaultClientConfig(
				torrent.ClientConfigPeerID(string(peerid[:])),
				torrent.ClientConfigSeed(true),
				torrent.ClientConfigInfoLogger(log.New(io.Discard, "[torrent] ", log.Flags())),
				torrent.ClientConfigMuxer(tm),
				torrent.ClientConfigBucketLimit(32),
				bootstrap,
			),
		),
	)
	if err != nil {
		return errorsx.Wrap(err, "unable to setup torrent client")
	}

	rootstore := fsx.DirVirtual(userx.DefaultDataDirectory(userx.DefaultRelRoot()))
	mediastore := fsx.DirVirtual(env.MediaDir())
	// tstore := blockcache.NewTorrentFromVirtualFS(mediastore)
	torrentdir := env.TorrentDir()
	tstore := storage.NewFileByInfoHash(torrentdir)

	dwatcher, err := downloads.NewDirectoryWatcher(dctx, db, rootstore, tclient, tstore)
	if err != nil {
		return errorsx.Wrap(err, "unable to setup directory monitoring for torrents")
	}

	if err = dwatcher.Add(userx.DefaultDownloadDirectory()); err != nil {
		return errorsx.Wrap(err, "unable to add the download directory to be watched")
	}

	{
		current, _ := slicesx.First(tclient.DhtServers()...)
		if peers, err := current.AddNodesFromFile(torrentpeers); err == nil {
			log.Println("added peers", peers)
		} else {
			log.Println("unable to read peers", err)
		}
	}

	for _, d := range tclient.DhtServers() {
		go dhtx.RecordBootstrapNodes(dctx, time.Minute, d, torrentpeers)
		go d.TableMaintainer()
		go d.Bootstrap(dctx)
	}

	go PrintStatistics(dctx, db)

	go timex.NowAndEvery(gctx.Context, time.Minute, func(ctx context.Context) error {
		_, err := db.ExecContext(ctx, "PRAGMA create_fts_index('library_metadata', 'id', 'description', overwrite = 1);")
		if err != nil {
			log.Println("failed to refresh library_metadata fts index", err)
		}
		return nil
	})

	if t.AutoDiscovery {
		go func() {
			if err := AutoDiscovery(dctx, db, tclient, tstore); err != nil {
				asyncfailure(errorsx.Wrap(err, "autodiscovery from peers failed"))
				return
			}
		}()
	} else {
		log.Println("autodiscovery is disabled, to enable add --auto-discovery flag, this is an alpha feature.")
	}

	go func() {
		if err := DiscoverFromRSSFeeds(dctx, db, rootstore, tclient, tstore); err != nil {
			asyncfailure(errorsx.Wrap(err, "autodiscovery of RSS feeds failed"))
			return
		}
	}()

	go ResumeDownloads(dctx, db, rootstore, tclient, tstore)

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

	media.NewHTTPLibrary(db, mediastore).Bind(httpmux.PathPrefix("/m").Subrouter())
	media.NewHTTPDiscovered(db, tclient, tstore).Bind(httpmux.PathPrefix("/d").Subrouter())
	media.NewHTTPRSSFeed(db).Bind(httpmux.PathPrefix("/rss").Subrouter())

	tlspem := envx.String(userx.DefaultCacheDirectory(userx.DefaultRelRoot(), "tls.pem"), env.DaemonTLSPEM)
	if err = tlsx.SelfSignedLocalHostTLS(tlspem); err != nil {
		return err
	}

	_ = httpmux.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		if uri, err := route.URLPath(); err == nil {
			log.Println("Route", errorsx.Zero(route.GetPathTemplate()), errorsx.Zero(route.GetMethods()), uri.String())
		}

		return nil
	})

	if !t.DisableMDNS {
		if err := MulticastService(dctx, httpbind); err != nil {
			return errorsx.Wrap(err, "unable to setup multicast service")
		}
	}

	log.Println("https listening on:", httpbind.Addr().String(), tlspem)
	if cause := http.ServeTLS(httpbind, httpmux, tlspem, tlspem); cause != nil {
		return errorsx.Wrap(cause, "http server stopped")
	}

	// report any async failures.
	return dctx.Err()
}
