package deeppool_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/deeppool"
	"github.com/retrovibed/retrovibed/retroapi/internal/httpx"
	"github.com/retrovibed/retrovibed/retroapi/internal/testx"
	"github.com/stretchr/testify/require"
)

func TestMediaIDClean(t *testing.T) {
	t.Run("sends text and returns cleaned result", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/clean/", r.URL.Path)

			var req deeppool.CleanIDRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
			require.Equal(t, "Barton Fink (1991) DVDRip XViD", req.Text)

			require.NoError(t, json.NewEncoder(w).Encode(&deeppool.CleanIDResponse{Text: "Barton Fink 1991"}))
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)
		result, err := deeppool.NewMediaID(c).Clean(t.Context(), "Barton Fink (1991) DVDRip XViD")
		require.NoError(t, err)
		require.Equal(t, "Barton Fink 1991", result)
	})

	t.Run("returns error on server failure", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)
		_, err := deeppool.NewMediaID(c).Clean(t.Context(), "some text")
		require.Error(t, err)
	})
}
