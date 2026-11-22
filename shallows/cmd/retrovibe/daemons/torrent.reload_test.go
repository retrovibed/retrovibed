package daemons

import (
	"context"
	"testing"
	"time"

	"github.com/james-lawrence/torrent"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/stretchr/testify/require"
)

func TestReload(t *testing.T) {
	t.Run("returns nil immediately", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)
		tr := newTestTorrenting(t, q)

		ctx, cancel := context.WithCancel(t.Context())

		cfg := AutoTorrentSettings(&TorrentSettings{})
		disc := &DiscoverySettings{}

		require.NoError(t, tr.Reload(ctx, cfg, &DHTSettings{}, disc))
		defer func() {
			cancel()
			<-tr.Stopped()
		}()
	})

	t.Run("initializes client asynchronously", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)
		tr := newTestTorrenting(t, q)

		ctx, cancel := context.WithCancel(t.Context())

		cfg := AutoTorrentSettings(&TorrentSettings{})
		disc := &DiscoverySettings{}

		require.NoError(t, tr.Reload(ctx, cfg, &DHTSettings{}, disc))
		defer func() {
			cancel()
			<-tr.Stopped()
		}()

		require.Eventually(t, func() bool {
			return tr._tclient.Load() != nil
		}, 10*time.Second, 50*time.Millisecond)
	})

	t.Run("reinitializes client on broadcast", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)
		tr := newTestTorrenting(t, q)

		ctx, cancel := context.WithCancel(t.Context())

		cfg := AutoTorrentSettings(&TorrentSettings{})
		disc := &DiscoverySettings{}

		require.NoError(t, tr.Reload(ctx, cfg, &DHTSettings{}, disc))
		defer func() {
			cancel()
			<-tr.Stopped()
		}()

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

	t.Run("a single broadcast during init is not lost", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)
		tr := newTestTorrenting(t, q)

		ctx, cancel := context.WithCancel(t.Context())

		cfg := AutoTorrentSettings(&TorrentSettings{})
		disc := &DiscoverySettings{}

		require.NoError(t, tr.Reload(ctx, cfg, &DHTSettings{}, disc))
		defer func() {
			cancel()
			<-tr.Stopped()
		}()

		// let the listener goroutine reach cond.Wait, then signal exactly once while
		// the loop is (most likely) still inside the first Init. a bare cond.Wait
		// loop drops this signal; the latched channel must keep it.
		time.Sleep(100 * time.Millisecond)
		tr.cond.Broadcast()

		var first *torrent.Client
		require.Eventually(t, func() bool {
			c := tr._tclient.Load()
			if first == nil {
				first = c
				return false
			}
			return c != nil && c != first
		}, 20*time.Second, 50*time.Millisecond)
	})

	t.Run("stops reloading after context is cancelled", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)
		tr := newTestTorrenting(t, q)

		ctx, cancel := context.WithCancel(t.Context())

		cfg := AutoTorrentSettings(&TorrentSettings{})
		disc := &DiscoverySettings{}

		require.NoError(t, tr.Reload(ctx, cfg, &DHTSettings{}, disc))

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

		<-tr.Stopped()
	})
}
