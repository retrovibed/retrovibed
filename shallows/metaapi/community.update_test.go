package metaapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestCommunityUpdate(t *testing.T) {
	t.Run("updates community successfully", func(t *testing.T) {
		var expected meta.CommunityUpdateRequest
		require.NoError(t, testx.Fake(&expected))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPut, r.Method)
			require.Contains(t, r.URL.Path, "/c/")
			require.NoError(t, json.NewEncoder(w).Encode(&meta.CommunityUpdateResponse{Community: expected.Community}))
		}))

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)
		resp, err := metaapi.CommunityUpdate(t.Context(), c, "test-community", &expected)
		require.NoError(t, err)

		require.Equal(t, expected.Community.Id, resp.Community.Id)
		require.Equal(t, expected.Community.Domain, resp.Community.Domain)
		require.Equal(t, expected.Community.Description, resp.Community.Description)
		require.Equal(t, expected.Community.Mimetype, resp.Community.Mimetype)
	})

	t.Run("updates hidden field", func(t *testing.T) {
		var expected meta.CommunityUpdateRequest
		require.NoError(t, testx.Fake(&expected))
		expected.Community.Hidden = true

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var received meta.CommunityUpdateRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&received))
			require.True(t, received.Community.Hidden)
			require.NoError(t, json.NewEncoder(w).Encode(&meta.CommunityUpdateResponse{Community: expected.Community}))
		}))

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)
		resp, err := metaapi.CommunityUpdate(t.Context(), c, "test-community", &expected)
		require.NoError(t, err)
		require.True(t, resp.Community.Hidden)
	})

	t.Run("includes domain in path", func(t *testing.T) {
		var expected meta.CommunityUpdateRequest
		require.NoError(t, testx.Fake(&expected))

		domainOrId := "my-test-domain"

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/c/"+domainOrId, r.URL.Path)
			require.NoError(t, json.NewEncoder(w).Encode(&meta.CommunityUpdateResponse{Community: expected.Community}))
		}))

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)
		_, err := metaapi.CommunityUpdate(t.Context(), c, domainOrId, &expected)
		require.NoError(t, err)
	})

	t.Run("returns error on server failure", func(t *testing.T) {
		var expected meta.CommunityUpdateRequest
		require.NoError(t, testx.Fake(&expected))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)
		_, err := metaapi.CommunityUpdate(t.Context(), c, "test-community", &expected)
		require.Error(t, err)
	})
}
