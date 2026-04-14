package daemons_test

import (
	"testing"

	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/autobind"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/james-lawrence/torrent/storage"
	"github.com/james-lawrence/torrent/torrenttest"
	"github.com/retrovibed/retrovibed/shallows/cmd/retrovibed/daemons"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/torrenttestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestMediaLocate(t *testing.T) {
	t.Run("should be able to locate and create torrent metadata from distributed discovery", func(t *testing.T) {
		var (
			k library.Known
			l library.Locate
			d ddisc.Discovered
		)
		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&k, library.KnownOptionTestDefaults))
		require.NoError(t, library.KnownInsertWithDefaults(t.Context(), q, k).Scan(&k))
		require.NoError(t, library.LocateInsertWithDefaults(t.Context(), q, library.Locate{KnownMediaID: k.UID}).Scan(&l))

		seedir := t.TempDir()
		mcache := torrent.NewMetadataCache(seedir)
		info, _, err := torrenttest.Random(seedir, 128*bytesx.KiB, metainfo.OptionDisplayName(k.Title))
		require.NoError(t, err)
		md, err := torrent.NewFromInfo(info, torrent.OptionStorage(storage.NewFile(seedir)))
		require.NoError(t, err)
		require.NoError(t, mcache.Write(md))

		tclient := torrenttestx.Client(
			t,
			autobind.NewLoopback(
				autobind.EnableDHT(torrenttestx.QuickDHT(t)),
			),
			mcache,
			storage.NewFile(seedir),
		)
		defer tclient.Close()

		id := int160.New(testx.Must(metainfo.Encode(info))(t))

		d = ddisc.NewDiscovered(
			&id,
			ddisc.DiscoveredOptionIndex(true),
			ddisc.DiscoveredOptionMimetype(mimex.Binary),
			ddisc.DiscoveredOptionFromTorrentInfo(info),
		)

		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(t.Context(), q, d).Scan(&d))

		require.Equal(t, 0, testx.Must(sqlx.Count(t.Context(), q, "SELECT COUNT(*) FROM torrents_metadata"))(t))

		require.NoError(t, daemons.LocateTorrentMedia(t.Context(), q, tclient))

		require.Equal(t, 1, testx.Must(sqlx.Count(t.Context(), q, "SELECT COUNT(*) FROM torrents_metadata WHERE initiated_at <= NOW()"))(t))
	})
}
