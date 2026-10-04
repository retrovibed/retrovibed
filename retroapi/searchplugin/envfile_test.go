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

	t.Run("skips blank values so plugin defaults apply", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "plugin.env")
		require.NoError(t, os.WriteFile(path, []byte("FOO=bar\nEMPTY=\nSPACES=   \n"), 0600))

		pairs, err := readEnvFile(path)
		require.NoError(t, err)
		require.Equal(t, []string{"FOO=bar"}, pairs)
	})

	t.Run("strips trailing inline comment from value", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "plugin.env")
		require.NoError(t, os.WriteFile(path, []byte("UNIT3D_VERBOSE=1 # set to 1 to log the raw http request/response for every api call\n"), 0600))

		pairs, err := readEnvFile(path)
		require.NoError(t, err)
		require.Equal(t, []string{"UNIT3D_VERBOSE=1"}, pairs)
	})

	t.Run("skips blank values followed by an inline comment", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "plugin.env")
		require.NoError(t, os.WriteFile(path, []byte("UNIT3D_APIKEY= # api key for requests; required unless baked in via -X main.apiKey\nUNIT3D_DOMAIN= # base url for the unit3d api; defaults to https://yu-scene.net\nUNIT3D_VERBOSE=1 # set to 1 to log the raw http request/response for every api call\n"), 0600))

		pairs, err := readEnvFile(path)
		require.NoError(t, err)
		require.Equal(t, []string{"UNIT3D_VERBOSE=1"}, pairs)
	})

	t.Run("skips malformed lines", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "plugin.env")
		require.NoError(t, os.WriteFile(path, []byte("FOO=bar\nnotapair\n"), 0600))

		pairs, err := readEnvFile(path)
		require.NoError(t, err)
		require.Equal(t, []string{"FOO=bar"}, pairs)
	})
}
