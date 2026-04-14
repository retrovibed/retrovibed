package ddiscapi_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/internal/grpcx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPeerToTrackingPeerCodex(t *testing.T) {
	compare := func(t *testing.T, mp tracking.Peer, p *ddiscapi.Peer) {
		assert.Equal(t, mp.ID, p.Id)
		assert.Equal(t, grpcx.EncodeTime(mp.CreatedAt), p.CreatedAt)
		assert.Equal(t, mp.Description, p.Description)
		assert.Equal(t, mp.Peer, p.Infohash)
		assert.Equal(t, mp.Ddisc, p.Ddisc)
		assert.Equal(t, mp.DdiscPartition, p.Partition)
		assert.Equal(t, mp.DdiscSyncoffset, p.Syncoffset)
	}

	t.Run("encode", func(t *testing.T) {
		mp := tracking.Peer{}
		require.NoError(t, testx.Fake(&mp, tracking.PeerOptionTestDefaults))

		p, err := ddiscapi.NewPeerFromTrackingPeer(mp)
		require.NoError(t, err)
		require.NotNil(t, p)
		compare(t, mp, p)
	})

	t.Run("decode", func(t *testing.T) {
		_p := tracking.Peer{}
		require.NoError(t, testx.Fake(&_p, tracking.PeerOptionTestDefaults))

		mp, err := ddiscapi.NewPeerFromTrackingPeer(_p)
		require.NoError(t, err)

		p, err := ddiscapi.NewTrackingPeerFromPeer(mp)
		require.NoError(t, err)
		require.NotNil(t, p)
		compare(t, p, mp)
	})
}
