package daemons_test

import (
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/bencode"
	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/james-lawrence/torrent/storage"
	"github.com/james-lawrence/torrent/torrenttest"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/retrovibe/daemons"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/torrenttestx"
	"github.com/stretchr/testify/require"
)

func TestDiscoverMedia(t *testing.T) {
	t.Run("successfully identifies and indexes media discovered via a real torrent download", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		target := torrenttestx.QuickDHT(t, dht.OptionBootstrapNodesNone)
		seeder := torrenttestx.QuickClientWithDHT(t, target)
		defer seeder.Close()

		// s only needs a starting node so DiscoverMedia's own peer-lookup gate
		// (torrentx.Peers) doesn't fail outright with "no initial nodes" - real
		// BEP5 get_peers/announce_peer discovery for a torrent that's already
		// fully seeded proved unreliable to drive deterministically within a
		// short test window (the seeder's announce can take the better part of
		// a minute to land). The actual transfer below is wired directly via
		// TuneClientPeer instead.
		s := torrenttestx.QuickDHT(
			t,
			dht.OptionBootstrapNodesNone,
			dht.OptionBootstrapFixedAddrs(dht.NewAddr(target.DynamicAddrPort())),
		)
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

		// DiscoverMedia creates and drops its own *torrent for this infohash
		// internally (once to fetch info, again to download the content), so
		// rather than wiring the peer connection once, keep re-arming a direct
		// connection to the seeder for as long as the test runs. cl.Start is a
		// cheap cache hit whenever the torrent is already registered, and only
		// does real work (reconnecting) right after DiscoverMedia drops it.
		consumerStorage := storage.NewFile(t.TempDir())
		go func() {
			for ctx.Err() == nil {
				if md, err := torrent.NewFromInfo(info, torrent.OptionStorage(consumerStorage)); err == nil {
					_, _, err = consumer.Start(md, torrent.TuneClientPeer(seeder), torrent.TuneAnnounceUntilComplete, torrent.TuneNewConns)
				}
				time.Sleep(200 * time.Millisecond)
			}
		}()

		disc := ddisc.NewDiscovered(&id, ddisc.DiscoveredOptionFromTorrentInfo(info))
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, disc).Scan(&disc))

		go func() {
			errorsx.Log(daemons.DiscoverMedia(
				ctx, q, s, consumer,
				daemons.DiscoverMediaOptionFrequency(10*time.Millisecond),
				daemons.DiscoverMediaOptionPeerTimeout(5*time.Second),
				daemons.DiscoverMediaOptionInfoTimeout(10*time.Second),
			))
		}()

		require.Eventually(t, func() bool {
			var updated ddisc.Discovered
			if err := ddisc.DiscoveredFindByID(ctx, q, disc.ID).Scan(&updated); err != nil {
				return false
			}
			return updated.KnownMediaID == uuid.Max.String() && updated.Mimetype != ""
		}, 30*time.Second, 100*time.Millisecond, "expected discovered media to be identified and indexed")
	})

	t.Run("pushes a row into cooldown when its content cannot be located within the configured timeouts", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		dhts := torrenttestx.QuickDHT(t)
		tclient := torrenttestx.QuickClientWithDHT(t, dhts)
		defer tclient.Close()

		id := int160.Random()
		disc := ddisc.NewDiscovered(&id)
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, disc).Scan(&disc))

		go func() {
			_ = daemons.DiscoverMedia(
				ctx, q, dhts, tclient,
				daemons.DiscoverMediaOptionFrequency(10*time.Millisecond),
				daemons.DiscoverMediaOptionPeerTimeout(500*time.Millisecond),
				daemons.DiscoverMediaOptionInfoTimeout(2*time.Second),
			)
		}()

		require.Eventually(t, func() bool {
			var updated ddisc.Discovered
			if err := ddisc.DiscoveredFindByID(ctx, q, disc.ID).Scan(&updated); err != nil {
				return false
			}
			return updated.Attempts > disc.Attempts && updated.KnownMediaID == uuid.Nil.String()
		}, 10*time.Second, 100*time.Millisecond, "expected unreachable media to be pushed into cooldown")
	})

	t.Run("leaves a row untouched when its next check is still in the future", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		dhts := torrenttestx.QuickDHT(t)
		tclient := torrenttestx.QuickClientWithDHT(t, dhts)
		defer tclient.Close()

		id := int160.Random()
		disc := ddisc.NewDiscovered(&id)
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, disc).Scan(&disc))
		_, err := q.ExecContext(ctx, `UPDATE ddisc_media SET next_check_at = NOW() + INTERVAL '1 hour' WHERE id = $1`, disc.ID)
		require.NoError(t, err)

		go func() {
			_ = daemons.DiscoverMedia(
				ctx, q, dhts, tclient,
				daemons.DiscoverMediaOptionFrequency(10*time.Millisecond),
				daemons.DiscoverMediaOptionPeerTimeout(500*time.Millisecond),
				daemons.DiscoverMediaOptionInfoTimeout(2*time.Second),
			)
		}()

		require.Never(t, func() bool {
			var updated ddisc.Discovered
			if err := ddisc.DiscoveredFindByID(ctx, q, disc.ID).Scan(&updated); err != nil {
				return false
			}
			return updated.Attempts != disc.Attempts || updated.KnownMediaID != disc.KnownMediaID
		}, 2*time.Second, 100*time.Millisecond, "row not yet due should not be touched by DiscoverMedia")
	})
}
