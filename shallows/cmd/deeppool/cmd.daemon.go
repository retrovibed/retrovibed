package main

import (
	"database/sql"
	"embed"
	"io"
	"io/fs"
	"log"
	"time"

	"github.com/pressly/goose/v3"

	"github.com/james-lawrence/deeppool/cmd/cmdopts"
	"github.com/james-lawrence/deeppool/cmd/deeppool/daemons"
	"github.com/james-lawrence/deeppool/downloads"
	"github.com/james-lawrence/deeppool/internal/x/errorsx"
	"github.com/james-lawrence/deeppool/internal/x/goosex"
	"github.com/james-lawrence/deeppool/internal/x/slicesx"
	"github.com/james-lawrence/deeppool/internal/x/timex"
	"github.com/james-lawrence/deeppool/internal/x/torrentx"
	"github.com/james-lawrence/deeppool/internal/x/userx"
	"github.com/james-lawrence/torrent/dht/v2"
	"github.com/james-lawrence/torrent/storage"

	_ "github.com/marcboeker/go-duckdb"

	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/bep0051"
)

//go:embed .migrations/*.sql
var embedsqlite embed.FS

type cmdDaemon struct{}

func (t cmdDaemon) Run(ctx *cmdopts.Global, peerID *cmdopts.PeerID) (err error) {
	var (
		db           *sql.DB
		torrentpeers = userx.DefaultCacheDirectory(userx.DefaultRelRoot(), "torrent.peers")
	)

	if db, err = sql.Open("duckdb", "dpool.db"); err != nil {
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
				torrent.ClientConfigPeerID(string(peerID[:])),
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

	<-ctx.Context.Done()
	return nil
}
