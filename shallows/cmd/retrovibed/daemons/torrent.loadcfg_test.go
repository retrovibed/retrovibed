package daemons

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/stretchr/testify/require"
)

func TestLoadcfg(t *testing.T) {
	t.Run("generates and round-trips config when file does not exist", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "torrent.cfg")
		cfg := &TorrentSettings{}
		require.NoError(t, testx.Fake(cfg))
		expected := cfg.MaximumRequests

		tr := _torrenting{cfgpath: path}
		require.NoError(t, tr.loadcfg(path, cfg))

		_, err := os.Stat(path)
		require.NoError(t, err, "config file should have been created")
		require.Equal(t, expected, cfg.MaximumRequests)
	})

	t.Run("reads existing config and preserves uint64 maximum_requests", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "torrent.cfg")

		src := &TorrentSettings{}
		require.NoError(t, testx.Fake(src))

		tr := _torrenting{cfgpath: path}
		require.NoError(t, tr.loadcfg(path, src))

		dst := &TorrentSettings{}
		require.NoError(t, tr.loadcfg(path, dst))
		require.Equal(t, src.MaximumRequests, dst.MaximumRequests)
	})

	t.Run("sub-message fields retain defaults when not specified in config", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "torrent.cfg")
		require.NoError(t, os.WriteFile(path, []byte(`{}`), 0600))

		tr := _torrenting{cfgpath: path}
		dst := &TorrentSettings{
			Peers:    &Peers{Min: 32, Max: 64},
			Upload:   &Limit{Rate: 100, Burst: 10},
			Download: &Limit{Rate: 200, Burst: 20},
			Inbound:  &Limit{Rate: 300, Burst: 30},
			Outbound: &Limit{Rate: 400, Burst: 40},
		}
		require.NoError(t, tr.loadcfg(path, dst))
		require.NotNil(t, dst.Peers, "Peers should not be nil after loading empty config")
		require.NotNil(t, dst.Upload, "Upload should not be nil after loading empty config")
		require.NotNil(t, dst.Download, "Download should not be nil after loading empty config")
		require.NotNil(t, dst.Inbound, "Inbound should not be nil after loading empty config")
		require.NotNil(t, dst.Outbound, "Outbound should not be nil after loading empty config")
	})

	t.Run("reads pre-existing protojson config with uint64 encoded as string", func(t *testing.T) {
		// Regression: protojson encodes uint64 as a quoted string which standard
		// encoding/json cannot unmarshal into uint64.
		path := filepath.Join(t.TempDir(), "torrent.cfg")
		require.NoError(t, os.WriteFile(path, []byte(`{"maximum_requests":"1024"}`), 0600))

		tr := _torrenting{cfgpath: path}
		dst := &TorrentSettings{}
		require.NoError(t, tr.loadcfg(path, dst))
		require.Equal(t, uint64(1024), dst.MaximumRequests)
	})
}
