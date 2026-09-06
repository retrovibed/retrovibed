package communityapi

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/retroapi/bytesx"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/internal/torrentx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestSyncPublishedContentItem(t *testing.T) {
	t.Run("propagates description to the stored published content", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		pc := &PublishedContent{
			Id:          uuid.Must(uuid.NewV7()).String(),
			CommunityId: uuid.Must(uuid.NewV7()).String(),
			MagnetUri:   "magnet:?xt=urn:btih:0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33",
			Title:       "Test Title",
			Description: "a detailed description of the content",
		}

		require.NoError(t, SyncPublishedContentItem(ctx, q, pc, false))

		stored, err := sqlx.ScanOne(community.PublishedContentSearch(ctx, q, community.PublishedContentSearchBuilder().Where(community.PublishedContentQueryCommunityID(pc.CommunityId))))
		require.NoError(t, err)
		require.Equal(t, pc.Description, stored.Description)
	})

	t.Run("round-trips proto through the database unchanged", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		pc := &PublishedContent{
			Id:            uuid.Must(uuid.NewV7()).String(),
			CommunityId:   uuid.Must(uuid.NewV7()).String(),
			MagnetUri:     "magnet:?xt=urn:btih:2beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a35",
			Title:         "Test Title",
			Description:   "a detailed description of the content",
			OauthGoogleId: uuid.Must(uuid.NewV7()).String(),
			Bytes:         uint64(3 * bytesx.MiB),
			Mimetype:      "video/mp4",
		}

		require.NoError(t, SyncPublishedContentItem(ctx, q, pc, false))

		stored, err := sqlx.ScanOne(community.PublishedContentSearch(ctx, q, community.PublishedContentSearchBuilder().Where(community.PublishedContentQueryCommunityID(pc.CommunityId))))
		require.NoError(t, err)

		roundtripped := &PublishedContent{}
		PublishedContentOptionFromDB(stored)(roundtripped)

		require.Equal(t, pc.CommunityId, roundtripped.CommunityId)
		require.Equal(t, pc.MagnetUri, roundtripped.MagnetUri)
		require.Equal(t, pc.Title, roundtripped.Title)
		require.Equal(t, pc.Description, roundtripped.Description)
		require.Equal(t, pc.OauthGoogleId, roundtripped.OauthGoogleId)
		require.Equal(t, pc.Bytes, roundtripped.Bytes)
		require.Equal(t, pc.Mimetype, roundtripped.Mimetype)
	})

	t.Run("stores items separately when multiple lack a local library association", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		communityID := uuid.Must(uuid.NewV7()).String()

		pc1 := &PublishedContent{
			Id:          uuid.Must(uuid.NewV7()).String(),
			CommunityId: communityID,
			MagnetUri:   "magnet:?xt=urn:btih:3beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a36",
			Title:       "Test Title 1",
		}
		pc2 := &PublishedContent{
			Id:          uuid.Must(uuid.NewV7()).String(),
			CommunityId: communityID,
			MagnetUri:   "magnet:?xt=urn:btih:4beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a37",
			Title:       "Test Title 2",
		}

		require.NoError(t, SyncPublishedContentItem(ctx, q, pc1, false))
		require.NoError(t, SyncPublishedContentItem(ctx, q, pc2, false))

		var stored []community.PublishedContent
		require.NoError(t, sqlx.ScanInto(community.PublishedContentSearch(ctx, q, community.PublishedContentSearchBuilder().Where(community.PublishedContentQueryCommunityID(communityID))), &stored))
		require.Len(t, stored, 2, "items with distinct magnet URIs should be stored separately even without a local library association")
	})

	t.Run("carries the published mimetype onto the torrent row", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		pc := &PublishedContent{
			Id:          uuid.Must(uuid.NewV7()).String(),
			CommunityId: uuid.Must(uuid.NewV7()).String(),
			MagnetUri:   "magnet:?xt=urn:btih:5beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a38&dn=lemmy.wasm",
			Title:       "lemmy.wasm",
			Mimetype:    mimex.RetrovibedPublishModule,
		}

		require.NoError(t, SyncPublishedContentItem(ctx, q, pc, false))

		// resolved the same way SyncPublishedContentItem derives the row's identity.
		magnet, err := metainfo.ParseMagnetURI(pc.MagnetUri)
		require.NoError(t, err)

		var tmd tracking.Metadata
		require.NoError(t, tracking.MetadataFindByID(ctx, q, torrentx.HashUID(new(int160.FromByteArray(magnet.InfoHash)))).Scan(&tmd))

		require.Equal(t, mimex.RetrovibedPublishModule, tmd.Mimetype)
		// a publish module is machinery, not media; MetadataOptionAutoHidden keeps it
		// out of the media listings, but only when it can see the real mimetype.
		require.NotEqual(t, timex.Inf(), tmd.HiddenAt)
	})

	t.Run("is visible to the publish module import selector", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		pc := &PublishedContent{
			Id:          uuid.Must(uuid.NewV7()).String(),
			CommunityId: uuid.Must(uuid.NewV7()).String(),
			MagnetUri:   "magnet:?xt=urn:btih:6beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a39&dn=lemmy.wasm",
			Title:       "lemmy.wasm",
			Mimetype:    mimex.RetrovibedPublishModule,
		}

		require.NoError(t, SyncPublishedContentItem(ctx, q, pc, false))

		// the predicate PublishPluginTorrentImport selects on; a synced module that
		// misses it never reaches publish.d.
		var stored []tracking.Metadata
		require.NoError(t, sqlx.ScanInto(tracking.MetadataSearch(ctx, q, tracking.MetadataSearchBuilder().Where(tracking.MetadataQueryPublishModule())), &stored))
		require.Len(t, stored, 1)
	})

	t.Run("repairs the mimetype of a row recorded before the mimetype was carried", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		pc := &PublishedContent{
			Id:          uuid.Must(uuid.NewV7()).String(),
			CommunityId: uuid.Must(uuid.NewV7()).String(),
			MagnetUri:   "magnet:?xt=urn:btih:7beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a40&dn=noop.wasm",
			Title:       "noop.wasm",
			Mimetype:    mimex.RetrovibedPublishModule,
		}

		magnet, err := metainfo.ParseMagnetURI(pc.MagnetUri)
		require.NoError(t, err)

		// the row an earlier sync left behind, before the published mimetype was
		// carried onto it. resyncing is the only repair path a user has, so the
		// insert has to correct the existing row rather than leave it as media.
		stale := tracking.NewMetadata(new(int160.FromByteArray(magnet.InfoHash)), tracking.MetadataOptionFromMagnet(&magnet))
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, stale).Scan(&stale))
		require.Equal(t, mimex.Bittorrent, stale.Mimetype)
		require.Equal(t, timex.Inf(), stale.HiddenAt)

		require.NoError(t, SyncPublishedContentItem(ctx, q, pc, false))

		var updated tracking.Metadata
		require.NoError(t, tracking.MetadataFindByID(ctx, q, stale.ID).Scan(&updated))
		require.Equal(t, mimex.RetrovibedPublishModule, updated.Mimetype)
		// a stale row is unhidden too, so the repair has to take the module out of
		// the media listings, not just make it installable.
		require.NotEqual(t, timex.Inf(), updated.HiddenAt)
	})
}

func TestSyncSubscriptions(t *testing.T) {
	t.Run("iterates community table without error", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		sub := community.Community{
			ID:                         uuid.Must(uuid.NewV7()).String(),
			AccountID:                  uuid.Nil.String(),
			SyncCursorPublishedContent: uuid.Nil.String(),
			LastSyncAt:                 time.Now(),
		}
		require.NoError(t, community.CommunityInsertWithDefaults(ctx, q, sub).Scan(&sub))

		client := NewDeeppoolPublished(httptestx.NewTestClient(func(req *http.Request) *http.Response {
			body, _ := json.Marshal(&PublishedContentSearchResponse{})
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(string(body))),
				Header:     make(http.Header),
			}
		}))

		require.NoError(t, syncSubscriptions(ctx, q, client))
	})

	t.Run("updates last_sync_at after syncing", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		before := time.Now().Add(-time.Hour)
		sub := community.Community{
			ID:                         uuid.Must(uuid.NewV7()).String(),
			AccountID:                  uuid.Nil.String(),
			SyncCursorPublishedContent: uuid.Nil.String(),
			LastSyncAt:                 before,
		}
		require.NoError(t, community.CommunityInsertWithDefaults(ctx, q, sub).Scan(&sub))

		client := NewDeeppoolPublished(httptestx.NewTestClient(func(req *http.Request) *http.Response {
			body, _ := json.Marshal(&PublishedContentSearchResponse{})
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(string(body))),
				Header:     make(http.Header),
			}
		}))

		require.NoError(t, syncSubscriptions(ctx, q, client))

		var updated community.Community
		require.NoError(t, community.CommunityFindByID(ctx, q, sub.ID).Scan(&updated))
		require.True(t, updated.LastSyncAt.After(before))
	})
}
