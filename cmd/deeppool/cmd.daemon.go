package main

import (
	"context"
	"database/sql"
	"embed"
	"io"
	"io/fs"
	"iter"
	"log"
	"net"
	"net/netip"
	"slices"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/anacrolix/stm/rate"
	"github.com/pressly/goose/v3"

	"github.com/james-lawrence/deeppool/cmd/cmdopts"
	"github.com/james-lawrence/deeppool/downloads"
	"github.com/james-lawrence/deeppool/internal/x/errorsx"
	"github.com/james-lawrence/deeppool/internal/x/goosex"
	"github.com/james-lawrence/deeppool/internal/x/langx"
	"github.com/james-lawrence/deeppool/internal/x/md5x"
	"github.com/james-lawrence/deeppool/internal/x/netipx"
	"github.com/james-lawrence/deeppool/internal/x/slicesx"
	"github.com/james-lawrence/deeppool/internal/x/sqlx"
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

func (t cmdDaemon) Run(ctx *cmdopts.Global, peerID *cmdopts.PeerID) (err error) {
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

	go func() {
		dht, ok := slicesx.First(tclient.DhtServers()...)
		if !ok {
			log.Println("No DHT servers")
			return
		}

		log.Println("autodiscovery of hashes initiated")
		defer log.Println("autodiscovery of hashes completed")
		for id, err := range Auto(ctx.Context, db, dht) {
			if err != nil {
				log.Println("autodiscovery failed", err)
				return
			}

			metadata, err := torrent.New(id)
			if err != nil {
				log.Println("unable create metadata from hash, skipping", id.HexString(), err)
				continue
			}

			func() {
				dctx, done := context.WithTimeout(ctx.Context, time.Minute)
				defer done()

				log.Println("retrieving info initiated", id.HexString())
				info, err := tclient.Info(dctx, metadata)
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

	go func() {
		dht, ok := slicesx.First(tclient.DhtServers()...)
		if !ok {
			log.Println("No DHT servers")
			return
		}
		log.Println("autodiscovery of samplable peers initiated")
		defer log.Println("autodiscovery of samplable peers completed")
		if err := FindBep51Peers(ctx.Context, db, dht); err != nil {
			log.Println("peer locating failed", err)
		}
	}()

	<-ctx.Context.Done()
	return nil
}

// request samples from the domain space.
func Auto(ctx context.Context, db sqlx.Queryer, s *dht.Server) iter.Seq2[metainfo.Hash, error] {
	return func(yield func(metainfo.Hash, error) bool) {
		runsample := func(ctx context.Context, p tracking.Peer) (_ bool, err error) {
			var (
				resp bep0051.Response
			)

			defer func() {
				if err == nil {
					return
				}
				log.Println("marking peer as failed")
				if err := tracking.PeerMarkNextCheck(ctx, db, langx.Clone(p, tracking.PeerOptionBEP51(p.Bep51Available, p.Bep51TTL))).Scan(&p); err != nil {
					log.Println(errorsx.Wrapf(err, "unable update peer record: %s", p.IP))
				}
			}()

			req, b, err := bep0051.NewRequestBinary(s.ID(), krpc.ID(p.Peer))
			if err != nil {
				return true, errorsx.Wrapf(err, "unable to prepare sample request: %s", p.IP)
			}
			dst := dht.NewAddr(net.UDPAddrFromAddrPort(netip.AddrPortFrom(netipx.AddrFromSlice(net.ParseIP(p.IP)), p.Port)))

			log.Println("attempting to sample", p.IP, dst.String())
			encoded, _, err := s.QueryContext(ctx, dst, req.Q, req.T, b)
			if err != nil {
				return true, errorsx.Wrapf(err, "query failed: %s", dst.String())
			}

			if err := bencode.Unmarshal(encoded, &resp); err != nil {
				return true, errorsx.Wrapf(err, "unable to deserialized sample response: %s", p.IP)
			}

			for id := range slices.Chunk(resp.R.Sample, 20) {
				var (
					unknown tracking.UnknownHash
				)
				if !yield(metainfo.Hash(id), nil) {
					return false, errorsx.New("yield failed")
				}

				if err = tracking.UnknownHashInsertWithDefaults(ctx, db, tracking.NewUnknownHash(metainfo.Hash(id))).Scan(&unknown); err != nil {
					return true, errorsx.Wrapf(err, "unable to track hash: %s", md5x.Digest(id))
				}
			}

			if err := tracking.PeerMarkNextCheck(ctx, db, langx.Clone(p, tracking.PeerOptionBEP51(uint64(resp.R.Available), uint16(resp.R.Interval)))).Scan(&p); err != nil {
				return true, errorsx.Wrapf(err, "unable update peer record: %s", p.IP)
			}

			return true, nil
		}

		for {
			q := tracking.PeerSearchBuilder().Where(
				squirrel.And{
					tracking.PeerQueryHasInfoHashes(),
					tracking.PeerQueryNeedsCheck(),
				},
			)

			scanner := tracking.PeerSearch(ctx, db, q)
			defer scanner.Close()

			for scanner.Next() {
				var (
					p tracking.Peer
				)

				if err := scanner.Scan(&p); err != nil {
					yield(metainfo.Hash{}, err)
					return
				}

				if ok, err := runsample(ctx, p); err != nil && ok {
					log.Println(err)
					continue
				} else if err != nil {
					log.Println(err)
					break
				}
			}

			if err := scanner.Err(); err != nil {
				yield(metainfo.Hash{}, err)
				return
			}
		}
	}
}

// randomly samples nodes from the dht.
func FindBep51Peers(ctx context.Context, db sqlx.Queryer, s *dht.Server) (err error) {
	l := rate.NewLimiter(rate.Every(100*time.Millisecond), 1)

	recordinterestingpeer := func(ctx context.Context, db sqlx.Queryer, s *dht.Server, n krpc.NodeInfo) error {
		var (
			resp bep0051.Response
			peer tracking.Peer
		)

		req, b, err := bep0051.NewRequestBinary(s.ID(), n.ID)
		if err != nil {
			return errorsx.Wrapf(err, "unable to generate sample request: %s", n.ID)
		}
		dst := dht.NewAddr(n.Addr.UDP())

		dctx, done := context.WithTimeout(ctx, 5*time.Second)
		defer done()

		encoded, _, err := s.QueryContext(dctx, dst, req.Q, req.T, b)
		if err != nil {
			return errorsx.Wrap(err, "query failed")
		}

		if err := bencode.Unmarshal(encoded, &resp); err != nil {
			return errorsx.Wrapf(err, "unable to generate sample response: %s", n.ID)
		}

		peer = tracking.NewPeer(n, tracking.PeerOptionBEP51(uint64(resp.R.Available), uint16(resp.R.Interval)))

		// if they have no hashes they are not interesting.
		if resp.R.Available == 0 {
			return nil
		}

		// track peers with large libraries.
		if err := tracking.PeerInsertWithDefaults(ctx, db, peer).Scan(&peer); err != nil {
			return errorsx.Wrapf(err, "unable to record interesting peer %s", n.ID)
		} else {
			log.Println("interesting peer", peer.ID, resp.Y, peer.Bep51, peer.Bep51TTL, peer.Bep51Available)
		}

		return nil
	}

	for err = l.Wait(ctx); err == nil; err = l.Wait(ctx) {
		if s.NumNodes() < 64 {
			log.Println("minimum nodes not available, waiting", s.NumNodes())
			time.Sleep(time.Second)
			continue
		}

		log.Println("locating samplable peers", s.NumNodes(), "available")

		for _, n := range s.Nodes() {
			if err := recordinterestingpeer(ctx, db, s, n); err != nil {
				log.Println(err)
			}
		}
	}

	return err
}
