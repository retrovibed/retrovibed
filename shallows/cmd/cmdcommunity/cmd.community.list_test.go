package cmdcommunity

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/communityapi"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommunityList(t *testing.T) {
	t.Run("returns error on community info failure", func(t *testing.T) {
		ctx, cancel := testx.Context(t)
		defer cancel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		cmd := cmdCommunityList{
			Name: "nonexistent-community",
		}
		gctx := &cmdopts.Global{
			Context:  ctx,
			Shutdown: cancel,
			Cleanup:  &sync.WaitGroup{},
		}
		dpc := cmdopts.DeeppoolClientTest{Client: c}

		err := cmd.Run(gctx, dpc)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to locate community")
	})

	t.Run("lists community with no published content", func(t *testing.T) {
		ctx, cancel := testx.Context(t)
		defer cancel()

		_ = sqltestx.Metadatabase(t)

		expected := &communityapi.Community{
			Id:          "00000000-0000-0000-0000-000000000001",
			Domain:      "test-community",
			Description: "test description",
		}

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			assert.NoError(t, json.NewEncoder(w).Encode(&communityapi.CommunityFindResponse{Community: expected}))
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		cmd := cmdCommunityList{
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
	})
}
