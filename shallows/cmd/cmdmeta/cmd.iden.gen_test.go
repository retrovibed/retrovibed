package cmdmeta_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/cmd/cmdmeta"
	cmdtestx "github.com/retrovibed/retrovibed/shallows/cmd/cmdtestx/v2"
	"github.com/stretchr/testify/require"
)

func TestGenerateID(t *testing.T) {
	harness := cmdtestx.New(t, &cmdmeta.Identity{}).Options(cmdtestx.Standard(t))

	t.Run("generates identity at a different path", func(t *testing.T) {
		keypath := filepath.Join(t.TempDir(), "custom-id")

		require.NoError(t, harness.Run("generate", "--private-key-path", keypath, "seed"))

		require.FileExists(t, keypath)
		require.FileExists(t, keypath+".pub")
	})

	t.Run("same seed at different paths generates the same identity", func(t *testing.T) {
		dir := t.TempDir()
		first := filepath.Join(dir, "first")
		second := filepath.Join(dir, "second")

		require.NoError(t, harness.Run("generate", "--private-key-path", first, "seed"))
		require.NoError(t, harness.Run("generate", "--private-key-path", second, "seed"))

		firstpub, err := os.ReadFile(first + ".pub")
		require.NoError(t, err)
		secondpub, err := os.ReadFile(second + ".pub")
		require.NoError(t, err)
		require.Equal(t, firstpub, secondpub)
	})

	t.Run("does not touch an identity at another path", func(t *testing.T) {
		dir := t.TempDir()
		existing := filepath.Join(dir, "existing")
		other := filepath.Join(dir, "other")

		require.NoError(t, harness.Run("generate", "--private-key-path", existing, "seed-a"))
		before, err := os.ReadFile(existing)
		require.NoError(t, err)

		require.NoError(t, harness.Run("generate", "--private-key-path", other, "seed-b"))

		after, err := os.ReadFile(existing)
		require.NoError(t, err)
		require.Equal(t, before, after)
		require.FileExists(t, other)
	})
}
