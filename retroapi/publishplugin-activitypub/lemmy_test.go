package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	t.Run("login exchanges credentials for a token and caches it", func(t *testing.T) {
		var received loginRequest

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/api/v3/user/login", r.URL.Path)
			require.NoError(t, json.NewDecoder(r.Body).Decode(&received))
			require.NoError(t, json.NewEncoder(w).Encode(loginResponse{JWT: "token.from.lemmy"}))
		}))
		defer srv.Close()

		cachedir := t.TempDir()
		client, err := NewClient(srv.URL, OptionHTTPClient(srv.Client()), OptionCacheDir(cachedir))
		require.NoError(t, err)
		require.False(t, client.Authenticated())

		require.NoError(t, client.Login(context.Background(), "publisher", "hunter2", "123456"))
		require.True(t, client.Authenticated())
		require.Equal(t, loginRequest{UsernameOrEmail: "publisher", Password: "hunter2", TOTP: "123456"}, received)

		cached, err := os.ReadFile(filepath.Join(cachedir, "session.json"))
		require.NoError(t, err)
		require.JSONEq(t, `{"jwt":"token.from.lemmy"}`, string(cached))
	})

	t.Run("login rejects a token-less success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.NoError(t, json.NewEncoder(w).Encode(loginResponse{}))
		}))
		defer srv.Close()

		client, err := NewClient(srv.URL, OptionHTTPClient(srv.Client()))
		require.NoError(t, err)

		require.Error(t, client.Login(context.Background(), "publisher", "hunter2", ""))
	})

	t.Run("restore session adopts a cached token", func(t *testing.T) {
		cachedir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(cachedir, "session.json"), []byte(`{"jwt":"cached.token"}`), 0600))

		client, err := NewClient("https://lemmy.invalid", OptionCacheDir(cachedir))
		require.NoError(t, err)

		require.True(t, client.RestoreSession())
		require.True(t, client.Authenticated())
	})

	t.Run("restore session tolerates a missing or corrupt cache", func(t *testing.T) {
		cachedir := t.TempDir()

		client, err := NewClient("https://lemmy.invalid", OptionCacheDir(cachedir))
		require.NoError(t, err)
		require.False(t, client.RestoreSession())

		require.NoError(t, os.WriteFile(filepath.Join(cachedir, "session.json"), []byte("not json"), 0600))
		require.False(t, client.RestoreSession())
		require.False(t, client.Authenticated())
	})

	t.Run("resolve community maps a name to an id", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/api/v3/community", r.URL.Path)
			require.Equal(t, "movies@lemmy.ml", r.URL.Query().Get("name"))
			require.Equal(t, "Bearer preset.token", r.Header.Get("Authorization"))

			var resp communityResponse
			resp.CommunityView.Community.ID = 42
			require.NoError(t, json.NewEncoder(w).Encode(resp))
		}))
		defer srv.Close()

		client, err := NewClient(srv.URL, OptionHTTPClient(srv.Client()), OptionToken("preset.token"))
		require.NoError(t, err)

		id, err := client.ResolveCommunity(context.Background(), "movies@lemmy.ml")
		require.NoError(t, err)
		require.Equal(t, int64(42), id)
	})

	t.Run("resolve community errors when lemmy returns nothing", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.NoError(t, json.NewEncoder(w).Encode(communityResponse{}))
		}))
		defer srv.Close()

		client, err := NewClient(srv.URL, OptionHTTPClient(srv.Client()))
		require.NoError(t, err)

		_, err = client.ResolveCommunity(context.Background(), "missing")
		require.Error(t, err)
	})

	t.Run("create post carries the magnet uri as the post url", func(t *testing.T) {
		var received CreatePost

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/api/v3/post", r.URL.Path)
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))
			require.NoError(t, json.NewDecoder(r.Body).Decode(&received))

			var resp postResponse
			resp.PostView.Post = Post{ID: 7, ApID: "https://lemmy.invalid/post/7", Name: received.Name}
			require.NoError(t, json.NewEncoder(w).Encode(resp))
		}))
		defer srv.Close()

		client, err := NewClient(srv.URL, OptionHTTPClient(srv.Client()), OptionToken("preset.token"))
		require.NoError(t, err)

		post, err := client.CreatePost(context.Background(), CreatePost{
			Name:        "ubuntu 24.04",
			CommunityID: 42,
			URL:         "magnet:?xt=urn:btih:0123456789abcdef",
			Body:        "an operating system",
		})
		require.NoError(t, err)
		require.Equal(t, int64(7), post.ID)
		require.Equal(t, "https://lemmy.invalid/post/7", post.ApID)
		require.Equal(t, "magnet:?xt=urn:btih:0123456789abcdef", received.URL)
		require.Equal(t, int64(42), received.CommunityID)
	})

	t.Run("a 401 surfaces as ErrUnauthorized", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error":"not_logged_in"}`))
		}))
		defer srv.Close()

		client, err := NewClient(srv.URL, OptionHTTPClient(srv.Client()), OptionToken("stale.token"))
		require.NoError(t, err)

		_, err = client.CreatePost(context.Background(), CreatePost{Name: "x", CommunityID: 1})
		require.ErrorIs(t, err, ErrUnauthorized)
	})

	t.Run("a non-2xx preserves lemmy's own error text", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("couldnt_find_community"))
		}))
		defer srv.Close()

		client, err := NewClient(srv.URL, OptionHTTPClient(srv.Client()))
		require.NoError(t, err)

		_, err = client.ResolveCommunity(context.Background(), "missing")
		require.ErrorContains(t, err, "couldnt_find_community")
	})

	t.Run("upload image posts multipart and returns the hosted url", func(t *testing.T) {
		var uploaded []byte

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/pictrs/image", r.URL.Path)

			part, header, err := r.FormFile("images[]")
			require.NoError(t, err)
			defer part.Close()
			require.Equal(t, "poster.png", header.Filename)

			uploaded, err = io.ReadAll(part)
			require.NoError(t, err)

			_, _ = w.Write([]byte(`{"msg":"ok","files":[{"file":"abc123.png","delete_token":"tok"}]}`))
		}))
		defer srv.Close()

		media := filepath.Join(t.TempDir(), "poster.png")
		require.NoError(t, os.WriteFile(media, []byte("not really a png"), 0600))

		client, err := NewClient(srv.URL, OptionHTTPClient(srv.Client()), OptionToken("preset.token"))
		require.NoError(t, err)

		hosted, err := client.UploadImage(context.Background(), media)
		require.NoError(t, err)
		require.Equal(t, srv.URL+"/pictrs/image/abc123.png", hosted)
		require.Equal(t, []byte("not really a png"), uploaded)
	})

	t.Run("upload image errors when nothing was stored", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{"msg":"rejected","files":[]}`))
		}))
		defer srv.Close()

		media := filepath.Join(t.TempDir(), "poster.png")
		require.NoError(t, os.WriteFile(media, []byte("not really a png"), 0600))

		client, err := NewClient(srv.URL, OptionHTTPClient(srv.Client()))
		require.NoError(t, err)

		_, err = client.UploadImage(context.Background(), media)
		require.ErrorContains(t, err, "rejected")
	})

	t.Run("a relative instance url is rejected", func(t *testing.T) {
		_, err := NewClient("lemmy.ml")
		require.Error(t, err)
	})
}
