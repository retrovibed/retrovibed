package daemons_test

import (
	"net"
	"path/filepath"
	"testing"
	"time"

	"github.com/james-lawrence/torrent/bep0051"
	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/dht/krpc"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/retrovibe/daemons"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

// exercises the exact wiring used in production: a dht server with
// bep0051.Query bound to bep0051.NewEndpoint(daemons.NewSampler(...)),
// queried over the wire by a separate dht server via
// bep0051.LatestSampleForNodeInfo.
func TestBEP51EndpointServesSamplerOverDHT(t *testing.T) {
	t.Run("returns the sampler's infohashes to a real sample_infohashes query", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		cachepath := filepath.Join(t.TempDir(), "sample.cache")

		one := tracking.NewMetadata(langx.Autoptr(int160.Random()))
		one.Private = false
		one.Seeding = true
		require.NoError(t, tracking.MetadataInsertWithDefaults(t.Context(), q, one).Scan(&one))

		two := tracking.NewMetadata(langx.Autoptr(int160.Random()))
		two.Private = false
		two.Seeding = true
		require.NoError(t, tracking.MetadataInsertWithDefaults(t.Context(), q, two).Scan(&two))

		// noise that must never appear in the sample.
		priv := tracking.NewMetadata(langx.Autoptr(int160.Random()))
		priv.Private = true
		priv.Seeding = true
		require.NoError(t, tracking.MetadataInsertWithDefaults(t.Context(), q, priv).Scan(&priv))

		target, err := dht.NewServer(32, dht.OptionMuxer(
			dht.DefaultMuxer().Method(bep0051.Query, bep0051.NewEndpoint(daemons.NewSampler(q, time.Hour, cachepath))),
		))
		require.NoError(t, err)
		t.Cleanup(target.Close)
		targetpc, err := net.ListenPacket("udp4", "127.0.0.1:0")
		require.NoError(t, err)
		require.NoError(t, target.Serve(t.Context(), targetpc))

		client, err := dht.NewServer(32)
		require.NoError(t, err)
		t.Cleanup(client.Close)
		clientpc, err := net.ListenPacket("udp4", "127.0.0.1:0")
		require.NoError(t, err)
		require.NoError(t, client.Serve(t.Context(), clientpc))

		targetAddr := target.DynamicAddrPort()
		info := krpc.NewInfo(target.ID(targetAddr).AsByteArray(), krpc.NewNodeAddrFromAddrPort(targetAddr))

		sample, err := bep0051.LatestSampleForNodeInfo(ctx, client, info)
		require.NoError(t, err)
		require.EqualValues(t, uint(time.Hour/time.Second), sample.Interval)
		require.EqualValues(t, 2, sample.Available)
		require.Len(t, sample.Sample, 2*20)

		expected := map[string]struct{}{
			string(one.Infohash): {},
			string(two.Infohash): {},
		}
		for i := 0; i < len(sample.Sample)/20; i++ {
			_, ok := expected[string(sample.Sample[i*20:(i+1)*20])]
			require.True(t, ok, "unexpected infohash in sample")
		}
	})

	t.Run("returns zero hashes when nothing qualifies", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		cachepath := filepath.Join(t.TempDir(), "sample.cache")

		target, err := dht.NewServer(32, dht.OptionMuxer(
			dht.DefaultMuxer().Method(bep0051.Query, bep0051.NewEndpoint(daemons.NewSampler(q, time.Hour, cachepath))),
		))
		require.NoError(t, err)
		t.Cleanup(target.Close)
		targetpc, err := net.ListenPacket("udp4", "127.0.0.1:0")
		require.NoError(t, err)
		require.NoError(t, target.Serve(t.Context(), targetpc))

		client, err := dht.NewServer(32)
		require.NoError(t, err)
		t.Cleanup(client.Close)
		clientpc, err := net.ListenPacket("udp4", "127.0.0.1:0")
		require.NoError(t, err)
		require.NoError(t, client.Serve(t.Context(), clientpc))

		targetAddr := target.DynamicAddrPort()
		info := krpc.NewInfo(target.ID(targetAddr).AsByteArray(), krpc.NewNodeAddrFromAddrPort(targetAddr))

		sample, err := bep0051.LatestSampleForNodeInfo(ctx, client, info)
		require.NoError(t, err)
		require.Zero(t, sample.Available)
		require.Empty(t, sample.Sample)
	})
}
