package ddisctorrent

import (
	"context"
	"iter"
	"log"
	"net/netip"

	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/torrentx"
)

// closestPeersPerLookup bounds how many DHT routing-table nodes are asked
// per partition lookup, mirroring DiscoverDHTBEP51Peers's existing usage
// (shallows/cmd/retrovibed/daemons/torrent.sampling.go).
const closestPeersPerLookup = 256

// NewPartitionStrategy looks up peers on the DHT swarm for the partition a
// known-media-id maps to (without announcing this node onto it) and fires a
// fire-and-forget MethodSearch query at each. Never yields synchronously —
// any response lands in ddisc_media asynchronously via the already-registered
// MethodMedia responder, picked up by a later ddisc.LocalStrategy call.
func NewPartitionStrategy(dhts *dht.Server, partitions *ddisc.Partition) ddisc.DiscoverStrategy {
	return partitionStrategy{dhts: dhts, partitions: partitions}
}

type partitionStrategy struct {
	dhts       *dht.Server
	partitions *ddisc.Partition
}

func (t partitionStrategy) Discover(ctx context.Context, req ddisc.DiscoverRequest) iterx.Seq[ddisc.Discovered] {
	return &partitionSeq{cfg: t, req: req}
}

type partitionSeq struct {
	cfg partitionStrategy
	req ddisc.DiscoverRequest
	err error
}

func (t *partitionSeq) Each(ctx context.Context) iter.Seq[ddisc.Discovered] {
	return func(yield func(ddisc.Discovered) bool) {
		target := t.cfg.partitions.Max([]byte(t.req.KnownMediaID))
		id := int160.New(target.Bytes())

		searchreq, err := NewSearchRequest(t.cfg.dhts.ID(t.cfg.dhts.DynamicAddrPort()), t.req.KnownMediaID)
		if err != nil {
			t.err = err
			return
		}

		closest := t.cfg.dhts.ClosestGoodNodeInfos(t.cfg.dhts.DynamicAddrPort(), closestPeersPerLookup, id)

		seen := make(map[netip.AddrPort]struct{}, len(closest))
		for _, node := range closest {
			ret := t.cfg.dhts.GetPeers(ctx, dht.NewAddr(node.Addr.AddrPort), id, false)
			if ret.Err != nil {
				continue
			}

			for _, n := range torrentx.NodesFromReply(ret) {
				addr := n.Addr.AddrPort
				if _, ok := seen[addr]; ok {
					continue
				}
				seen[addr] = struct{}{}

				if qret := t.cfg.dhts.Query(ctx, dht.NewAddr(addr), searchreq); qret.Err != nil {
					log.Println("partition peer search query failed", addr, qret.Err)
				}
			}
		}
	}
}

func (t *partitionSeq) Err() error {
	return t.err
}
