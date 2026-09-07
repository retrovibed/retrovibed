package publishplugin_test

import (
	"bytes"
	"crypto/md5"
	"os"
	"path/filepath"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/internal/md5x"
	"github.com/retrovibed/retrovibed/retroapi/publishplugin"
	"github.com/stretchr/testify/require"
)

func TestIdentity(t *testing.T) {
	wasm := append([]byte{0x00, 0x61, 0x73, 0x6D}, bytes.Repeat([]byte{0x42}, 64)...)

	digest := func(t *testing.T, b []byte) string {
		t.Helper()
		d := md5.New()
		_, err := d.Write(b)
		require.NoError(t, err)
		return md5x.FormatUUID(d)
	}

	t.Run("a canonically named module is its content digest", func(t *testing.T) {
		dir := t.TempDir()
		content := digest(t, wasm)
		require.NoError(t, os.WriteFile(filepath.Join(dir, content+".wasm"), wasm, 0o600))

		id, err := publishplugin.Identity(filepath.Join(dir, content+".wasm"))
		require.NoError(t, err)
		require.Equal(t, content, id)
	})

	t.Run("a named module folds its name in", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "lemmy.wasm"), wasm, 0o600))

		id, err := publishplugin.Identity(filepath.Join(dir, "lemmy.wasm"))
		require.NoError(t, err)
		require.Equal(t, md5x.FormatUUID(md5x.Digest(digest(t, wasm), "lemmy")), id)
		require.NotEqual(t, digest(t, wasm), id)
	})

	t.Run("the same bytes under the same name are stable", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "lemmy.wasm"), wasm, 0o600))

		first, err := publishplugin.Identity(filepath.Join(dir, "lemmy.wasm"))
		require.NoError(t, err)

		// a reinstall of the identical module - the id has to survive it.
		require.NoError(t, os.WriteFile(filepath.Join(dir, "lemmy.wasm"), wasm, 0o600))

		second, err := publishplugin.Identity(filepath.Join(dir, "lemmy.wasm"))
		require.NoError(t, err)
		require.Equal(t, first, second)
	})

	t.Run("different bytes under the same name differ", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "lemmy.wasm"), wasm, 0o600))

		before, err := publishplugin.Identity(filepath.Join(dir, "lemmy.wasm"))
		require.NoError(t, err)

		upgraded := append([]byte{0x00, 0x61, 0x73, 0x6D}, bytes.Repeat([]byte{0x43}, 128)...)
		require.NoError(t, os.WriteFile(filepath.Join(dir, "lemmy.wasm"), upgraded, 0o600))

		after, err := publishplugin.Identity(filepath.Join(dir, "lemmy.wasm"))
		require.NoError(t, err)
		require.NotEqual(t, before, after)
	})

	t.Run("each symlink to a shared module gets its own identity", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "lemmy.wasm"), wasm, 0o600))
		require.NoError(t, os.Symlink("lemmy.wasm", filepath.Join(dir, "lemmy-movies.wasm")))
		require.NoError(t, os.Symlink("lemmy.wasm", filepath.Join(dir, "lemmy-music.wasm")))

		lemmy, err := publishplugin.Identity(filepath.Join(dir, "lemmy.wasm"))
		require.NoError(t, err)
		movies, err := publishplugin.Identity(filepath.Join(dir, "lemmy-movies.wasm"))
		require.NoError(t, err)
		music, err := publishplugin.Identity(filepath.Join(dir, "lemmy-music.wasm"))
		require.NoError(t, err)

		require.NotEqual(t, lemmy, movies)
		require.NotEqual(t, lemmy, music)
		require.NotEqual(t, movies, music)
	})

	t.Run("a missing module errors", func(t *testing.T) {
		_, err := publishplugin.Identity(filepath.Join(t.TempDir(), "gone.wasm"))
		require.Error(t, err)
	})
}
