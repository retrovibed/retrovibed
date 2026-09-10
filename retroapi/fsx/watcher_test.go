package fsx_test

import (
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/retrovibed/retrovibed/retroapi/fsx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
)

// waitBroadcast parks a waiter on cond, runs trigger once the waiter is
// parked, and reports whether a broadcast arrived within timeout.
func waitBroadcast(t *testing.T, cond *sync.Cond, timeout time.Duration, trigger func()) bool {
	t.Helper()

	done := make(chan struct{})

	go func() {
		cond.L.Lock()
		cond.Wait()
		cond.L.Unlock()
		close(done)
	}()

	// Give the waiter time to park so the broadcast is not missed.
	time.Sleep(25 * time.Millisecond)

	if trigger != nil {
		trigger()
	}

	select {
	case <-done:
		return true
	case <-time.After(timeout):
		return false
	}
}

func TestWatch(t *testing.T) {
	t.Run("broadcasts when a new file is created in the watched directory", func(t *testing.T) {
		ctx, cancel := testx.Context(t)
		defer cancel()

		dir := t.TempDir()
		cond := sync.NewCond(&sync.Mutex{})

		require.NoError(t, fsx.Watch(ctx, cond, dir))

		ok := waitBroadcast(t, cond, 3*time.Second, func() {
			require.NoError(t, fsx.Touch(0600, filepath.Join(dir, "newfile")))
		})
		require.True(t, ok, "expected a broadcast after touching a new file in the watched directory")
	})

	t.Run("broadcasts for events in each of multiple watched directories", func(t *testing.T) {
		ctx, cancel := testx.Context(t)
		defer cancel()

		dirA := t.TempDir()
		dirB := t.TempDir()
		cond := sync.NewCond(&sync.Mutex{})

		require.NoError(t, fsx.Watch(ctx, cond, dirA, dirB))

		ok := waitBroadcast(t, cond, 3*time.Second, func() {
			require.NoError(t, fsx.Touch(0600, filepath.Join(dirB, "newfile")))
		})
		require.True(t, ok, "expected a broadcast after touching a new file in the second watched directory")
	})

	t.Run("touching an existing file does not broadcast", func(t *testing.T) {
		ctx, cancel := testx.Context(t)
		defer cancel()

		dir := t.TempDir()
		path := filepath.Join(dir, "existing")
		require.NoError(t, fsx.Touch(0600, path))

		cond := sync.NewCond(&sync.Mutex{})
		require.NoError(t, fsx.Watch(ctx, cond, dir))

		ok := waitBroadcast(t, cond, 500*time.Millisecond, func() {
			require.NoError(t, fsx.Touch(0600, path))
		})
		require.False(t, ok, "expected no broadcast when re-touching an existing file")
	})

	t.Run("watching a non-existent path does not fail", func(t *testing.T) {
		ctx, cancel := testx.Context(t)
		defer cancel()

		dir := t.TempDir()
		cond := sync.NewCond(&sync.Mutex{})

		require.NoError(t, fsx.Watch(ctx, cond, filepath.Join(dir, "does-not-exist")))
	})

	t.Run("stops broadcasting after the context is cancelled", func(t *testing.T) {
		ctx, cancel := testx.Context(t)

		dir := t.TempDir()
		cond := sync.NewCond(&sync.Mutex{})

		require.NoError(t, fsx.Watch(ctx, cond, dir))

		cancel()
		// Allow the watcher goroutine to exit and close the watcher.
		time.Sleep(150 * time.Millisecond)

		ok := waitBroadcast(t, cond, 500*time.Millisecond, func() {
			require.NoError(t, fsx.Touch(0600, filepath.Join(dir, "newfile")))
		})
		require.False(t, ok, "expected no broadcast after the context was cancelled")
	})
}
