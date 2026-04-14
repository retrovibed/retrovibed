package deeppool_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/deeppool"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/stretchr/testify/require"
)

func TestYouTubeExchange(t *testing.T) {
	t.Run("exchanges code for tokens", func(t *testing.T) {
		expected := deeppool.GoogleTokenResponse{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			TokenType:    "Bearer",
			ExpiresIn:    3600,
			Scope:        "https://www.googleapis.com/auth/youtube.upload",
		}

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/oauth2/proxy/google/token", r.URL.Path)

			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			form, err := url.ParseQuery(string(body))
			require.NoError(t, err)
			require.Equal(t, "authorization_code", form.Get("grant_type"))
			require.Equal(t, "test-code", form.Get("code"))
			require.Equal(t, "https://example.com/callback", form.Get("redirect_uri"))

			require.NoError(t, json.NewEncoder(w).Encode(&expected))
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)
		tok, err := deeppool.NewYouTube(c).Exchange(t.Context(), "test-code", "https://example.com/callback")
		require.NoError(t, err)
		require.Equal(t, expected.AccessToken, tok.AccessToken)
		require.Equal(t, expected.RefreshToken, tok.RefreshToken)
		require.Equal(t, expected.ExpiresIn, tok.ExpiresIn)
	})

	t.Run("returns error on server failure", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)
		_, err := deeppool.NewYouTube(c).Exchange(t.Context(), "test-code", "https://example.com/callback")
		require.Error(t, err)
	})
}

func TestYouTubeRefresh(t *testing.T) {
	t.Run("refreshes tokens", func(t *testing.T) {
		expected := deeppool.GoogleTokenResponse{
			AccessToken: "new-access-token",
			TokenType:   "Bearer",
			ExpiresIn:   3600,
		}

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/oauth2/proxy/google/token", r.URL.Path)

			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			form, err := url.ParseQuery(string(body))
			require.NoError(t, err)
			require.Equal(t, "refresh_token", form.Get("grant_type"))
			require.Equal(t, "test-refresh-token", form.Get("refresh_token"))

			require.NoError(t, json.NewEncoder(w).Encode(&expected))
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)
		tok, err := deeppool.NewYouTube(c).Refresh(t.Context(), "test-refresh-token")
		require.NoError(t, err)
		require.Equal(t, expected.AccessToken, tok.AccessToken)
		require.Equal(t, expected.ExpiresIn, tok.ExpiresIn)
	})

	t.Run("returns error on server failure", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)
		_, err := deeppool.NewYouTube(c).Refresh(t.Context(), "test-refresh-token")
		require.Error(t, err)
	})
}
