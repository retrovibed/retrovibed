package cmdcommunity

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"

	"github.com/retrovibed/retrovibed/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/internal/httpx"
	"github.com/retrovibed/retrovibed/internal/testx"
	"github.com/retrovibed/retrovibed/meta"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommunityUpdate(t *testing.T) {
	t.Run("updates community successfully", func(t *testing.T) {
		ctx, cancel := testx.Context(t)
		defer cancel()

		expected := &meta.Community{
			Id:          "test-id",
			Domain:      "test-community",
			Description: "updated description",
			Mimetype:    "video/mp4",
		}

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPut, r.Method)
			require.Equal(t, "/c/test-community", r.URL.Path)
			assert.NoError(t, json.NewEncoder(w).Encode(&meta.CommunityUpdateResponse{Community: expected}))
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		cmd := cmdCommunityUpdate{
			Name:        "test-community",
			Description: "updated description",
			Mimetype:    "video/mp4",
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

		cmd := cmdCommunityUpdate{
			Name:        "test-community",
			Description: "updated description",
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

		var receivedReq meta.CommunityUpdateRequest

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.NoError(t, json.NewDecoder(r.Body).Decode(&receivedReq))
			assert.NoError(t, json.NewEncoder(w).Encode(&meta.CommunityUpdateResponse{
				Community: &meta.Community{
					Id:          "updated-id",
					Description: receivedReq.Community.Description,
					Mimetype:    receivedReq.Community.Mimetype,
				},
			}))
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		cmd := cmdCommunityUpdate{
			Name:        "my-community",
			Description: "new description",
			Mimetype:    "audio/mpeg",
		}
		gctx := &cmdopts.Global{
			Context:  ctx,
			Shutdown: cancel,
			Cleanup:  &sync.WaitGroup{},
		}
		dpc := cmdopts.DeeppoolClientTest{Client: c}

		err := cmd.Run(gctx, dpc)
		require.NoError(t, err)

		require.Equal(t, "new description", receivedReq.Community.Description)
		require.Equal(t, "audio/mpeg", receivedReq.Community.Mimetype)
	})

	t.Run("includes community name in path", func(t *testing.T) {
		ctx, cancel := testx.Context(t)
		defer cancel()

		communityName := "my-special-community"

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/c/"+communityName, r.URL.Path)
			assert.NoError(t, json.NewEncoder(w).Encode(&meta.CommunityUpdateResponse{
				Community: &meta.Community{},
			}))
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		cmd := cmdCommunityUpdate{
			Name: communityName,
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
}
