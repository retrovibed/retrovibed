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

func TestDiscoveryCreate(t *testing.T) {
	t.Run("creates a discovery entry", func(t *testing.T) {
		var (
			discovery ddiscapi.Discovery
			expected  ddiscapi.DiscoveryCreateResponse
		)
		require.NoError(t, testx.Fake(&discovery))
		expected.Discovery = &discovery

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/ddisc/discovery/", r.URL.Path)
			require.NoError(t, json.NewEncoder(w).Encode(&expected))
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		resp, err := ddiscapi.DiscoveryCreate(t.Context(), c, "example.com", &ddiscapi.DiscoveryCreateRequest{Discovery: &discovery})
		require.NoError(t, err)
		require.Equal(t, expected.Discovery.Id, resp.Discovery.Id)
	})

	t.Run("returns error on server failure", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		_, err := ddiscapi.DiscoveryCreate(t.Context(), c, "example.com", &ddiscapi.DiscoveryCreateRequest{})
		require.Error(t, err)
	})
}
