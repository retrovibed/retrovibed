package searchplugin

import (
	"context"
	"os/exec"
	"path/filepath"
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

	seq := r.Search(ctx, "video", "ubuntu")

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

	seq := r.Search(ctx, "video", "ubuntu")
	var count int
	for range seq.Each(ctx) {
		count++
	}
	require.NoError(t, seq.Err())
	require.Equal(t, 1, count)
}
