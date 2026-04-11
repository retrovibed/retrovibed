package metaapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/retrovibed/retrovibed/internal/httpx"
	"github.com/retrovibed/retrovibed/internal/testx"
	"github.com/retrovibed/retrovibed/meta"
	"github.com/retrovibed/retrovibed/metaapi"
	"github.com/stretchr/testify/require"
)

func TestCommunityCreate(t *testing.T) {
	t.Run("example 1", func(t *testing.T) {
		var (
			expected meta.CommunityCreateRequest
		)

		require.NoError(t, testx.Fake(&expected))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.NoError(t, json.NewEncoder(w).Encode(&meta.CommunityCreateResponse{Community: expected.Community}))
		}))

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)
		resp, err := metaapi.CommunityCreate(t.Context(), c, &expected)
		require.NoError(t, err)

		require.Equal(t, expected.Community.Id, resp.Community.Id)
		require.Equal(t, expected.Community.Domain, resp.Community.Domain)
		require.Equal(t, expected.Community.Bytes, resp.Community.Bytes)
		require.Equal(t, expected.Community.Description, resp.Community.Description)
		require.Equal(t, expected.Community.Mimetype, resp.Community.Mimetype)
	})
}
