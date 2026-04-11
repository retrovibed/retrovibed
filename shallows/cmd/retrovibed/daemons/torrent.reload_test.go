package daemons

import (
	"context"
	"testing"
	"time"

	"github.com/retrovibed/retrovibed/internal/sqltestx"
	"github.com/stretchr/testify/require"
)

func TestReload(t *testing.T) {
	t.Run("returns nil immediately", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)
		tr := newTestTorrenting(t, q)

		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()

		cfg := AutoTorrentSettings(&TorrentSettings{})
		disc := &DiscoverySettings{}

		require.NoError(t, tr.Reload(ctx, cfg, disc))
	})

	t.Run("initializes client asynchronously", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)
		tr := newTestTorrenting(t, q)

		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()

		cfg := AutoTorrentSettings(&TorrentSettings{})
		disc := &DiscoverySettings{}

		require.NoError(t, tr.Reload(ctx, cfg, disc))

		require.Eventually(t, func() bool {
			return tr._tclient.Load() != nil
		}, 10*time.Second, 50*time.Millisecond)
	})

	t.Run("reinitializes client on broadcast", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)
		tr := newTestTorrenting(t, q)

		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()

		cfg := AutoTorrentSettings(&TorrentSettings{})
		disc := &DiscoverySettings{}

		require.NoError(t, tr.Reload(ctx, cfg, disc))

		// wait for first init to complete
		require.Eventually(t, func() bool {
			return tr._tclient.Load() != nil
		}, 10*time.Second, 50*time.Millisecond)

		first := tr._tclient.Load()

		// keep broadcasting every tick until the goroutine wakes and re-inits
		require.Eventually(t, func() bool {
			tr.cond.Broadcast()
			return tr._tclient.Load() != first
		}, 10*time.Second, 50*time.Millisecond)

		require.NotNil(t, tr._tclient.Load())
	})

	t.Run("stops reloading after context is cancelled", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)
		tr := newTestTorrenting(t, q)

		ctx, cancel := context.WithCancel(t.Context())

		cfg := AutoTorrentSettings(&TorrentSettings{})
		disc := &DiscoverySettings{}

		require.NoError(t, tr.Reload(ctx, cfg, disc))

		// wait for first init to complete
		require.Eventually(t, func() bool {
			return tr._tclient.Load() != nil
		}, 10*time.Second, 50*time.Millisecond)

		first := tr._tclient.Load()

		// cancel context then keep broadcasting — the loop should exit without re-init
		cancel()
		require.Never(t, func() bool {
			tr.cond.Broadcast()
			return tr._tclient.Load() != first
		}, time.Second, 50*time.Millisecond)
	})
}
