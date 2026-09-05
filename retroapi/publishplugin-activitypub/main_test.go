package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/envfile"
	"github.com/retrovibed/retrovibed/retroapi/publishplugin"
	"github.com/stretchr/testify/require"
)

func TestEnvironment(t *testing.T) {
	t.Run("every declared variable parses with a hint", func(t *testing.T) {
		declared := envfile.Parse(environment)
		require.NotEmpty(t, declared)

		for _, v := range declared {
			require.NotEmptyf(t, v.Hint, "%s is declared without a description", v.Key)
		}
	})

	t.Run("declares the variables publishCmd reads", func(t *testing.T) {
		keys := make([]string, 0, len(envfile.Parse(environment)))
		for _, v := range envfile.Parse(environment) {
			keys = append(keys, v.Key)
		}

		require.Equal(t, []string{
			"LEMMY_INSTANCE",
			"LEMMY_COMMUNITY",
			"LEMMY_USERNAME",
			"LEMMY_PASSWORD",
			"LEMMY_TOKEN",
			"LEMMY_TOTP",
			"LEMMY_NSFW",
			"LEMMY_LANGUAGE_ID",
			"LEMMY_THUMBNAIL_MAX",
		}, keys)
	})
}

// TestPluginUnderRegistry builds this plugin the way it actually ships and
// runs it through the real sandbox. Nothing else here proves the wasip1
// build works: the module imports wasinet's host functions, so a bare
// wasm runtime cannot instantiate it - only a registry that supplies them
// can.
func TestPluginUnderRegistry(t *testing.T) {
	t.Run("env runs in the sandbox and declares what publishCmd reads", func(t *testing.T) {
		wasmPath := filepath.Join(t.TempDir(), "lemmy.wasm")

		build := exec.Command("go", "build", "-o", wasmPath, ".")
		build.Env = append(build.Environ(), "GOOS=wasip1", "GOARCH=wasm")
		out, err := build.CombinedOutput()
		require.NoError(t, err, string(out))

		ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
		defer cancel()

		reg, err := publishplugin.NewRegistry(ctx, publishplugin.OptionConfigDir(t.TempDir()), publishplugin.OptionCacheDir(t.TempDir()))
		require.NoError(t, err)
		require.NoError(t, reg.Load(ctx, wasmPath))

		declared, err := reg.Environment(ctx, wasmPath)
		require.NoError(t, err)
		require.Equal(t, envfile.Parse(environment), envfile.Parse(string(declared)))
	})
}

func TestThumbnail(t *testing.T) {
	t.Run("uploads an image within the limit", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{"msg":"ok","files":[{"file":"abc123.png","delete_token":"tok"}]}`))
		}))
		defer srv.Close()

		media := filepath.Join(t.TempDir(), "poster.png")
		require.NoError(t, os.WriteFile(media, []byte("not really a png"), 0600))

		client, err := NewClient(srv.URL, OptionHTTPClient(srv.Client()))
		require.NoError(t, err)

		cmd := publishCmd{Media: media, Mimetype: "image/png", ThumbnailMax: "8 MB"}
		require.Equal(t, srv.URL+"/pictrs/image/abc123.png", cmd.thumbnail(context.Background(), client))
	})

	t.Run("skips video, the common case for torrented media", func(t *testing.T) {
		media := filepath.Join(t.TempDir(), "movie.mkv")
		require.NoError(t, os.WriteFile(media, []byte("not really a movie"), 0600))

		client, err := NewClient("https://lemmy.invalid")
		require.NoError(t, err)

		cmd := publishCmd{Media: media, Mimetype: "video/x-matroska", ThumbnailMax: "8 MB"}
		require.Equal(t, "", cmd.thumbnail(context.Background(), client))
	})

	t.Run("skips an image over the limit", func(t *testing.T) {
		media := filepath.Join(t.TempDir(), "poster.png")
		require.NoError(t, os.WriteFile(media, make([]byte, 2048), 0600))

		client, err := NewClient("https://lemmy.invalid")
		require.NoError(t, err)

		cmd := publishCmd{Media: media, Mimetype: "image/png", ThumbnailMax: "1 KB"}
		require.Equal(t, "", cmd.thumbnail(context.Background(), client))
	})

	t.Run("skips when no media was mounted", func(t *testing.T) {
		client, err := NewClient("https://lemmy.invalid")
		require.NoError(t, err)

		cmd := publishCmd{Mimetype: "image/png", ThumbnailMax: "8 MB"}
		require.Equal(t, "", cmd.thumbnail(context.Background(), client))
	})

	t.Run("skips rather than fails when the upload errors", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer srv.Close()

		media := filepath.Join(t.TempDir(), "poster.png")
		require.NoError(t, os.WriteFile(media, []byte("not really a png"), 0600))

		client, err := NewClient(srv.URL, OptionHTTPClient(srv.Client()))
		require.NoError(t, err)

		cmd := publishCmd{Media: media, Mimetype: "image/png", ThumbnailMax: "8 MB"}
		require.Equal(t, "", cmd.thumbnail(context.Background(), client))
	})
}
