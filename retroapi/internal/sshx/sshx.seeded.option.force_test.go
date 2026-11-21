package sshx_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/internal/sshx"
	"github.com/stretchr/testify/require"
)

func TestSeededOptionForce(t *testing.T) {
	t.Run("true overrides a false force argument", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "id")

		_, err := sshx.Seeded(context.Background(), "seed1", false, path)
		require.NoError(t, err)

		_, err = sshx.Seeded(context.Background(), "seed2", false, path, sshx.SeededOptionForce(true))
		require.NoError(t, err)
	})

	t.Run("false overrides a true force argument", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "id")

		_, err := sshx.Seeded(context.Background(), "seed1", false, path)
		require.NoError(t, err)

		_, err = sshx.Seeded(context.Background(), "seed2", true, path, sshx.SeededOptionForce(false))
		require.Error(t, err)
	})
}
