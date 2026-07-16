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

// writeFixturePlugin writes a minimal Go module outside the retrovibed
// go.work workspace (so `go build` inside it never picks up the parent
// workspace) that only needs to compile to a valid wasm module - it does
// not need to implement the actual searchplugin protocol.
func writeFixturePlugin(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module searchplugintestdata\n\ngo 1.24\n"), 0600))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\n\nfunc main() {}\n"), 0600))

	return dir
}

func TestSearchPluginInstall(t *testing.T) {
	t.Run("compiles and installs a plugin with configuration", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", t.TempDir())
		dir := writeFixturePlugin(t)

		require.NoError(t, cmdtestx.Execute(t, genparser(t), "search", "plugin", "install", dir,
			"--name", "noop",
			"-e", "FOO=bar",
		))

		wasm, err := os.ReadFile(filepath.Join(searchplugin.SearchPluginDir(), "noop.wasm"))
		require.NoError(t, err)
		require.Equal(t, []byte{0x00, 0x61, 0x73, 0x6D}, wasm[:4])

		pairs, err := envx.FromPath(filepath.Join(searchplugin.SearchPluginDir(), "noop.env"))
		require.NoError(t, err)
		require.Equal(t, []string{"FOO=bar"}, pairs)
	})

	t.Run("trims a .wasm suffix passed to --name", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", t.TempDir())
		dir := writeFixturePlugin(t)

		require.NoError(t, cmdtestx.Execute(t, genparser(t), "search", "plugin", "install", dir,
			"--name", "noop.wasm",
		))

		_, err := os.Stat(filepath.Join(searchplugin.SearchPluginDir(), "noop.wasm"))
		require.NoError(t, err)
		_, err = os.Stat(filepath.Join(searchplugin.SearchPluginDir(), "noop.wasm.wasm"))
		require.True(t, os.IsNotExist(err))
	})

	t.Run("defaults name to the repository's base name", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", t.TempDir())
		dir := writeFixturePlugin(t)
		renamed := filepath.Join(filepath.Dir(dir), "myplugin")
		require.NoError(t, os.Rename(dir, renamed))

		require.NoError(t, cmdtestx.Execute(t, genparser(t), "search", "plugin", "install", renamed))

		_, err := os.Stat(filepath.Join(searchplugin.SearchPluginDir(), "myplugin.wasm"))
		require.NoError(t, err)
	})
}
