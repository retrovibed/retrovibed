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
	"github.com/stretchr/testify/require"
)

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
