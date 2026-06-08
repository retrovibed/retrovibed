package library_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestImportFilesystemDryRun(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()

	count := 0
	for _, err := range library.ImportFilesystem(ctx, library.ImportFileDryRun, os.DirFS(testx.Fixture()).(fs.StatFS), "tree.example.1") {
		require.NoError(t, err)
		count++
	}

	require.Equal(t, 2, count)
}

func TestImportFilesystemCopy(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()

	tmpdir := t.TempDir()
	vfs := fsx.DirVirtual(tmpdir)
	count := 0
	for tx, err := range library.ImportFilesystem(ctx, library.ImportCopyFile(vfs), os.DirFS(testx.Fixture()).(fs.StatFS), "tree.example.1") {
		require.NoError(t, err)
		require.Equal(t, testx.ReadMD5(testx.Fixture(tx.Path)), testx.ReadMD5(filepath.Join(tmpdir, md5x.FormatUUID(tx.MD5))))
		count++
	}

	require.Equal(t, 2, count)
}

func TestImportFilesystemSymlink(t *testing.T) {
	t.Run("creates symlinks named by md5 uuid with matching content", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		tmpdir := t.TempDir()
		vfs := fsx.DirVirtual(tmpdir)
		fixvfs := fsx.DirVirtual(testx.Must(filepath.Abs(testx.Fixture()))(t))
		count := 0

		for tx, err := range library.ImportFilesystem(ctx, library.ImportSymlinkFile(fixvfs, vfs), os.DirFS(testx.Fixture()).(fs.StatFS), "tree.example.1") {
			require.NoError(t, err)
			require.Equal(t, testx.ReadMD5(testx.Fixture(tx.Path)), testx.ReadMD5(filepath.Join(tmpdir, md5x.FormatUUID(tx.MD5))))
			count++
		}

		require.Equal(t, 2, count)
	})

	t.Run("symlink points to original file location", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		tmpdir := t.TempDir()
		vfs := fsx.DirVirtual(tmpdir)
		absfix := testx.Must(filepath.Abs(testx.Fixture()))(t)
		fixvfs := fsx.DirVirtual(absfix)

		for tx, err := range library.ImportFilesystem(ctx, library.ImportSymlinkFile(fixvfs, vfs), os.DirFS(testx.Fixture()).(fs.StatFS), "tree.example.1") {
			require.NoError(t, err)
			target, lerr := os.Readlink(filepath.Join(tmpdir, md5x.FormatUUID(tx.MD5)))
			require.NoError(t, lerr)
			require.Equal(t, filepath.Join(absfix, tx.Path), target)
		}
	})

	t.Run("idempotent when re-importing the same files", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		tmpdir := t.TempDir()
		vfs := fsx.DirVirtual(tmpdir)
		fixvfs := fsx.DirVirtual(testx.Must(filepath.Abs(testx.Fixture()))(t))
		src := os.DirFS(testx.Fixture()).(fs.StatFS)

		count := 0
		for tx, err := range library.ImportFilesystem(ctx, library.ImportSymlinkFile(fixvfs, vfs), src, "tree.example.1") {
			require.NoError(t, err)
			count++
			_ = tx
		}
		require.Equal(t, 2, count)

		count = 0
		for tx, err := range library.ImportFilesystem(ctx, library.ImportSymlinkFile(fixvfs, vfs), src, "tree.example.1") {
			require.NoError(t, err)
			count++
			_ = tx
		}
		require.Equal(t, 2, count)
	})

	t.Run("single file import", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		tmpdir := t.TempDir()
		vfs := fsx.DirVirtual(tmpdir)
		fixvfs := fsx.DirVirtual(testx.Must(filepath.Abs(testx.Fixture()))(t))
		count := 0

		for tx, err := range library.ImportFilesystem(ctx, library.ImportSymlinkFile(fixvfs, vfs), os.DirFS(testx.Fixture()).(fs.StatFS), "tree.example.1/1.txt") {
			require.NoError(t, err)
			require.Equal(t, testx.ReadMD5(testx.Fixture(tx.Path)), testx.ReadMD5(filepath.Join(tmpdir, md5x.FormatUUID(tx.MD5))))
			count++
		}

		require.Equal(t, 1, count)
	})
}
