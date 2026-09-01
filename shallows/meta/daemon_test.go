package meta_test

import (
	"net"
	"os"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/stretchr/testify/require"
)

func TestDaemonFromHostPort(t *testing.T) {
	d := meta.DaemonFromHost()

	host, port, err := net.SplitHostPort(d.Hostname)
	require.NoError(t, err)
	require.Equal(t, "9998", port)

	hostname, err := os.Hostname()
	require.NoError(t, err)
	require.Equal(t, hostname, host)
}

func TestDaemonDownload(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()
	q := sqltestx.Metadatabase(t)

	var (
		current meta.Daemon
		last    meta.Daemon
	)

	for i := range 3 {
		require.NoError(t, testx.Fake(&last, meta.DaemonOptionTestDefaults, meta.DaemonOptionMaybeID, func(d *meta.Daemon) {
			d.Downloads = i == 0
		}))
		require.NoError(t, meta.DaemonInsertWithDefaults(ctx, q, last).Scan(&last))
	}

	require.NoError(t, meta.DaemonFindByDownload(ctx, q).Scan(&current))
	require.NotEqualValues(t, current.ID, last.ID)
	require.False(t, last.Downloads)
	require.EqualValues(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM meta_daemons WHERE downloads"))
	require.NoError(t, meta.DaemonDownload(ctx, q, last.ID).Scan(&last))
	require.True(t, last.Downloads)
	require.NoError(t, meta.DaemonFindByDownload(ctx, q).Scan(&current))
	require.EqualValues(t, current.ID, last.ID)
	require.EqualValues(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM meta_daemons WHERE downloads"))
}

func TestDaemonTouch(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()
	q := sqltestx.Metadatabase(t)

	var (
		current meta.Daemon
		last    meta.Daemon
	)

	for i := range 3 {
		require.NoError(t, testx.Fake(&last, meta.DaemonOptionTestDefaults, meta.DaemonOptionMaybeID, func(d *meta.Daemon) {
			d.Default = i == 0
		}))
		require.NoError(t, meta.DaemonInsertWithDefaults(ctx, q, last).Scan(&last))
	}

	require.NoError(t, meta.DaemonFindDefault(ctx, q).Scan(&current))
	require.NotEqualValues(t, current.ID, last.ID)
	require.False(t, last.Default)
	require.EqualValues(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM meta_daemons WHERE \"default\""))
	require.NoError(t, meta.DaemonTouch(ctx, q, last.ID).Scan(&last))
	require.True(t, last.Default)
	require.NoError(t, meta.DaemonFindDefault(ctx, q).Scan(&current))
	require.EqualValues(t, current.ID, last.ID)
	require.EqualValues(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM meta_daemons WHERE \"default\""))
}
