package main

import (
	"context"
	"database/sql"
	"embed"
	"io"
	"io/fs"
	"iter"
	"log"
	"slices"
	"time"

	"github.com/anacrolix/stm/rate"
	"github.com/pressly/goose/v3"

	"github.com/james-lawrence/deeppool/cmd/cmdopts"
	"github.com/james-lawrence/deeppool/downloads"
	"github.com/james-lawrence/deeppool/internal/x/errorsx"
	"github.com/james-lawrence/deeppool/internal/x/goosex"
	"github.com/james-lawrence/deeppool/internal/x/slicesx"
	"github.com/james-lawrence/deeppool/internal/x/torrentx"
	"github.com/james-lawrence/deeppool/internal/x/userx"
	"github.com/james-lawrence/deeppool/tracking"
	"github.com/james-lawrence/torrent/bencode"
	"github.com/james-lawrence/torrent/dht/v2"
	"github.com/james-lawrence/torrent/metainfo"

	_ "github.com/marcboeker/go-duckdb"

	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/bep0051"
	"github.com/james-lawrence/torrent/dht/v2/krpc"
)

//go:embed .migrations/*.sql
var embedsqlite embed.FS

type cmdDaemon struct{}

func (t cmdDaemon) Run(ctx *cmdopts.Global) (err error) {
	var (
		db *sql.DB
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

	tm := dht.DefaultMuxer().
		Method(bep0051.Query, bep0051.NewEndpoint(bep0051.EmptySampler{}))
	tclient, err := tnetwork.Bind(
		torrent.NewClient(
			torrent.NewDefaultClientConfig(
				torrent.ClientConfigPeerID(string([]byte{45, 71, 84, 48, 48, 48, 50, 45, 169, 218, 156, 162, 223, 141, 136, 209, 85, 207, 231, 113})),
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

	go func() {
		dht, ok := slicesx.First(tclient.DhtServers()...)
		if !ok {
			log.Println("No DHT servers")
			return
		}
		log.Println("autodiscovery of hashes initiated")
		defer log.Println("autodiscovery of hashes completed")
		for id, err := range Auto(ctx.Context, dht) {
			if err != nil {
				log.Println("autodiscovery failed", err)
				return
			}

			metadata, err := torrent.New(id)
			if err != nil {
				log.Println("unable create metadata from hash, skipping", id.HexString(), err)
				continue
			}

			go func() {
				log.Println("retrieving info initiated", id.HexString())
				info, err := tclient.Info(ctx.Context, metadata)
				if err != nil {
					log.Println("unable retrieve info from metadata, skipping", id.HexString(), err)
					return
				}
				log.Println("retrieving info completed", id.HexString(), info.Name, info.Length, info.TotalLength())

				var md tracking.Metadata
				if err = tracking.MetadataInsertWithDefaults(ctx.Context, db, tracking.NewMetadata(&metadata.InfoHash, tracking.MetadataOptionFromInfo(info))).Scan(&md); err != nil {
					log.Println("unable to record metadata", err)
					return
				}
			}()
		}
	}()

	<-ctx.Context.Done()
	return nil
}

// randomly samples nodes from the dht.
func Auto(ctx context.Context, db *sql.DB, s *dht.Server) iter.Seq2[metainfo.Hash, error] {
	return func(yield func(metainfo.Hash, error) bool) {
		var (
			err error
		)
		l := rate.NewLimiter(rate.Every(100*time.Millisecond), 1)

		for err = l.Wait(ctx); err == nil; err = l.Wait(ctx) {
			if s.NumNodes() < 64 {
				log.Println("minimum nodes not available, waiting", s.NumNodes())
				time.Sleep(time.Second)
				continue
			}

			nodes := s.MakeReturnNodes(dht.Int160FromByteArray(krpc.RandomID()), func(na krpc.NodeAddr) bool { return true })
			log.Println("sampling from", len(nodes), "nodes", s.NumNodes(), "available")

			for _, n := range nodes {
				var (
					resp bep0051.Response
				)
				req, b, err := bep0051.NewRequestBinary(s.ID(), n.ID)
				if err != nil {
					log.Println("unable to serialize sample request, skipping", n.ID, err)
					continue
				}
				dst := dht.NewAddr(n.Addr.UDP())

				log.Println("Requesting samples", dst.String())
				encoded, _, err := s.QueryContext(ctx, dst, req.Q, req.T, b)
				if err != nil {
					log.Println("query failed", err)
					continue
				}

				if err := bencode.Unmarshal(encoded, &resp); err != nil {
					log.Println("unable to deserialized sample response, skipping", n.ID, err)
					continue
				}

				tracking.PeerInsertWithDefaults(ctx, db)
				// track peers with large libraries.
				log.Println("interesting peer", resp.R.ID, resp.Y, resp.R.Interval, resp.R.Available)

				for id := range slices.Chunk(resp.R.Sample, 20) {
					if !yield(metainfo.Hash(id), nil) {
						return
					}
				}
			}
		}
		errorsx.MaybeLog(err)
	}
}
