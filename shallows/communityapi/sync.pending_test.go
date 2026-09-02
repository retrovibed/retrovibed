package communityapi

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/james-lawrence/torrent/torrenttest"
	"github.com/retrovibed/retrovibed/retroapi/blockcache"
	"github.com/retrovibed/retrovibed/retroapi/deeppool"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/publishplugin"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

type mockMetrics struct {
	published []*PublishContentRequest
}

func (m *mockMetrics) Publish(ctx context.Context, req *PublishContentRequest, torrent io.Reader) (*PublishContentResponse, error) {
	m.published = append(m.published, req)
	return &PublishContentResponse{}, nil
}

type mockFeedPublisher struct {
	feeds []string
}

func (m *mockFeedPublisher) Find(ctx context.Context, communityID string) (*Community, error) {
	return &Community{Id: communityID}, nil
}

func (m *mockFeedPublisher) UploadFeed(ctx context.Context, communityID string, feed io.Reader) error {
	m.feeds = append(m.feeds, communityID)
	return nil
}

type mockArchiver struct {
	archiveID string
}

func (m *mockArchiver) Upload(ctx context.Context, mimetype string, r io.Reader) (*deeppool.Media, error) {
	if _, err := io.Copy(io.Discard, r); err != nil {
		return nil, err
	}
	return &deeppool.Media{Id: m.archiveID}, nil
}

func TestSyncPendingToDeeppool(t *testing.T) {
	t.Run("skips content without archive", func(t *testing.T) {
		var (
			ctx, done   = testx.Context(t)
			q           = sqltestx.Metadatabase(t)
			mock        = &mockMetrics{}
			feeds       = &mockFeedPublisher{}
			communityID = uuid.Must(uuid.NewV7()).String()
			libraryID   = uuid.Must(uuid.NewV7()).String()
		)
		defer done()

		lmd := library.Metadata{
			ID:             libraryID,
			Description:    "test media",
			ArchiveID:      uuid.Max.String(),
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

		err := SyncPendingToDeeppool(ctx, q, http.DefaultClient, mock, feeds, publishplugin.Unimplemented{}, nil, fsx.DirVirtual(t.TempDir()), fsx.DirVirtual(t.TempDir()))
		require.NoError(t, err)
		require.Empty(t, feeds.feeds)
		require.Empty(t, mock.published)
	})

	t.Run("syncs archived content to deeppool", func(t *testing.T) {
		var (
			ctx, done   = testx.Context(t)
			q           = sqltestx.Metadatabase(t)
			mock        = &mockMetrics{}
			feeds       = &mockFeedPublisher{}
			communityID = uuid.Must(uuid.NewV7()).String()
			libraryID   = uuid.Must(uuid.NewV7()).String()
			archiveID   = uuid.Must(uuid.NewV7()).String()
		)
		defer done()

		lmd := library.Metadata{
			ID:             libraryID,
			Description:    "test media",
			ArchiveID:      archiveID,
			Mimetype:       "video/mp4",
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

		err := SyncPendingToDeeppool(ctx, q, http.DefaultClient, mock, feeds, publishplugin.Unimplemented{}, nil, fsx.DirVirtual(t.TempDir()), fsx.DirVirtual(t.TempDir()))
		require.NoError(t, err)
		require.Contains(t, feeds.feeds, communityID)
		require.Len(t, mock.published, 1)
		require.Equal(t, pc.ID, mock.published[0].PublishedContent.Id)
		require.Equal(t, archiveID, mock.published[0].PublishedContent.ArchivedId)

		var updated community.PublishedContent
		require.NoError(t, community.PublishedContentFindByID(ctx, q, pc.ID).Scan(&updated))
		require.True(t, updated.PublishedAt.Before(time.Now()))
		require.True(t, updated.PublishedAt.After(time.Now().Add(-time.Minute)))
	})

	t.Run("skips already synced content", func(t *testing.T) {
		var (
			ctx, done   = testx.Context(t)
			q           = sqltestx.Metadatabase(t)
			mock        = &mockMetrics{}
			feeds       = &mockFeedPublisher{}
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

		err := SyncPendingToDeeppool(ctx, q, http.DefaultClient, mock, feeds, publishplugin.Unimplemented{}, nil, fsx.DirVirtual(t.TempDir()), fsx.DirVirtual(t.TempDir()))
		require.NoError(t, err)
		require.Empty(t, feeds.feeds)
		require.Empty(t, mock.published)
	})

	t.Run("generates torrent for content with empty magnet uri", func(t *testing.T) {
		var (
			ctx, done   = testx.Context(t)
			q           = sqltestx.Metadatabase(t)
			mock        = &mockMetrics{}
			feeds       = &mockFeedPublisher{}
			communityID = uuid.Must(uuid.NewV7()).String()
			libraryID   = uuid.Must(uuid.NewV7()).String()
			archiveID   = uuid.Must(uuid.NewV7()).String()
			mediaDir    = t.TempDir()
			torrentDir  = t.TempDir()
			testContent = []byte("test media content for torrent generation")
		)
		defer done()

		mediaPath := filepath.Join(mediaDir, libraryID)
		require.NoError(t, os.MkdirAll(mediaPath, 0700))
		require.NoError(t, os.WriteFile(filepath.Join(mediaPath, "0"), testContent, 0600))

		lmd := library.Metadata{
			ID:             libraryID,
			Description:    "test media",
			ArchiveID:      archiveID,
			Bytes:          uint64(len(testContent)),
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
			p.MagnetURI = ""
			p.LibraryID = libraryID
			p.PublishMode = int32(PublishMode_SYNDICATED)
		}))
		require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, pc).Scan(&pc))

		err := SyncPendingToDeeppool(ctx, q, http.DefaultClient, mock, feeds, publishplugin.Unimplemented{}, nil, fsx.DirVirtual(mediaDir), fsx.DirVirtual(torrentDir))
		require.NoError(t, err)
		require.Contains(t, feeds.feeds, communityID)
		require.Len(t, mock.published, 1)
		require.True(t, strings.HasPrefix(mock.published[0].PublishedContent.MagnetUri, "magnet:?"))

		var updated community.PublishedContent
		require.NoError(t, community.PublishedContentFindByID(ctx, q, pc.ID).Scan(&updated))
		require.True(t, strings.HasPrefix(updated.MagnetURI, "magnet:?"))
	})

	t.Run("full async flow: publish without torrent then sync generates and archives", func(t *testing.T) {
		var (
			ctx, done   = testx.Context(t)
			q           = sqltestx.Metadatabase(t)
			mock        = &mockMetrics{}
			feeds       = &mockFeedPublisher{}
			communityID = uuid.Must(uuid.NewV7()).String()
			libraryID   = uuid.Must(uuid.NewV7()).String()
			mediaDir    = t.TempDir()
			torrentDir  = t.TempDir()
			testContent = []byte("test media content for full async flow")
		)
		defer done()

		mediaPath := filepath.Join(mediaDir, libraryID)
		require.NoError(t, os.MkdirAll(mediaPath, 0700))
		require.NoError(t, os.WriteFile(filepath.Join(mediaPath, "0"), testContent, 0600))

		lmd := library.Metadata{
			ID:             libraryID,
			Description:    "test media for async",
			ArchiveID:      uuid.Max.String(),
			Bytes:          uint64(len(testContent)),
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
			p.MagnetURI = ""
			p.LibraryID = libraryID
			p.PublishMode = int32(PublishMode_SYNDICATED)
		}))
		require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, pc).Scan(&pc))

		err := SyncPendingToDeeppool(ctx, q, http.DefaultClient, mock, feeds, publishplugin.Unimplemented{}, nil, fsx.DirVirtual(mediaDir), fsx.DirVirtual(torrentDir))
		require.NoError(t, err)
		require.Empty(t, feeds.feeds)
		require.Empty(t, mock.published)

		var afterFirstSync community.PublishedContent
		require.NoError(t, community.PublishedContentFindByID(ctx, q, pc.ID).Scan(&afterFirstSync))
		require.True(t, strings.HasPrefix(afterFirstSync.MagnetURI, "magnet:?"))

		archiveID := uuid.Must(uuid.NewV7()).String()
		require.NoError(t, library.MetadataArchivedByID(ctx, q, libraryID, archiveID, 0).Scan(&lmd))

		err = SyncPendingToDeeppool(ctx, q, http.DefaultClient, mock, feeds, publishplugin.Unimplemented{}, nil, fsx.DirVirtual(mediaDir), fsx.DirVirtual(torrentDir))
		require.NoError(t, err)
		require.Contains(t, feeds.feeds, communityID)
		require.Len(t, mock.published, 1)
		require.Equal(t, archiveID, mock.published[0].PublishedContent.ArchivedId)
		require.True(t, strings.HasPrefix(mock.published[0].PublishedContent.MagnetUri, "magnet:?"))

		var finalPC community.PublishedContent
		require.NoError(t, community.PublishedContentFindByID(ctx, q, pc.ID).Scan(&finalPC))
		require.True(t, finalPC.PublishedAt.Before(time.Now()))
	})

	t.Run("uses existing torrent when torrent_id is set", func(t *testing.T) {
		var (
			ctx, done   = testx.Context(t)
			q           = sqltestx.Metadatabase(t)
			mock        = &mockMetrics{}
			feeds       = &mockFeedPublisher{}
			communityID = uuid.Must(uuid.NewV7()).String()
			libraryID   = uuid.Must(uuid.NewV7()).String()
			archiveID   = uuid.Must(uuid.NewV7()).String()
			mediaDir    = t.TempDir()
			torrentDir  = t.TempDir()
			ih          = int160.FromBytes([]byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd})
		)
		defer done()

		tmd := tracking.NewMetadata(
			&ih,
			tracking.MetadataOptionDescription("existing torrent"),
			tracking.MetadataOptionKnownMediaID(uuid.Nil.String()),
			tracking.MetadataOptionBytes(1024),
			tracking.MetadataOptionDownloaded(1024),
			tracking.MetadataOptionAvailable(1024),
			tracking.MetadataOptionCompleted,
			tracking.MetadataOptionAutoSeeding,
		)
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, tmd).Scan(&tmd))

		torrentFile, err := os.Create(filepath.Join(torrentDir, ih.String()+tracking.TorrentSuffix))
		require.NoError(t, err)
		mi := metainfo.MetaInfo{}
		require.NoError(t, mi.Write(torrentFile))
		require.NoError(t, torrentFile.Close())

		lmd := library.Metadata{
			ID:             libraryID,
			Description:    "test media with existing torrent",
			ArchiveID:      archiveID,
			TorrentID:      tmd.ID,
			KnownMediaID:   uuid.Nil.String(),
			DirectoryID:    uuid.Nil.String(),
			Bytes:          1024,
			EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
		}
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		var pc community.PublishedContent
		require.NoError(t, testx.Fake(&pc, community.PublishedContentOptionTestDefaults, func(p *community.PublishedContent) {
			p.CommunityID = communityID
			p.KnownMediaID = uuid.Must(uuid.NewV7()).String()
			p.MagnetURI = ""
			p.LibraryID = libraryID
			p.PublishMode = int32(PublishMode_SYNDICATED)
		}))
		require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, pc).Scan(&pc))

		err = SyncPendingToDeeppool(ctx, q, http.DefaultClient, mock, feeds, publishplugin.Unimplemented{}, nil, fsx.DirVirtual(mediaDir), fsx.DirVirtual(torrentDir))
		require.NoError(t, err)
		require.Contains(t, feeds.feeds, communityID)
		require.Len(t, mock.published, 1)

		var updated community.PublishedContent
		require.NoError(t, community.PublishedContentFindByID(ctx, q, pc.ID).Scan(&updated))
		require.Contains(t, updated.MagnetURI, "aabbccddeeff00112233445566778899aabbccdd")
	})

	t.Run("skips unlisted content", func(t *testing.T) {
		var (
			ctx, done   = testx.Context(t)
			q           = sqltestx.Metadatabase(t)
			mock        = &mockMetrics{}
			feeds       = &mockFeedPublisher{}
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
			p.PublishMode = int32(PublishMode_UNLISTED)
		}))
		require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, pc).Scan(&pc))

		err := SyncPendingToDeeppool(ctx, q, http.DefaultClient, mock, feeds, publishplugin.Unimplemented{}, nil, fsx.DirVirtual(t.TempDir()), fsx.DirVirtual(t.TempDir()))
		require.NoError(t, err)
		require.Empty(t, feeds.feeds)
		require.Empty(t, mock.published)
	})

	t.Run("marks existing torrent as completed and seeding", func(t *testing.T) {
		var (
			ctx, done   = testx.Context(t)
			q           = sqltestx.Metadatabase(t)
			mock        = &mockMetrics{}
			feeds       = &mockFeedPublisher{}
			communityID = uuid.Must(uuid.NewV7()).String()
			libraryID   = uuid.Must(uuid.NewV7()).String()
			archiveID   = uuid.Must(uuid.NewV7()).String()
			mvfs        = fsx.DirVirtual(t.TempDir())
			tvfs        = fsx.DirVirtual(t.TempDir())
		)
		defer done()

		// generate torrent data using the standard test helper
		dataDir := t.TempDir()
		info, _, err := torrenttest.Random(dataDir, 128*bytesx.KiB)
		require.NoError(t, err)

		infobytes, err := metainfo.Encode(*info)
		require.NoError(t, err)

		mi := metainfo.MetaInfo{InfoBytes: infobytes}
		infohash := mi.HashInfoBytes()
		ih := int160.FromByteArray(infohash)
		infohashHex := ih.String()

		// write the .torrent file into tvfs
		f, err := os.Create(tvfs.Path(infohashHex + tracking.TorrentSuffix))
		require.NoError(t, err)
		require.NoError(t, mi.Write(f))
		require.NoError(t, f.Close())

		// copy the flat torrent file into the block cache in tvfs
		src, err := os.Open(filepath.Join(dataDir, infohashHex))
		require.NoError(t, err)
		defer src.Close()
		bc, err := blockcache.NewDirectoryCache(tvfs.Path(infohashHex))
		require.NoError(t, err)
		_, err = io.Copy(io.NewOffsetWriter(bc, 0), src)
		require.NoError(t, err)

		// create tracking metadata that is not yet seeding
		tmd := tracking.NewMetadata(&ih, tracking.MetadataOptionBytes(info.TotalLength()))
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, tmd).Scan(&tmd))
		require.False(t, tmd.Seeding)
		require.Equal(t, timex.Inf(), tmd.CompletedAt)

		lmd := library.Metadata{
			ID:             libraryID,
			Description:    "test media with incomplete torrent",
			ArchiveID:      archiveID,
			TorrentID:      tmd.ID,
			KnownMediaID:   uuid.Nil.String(),
			DirectoryID:    uuid.Nil.String(),
			Bytes:          uint64(info.TotalLength()),
			EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
		}
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		var pc community.PublishedContent
		require.NoError(t, testx.Fake(&pc, community.PublishedContentOptionTestDefaults, func(p *community.PublishedContent) {
			p.CommunityID = communityID
			p.KnownMediaID = uuid.Must(uuid.NewV7()).String()
			p.MagnetURI = ""
			p.LibraryID = libraryID
			p.PublishMode = int32(PublishMode_SYNDICATED)
		}))
		require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, pc).Scan(&pc))

		err = SyncPendingToDeeppool(ctx, q, http.DefaultClient, mock, feeds, publishplugin.Unimplemented{}, nil, mvfs, tvfs)
		require.NoError(t, err)

		var afterSync tracking.Metadata
		require.NoError(t, tracking.MetadataFindByID(ctx, q, tmd.ID).Scan(&afterSync))
		require.True(t, afterSync.Seeding)
		require.NotEqual(t, timex.Inf(), afterSync.CompletedAt)
	})

	t.Run("same content produces same magnet uri", func(t *testing.T) {
		var (
			ctx, done   = testx.Context(t)
			q           = sqltestx.Metadatabase(t)
			mock        = &mockMetrics{}
			feeds       = &mockFeedPublisher{}
			communityID = uuid.Must(uuid.NewV7()).String()
			libraryID   = uuid.Must(uuid.NewV7()).String()
			mediaDir    = t.TempDir()
			torrentDir  = t.TempDir()
			testContent = []byte("test media content for idempotency check")
		)
		defer done()

		mediaPath := filepath.Join(mediaDir, libraryID)
		require.NoError(t, os.MkdirAll(mediaPath, 0700))
		require.NoError(t, os.WriteFile(filepath.Join(mediaPath, "0"), testContent, 0600))

		lmd := library.Metadata{
			ID:             libraryID,
			Description:    "test media",
			ArchiveID:      uuid.Max.String(),
			Bytes:          uint64(len(testContent)),
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
			p.MagnetURI = ""
			p.LibraryID = libraryID
			p.PublishMode = int32(PublishMode_SYNDICATED)
		}))
		require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, pc).Scan(&pc))

		mvfs := fsx.DirVirtual(mediaDir)
		tvfs := fsx.DirVirtual(torrentDir)

		err := SyncPendingToDeeppool(ctx, q, http.DefaultClient, mock, feeds, publishplugin.Unimplemented{}, nil, mvfs, tvfs)
		require.NoError(t, err)

		var afterFirst community.PublishedContent
		require.NoError(t, community.PublishedContentFindByID(ctx, q, pc.ID).Scan(&afterFirst))
		require.NotEmpty(t, afterFirst.MagnetURI)

		err = SyncPendingToDeeppool(ctx, q, http.DefaultClient, mock, feeds, publishplugin.Unimplemented{}, nil, mvfs, tvfs)
		require.NoError(t, err)

		var afterSecond community.PublishedContent
		require.NoError(t, community.PublishedContentFindByID(ctx, q, pc.ID).Scan(&afterSecond))
		require.Equal(t, afterFirst.MagnetURI, afterSecond.MagnetURI)
		require.True(t, afterSecond.UpdatedAt.After(afterFirst.UpdatedAt))
	})

	t.Run("different content produces different magnet uri", func(t *testing.T) {
		var (
			ctx, done   = testx.Context(t)
			q           = sqltestx.Metadatabase(t)
			mock        = &mockMetrics{}
			feeds       = &mockFeedPublisher{}
			communityID = uuid.Must(uuid.NewV7()).String()
			libraryID   = uuid.Must(uuid.NewV7()).String()
			mediaDir    = t.TempDir()
			torrentDir  = t.TempDir()
			testContent = []byte("test media content for description change test")
		)
		defer done()

		mediaPath := filepath.Join(mediaDir, libraryID)
		require.NoError(t, os.MkdirAll(mediaPath, 0700))
		require.NoError(t, os.WriteFile(filepath.Join(mediaPath, "0"), testContent, 0600))

		lmd := library.Metadata{
			ID:             libraryID,
			Description:    "description alpha",
			ArchiveID:      uuid.Max.String(),
			Bytes:          uint64(len(testContent)),
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
			p.LibraryID = libraryID
			p.PublishMode = int32(PublishMode_SYNDICATED)
		}))
		require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, pc).Scan(&pc))

		mvfs := fsx.DirVirtual(mediaDir)
		tvfs := fsx.DirVirtual(torrentDir)
		err := SyncPendingToDeeppool(ctx, q, http.DefaultClient, mock, feeds, publishplugin.Unimplemented{}, nil, mvfs, tvfs)
		require.NoError(t, err)

		require.Equal(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM torrents_metadata"))
		require.Equal(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM library_metadata"))

		var afterFirst community.PublishedContent
		require.NoError(t, community.PublishedContentFindByID(ctx, q, pc.ID).Scan(&afterFirst))
		require.NotEmpty(t, afterFirst.MagnetURI)

		require.NoError(t, library.MetadataUpdateDescriptionByID(ctx, q, libraryID, "description beta").Scan(&lmd))

		err = SyncPendingToDeeppool(ctx, q, http.DefaultClient, mock, feeds, publishplugin.Unimplemented{}, nil, mvfs, tvfs)
		require.NoError(t, err)

		require.Equal(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM torrents_metadata"))
		require.Equal(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM library_metadata"))

		var afterSecond community.PublishedContent
		require.NoError(t, community.PublishedContentFindByID(ctx, q, pc.ID).Scan(&afterSecond))
		require.NotEmpty(t, afterSecond.MagnetURI)
		require.NotEqual(t, afterFirst.MagnetURI, afterSecond.MagnetURI)
	})

	t.Run("syncs listed content without metrics publish", func(t *testing.T) {
		var (
			ctx, done   = testx.Context(t)
			q           = sqltestx.Metadatabase(t)
			mock        = &mockMetrics{}
			feeds       = &mockFeedPublisher{}
			communityID = uuid.Must(uuid.NewV7()).String()
			libraryID   = uuid.Must(uuid.NewV7()).String()
		)
		defer done()

		lmd := library.Metadata{
			ID:             libraryID,
			Description:    "test media",
			ArchiveID:      uuid.Nil.String(),
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
			p.PublishMode = int32(PublishMode_LISTED)
		}))
		require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, pc).Scan(&pc))

		err := SyncPendingToDeeppool(ctx, q, http.DefaultClient, mock, feeds, publishplugin.Unimplemented{}, nil, fsx.DirVirtual(t.TempDir()), fsx.DirVirtual(t.TempDir()))
		require.NoError(t, err)
		require.Contains(t, feeds.feeds, communityID)
		require.Empty(t, mock.published)

		var updated community.PublishedContent
		require.NoError(t, community.PublishedContentFindByID(ctx, q, pc.ID).Scan(&updated))
		require.True(t, updated.PublishedAt.Before(time.Now()))
		require.True(t, updated.PublishedAt.After(time.Now().Add(-time.Minute)))
	})

	t.Run("prefers known mimetype over library mimetype", func(t *testing.T) {
		var (
			ctx, done   = testx.Context(t)
			q           = sqltestx.Metadatabase(t)
			mock        = &mockMetrics{}
			feeds       = &mockFeedPublisher{}
			communityID = uuid.Must(uuid.NewV7()).String()
			libraryID   = uuid.Must(uuid.NewV7()).String()
			archiveID   = uuid.Must(uuid.NewV7()).String()
		)
		defer done()

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		known.Mimetype = mimex.Video
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		lmd := library.Metadata{
			ID:             libraryID,
			Description:    "test media",
			ArchiveID:      archiveID,
			Mimetype:       "video/mp4",
			TorrentID:      uuid.Nil.String(),
			KnownMediaID:   known.UID,
			DirectoryID:    uuid.Nil.String(),
			EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
		}
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		var pc community.PublishedContent
		require.NoError(t, testx.Fake(&pc, community.PublishedContentOptionTestDefaults, func(p *community.PublishedContent) {
			p.CommunityID = communityID
			p.KnownMediaID = known.UID
			p.MagnetURI = "magnet:?xt=urn:btih:0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33"
			p.LibraryID = libraryID
			p.PublishMode = int32(PublishMode_SYNDICATED)
		}))
		require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, pc).Scan(&pc))

		err := SyncPendingToDeeppool(ctx, q, http.DefaultClient, mock, feeds, publishplugin.Unimplemented{}, nil, fsx.DirVirtual(t.TempDir()), fsx.DirVirtual(t.TempDir()))
		require.NoError(t, err)
		require.Contains(t, feeds.feeds, communityID)
		require.Len(t, mock.published, 1)
		require.Equal(t, mimex.Video, mock.published[0].PublishedContent.Mimetype)
	})

	t.Run("passes encryption_seed from library metadata", func(t *testing.T) {
		var (
			ctx, done      = testx.Context(t)
			q              = sqltestx.Metadatabase(t)
			mock           = &mockMetrics{}
			feeds          = &mockFeedPublisher{}
			communityID    = uuid.Must(uuid.NewV7()).String()
			libraryID      = uuid.Must(uuid.NewV7()).String()
			archiveID      = uuid.Must(uuid.NewV7()).String()
			encryptionSeed = uuid.Must(uuid.NewV4()).String()
		)
		defer done()

		lmd := library.Metadata{
			ID:             libraryID,
			Description:    "test media",
			ArchiveID:      archiveID,
			Mimetype:       "video/mp4",
			TorrentID:      uuid.Nil.String(),
			KnownMediaID:   uuid.Nil.String(),
			DirectoryID:    uuid.Nil.String(),
			EncryptionSeed: encryptionSeed,
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

		err := SyncPendingToDeeppool(ctx, q, http.DefaultClient, mock, feeds, publishplugin.Unimplemented{}, nil, fsx.DirVirtual(t.TempDir()), fsx.DirVirtual(t.TempDir()))
		require.NoError(t, err)
		require.Contains(t, feeds.feeds, communityID)
		require.Len(t, mock.published, 1)
		require.Equal(t, encryptionSeed, mock.published[0].PublishedContent.EncryptionSeed)
	})

	t.Run("skips syndicated content with archival in progress", func(t *testing.T) {
		var (
			ctx, done   = testx.Context(t)
			q           = sqltestx.Metadatabase(t)
			mock        = &mockMetrics{}
			feeds       = &mockFeedPublisher{}
			communityID = uuid.Must(uuid.NewV7()).String()
			libraryID   = uuid.Must(uuid.NewV7()).String()
		)
		defer done()

		lmd := library.Metadata{
			ID:             libraryID,
			Description:    "test media",
			ArchiveID:      uuid.Max.String(),
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

		err := SyncPendingToDeeppool(ctx, q, http.DefaultClient, mock, feeds, publishplugin.Unimplemented{}, nil, fsx.DirVirtual(t.TempDir()), fsx.DirVirtual(t.TempDir()))
		require.NoError(t, err)
		require.Empty(t, feeds.feeds)
		require.Empty(t, mock.published)

		// published_at must remain unset (item is still pending)
		var updated community.PublishedContent
		require.NoError(t, community.PublishedContentFindByID(ctx, q, pc.ID).Scan(&updated))
		require.True(t, updated.PublishedAt.Equal(timex.Inf()))
	})

	t.Run("archives content with nil archive_id then publishes to deeppool", func(t *testing.T) {
		var (
			ctx, done   = testx.Context(t)
			q           = sqltestx.Metadatabase(t)
			mock        = &mockMetrics{}
			feeds       = &mockFeedPublisher{}
			archiveID   = uuid.Must(uuid.NewV7()).String()
			archiver    = &mockArchiver{archiveID: archiveID}
			communityID = uuid.Must(uuid.NewV7()).String()
			libraryID   = uuid.Must(uuid.NewV7()).String()
			mediaDir    = t.TempDir()
			torrentDir  = t.TempDir()
			testContent = []byte("test media content for inline archival")
		)
		defer done()

		mediaPath := filepath.Join(mediaDir, libraryID)
		require.NoError(t, os.MkdirAll(mediaPath, 0700))
		require.NoError(t, os.WriteFile(filepath.Join(mediaPath, "0"), testContent, 0600))

		lmd := library.Metadata{
			ID:             libraryID,
			Description:    "test media",
			ArchiveID:      uuid.Nil.String(),
			Bytes:          uint64(len(testContent)),
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
			p.MagnetURI = ""
			p.LibraryID = libraryID
			p.PublishMode = int32(PublishMode_SYNDICATED)
		}))
		require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, pc).Scan(&pc))

		err := SyncPendingToDeeppool(ctx, q, http.DefaultClient, mock, feeds, publishplugin.Unimplemented{}, archiver, fsx.DirVirtual(mediaDir), fsx.DirVirtual(torrentDir))
		require.NoError(t, err)
		require.Contains(t, feeds.feeds, communityID)
		require.Len(t, mock.published, 1)
		require.Equal(t, archiveID, mock.published[0].PublishedContent.ArchivedId)

		// library metadata must have the archive_id set
		var updatedLMD library.Metadata
		require.NoError(t, library.MetadataFindByID(ctx, q, libraryID).Scan(&updatedLMD))
		require.Equal(t, archiveID, updatedLMD.ArchiveID)

		// published_at must be set
		var updatedPC community.PublishedContent
		require.NoError(t, community.PublishedContentFindByID(ctx, q, pc.ID).Scan(&updatedPC))
		require.True(t, updatedPC.PublishedAt.Before(time.Now()))
		require.True(t, updatedPC.PublishedAt.After(time.Now().Add(-time.Minute)))
	})

	t.Run("does not leak a still-pending sibling into a feed regenerated mid-sync", func(t *testing.T) {
		var (
			ctx, done   = testx.Context(t)
			q           = sqltestx.Metadatabase(t)
			mock        = &mockMetrics{}
			feeds       = &mockFeedPublisher{}
			communityID = uuid.Must(uuid.NewV7()).String()
			mediaDir    = t.TempDir()
			torrentDir  = t.TempDir()
			testContent = []byte("test media content for mid-sync feed race")
		)
		defer done()

		// pc1: fully resolvable, completes within this sync pass and
		// triggers a feed regeneration for communityID.
		var lmd1 library.Metadata
		require.NoError(t, testx.Fake(&lmd1, library.MetadataOptionTestDefaults, func(m *library.Metadata) {
			m.ID = uuid.Must(uuid.NewV7()).String()
			m.Description = "ready media"
			m.Bytes = uint64(len(testContent))
		}))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd1).Scan(&lmd1))

		mediaPath := filepath.Join(mediaDir, lmd1.ID)
		require.NoError(t, os.MkdirAll(mediaPath, 0700))
		require.NoError(t, os.WriteFile(filepath.Join(mediaPath, "0"), testContent, 0600))

		var pc1 community.PublishedContent
		require.NoError(t, testx.Fake(&pc1, community.PublishedContentOptionTestDefaults, func(p *community.PublishedContent) {
			p.CommunityID = communityID
			p.MagnetURI = ""
			p.LibraryID = lmd1.ID
			p.PublishMode = int32(PublishMode_LISTED)
		}))
		require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, pc1).Scan(&pc1))

		// pc2: same community, but its media is never written to mediaDir,
		// so ensureTorrent fails every pass and it never reaches
		// published_at - it stays permanently pending, exactly like a
		// sibling row mid-flight when a regeneration fires.
		var lmd2 library.Metadata
		require.NoError(t, testx.Fake(&lmd2, library.MetadataOptionTestDefaults, func(m *library.Metadata) {
			m.ID = uuid.Must(uuid.NewV7()).String()
			m.Description = "pending media"
			m.Bytes = bytesx.KiB
		}))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd2).Scan(&lmd2))

		var pc2 community.PublishedContent
		require.NoError(t, testx.Fake(&pc2, community.PublishedContentOptionTestDefaults, func(p *community.PublishedContent) {
			p.CommunityID = communityID
			p.MagnetURI = ""
			p.LibraryID = lmd2.ID
			p.PublishMode = int32(PublishMode_SYNDICATED)
		}))
		require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, pc2).Scan(&pc2))

		err := SyncPendingToDeeppool(ctx, q, http.DefaultClient, mock, feeds, publishplugin.Unimplemented{}, nil, fsx.DirVirtual(mediaDir), fsx.DirVirtual(torrentDir))
		require.NoError(t, err)
		require.Contains(t, feeds.feeds, communityID)

		// pc2 must still be pending: this is the state a mid-sync feed
		// regeneration would have observed for it.
		var updatedPC2 community.PublishedContent
		require.NoError(t, community.PublishedContentFindByID(ctx, q, pc2.ID).Scan(&updatedPC2))
		require.True(t, updatedPC2.PublishedAt.Equal(timex.Inf()))

		items, err := buildFeedItems(ctx, q, &Community{Id: communityID})
		require.NoError(t, err)
		require.Len(t, items, 1)
		require.Equal(t, pc1.ID, items[0].Guid)
	})
}
