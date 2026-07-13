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

func TestLocateCreate(t *testing.T) {
	t.Run("submits a locate request", func(t *testing.T) {
		var (
			locate   ddiscapi.Locate
			expected ddiscapi.LocateCreateResponse
		)
		require.NoError(t, testx.Fake(&locate))
		expected.Locate = &locate

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/l/", r.URL.Path)
			require.NoError(t, json.NewEncoder(w).Encode(&expected))
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		resp, err := ddiscapi.LocateCreate(t.Context(), c, "example.com", &ddiscapi.LocateCreateRequest{Locate: &locate})
		require.NoError(t, err)
		require.Equal(t, expected.Locate.Id, resp.Locate.Id)
	})

	t.Run("returns error on server failure", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		_, err := ddiscapi.LocateCreate(t.Context(), c, "example.com", &ddiscapi.LocateCreateRequest{})
		require.Error(t, err)
	})
}
