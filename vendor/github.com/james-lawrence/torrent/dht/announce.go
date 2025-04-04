package dht

// get_peers and announce_peers.

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/anacrolix/chansync"
	"github.com/anacrolix/chansync/events"

	"github.com/james-lawrence/torrent/dht/int160"
	dhtutil "github.com/james-lawrence/torrent/dht/k-nearest-nodes"
	"github.com/james-lawrence/torrent/dht/krpc"
	"github.com/james-lawrence/torrent/dht/traversal"
)

// Maintains state for an ongoing Announce operation. An Announce is started by calling
// Server.Announce.
type Announce struct {
	Peers chan PeersValues

	server   *Server
	infoHash int160.T // Target

	announcePeerOpts *AnnouncePeerOpts
	scrape           bool
	peerAnnounced    chansync.SetOnce

	traversal *traversal.Operation

	closed chansync.SetOnce
}

func (a *Announce) String() string {
	return fmt.Sprintf("%[1]T %[1]p of %v on %v", a, a.infoHash, a.server)
}

// Returns the number of distinct remote addresses the announce has queried.
func (a *Announce) NumContacted() uint32 {
	return atomic.LoadUint32(&a.traversal.Stats().NumAddrsTried)
}

func (a *Announce) TraversalStats() TraversalStats {
	return *a.traversal.Stats()
}

// Server.Announce option
type AnnounceOpt func(a *Announce)

// Scrape BEP 33 bloom filters in queries.
func Scrape() AnnounceOpt {
	return func(a *Announce) {
		a.scrape = true
	}
}

// Arguments for announce_peer from a Server.Announce.
type AnnouncePeerOpts struct {
	// The peer port that we're announcing.
	Port int
	// The peer port should be determined by the receiver to be the source port of the query packet.
	ImpliedPort bool
}

// Finish an Announce get_peers traversal with an announce of a local peer.
func AnnouncePeer(implied bool, port int) AnnounceOpt {
	return func(a *Announce) {
		if port == 0 && !implied {
			return
		}

		a.announcePeerOpts = &AnnouncePeerOpts{
			Port:        port,
			ImpliedPort: implied,
		}
	}
}

// Traverses the DHT graph toward nodes that store peers for the infohash, streaming them to the
// caller.
func (s *Server) AnnounceTraversal(ctx context.Context, infoHash [20]byte, opts ...AnnounceOpt) (_ *Announce, err error) {
	a := &Announce{
		Peers:    make(chan PeersValues),
		server:   s,
		infoHash: int160.FromByteArray(infoHash),
	}
	for _, opt := range opts {
		opt(a)
	}
	a.traversal = traversal.Start(traversal.OperationInput{
		Target:     infoHash,
		DoQuery:    a.getPeers,
		NodeFilter: s.TraversalNodeFilter,
		DataFilter: func(data any) bool {
			_, ok := data.(string)
			return ok
		},
	})
	nodes, err := s.TraversalStartingNodes()
	if err != nil {
		a.traversal.Stop()
		return
	}
	a.traversal.AddNodes(nodes)
	go func() {
		select {
		case <-a.traversal.Stalled():
		case <-ctx.Done():
		}

		a.traversal.Stop()
		<-a.traversal.Stopped()
		if a.announcePeerOpts != nil {
			a.announceClosest(ctx)
		}
		a.peerAnnounced.Set()
		close(a.Peers)
	}()

	return a, nil
}

func (a *Announce) announceClosest(ctx context.Context) {
	var wg sync.WaitGroup
	a.traversal.Closest().Range(func(elem dhtutil.Elem) {
		wg.Add(1)
		go func() {
			a.logger().Printf("announce_peer to %v: %v\n",
				elem, a.announcePeer(ctx, elem),
			)
			wg.Done()
		}()
	})
	wg.Wait()
}

func (a *Announce) announcePeer(ctx context.Context, peer dhtutil.Elem) error {
	ctx, done := context.WithTimeout(ctx, 5*time.Second)
	defer done()

	go func() {
		select {
		case <-a.closed.Done():
		case <-ctx.Done():
		}
	}()
	return a.server.announcePeer(
		ctx,
		NewAddr(peer.Addr.UDP()),
		a.infoHash,
		a.announcePeerOpts.Port,
		peer.Data.(string),
		a.announcePeerOpts.ImpliedPort,
	).Err
}

func (a *Announce) getPeers(ctx context.Context, addr krpc.NodeAddr) traversal.QueryResult {
	res := a.server.GetPeers(ctx, NewAddr(addr.UDP()), a.infoHash, a.scrape, QueryRateLimiting{})
	if r := res.Reply.R; r != nil {
		peersValues := PeersValues{
			Peers: r.Values,
			NodeInfo: krpc.NodeInfo{
				Addr: addr,
				ID:   r.ID,
			},
			Return: *r,
		}
		select {
		case a.Peers <- peersValues:
		case <-a.traversal.Stopped():
		}
	}
	return res.TraversalQueryResult(addr)
}

// Corresponds to the "values" key in a get_peers KRPC response. A list of
// peers that a node has reported as being in the swarm for a queried info
// hash.
type PeersValues struct {
	Peers         []Peer // Peers given in get_peers response.
	krpc.NodeInfo        // The node that gave the response.
	krpc.Return
}

// Stop the announce.
func (a *Announce) Close() {
	a.StopTraversing()
	// This will prevent peer announces from proceeding.
	a.closed.Set()
}

func (a *Announce) logger() *log.Logger {
	return a.server.logger()
}

// Halts traversal, but won't block peer announcing.
func (a *Announce) StopTraversing() {
	a.traversal.Stop()
}

// Traversal and peer announcing steps are done.
func (a *Announce) Finished() events.Done {
	// This is the last step in an announce.
	return a.peerAnnounced.Done()
}
