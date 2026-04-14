package metaapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestCommunityPublish(t *testing.T) {
	t.Run("publishes feed successfully", func(t *testing.T) {
		var expected meta.CommunityUploadResponse
		require.NoError(t, testx.Fake(&expected))

		communityID := "test-community"
		feedContent := `<?xml version="1.0" encoding="UTF-8"?><rss version="2.0"><channel><title>Test</title></channel></rss>`

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/c/"+communityID, r.URL.Path)
			require.Contains(t, r.Header.Get("Content-Type"), "multipart/form-data")
			require.NoError(t, json.NewEncoder(w).Encode(&expected))
		}))

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)
		_, err := metaapi.CommunityPublish(t.Context(), c, communityID, strings.NewReader(feedContent))
		require.NoError(t, err)
	})

	t.Run("returns error on server failure", func(t *testing.T) {
		feedContent := `<?xml version="1.0" encoding="UTF-8"?><rss version="2.0"><channel><title>Test</title></channel></rss>`

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)
		_, err := metaapi.CommunityPublish(t.Context(), c, "test-community", strings.NewReader(feedContent))
		require.Error(t, err)
	})
}
