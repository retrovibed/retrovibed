package sshx_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/internal/sshx"
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

func TestLoad(t *testing.T) {
	t.Run("loads an existing cached identity", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "id")

		seeded, err := sshx.Seeded(context.Background(), "testseed", false, path)
		require.NoError(t, err)

		loaded, err := sshx.Load(path)
		require.NoError(t, err)
		require.Equal(t, seeded.PublicKey().Marshal(), loaded.PublicKey().Marshal())
	})

	t.Run("errors when no identity exists", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "id")

		_, err := sshx.Load(path)
		require.Error(t, err)
	})
}

func TestSeeded(t *testing.T) {
	t.Run("creates key at new path", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "id")

		s, err := sshx.Seeded(context.Background(), "testseed", false, path)
		require.NoError(t, err)
		require.NotNil(t, s)
		require.FileExists(t, path)
		require.FileExists(t, path+".pub")
	})

	t.Run("returns error when key exists and force is false", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "id")

		_, err := sshx.Seeded(context.Background(), "testseed", false, path)
		require.NoError(t, err)

		_, err = sshx.Seeded(context.Background(), "testseed", false, path)
		require.Error(t, err)
		require.Contains(t, err.Error(), "an identity already exists")
		require.Contains(t, err.Error(), "--force")
	})

	t.Run("backs up existing keys when force is true", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "id")

		_, err := sshx.Seeded(context.Background(), "seed1", false, path)
		require.NoError(t, err)

		origPriv, err := os.ReadFile(path)
		require.NoError(t, err)

		s, err := sshx.Seeded(context.Background(), "seed2", true, path)
		require.NoError(t, err)
		require.NotNil(t, s)

		require.FileExists(t, path)

		matches, err := filepath.Glob(filepath.Join(dir, "id.[0-9]*"))
		require.NoError(t, err)
		require.Len(t, matches, 1)

		backupContent, err := os.ReadFile(matches[0])
		require.NoError(t, err)
		require.Equal(t, origPriv, backupContent)
	})

	t.Run("same seed produces same public key", func(t *testing.T) {
		dir1 := t.TempDir()
		dir2 := t.TempDir()

		s1, err := sshx.Seeded(context.Background(), "deterministic", false, filepath.Join(dir1, "id"))
		require.NoError(t, err)

		s2, err := sshx.Seeded(context.Background(), "deterministic", false, filepath.Join(dir2, "id"))
		require.NoError(t, err)

		require.Equal(t, s1.PublicKey().Marshal(), s2.PublicKey().Marshal())
	})

	t.Run("different seeds produce different keys", func(t *testing.T) {
		dir1 := t.TempDir()
		dir2 := t.TempDir()

		s1, err := sshx.Seeded(context.Background(), "seed-a", false, filepath.Join(dir1, "id"))
		require.NoError(t, err)

		s2, err := sshx.Seeded(context.Background(), "seed-b", false, filepath.Join(dir2, "id"))
		require.NoError(t, err)

		require.NotEqual(t, s1.PublicKey().Marshal(), s2.PublicKey().Marshal())
	})
}
