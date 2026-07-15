package ddisctorrent_test

import (
	"context"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/ddisc/ddisctorrent"
	"github.com/retrovibed/retrovibed/shallows/internal/cryptox"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/torrenttestx"
	"github.com/stretchr/testify/require"
)

func TestDiscoveredProtocol(t *testing.T) {
	t.Run("should receive results of the discovered protocol", func(t *testing.T) {
		const (
			records = 12
		)
		pdb := sqltestx.Metadatabase(t)
		defer pdb.Close()

		partitions := ddisc.Partitions(16, cryptox.NewChaCha8(t.Name()))

		ptm := dht.DefaultMuxer().
			Method(ddisctorrent.MethodDisc, ddisctorrent.NewDiscovered(pdb, partitions))
		pdht, err := dht.NewServer(
			32,
			dht.OptionMuxer(ptm),
		)
		require.NoError(t, err)
		tpeer := torrenttestx.QuickClientWithDHT(
			t,
			pdht,
			torrent.ClientConfigExtension(ddisctorrent.ExtensionName),
		)
		defer tpeer.Close()

		cdb := sqltestx.Metadatabase(t)
		defer cdb.Close()

		ctm := dht.DefaultMuxer().
			Method(ddisctorrent.MethodMedia, ddisctorrent.NewMediaRecorder(cdb))
		cdht, err := dht.NewServer(
			32,
			dht.OptionMuxer(ctm),
		)
		require.NoError(t, err)
		tclient := torrenttestx.QuickClientWithDHT(
			t,
			cdht,
			torrent.ClientConfigExtension(ddisctorrent.ExtensionName),
		)
		defer tclient.Close()

		genrecord := func(ctx context.Context, q sqlx.Queryer) {
			id := int160.Random()
			d := ddisc.NewDiscovered(
				&id,
				ddisc.DiscoveredOptionKnownMedia(uuid.Nil.String()),
				ddisc.DiscoveredOptionMimetype(mimex.Binary),
			)

			require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))
			require.Equal(t, uuid.Nil.String(), d.KnownMediaID)
		}
		for range records {
			genrecord(t.Context(), pdb)
		}

		req, err := ddisctorrent.NewDiscoveredRequest(cdht.ID(cdht.DynamicAddrPort()).AsByteArray(), 100, uuid.Nil.String())
		require.NoError(t, err)

		ret := cdht.Query(t.Context(), dht.NewAddr(pdht.DynamicAddrPort()), req)
		require.NoError(t, ret.Err)

		require.Eventually(t, func() bool {
			return sqltestx.Count(t, cdb, "SELECT COUNT(*) FROM ddisc_media") == records
		}, time.Second, 100*time.Millisecond)
	})
}
