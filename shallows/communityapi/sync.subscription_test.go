package communityapi

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
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
