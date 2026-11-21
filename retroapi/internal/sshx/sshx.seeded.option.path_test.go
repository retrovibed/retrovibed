package sshx_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/internal/sshx"
	"github.com/stretchr/testify/require"
)

func TestSeededOptionPath(t *testing.T) {
	t.Run("overrides the path argument", func(t *testing.T) {
		dir := t.TempDir()
		ignored := filepath.Join(dir, "ignored")
		path := filepath.Join(dir, "id")

		_, err := sshx.Seeded(context.Background(), "seed", false, ignored, sshx.SeededOptionPath(path))
		require.NoError(t, err)
		require.FileExists(t, path)
		require.FileExists(t, path+".pub")
		require.NoFileExists(t, ignored)
	})
}
