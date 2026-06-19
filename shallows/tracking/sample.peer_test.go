package tracking_test

import (
	"database/sql"
	"net/netip"
	"testing"

	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestSamplePeerAlwaysInsertsAtRateOne(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)

	id := int160.Random()
	addr := netip.MustParseAddrPort("1.2.3.4:6881")

	require.NoError(t, tracking.SamplePeer(ctx, q, 1, id, addr))

	var stored tracking.Peer
	require.NoError(t, tracking.PeerFindByID(ctx, q, tracking.PeerUID(id)).Scan(&stored))
	require.Equal(t, addr.Addr(), stored.IP)
	require.Equal(t, addr.Port(), stored.Port)
}

func TestSamplePeerNeverInsertsAtRateZero(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)

	id := int160.Random()
	addr := netip.MustParseAddrPort("1.2.3.4:6881")

	require.NoError(t, tracking.SamplePeer(ctx, q, 0, id, addr))

	var stored tracking.Peer
	require.ErrorIs(t, tracking.PeerFindByID(ctx, q, tracking.PeerUID(id)).Scan(&stored), sql.ErrNoRows)
}
