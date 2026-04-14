package cmdcommunity

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommunityAdd(t *testing.T) {
	t.Run("adds content to community successfully", func(t *testing.T) {
		ctx, cancel := testx.Context(t)
		defer cancel()

		expected := &meta.PublishedContent{
			Id:           "published-id",
			CommunityId:  "test-community",
			MagnetUri:    "magnet:?xt=urn:btih:abc123",
			KnownMediaId: "known-media-123",
		}

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Contains(t, r.URL.Path, "/c/test-community/publish")
			assert.NoError(t, json.NewEncoder(w).Encode(&meta.PublishContentResponse{PublishedContent: expected}))
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		cmd := cmdCommunityAdd{
			Community:    "test-community",
			MagnetURI:    "magnet:?xt=urn:btih:abc123",
			KnownMediaID: "known-media-123",
		}
		gctx := &cmdopts.Global{
			Context:  ctx,
			Shutdown: cancel,
			Cleanup:  &sync.WaitGroup{},
		}
		dpc := cmdopts.DeeppoolClientTest{Client: c}

		err := cmd.Run(gctx, dpc)
		require.NoError(t, err)
	})

	t.Run("returns error on server failure", func(t *testing.T) {
		ctx, cancel := testx.Context(t)
		defer cancel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		cmd := cmdCommunityAdd{
			Community: "test-community",
			MagnetURI: "magnet:?xt=urn:btih:abc123",
		}
		gctx := &cmdopts.Global{
			Context:  ctx,
			Shutdown: cancel,
			Cleanup:  &sync.WaitGroup{},
		}
		dpc := cmdopts.DeeppoolClientTest{Client: c}

		err := cmd.Run(gctx, dpc)
		require.Error(t, err)
	})

	t.Run("sends correct request body", func(t *testing.T) {
		ctx, cancel := testx.Context(t)
		defer cancel()

		var receivedReq meta.PublishContentRequest

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.NoError(t, json.NewDecoder(r.Body).Decode(&receivedReq))
			assert.NoError(t, json.NewEncoder(w).Encode(&meta.PublishContentResponse{
				PublishedContent: &meta.PublishedContent{
					Id: "created-id",
				},
			}))
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		cmd := cmdCommunityAdd{
			Community:    "my-community",
			MagnetURI:    "magnet:?xt=urn:btih:xyz789",
			KnownMediaID: "tmdb-12345",
			ArchivedID:   "archive-abc",
		}
		gctx := &cmdopts.Global{
			Context:  ctx,
			Shutdown: cancel,
			Cleanup:  &sync.WaitGroup{},
		}
		dpc := cmdopts.DeeppoolClientTest{Client: c}

		err := cmd.Run(gctx, dpc)
		require.NoError(t, err)

		require.Equal(t, "magnet:?xt=urn:btih:xyz789", receivedReq.PublishedContent.MagnetUri)
		require.Equal(t, "tmdb-12345", receivedReq.PublishedContent.KnownMediaId)
		require.Equal(t, "archive-abc", receivedReq.PublishedContent.ArchivedId)
	})
}
