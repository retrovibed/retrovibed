package communityapi

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/bytesx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
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
