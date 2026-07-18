package ddisc_test

import (
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/torrenttest"
	"github.com/retrovibed/retrovibed/retroapi/ddiscapi"
	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/retroapi/uuidx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/cryptox"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindMedia(t *testing.T) {
	t.Run("return any media that matches the requesting peer since the sync id", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		tmpdir := t.TempDir()

		q := sqltestx.Metadatabase(t)

		partitions := ddisc.Partitions(16, cryptox.NewChaCha8(t.Name()))

		p0 := uuid.FromStringOrNil("be526731-7bfc-40f5-8984-73f8ae033ead")
		n0 := partitions.Max(p0.Bytes())

		genrecord := func(idx int) {
			id := int160.Random()
			info, _, err := torrenttest.Random(tmpdir, 128*bytesx.KiB)
			require.NoError(t, err)

			d := ddisc.NewDiscovered(
				&id,
				ddisc.DiscoveredOptionKnownMedia(ddiscapi.ImportedMediaUUID("", uuid.FromStringOrNil(uuidx.WithSuffix(idx%32))).String()),
				ddisc.DiscoveredOptionMimetype(mimex.Binary),
				ddisc.DiscoveredOptionFromTorrentInfo(info),
				ddisc.DiscoveredOptionPartitionAuto(partitions),
				ddisc.DiscoveredOptionAutoMagnet,
			)

			require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))
			require.NotEqual(t, uuid.Nil.String(), d.KnownMediaID)
			require.NotEqual(t, uuid.Nil.String(), d.Partition)
		}
		for idx := range 128 {
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
				ddisc.DiscoveredOptionAutoMagnet,
			)
			require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))
			require.Equal(t, uuid.Nil.String(), d.KnownMediaID)
			require.Equal(t, uuid.Nil.String(), d.Partition)
		}

		require.EqualValues(t, 32, sqltestx.Count(t, q, "SELECT COUNT(*) FROM ddisc_media WHERE partition = '00000000-0000-0000-0000-000000000000'"))
		require.EqualValues(t, 128, sqltestx.Count(t, q, "SELECT COUNT(*) FROM ddisc_media WHERE partition != '00000000-0000-0000-0000-000000000000'"))
		require.EqualValues(t, 4, sqltestx.Count(t, q, "SELECT COUNT(*) FROM ddisc_media WHERE partition = ?", n0.String()))

		assertbatch := func(idx int, n int, syncer iterx.Seq[ddisc.Discovered]) {
			var (
				count = 0
			)
			for d := range syncer.Each(t.Context()) {
				require.NotEqual(t, uuid.Nil.String(), d.KnownMediaID)
				count++
			}
			require.NoError(t, syncer.Err())
			assert.Equal(t, n, count, "partition should have received n records of data %d", idx)
		}

		for idx, n := range []int{4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4} {
			assertbatch(idx, n, ddisc.FindMedia(q, ddiscapi.ImportedMediaUUID("", uuid.FromStringOrNil(uuidx.WithSuffix(idx))).String()))
		}

	})

	t.Run("never returns private media", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		tmpdir := t.TempDir()

		q := sqltestx.Metadatabase(t)

		kid := ddiscapi.ImportedMediaUUID("", uuid.FromStringOrNil(uuidx.WithSuffix(0))).String()

		for range 16 {
			id := int160.Random()
			info, _, err := torrenttest.Random(tmpdir, 128*bytesx.KiB)
			require.NoError(t, err)
			priv := true
			info.Private = &priv

			d := ddisc.NewDiscovered(
				&id,
				ddisc.DiscoveredOptionKnownMedia(kid),
				ddisc.DiscoveredOptionMimetype(mimex.Binary),
				ddisc.DiscoveredOptionFromTorrentInfo(info),
				ddisc.DiscoveredOptionAutoMagnet,
			)
			require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))
			require.True(t, d.Private)
		}

		count := 0
		seq := ddisc.FindMedia(q, kid)
		for d := range seq.Each(t.Context()) {
			require.False(t, d.Private, "private media must never be returned from a peer search")
			count++
		}
		require.NoError(t, seq.Err())
		require.Zero(t, count, "expected no private media to be returned")
	})
}
