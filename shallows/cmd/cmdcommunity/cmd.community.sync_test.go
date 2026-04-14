package cmdcommunity

import (
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestCommunitySync(t *testing.T) {
	t.Run("syncs published content from deeppool", func(t *testing.T) {
		ctx, cancel := testx.Context(t)
		defer cancel()

		db := sqltestx.Metadatabase(t)
		defer db.Close()

		communityID := uuid.Must(uuid.NewV7()).String()
		knownMediaID := uuid.Must(uuid.NewV7()).String()
		publishedContentID := uuid.Must(uuid.NewV7()).String()
		libraryID := uuid.Must(uuid.NewV7()).String()

		err := community.SyncPublishedContentItem(ctx, db, &meta.PublishedContent{
			Id:           publishedContentID,
			CommunityId:  communityID,
			KnownMediaId: knownMediaID,
			MagnetUri:    "magnet:?xt=urn:btih:0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33&dn=testfile.txt",
			LibraryId:    libraryID,
		}, false)
		require.NoError(t, err)

		require.Equal(t, 1, errorsx.Zero(sqlx.Count(ctx, db, "SELECT COUNT(*) FROM torrents_metadata")))
		require.Equal(t, 1, errorsx.Zero(sqlx.Count(ctx, db, "SELECT COUNT(*) FROM published_content")))

		var tmeta tracking.Metadata
		require.NoError(t, tracking.MetadataFindByID(ctx, db, errorsx.Must(sqlx.String(ctx, db, "SELECT id::text FROM torrents_metadata"))).Scan(&tmeta))
		require.Equal(t, "testfile.txt", tmeta.Description)
		require.EqualValues(t, 0x0b, tmeta.Infohash[0])
		require.EqualValues(t, 0xee, tmeta.Infohash[1])
		require.EqualValues(t, 0xc7, tmeta.Infohash[2])

		require.Equal(t, 1, errorsx.Zero(sqlx.Count(ctx, db, "SELECT COUNT(*) FROM published_content")))
	})

	t.Run("syncs multiple published items", func(t *testing.T) {
		ctx, cancel := testx.Context(t)
		defer cancel()

		db := sqltestx.Metadatabase(t)
		defer db.Close()

		communityID := uuid.Must(uuid.NewV7()).String()

		items := []*meta.PublishedContent{
			{
				Id:           uuid.Must(uuid.NewV7()).String(),
				CommunityId:  communityID,
				KnownMediaId: uuid.Must(uuid.NewV7()).String(),
				MagnetUri:    "magnet:?xt=urn:btih:1111111111111111111111111111111111111111&dn=file1.txt",
				LibraryId:    uuid.Must(uuid.NewV7()).String(),
			},
			{
				Id:           uuid.Must(uuid.NewV7()).String(),
				CommunityId:  communityID,
				KnownMediaId: uuid.Must(uuid.NewV7()).String(),
				MagnetUri:    "magnet:?xt=urn:btih:2222222222222222222222222222222222222222&dn=file2.txt",
				LibraryId:    uuid.Must(uuid.NewV7()).String(),
			},
			{
				Id:           uuid.Must(uuid.NewV7()).String(),
				CommunityId:  communityID,
				KnownMediaId: uuid.Must(uuid.NewV7()).String(),
				MagnetUri:    "magnet:?xt=urn:btih:3333333333333333333333333333333333333333&dn=file3.txt",
				LibraryId:    uuid.Must(uuid.NewV7()).String(),
			},
		}

		for _, item := range items {
			require.NoError(t, community.SyncPublishedContentItem(ctx, db, item, false))
		}

		require.Equal(t, 3, errorsx.Zero(sqlx.Count(ctx, db, "SELECT COUNT(*) FROM torrents_metadata")))
		require.Equal(t, 3, errorsx.Zero(sqlx.Count(ctx, db, "SELECT COUNT(*) FROM published_content")))
	})

	t.Run("handles duplicate sync idempotently", func(t *testing.T) {
		ctx, cancel := testx.Context(t)
		defer cancel()

		db := sqltestx.Metadatabase(t)
		defer db.Close()

		communityID := uuid.Must(uuid.NewV7()).String()
		publishedContentID := uuid.Must(uuid.NewV7()).String()
		libraryID := uuid.Must(uuid.NewV7()).String()

		pc := &meta.PublishedContent{
			Id:           publishedContentID,
			CommunityId:  communityID,
			KnownMediaId: uuid.Must(uuid.NewV7()).String(),
			MagnetUri:    "magnet:?xt=urn:btih:4444444444444444444444444444444444444444&dn=duplicate.txt",
			LibraryId:    libraryID,
		}

		require.NoError(t, community.SyncPublishedContentItem(ctx, db, pc, false))
		require.Equal(t, 1, errorsx.Zero(sqlx.Count(ctx, db, "SELECT COUNT(*) FROM torrents_metadata")))
		require.Equal(t, 1, errorsx.Zero(sqlx.Count(ctx, db, "SELECT COUNT(*) FROM published_content")))

		pc2 := &meta.PublishedContent{
			Id:           uuid.Must(uuid.NewV7()).String(),
			CommunityId:  communityID,
			KnownMediaId: uuid.Must(uuid.NewV7()).String(),
			MagnetUri:    "magnet:?xt=urn:btih:4444444444444444444444444444444444444444&dn=duplicate.txt",
			LibraryId:    uuid.Must(uuid.NewV7()).String(),
		}
		require.NoError(t, community.SyncPublishedContentItem(ctx, db, pc2, false))
		require.Equal(t, 1, errorsx.Zero(sqlx.Count(ctx, db, "SELECT COUNT(*) FROM torrents_metadata")), "should not create duplicate torrent metadata")
		require.Equal(t, 2, errorsx.Zero(sqlx.Count(ctx, db, "SELECT COUNT(*) FROM published_content")), "should create separate published content records")
	})

	t.Run("rejects invalid magnet uri", func(t *testing.T) {
		ctx, cancel := testx.Context(t)
		defer cancel()

		db := sqltestx.Metadatabase(t)
		defer db.Close()

		communityID := uuid.Must(uuid.NewV7()).String()

		pc := &meta.PublishedContent{
			Id:           uuid.Must(uuid.NewV7()).String(),
			CommunityId:  communityID,
			KnownMediaId: uuid.Must(uuid.NewV7()).String(),
			MagnetUri:    "not-a-valid-magnet-uri",
			LibraryId:    uuid.Must(uuid.NewV7()).String(),
		}

		err := community.SyncPublishedContentItem(ctx, db, pc, false)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to parse magnet URI")
	})

	t.Run("syncs with autodownload enabled", func(t *testing.T) {
		ctx, cancel := testx.Context(t)
		defer cancel()

		db := sqltestx.Metadatabase(t)
		defer db.Close()

		communityID := uuid.Must(uuid.NewV7()).String()
		publishedContentID := uuid.Must(uuid.NewV7()).String()

		pc := &meta.PublishedContent{
			Id:           publishedContentID,
			CommunityId:  communityID,
			KnownMediaId: uuid.Must(uuid.NewV7()).String(),
			MagnetUri:    "magnet:?xt=urn:btih:5555555555555555555555555555555555555555&dn=autodownload.txt",
			LibraryId:    uuid.Must(uuid.NewV7()).String(),
		}

		require.NoError(t, community.SyncPublishedContentItem(ctx, db, pc, true))

		var tmeta tracking.Metadata
		require.NoError(t, tracking.MetadataFindByID(ctx, db, errorsx.Must(sqlx.String(ctx, db, "SELECT id::text FROM torrents_metadata"))).Scan(&tmeta))
		require.False(t, tmeta.InitiatedAt.IsZero(), "autodownload should set initiated_at")
	})

	t.Run("preserves known media id from published content", func(t *testing.T) {
		ctx, cancel := testx.Context(t)
		defer cancel()

		db := sqltestx.Metadatabase(t)
		defer db.Close()

		communityID := uuid.Must(uuid.NewV7()).String()
		knownMediaID := uuid.Must(uuid.NewV7()).String()
		publishedContentID := uuid.Must(uuid.NewV7()).String()

		pc := &meta.PublishedContent{
			Id:           publishedContentID,
			CommunityId:  communityID,
			KnownMediaId: knownMediaID,
			MagnetUri:    "magnet:?xt=urn:btih:6666666666666666666666666666666666666666&dn=knownmedia.txt",
			LibraryId:    uuid.Must(uuid.NewV7()).String(),
		}

		require.NoError(t, community.SyncPublishedContentItem(ctx, db, pc, false))

		var tmeta tracking.Metadata
		require.NoError(t, tracking.MetadataFindByID(ctx, db, errorsx.Must(sqlx.String(ctx, db, "SELECT id::text FROM torrents_metadata"))).Scan(&tmeta))
		require.Equal(t, knownMediaID, tmeta.KnownMediaID)
	})
}
