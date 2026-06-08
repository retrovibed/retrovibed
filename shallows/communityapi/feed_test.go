package communityapi

import (
	"bytes"
	"slices"
	"testing"
	"time"

	"github.com/retrovibed/retrovibed/shallows/community"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"

	"github.com/retrovibed/retrovibed/shallows/rss"
	"github.com/stretchr/testify/require"
)

func TestFeedGeneration(t *testing.T) {
	t.Run("generates RSS feed XML", func(t *testing.T) {
		var (
			ctx, done   = testx.Context(t)
			q           = sqltestx.Metadatabase(t)
			communityID = uuid.Must(uuid.NewV7()).String()
			libraryID1  = uuid.Must(uuid.NewV7()).String()
			libraryID2  = uuid.Must(uuid.NewV7()).String()
			knownID     = uuid.Must(uuid.NewV7()).String()
		)
		defer done()

		known := library.Known{
			ID:            knownID,
			UID:           knownID,
			Md5:           uuid.Must(uuid.NewV4()).String(),
			Title:         "Test Movie",
			OriginalTitle: "Test Movie Original",
			Overview:      "A test movie for feed generation",
		}
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		lmd1 := library.Metadata{
			ID:             libraryID1,
			Description:    "test-media-1.mp4",
			Bytes:          1024 * 1024 * 100,
			Mimetype:       "video/mp4",
			TorrentID:      uuid.Nil.String(),
			KnownMediaID:   known.UID,
			ArchiveID:      uuid.Must(uuid.NewV7()).String(),
			EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
		}
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd1).Scan(&lmd1))

		lmd2 := library.Metadata{
			ID:             libraryID2,
			Description:    "unknown-media.mkv",
			Bytes:          1024 * 1024 * 500,
			Mimetype:       "video/x-matroska",
			TorrentID:      uuid.Nil.String(),
			KnownMediaID:   uuid.Nil.String(),
			ArchiveID:      uuid.Must(uuid.NewV7()).String(),
			EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
		}
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd2).Scan(&lmd2))

		pc1 := community.NewPublishedContent(community.PublishedContent{
			CommunityID:  communityID,
			KnownMediaID: known.UID,
			MagnetURI:    "magnet:?xt=urn:btih:0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33&dn=Test+Movie",
			LibraryID:    libraryID1,
			PublishMode:  int32(PublishMode_LISTED),
		})
		require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, pc1).Scan(&pc1))
		require.NoError(t, community.PublishedContentUpdatePublishedAt(ctx, q, pc1.ID, time.Now().Add(-time.Hour)).Scan(&pc1))

		pc2 := community.NewPublishedContent(community.PublishedContent{
			CommunityID: communityID,
			MagnetURI:   "magnet:?xt=urn:btih:1234567890abcdef1234567890abcdef12345678&dn=Unknown+Media",
			LibraryID:   libraryID2,
			PublishMode: int32(PublishMode_SYNDICATED),
		})
		require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, pc2).Scan(&pc2))
		require.NoError(t, community.PublishedContentUpdatePublishedAt(ctx, q, pc2.ID, time.Now()).Scan(&pc2))

		community := &Community{
			Id:          communityID,
			Domain:      "testcommunity",
			Description: "A test community for RSS feed generation",
			Entropy:     "test-entropy-value",
			Mimetype:    "video/*",
		}

		items, err := buildFeedItems(ctx, q, community)
		require.NoError(t, err)
		require.Len(t, items, 2)

		buf := new(bytes.Buffer)
		channel := rss.Channel{
			Title:         community.Domain,
			Link:          "https://testcommunity.community.retrovibe.space",
			Description:   community.Description,
			TTL:           feedDefaultTTL,
			LastBuildDate: time.Now().UTC(),
			Language:      feedDefaultLanguage,
			Retrovibed:    &rss.Retrovibed{Entropy: community.Entropy, Mimetype: community.Mimetype},
		}
		require.NoError(t, rss.Generator().Generate(buf, channel, slices.Values(items)))

		t.Logf("\n%s", buf.String())
	})

	t.Run("excludes unlisted content from feed", func(t *testing.T) {
		var (
			ctx, done   = testx.Context(t)
			q           = sqltestx.Metadatabase(t)
			communityID = uuid.Must(uuid.NewV7()).String()
			libraryID1  = uuid.Must(uuid.NewV7()).String()
			libraryID2  = uuid.Must(uuid.NewV7()).String()
		)
		defer done()

		lmd1 := library.Metadata{
			ID:             libraryID1,
			Description:    "unlisted-media.mp4",
			Bytes:          1024,
			Mimetype:       "video/mp4",
			TorrentID:      uuid.Nil.String(),
			KnownMediaID:   uuid.Nil.String(),
			ArchiveID:      uuid.Nil.String(),
			EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
		}
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd1).Scan(&lmd1))

		lmd2 := library.Metadata{
			ID:             libraryID2,
			Description:    "hosted-media.mp4",
			Bytes:          2048,
			Mimetype:       "video/mp4",
			TorrentID:      uuid.Nil.String(),
			KnownMediaID:   uuid.Nil.String(),
			ArchiveID:      uuid.Nil.String(),
			EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
		}
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd2).Scan(&lmd2))

		pc1 := community.NewPublishedContent(community.PublishedContent{
			CommunityID: communityID,
			MagnetURI:   "magnet:?xt=urn:btih:0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33",
			LibraryID:   libraryID1,
			PublishMode: int32(PublishMode_UNLISTED),
		})
		require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, pc1).Scan(&pc1))
		require.NoError(t, community.PublishedContentUpdatePublishedAt(ctx, q, pc1.ID, time.Now()).Scan(&pc1))

		pc2 := community.NewPublishedContent(community.PublishedContent{
			CommunityID: communityID,
			MagnetURI:   "magnet:?xt=urn:btih:1234567890abcdef1234567890abcdef12345678",
			LibraryID:   libraryID2,
			PublishMode: int32(PublishMode_LISTED),
		})
		require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, pc2).Scan(&pc2))
		require.NoError(t, community.PublishedContentUpdatePublishedAt(ctx, q, pc2.ID, time.Now()).Scan(&pc2))

		community := &Community{
			Id:          communityID,
			Domain:      "testcommunity",
			Description: "test",
			Entropy:     "test-entropy",
			Mimetype:    "video/*",
		}

		items, err := buildFeedItems(ctx, q, community)
		require.NoError(t, err)
		require.Len(t, items, 1)
		require.Equal(t, pc2.ID, items[0].Guid)
	})
}
