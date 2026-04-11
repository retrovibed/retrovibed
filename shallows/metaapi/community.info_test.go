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

func TestCommunityInfo(t *testing.T) {
	t.Run("example 1", func(t *testing.T) {
		var (
			expected meta.Community
		)

		require.NoError(t, testx.Fake(&expected))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.NoError(t, json.NewEncoder(w).Encode(&meta.CommunityFindResponse{Community: &expected}))
		}))

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)
		resp, err := metaapi.CommunityInfo(t.Context(), c, "derp")
		require.NoError(t, err)

		require.Equal(t, expected.Id, resp.Community.Id)
		require.Equal(t, expected.Domain, resp.Community.Domain)
		require.Equal(t, expected.Bytes, resp.Community.Bytes)
		require.Equal(t, expected.Description, resp.Community.Description)
		require.Equal(t, expected.Mimetype, resp.Community.Mimetype)
	})
}
