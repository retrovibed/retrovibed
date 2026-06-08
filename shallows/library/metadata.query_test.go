package library_test

import (
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestMetadataQueryHasKnownMedia(t *testing.T) {
	realID := uuid.Must(uuid.NewV4()).String()

	t.Run("true returns only items with a resolved known_media_id", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var withKnown, withNil, withMax library.Metadata
		require.NoError(t, testx.Fake(&withKnown, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID, library.MetadataOptionKnownMediaID(realID)))
		require.NoError(t, testx.Fake(&withNil, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID)) // nil UUID
		require.NoError(t, testx.Fake(&withMax, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID, library.MetadataOptionKnownMediaID(uuid.Max.String())))

		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, withKnown).Scan(&withKnown))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, withNil).Scan(&withNil))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, withMax).Scan(&withMax))

		q := library.MetadataSearchBuilder().Where(library.MetadataQueryHasKnownMedia(true))
		scanner := sqlx.Scan(library.MetadataSearch(ctx, db, q))
		var results []library.Metadata
		for md := range scanner.Iter() {
			results = append(results, md)
		}
		require.NoError(t, scanner.Err())
		require.Len(t, results, 1)
		require.Equal(t, withKnown.ID, results[0].ID)
	})

	t.Run("false returns only items without a resolved known_media_id", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var withKnown, withNil, withMax library.Metadata
		require.NoError(t, testx.Fake(&withKnown, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID, library.MetadataOptionKnownMediaID(realID)))
		require.NoError(t, testx.Fake(&withNil, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
		require.NoError(t, testx.Fake(&withMax, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID, library.MetadataOptionKnownMediaID(uuid.Max.String())))

		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, withKnown).Scan(&withKnown))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, withNil).Scan(&withNil))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, withMax).Scan(&withMax))

		q := library.MetadataSearchBuilder().Where(library.MetadataQueryHasKnownMedia(false))
		scanner := sqlx.Scan(library.MetadataSearch(ctx, db, q))
		var results []library.Metadata
		for md := range scanner.Iter() {
			results = append(results, md)
		}
		require.NoError(t, scanner.Err())
		require.Len(t, results, 2)
	})
}

func TestMetadataQueryHasTorrent(t *testing.T) {
	torrentID := uuid.Must(uuid.NewV4()).String()

	t.Run("true returns only items with a torrent_id set", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var withTorrent, withoutTorrent library.Metadata
		require.NoError(t, testx.Fake(&withTorrent, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID, library.MetadataOptionTorrentID(torrentID)))
		require.NoError(t, testx.Fake(&withoutTorrent, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))

		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, withTorrent).Scan(&withTorrent))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, withoutTorrent).Scan(&withoutTorrent))

		q := library.MetadataSearchBuilder().Where(library.MetadataQueryHasTorrent(true))
		scanner := sqlx.Scan(library.MetadataSearch(ctx, db, q))
		var results []library.Metadata
		for md := range scanner.Iter() {
			results = append(results, md)
		}
		require.NoError(t, scanner.Err())
		require.Len(t, results, 1)
		require.Equal(t, withTorrent.ID, results[0].ID)
	})

	t.Run("false returns only items without a torrent_id", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var withTorrent, withoutTorrent library.Metadata
		require.NoError(t, testx.Fake(&withTorrent, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID, library.MetadataOptionTorrentID(torrentID)))
		require.NoError(t, testx.Fake(&withoutTorrent, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))

		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, withTorrent).Scan(&withTorrent))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, withoutTorrent).Scan(&withoutTorrent))

		q := library.MetadataSearchBuilder().Where(library.MetadataQueryHasTorrent(false))
		scanner := sqlx.Scan(library.MetadataSearch(ctx, db, q))
		var results []library.Metadata
		for md := range scanner.Iter() {
			results = append(results, md)
		}
		require.NoError(t, scanner.Err())
		require.Len(t, results, 1)
		require.Equal(t, withoutTorrent.ID, results[0].ID)
	})
}
