package sshx_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/internal/sshx"
	"github.com/stretchr/testify/require"
)

func TestUnseeded(t *testing.T) {
	t.Run("removes seeded key files", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "id")

		_, err := sshx.Seeded(context.Background(), "testseed", false, path)
		require.NoError(t, err)
		require.FileExists(t, path)
		require.FileExists(t, path+".pub")

		err = sshx.Unseeded(path)
		require.NoError(t, err)
		require.NoFileExists(t, path)
		require.NoFileExists(t, path+".pub")
	})

	t.Run("succeeds when files do not exist", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "nonexistent")

		err := sshx.Unseeded(path)
		require.NoError(t, err)
	})
}
