package searchplugin

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadEnvFile(t *testing.T) {
	t.Run("missing file returns nil no error", func(t *testing.T) {
		pairs, err := readEnvFile(filepath.Join(t.TempDir(), "missing.env"))
		require.NoError(t, err)
		require.Nil(t, pairs)
	})

	t.Run("skips blank and comment lines", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "plugin.env")
		require.NoError(t, os.WriteFile(path, []byte("FOO=bar\n\n# a comment\nBAZ=qux\n"), 0600))

		pairs, err := readEnvFile(path)
		require.NoError(t, err)
		require.Equal(t, []string{"FOO=bar", "BAZ=qux"}, pairs)
	})

	t.Run("skips malformed lines", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "plugin.env")
		require.NoError(t, os.WriteFile(path, []byte("FOO=bar\nnotapair\n"), 0600))

		pairs, err := readEnvFile(path)
		require.NoError(t, err)
		require.Equal(t, []string{"FOO=bar"}, pairs)
	})
}
