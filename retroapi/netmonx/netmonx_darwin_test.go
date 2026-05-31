//go:build darwin

package netmonx

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDarwinWatcherStarts(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	notify := make(chan struct{}, 1)
	err := startWatcher(ctx, notify)
	require.NoError(t, err)
}

func TestDarwinDefaultRouteInterface(t *testing.T) {
	iface := defaultRouteInterface()
	require.NotEmpty(t, iface, "expected a default route interface on macOS")
}

func TestDarwinMonitorInitialState(t *testing.T) {
	m, err := New()
	require.NoError(t, err)
	defer m.Close()

	require.NoError(t, m.Err())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var got *State
	for delta := range m.Each(ctx) {
		got = delta.New
		break
	}
	require.NotNil(t, got)
	require.NotEmpty(t, got.Interfaces, "expected at least one interface with IPs")
}
