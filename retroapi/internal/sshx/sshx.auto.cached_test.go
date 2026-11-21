package sshx_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/internal/sshx"
	"github.com/stretchr/testify/require"
)

func TestAutoCached(t *testing.T) {
	t.Run("generates and writes the key pair when no cache exists", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "id")

		s, err := sshx.AutoCached(sshx.NewKeyGenSeeded("seed"), path)
		require.NoError(t, err)
		require.NotNil(t, s)
		require.FileExists(t, path)
		require.FileExists(t, path+".pub")
	})

	t.Run("writes files with restrictive permissions", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "id")

		_, err := sshx.AutoCached(sshx.NewKeyGenSeeded("seed"), path)
		require.NoError(t, err)

		info, err := os.Stat(path)
		require.NoError(t, err)
		require.Equal(t, os.FileMode(0600), info.Mode().Perm())

		info, err = os.Stat(path + ".pub")
		require.NoError(t, err)
		require.Equal(t, os.FileMode(0600), info.Mode().Perm())
	})

	t.Run("returns the cached key instead of generating a new one", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "id")

		s1, err := sshx.AutoCached(sshx.NewKeyGenSeeded("seed-a"), path)
		require.NoError(t, err)

		s2, err := sshx.AutoCached(sshx.NewKeyGenSeeded("seed-b"), path)
		require.NoError(t, err)

		require.Equal(t, s1.PublicKey().Marshal(), s2.PublicKey().Marshal())
	})

	t.Run("regenerates when the cached key is invalid", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "id")
		require.NoError(t, os.WriteFile(path, []byte("not a key"), 0600))

		s, err := sshx.AutoCached(sshx.NewKeyGenSeeded("seed"), path)
		require.NoError(t, err)
		require.NotNil(t, s)

		expected, err := sshx.SignerFromGenerator(sshx.NewKeyGenSeeded("seed"))
		require.NoError(t, err)
		require.Equal(t, expected.PublicKey().Marshal(), s.PublicKey().Marshal())
	})

	t.Run("creates missing parent directories", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "missing", "id")

		s, err := sshx.AutoCached(sshx.NewKeyGenSeeded("seed"), path)
		require.NoError(t, err)
		require.NotNil(t, s)
		require.FileExists(t, path)
		require.FileExists(t, path+".pub")
	})

	t.Run("returns an error when the parent directory cannot be created due to permissions", func(t *testing.T) {
		// the readonly directory is owned by the caller and would be chmod'd writable if it were the
		// parent of the key, so the key lives in a child that can't be created beneath it.
		dir := filepath.Join(t.TempDir(), "readonly")
		require.NoError(t, os.Mkdir(dir, 0500))

		_, err := sshx.AutoCached(sshx.NewKeyGenSeeded("seed"), filepath.Join(dir, "child", "id"))
		require.Error(t, err)
		require.NoDirExists(t, filepath.Join(dir, "child"))
	})
}
