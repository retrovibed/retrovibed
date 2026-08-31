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

func TestCommunityUpdate(t *testing.T) {
	t.Run("updates community successfully", func(t *testing.T) {
		ctx, cancel := testx.Context(t)
		defer cancel()

		current := &communityapi.Community{
			Id:          "test-id",
			Url:         "https://test-community.community.retrovibe.space",
			Description: "original description",
			Mimetype:    "video/mp4",
		}
		updated := &communityapi.Community{
			Id:          "test-id",
			Url:         "https://test-community.community.retrovibe.space",
			Description: "updated description",
			Mimetype:    "video/mp4",
		}

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/c/test-community", r.URL.Path)
			switch r.Method {
			case http.MethodGet:
				assert.NoError(t, json.NewEncoder(w).Encode(&communityapi.CommunityFindResponse{Community: current}))
			case http.MethodPut:
				assert.NoError(t, json.NewEncoder(w).Encode(&communityapi.CommunityUpdateResponse{Community: updated}))
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		cmd := cmdCommunityUpdate{
			Name:        "test-community",
			Description: new("updated description"),
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

	t.Run("returns error when info fails", func(t *testing.T) {
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
			Description: new("updated description"),
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

	t.Run("returns error when update fails", func(t *testing.T) {
		ctx, cancel := testx.Context(t)
		defer cancel()

		current := &communityapi.Community{Id: "test-id", Url: "https://test-community.community.retrovibe.space"}

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				assert.NoError(t, json.NewEncoder(w).Encode(&communityapi.CommunityFindResponse{Community: current}))
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		cmd := cmdCommunityUpdate{
			Name:        "test-community",
			Description: new("updated description"),
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

	t.Run("preserves current values for unset fields", func(t *testing.T) {
		ctx, cancel := testx.Context(t)
		defer cancel()

		current := &communityapi.Community{
			Id:                 "test-id",
			Url:                "https://test-community.community.retrovibe.space",
			Description:        "original description",
			Mimetype:           "video/mp4",
			DefaultPublishMode: communityapi.PublishMode_LISTED,
			DefaultTtl:         3600,
			DefaultLanguage:    "en",
			Hidden:             true,
			Adult:              true,
		}

		var receivedReq communityapi.CommunityUpdateRequest

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				assert.NoError(t, json.NewEncoder(w).Encode(&communityapi.CommunityFindResponse{Community: current}))
			case http.MethodPut:
				assert.NoError(t, json.NewDecoder(r.Body).Decode(&receivedReq))
				assert.NoError(t, json.NewEncoder(w).Encode(&communityapi.CommunityUpdateResponse{Community: &communityapi.Community{}}))
			}
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		// only override description; all other fields should come from current
		cmd := cmdCommunityUpdate{
			Name:        "test-community",
			Description: new("new description"),
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
		require.Equal(t, "video/mp4", receivedReq.Community.Mimetype)
		require.Equal(t, communityapi.PublishMode_LISTED, receivedReq.Community.DefaultPublishMode)
		require.Equal(t, uint64(3600), receivedReq.Community.DefaultTtl)
		require.Equal(t, "en", receivedReq.Community.DefaultLanguage)
		require.True(t, receivedReq.Community.Hidden)
		require.True(t, receivedReq.Community.Adult)
	})

	t.Run("overrides default publish mode", func(t *testing.T) {
		ctx, cancel := testx.Context(t)
		defer cancel()

		current := &communityapi.Community{
			Id:                 "test-id",
			Url:                "https://test-community.community.retrovibe.space",
			DefaultPublishMode: communityapi.PublishMode_UNLISTED,
		}

		var receivedReq communityapi.CommunityUpdateRequest

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				assert.NoError(t, json.NewEncoder(w).Encode(&communityapi.CommunityFindResponse{Community: current}))
			case http.MethodPut:
				assert.NoError(t, json.NewDecoder(r.Body).Decode(&receivedReq))
				assert.NoError(t, json.NewEncoder(w).Encode(&communityapi.CommunityUpdateResponse{Community: &communityapi.Community{}}))
			}
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		cmd := cmdCommunityUpdate{
			Name:               "test-community",
			DefaultPublishMode: new("SYNDICATED"),
		}
		gctx := &cmdopts.Global{
			Context:  ctx,
			Shutdown: cancel,
			Cleanup:  &sync.WaitGroup{},
		}
		dpc := cmdopts.DeeppoolClientTest{Client: c}

		err := cmd.Run(gctx, dpc)
		require.NoError(t, err)
		require.Equal(t, communityapi.PublishMode_SYNDICATED, receivedReq.Community.DefaultPublishMode)
	})

	t.Run("includes community name in path", func(t *testing.T) {
		ctx, cancel := testx.Context(t)
		defer cancel()

		communityName := "my-special-community"

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/c/"+communityName, r.URL.Path)
			switch r.Method {
			case http.MethodGet:
				assert.NoError(t, json.NewEncoder(w).Encode(&communityapi.CommunityFindResponse{Community: &communityapi.Community{}}))
			case http.MethodPut:
				assert.NoError(t, json.NewEncoder(w).Encode(&communityapi.CommunityUpdateResponse{Community: &communityapi.Community{}}))
			}
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		cmd := cmdCommunityUpdate{Name: communityName}
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
