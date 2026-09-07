package communityapi

import (
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/internal/pqueuetestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestSyncPendingToDeeppool(t *testing.T) {
	t.Run("enqueues pending content", func(t *testing.T) {
		var (
			ctx, done   = testx.Context(t)
			q           = sqltestx.Metadatabase(t)
			wq          = pqueuetestx.NewQueue()
			communityID = uuid.Must(uuid.NewV7()).String()
			libraryID   = uuid.Must(uuid.NewV7()).String()
			archiveID   = uuid.Must(uuid.NewV7()).String()
		)
		defer done()

		lmd := library.Metadata{
			ID:             libraryID,
			Description:    "test media",
			ArchiveID:      archiveID,
			TorrentID:      uuid.Nil.String(),
			KnownMediaID:   uuid.Nil.String(),
			DirectoryID:    uuid.Nil.String(),
			EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
		}
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		var pc community.PublishedContent
		require.NoError(t, testx.Fake(&pc, community.PublishedContentOptionTestDefaults, func(p *community.PublishedContent) {
			p.CommunityID = communityID
			p.KnownMediaID = uuid.Must(uuid.NewV7()).String()
			p.MagnetURI = "magnet:?xt=urn:btih:0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33"
			p.LibraryID = libraryID
			p.PublishMode = int32(PublishMode_SYNDICATED)
		}))
		require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, pc).Scan(&pc))

		require.NoError(t, SyncPendingToDeeppool(ctx, q, wq))
		require.Equal(t, 1, wq.Len())
	})

	t.Run("skips already synced content", func(t *testing.T) {
		var (
			ctx, done   = testx.Context(t)
			q           = sqltestx.Metadatabase(t)
			wq          = pqueuetestx.NewQueue()
			communityID = uuid.Must(uuid.NewV7()).String()
			libraryID   = uuid.Must(uuid.NewV7()).String()
			archiveID   = uuid.Must(uuid.NewV7()).String()
		)
		defer done()

		lmd := library.Metadata{
			ID:             libraryID,
			Description:    "test media",
			ArchiveID:      archiveID,
			TorrentID:      uuid.Nil.String(),
			KnownMediaID:   uuid.Nil.String(),
			DirectoryID:    uuid.Nil.String(),
			EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
		}
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		var pc community.PublishedContent
		require.NoError(t, testx.Fake(&pc, community.PublishedContentOptionTestDefaults, func(p *community.PublishedContent) {
			p.CommunityID = communityID
			p.KnownMediaID = uuid.Must(uuid.NewV7()).String()
			p.MagnetURI = "magnet:?xt=urn:btih:0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33"
			p.LibraryID = libraryID
			p.PublishMode = int32(PublishMode_SYNDICATED)
		}))
		require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, pc).Scan(&pc))
		require.NoError(t, community.PublishedContentUpdatePublishedAt(ctx, q, pc.ID, time.Now()).Scan(&pc))

		require.NoError(t, SyncPendingToDeeppool(ctx, q, wq))
		require.Equal(t, 0, wq.Len())
	})
}
