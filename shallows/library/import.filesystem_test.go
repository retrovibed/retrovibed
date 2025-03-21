package library_test

import (
	"path/filepath"
	"testing"

	"github.com/retrovibed/retrovibed/internal/fsx"
	"github.com/retrovibed/retrovibed/internal/md5x"
	"github.com/retrovibed/retrovibed/internal/testx"
	"github.com/retrovibed/retrovibed/library"
	"github.com/stretchr/testify/require"
)

func TestImportFilesystemDryRun(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()

	count := 0
	for _, err := range library.ImportFilesystem(ctx, library.ImportFileDryRun, testx.Fixture("tree.example.1")) {
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
	for tx, err := range library.ImportFilesystem(ctx, library.ImportCopyFile(vfs), testx.Fixture("tree.example.1")) {
		require.NoError(t, err)
		require.Equal(t, testx.ReadMD5(tx.Path), testx.ReadMD5(filepath.Join(tmpdir, md5x.FormatString(tx.MD5))))
		count++
	}

	require.Equal(t, 2, count)
}

func TestImportFilesystemSymlink(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()

	tmpdir := t.TempDir()
	vfs := fsx.DirVirtual(tmpdir)
	count := 0
	for tx, err := range library.ImportFilesystem(ctx, library.ImportSymlinkFile(vfs), testx.Must(filepath.Abs(testx.Fixture("tree.example.1")))(t)) {
		require.NoError(t, err)
		require.Equal(t, testx.ReadMD5(tx.Path), testx.ReadMD5(filepath.Join(tmpdir, md5x.FormatString(tx.MD5))))
		count++
	}

	require.Equal(t, 2, count)
}
