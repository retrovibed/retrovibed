package cmdcommunity

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/communityapi"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommunityCreate(t *testing.T) {
	t.Run("creates community successfully", func(t *testing.T) {
		ctx, cancel := testx.Context(t)
		defer cancel()

		expected := &communityapi.Community{
			Id:          "test-id",
			Url:         "https://test-community.community.retrovibe.space",
			Description: "test description",
			Mimetype:    "application/octet-stream",
		}

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/c/", r.URL.Path)
			assert.NoError(t, json.NewEncoder(w).Encode(&communityapi.CommunityCreateResponse{Community: expected}))
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		cmd := cmdCommunityCreate{
			Name:        "test-community",
			Description: "test description",
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

		cmd := cmdCommunityCreate{
			Name:        "test-community",
			Description: "test description",
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

		var receivedReq communityapi.CommunityCreateRequest

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.NoError(t, json.NewDecoder(r.Body).Decode(&receivedReq))
			assert.NoError(t, json.NewEncoder(w).Encode(&communityapi.CommunityCreateResponse{
				Community: &communityapi.Community{
					Id:          "created-id",
					Url:         receivedReq.Community.Url,
					Description: receivedReq.Community.Description,
					Mimetype:    receivedReq.Community.Mimetype,
				},
			}))
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		cmd := cmdCommunityCreate{
			Name:        "my-community",
			Description: "my description",
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

		require.Equal(t, "https://my-community.community.retrovibe.space", receivedReq.Community.Url)
		require.Equal(t, "my description", receivedReq.Community.Description)
		require.Equal(t, "video/mp4", receivedReq.Community.Mimetype)
	})

	t.Run("uses default mimetype when not specified", func(t *testing.T) {
		ctx, cancel := testx.Context(t)
		defer cancel()

		var receivedReq communityapi.CommunityCreateRequest

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.NoError(t, json.NewDecoder(r.Body).Decode(&receivedReq))
			assert.NoError(t, json.NewEncoder(w).Encode(&communityapi.CommunityCreateResponse{
				Community: &communityapi.Community{},
			}))
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		cmd := cmdCommunityCreate{
			Name: "test-community",
		}
		gctx := &cmdopts.Global{
			Context:  ctx,
			Shutdown: cancel,
			Cleanup:  &sync.WaitGroup{},
		}
		dpc := cmdopts.DeeppoolClientTest{Client: c}

		err := cmd.Run(gctx, dpc)
		require.NoError(t, err)

		require.Equal(t, "application/octet-stream", receivedReq.Community.Mimetype)
	})
}
