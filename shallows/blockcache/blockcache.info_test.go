package blockcache_test

import (
	"io/fs"
	"testing"
	"time"

	"github.com/retrovibed/retrovibed/blockcache"
	"github.com/stretchr/testify/require"
)

func TestInfo(t *testing.T) {
	testInfo := func(m fs.FileMode, len uint16, ts time.Time, name string) blockcache.Info {
		return blockcache.Info{File: blockcache.NewFile(NewByteArrayCache(len), ts, name, uint64(len), m)}
	}

	t.Run("Directory Behavior", func(t *testing.T) {
		inf := testInfo(fs.ModeDir|(0700&fs.ModePerm), 0, time.Date(2024, time.January, 1, 10, 0, 0, 0, time.UTC), "d")

		t.Run("Name", func(t *testing.T) {
			require.Equal(t, "d", inf.Name())
		})

		t.Run("IsDir", func(t *testing.T) {
			require.True(t, inf.IsDir())
		})

		t.Run("Size", func(t *testing.T) {
			require.EqualValues(t, 0, inf.Size())
		})

		t.Run("Mode", func(t *testing.T) {
			require.Equal(t, fs.ModeDir|(0700&fs.ModePerm), inf.Mode())
			require.True(t, inf.Mode().IsDir())
			require.False(t, inf.Mode().IsRegular())
			require.False(t, inf.Mode()&fs.ModeSymlink != 0)
		})

		t.Run("ModTime", func(t *testing.T) {
			require.Equal(t, time.Date(2024, time.January, 1, 10, 0, 0, 0, time.UTC), inf.ModTime())
		})

		t.Run("Sys", func(t *testing.T) {
			require.Nil(t, inf.Sys())
		})
	})

	t.Run("Regular File Behavior", func(t *testing.T) {
		inf := testInfo(0600&fs.ModePerm, 12345, time.Date(2024, time.February, 15, 14, 30, 0, 0, time.UTC), "reg")

		t.Run("Name", func(t *testing.T) {
			require.Equal(t, "reg", inf.Name())
		})

		t.Run("IsDir", func(t *testing.T) {
			require.False(t, inf.IsDir())
		})

		t.Run("Size", func(t *testing.T) {
			require.EqualValues(t, 12345, inf.Size())
		})

		t.Run("Mode", func(t *testing.T) {
			require.Equal(t, 0600&fs.ModePerm, inf.Mode())
			require.True(t, inf.Mode().IsRegular())
			require.False(t, inf.Mode().IsDir())
			require.False(t, inf.Mode()&fs.ModeSymlink != 0)
		})

		t.Run("ModTime", func(t *testing.T) {
			require.Equal(t, time.Date(2024, time.February, 15, 14, 30, 0, 0, time.UTC), inf.ModTime())
		})

		t.Run("Sys", func(t *testing.T) {
			require.Nil(t, inf.Sys())
		})
	})

	t.Run("Symlink Behavior", func(t *testing.T) {
		inf := testInfo(fs.ModeSymlink|(0777&fs.ModePerm), 20, time.Date(2024, time.March, 1, 11, 0, 0, 0, time.UTC), "sym")

		t.Run("Name", func(t *testing.T) {
			require.Equal(t, "sym", inf.Name())
		})

		t.Run("IsDir", func(t *testing.T) {
			require.False(t, inf.IsDir())
		})

		t.Run("Size", func(t *testing.T) {
			require.EqualValues(t, 20, inf.Size())
		})

		t.Run("Mode", func(t *testing.T) {
			m := inf.Mode()
			require.True(t, m&fs.ModeSymlink != 0)
			require.False(t, m.IsRegular())
			require.False(t, m.IsDir())
		})

		t.Run("ModTime", func(t *testing.T) {
			require.Equal(t, time.Date(2024, time.March, 1, 11, 0, 0, 0, time.UTC), inf.ModTime())
		})

		t.Run("Sys", func(t *testing.T) {
			require.Nil(t, inf.Sys())
		})
	})
}
