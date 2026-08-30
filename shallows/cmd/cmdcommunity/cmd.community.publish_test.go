package cmdcommunity

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"strings"
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

	t.Run("uploaded feed uses domain as title and url as link when url is set", func(t *testing.T) {
		ctx, cancel := testx.Context(t)
		defer cancel()

		communityInfo := &communityapi.Community{
			Id:          "test-id",
			Domain:      "mysite",
			Url:         "https://mysite.example.com",
			Description: "test description",
			Mimetype:    "video/mp4",
		}

		var uploadedFeed []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				assert.NoError(t, json.NewEncoder(w).Encode(&communityapi.CommunityFindResponse{Community: communityInfo}))
				return
			}
			if r.Method == http.MethodPost {
				reader, err := r.MultipartReader()
				assert.NoError(t, err)
				part, err := reader.NextPart()
				assert.NoError(t, err)
				uploadedFeed, err = io.ReadAll(part)
				assert.NoError(t, err)
				assert.NoError(t, json.NewEncoder(w).Encode(&communityapi.CommunityUploadResponse{Community: communityInfo}))
				return
			}
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		cmd := cmdCommunityPublish{
			Name:   "mysite",
			DryRun: false,
		}
		gctx := &cmdopts.Global{
			Context:  ctx,
			Shutdown: cancel,
			Cleanup:  &sync.WaitGroup{},
		}
		dpc := cmdopts.DeeppoolClientTest{Client: c}

		require.NoError(t, cmd.Run(gctx, dpc))
		require.Contains(t, string(uploadedFeed), "<title>mysite</title>")
		require.Contains(t, string(uploadedFeed), "<link>https://mysite.example.com</link>")
	})

	t.Run("uploaded feed falls back to the canonical hosted url when url is blank", func(t *testing.T) {
		ctx, cancel := testx.Context(t)
		defer cancel()

		communityInfo := &communityapi.Community{
			Id:          "test-id",
			Domain:      "mysite",
			Description: "test description",
			Mimetype:    "video/mp4",
		}

		var uploadedFeed []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				assert.NoError(t, json.NewEncoder(w).Encode(&communityapi.CommunityFindResponse{Community: communityInfo}))
				return
			}
			if r.Method == http.MethodPost {
				reader, err := r.MultipartReader()
				assert.NoError(t, err)
				part, err := reader.NextPart()
				assert.NoError(t, err)
				uploadedFeed, err = io.ReadAll(part)
				assert.NoError(t, err)
				assert.NoError(t, json.NewEncoder(w).Encode(&communityapi.CommunityUploadResponse{Community: communityInfo}))
				return
			}
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		cmd := cmdCommunityPublish{
			Name:   "mysite",
			DryRun: false,
		}
		gctx := &cmdopts.Global{
			Context:  ctx,
			Shutdown: cancel,
			Cleanup:  &sync.WaitGroup{},
		}
		dpc := cmdopts.DeeppoolClientTest{Client: c}

		require.NoError(t, cmd.Run(gctx, dpc))
		require.Contains(t, string(uploadedFeed), "<title>mysite</title>")
		require.Contains(t, string(uploadedFeed), "<link>https://mysite.community.retrovibe.space</link>")
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

func TestCommunityPublishItemsLink(t *testing.T) {
	const infohash = `{"id":"0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33","description":"a torrent"}` + "\n"

	t.Run("uses the community url when set", func(t *testing.T) {
		cmd := cmdCommunityPublish{}
		c := &communityapi.Community{
			Domain: "mysite",
			Url:    "https://mysite.example.com",
		}

		items := slices.Collect(cmd.items(c, strings.NewReader(infohash)))
		require.Len(t, items, 1)
		require.Equal(t, "https://mysite.example.com/0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33", items[0].Link)
	})

	t.Run("falls back to the canonical hosted url when url is blank", func(t *testing.T) {
		cmd := cmdCommunityPublish{}
		c := &communityapi.Community{
			Domain: "mysite",
		}

		items := slices.Collect(cmd.items(c, strings.NewReader(infohash)))
		require.Len(t, items, 1)
		require.Equal(t, "https://mysite.community.retrovibe.space/0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33", items[0].Link)
	})
}
