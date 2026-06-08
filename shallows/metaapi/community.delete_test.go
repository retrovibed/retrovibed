package metaapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/communityapi"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestCommunityDelete(t *testing.T) {
	t.Run("deletes community successfully", func(t *testing.T) {
		var expected communityapi.CommunityDeleteResponse
		require.NoError(t, testx.Fake(&expected))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodDelete, r.Method)
			require.Contains(t, r.URL.Path, "/c/")
			require.NoError(t, json.NewEncoder(w).Encode(&expected))
		}))

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)
		_, err := metaapi.CommunityDelete(t.Context(), c, "test-community")
		require.NoError(t, err)
	})

	t.Run("includes domain in path", func(t *testing.T) {
		var expected communityapi.CommunityDeleteResponse
		require.NoError(t, testx.Fake(&expected))

		domainOrId := "my-test-domain"

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/c/"+domainOrId, r.URL.Path)
			require.NoError(t, json.NewEncoder(w).Encode(&expected))
		}))

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)
		_, err := metaapi.CommunityDelete(t.Context(), c, domainOrId)
		require.NoError(t, err)
	})

	t.Run("returns error on server failure", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)
		_, err := metaapi.CommunityDelete(t.Context(), c, "test-community")
		require.Error(t, err)
	})
}
