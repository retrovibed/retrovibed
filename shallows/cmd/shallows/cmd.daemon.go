package main

import (
	"database/sql"
	"embed"
	"io"
	"io/fs"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/justinas/alice"
	"github.com/pressly/goose/v3"
	"golang.org/x/crypto/ssh"

	"github.com/james-lawrence/deeppool/cmd/cmdopts"
	"github.com/james-lawrence/deeppool/cmd/shallows/daemons"
	"github.com/james-lawrence/deeppool/downloads"
	"github.com/james-lawrence/deeppool/internal/env"
	"github.com/james-lawrence/deeppool/internal/x/envx"
	"github.com/james-lawrence/deeppool/internal/x/errorsx"
	"github.com/james-lawrence/deeppool/internal/x/goosex"
	"github.com/james-lawrence/deeppool/internal/x/httpx"
	"github.com/james-lawrence/deeppool/internal/x/slicesx"
	"github.com/james-lawrence/deeppool/internal/x/timex"
	"github.com/james-lawrence/deeppool/internal/x/tlsx"
	"github.com/james-lawrence/deeppool/internal/x/torrentx"
	"github.com/james-lawrence/deeppool/internal/x/userx"
	"github.com/james-lawrence/deeppool/media"
	"github.com/james-lawrence/torrent/dht/v2"
	"github.com/james-lawrence/torrent/dht/v2/krpc"
	"github.com/james-lawrence/torrent/storage"

	_ "github.com/marcboeker/go-duckdb"

	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/bep0051"

	"github.com/gorilla/mux"
)

//go:embed .migrations/*.sql
var embedsqlite embed.FS

type cmdDaemon struct{}

func (t cmdDaemon) Run(ctx *cmdopts.Global, id *cmdopts.SSHID) (err error) {
	var (
		db           *sql.DB
		torrentpeers = userx.DefaultCacheDirectory(userx.DefaultRelRoot(), "torrent.peers")
		dbpath       = userx.DefaultConfigDir(userx.DefaultRelRoot(), "dpool.db")
		peerid       = krpc.IdFromString(ssh.FingerprintSHA256(id.PublicKey()))
		httpbind     net.Listener
	)

	if db, err = sql.Open("duckdb", dbpath); err != nil {
		return errorsx.Wrap(err, "unable to open db")
	}
	defer db.Close()

	{
		mprov, err := goose.NewProvider("", db, errorsx.Must(fs.Sub(embedsqlite, ".migrations")), goose.WithStore(goosex.DuckdbStore{}))
		if err != nil {
			return errorsx.Wrap(err, "unable to build migration provider")
		}

		if _, err := mprov.Up(ctx.Context); err != nil {
			return errorsx.Wrap(err, "unable to run migrations")
		}
	}

	tnetwork, err := torrentx.Autosocket(9999)
	if err != nil {
		return errorsx.Wrap(err, "unable to setup torrent socket")
	}

	torrentdir := userx.DefaultDataDirectory(userx.DefaultRelRoot(), "torrents")
	tm := dht.DefaultMuxer().
		Method(bep0051.Query, bep0051.NewEndpoint(bep0051.EmptySampler{}))
	tclient, err := tnetwork.Bind(
		torrent.NewClient(
			torrent.NewDefaultClientConfig(
				torrent.ClientConfigPeerID(string(peerid[:])),
				torrent.ClientConfigSeed(true),
				torrent.ClientConfigInfoLogger(log.New(io.Discard, "", log.Flags())),
				torrent.ClientConfigMuxer(tm),
			),
		),
	)
	if err != nil {
		return errorsx.Wrap(err, "unable to setup torrent client")
	}

	dwatcher, err := downloads.NewDirectoryWatcher(ctx.Context, db, tclient)
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

	go timex.Every(time.Minute, func() {
		current := slicesx.FirstOrZero(tclient.DhtServers()...).Nodes()
		log.Println("saving torrent peers", len(current))
		errorsx.Log(
			errorsx.Wrap(
				dht.WriteNodesToFile(current, torrentpeers),
				"unable to persist peers",
			),
		)
	})

	go daemons.PrintStatistics(ctx.Context, db)

	go func() {
		dht, ok := slicesx.First(tclient.DhtServers()...)
		if !ok {
			log.Println("No DHT servers")
			return
		}

		log.Println("auto retrieval of torrent info initiated")
		defer log.Println("auto retrieval of torrent info completed")

		if err := daemons.DiscoverDHTMetadata(ctx.Context, db, dht, tclient, storage.NewFile(torrentdir)); err != nil {
			log.Println("resolving info hashes has failed", err)
			panic(err)
		}
	}()

	go func() {
		dht, ok := slicesx.First(tclient.DhtServers()...)
		if !ok {
			log.Println("No DHT servers")
			return
		}

		log.Println("autodiscovery of hashes initiated")
		defer log.Println("autodiscovery of hashes completed")
		if err := daemons.DiscoverDHTInfoHashes(ctx.Context, db, dht); err != nil {
			log.Println("autodiscovery of hashes failed", err)
			return
		}
	}()

	go func() {
		dht, ok := slicesx.First(tclient.DhtServers()...)
		if !ok {
			log.Println("No DHT servers")
			return
		}
		log.Println("autodiscovery of samplable peers initiated")
		defer log.Println("autodiscovery of samplable peers completed")
		if err := daemons.DiscoverDHTBEP51Peers(ctx.Context, db, dht); err != nil {
			log.Println("peer locating failed", err)
		}
	}()

	go func() {
		if err := daemons.DiscoverFromRSSFeeds(ctx.Context, db); err != nil {
			log.Println("autodiscovery of RSS feeds failed", err)
		}
	}()

	httpmux := mux.NewRouter()
	httpmux.NotFoundHandler = httpx.NotFound(alice.New())
	httpmux.Use(
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

	media.NewHTTPDiscovered(db).Bind(httpmux.PathPrefix("/d").Subrouter())

	if httpbind, err = net.Listen("tcp", ":9998"); err != nil {
		return err
	}

	tlspem := envx.String(userx.DefaultCacheDirectory(userx.DefaultRelRoot(), "tls.pem"), env.DaemonTLSPEM)
	if err = tlsx.SelfSignedLocalHostTLS(tlspem); err != nil {
		return err
	}

	go func() {
		<-ctx.Context.Done()
		httpbind.Close()
	}()

	_ = httpmux.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		if uri, err := route.URLPath(); err == nil {
			log.Println("Route", uri.String())
		}

		return nil
	})

	log.Println("https listening on:", httpbind.Addr().String(), tlspem)
	if cause := http.ServeTLS(httpbind, httpmux, tlspem, tlspem); cause != nil {
		log.Println("http server stopped", cause)
	}

	return nil
}
