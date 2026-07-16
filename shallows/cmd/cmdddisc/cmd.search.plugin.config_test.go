package cmdddisc_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/searchplugin"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdtestx"
	"github.com/retrovibed/retrovibed/shallows/internal/envx"
	"github.com/stretchr/testify/require"
)

func TestSearchPluginConfig(t *testing.T) {
	t.Run("merges into an existing env file", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", t.TempDir())

		envpath := filepath.Join(searchplugin.SearchPluginDir(), "noop.env")
		require.NoError(t, os.MkdirAll(filepath.Dir(envpath), 0o700))
		require.NoError(t, envx.WriteFile(envpath, []string{"FOO=bar"}))

		require.NoError(t, cmdtestx.Execute(t, genparser(t), "search", "plugin", "config", "noop",
			"-e", "FOO=updated",
			"-e", "BAZ=qux",
		))

		pairs, err := envx.FromPath(envpath)
		require.NoError(t, err)
		require.Equal(t, []string{"FOO=updated", "BAZ=qux"}, pairs)
	})

	t.Run("creates the env file when none exists", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", t.TempDir())

		require.NoError(t, cmdtestx.Execute(t, genparser(t), "search", "plugin", "config", "noop",
			"-e", "FOO=bar",
		))

		envpath := filepath.Join(searchplugin.SearchPluginDir(), "noop.env")
		pairs, err := envx.FromPath(envpath)
		require.NoError(t, err)
		require.Equal(t, []string{"FOO=bar"}, pairs)
	})

	t.Run("accepts a .wasm suffixed name", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", t.TempDir())

		require.NoError(t, cmdtestx.Execute(t, genparser(t), "search", "plugin", "config", "noop.wasm",
			"-e", "FOO=bar",
		))

		envpath := filepath.Join(searchplugin.SearchPluginDir(), "noop.env")
		pairs, err := envx.FromPath(envpath)
		require.NoError(t, err)
		require.Equal(t, []string{"FOO=bar"}, pairs)
	})
}
