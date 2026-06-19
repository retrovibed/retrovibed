package tracking_test

import (
	"net/netip"
	"testing"
	"time"

	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestPeerUIDStableAcrossRealSecure(t *testing.T) {
	base := int160.Random()
	addrA := netip.MustParseAddr("1.2.3.4")
	addrB := netip.MustParseAddr("8.8.8.8")

	securedA := base.Secure(addrA)
	securedB := base.Secure(addrB)

	// sanity: the two addresses actually produce different secure prefixes
	require.NotEqual(t, securedA.Bytes()[:3], securedB.Bytes()[:3])

	require.Equal(t, tracking.PeerUID(base), tracking.PeerUID(securedA))
	require.Equal(t, tracking.PeerUID(base), tracking.PeerUID(securedB))
}

func TestPeerUIDChangesWithStableSuffix(t *testing.T) {
	var a [20]byte
	for i := range a {
		a[i] = byte(i)
	}
	b := a
	b[10] ^= 0xff // flip a byte within the stable suffix

	idA := int160.FromByteArray(a)
	idB := int160.FromByteArray(b)

	require.NotEqual(t, tracking.PeerUID(idA), tracking.PeerUID(idB))
}

func TestNewPeerIDStableAcrossSecurePrefix(t *testing.T) {
	var suffix [20]byte
	for i := 3; i < len(suffix); i++ {
		suffix[i] = byte(i)
	}

	a := suffix
	a[0], a[1], a[2] = 0x01, 0x02, 0x03
	b := suffix
	b[0], b[1], b[2] = 0xfe, 0xfd, 0xfc

	peerA := tracking.NewPeer(int160.FromByteArray(a))
	peerB := tracking.NewPeer(int160.FromByteArray(b))

	require.Equal(t, peerA.ID, peerB.ID)
	require.NotEqual(t, peerA.Peer, peerB.Peer)
}

// Regression: re-inserting a peer that already exists (same id, e.g. seen
// again with a new secure prefix/address) must upsert the latest 'peer',
// 'port', 'ip', and 'tombstoned_at' onto the existing row rather than
// leaving stale values.
func TestPeerInsertWithDefaultsUpsertsLatestPeerPortIP(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)

	var suffix [20]byte
	for i := 3; i < len(suffix); i++ {
		suffix[i] = byte(i)
	}
	original := suffix
	original[0], original[1], original[2] = 0x01, 0x02, 0x03
	updated := suffix
	updated[0], updated[1], updated[2] = 0xfe, 0xfd, 0xfc

	originalTombstone := time.Now().Add(1 * time.Hour).Truncate(time.Millisecond)
	updatedTombstone := time.Now().Add(48 * time.Hour).Truncate(time.Millisecond)

	originalPeer := tracking.NewPeer(
		int160.FromByteArray(original),
		tracking.PeerOptionIP(netip.MustParseAddrPort("1.2.3.4:6881")),
		tracking.PeerOptionTombstone(originalTombstone),
	)
	require.NoError(t, tracking.PeerInsertWithDefaults(ctx, q, originalPeer).Scan(&originalPeer))

	updatedPeer := tracking.NewPeer(
		int160.FromByteArray(updated),
		tracking.PeerOptionIP(netip.MustParseAddrPort("5.6.7.8:7000")),
		tracking.PeerOptionTombstone(updatedTombstone),
	)
	require.Equal(t, originalPeer.ID, updatedPeer.ID, "sanity: same logical peer")

	// capture the expected values before Scan overwrites updatedPeer with
	// whatever was actually persisted.
	expectedPeer, expectedPort, expectedIP := updatedPeer.Peer, updatedPeer.Port, updatedPeer.IP

	require.NoError(t, tracking.PeerInsertWithDefaults(ctx, q, updatedPeer).Scan(&updatedPeer))

	var stored tracking.Peer
	require.NoError(t, tracking.PeerFindByID(ctx, q, originalPeer.ID).Scan(&stored))

	require.Equal(t, expectedPeer, stored.Peer)
	require.Equal(t, expectedPort, stored.Port)
	require.Equal(t, expectedIP, stored.IP)
	require.True(t, updatedTombstone.Equal(stored.TombstonedAt), "expected %s, got %s", updatedTombstone, stored.TombstonedAt)
}
