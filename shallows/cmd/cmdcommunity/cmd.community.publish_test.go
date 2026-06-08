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

func TestCommunityPublish(t *testing.T) {
	t.Run("dry run does not upload", func(t *testing.T) {
		ctx, cancel := testx.Context(t)
		defer cancel()

		communityInfo := &communityapi.Community{
			Id:          "test-id",
			Domain:      "test-community",
			Description: "test description",
			Mimetype:    "video/mp4",
		}

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				assert.NoError(t, json.NewEncoder(w).Encode(&communityapi.CommunityFindResponse{Community: communityInfo}))
				return
			}
			t.Fatal("should not upload in dry run mode")
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		cmd := cmdCommunityPublish{
			Name:   "test-community",
			DryRun: true,
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

	t.Run("uploads when dry run is false", func(t *testing.T) {
		ctx, cancel := testx.Context(t)
		defer cancel()

		communityInfo := &communityapi.Community{
			Id:          "test-id",
			Domain:      "test-community",
			Description: "test description",
			Mimetype:    "video/mp4",
		}

		uploaded := false
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				assert.NoError(t, json.NewEncoder(w).Encode(&communityapi.CommunityFindResponse{Community: communityInfo}))
				return
			}
			if r.Method == http.MethodPost {
				uploaded = true
				assert.NoError(t, json.NewEncoder(w).Encode(&communityapi.CommunityUploadResponse{Community: communityInfo}))
				return
			}
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		cmd := cmdCommunityPublish{
			Name:   "test-community",
			DryRun: false,
		}
		gctx := &cmdopts.Global{
			Context:  ctx,
			Shutdown: cancel,
			Cleanup:  &sync.WaitGroup{},
		}
		dpc := cmdopts.DeeppoolClientTest{Client: c}

		err := cmd.Run(gctx, dpc)
		require.NoError(t, err)
		require.True(t, uploaded, "should upload when dry run is false")
	})

	t.Run("returns error on community info failure", func(t *testing.T) {
		ctx, cancel := testx.Context(t)
		defer cancel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		cmd := cmdCommunityPublish{
			Name:   "nonexistent-community",
			DryRun: false,
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
}
