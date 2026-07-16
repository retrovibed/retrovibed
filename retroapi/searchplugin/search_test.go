package searchplugin

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRunSearchJobInjectsSiblingEnvFile(t *testing.T) {
	wasmPath := filepath.Join(t.TempDir(), "env.wasm")

	build := exec.Command("go", "build", "-o", wasmPath, "./testdata/envplugin")
	build.Env = append(build.Environ(), "GOOS=wasip1", "GOARCH=wasm")
	out, err := build.CombinedOutput()
	require.NoError(t, err, string(out))

	envPath := strings.TrimSuffix(wasmPath, ".wasm") + ".env"
	require.NoError(t, os.WriteFile(envPath, []byte("PLUGIN_TOKEN=injected-value\n"), 0600))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	r, err := newRegistry(ctx, defaultSocket())
	require.NoError(t, err)
	require.NoError(t, r.Load(ctx, wasmPath))

	seq := r.Search(ctx, []string{"video"}, "ubuntu")

	var results []string
	for imp := range seq.Each(ctx) {
		results = append(results, imp.Magnet)
	}
	require.NoError(t, seq.Err())
	require.Len(t, results, 1)
	require.Equal(t, "magnet:?xt=urn:btih:1111111111111111111111111111111111111111&dn=injected-value", results[0])
}

func TestRunSearchJobToleratesMissingEnvFile(t *testing.T) {
	wasmPath := filepath.Join(t.TempDir(), "env.wasm")

	build := exec.Command("go", "build", "-o", wasmPath, "./testdata/envplugin")
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
	}
	require.NoError(t, seq.Err())
	require.Len(t, results, 1)
	require.Equal(t, "magnet:?xt=urn:btih:1111111111111111111111111111111111111111&dn=", results[0])
}
