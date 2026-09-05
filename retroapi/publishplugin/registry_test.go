package publishplugin

import (
	"context"
	"errors"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/envfile"
	"github.com/stretchr/testify/require"
)

func TestRegistry(t *testing.T) {
	t.Run("satisfies T", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		r, err := newRegistry(ctx, defaultSocket())
		require.NoError(t, err)

		var reg T = r
		_, err = reg.Publish(ctx, filepath.Join(t.TempDir(), "missing.wasm"), Request{})
		require.ErrorIs(t, err, ErrNotLoaded)
	})

	t.Run("publish decodes plugin output", func(t *testing.T) {
		wasmPath := filepath.Join(t.TempDir(), "echo.wasm")

		build := exec.Command("go", "build", "-o", wasmPath, "./.fixtures/echopublisher")
		build.Env = append(build.Environ(), "GOOS=wasip1", "GOARCH=wasm")
		out, err := build.CombinedOutput()
		require.NoError(t, err, string(out))

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		r, err := newRegistry(ctx, defaultSocket())
		require.NoError(t, err)
		require.NoError(t, r.Load(ctx, wasmPath))

		result, err := r.Publish(ctx, wasmPath, Request{Title: "hello", Mimetype: "video/mp4"})
		require.NoError(t, err)
		require.Equal(t, "https://example.invalid/echo/hello", result.URL)
		require.Equal(t, "echo", result.ExternalID)
		require.Equal(t, "published", result.Status)
	})

	t.Run("publish forwards the link to the plugin", func(t *testing.T) {
		wasmPath := filepath.Join(t.TempDir(), "echo.wasm")

		build := exec.Command("go", "build", "-o", wasmPath, "./.fixtures/echopublisher")
		build.Env = append(build.Environ(), "GOOS=wasip1", "GOARCH=wasm")
		out, err := build.CombinedOutput()
		require.NoError(t, err, string(out))

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		r, err := newRegistry(ctx, defaultSocket())
		require.NoError(t, err)
		require.NoError(t, r.Load(ctx, wasmPath))

		result, err := r.Publish(ctx, wasmPath, Request{Title: "hello", Link: "magnet:?xt=urn:btih:0123456789abcdef"})
		require.NoError(t, err)
		require.Equal(t, "magnet:?xt=urn:btih:0123456789abcdef", result.ExternalID)
	})

	t.Run("environment returns the plugin's declaration", func(t *testing.T) {
		wasmPath := filepath.Join(t.TempDir(), "echo.wasm")

		build := exec.Command("go", "build", "-o", wasmPath, "./.fixtures/echopublisher")
		build.Env = append(build.Environ(), "GOOS=wasip1", "GOARCH=wasm")
		out, err := build.CombinedOutput()
		require.NoError(t, err, string(out))

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		r, err := newRegistry(ctx, defaultSocket())
		require.NoError(t, err)
		require.NoError(t, r.Load(ctx, wasmPath))

		var reg E = r
		declared, err := reg.Environment(ctx, wasmPath)
		require.NoError(t, err)
		require.Equal(t, []envfile.Variable{
			{Key: "ECHO_STATUS", Value: "published", Hint: "status reported for every echo publish"},
		}, envfile.Parse(string(declared)))
	})

	t.Run("environment reports not loaded", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		r, err := newRegistry(ctx, defaultSocket())
		require.NoError(t, err)

		_, err = r.Environment(ctx, filepath.Join(t.TempDir(), "missing.wasm"))
		require.ErrorIs(t, err, ErrNotLoaded)
	})

	t.Run("publish surfaces a failing plugin as an error", func(t *testing.T) {
		wasmPath := filepath.Join(t.TempDir(), "fail.wasm")

		build := exec.Command("go", "build", "-o", wasmPath, "./.fixtures/failpublisher")
		build.Env = append(build.Environ(), "GOOS=wasip1", "GOARCH=wasm")
		out, err := build.CombinedOutput()
		require.NoError(t, err, string(out))

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		r, err := newRegistry(ctx, defaultSocket())
		require.NoError(t, err)
		require.NoError(t, r.Load(ctx, wasmPath))

		_, err = r.Publish(ctx, wasmPath, Request{Title: "hello"})
		require.Error(t, err)
	})

	t.Run("unload then publish reports not loaded", func(t *testing.T) {
		wasmPath := filepath.Join(t.TempDir(), "echo.wasm")

		build := exec.Command("go", "build", "-o", wasmPath, "./.fixtures/echopublisher")
		build.Env = append(build.Environ(), "GOOS=wasip1", "GOARCH=wasm")
		out, err := build.CombinedOutput()
		require.NoError(t, err, string(out))

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		r, err := newRegistry(ctx, defaultSocket())
		require.NoError(t, err)
		require.NoError(t, r.Load(ctx, wasmPath))
		r.Unload(wasmPath)

		_, err = r.Publish(ctx, wasmPath, Request{Title: "hello"})
		require.ErrorIs(t, err, ErrNotLoaded)
	})
}

func TestPluginDirs(t *testing.T) {
	t.Run("PluginConfigDir and PluginCacheDir build expected paths", func(t *testing.T) {
		require.Equal(t, filepath.Join("/tmp/config", "publish.d", "noop.env"), EnvPath(filepath.Join("/tmp/config", "publish.d", "noop.wasm")))
		require.Equal(t, filepath.Join("/tmp/config", "publish.d"), PublishPluginDir("/tmp/config"))
		require.Equal(t, filepath.Join("/tmp/config", "publish.d", "noop.config.d"), PluginConfigDir("/tmp/config", "noop"))
		require.Equal(t, filepath.Join("/tmp/cache", "publish.d", "noop.cache.d"), PluginCacheDir("/tmp/cache", "noop"))

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		r, err := newRegistry(ctx, defaultSocket())
		require.NoError(t, err)
		r.configDir = "/tmp/config"
		r.cacheDir = "/tmp/cache"

		require.Equal(t, filepath.Join("/tmp/config", "publish.d", "noop.config.d"), r.PluginConfigDir("noop"))
		require.Equal(t, filepath.Join("/tmp/cache", "publish.d", "noop.cache.d"), r.PluginCacheDir("noop"))
	})
}

func TestUnimplemented(t *testing.T) {
	t.Run("publish returns errUnsupported", func(t *testing.T) {
		var reg T = Unimplemented{}

		_, err := reg.Publish(context.Background(), "irrelevant", Request{})
		require.ErrorIs(t, err, errors.ErrUnsupported)
	})

	t.Run("environment returns errUnsupported", func(t *testing.T) {
		var reg E = Unimplemented{}

		_, err := reg.Environment(context.Background(), "irrelevant")
		require.ErrorIs(t, err, errors.ErrUnsupported)
	})
}
