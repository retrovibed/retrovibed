package envx

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFromPath(t *testing.T) {
	t.Run("missing file returns nil no error", func(t *testing.T) {
		pairs, err := FromPath(filepath.Join(t.TempDir(), "missing.env"))
		require.NoError(t, err)
		require.Nil(t, pairs)
	})

	t.Run("reads existing file", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "existing.env")
		require.NoError(t, os.WriteFile(path, []byte("A=1\nB=2\n"), 0600))

		pairs, err := FromPath(path)
		require.NoError(t, err)
		require.Equal(t, []string{"A=1", "B=2"}, pairs)
	})
}

func TestMerge(t *testing.T) {
	t.Run("overrides existing key in place", func(t *testing.T) {
		got := Merge([]string{"A=1", "B=2"}, "A=updated")
		require.Equal(t, []string{"A=updated", "B=2"}, got)
	})

	t.Run("appends new key in order given", func(t *testing.T) {
		got := Merge([]string{"A=1"}, "B=2", "C=3")
		require.Equal(t, []string{"A=1", "B=2", "C=3"}, got)
	})

	t.Run("skips malformed updates", func(t *testing.T) {
		got := Merge([]string{"A=1"}, "notapair")
		require.Equal(t, []string{"A=1"}, got)
	})

	t.Run("does not mutate existing slice", func(t *testing.T) {
		existing := []string{"A=1"}
		Merge(existing, "A=updated")
		require.Equal(t, []string{"A=1"}, existing)
	})
}

func TestWriteFile(t *testing.T) {
	t.Run("round trips through FromPath", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "out.env")
		require.NoError(t, WriteFile(path, []string{"A=1", "B=2"}))

		pairs, err := FromPath(path)
		require.NoError(t, err)
		require.Equal(t, []string{"A=1", "B=2"}, pairs)
	})
}
