package ddisc_test

import (
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/torrenttest"
	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/cryptox"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/uuidx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestSyncPartition(t *testing.T) {
	t.Run("return any media that matches the requesting peer since the sync id", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		tmpdir := t.TempDir()

		q := sqltestx.Metadatabase(t)

		partitions := ddisc.Partitions(16, cryptox.NewChaCha8(t.Name()))
		p0 := uuid.FromStringOrNil("d7e5e5df-b2ad-416a-97c6-c634465cec75")
		p1 := uuid.FromStringOrNil("0c67c93f-a940-4d55-aab0-aaacaea2b050")

		n0 := partitions.Max(p0.Bytes())
		n1 := partitions.Max(p1.Bytes())

		genrecord := func(idx int) {
			id := int160.Random()
			info, _, err := torrenttest.Random(tmpdir, 128*bytesx.KiB)
			require.NoError(t, err)

			d := ddisc.NewDiscovered(
				&id,
				ddisc.DiscoveredOptionKnownMedia(library.KnownImportedUUID("", uuid.FromStringOrNil(uuidx.WithSuffix(idx))).String()),
				ddisc.DiscoveredOptionMimetype(mimex.Binary),
				ddisc.DiscoveredOptionFromTorrentInfo(info),
				ddisc.DiscoveredOptionPartitionAuto(partitions),
			)

			require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))
			require.NotEqual(t, uuid.Nil.String(), d.KnownMediaID)
			require.NotEqual(t, uuid.Nil.String(), d.Partition)
		}
		for idx := range 64 {
			genrecord(idx)
		}

		for range 32 {
			// generate unknown media
			id := int160.Random()
			info, _, err := torrenttest.Random(tmpdir, 128*bytesx.KiB)
			require.NoError(t, err)

			d := ddisc.NewDiscovered(
				&id,
				ddisc.DiscoveredOptionMimetype(mimex.Binary),
				ddisc.DiscoveredOptionFromTorrentInfo(info),
				ddisc.DiscoveredOptionPartitionAuto(partitions),
			)
			require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))
			require.Equal(t, uuid.Nil.String(), d.KnownMediaID)
			require.Equal(t, uuid.Nil.String(), d.Partition)
		}

		require.InDelta(t, 64, sqltestx.Count(t, q, "SELECT COUNT(*) FROM ddisc_media WHERE partition != '00000000-0000-0000-0000-000000000000'"), 3)
		require.InDelta(t, 32, sqltestx.Count(t, q, "SELECT COUNT(*) FROM ddisc_media WHERE partition = '00000000-0000-0000-0000-000000000000'"), 3)

		assertbatch := func(last ddisc.Discovered, n int, syncer iterx.Seq[ddisc.Discovered]) ddisc.Discovered {
			var (
				count = 0
			)
			for d := range syncer.Each(t.Context()) {
				require.NotEqual(t, uuid.Nil.String(), d.KnownMediaID)
				require.GreaterOrEqual(t, d.SyncUID, last.SyncUID)
				last = d
				count++
			}
			require.NoError(t, syncer.Err())
			require.Equal(t, n, count, "partition should have received n records of data")
			return last
		}

		var (
			last = ddisc.Discovered{
				SyncUID: "00000000-0000-0000-0000-000000000000",
			}
		)

		prev := last
		last = assertbatch(prev, 9, ddisc.SyncPartition(q, n0.String(), prev.SyncUID))
		require.NotEqual(t, uuid.Nil.String(), last.SyncUID)
		last = assertbatch(prev, 5, ddisc.SyncPartition(q, n1.String(), prev.SyncUID))
		require.NotEqual(t, uuid.Nil.String(), last.SyncUID)

		for idx := range 128 {
			genrecord(idx + 64)
		}

		prev = last
		last = assertbatch(prev, 7, ddisc.SyncPartition(q, n0.String(), prev.SyncUID))
		require.NotEqual(t, uuid.Nil.String(), last.SyncUID)
		assertbatch(ddisc.Discovered{SyncUID: uuid.Nil.String()}, 16, ddisc.SyncPartition(q, n0.String(), uuid.Nil.String()))
	})
}
