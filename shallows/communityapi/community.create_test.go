package communityapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/communityapi"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/stretchr/testify/require"
)

func TestCommunityCreate(t *testing.T) {
	t.Run("example 1", func(t *testing.T) {
		var (
			expected communityapi.CommunityCreateRequest
		)

		require.NoError(t, testx.Fake(&expected))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.NoError(t, json.NewEncoder(w).Encode(&communityapi.CommunityCreateResponse{Community: expected.Community}))
		}))

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)
		resp, err := communityapi.CommunityCreate(t.Context(), c, &expected)
		require.NoError(t, err)

		require.Equal(t, expected.Community.Id, resp.Community.Id)
		require.Equal(t, expected.Community.Url, resp.Community.Url)
		require.Equal(t, expected.Community.Bytes, resp.Community.Bytes)
		require.Equal(t, expected.Community.Description, resp.Community.Description)
		require.Equal(t, expected.Community.Mimetype, resp.Community.Mimetype)
	})
}
