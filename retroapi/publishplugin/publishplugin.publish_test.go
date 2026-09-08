package publishplugin

import (
	"bytes"
	"context"
	"crypto/rand"
	"image"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/internal/md5x"
	"github.com/stretchr/testify/require"
)

func TestPublish(t *testing.T) {
	t.Run("invokes the named plugin and decodes its result", func(t *testing.T) {
		wasmPath := filepath.Join(t.TempDir(), "echo.wasm")

		build := exec.Command("go", "build", "-o", wasmPath, "./.fixtures/echopublisher")
		build.Env = append(build.Environ(), "GOOS=wasip1", "GOARCH=wasm")
		out, err := build.CombinedOutput()
		require.NoError(t, err, string(out))

		mediadir := t.TempDir()
		img := image.NewRGBA(image.Rect(0, 0, 8, 8))
		_, err = rand.Read(img.Pix)
		require.NoError(t, err)

		var encoded bytes.Buffer
		require.NoError(t, png.Encode(&encoded, img))
		require.NoError(t, os.WriteFile(filepath.Join(mediadir, "media.png"), encoded.Bytes(), 0600))
		externalID := md5x.FormatUUID(md5x.Digest(encoded.Bytes()))
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		r, err := newRegistry(ctx, defaultSocket(), OptionConfigDir(t.TempDir()), OptionCacheDir(t.TempDir()))
		require.NoError(t, err)
		require.NoError(t, r.Load(ctx, wasmPath))

		result, err := r.Publish(ctx, wasmPath, Request{
			Title:       "hello",
			Description: "a happy publish",
			Mimetype:    "video/mp4",
			CommunityID: "77b32e4f-4b1e-4b8f-9c0c-2e3f3a4b5c6d",
			Magnet:      "magnet:?xt=urn:btih:0123456789abcdef",
			Directory:   mediadir,
			MediaPath:   "media.png",
		})
		require.NoError(t, err)
		require.Equal(t, &Result{
			URL:        "https://example.invalid/echo/hello",
			ExternalID: externalID,
		}, result)
	})
}
