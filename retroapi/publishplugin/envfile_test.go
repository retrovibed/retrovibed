package publishplugin

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

	t.Run("skips blank values so plugin defaults apply", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "plugin.env")
		require.NoError(t, os.WriteFile(path, []byte("FOO=bar\n# a comment\nEMPTY=\nSPACES=   \n"), 0600))

		pairs, err := readEnvFile(path)
		require.NoError(t, err)
		require.Equal(t, []string{"FOO=bar"}, pairs)
	})
}
