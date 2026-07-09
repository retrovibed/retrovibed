package ddiscapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestPeerDelete(t *testing.T) {
	t.Run("deletes peer successfully", func(t *testing.T) {
		var expected ddiscapi.PeerDeleteResponse
		require.NoError(t, testx.Fake(&expected))

		infohash := int160.Random()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodDelete, r.Method)
			require.Equal(t, "/ddisc/"+tracking.PeerUID(infohash), r.URL.Path)
			require.NoError(t, json.NewEncoder(w).Encode(&expected))
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		_, err := ddiscapi.PeerDelete(t.Context(), c, "example.com", infohash.String())
		require.NoError(t, err)
	})

	t.Run("returns error on invalid peer id", func(t *testing.T) {
		c := &http.Client{}

		_, err := ddiscapi.PeerDelete(t.Context(), c, "example.com", "not-hex")
		require.Error(t, err)
	})

	t.Run("returns error on server failure", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		_, err := ddiscapi.PeerDelete(t.Context(), c, "example.com", int160.Random().String())
		require.Error(t, err)
	})
}
