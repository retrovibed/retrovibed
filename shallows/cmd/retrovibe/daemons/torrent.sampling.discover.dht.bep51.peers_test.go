package daemons_test

import (
	"net"
	"testing"
	"time"

	"github.com/james-lawrence/torrent/bep0051"
	"github.com/james-lawrence/torrent/dht"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/retrovibe/daemons"
	"github.com/retrovibed/retrovibed/shallows/internal/int160x"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

// fakeSampler implements bep0051.Sampler with a fixed, canned response.
type fakeSampler struct {
	total  uint
	sample []byte
}

func (f fakeSampler) Snapshot(max int) (ttl uint, total uint, sample []byte) {
	return 600, f.total, f.sample
}

func TestDiscoverDHTBEP51Peers(t *testing.T) {
	t.Run("records a peer discovered through the dht that reports available bep51 hashes", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		// a discovers relay via the dht; relay knows about target and will
		// surface it in a get_peers reply. target answers bep51 sample
		// queries directly with a non-empty sample.
		a, err := dht.NewServer(32)
		require.NoError(t, err)
		t.Cleanup(a.Close)
		apc, err := net.ListenPacket("udp4", "127.0.0.1:0")
		require.NoError(t, err)
		require.NoError(t, a.Serve(t.Context(), apc))

		relay, err := dht.NewServer(32)
		require.NoError(t, err)
		t.Cleanup(relay.Close)
		relaypc, err := net.ListenPacket("udp4", "127.0.0.1:0")
		require.NoError(t, err)
		require.NoError(t, relay.Serve(t.Context(), relaypc))

		target, err := dht.NewServer(32, dht.OptionMuxer(
			dht.DefaultMuxer().Method(bep0051.Query, bep0051.NewEndpoint(fakeSampler{
				total:  1,
				sample: make([]byte, 20),
			})),
		))
		require.NoError(t, err)
		t.Cleanup(target.Close)
		targetpc, err := net.ListenPacket("udp4", "127.0.0.1:0")
		require.NoError(t, err)
		require.NoError(t, target.Serve(t.Context(), targetpc))

		require.NoError(t, a.Ping(relay.DynamicAddrPort()).Err)
		require.NoError(t, relay.Ping(target.DynamicAddrPort()).Err)

		// Targeting the relay's own id trips the dht's "target == root"
		// special case in its bucket scan, forcing a full table scan
		// instead of one bounded by the query's bucket index. That makes
		// the relay deterministically surface every good node it knows
		// about - including target - regardless of where random node ids
		// happen to land relative to each other.
		relayID := relay.ID(relay.DynamicAddrPort())

		go func() {
			_ = daemons.DiscoverDHTBEP51Peers(ctx, q, a, int160x.NewRangeFixed(relayID))
		}()

		targetPeerID := tracking.PeerUID(target.ID(target.DynamicAddrPort()))

		require.Eventually(t, func() bool {
			var peer tracking.Peer
			if err := tracking.PeerFindByID(ctx, q, targetPeerID).Scan(&peer); err != nil {
				return false
			}
			return peer.Bep51 && peer.Bep51Available > 0
		}, 10*time.Second, 100*time.Millisecond, "expected target peer to be recorded as interesting")
	})

	t.Run("does not record a peer reporting zero available hashes", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		a, err := dht.NewServer(32)
		require.NoError(t, err)
		t.Cleanup(a.Close)
		apc, err := net.ListenPacket("udp4", "127.0.0.1:0")
		require.NoError(t, err)
		require.NoError(t, a.Serve(t.Context(), apc))

		relay, err := dht.NewServer(32)
		require.NoError(t, err)
		t.Cleanup(relay.Close)
		relaypc, err := net.ListenPacket("udp4", "127.0.0.1:0")
		require.NoError(t, err)
		require.NoError(t, relay.Serve(t.Context(), relaypc))

		target, err := dht.NewServer(32, dht.OptionMuxer(
			dht.DefaultMuxer().Method(bep0051.Query, bep0051.NewEndpoint(fakeSampler{
				total:  0,
				sample: []byte{},
			})),
		))
		require.NoError(t, err)
		t.Cleanup(target.Close)
		targetpc, err := net.ListenPacket("udp4", "127.0.0.1:0")
		require.NoError(t, err)
		require.NoError(t, target.Serve(t.Context(), targetpc))

		require.NoError(t, a.Ping(relay.DynamicAddrPort()).Err)
		require.NoError(t, relay.Ping(target.DynamicAddrPort()).Err)

		relayID := relay.ID(relay.DynamicAddrPort())

		go func() {
			_ = daemons.DiscoverDHTBEP51Peers(ctx, q, a, int160x.NewRangeFixed(relayID))
		}()

		targetPeerID := tracking.PeerUID(target.ID(target.DynamicAddrPort()))

		// give the discovery loop a chance to run before asserting the
		// negative; it runs its first iteration immediately on start.
		require.Never(t, func() bool {
			var peer tracking.Peer
			return tracking.PeerFindByID(ctx, q, targetPeerID).Scan(&peer) == nil
		}, 2*time.Second, 100*time.Millisecond, "peer with zero available hashes should not be recorded")
	})
}
