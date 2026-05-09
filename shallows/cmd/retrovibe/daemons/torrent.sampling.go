package daemons

import (
	"context"
	"errors"
	"log"
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
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/asynccompute"
	"github.com/retrovibed/retrovibed/shallows/internal/backoffx"
	"github.com/retrovibed/retrovibed/shallows/internal/contextx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/internal/torrentx"
	"github.com/retrovibed/retrovibed/shallows/internal/userx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"golang.org/x/time/rate"
)

// discover peers in the dht who support bep51.
func DiscoverDHTBEP51Peers(ctx context.Context, q sqlx.Queryer, s *dht.Server) (err error) {
	l := rate.NewLimiter(rate.Every(10*time.Second), 1)

	recordinterestingpeer := func(ctx context.Context, db sqlx.Queryer, s *dht.Server, n krpc.NodeInfo) (err error) {
		var (
			peer        tracking.Peer
			bep51       tracking.PeerOption = tracking.PeerOptionNoop
			ddiscpeer   tracking.PeerOption = tracking.PeerOptionNoop
			interesting bool
		)

		{
			failed := func() error {
				var (
					sampled *bep0051.Sample
				)

				dctx, done := context.WithTimeout(ctx, 30*time.Second)
				defer done()

				if sampled, err = bep0051.LatestSampleForNodeInfo(dctx, s, n); err != nil {
					return err
				}

				bep51 = tracking.PeerOptionBEP51(uint64(sampled.Available), uint16(sampled.Interval))
				// if they have hashes they are not interesting.
				interesting = interesting && sampled.Available > 0
				return nil
			}()
			errorsx.Log(errorsx.Wrap(failed, "failure checking if peer supports bep51"))
		}

		// nothing about this peer is interesting.
		if !interesting {
			return nil
		}

		peer = tracking.NewPeerFromInfo(
			n,
			bep51,
			ddiscpeer,
			tracking.PeerOptionTombstone(time.Now().Add(30*24*time.Hour)),
		)

		// track peers with large libraries.
		if err := tracking.PeerInsertWithDefaults(ctx, db, peer).Scan(&peer); err != nil {
			return errorsx.Wrapf(err, "unable to record interesting peer %s", n.ID)
		} else if peer.CreatedAt.Before(peer.UpdatedAt) {
			log.Println("interesting peer", peer.ID, peer.Bep51, peer.Bep51TTL, peer.Bep51Available, peer.CreatedAt, peer.CreatedAt.Equal(peer.UpdatedAt))
		}

		return nil
	}

	for err = l.Wait(ctx); err == nil; err = l.Wait(ctx) {
		log.Println("locating samplable peers", s.NumNodes(), "available")

		if err := sqlx.Discard(sqlx.Scan(tracking.PeerClearTombstoned(ctx, q, time.Now()))); err != nil {
			log.Println("unable to clear tombstoned peers", err)
		}

		target := int160.Random()
		peers := s.ClosestGoodNodeInfos(s.DynamicAddrPort(), 256, target)

		for _, dis := range peers {
			ret := s.GetPeers(ctx, dht.NewAddr(dis.Addr.AddrPort), target, false)
			if ret.Err != nil {
				errorsx.Log(errorsx.Ignore(errorsx.Wrap(ret.Err, "failed to discover peers from dht"), context.Canceled))
				continue
			}

			for _, n := range torrentx.NodesFromReply(ret) {
				if err := recordinterestingpeer(ctx, q, s, n); err != nil {
					log.Println(err)
					continue
				}
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

		dst := netip.AddrPortFrom(p.IP, p.Port)
		dstaddr := dht.NewAddr(dst)

		log.Println("infohash sample initiated", p.IP, dstaddr.String())
		defer log.Println("infohash sample completed", p.IP, dstaddr.String())

		defer func() {
			if err == nil {
				return
			}

			log.Println("marking peer as failed", err)
			if err := tracking.PeerMarkNextCheck(ctx, db, langx.Clone(p, tracking.PeerOptionBEP51(p.Bep51Available, p.Bep51TTL))).Scan(&p); err != nil {
				log.Println(errorsx.Wrapf(err, "unable update peer record: %s", p.IP))
			}
		}()

		qi, err := bep0051.NewRequest(s.ID(dst), krpc.ID(p.Peer))
		if err != nil {
			return errorsx.Wrapf(err, "unable to prepare sample request: %s", p.IP)
		}

		ret := s.Query(ctx, dstaddr, qi)
		if ret.Err != nil {
			return errorsx.Wrapf(err, "query failed: %s", dstaddr.String())
		}

		if err := bencode.Unmarshal(ret.Raw, &resp); err != nil {
			return errorsx.Wrapf(err, "unable to deserialized sample response: %s", p.IP)
		}

		for id := range slices.Chunk(resp.R.Sample, 20) {
			var (
				known      tracking.Metadata
				discovered ddisc.Discovered
				unknown    tracking.UnknownHash
			)

			id := int160.FromBytes(id)

			if err := ddisc.DiscoveredFindByID(ctx, db, torrentx.HashUID(langx.Autoptr(id))).Scan(&discovered); err == nil {
				continue
			} else if sqlx.IgnoreNoRows(err) != nil {
				return errorsx.Wrap(err, "unable to determine if infohash is in discovered")
			}

			if err := tracking.MetadataFindByID(ctx, db, torrentx.HashUID(langx.Autoptr(id))).Scan(&known); err == nil {
				continue
			} else if sqlx.IgnoreNoRows(err) != nil {
				return errorsx.Wrap(err, "unable to determine if infohash is known")
			}

			unknown = tracking.NewUnknownHash(
				id,
				tracking.OptionUnknownHashPeer(int160.FromBytes(p.Peer), dst),
			)

			if err = tracking.UnknownHashInsertWithDefaults(ctx, db, unknown).Scan(&unknown); err != nil {
				return errorsx.Wrapf(err, "unable to track hash: %s", torrentx.HashUID(langx.Autoptr(id)))
			}
		}

		p = langx.Clone(
			p,
			tracking.PeerOptionBEP51(uint64(resp.R.Available), uint16(resp.R.Interval)),
			tracking.PeerOptionTombstone(time.Now().Add(30*24*time.Hour)),
		)

		if err := tracking.PeerMarkNextCheck(ctx, db, p).Scan(&p); err != nil {
			return errorsx.Wrapf(err, "unable update peer record: %s", p.IP)
		}

		return nil
	}

	const workloads = uint16(12)

	querypeers := func(pool *asynccompute.Pool[tracking.Peer]) error {
		q := tracking.PeerSearchBuilder().Where(
			squirrel.And{
				tracking.PeerQueryHasInfoHashes(),
				tracking.PeerQueryNeedsCheck(),
			},
		).Limit(uint64(workloads) * 2)

		s := sqlx.Scan(tracking.PeerSearch(ctx, db, q))

		for p := range s.Iter() {
			errorsx.Log(pool.Run(ctx, p))
		}

		return s.Err()
	}

	l := rate.NewLimiter(rate.Every(10*time.Second), 1)
	getpending := func() int {
		return errorsx.Zero(sqlx.Count(ctx, db, "SELECT COUNT (*) FROM torrents_unknown_infohashes WHERE next_check < NOW()"))
	}

	samplers := asynccompute.New(func(ctx context.Context, unk tracking.Peer) error {
		return errorsx.Wrap(contextx.IgnoreDeadlineExceeded(runsample(ctx, unk)), "failed to retrieve metadata")
	}, asynccompute.Backlog[tracking.Peer](workloads))

	for err, pending := l.Wait(ctx), getpending(); err == nil; err, pending = l.Wait(ctx), getpending() {
		if pending < 100 {
			log.Println("querying peers for info hashes", pending, "< 100")
		} else {
			continue
		}

		if err := querypeers(samplers); err != nil {
			log.Println("failed to query peers", err)
		}
	}

	return errorsx.Compact(ctx.Err(), asynccompute.Shutdown(context.Background(), samplers))
}

func WaitForMinimumNodes(ctx context.Context, min int, dhts *dht.Server, do func()) {
	b := backoffx.New(backoffx.Exponential(time.Second), backoffx.Maximum(time.Minute))
	for attempts := 0; ; attempts++ {
		if dhts.NumNodes() > 32 {
			break
		}

		log.Printf("minimum nodes not available, waiting %p %d\n", dhts, dhts.NumNodes())
		time.Sleep(b.Backoff(attempts))
	}

	do()
}

// request samples from the domain space.
func DiscoverDHTMetadata(ctx context.Context, workloads uint64, db sqlx.Queryer, tclient *torrent.Client, blocked ddisc.Filter) error {
	l := rate.NewLimiter(rate.Every(10*time.Second), 1)
	timeouts := uint64(0)
	ttstore := storage.NewFile(userx.DefaultCacheDirectory("torrentddisc"))
	runsample := func(ctx context.Context, timeout time.Duration, unk tracking.UnknownHash) (err error) {
		var (
			unknown tracking.UnknownHash
			disc    ddisc.Discovered
		)

		timeout = timeout + backoffx.DynamicHashDuration(timeout, unk.ID)

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

		metadata, err := torrent.New(metainfo.Hash(unk.Infohash), torrent.OptionStorage(ttstore))
		if err != nil {
			return errorsx.Wrapf(err, "unable to create metadata from infohash %s", unk.ID)
		}
		defer tclient.Stop(metadata)

		dctx, done := context.WithTimeout(ctx, timeout)
		defer done()

		peeropt := torrent.TuneNoop
		if len(unk.Peer) > 0 && unk.Port > 0 && unk.IP.IsValid() {
			peeropt = torrent.TunePeers(torrent.NewPeer(int160.FromBytes(unk.Peer), netip.AddrPortFrom(netip.AddrFrom16(unk.IP.As16()), unk.Port)))
		}
		log.Println("metadata lookup initiated", metadata.ID, unk.IP, unk.Port)
		info, err := tclient.Info(
			dctx,
			metadata,
			peeropt,
			torrent.TuneAnnounceUntilComplete,
		)
		if contextx.IsDeadlineExceeded(err) {
			return langx.FirstNonZero(errorsx.Wrap(tracking.UnknownHashCooldown(ctx, db, unk).Scan(&unk), "unable to cooldown unknown hash"), err)
		}

		if err != nil {
			return errorsx.Wrapf(err, "unable to download metadata for infohash %s", unk.ID)
		}

		log.Println("metadata lookup completed", metadata.ID, unk.IP, unk.Port)

		id := metadata.ID

		if err = ddisc.DiscoveredInsertWithDefaults(
			ctx,
			db,
			ddisc.NewDiscovered(
				&id,
				ddisc.DiscoveredOptionIndex(!blocked(id.Bytes())),
				ddisc.DiscoveredOptionMimetype(mimex.Binary),
				ddisc.DiscoveredOptionFromTorrentInfo(info),
			),
		).Scan(&disc); err != nil {
			return errorsx.Wrap(err, "unable to insert discovered record")
		}

		if err := tracking.UnknownHashDeleteByID(ctx, db, torrentx.HashUID(&id)).Scan(&unknown); err != nil {
			return errorsx.Wrapf(err, "unable to delete unknown infohash: %s", unk.ID)
		}

		return nil
	}

	locatehashed := func(ctx context.Context) sqlx.Iter[tracking.UnknownHash] {
		// consider newest unknown hashes first.
		q := tracking.UnknownSearchBuilder().Where(
			squirrel.And{
				tracking.UnknownHashQueryNeedsCheck(),
			},
		).OrderBy("attempts ASC, created_at DESC").Limit(workloads * 10)
		return sqlx.Scan(tracking.UnknownSearch(ctx, db, q))
	}

	buff := make(chan tracking.UnknownHash, workloads)
	for i := range workloads {
		go func(i uint64) {
			time.Sleep(backoffx.DynamicHashDuration(10*time.Second, strconv.FormatInt(int64(i), 36)))
			bs := backoffx.New(backoffx.Exponential(1*time.Second), backoffx.Minimum(5*time.Second), backoffx.Maximum(45*time.Second))
			for unk := range buff {
				if err := runsample(ctx, bs.Backoff(int(unk.Attempts)), unk); errorsx.Ignore(err, context.DeadlineExceeded) != nil {
					log.Println("failed to retrieve metadata", unk.ID, err)
					continue
				}
			}
		}(i)
	}

	bs := backoffx.New(backoffx.Exponential(time.Second), backoffx.Maximum(1*time.Minute))
	for attempts := 0; ; attempts += 1 {
		s := locatehashed(ctx)
		log.Println("locate hashed initiated")
		for unk := range s.Iter() {
			select {
			case buff <- unk:
				attempts = -1
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		if err := s.Err(); errors.Is(err, context.Canceled) {
			return err
		} else if err != nil {
			log.Println("locate hashed failed", err)
			continue
		}
		log.Println("locate hashed completed")

		log.Println("sleeping for", bs.Backoff(attempts))
		select {
		case <-time.After(bs.Backoff(attempts)):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func PrintStatistics(ctx context.Context, q sqlx.Queryer) {
	timex.NowAndEveryVoid(ctx, 5*time.Minute, func(ctx context.Context) {
		type stats struct {
			Pending   int
			Available int
			Offload   int
			Indexing  int
			Indexed   int
			Peers     int
			RSS       int
		}

		m := stats{
			Pending:   errorsx.Zero(sqlx.Count(ctx, q, "SELECT COUNT (*) FROM torrents_unknown_infohashes WHERE next_check < NOW()")),
			Indexing:  errorsx.Zero(sqlx.Count(ctx, q, "SELECT COUNT (*) FROM ddisc_media WHERE known_media_id = 'ffffffff-ffff-ffff-ffff-ffffffffffff'")),
			Offload:   errorsx.Zero(sqlx.Count(ctx, q, "SELECT COUNT (*) FROM ddisc_media WHERE known_media_id = '00000000-0000-0000-0000-000000000000'")),
			Indexed:   errorsx.Zero(sqlx.Count(ctx, q, "SELECT COUNT (*) FROM ddisc_media WHERE known_media_id NOT IN ('00000000-0000-0000-0000-000000000000', 'ffffffff-ffff-ffff-ffff-ffffffffffff')")),
			Available: errorsx.Zero(sqlx.Count(ctx, q, "SELECT COUNT (*) FROM torrents_metadata")),
			Peers:     errorsx.Zero(sqlx.Count(ctx, q, "SELECT COUNT (*) FROM torrents_peers WHERE next_check < NOW()")),
			RSS:       errorsx.Zero(sqlx.Count(ctx, q, "SELECT COUNT (*) FROM torrents_feed_rss WHERE next_check < NOW()")),
		}

		log.Println("status", spew.Sdump(m))
	})
}

func AutoDiscovery(ctx context.Context, q sqlx.Queryer, dhts *dht.Server, tstore storage.ClientImpl) error {
	go func() {
		log.Println("autodiscovery of hashes initiated")
		defer log.Println("autodiscovery of hashes completed")
		if err := DiscoverDHTInfoHashes(ctx, q, dhts); err != nil {
			log.Println("autodiscovery of hashes failed", err)
			return
		}
	}()

	go func() {
		log.Println("autodiscovery of samplable peers initiated")
		defer log.Println("autodiscovery of samplable peers completed")
		if err := DiscoverDHTBEP51Peers(ctx, q, dhts); err != nil {
			log.Println("peer locating failed", err)
		}
	}()

	return nil
}
