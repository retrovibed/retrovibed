package blockcache_test

import (
	"bytes"
	"io"
	"io/fs"
	"testing"
	"time"

	"github.com/retrovibed/retrovibed/shallows/blockcache"
	"github.com/stretchr/testify/require"
)

func TestFile(t *testing.T) {
	t.Run("Seek", func(t *testing.T) {
		t.Run("SeekStart positive offset", func(t *testing.T) {
			f := blockcache.NewFile(NewByteArrayCache(uint16(100)), time.Now(), "", 100, 0, blockcache.WithInitialIndex(50))
			newPos, err := f.Seek(25, io.SeekStart)
			require.NoError(t, err)
			require.EqualValues(t, 25, newPos)
			require.EqualValues(t, 25, blockcache.FileIndex(f))
		})

		t.Run("SeekStart zero offset", func(t *testing.T) {
			f := blockcache.NewFile(NewByteArrayCache(uint16(100)), time.Now(), "", 100, 0, blockcache.WithInitialIndex(50))
			newPos, err := f.Seek(0, io.SeekStart)
			require.NoError(t, err)
			require.EqualValues(t, 0, newPos)
			require.EqualValues(t, 0, blockcache.FileIndex(f))
		})

		t.Run("SeekStart negative offset (should error)", func(t *testing.T) {
			f := blockcache.NewFile(NewByteArrayCache(uint16(100)), time.Now(), "", 100, 0, blockcache.WithInitialIndex(50))
			newPos, err := f.Seek(-10, io.SeekStart)
			require.Error(t, err)
			require.Contains(t, err.Error(), "seek to negative index")
			require.EqualValues(t, 50, blockcache.FileIndex(f))
			require.Zero(t, newPos)
		})

		t.Run("SeekEnd positive offset (beyond logical end)", func(t *testing.T) {
			f := blockcache.NewFile(NewByteArrayCache(uint16(100)), time.Now(), "", 100, 0, blockcache.WithInitialIndex(50))
			newPos, err := f.Seek(10, io.SeekEnd)
			require.NoError(t, err)
			require.EqualValues(t, 110, newPos)
			require.EqualValues(t, 110, blockcache.FileIndex(f))
		})

		t.Run("SeekEnd - zero offset (at logical end)", func(t *testing.T) {
			f := blockcache.NewFile(NewByteArrayCache(uint16(100)), time.Now(), "", 100, 0, blockcache.WithInitialIndex(50))
			newPos, err := f.Seek(0, io.SeekEnd)
			require.NoError(t, err)
			require.EqualValues(t, 100, newPos)
			require.EqualValues(t, 100, blockcache.FileIndex(f))
		})

		t.Run("SeekEnd negative offset (before logical end)", func(t *testing.T) {
			f := blockcache.NewFile(NewByteArrayCache(uint16(100)), time.Now(), "", 100, 0, blockcache.WithInitialIndex(50))
			newPos, err := f.Seek(-20, io.SeekEnd)
			require.NoError(t, err)
			require.EqualValues(t, 80, newPos)
			require.EqualValues(t, 80, blockcache.FileIndex(f))
		})

		t.Run("SeekEnd offset causing negative result (should error)", func(t *testing.T) {
			f := blockcache.NewFile(NewByteArrayCache(uint16(100)), time.Now(), "", 100, 0, blockcache.WithInitialIndex(50))
			newPos, err := f.Seek(-101, io.SeekEnd)
			require.Error(t, err)
			require.Contains(t, err.Error(), "seek to negative index")
			require.EqualValues(t, 50, blockcache.FileIndex(f))
			require.Zero(t, newPos)
		})

		t.Run("SeekCurrent - positive offset", func(t *testing.T) {
			f := blockcache.NewFile(NewByteArrayCache(uint16(100)), time.Now(), "", 100, 0, blockcache.WithInitialIndex(50))
			newPos, err := f.Seek(20, io.SeekCurrent)
			require.NoError(t, err)
			require.EqualValues(t, 70, newPos)
			require.EqualValues(t, 70, blockcache.FileIndex(f))
		})

		t.Run("SeekCurrent - zero offset", func(t *testing.T) {
			f := blockcache.NewFile(NewByteArrayCache(uint16(100)), time.Now(), "", 100, 0, blockcache.WithInitialIndex(50))
			newPos, err := f.Seek(0, io.SeekCurrent)
			require.NoError(t, err)
			require.EqualValues(t, 50, newPos)
			require.EqualValues(t, 50, blockcache.FileIndex(f))
		})

		t.Run("SeekCurrent - negative offset", func(t *testing.T) {
			f := blockcache.NewFile(NewByteArrayCache(uint16(100)), time.Now(), "", 100, 0, blockcache.WithInitialIndex(50))
			newPos, err := f.Seek(-30, io.SeekCurrent)
			require.NoError(t, err)
			require.EqualValues(t, 20, newPos)
			require.EqualValues(t, 20, blockcache.FileIndex(f))
		})

		t.Run("SeekCurrent - offset causing negative result (should error)", func(t *testing.T) {
			f := blockcache.NewFile(NewByteArrayCache(uint16(100)), time.Now(), "", 100, 0, blockcache.WithInitialIndex(50))
			newPos, err := f.Seek(-51, io.SeekCurrent)
			require.Error(t, err)
			require.Contains(t, err.Error(), "seek to negative index")
			require.EqualValues(t, 50, blockcache.FileIndex(f))
			require.Zero(t, newPos)
		})

		t.Run("Invalid whence (should default to SeekCurrent behavior)", func(t *testing.T) {
			f := blockcache.NewFile(NewByteArrayCache(uint16(100)), time.Now(), "", 100, 0, blockcache.WithInitialIndex(50))
			newPos, err := f.Seek(15, 999)
			require.NoError(t, err)
			require.EqualValues(t, 65, newPos)
			require.EqualValues(t, 65, blockcache.FileIndex(f))
		})
	})

	t.Run("ReadAt", func(t *testing.T) {
		initialContent := []byte("Hello, World! This is some test data.")
		fileLength := uint64(len(initialContent))

		t.Run("Read full content from start", func(t *testing.T) {
			bCache := NewByteArrayCache(uint16(fileLength))
			_, err := bCache.WriteAt(initialContent, 0)
			require.NoError(t, err)
			f := blockcache.NewFile(bCache, time.Now(), "testfile.txt", fileLength, 0) // No offset, no initial index needed here
			buf := make([]byte, fileLength)
			n, err := f.ReadAt(buf, 0)
			require.NoError(t, err)
			require.Equal(t, int(fileLength), n)
			require.Equal(t, initialContent, buf)
		})

		t.Run("Read part of content from middle", func(t *testing.T) {
			bCache := NewByteArrayCache(uint16(fileLength))
			_, err := bCache.WriteAt(initialContent, 0)
			require.NoError(t, err)
			f := blockcache.NewFile(bCache, time.Now(), "testfile.txt", fileLength, 0)
			buf := make([]byte, 5)
			n, err := f.ReadAt(buf, 7)
			require.NoError(t, err)
			require.Equal(t, 5, n)
			require.Equal(t, []byte("World"), buf)
		})

		t.Run("Read beyond end of file (partial read + EOF)", func(t *testing.T) {
			bCache := NewByteArrayCache(uint16(fileLength))
			_, err := bCache.WriteAt(initialContent, 0)
			require.NoError(t, err)
			f := blockcache.NewFile(bCache, time.Now(), "testfile.txt", fileLength, 0)
			buf := make([]byte, 20)
			n, err := f.ReadAt(buf, 20)
			require.ErrorIs(t, err, io.EOF)
			require.Equal(t, int(fileLength-20), n)
			require.Equal(t, initialContent[20:], buf[:n])
		})

		t.Run("Read at exact end of file (EOF)", func(t *testing.T) {
			bCache := NewByteArrayCache(uint16(fileLength))
			_, err := bCache.WriteAt(initialContent, 0)
			require.NoError(t, err)
			f := blockcache.NewFile(bCache, time.Now(), "testfile.txt", fileLength, 0)
			buf := make([]byte, 10)
			n, err := f.ReadAt(buf, int64(fileLength))
			require.ErrorIs(t, err, io.EOF)
			require.Zero(t, n)
		})

		t.Run("Read beyond file (EOF immediately)", func(t *testing.T) {
			bCache := NewByteArrayCache(uint16(fileLength))
			_, err := bCache.WriteAt(initialContent, 0)
			require.NoError(t, err)
			f := blockcache.NewFile(bCache, time.Now(), "testfile.txt", fileLength, 0)
			buf := make([]byte, 10)
			n, err := f.ReadAt(buf, int64(fileLength+5))
			require.ErrorIs(t, err, io.EOF)
			require.Zero(t, n)
		})

		t.Run("Read with negative offset (should error)", func(t *testing.T) {
			bCache := NewByteArrayCache(uint16(fileLength))
			_, err := bCache.WriteAt(initialContent, 0)
			require.NoError(t, err)
			f := blockcache.NewFile(bCache, time.Now(), "testfile.txt", fileLength, 0)
			buf := make([]byte, 10)
			n, err := f.ReadAt(buf, -1)
			require.Error(t, err)
			require.Contains(t, err.Error(), "negative offset")
			require.Zero(t, n)
		})

		t.Run("Read an offset twice", func(t *testing.T) {
			bCache := NewByteArrayCache(uint16(fileLength))
			_, err := bCache.WriteAt(initialContent, 0)
			require.NoError(t, err)
			f := blockcache.NewFile(bCache, time.Now(), "testfile.txt", fileLength, 0)
			buf := make([]byte, 20)
			n, err := f.ReadAt(buf, 20)
			require.ErrorIs(t, err, io.EOF)
			require.Equal(t, int(fileLength-20), n)
			require.Equal(t, initialContent[20:], buf[:n])

			n, err = f.ReadAt(buf, 20)
			require.ErrorIs(t, err, io.EOF)
			require.Equal(t, int(fileLength-20), n)
			require.Equal(t, initialContent[20:], buf[:n])
		})
	})

	t.Run("WriteAt", func(t *testing.T) {
		initialCapacity := uint64(100)
		initialData := make([]byte, initialCapacity)

		t.Run("Write to start", func(t *testing.T) {
			bCache := NewByteArrayCache(uint16(initialCapacity))
			_, err := bCache.WriteAt(initialData, 0) // Populate with initialData for correct state
			require.NoError(t, err)
			f := blockcache.NewFile(bCache, time.Now(), "testfile.txt", initialCapacity, 0) // No offset, no initial index
			dataToWrite := []byte("START")
			n, err := f.WriteAt(dataToWrite, 0)
			require.NoError(t, err)
			require.Equal(t, len(dataToWrite), n)

			// Verify the content was written
			readBuf := make([]byte, len(dataToWrite))
			_, err = f.ReadAt(readBuf, 0)
			require.NoError(t, err)
			require.Equal(t, dataToWrite, readBuf)
		})

		t.Run("Write to middle", func(t *testing.T) {
			bCache := NewByteArrayCache(uint16(initialCapacity))
			_, err := bCache.WriteAt(initialData, 0)
			require.NoError(t, err)
			f := blockcache.NewFile(bCache, time.Now(), "testfile.txt", initialCapacity, 0)
			dataToWrite := []byte("MIDDLE")
			offset := int64(10)
			n, err := f.WriteAt(dataToWrite, offset)
			require.NoError(t, err)
			require.Equal(t, len(dataToWrite), n)

			readBuf := make([]byte, len(dataToWrite))
			_, err = f.ReadAt(readBuf, offset)
			require.NoError(t, err)
			require.Equal(t, dataToWrite, readBuf)
		})

		t.Run("Write exactly to end (fill capacity)", func(t *testing.T) {
			bCache := NewByteArrayCache(uint16(initialCapacity))
			_, err := bCache.WriteAt(initialData, 0)
			require.NoError(t, err)
			f := blockcache.NewFile(bCache, time.Now(), "testfile.txt", initialCapacity, 0)
			dataToWrite := bytes.Repeat([]byte("A"), int(initialCapacity))
			n, err := f.WriteAt(dataToWrite, 0)
			require.NoError(t, err)
			require.Equal(t, int(initialCapacity), n)
		})

		t.Run("Write beyond capacity (should error)", func(t *testing.T) {
			bCache := NewByteArrayCache(uint16(initialCapacity))
			_, err := bCache.WriteAt(initialData, 0)
			require.NoError(t, err)
			f := blockcache.NewFile(bCache, time.Now(), "testfile.txt", initialCapacity, 0)
			dataToWrite := []byte("TOO MUCH")
			offset := int64(initialCapacity - 2)
			n, err := f.WriteAt(dataToWrite, offset)
			require.Error(t, err)
			require.Contains(t, err.Error(), "write beyond cache capacity")
			require.Zero(t, n)
		})

		t.Run("Write with non-zero Offset writes to correct cache position", func(t *testing.T) {
			// cache is split into two logical files: bytes 0-9 (fileA) and bytes 10-19 (fileB)
			cache := NewByteArrayCache(20)
			fileA := blockcache.NewFile(cache, time.Now(), "file_a.bin", 10, 0600)
			fileB := blockcache.NewFile(cache, time.Now(), "file_b.bin", 10, 0600, blockcache.WithInitialOffset(10))

			data := []byte("HELLO")
			_, err := fileB.WriteAt(data, 0)
			require.NoError(t, err)

			// read back through fileB — must see the written data
			buf := make([]byte, 5)
			_, err = fileB.ReadAt(buf, 0)
			require.NoError(t, err)
			require.Equal(t, data, buf)

			// fileA's region must be untouched
			bufA := make([]byte, 10)
			_, err = fileA.ReadAt(bufA, 0)
			require.NoError(t, err)
			require.Equal(t, make([]byte, 10), bufA)
		})

		t.Run("Write with negative offset (should error)", func(t *testing.T) {
			bCache := NewByteArrayCache(uint16(initialCapacity))
			_, err := bCache.WriteAt(initialData, 0)
			require.NoError(t, err)
			f := blockcache.NewFile(bCache, time.Now(), "testfile.txt", initialCapacity, 0)
			dataToWrite := []byte("TEST")
			n, err := f.WriteAt(dataToWrite, -1)
			require.Error(t, err)
			require.Contains(t, err.Error(), "negative offset")
			require.Zero(t, n)
		})
	})

	t.Run("Read", func(t *testing.T) {
		initialContent := []byte("ABCDEFGHIJ")
		fileLength := uint64(len(initialContent))

		t.Run("Read from current index", func(t *testing.T) {
			bCache := NewByteArrayCache(uint16(fileLength))
			_, err := bCache.WriteAt(initialContent, 0)
			require.NoError(t, err)
			f := blockcache.NewFile(bCache, time.Now(), "testfile.txt", fileLength, 0, blockcache.WithInitialIndex(3))
			buf := make([]byte, 4)
			n, err := f.Read(buf)
			require.NoError(t, err)
			require.Equal(t, 4, n)
			require.Equal(t, []byte("DEFG"), buf)
			require.EqualValues(t, 3+4, blockcache.FileIndex(f))
		})

		t.Run("Read to EOF", func(t *testing.T) {
			bCache := NewByteArrayCache(uint16(fileLength))
			_, err := bCache.WriteAt(initialContent, 0)
			require.NoError(t, err)
			f := blockcache.NewFile(bCache, time.Now(), "testfile.txt", fileLength, 0, blockcache.WithInitialIndex(7))
			buf := make([]byte, 5)
			n, err := f.Read(buf)
			require.ErrorIs(t, err, io.EOF)
			require.Equal(t, 3, n)
			require.Equal(t, []byte("HIJ\x00\x00"), buf)
			require.EqualValues(t, 7+3, blockcache.FileIndex(f))
		})

		t.Run("Read from EOF", func(t *testing.T) {
			bCache := NewByteArrayCache(uint16(fileLength))
			_, err := bCache.WriteAt(initialContent, 0)
			require.NoError(t, err)
			f := blockcache.NewFile(bCache, time.Now(), "testfile.txt", fileLength, 0, blockcache.WithInitialIndex(fileLength))
			buf := make([]byte, 5)
			n, err := f.Read(buf)
			require.ErrorIs(t, err, io.EOF)
			require.Zero(t, n)
			require.EqualValues(t, fileLength, blockcache.FileIndex(f))
		})
	})

	t.Run("Stat", func(t *testing.T) {
		name := "example.txt"
		length := uint64(500)
		modTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
		mode := fs.FileMode(0644)
		f := blockcache.NewFile(NewByteArrayCache(uint16(length)), modTime, name, length, mode)

		info, err := f.Stat()
		require.NoError(t, err)
		require.NotNil(t, info)
		require.Equal(t, name, info.Name())
		require.EqualValues(t, length, info.Size())
		require.Equal(t, modTime, info.ModTime())
		require.Equal(t, mode, info.Mode())
		require.False(t, info.IsDir())
		require.Nil(t, info.Sys())
	})

	t.Run("Info", func(t *testing.T) {
		name := "example-dir"
		modTime := time.Date(2023, 2, 1, 10, 0, 0, 0, time.UTC)
		mode := fs.ModeDir | (0755 & fs.ModePerm)
		f := blockcache.NewFile(NewByteArrayCache(uint16(0)), modTime, name, 0, mode)

		dirEntryInfo, err := f.Info()
		require.NoError(t, err)
		require.NotNil(t, dirEntryInfo)
		require.Equal(t, name, dirEntryInfo.Name())
		require.True(t, dirEntryInfo.IsDir())
		require.EqualValues(t, 0, dirEntryInfo.Size())
		require.Equal(t, modTime, dirEntryInfo.ModTime())
		require.Equal(t, mode, dirEntryInfo.Mode())
		require.Nil(t, dirEntryInfo.Sys())
	})

	t.Run("IsDir_DirEntry", func(t *testing.T) {
		fDir := blockcache.NewFile(NewByteArrayCache(uint16(0)), time.Now(), "test_dir", 0, fs.ModeDir|0755)
		fFile := blockcache.NewFile(NewByteArrayCache(uint16(10)), time.Now(), "test_file", 10, 0644)

		require.True(t, fDir.IsDir())
		require.False(t, fFile.IsDir())
	})

	t.Run("Name_DirEntry", func(t *testing.T) {
		f := blockcache.NewFile(NewByteArrayCache(uint16(0)), time.Now(), "/path/to/my/file.txt", 10, 0644)
		require.Equal(t, "file.txt", f.Name())

		fDir := blockcache.NewFile(NewByteArrayCache(uint16(0)), time.Now(), "/path/to/my/directory/", 0, fs.ModeDir|0755)
		require.Equal(t, "directory", fDir.Name())
	})

	t.Run("Type_DirEntry", func(t *testing.T) {
		f := blockcache.NewFile(NewByteArrayCache(uint16(10)), time.Now(), "file.txt", 10, 0644)
		require.Equal(t, fs.FileMode(0), f.Type())

		fDir := blockcache.NewFile(NewByteArrayCache(uint16(0)), time.Now(), "dir", 0, fs.ModeDir|0755)
		require.Equal(t, fs.ModeDir, fDir.Type())

		fSymlink := blockcache.NewFile(NewByteArrayCache(uint16(0)), time.Now(), "symlink", 0, fs.ModeSymlink|0777)
		require.Equal(t, fs.ModeSymlink, fSymlink.Type())
	})

	t.Run("Close", func(t *testing.T) {
		f := blockcache.NewFile(NewByteArrayCache(uint16(10)), time.Now(), "closable.txt", 10, 0644)
		err := f.Close()
		require.NoError(t, err)
	})
}
