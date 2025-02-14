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
	"runtime"
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
	"github.com/james-lawrence/deeppool/internal/x/timex"
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

	tm := dht.DefaultMuxer().
		Method(bep0051.Query, bep0051.NewEndpoint(bep0051.EmptySampler{}))
	tclient, err := tnetwork.Bind(
		torrent.NewClient(
			torrent.NewDefaultClientConfig(
				torrent.ClientConfigPeerID(string(peerID[:])),
				torrent.ClientConfigSeed(true),
				torrent.ClientConfigInfoLogger(log.New(io.Discard, "", log.Flags())),
				torrent.ClientConfigMuxer(tm),
				torrent.ClientConfigStorageDir(userx.DefaultDataDirectory(userx.DefaultRelRoot(), "torrents")),
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

	go func() {
		dht, ok := slicesx.First(tclient.DhtServers()...)
		if !ok {
			log.Println("No DHT servers")
			return
		}

		log.Println("auto retrieval of torrent info initiated")
		defer log.Println("auto retrieval of torrent info completed")

		if err := ResolveUnknownInfoHashes(ctx.Context, db, dht, tclient); err != nil {
			log.Println("resolving info hashes has failed", err)
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
		if err := SampleHashes(ctx.Context, db, dht); err != nil {
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
		if err := FindBep51Peers(ctx.Context, db, dht); err != nil {
			log.Println("peer locating failed", err)
		}
	}()

	<-ctx.Context.Done()
	return nil
}

// request samples from the domain space.
func ResolveUnknownInfoHashes(ctx context.Context, db sqlx.Queryer, s *dht.Server, tclient *torrent.Client) error {
	runsample := func(ctx context.Context, unk tracking.UnknownHash) (err error) {
		var (
			unknown tracking.UnknownHash
			md      tracking.Metadata
		)
		dctx, done := context.WithTimeout(ctx, time.Minute)
		defer done()

		log.Println("locate infohash initiated", unk.ID)
		defer log.Println("locate infohash completed", unk.ID)

		metadata, err := torrent.New(metainfo.Hash(unk.Infohash))
		if err != nil {
			return errorsx.Wrapf(err, "unable to create metadata from infohash %s", unk.ID)
		}

		info, err := tclient.Info(dctx, metadata)
		if err != nil {
			return errorsx.Wrapf(err, "unable to download metadata for infohash %s", unk.ID)
		}
		defer tclient.Stop(metadata)

		log.Println("retrieving info completed", unk.ID, info.Name, info.Length, info.TotalLength())

		if err = tracking.MetadataInsertWithDefaults(ctx, db, tracking.NewMetadata(&metadata.InfoHash, tracking.MetadataOptionFromInfo(info))).Scan(&md); err != nil {
			return errorsx.Wrap(err, "unable to insert metadata")
		}

		if err := tracking.UnknownHashDeleteByID(ctx, db, tracking.HashUID(&metadata.InfoHash)).Scan(&unknown); err != nil {
			return errorsx.Wrapf(err, "unable to delete unknown infohash: %s", unk.ID)
		}

		return nil
	}

	locatehashed := func(ctx context.Context) iter.Seq2[tracking.UnknownHash, error] {
		return func(yield func(tracking.UnknownHash, error) bool) {
			// consider newest unknown hashes first.
			q := tracking.UnknownSearchBuilder().OrderBy("created_at DESC").Limit(128)
			scanner := tracking.UnknownSearch(ctx, db, q)
			defer scanner.Close()

			for scanner.Next() {
				var (
					p tracking.UnknownHash
				)

				if err := scanner.Scan(&p); err != nil {
					yield(tracking.UnknownHash{}, err)
					return
				}

				if !yield(p, nil) {
					return
				}
			}

			if err := scanner.Err(); err != nil {
				yield(tracking.UnknownHash{}, err)
				return
			}
		}
	}

	// pool := pond.NewPool(1)
	// defer pool.StopAndWait()
	buff := make(chan tracking.UnknownHash, runtime.NumCPU()*4)
	for i := 0; i < runtime.NumCPU()*4; i++ {
		go func() {
			for unk := range buff {
				if err := runsample(ctx, unk); err != nil {
					log.Println("failed to retrieve metadata", unk.ID, err)
				}
			}
		}()
	}

	for {
		for unk, err := range locatehashed(ctx) {
			if err != nil {
				log.Println("locating pending info hashes failed", err)
				break
			}

			select {
			case buff <- unk:
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}
}

// request samples from the domain space.
func SampleHashes(ctx context.Context, db sqlx.Queryer, s *dht.Server) error {
	runsample := func(ctx context.Context, p tracking.Peer) (err error) {
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
			return errorsx.Wrapf(err, "unable to prepare sample request: %s", p.IP)
		}
		dst := dht.NewAddr(net.UDPAddrFromAddrPort(netip.AddrPortFrom(netipx.AddrFromSlice(net.ParseIP(p.IP)), p.Port)))

		log.Println("infohash sample initiated", p.IP, dst.String())
		defer log.Println("infohash sample completed", p.IP, dst.String())

		encoded, _, err := s.QueryContext(ctx, dst, req.Q, req.T, b)
		if err != nil {
			return errorsx.Wrapf(err, "query failed: %s", dst.String())
		}

		if err := bencode.Unmarshal(encoded, &resp); err != nil {
			return errorsx.Wrapf(err, "unable to deserialized sample response: %s", p.IP)
		}

		for id := range slices.Chunk(resp.R.Sample, 20) {
			var (
				known   tracking.Metadata
				unknown tracking.UnknownHash
			)

			if err := tracking.MetadataFindByID(ctx, db, tracking.HashUID(langx.Autoptr(metainfo.Hash(id)))).Scan(&known); err == nil {
				continue
			} else if sqlx.IgnoreNoRows(err) != nil {
				return errorsx.Wrap(err, "unable to determine if infohash is known")
			}

			if err = tracking.UnknownHashInsertWithDefaults(ctx, db, tracking.NewUnknownHash(metainfo.Hash(id))).Scan(&unknown); err != nil {
				return errorsx.Wrapf(err, "unable to track hash: %s", md5x.Digest(id))
			}
		}

		if err := tracking.PeerMarkNextCheck(ctx, db, langx.Clone(p, tracking.PeerOptionBEP51(uint64(resp.R.Available), uint16(resp.R.Interval)))).Scan(&p); err != nil {
			return errorsx.Wrapf(err, "unable update peer record: %s", p.IP)
		}

		return nil
	}

	querypeers := func() error {
		q := tracking.PeerSearchBuilder().Where(
			squirrel.And{
				tracking.PeerQueryHasInfoHashes(),
				tracking.PeerQueryNeedsCheck(),
			},
		).Limit(128)

		scanner := tracking.PeerSearch(ctx, db, q)
		defer scanner.Close()

		for scanner.Next() {
			var (
				p tracking.Peer
			)

			if err := scanner.Scan(&p); err != nil {
				return err
			}

			if err := runsample(ctx, p); err != nil {
				log.Println(err)
				continue
			}
		}

		if err := scanner.Err(); err != nil {
			return err
		}

		return nil
	}

	l := rate.NewLimiter(rate.Every(10*time.Second), 1)
	for err := l.Wait(ctx); err == nil; err = l.Wait(ctx) {
		unknownc := errorsx.Zero(sqlx.Count(ctx, db, "SELECT COUNT (*) FROM torrents_unknown_infohashes"))
		metadatac := errorsx.Zero(sqlx.Count(ctx, db, "SELECT COUNT (*) FROM torrents_metadata"))
		log.Println("sampling cycle", unknownc, metadatac)

		if err := querypeers(); err != nil {
			log.Println("failed to query peers", err)
		}
	}

	return nil
}

// randomly samples nodes from the dht.
func FindBep51Peers(ctx context.Context, db sqlx.Queryer, s *dht.Server) (err error) {
	for {
		if s.NumNodes() > 64 {
			break
		}
		log.Println("minimum nodes not available, waiting", s.NumNodes())
		time.Sleep(time.Second)
	}

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
			log.Println("interesting peer", peer.ID, resp.Y, peer.Bep51, peer.Bep51TTL, peer.Bep51Available, peer.CreatedAt, peer.UpdatedAt)
		}

		return nil
	}

	for err = l.Wait(ctx); err == nil; err = l.Wait(ctx) {
		log.Println("locating samplable peers", s.NumNodes(), "available")

		for _, n := range s.Nodes() {
			if err := recordinterestingpeer(ctx, db, s, n); err != nil {
				log.Println(err)
			}
		}
	}

	return err
}
