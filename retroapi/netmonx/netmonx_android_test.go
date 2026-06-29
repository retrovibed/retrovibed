//go:build android

package netmonx

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAndroidWatcherUnavailable(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	notify := make(chan struct{}, 1)
	err := startWatcher(ctx, notify)
	require.Error(t, err, "expected native watcher to be unavailable on android")

	// New() must still succeed via polling fallback.
	m, err := New()
	require.NoError(t, err)
	defer m.Close()
	require.NoError(t, m.Err())
}

func TestAndroidDefaultRouteInterface(t *testing.T) {
	iface := defaultRouteInterface()
	require.NotEmpty(t, iface, "expected a default route interface on android")
}
