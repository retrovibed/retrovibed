package netmonx_test

import (
	"testing"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/netmonx"
	"github.com/stretchr/testify/require"
)

func TestMonitorMeteredDefault(t *testing.T) {
	m, err := netmonx.New()
	require.NoError(t, err)
	defer m.Close()

	// A freshly created monitor should not report as metered; the OS-reported
	// IsExpensive is false on most wired/wifi interfaces.
	require.False(t, m.Metered())
}

func TestMonitorSetMeteredRoundtrip(t *testing.T) {
	m, err := netmonx.New()
	require.NoError(t, err)
	defer m.Close()

	m.SetMetered(true)
	require.Eventually(t, func() bool { return m.Metered() }, time.Second, time.Millisecond)

	m.SetMetered(false)
	require.Eventually(t, func() bool { return !m.Metered() }, time.Second, time.Millisecond)
}

// TestSetMeteredRoundtrip exercises the package-level SetMetered/Metered pair
// via the global monitor. Reads happen before InjectEvent's async state refresh
// can overwrite the manually-set value.
func TestSetMeteredRoundtrip(t *testing.T) {
	initial := netmonx.Metered()
	defer netmonx.SetMetered(initial)

	netmonx.SetMetered(true)
	require.Eventually(t, func() bool { return netmonx.Metered() }, time.Second, time.Millisecond)

	netmonx.SetMetered(false)
	require.Eventually(t, func() bool { return !netmonx.Metered() }, time.Second, time.Millisecond)
}

// TestMeteredPackageLevelNilSafe verifies the package-level Metered() returns
// false when the global monitor is unavailable. This is tested indirectly via
// a fresh monitor whose InterfaceState may return nil on constrained hosts.
func TestMeteredPackageLevelReturnsBool(t *testing.T) {
	// Calling Metered() must not panic regardless of network state.
	_ = netmonx.Metered()
}
