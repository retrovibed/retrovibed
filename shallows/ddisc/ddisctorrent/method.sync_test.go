package ddisctorrent_test

import (
	"context"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/ddisc"
	"github.com/retrovibed/retrovibed/ddisc/ddisctorrent"
	"github.com/retrovibed/retrovibed/internal/mimex"
	"github.com/retrovibed/retrovibed/internal/sqltestx"
	"github.com/retrovibed/retrovibed/internal/sqlx"
	"github.com/retrovibed/retrovibed/internal/torrenttestx"
	"github.com/retrovibed/retrovibed/internal/uuidx"
	"github.com/retrovibed/retrovibed/library"
	"github.com/stretchr/testify/require"
)

func TestSyncProtocol(t *testing.T) {
	t.Run("should receive results of a sync", func(t *testing.T) {
		const (
			records = 12
		)
		pdb := sqltestx.Metadatabase(t)
		defer pdb.Close()

		ptm := dht.DefaultMuxer().
			Method(ddisctorrent.MethodSync, ddisctorrent.NewSync(pdb))
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
				ddisc.DiscoveredOptionPartition(uuidx.WithSuffix(idx%2)),
			)

			require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))
			require.NotEqual(t, uuid.Nil.String(), d.KnownMediaID)
		}
		for idx := range records {
			genrecord(t.Context(), pdb, idx)
		}

		req, err := ddisctorrent.NewSyncRequest(cdht.ID(), uuidx.WithSuffix(0), uuid.Nil.String())
		require.NoError(t, err)

		ret := cdht.Query(t.Context(), dht.NewAddr(pdht.AddrPort()), req)
		require.NoError(t, ret.Err)

		require.Eventually(t, func() bool {
			return sqltestx.Count(t, cdb, "SELECT COUNT(*) FROM ddisc_media") == (records / 2)
		}, time.Second, 100*time.Millisecond)
	})
}
