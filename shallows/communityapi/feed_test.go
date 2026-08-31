package communityapi

import (
	"bytes"
	"context"
	"io"
	"slices"
	"testing"
	"time"

	"github.com/retrovibed/retrovibed/shallows/community"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"

	"github.com/retrovibed/retrovibed/shallows/rss"
	"github.com/stretchr/testify/require"
)

type fakeFeedPublisher struct {
	community *Community
	uploaded  []byte
}

func (f *fakeFeedPublisher) Find(ctx context.Context, communityID string) (*Community, error) {
	return f.community, nil
}

func (f *fakeFeedPublisher) UploadFeed(ctx context.Context, communityID string, feed io.Reader) error {
	b, err := io.ReadAll(feed)
	if err != nil {
		return err
	}
	f.uploaded = b
	return nil
}

func TestRegenerateFeed(t *testing.T) {
	t.Run("uses the community subdomain label as title when url is hosted on our domain", func(t *testing.T) {
		var (
			ctx, done   = testx.Context(t)
			q           = sqltestx.Metadatabase(t)
			communityID = uuid.Must(uuid.NewV7()).String()
		)
		defer done()

		publisher := &fakeFeedPublisher{
			community: &Community{
				Id:          communityID,
				Url:         "https://mysite.community.retrovibe.space",
				Description: "test",
				Entropy:     "test-entropy",
				Mimetype:    "video/*",
			},
		}

		require.NoError(t, RegenerateFeed(ctx, q, publisher, communityID))
		require.Contains(t, string(publisher.uploaded), "<title>mysite</title>")
		require.Contains(t, string(publisher.uploaded), "<link>https://mysite.community.retrovibe.space</link>")
	})

	t.Run("uses the fqdn as title when it is a custom domain", func(t *testing.T) {
		var (
			ctx, done   = testx.Context(t)
			q           = sqltestx.Metadatabase(t)
			communityID = uuid.Must(uuid.NewV7()).String()
		)
		defer done()

		publisher := &fakeFeedPublisher{
			community: &Community{
				Id:          communityID,
				Url:         "https://mysite.example.com",
				Description: "test",
				Entropy:     "test-entropy",
				Mimetype:    "video/*",
			},
		}

		// CommunityDomainFromURL only extracts a short label for our own hosted
		// subdomains; a custom domain has no reliable label to derive so its
		// fqdn is used as the title instead.
		require.NoError(t, RegenerateFeed(ctx, q, publisher, communityID))
		require.Contains(t, string(publisher.uploaded), "<title>mysite.example.com</title>")
		require.Contains(t, string(publisher.uploaded), "<link>https://mysite.example.com</link>")
	})
}

func TestFeedGeneration(t *testing.T) {
	t.Run("generates RSS feed XML", func(t *testing.T) {
		var (
			ctx, done   = testx.Context(t)
			q           = sqltestx.Metadatabase(t)
			communityID = uuid.Must(uuid.NewV7()).String()
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

		var lmd1 library.Metadata
		require.NoError(t, testx.Fake(&lmd1, library.MetadataOptionTestDefaults, func(m *library.Metadata) {
			m.ID = uuid.Must(uuid.NewV7()).String()
			m.Description = "test-media-1.mp4"
			m.Bytes = 100 * bytesx.MiB
			m.Mimetype = "video/mp4"
			m.KnownMediaID = known.UID
		}))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd1).Scan(&lmd1))

		var lmd2 library.Metadata
		require.NoError(t, testx.Fake(&lmd2, library.MetadataOptionTestDefaults, func(m *library.Metadata) {
			m.ID = uuid.Must(uuid.NewV7()).String()
			m.Description = "unknown-media.mkv"
			m.Bytes = 500 * bytesx.MiB
			m.Mimetype = "video/x-matroska"
		}))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd2).Scan(&lmd2))

		var pc1 community.PublishedContent
		require.NoError(t, testx.Fake(&pc1, community.PublishedContentOptionTestDefaults, func(p *community.PublishedContent) {
			p.CommunityID = communityID
			p.KnownMediaID = known.UID
			p.MagnetURI = "magnet:?xt=urn:btih:0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33&dn=Test+Movie"
			p.LibraryID = lmd1.ID
			p.PublishMode = int32(PublishMode_LISTED)
		}))
		require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, pc1).Scan(&pc1))
		require.NoError(t, community.PublishedContentUpdatePublishedAt(ctx, q, pc1.ID, time.Now().Add(-time.Hour)).Scan(&pc1))

		var pc2 community.PublishedContent
		require.NoError(t, testx.Fake(&pc2, community.PublishedContentOptionTestDefaults, func(p *community.PublishedContent) {
			p.CommunityID = communityID
			p.MagnetURI = "magnet:?xt=urn:btih:1234567890abcdef1234567890abcdef12345678&dn=Unknown+Media"
			p.LibraryID = lmd2.ID
			p.PublishMode = int32(PublishMode_SYNDICATED)
		}))
		require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, pc2).Scan(&pc2))
		require.NoError(t, community.PublishedContentUpdatePublishedAt(ctx, q, pc2.ID, time.Now()).Scan(&pc2))

		community := &Community{
			Id:          communityID,
			Url:         "https://testcommunity.community.retrovibe.space",
			Description: "A test community for RSS feed generation",
			Entropy:     "test-entropy-value",
			Mimetype:    "video/*",
		}

		items, err := buildFeedItems(ctx, q, community)
		require.NoError(t, err)
		require.Len(t, items, 2)

		buf := new(bytes.Buffer)
		channel := rss.Channel{
			Title:         "testcommunity",
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
		)
		defer done()

		var lmd1 library.Metadata
		require.NoError(t, testx.Fake(&lmd1, library.MetadataOptionTestDefaults, func(m *library.Metadata) {
			m.ID = uuid.Must(uuid.NewV7()).String()
			m.Description = "unlisted-media.mp4"
			m.Bytes = bytesx.KiB
			m.Mimetype = "video/mp4"
		}))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd1).Scan(&lmd1))

		var lmd2 library.Metadata
		require.NoError(t, testx.Fake(&lmd2, library.MetadataOptionTestDefaults, func(m *library.Metadata) {
			m.ID = uuid.Must(uuid.NewV7()).String()
			m.Description = "hosted-media.mp4"
			m.Bytes = 2 * bytesx.KiB
			m.Mimetype = "video/mp4"
		}))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd2).Scan(&lmd2))

		var pc1 community.PublishedContent
		require.NoError(t, testx.Fake(&pc1, community.PublishedContentOptionTestDefaults, func(p *community.PublishedContent) {
			p.CommunityID = communityID
			p.MagnetURI = "magnet:?xt=urn:btih:0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33"
			p.LibraryID = lmd1.ID
			p.PublishMode = int32(PublishMode_UNLISTED)
		}))
		require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, pc1).Scan(&pc1))
		require.NoError(t, community.PublishedContentUpdatePublishedAt(ctx, q, pc1.ID, time.Now()).Scan(&pc1))

		var pc2 community.PublishedContent
		require.NoError(t, testx.Fake(&pc2, community.PublishedContentOptionTestDefaults, func(p *community.PublishedContent) {
			p.CommunityID = communityID
			p.MagnetURI = "magnet:?xt=urn:btih:1234567890abcdef1234567890abcdef12345678"
			p.LibraryID = lmd2.ID
			p.PublishMode = int32(PublishMode_LISTED)
		}))
		require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, pc2).Scan(&pc2))
		require.NoError(t, community.PublishedContentUpdatePublishedAt(ctx, q, pc2.ID, time.Now()).Scan(&pc2))

		community := &Community{
			Id:          communityID,
			Url:         "https://testcommunity.community.retrovibe.space",
			Description: "test",
			Entropy:     "test-entropy",
			Mimetype:    "video/*",
		}

		items, err := buildFeedItems(ctx, q, community)
		require.NoError(t, err)
		require.Len(t, items, 1)
		require.Equal(t, pc2.ID, items[0].Guid)
	})

	t.Run("excludes pending (not yet synced) content from feed", func(t *testing.T) {
		var (
			ctx, done   = testx.Context(t)
			q           = sqltestx.Metadatabase(t)
			communityID = uuid.Must(uuid.NewV7()).String()
		)
		defer done()

		var lmd1 library.Metadata
		require.NoError(t, testx.Fake(&lmd1, library.MetadataOptionTestDefaults, func(m *library.Metadata) {
			m.ID = uuid.Must(uuid.NewV7()).String()
			m.Description = "synced-media.mp4"
			m.Bytes = bytesx.KiB
			m.Mimetype = "video/mp4"
		}))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd1).Scan(&lmd1))

		var lmd2 library.Metadata
		require.NoError(t, testx.Fake(&lmd2, library.MetadataOptionTestDefaults, func(m *library.Metadata) {
			m.ID = uuid.Must(uuid.NewV7()).String()
			m.Description = "pending-media.mp4"
			m.Bytes = 2 * bytesx.KiB
			m.Mimetype = "video/mp4"
		}))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd2).Scan(&lmd2))

		var pc1 community.PublishedContent
		require.NoError(t, testx.Fake(&pc1, community.PublishedContentOptionTestDefaults, func(p *community.PublishedContent) {
			p.CommunityID = communityID
			p.MagnetURI = "magnet:?xt=urn:btih:0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33"
			p.LibraryID = lmd1.ID
			p.PublishMode = int32(PublishMode_LISTED)
		}))
		require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, pc1).Scan(&pc1))
		require.NoError(t, community.PublishedContentUpdatePublishedAt(ctx, q, pc1.ID, time.Now()).Scan(&pc1))

		// still mid-sync: no magnet_uri yet and published_at left at its
		// default 'infinity', simulating a sibling row caught by a feed
		// regeneration that fires while it's still being processed.
		var pc2 community.PublishedContent
		require.NoError(t, testx.Fake(&pc2, community.PublishedContentOptionTestDefaults, func(p *community.PublishedContent) {
			p.CommunityID = communityID
			p.MagnetURI = ""
			p.LibraryID = lmd2.ID
			p.PublishMode = int32(PublishMode_LISTED)
		}))
		require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, pc2).Scan(&pc2))

		community := &Community{
			Id:          communityID,
			Url:         "https://testcommunity.community.retrovibe.space",
			Description: "test",
			Entropy:     "test-entropy",
			Mimetype:    "video/*",
		}

		items, err := buildFeedItems(ctx, q, community)
		require.NoError(t, err)
		require.Len(t, items, 1)
		require.Equal(t, pc1.ID, items[0].Guid)
	})

	t.Run("excludes published content missing a magnet uri", func(t *testing.T) {
		var (
			ctx, done   = testx.Context(t)
			q           = sqltestx.Metadatabase(t)
			communityID = uuid.Must(uuid.NewV7()).String()
		)
		defer done()

		var lmd1 library.Metadata
		require.NoError(t, testx.Fake(&lmd1, library.MetadataOptionTestDefaults, func(m *library.Metadata) {
			m.ID = uuid.Must(uuid.NewV7()).String()
			m.Description = "synced-media.mp4"
			m.Bytes = bytesx.KiB
			m.Mimetype = "video/mp4"
		}))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd1).Scan(&lmd1))

		var lmd2 library.Metadata
		require.NoError(t, testx.Fake(&lmd2, library.MetadataOptionTestDefaults, func(m *library.Metadata) {
			m.ID = uuid.Must(uuid.NewV7()).String()
			m.Description = "urlless-media.mp4"
			m.Bytes = 2 * bytesx.KiB
			m.Mimetype = "video/mp4"
		}))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd2).Scan(&lmd2))

		var pc1 community.PublishedContent
		require.NoError(t, testx.Fake(&pc1, community.PublishedContentOptionTestDefaults, func(p *community.PublishedContent) {
			p.CommunityID = communityID
			p.MagnetURI = "magnet:?xt=urn:btih:0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33"
			p.LibraryID = lmd1.ID
			p.PublishMode = int32(PublishMode_LISTED)
		}))
		require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, pc1).Scan(&pc1))
		require.NoError(t, community.PublishedContentUpdatePublishedAt(ctx, q, pc1.ID, time.Now()).Scan(&pc1))

		// published_at has already been stamped (e.g. by a stale sync) but the
		// magnet uri never landed, so the item has no usable link/enclosure.
		var pc2 community.PublishedContent
		require.NoError(t, testx.Fake(&pc2, community.PublishedContentOptionTestDefaults, func(p *community.PublishedContent) {
			p.CommunityID = communityID
			p.MagnetURI = ""
			p.LibraryID = lmd2.ID
			p.PublishMode = int32(PublishMode_LISTED)
		}))
		require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, pc2).Scan(&pc2))
		require.NoError(t, community.PublishedContentUpdatePublishedAt(ctx, q, pc2.ID, time.Now()).Scan(&pc2))

		community := &Community{
			Id:          communityID,
			Url:         "https://testcommunity.community.retrovibe.space",
			Description: "test",
			Entropy:     "test-entropy",
			Mimetype:    "video/*",
		}

		items, err := buildFeedItems(ctx, q, community)
		require.NoError(t, err)
		require.Len(t, items, 1)
		require.Equal(t, pc1.ID, items[0].Guid)
	})
}
