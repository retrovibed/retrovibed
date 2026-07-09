package ddiscapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/stretchr/testify/require"
)

func TestMediaSearch(t *testing.T) {
	t.Run("returns decoded results", func(t *testing.T) {
		var expected ddiscapi.MediaSearchResponse
		require.NoError(t, testx.Fake(&expected))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/ddisc/media/", r.URL.Path)
			require.NoError(t, json.NewEncoder(w).Encode(&expected))
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		resp, err := ddiscapi.MediaSearch(t.Context(), c, "example.com", &ddiscapi.MediaSearchRequest{Query: "derp"})
		require.NoError(t, err)
		require.Equal(t, len(expected.Items), len(resp.Items))
	})

	t.Run("returns error on server failure", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		_, err := ddiscapi.MediaSearch(t.Context(), c, "example.com", &ddiscapi.MediaSearchRequest{})
		require.Error(t, err)
	})
}
