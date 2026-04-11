package asyncx_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/retrovibed/retrovibed/internal/asyncx"
	"github.com/stretchr/testify/require"
)

func TestWatchFiles(t *testing.T) {
	t.Run("watch files that don't exist", func(t *testing.T) {
		tempDir := t.TempDir()

		wakeup := asyncx.NewWakeup(t.Context())
		defer wakeup.Close()

		file1 := filepath.Join(tempDir, "file1.txt")
		file2 := filepath.Join(tempDir, "file2.txt")

		err := asyncx.WatchFiles(t.Context(), wakeup, func(e fsnotify.Event) bool { return true }, file1, file2)
		require.NoError(t, err)

		_, err = os.Stat(file1)
		require.NoError(t, err)
		_, err = os.Stat(file2)
		require.NoError(t, err)

		go func() {
			time.Sleep(100 * time.Millisecond)
			err := os.WriteFile(file1, []byte("test"), 0600)
			require.NoError(t, err)
		}()

		select {
		case <-wakeup.C:
		case <-t.Context().Done():
			require.Fail(t, "wakeup was not triggered")
		}
	})

	t.Run("watch existing files", func(t *testing.T) {
		tempDir := t.TempDir()

		file1 := filepath.Join(tempDir, "file1.txt")
		file2 := filepath.Join(tempDir, "file2.txt")
		err := os.WriteFile(file1, []byte("test1"), 0600)
		require.NoError(t, err)
		err = os.WriteFile(file2, []byte("test2"), 0600)
		require.NoError(t, err)

		wakeup := asyncx.NewWakeup(t.Context())
		defer wakeup.Close()

		err = asyncx.WatchFiles(t.Context(), wakeup, func(e fsnotify.Event) bool { return true }, file1, file2)
		require.NoError(t, err)

		go func() {
			time.Sleep(100 * time.Millisecond)
			err := os.WriteFile(file1, []byte("test1_updated"), 0600)
			require.NoError(t, err)
		}()

		select {
		case <-wakeup.C:
		case <-t.Context().Done():
			require.Fail(t, "wakeup was not triggered")
		}
	})
}

func TestWatchDirectories(t *testing.T) {
	t.Run("watch directories that don't exist", func(t *testing.T) {
		tempDir := t.TempDir()

		wakeup := asyncx.NewWakeup(t.Context())
		defer wakeup.Close()

		dir1 := filepath.Join(tempDir, "dir1")
		dir2 := filepath.Join(tempDir, "dir2")

		err := asyncx.WatchDirectories(t.Context(), wakeup, func(e fsnotify.Event) bool { return true }, dir1, dir2)
		require.NoError(t, err)

		_, err = os.Stat(dir1)
		require.NoError(t, err)
		_, err = os.Stat(dir2)
		require.NoError(t, err)

		go func() {
			time.Sleep(100 * time.Millisecond)
			fileInDir1 := filepath.Join(dir1, "test.txt")
			err := os.WriteFile(fileInDir1, []byte("test"), 0600)
			require.NoError(t, err)
		}()

		select {
		case <-wakeup.C:
		case <-t.Context().Done():
			require.Fail(t, "wakeup was not triggered")
		}
	})

	t.Run("watch existing directories", func(t *testing.T) {
		tempDir := t.TempDir()

		dir1 := filepath.Join(tempDir, "dir1")
		dir2 := filepath.Join(tempDir, "dir2")
		err := os.MkdirAll(dir1, 0700)
		require.NoError(t, err)
		err = os.MkdirAll(dir2, 0700)
		require.NoError(t, err)

		wakeup := asyncx.NewWakeup(t.Context())
		defer wakeup.Close()

		err = asyncx.WatchDirectories(t.Context(), wakeup, func(e fsnotify.Event) bool { return true }, dir1, dir2)
		require.NoError(t, err)

		go func() {
			time.Sleep(100 * time.Millisecond)
			fileInDir1 := filepath.Join(dir1, "test.txt")
			err := os.WriteFile(fileInDir1, []byte("test"), 0600)
			require.NoError(t, err)
		}()

		select {
		case <-wakeup.C:
		case <-t.Context().Done():
			require.Fail(t, "wakeup was not triggered")
		}
	})
}
