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
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/torrenttestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestSearchProtocol(t *testing.T) {
	t.Run("should receive results of a search", func(t *testing.T) {
		const (
			records = 12
		)
		pdb := sqltestx.Metadatabase(t)
		defer pdb.Close()

		ptm := dht.DefaultMuxer().
			Method(ddisctorrent.MethodSearch, ddisctorrent.NewSearch(pdb))
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

		knownmedia := library.KnownImportedUUID(t.Name(), uuid.Must(uuid.NewV7()))
		genrecord := func(ctx context.Context, q sqlx.Queryer, idx int) {
			id := int160.Random()
			d := ddisc.NewDiscovered(
				&id,
				ddisc.DiscoveredOptionKnownMedia(knownmedia.String()),
				ddisc.DiscoveredOptionMimetype(mimex.Binary),
			)

			require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))
			require.NotEqual(t, uuid.Nil.String(), d.KnownMediaID)
		}
		for idx := range records {
			genrecord(t.Context(), pdb, idx)
		}

		req, err := ddisctorrent.NewSearchRequest(cdht.ID(cdht.DynamicAddrPort()), knownmedia.String())
		require.NoError(t, err)

		ret := cdht.Query(t.Context(), dht.NewAddr(pdht.DynamicAddrPort()), req)
		require.NoError(t, ret.Err)

		require.Eventually(t, func() bool {
			return sqltestx.Count(t, cdb, "SELECT COUNT(*) FROM ddisc_media") == records
		}, time.Second, 100*time.Millisecond)
	})
}
