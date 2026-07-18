package daemons_test

import (
	"net"
	"testing"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/james-lawrence/torrent/bep0051"
	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/retrovibe/daemons"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/torrentx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

// newDHTPair starts a querying server `s` and a target server that answers
// bep0051 sample requests with the provided fakeSampler, returning both
// along with the target's dht id.
func newDHTInfoHashesPair(t *testing.T, sampler fakeSampler) (s *dht.Server, target *dht.Server, targetID int160.T) {
	t.Helper()

	s, err := dht.NewServer(32)
	require.NoError(t, err)
	t.Cleanup(s.Close)
	spc, err := net.ListenPacket("udp4", "127.0.0.1:0")
	require.NoError(t, err)
	require.NoError(t, s.Serve(t.Context(), spc))

	target, err = dht.NewServer(32, dht.OptionMuxer(
		dht.DefaultMuxer().Method(bep0051.Query, bep0051.NewEndpoint(sampler)),
	))
	require.NoError(t, err)
	t.Cleanup(target.Close)
	targetpc, err := net.ListenPacket("udp4", "127.0.0.1:0")
	require.NoError(t, err)
	require.NoError(t, target.Serve(t.Context(), targetpc))

	targetID = target.ID(target.DynamicAddrPort())

	return s, target, targetID
}

func TestDiscoverDHTInfoHashes(t *testing.T) {
	t.Run("records a previously unseen infohash sampled from a peer", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		infohash := int160.Random()

		s, target, targetID := newDHTInfoHashesPair(t, fakeSampler{
			total:  1,
			sample: infohash.Bytes(),
		})

		peer := tracking.NewPeer(targetID, tracking.PeerOptionIP(target.DynamicAddrPort()), tracking.PeerOptionBEP51(1, 600))
		require.NoError(t, tracking.PeerInsertWithDefaults(ctx, q, peer).Scan(&peer))

		go func() {
			_ = daemons.DiscoverDHTInfoHashes(ctx, q, s)
		}()

		expectedID := torrentx.HashUID(&infohash)
		sql, args, err := tracking.UnknownSearchBuilder().RemoveColumns().Columns("COUNT(*)").Where(squirrel.Eq{"id": expectedID}).ToSql()
		require.NoError(t, err)

		require.Eventually(t, func() bool {
			return sqltestx.Count(t, q, sql, args...) == 1
		}, 10*time.Second, 100*time.Millisecond, "expected sampled infohash to be recorded as unknown")
	})

	t.Run("does not record an infohash that is already known metadata", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		infohash := int160.Random()

		known := tracking.NewMetadata(&infohash)
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, known).Scan(&known))

		s, target, targetID := newDHTInfoHashesPair(t, fakeSampler{
			total:  1,
			sample: infohash.Bytes(),
		})

		peer := tracking.NewPeer(targetID, tracking.PeerOptionIP(target.DynamicAddrPort()), tracking.PeerOptionBEP51(1, 600))
		require.NoError(t, tracking.PeerInsertWithDefaults(ctx, q, peer).Scan(&peer))

		go func() {
			_ = daemons.DiscoverDHTInfoHashes(ctx, q, s)
		}()

		expectedID := torrentx.HashUID(&infohash)
		sql, args, err := tracking.UnknownSearchBuilder().RemoveColumns().Columns("COUNT(*)").Where(squirrel.Eq{"id": expectedID}).ToSql()
		require.NoError(t, err)

		require.Never(t, func() bool {
			return sqltestx.Count(t, q, sql, args...) == 1
		}, 2*time.Second, 100*time.Millisecond, "already known infohash should not be tracked as unknown")
	})

	t.Run("does not record an infohash that has already been discovered", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		infohash := int160.Random()

		discovered := ddisc.NewDiscovered(&infohash, ddisc.DiscoveredOptionAutoMagnet)
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, discovered).Scan(&discovered))

		s, target, targetID := newDHTInfoHashesPair(t, fakeSampler{
			total:  1,
			sample: infohash.Bytes(),
		})

		peer := tracking.NewPeer(targetID, tracking.PeerOptionIP(target.DynamicAddrPort()), tracking.PeerOptionBEP51(1, 600))
		require.NoError(t, tracking.PeerInsertWithDefaults(ctx, q, peer).Scan(&peer))

		go func() {
			_ = daemons.DiscoverDHTInfoHashes(ctx, q, s)
		}()

		expectedID := torrentx.HashUID(&infohash)
		sql, args, err := tracking.UnknownSearchBuilder().RemoveColumns().Columns("COUNT(*)").Where(squirrel.Eq{"id": expectedID}).ToSql()
		require.NoError(t, err)

		require.Never(t, func() bool {
			return sqltestx.Count(t, q, sql, args...) == 1
		}, 2*time.Second, 100*time.Millisecond, "already discovered infohash should not be tracked as unknown")
	})

	t.Run("pushes the peer's next check into the future after a successful sample", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		infohash := int160.Random()

		s, target, targetID := newDHTInfoHashesPair(t, fakeSampler{
			total:  3,
			sample: infohash.Bytes(),
		})

		peer := tracking.NewPeer(targetID, tracking.PeerOptionIP(target.DynamicAddrPort()), tracking.PeerOptionBEP51(1, 600))
		require.NoError(t, tracking.PeerInsertWithDefaults(ctx, q, peer).Scan(&peer))
		originalNextCheck := peer.NextCheck

		go func() {
			_ = daemons.DiscoverDHTInfoHashes(ctx, q, s)
		}()

		require.Eventually(t, func() bool {
			var updated tracking.Peer
			if err := tracking.PeerFindByID(ctx, q, peer.ID).Scan(&updated); err != nil {
				return false
			}
			return updated.NextCheck.After(originalNextCheck)
		}, 10*time.Second, 100*time.Millisecond, "expected peer's next check to advance after a successful sample")
	})
}
