package searchplugin

import (
	"context"
	"net"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRegistrySearchDecodesPluginOutput(t *testing.T) {
	wasmPath := filepath.Join(t.TempDir(), "echo.wasm")

	build := exec.Command("go", "build", "-o", wasmPath, "./testdata/echoplugin")
	build.Env = append(build.Environ(), "GOOS=wasip1", "GOARCH=wasm")
	out, err := build.CombinedOutput()
	require.NoError(t, err, string(out))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	r, err := newRegistry(ctx, defaultSocket())
	require.NoError(t, err)
	require.NoError(t, r.Load(ctx, wasmPath))

	seq := r.Search(ctx, []string{"video"}, "ubuntu")

	var results []string
	for imp := range seq.Each(ctx) {
		results = append(results, imp.Magnet)
		require.EqualValues(t, 42, imp.Health)
		require.Equal(t, "video", imp.Mimetype)
	}
	require.NoError(t, seq.Err())
	require.Len(t, results, 1)
	require.Equal(t, "magnet:?xt=urn:btih:1111111111111111111111111111111111111111&dn=ubuntu", results[0])
}

func TestRegistrySearchSkipsFailingPluginWithoutAbortingOthers(t *testing.T) {
	dir := t.TempDir()

	echoPath := filepath.Join(dir, "echo.wasm")
	build := exec.Command("go", "build", "-o", echoPath, "./testdata/echoplugin")
	build.Env = append(build.Environ(), "GOOS=wasip1", "GOARCH=wasm")
	out, err := build.CombinedOutput()
	require.NoError(t, err, string(out))

	failPath := filepath.Join(dir, "fail.wasm")
	build = exec.Command("go", "build", "-o", failPath, "./testdata/failplugin")
	build.Env = append(build.Environ(), "GOOS=wasip1", "GOARCH=wasm")
	out, err = build.CombinedOutput()
	require.NoError(t, err, string(out))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	r, err := newRegistry(ctx, defaultSocket())
	require.NoError(t, err)
	require.NoError(t, r.Load(ctx, echoPath))
	require.NoError(t, r.Load(ctx, failPath))

	seq := r.Search(ctx, []string{"video"}, "ubuntu")
	var count int
	for range seq.Each(ctx) {
		count++
	}
	require.NoError(t, seq.Err())
	require.Equal(t, 1, count)
}

func TestRegistryBlocksNonPublicConnections(t *testing.T) {
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
	build := exec.Command("go", "build", "-o", wasmPath, "./testdata/dialplugin")
	build.Env = append(build.Environ(), "GOOS=wasip1", "GOARCH=wasm")
	out, err := build.CombinedOutput()
	require.NoError(t, err, string(out))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	r, err := newRegistry(ctx, defaultSocket())
	require.NoError(t, err)
	require.NoError(t, r.Load(ctx, wasmPath))

	seq := r.Search(ctx, []string{"test"}, ln.Addr().String())

	var results []string
	for imp := range seq.Each(ctx) {
		results = append(results, imp.Magnet)
	}
	require.NoError(t, seq.Err())
	require.Len(t, results, 1)
	require.True(t, strings.HasPrefix(results[0], "blocked:"), "expected connection to loopback address to be blocked, got: %s", results[0])
	require.False(t, accepted.Load(), "listener should never have received a connection from the sandboxed plugin")
}
