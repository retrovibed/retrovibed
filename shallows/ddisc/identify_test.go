package ddisc_test

import (
	"testing"
	"time"

	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/bencode"
	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/james-lawrence/torrent/storage"
	"github.com/james-lawrence/torrent/torrenttest"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/torrenttestx"
	"github.com/stretchr/testify/require"
)

func TestIdentifyOne(t *testing.T) {
	t.Run("identifies media from a real torrent download without touching the database", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		target := torrenttestx.QuickDHT(t, dht.OptionBootstrapNodesNone)
		seeder := torrenttestx.QuickClientWithDHT(t, target)
		defer seeder.Close()

		s := torrenttestx.QuickDHT(t, dht.OptionBootstrapNodesNone, dht.OptionBootstrapFixedAddrs(dht.NewAddr(target.DynamicAddrPort())))
		consumer := torrenttestx.QuickClientWithDHT(t, s)
		defer consumer.Close()

		seedDir := t.TempDir()
		info, _, err := torrenttest.Random(seedDir, 32*bytesx.KiB)
		require.NoError(t, err)

		encoded, err := bencode.Marshal(info)
		require.NoError(t, err)
		hash := metainfo.NewHashFromBytes(encoded)
		id := int160.FromBytes(hash.Bytes())

		seedermd, err := torrent.NewFromInfo(info, torrent.OptionStorage(storage.NewFile(seedDir)))
		require.NoError(t, err)
		_, _, err = seeder.Start(seedermd, torrent.TuneAnnounceUntilComplete, torrent.TuneNewConns)
		require.NoError(t, err)

		consumerStorage := storage.NewFile(t.TempDir())
		go func() {
			for ctx.Err() == nil {
				if md, err := torrent.NewFromInfo(info, torrent.OptionStorage(consumerStorage)); err == nil {
					_, _, _ = consumer.Start(md, torrent.TuneClientPeer(seeder), torrent.TuneAnnounceUntilComplete, torrent.TuneNewConns)
				}
				time.Sleep(5 * time.Millisecond)
			}
		}()

		disc := ddisc.NewDiscovered(&id, ddisc.DiscoveredOptionFromTorrentInfo(info))

		result, err := ddisc.IdentifyOne(ctx, s, consumer, consumerStorage, 5*time.Second, 10*time.Second, disc)
		require.NoError(t, err)
		require.NotEmpty(t, result.Mimetype)
		require.Equal(t, disc.ID, result.ID, "IdentifyOne should not mutate the identity of the record")
	})
}
