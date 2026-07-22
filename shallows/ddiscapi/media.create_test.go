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

func TestMediaCreate(t *testing.T) {
	t.Run("creates a media entry", func(t *testing.T) {
		var (
			media    ddiscapi.Media
			expected ddiscapi.MediaCreateResponse
		)
		require.NoError(t, testx.Fake(&media))
		expected.Media = &media

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/ddisc/media/", r.URL.Path)
			require.NoError(t, json.NewEncoder(w).Encode(&expected))
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		resp, err := ddiscapi.MediaCreate(t.Context(), c, srv.URL, &ddiscapi.MediaCreateRequest{Media: &media})
		require.NoError(t, err)
		require.Equal(t, expected.Media.Id, resp.Media.Id)
	})

	t.Run("returns error on server failure", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		_, err := ddiscapi.MediaCreate(t.Context(), c, srv.URL, &ddiscapi.MediaCreateRequest{})
		require.Error(t, err)
	})
}
