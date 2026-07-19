package ddisc_test

import (
	"testing"
	"time"

	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

// Regression: only search-plugin-derived candidates (unconfirmed real
// infohash, see NewDiscoveredFromImport) should expire on their own -
// candidates built from a confirmed infohash are kept indefinitely.
func TestDiscoveredTombstonedAt(t *testing.T) {
	id := int160.Random()
	require.True(t, ddisc.NewDiscovered(&id).TombstonedAt.Equal(timex.Inf()))
	require.True(t, ddisc.NewDiscoveredFromKnown(id, library.Known{}).TombstonedAt.Equal(timex.Inf()))

	before := time.Now().Add(3 * time.Hour)
	got := ddisc.NewDiscoveredFromImport(&ddiscapi.Import{Uri: "https://tracker.example/download/1.torrent"}).TombstonedAt
	after := time.Now().Add(3 * time.Hour)
	require.False(t, got.Before(before))
	require.False(t, got.After(after))
}
