package searchplugin

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

func TestEnvironment(t *testing.T) {
	t.Run("returns the plugin's declaration", func(t *testing.T) {
		wasmPath := filepath.Join(t.TempDir(), "env.wasm")

		build := exec.Command("go", "build", "-o", wasmPath, "./.fixtures/envplugin")
		build.Env = append(build.Environ(), "GOOS=wasip1", "GOARCH=wasm")
		out, err := build.CombinedOutput()
		require.NoError(t, err, string(out))

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		r, err := newRegistry(ctx, defaultSocket(), OptionConfigDir(t.TempDir()), OptionCacheDir(t.TempDir()))
		require.NoError(t, err)
		require.NoError(t, r.Load(ctx, wasmPath))

		var reg E = r
		declared, err := reg.Environment(ctx, wasmPath)
		require.NoError(t, err)
		require.Equal(t, []envfile.Variable{
			{Key: "PLUGIN_TOKEN", Value: "", Hint: "token echoed back in results"},
		}, envfile.Parse(string(declared)))
	})

	t.Run("fails for a plugin without an env command", func(t *testing.T) {
		wasmPath := filepath.Join(t.TempDir(), "fail.wasm")

		build := exec.Command("go", "build", "-o", wasmPath, "./.fixtures/failplugin")
		build.Env = append(build.Environ(), "GOOS=wasip1", "GOARCH=wasm")
		out, err := build.CombinedOutput()
		require.NoError(t, err, string(out))

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		r, err := newRegistry(ctx, defaultSocket(), OptionConfigDir(t.TempDir()), OptionCacheDir(t.TempDir()))
		require.NoError(t, err)
		require.NoError(t, r.Load(ctx, wasmPath))

		_, err = r.Environment(ctx, wasmPath)
		require.Error(t, err)
	})

	t.Run("reports not loaded", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		r, err := newRegistry(ctx, defaultSocket())
		require.NoError(t, err)

		_, err = r.Environment(ctx, filepath.Join(t.TempDir(), "missing.wasm"))
		require.ErrorIs(t, err, ErrNotLoaded)
	})

	t.Run("unimplemented returns errUnsupported", func(t *testing.T) {
		var reg E = Unimplemented{}

		_, err := reg.Environment(context.Background(), "irrelevant")
		require.ErrorIs(t, err, errors.ErrUnsupported)
	})
}
