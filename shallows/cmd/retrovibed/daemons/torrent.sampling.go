package daemons

import (
	"context"
	"iter"
	"log"
	"net"
	"net/netip"
	"slices"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/davecgh/go-spew/spew"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/bencode"
	"github.com/james-lawrence/torrent/bep0051"
	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/dht/krpc"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/james-lawrence/torrent/storage"
	"github.com/retrovibed/retrovibed/internal/backoffx"
	"github.com/retrovibed/retrovibed/internal/contextx"
	"github.com/retrovibed/retrovibed/internal/debugx"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/langx"
	"github.com/retrovibed/retrovibed/internal/netipx"
	"github.com/retrovibed/retrovibed/internal/slicesx"
	"github.com/retrovibed/retrovibed/internal/sqlx"
	"github.com/retrovibed/retrovibed/internal/timex"
	"github.com/retrovibed/retrovibed/internal/torrentx"
	"github.com/retrovibed/retrovibed/tracking"
	"golang.org/x/time/rate"
)

// discover peers in the dht who support bep51.
func DiscoverDHTBEP51Peers(ctx context.Context, q sqlx.Queryer, s *dht.Server) (err error) {
	for {
		if s.NumNodes() > 32 {
			break
		}
		log.Println("minimum nodes not available, waiting", s.NumNodes())
		time.Sleep(time.Second)
	}

	l := rate.NewLimiter(rate.Every(10*time.Second), 1)

	recordinterestingpeer := func(ctx context.Context, db sqlx.Queryer, s *dht.Server, n krpc.NodeInfo) error {
		var (
			resp bep0051.Response
			peer tracking.Peer
		)

		qi, err := bep0051.NewRequest(s.ID(), n.ID)
		if err != nil {
			return errorsx.Wrapf(err, "unable to generate sample request: %s", n.ID)
		}
		dst := dht.NewAddr(n.Addr.UDP())

		dctx, done := context.WithTimeout(ctx, 30*time.Second)
		defer done()

		ret := s.Query(dctx, dst, qi)
		if ret.Err != nil {
			return errorsx.Wrap(err, "sample query failed")
		}

		if err := bencode.Unmarshal(ret.Raw, &resp); err != nil {
			if _, ok := err.(bencode.ErrUnusedTrailingBytes); !ok {
				return errorsx.Wrapf(err, "unable to deserialize sample response: %T %s", err, n.ID)
			}
		}

		peer = tracking.NewPeer(n, tracking.PeerOptionBEP51(uint64(resp.R.Available), uint16(resp.R.Interval)))

		// if they have no hashes they are not interesting.
		if resp.R.Available == 0 {
			return nil
		}

		// track peers with large libraries.
		if err := tracking.PeerInsertWithDefaults(ctx, db, peer).Scan(&peer); err != nil {
			return errorsx.Wrapf(err, "unable to record interesting peer %s", n.ID)
		} else if peer.CreatedAt.Before(peer.UpdatedAt) {
			log.Println("interesting peer", peer.ID, resp.Y, peer.Bep51, peer.Bep51TTL, peer.Bep51Available, peer.CreatedAt, peer.CreatedAt.Equal(peer.UpdatedAt))
		}

		return nil
	}

	for err = l.Wait(ctx); err == nil; err = l.Wait(ctx) {
		log.Println("locating samplable peers", s.NumNodes(), "available")

		target := int160.Random()
		dis, ok := slicesx.First(s.ClosestGoodNodeInfos(1, target)...)
		if !ok {
			continue
		}
		ret := s.GetPeers(ctx, dht.NewAddr(dis.Addr.UDP()), target, false, dht.QueryRateLimiting{})
		if ret.Err != nil {
			log.Println("failed to discover peers from dht", err)
			continue
		}

		for _, n := range torrentx.NodesFromReply(ret) {
			if err := recordinterestingpeer(ctx, q, s, n); err != nil {
				log.Println(err)
				continue
			}
		}
	}

	return err
}

// request samples from the domain space.
func DiscoverDHTInfoHashes(ctx context.Context, db sqlx.Queryer, s *dht.Server) error {
	runsample := func(ctx context.Context, p tracking.Peer) (err error) {
		var (
			resp bep0051.Response
		)

		defer func() {
			if err == nil {
				return
			}
			log.Println("marking peer as failed", err)
			if err := tracking.PeerMarkNextCheck(ctx, db, langx.Clone(p, tracking.PeerOptionBEP51(p.Bep51Available, p.Bep51TTL))).Scan(&p); err != nil {
				log.Println(errorsx.Wrapf(err, "unable update peer record: %s", p.IP))
			}
		}()

		qi, err := bep0051.NewRequest(s.ID(), krpc.ID(p.Peer))
		if err != nil {
			return errorsx.Wrapf(err, "unable to prepare sample request: %s", p.IP)
		}
		dst := dht.NewAddr(net.UDPAddrFromAddrPort(netip.AddrPortFrom(netipx.AddrFromSlice(net.ParseIP(p.IP)), p.Port)))

		log.Println("infohash sample initiated", p.IP, dst.String())
		defer log.Println("infohash sample completed", p.IP, dst.String())
		ret := s.Query(ctx, dst, qi)
		if ret.Err != nil {
			return errorsx.Wrapf(err, "query failed: %s", dst.String())
		}

		if err := bencode.Unmarshal(ret.Raw, &resp); err != nil {
			return errorsx.Wrapf(err, "unable to deserialized sample response: %s", p.IP)
		}

		log.Println("length of samples", len(resp.R.Sample)/20)
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
				return errorsx.Wrapf(err, "unable to track hash: %s", tracking.HashUID(langx.Autoptr(metainfo.Hash(id))))
			}
		}

		if err := tracking.PeerMarkNextCheck(ctx, db, langx.Clone(p, tracking.PeerOptionBEP51(uint64(resp.R.Available), uint16(resp.R.Interval)))).Scan(&p); err != nil {
			return errorsx.Wrapf(err, "unable update peer record: %s", p.IP)
		}

		return nil
	}

	const workloads = uint64(12)
	buff := make(chan tracking.Peer, workloads)
	for i := uint64(0); i < workloads; i++ {
		go func(i uint64) {
			time.Sleep(backoffx.DynamicHashDuration(time.Minute, strconv.FormatInt(int64(i), 36)))
			for unk := range buff {
				if err := runsample(ctx, unk); contextx.IgnoreDeadlineExceeded(err) != nil {
					log.Println("failed to retrieve metadata", unk.ID, err)
					continue
				}

				select {
				case <-ctx.Done():
					return
				default:
				}
			}
		}(i)
	}

	querypeers := func() error {
		q := tracking.PeerSearchBuilder().Where(
			squirrel.And{
				tracking.PeerQueryHasInfoHashes(),
				tracking.PeerQueryNeedsCheck(),
			},
		).Limit(workloads * 2)

		scanner := tracking.PeerSearch(ctx, db, q)
		defer scanner.Close()

		for scanner.Next() {
			var (
				p tracking.Peer
			)

			if err := scanner.Scan(&p); err != nil {
				return err
			}

			select {
			case buff <- p:
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		if err := scanner.Err(); err != nil {
			return err
		}

		return nil
	}

	l := rate.NewLimiter(rate.Every(10*time.Second), 1)
	getpending := func() int {
		return errorsx.Zero(sqlx.Count(ctx, db, "SELECT COUNT (*) FROM torrents_unknown_infohashes WHERE next_check < NOW()"))
	}

	for err, pending := l.Wait(ctx), getpending(); err == nil; err, pending = l.Wait(ctx), getpending() {
		if pending < 100 {
			log.Println("querying peers for info hashes", pending, "< 100")
		} else {
			continue
		}

		if err := querypeers(); err != nil {
			log.Println("failed to query peers", err)
		}
	}

	return ctx.Err()
}

// request samples from the domain space.
func DiscoverDHTMetadata(ctx context.Context, db sqlx.Queryer, s *dht.Server, tclient *torrent.Client, tstore storage.ClientImpl) error {
	for {
		if s.NumNodes() > 32 {
			break
		}
		log.Println("minimum nodes not available, waiting", s.NumNodes())
		time.Sleep(time.Second)
	}

	l := rate.NewLimiter(rate.Every(10*time.Second), 1)
	timeouts := uint64(0)
	workloads := uint64(384)

	runsample := func(ctx context.Context, timeout time.Duration, unk tracking.UnknownHash) (err error) {
		var (
			unknown tracking.UnknownHash
			md      tracking.Metadata
		)

		timeout = timeout + backoffx.DynamicHashDuration(timeout, unk.ID)
		dctx, done := context.WithTimeout(ctx, timeout)
		defer done()

		ts := time.Now()
		defer func() {
			st := time.Since(ts)

			if err == nil {
				log.Println("locate infohash completed", unk.ID, unk.Attempts, st, timeout)
				return
			}

			if l.Allow() && contextx.IsDeadlineExceeded(err) {
				tc := atomic.SwapUint64(&timeouts, 0) + 1
				log.Println("locate infohash timed out", tc, float32(tc)/float32(workloads), unk.ID, unk.Attempts, st, timeout)
				return
			} else if contextx.IsDeadlineExceeded(err) {
				atomic.AddUint64(&timeouts, 1)
				return
			}

			if err != nil {
				log.Println("locate infohash failed", unk.ID, unk.Attempts, st, timeout, err)
			}
		}()

		metadata, err := torrent.New(metainfo.Hash(unk.Infohash), torrent.OptionStorage(tstore))
		if err != nil {
			return errorsx.Wrapf(err, "unable to create metadata from infohash %s", unk.ID)
		}
		defer tclient.Stop(metadata)

		debugx.Println("initiating metadata lookup", metadata.ID)
		info, err := tclient.Info(dctx, metadata)
		if contextx.IsDeadlineExceeded(err) {
			return errorsx.Compact(tracking.UnknownHashCooldown(ctx, db, unk).Scan(&unk), err)
		}

		if err != nil {
			return errorsx.Wrapf(err, "unable to download metadata for infohash %s", unk.ID)
		}

		if err = tracking.MetadataInsertWithDefaults(ctx, db, tracking.NewMetadata(&metadata.ID, tracking.MetadataOptionFromInfo(info))).Scan(&md); err != nil {
			return errorsx.Wrap(err, "unable to insert metadata")
		}

		if err := tracking.UnknownHashDeleteByID(ctx, db, tracking.HashUID(&metadata.ID)).Scan(&unknown); err != nil {
			return errorsx.Wrapf(err, "unable to delete unknown infohash: %s", unk.ID)
		}

		return nil
	}

	locatehashed := func(ctx context.Context) iter.Seq2[tracking.UnknownHash, error] {
		return func(yield func(tracking.UnknownHash, error) bool) {
			log.Println("locate hashed initiated")
			defer log.Println("locate hashed completed")

			// consider newest unknown hashes first.
			q := tracking.UnknownSearchBuilder().Where(
				squirrel.And{
					tracking.UnknownHashQueryNeedsCheck(),
				},
			).OrderBy("attempts ASC, created_at DESC").Limit(workloads * 10)
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

	buff := make(chan tracking.UnknownHash, workloads)
	for i := uint64(0); i < workloads; i++ {
		go func(i uint64) {
			time.Sleep(backoffx.DynamicHashDuration(10*time.Second, strconv.FormatInt(int64(i), 36)))
			bs := backoffx.New(backoffx.Exponential(1*time.Second), backoffx.Minimum(5*time.Second), backoffx.Maximum(45*time.Second))
			for unk := range buff {
				if err := runsample(ctx, bs.Backoff(int(unk.Attempts)), unk); contextx.IgnoreDeadlineExceeded(err) != nil {
					log.Println("failed to retrieve metadata", unk.ID, err)
					continue
				}
			}
		}(i)
	}

	bs := backoffx.New(backoffx.Exponential(time.Second), backoffx.Maximum(1*time.Minute))
	for attempts := 0; ; attempts += 1 {
		for unk, err := range locatehashed(ctx) {
			if err != nil {
				log.Println("locating pending info hashes failed", err)
				continue
			}

			select {
			case buff <- unk:
				attempts = -1
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		_, err := db.ExecContext(ctx, "PRAGMA create_fts_index('torrents_metadata', 'id', 'description', overwrite = 1);")
		if err != nil {
			log.Println("failed to refresh fts index", err)
		}

		log.Println("sleeping for", bs.Backoff(attempts))
		select {
		case <-time.After(bs.Backoff(attempts)):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func PrintStatistics(ctx context.Context, q sqlx.Queryer) {
	timex.Every(30*time.Second, func() {
		type stats struct {
			Pending   int
			Available int
			Peers     int
			RSS       int
		}

		m := stats{
			Pending:   errorsx.Zero(sqlx.Count(ctx, q, "SELECT COUNT (*) FROM torrents_unknown_infohashes WHERE next_check < NOW()")),
			Available: errorsx.Zero(sqlx.Count(ctx, q, "SELECT COUNT (*) FROM torrents_metadata")),
			Peers:     errorsx.Zero(sqlx.Count(ctx, q, "SELECT COUNT (*) FROM torrents_peers WHERE next_check < NOW()")),
			RSS:       errorsx.Zero(sqlx.Count(ctx, q, "SELECT COUNT (*) FROM torrents_feed_rss WHERE next_check < NOW()")),
		}

		log.Println("status", spew.Sdump(m))
	})
}

func AutoDiscovery(ctx context.Context, q sqlx.Queryer, tclient *torrent.Client, tstore storage.ClientImpl) error {
	go func() {
		dht, ok := slicesx.First(tclient.DhtServers()...)
		if !ok {
			log.Println("No DHT servers")
			return
		}

		log.Println("auto retrieval of torrent info initiated")
		defer log.Println("auto retrieval of torrent info completed")

		if err := DiscoverDHTMetadata(ctx, q, dht, tclient, tstore); err != nil {
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
		if err := DiscoverDHTInfoHashes(ctx, q, dht); err != nil {
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
		if err := DiscoverDHTBEP51Peers(ctx, q, dht); err != nil {
			log.Println("peer locating failed", err)
		}
	}()
	return nil
}
