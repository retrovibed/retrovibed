package main

import (
	"context"
	"database/sql"
	"embed"
	"io"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/justinas/alice"
	"github.com/pressly/goose/v3"
	"golang.org/x/crypto/ssh"

	"github.com/james-lawrence/deeppool/cmd/cmdopts"
	"github.com/james-lawrence/deeppool/cmd/shallows/daemons"
	"github.com/james-lawrence/deeppool/downloads"
	"github.com/james-lawrence/deeppool/internal/env"
	"github.com/james-lawrence/deeppool/internal/x/contextx"
	"github.com/james-lawrence/deeppool/internal/x/dhtx"
	"github.com/james-lawrence/deeppool/internal/x/envx"
	"github.com/james-lawrence/deeppool/internal/x/errorsx"
	"github.com/james-lawrence/deeppool/internal/x/fsx"
	"github.com/james-lawrence/deeppool/internal/x/goosex"
	"github.com/james-lawrence/deeppool/internal/x/httpx"
	"github.com/james-lawrence/deeppool/internal/x/slicesx"
	"github.com/james-lawrence/deeppool/internal/x/tlsx"
	"github.com/james-lawrence/deeppool/internal/x/torrentx"
	"github.com/james-lawrence/deeppool/internal/x/userx"
	"github.com/james-lawrence/deeppool/media"
	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/dht/krpc"
	"github.com/james-lawrence/torrent/storage"

	_ "github.com/marcboeker/go-duckdb"

	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/bep0051"

	"github.com/gorilla/mux"
)

//go:embed .migrations/*.sql
var embedsqlite embed.FS

type cmdDaemon struct {
	AutoBootstrap bool             `flag:"" name:"auto-bootstrap" help:"bootstrap from a predefined set of peers" default:"false"`
	AutoDiscovery bool             `flag:"" name:"auto-discovery" help:"enable autodiscovery of content from peers" default:"false"`
	HTTP          cmdopts.Listener `flag:"" name:"http-address" help:"address to serve daemon api from" default:"tcp://:9998"`
}

func (t cmdDaemon) Run(gctx *cmdopts.Global, id *cmdopts.SSHID) (err error) {
	var (
		db           *sql.DB
		torrentpeers                            = userx.DefaultCacheDirectory(userx.DefaultRelRoot(), "torrent.peers")
		dbpath                                  = userx.DefaultConfigDir(userx.DefaultRelRoot(), "dpool.db")
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

	if db, err = sql.Open("duckdb", dbpath); err != nil {
		return errorsx.Wrap(err, "unable to open db")
	}
	defer db.Close()

	{
		mprov, err := goose.NewProvider("", db, errorsx.Must(fs.Sub(embedsqlite, ".migrations")), goose.WithStore(goosex.DuckdbStore{}))
		if err != nil {
			return errorsx.Wrap(err, "unable to build migration provider")
		}

		if _, err := mprov.Up(dctx); err != nil {
			return errorsx.Wrap(err, "unable to run migrations")
		}
	}

	go func() {
		errorsx.Log(errorsx.Wrap(daemons.PrepareDefaultFeeds(dctx, db), "unable to initialize default rss feeds"))
	}()

	tnetwork, err := torrentx.Autosocket(0)
	if err != nil {
		return errorsx.Wrap(err, "unable to setup torrent socket")
	}

	torrentdir := userx.DefaultDataDirectory(userx.DefaultRelRoot(), "torrents")

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

	tstore := storage.NewFileByInfoHash(torrentdir)

	dwatcher, err := downloads.NewDirectoryWatcher(dctx, db, tclient, tstore)
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

	go daemons.PrintStatistics(dctx, db)

	if t.AutoDiscovery {
		go func() {
			if err := daemons.AutoDiscovery(dctx, db, tclient, tstore); err != nil {
				asyncfailure(errorsx.Wrap(err, "autodiscovery from peers failed"))
				return
			}
		}()
	} else {
		log.Println("autodiscovery is disabled, to enable add --auto-discovery flag, this is an alpha feature.")
	}

	go func() {
		if err := daemons.DiscoverFromRSSFeeds(dctx, db, tclient, tstore); err != nil {
			asyncfailure(errorsx.Wrap(err, "autodiscovery of RSS feeds failed"))
			return
		}
	}()

	go daemons.ResumeDownloads(dctx, db, tclient, tstore)

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

	media.NewHTTPLibrary(db, tstore).Bind(httpmux.PathPrefix("/m").Subrouter())
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

	if err := daemons.MulticastService(dctx, httpbind); err != nil {
		return errorsx.Wrap(err, "unable to setup multicast service")
	}

	log.Println("https listening on:", httpbind.Addr().String(), tlspem)
	if cause := http.ServeTLS(httpbind, httpmux, tlspem, tlspem); cause != nil {
		return errorsx.Wrap(cause, "http server stopped")
	}

	// report any async failures.
	return dctx.Err()
}
