package searchplugin

import (
	"context"
	"errors"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/ddiscapi"
	"github.com/stretchr/testify/require"
)

func TestRegistry(t *testing.T) {
	t.Run("satisfies T", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		r, err := newRegistry(ctx, defaultSocket())
		require.NoError(t, err)

		var reg T = r

		seq := reg.Search(ctx, []string{"video"}, "ubuntu", false, false)
		for range seq.Each(ctx) {
		}
		require.NoError(t, seq.Err())
	})

	t.Run("satisfies R", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		r, err := newRegistry(ctx, defaultSocket())
		require.NoError(t, err)

		var reg R = r

		seq := reg.Recommend(ctx, []string{"video"}, 5, "", false, false)
		for range seq.Each(ctx) {
		}
		require.NoError(t, seq.Err())
	})

	t.Run("search decodes plugin output", func(t *testing.T) {
		wasmPath := filepath.Join(t.TempDir(), "echo.wasm")

		build := exec.Command("go", "build", "-o", wasmPath, "./.fixtures/echoplugin")
		build.Env = append(build.Environ(), "GOOS=wasip1", "GOARCH=wasm")
		out, err := build.CombinedOutput()
		require.NoError(t, err, string(out))

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		r, err := newRegistry(ctx, defaultSocket())
		require.NoError(t, err)
		require.NoError(t, r.Load(ctx, wasmPath))

		seq := r.Search(ctx, []string{"video"}, "ubuntu", false, false)

		var results []string
		for imp := range seq.Each(ctx) {
			results = append(results, imp.Uri)
			require.EqualValues(t, 42, imp.Health)
			require.Equal(t, "video", imp.Mimetype)
		}
		require.NoError(t, seq.Err())
		require.Len(t, results, 1)
		require.Equal(t, "magnet:?xt=urn:btih:1111111111111111111111111111111111111111&dn=ubuntu", results[0])
	})

	t.Run("search skips failing plugin without aborting others", func(t *testing.T) {
		dir := t.TempDir()

		echoPath := filepath.Join(dir, "echo.wasm")
		build := exec.Command("go", "build", "-o", echoPath, "./.fixtures/echoplugin")
		build.Env = append(build.Environ(), "GOOS=wasip1", "GOARCH=wasm")
		out, err := build.CombinedOutput()
		require.NoError(t, err, string(out))

		failPath := filepath.Join(dir, "fail.wasm")
		build = exec.Command("go", "build", "-o", failPath, "./.fixtures/failplugin")
		build.Env = append(build.Environ(), "GOOS=wasip1", "GOARCH=wasm")
		out, err = build.CombinedOutput()
		require.NoError(t, err, string(out))

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		r, err := newRegistry(ctx, defaultSocket())
		require.NoError(t, err)
		require.NoError(t, r.Load(ctx, echoPath))
		require.NoError(t, r.Load(ctx, failPath))

		seq := r.Search(ctx, []string{"video"}, "ubuntu", false, false)
		var count int
		for range seq.Each(ctx) {
			count++
		}
		require.NoError(t, seq.Err())
		require.Equal(t, 1, count)
	})

	t.Run("recommend decodes plugin output", func(t *testing.T) {
		wasmPath := filepath.Join(t.TempDir(), "recommend.wasm")

		build := exec.Command("go", "build", "-o", wasmPath, "./.fixtures/recommendplugin")
		build.Env = append(build.Environ(), "GOOS=wasip1", "GOARCH=wasm")
		out, err := build.CombinedOutput()
		require.NoError(t, err, string(out))

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		r, err := newRegistry(ctx, defaultSocket())
		require.NoError(t, err)
		require.NoError(t, r.Load(ctx, wasmPath))

		seq := r.Recommend(ctx, []string{"video"}, 3, "", false, false)

		var results []string
		for imp := range seq.Each(ctx) {
			results = append(results, imp.Uri)
			require.Equal(t, "video", imp.Mimetype)
			require.Equal(t, "limit=3", imp.Title)
		}
		require.NoError(t, seq.Err())
		require.Len(t, results, 1)
		require.Equal(t, "magnet:?xt=urn:btih:2222222222222222222222222222222222222222&dn=recommended", results[0])
	})

	t.Run("recommend omits --limit when zero", func(t *testing.T) {
		wasmPath := filepath.Join(t.TempDir(), "recommend.wasm")

		build := exec.Command("go", "build", "-o", wasmPath, "./.fixtures/recommendplugin")
		build.Env = append(build.Environ(), "GOOS=wasip1", "GOARCH=wasm")
		out, err := build.CombinedOutput()
		require.NoError(t, err, string(out))

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		r, err := newRegistry(ctx, defaultSocket())
		require.NoError(t, err)
		require.NoError(t, r.Load(ctx, wasmPath))

		seq := r.Recommend(ctx, []string{"video"}, 0, "", false, false)

		var results []*ddiscapi.Import
		for imp := range seq.Each(ctx) {
			results = append(results, imp)
		}
		require.NoError(t, seq.Err())
		require.Len(t, results, 1)
		require.Equal(t, "limit=", results[0].Title)
	})

	t.Run("blocks non-public connections", func(t *testing.T) {
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		require.NoError(t, err)
		defer ln.Close()

		var accepted atomic.Bool
		go func() {
			for {
				conn, err := ln.Accept()
				if err != nil {
					return
				}
				accepted.Store(true)
				conn.Close()
			}
		}()

		wasmPath := filepath.Join(t.TempDir(), "dial.wasm")
		build := exec.Command("go", "build", "-o", wasmPath, "./.fixtures/dialplugin")
		build.Env = append(build.Environ(), "GOOS=wasip1", "GOARCH=wasm")
		out, err := build.CombinedOutput()
		require.NoError(t, err, string(out))

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		r, err := newRegistry(ctx, defaultSocket())
		require.NoError(t, err)
		require.NoError(t, r.Load(ctx, wasmPath))

		seq := r.Search(ctx, []string{"test"}, ln.Addr().String(), false, false)

		var results []string
		for imp := range seq.Each(ctx) {
			results = append(results, imp.Uri)
		}
		require.NoError(t, seq.Err())
		require.Len(t, results, 1)
		require.True(t, strings.HasPrefix(results[0], "blocked:"), "expected connection to loopback address to be blocked, got: %s", results[0])
		require.False(t, accepted.Load(), "listener should never have received a connection from the sandboxed plugin")
	})
}

func TestPluginDirs(t *testing.T) {
	t.Run("PluginConfigDir and PluginCacheDir build expected paths", func(t *testing.T) {
		require.Equal(t, filepath.Join("/tmp/config", "search.d"), SearchPluginDir("/tmp/config"))
		require.Equal(t, filepath.Join("/tmp/config", "search.d", "noop.config.d"), PluginConfigDir("/tmp/config", "noop"))
		require.Equal(t, filepath.Join("/tmp/cache", "search.d", "noop.cache.d"), PluginCacheDir("/tmp/cache", "noop"))

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		r, err := newRegistry(ctx, defaultSocket())
		require.NoError(t, err)
		r.configDir = "/tmp/config"
		r.cacheDir = "/tmp/cache"

		require.Equal(t, filepath.Join("/tmp/config", "search.d", "noop.config.d"), r.PluginConfigDir("noop"))
		require.Equal(t, filepath.Join("/tmp/cache", "search.d", "noop.cache.d"), r.PluginCacheDir("noop"))
	})

	t.Run("runSearchJob creates the per-plugin config and cache directories", func(t *testing.T) {
		wasmPath := filepath.Join(t.TempDir(), "echo.wasm")

		build := exec.Command("go", "build", "-o", wasmPath, "./.fixtures/echoplugin")
		build.Env = append(build.Environ(), "GOOS=wasip1", "GOARCH=wasm")
		out, err := build.CombinedOutput()
		require.NoError(t, err, string(out))

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		r, err := newRegistry(ctx, defaultSocket())
		require.NoError(t, err)
		r.configDir = t.TempDir()
		r.cacheDir = t.TempDir()
		require.NoError(t, r.Load(ctx, wasmPath))

		seq := r.Search(ctx, []string{"video"}, "ubuntu", false, false)
		for range seq.Each(ctx) {
		}
		require.NoError(t, seq.Err())

		require.DirExists(t, r.PluginConfigDir("echo"))
		require.DirExists(t, r.PluginCacheDir("echo"))
	})

	t.Run("guest sees CONFIGURATION_DIRECTORY and CACHE_DIRECTORY through the mount", func(t *testing.T) {
		wasmPath := filepath.Join(t.TempDir(), "dirplugin.wasm")

		build := exec.Command("go", "build", "-o", wasmPath, "./.fixtures/dirplugin")
		build.Env = append(build.Environ(), "GOOS=wasip1", "GOARCH=wasm")
		out, err := build.CombinedOutput()
		require.NoError(t, err, string(out))

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		r, err := newRegistry(ctx, defaultSocket())
		require.NoError(t, err)
		r.configDir = t.TempDir()
		r.cacheDir = t.TempDir()
		require.NoError(t, r.Load(ctx, wasmPath))

		seq := r.Search(ctx, []string{"video"}, "ubuntu", false, false)

		var guestPaths []string
		for imp := range seq.Each(ctx) {
			guestPaths = append(guestPaths, imp.Uri)
		}
		require.NoError(t, seq.Err())
		require.Len(t, guestPaths, 1)
		require.Equal(t, guestPluginCacheDir+"/marker", guestPaths[0])

		marker, err := os.ReadFile(filepath.Join(r.PluginCacheDir("dirplugin"), "marker"))
		require.NoError(t, err)
		require.Equal(t, "dirplugin", string(marker))
	})
}

func TestUnimplemented(t *testing.T) {
	t.Run("search returns errUnsupported", func(t *testing.T) {
		var reg T = Unimplemented{}

		seq := reg.Search(context.Background(), []string{"video"}, "ubuntu", false, false)

		var count int
		for range seq.Each(context.Background()) {
			count++
		}
		require.Equal(t, 0, count)
		require.ErrorIs(t, seq.Err(), errors.ErrUnsupported)
	})

	t.Run("recommend returns errUnsupported", func(t *testing.T) {
		var reg R = Unimplemented{}

		seq := reg.Recommend(context.Background(), []string{"video"}, 5, "", false, false)

		var count int
		for range seq.Each(context.Background()) {
			count++
		}
		require.Equal(t, 0, count)
		require.ErrorIs(t, seq.Err(), errors.ErrUnsupported)
	})
}
