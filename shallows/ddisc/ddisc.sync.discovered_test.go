package ddisc_test

import (
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/torrenttest"
	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/cryptox"
	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/stretchr/testify/require"
)

func TestSyncDiscovered(t *testing.T) {
	t.Run("return any media that matches the requesting user since the sync id", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		tmpdir := t.TempDir()

		q := sqltestx.Metadatabase(t)

		partitions := ddisc.Partitions(128, cryptox.NewChaCha8(t.Name()))
		p0 := uuid.Must(uuid.NewV4())
		p1 := uuid.Must(uuid.NewV4())

		n := partitions.Max(p0.Bytes())

		block0 := ddisc.FilterRatio(cryptox.NewChaCha8(n[:]), 10)

		n1 := partitions.Max(p1.Bytes())

		block1 := ddisc.FilterRatio(cryptox.NewChaCha8(n1[:]), 10)

		for range 64 {
			id := int160.Random()
			info, _, err := torrenttest.Random(tmpdir, 128*bytesx.KiB)
			require.NoError(t, err)

			d := ddisc.NewDiscovered(
				&id,
				ddisc.DiscoveredOptionIndex(!block0.Filter(id.Bytes())),
				ddisc.DiscoveredOptionMimetype(mimex.Binary),
				ddisc.DiscoveredOptionFromTorrentInfo(info),
			)
			require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))
			require.Contains(t, []string{uuid.Nil.String(), uuid.Max.String()}, d.KnownMediaID)
		}

		require.InDelta(t, 7, sqltestx.Count(t, q, "SELECT COUNT(*) FROM ddisc_media WHERE known_media_id = 'FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF'"), 5)
		require.InDelta(t, 57, sqltestx.Count(t, q, "SELECT COUNT(*) FROM ddisc_media WHERE known_media_id = '00000000-0000-0000-0000-000000000000'"), 5)

		assertbatch := func(last ddisc.Discovered, syncer iterx.Seq[ddisc.Discovered]) ddisc.Discovered {
			var (
				count = 0
			)
			for d := range syncer.Each(t.Context()) {
				require.Contains(t, []string{uuid.Nil.String(), uuid.Max.String()}, d.KnownMediaID)
				require.GreaterOrEqual(t, d.SyncUID, last.SyncUID)
				last = d
				count++
			}
			require.NoError(t, syncer.Err())
			require.InDelta(t, 7, count, 6, "peer should only be receive ~10 percent of the data")
			return last
		}

		var (
			last = ddisc.Discovered{
				SyncUID: "00000000-0000-0000-0000-000000000000",
			}
		)
		last = assertbatch(last, ddisc.SyncDiscovered(q, block1, uuid.Nil.String()))
		require.NotEqual(t, uuid.Nil.String(), last.SyncUID)

		// add newly discovered data to the pool.
		for range 64 {
			id := int160.Random()
			info, _, err := torrenttest.Random(tmpdir, 128*bytesx.KiB)
			require.NoError(t, err)

			d := ddisc.NewDiscovered(
				&id,
				ddisc.DiscoveredOptionIndex(!block0.Filter(id.Bytes())),
				ddisc.DiscoveredOptionMimetype(mimex.Binary),
				ddisc.DiscoveredOptionFromTorrentInfo(info),
			)
			require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))
		}

		next := assertbatch(last, ddisc.SyncDiscovered(q, block1, last.SyncUID))
		require.NotEqual(t, last.SyncUID, next.SyncUID)
	})
}
