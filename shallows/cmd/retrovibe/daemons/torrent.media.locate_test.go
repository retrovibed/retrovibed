package daemons_test

import (
	"testing"

	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/autobind"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/james-lawrence/torrent/storage"
	"github.com/james-lawrence/torrent/torrenttest"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/retrovibe/daemons"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/torrenttestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestMediaLocate(t *testing.T) {
	t.Run("should be able to locate and create torrent metadata from distributed discovery", func(t *testing.T) {
		var (
			k library.Known
			l ddisc.Locate
			d ddisc.Discovered
		)
		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&k, library.KnownOptionTestDefaults))
		require.NoError(t, library.KnownInsertWithDefaults(t.Context(), q, k).Scan(&k))
		require.NoError(t, ddisc.LocateInsertWithDefaults(t.Context(), q, ddisc.NewLocate(k.Title, mimex.Binary, ddisc.LocateOptionKnownMedia(k.UID), ddisc.LocateOptionAutoDownload(true))).Scan(&l))

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
			ddisc.DiscoveredOptionKnownMedia(k.UID),
			ddisc.DiscoveredOptionMimetype(mimex.Binary),
			ddisc.DiscoveredOptionFromTorrentInfo(info),
		)

		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(t.Context(), q, d).Scan(&d))

		require.Equal(t, 0, testx.Must(sqlx.Count(t.Context(), q, "SELECT COUNT(*) FROM torrents_metadata"))(t))

		require.NoError(t, daemons.LocateMedia(t.Context(), q, tclient, &daemons.DiscoverySettings{LocateP2P: true}, nil, nil, nil, ddisc.DefaultPolicy()))

		require.Equal(t, 1, testx.Must(sqlx.Count(t.Context(), q, "SELECT COUNT(*) FROM torrents_metadata WHERE initiated_at <= NOW()"))(t))
		require.Equal(t, 1, testx.Must(sqlx.Count(t.Context(), q, "SELECT COUNT(*) FROM ddisc_locate WHERE tombstoned_at < 'infinity'"))(t))
	})

	t.Run("should record a recommendation instead of downloading when autodownload is disabled", func(t *testing.T) {
		var (
			k library.Known
			l ddisc.Locate
			d ddisc.Discovered
		)
		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&k, library.KnownOptionTestDefaults))
		require.NoError(t, library.KnownInsertWithDefaults(t.Context(), q, k).Scan(&k))
		require.NoError(t, ddisc.LocateInsertWithDefaults(t.Context(), q, ddisc.NewLocate(k.Title, mimex.Binary, ddisc.LocateOptionKnownMedia(k.UID))).Scan(&l))

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
			ddisc.DiscoveredOptionKnownMedia(k.UID),
			ddisc.DiscoveredOptionMimetype(mimex.Binary),
			ddisc.DiscoveredOptionFromTorrentInfo(info),
		)

		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(t.Context(), q, d).Scan(&d))

		require.NoError(t, daemons.LocateMedia(t.Context(), q, tclient, &daemons.DiscoverySettings{LocateP2P: true}, nil, nil, nil, ddisc.DefaultPolicy()))

		require.Equal(t, 0, testx.Must(sqlx.Count(t.Context(), q, "SELECT COUNT(*) FROM torrents_metadata WHERE initiated_at <= NOW()"))(t))
		require.Equal(t, 1, testx.Must(sqlx.Count(t.Context(), q, "SELECT COUNT(*) FROM library_recommendations"))(t))
		require.Equal(t, 1, testx.Must(sqlx.Count(t.Context(), q, "SELECT COUNT(*) FROM ddisc_locate WHERE tombstoned_at < 'infinity'"))(t))
	})

	t.Run("should re-locate when the same query is searched again after its recommendation was consumed", func(t *testing.T) {
		var (
			k library.Known
			l ddisc.Locate
			d ddisc.Discovered
		)
		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&k, library.KnownOptionTestDefaults))
		require.NoError(t, library.KnownInsertWithDefaults(t.Context(), q, k).Scan(&k))
		require.NoError(t, ddisc.LocateInsertWithDefaults(t.Context(), q, ddisc.NewLocate(k.Title, mimex.Binary, ddisc.LocateOptionKnownMedia(k.UID))).Scan(&l))

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
			ddisc.DiscoveredOptionKnownMedia(k.UID),
			ddisc.DiscoveredOptionMimetype(mimex.Binary),
			ddisc.DiscoveredOptionFromTorrentInfo(info),
		)

		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(t.Context(), q, d).Scan(&d))

		require.NoError(t, daemons.LocateMedia(t.Context(), q, tclient, &daemons.DiscoverySettings{LocateP2P: true}, nil, nil, nil, ddisc.DefaultPolicy()))

		var located ddisc.Locate
		require.NoError(t, ddisc.LocateFindByID(t.Context(), q, l.ID).Scan(&located))

		var rec library.Recommendation
		require.NoError(t, library.RecommendationFindByContentID(t.Context(), q, located.LocatedTorrentID).Scan(&rec))

		// simulate the user downloading the content: KnownMediaLocator._onTap
		// hard-deletes the recommendation once it's queued for download.
		require.NoError(t, library.RecommendationDeleteByID(t.Context(), q, rec.ID).Scan(&rec))
		require.Equal(t, 0, testx.Must(sqlx.Count(t.Context(), q, "SELECT COUNT(*) FROM library_recommendations"))(t))

		// the user searches the exact same query again - same deterministic
		// locate id, but the request should reopen for reprocessing rather
		// than silently handing back the now-dangling located_torrent_id.
		require.NoError(t, ddisc.LocateInsertWithDefaults(t.Context(), q, ddisc.NewLocate(k.Title, mimex.Binary, ddisc.LocateOptionKnownMedia(k.UID))).Scan(&l))

		require.NoError(t, daemons.LocateMedia(t.Context(), q, tclient, &daemons.DiscoverySettings{LocateP2P: true}, nil, nil, nil, ddisc.DefaultPolicy()))

		require.NoError(t, ddisc.LocateFindByID(t.Context(), q, l.ID).Scan(&located))
		require.NotEqual(t, "FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF", located.LocatedTorrentID)

		// GET /r/content/<locatedTorrentId> must resolve again - no permanent 404.
		require.Equal(t, 1, testx.Must(sqlx.Count(t.Context(), q, "SELECT COUNT(*) FROM library_recommendations"))(t))
		require.NoError(t, library.RecommendationFindByContentID(t.Context(), q, located.LocatedTorrentID).Scan(&rec))
	})

	t.Run("should not mark a locate as located when recording the recommendation fails", func(t *testing.T) {
		var (
			k library.Known
			l ddisc.Locate
			d ddisc.Discovered
		)
		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&k, library.KnownOptionTestDefaults))
		require.NoError(t, library.KnownInsertWithDefaults(t.Context(), q, k).Scan(&k))
		require.NoError(t, ddisc.LocateInsertWithDefaults(t.Context(), q, ddisc.NewLocate(k.Title, mimex.Binary, ddisc.LocateOptionKnownMedia(k.UID))).Scan(&l))

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
			ddisc.DiscoveredOptionKnownMedia(k.UID),
			ddisc.DiscoveredOptionMimetype(mimex.Binary),
			ddisc.DiscoveredOptionFromTorrentInfo(info),
		)
		// corrupt the id so the library_recommendations insert genuinely fails
		// (content_id is a UUID column; this fails DuckDB's implicit cast).
		d.ID = "not-a-uuid"

		require.Error(t, daemons.DiscoveredDownload(t.Context(), q, tclient, l, d))

		require.Equal(t, 0, testx.Must(sqlx.Count(t.Context(), q, "SELECT COUNT(*) FROM library_recommendations"))(t))
		require.Equal(t, 1, testx.Must(sqlx.Count(t.Context(), q, "SELECT COUNT(*) FROM ddisc_locate WHERE located_torrent_id = 'FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF'"))(t))
	})

	t.Run("should not query ddisc when p2p locate is disabled", func(t *testing.T) {
		var (
			k library.Known
			l ddisc.Locate
			d ddisc.Discovered
		)
		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&k, library.KnownOptionTestDefaults))
		require.NoError(t, library.KnownInsertWithDefaults(t.Context(), q, k).Scan(&k))
		require.NoError(t, ddisc.LocateInsertWithDefaults(t.Context(), q, ddisc.NewLocate(k.Title, mimex.Binary, ddisc.LocateOptionKnownMedia(k.UID))).Scan(&l))

		seedir := t.TempDir()
		mcache := torrent.NewMetadataCache(seedir)
		info, _, err := torrenttest.Random(seedir, 128*bytesx.KiB, metainfo.OptionDisplayName(k.Title))
		require.NoError(t, err)

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

		require.NoError(t, daemons.LocateMedia(t.Context(), q, tclient, &daemons.DiscoverySettings{LocateP2P: false}, nil, nil, nil, ddisc.DefaultPolicy()))

		require.Equal(t, 0, testx.Must(sqlx.Count(t.Context(), q, "SELECT COUNT(*) FROM torrents_metadata"))(t))
		require.Equal(t, 0, testx.Must(sqlx.Count(t.Context(), q, "SELECT COUNT(*) FROM ddisc_locate WHERE tombstoned_at < 'infinity'"))(t))
	})
}
