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

func TestPeerCreate(t *testing.T) {
	t.Run("creates a peer", func(t *testing.T) {
		var (
			peer     ddiscapi.Peer
			expected ddiscapi.PeerCreateResponse
		)
		require.NoError(t, testx.Fake(&peer))
		expected.Peer = &peer

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/ddisc/", r.URL.Path)
			require.NoError(t, json.NewEncoder(w).Encode(&expected))
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		resp, err := ddiscapi.PeerCreate(t.Context(), c, srv.URL, &ddiscapi.PeerCreateRequest{Peer: &peer})
		require.NoError(t, err)
		require.Equal(t, expected.Peer.Id, resp.Peer.Id)
	})

	t.Run("returns error on server failure", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		_, err := ddiscapi.PeerCreate(t.Context(), c, srv.URL, &ddiscapi.PeerCreateRequest{})
		require.Error(t, err)
	})
}
