package wireguardx

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// freshState returns an autohealState with lastRecovery backdated so the
// suppress guard does not interfere with the test scenario.
func freshState() autohealState {
	s := newAutohealState()
	s.lastRecovery = time.Now().Add(-1 * time.Hour)
	return s
}

func recentHandshake() int64 {
	return time.Now().Add(-10 * time.Second).Unix()
}

func staleHandshake(s autohealState) int64 {
	return time.Now().Add(-(s.handshakeExpiry + time.Minute)).Unix()
}

func TestNeedsRecovery(t *testing.T) {
	t.Run("no peer", func(t *testing.T) {
		state := freshState()
		curr := Statistics{PeerKey: "", LastHandshakeSec: recentHandshake(), KeepaliveInterval: 25}
		require.False(t, state.needsRecovery(curr))
	})

	t.Run("no handshake ever", func(t *testing.T) {
		state := freshState()
		curr := Statistics{PeerKey: "abc", LastHandshakeSec: 0, KeepaliveInterval: 25}
		require.False(t, state.needsRecovery(curr))
	})

	t.Run("suppressed: handshake older than lastRecovery", func(t *testing.T) {
		state := newAutohealState() // lastRecovery = time.Now()
		curr := Statistics{
			PeerKey:           "abc",
			KeepaliveInterval: 25,
			LastHandshakeSec:  time.Now().Add(-30 * time.Second).Unix(),
		}
		require.False(t, state.needsRecovery(curr))
	})

	t.Run("healthy: fresh handshake within expiry with keepalive", func(t *testing.T) {
		state := freshState()
		curr := Statistics{
			PeerKey:           "abc",
			KeepaliveInterval: 25,
			LastHandshakeSec:  recentHandshake(),
			TXBytes:           100,
			RXBytes:           100,
		}
		require.False(t, state.needsRecovery(curr))
	})

	t.Run("stale handshake with keepalive", func(t *testing.T) {
		state := freshState()
		curr := Statistics{
			PeerKey:           "abc",
			KeepaliveInterval: 25,
			LastHandshakeSec:  staleHandshake(state),
		}
		require.True(t, state.needsRecovery(curr))
	})

	t.Run("stale handshake without keepalive", func(t *testing.T) {
		state := freshState()
		curr := Statistics{
			PeerKey:           "abc",
			KeepaliveInterval: 0,
			LastHandshakeSec:  staleHandshake(state),
		}
		require.False(t, state.needsRecovery(curr))
	})

	t.Run("unbalanced pipe with baseline", func(t *testing.T) {
		state := freshState()
		state.prev = Statistics{TXBytes: 100, RXBytes: 500}
		curr := Statistics{
			PeerKey:          "abc",
			LastHandshakeSec: recentHandshake(),
			TXBytes:          200,
			RXBytes:          500,
		}
		require.True(t, state.needsRecovery(curr))
	})

	t.Run("unbalanced pipe without baseline", func(t *testing.T) {
		state := freshState()
		state.prev = Statistics{TXBytes: 0, RXBytes: 0}
		curr := Statistics{
			PeerKey:          "abc",
			LastHandshakeSec: recentHandshake(),
			TXBytes:          100,
			RXBytes:          0,
		}
		require.False(t, state.needsRecovery(curr))
	})
}
