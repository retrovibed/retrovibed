package tracking_test

import (
	"crypto/rand"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/james-lawrence/torrent/metainfo"
	"github.com/james-lawrence/torrent/torrenttest"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestInfoFromPath(t *testing.T) {
	t.Run("round trips a multi file torrent", func(t *testing.T) {
		dir := t.TempDir()
		info := testx.Must(torrenttest.Tree(
			dir, rand.Reader, 16*bytesx.KiB, 64*bytesx.KiB,
			[]string{"file1.mkv", "file2.mkv", "file3.mkv"},
		))(t)

		p := filepath.Join(t.TempDir(), info.Name)
		md := metainfo.MetaInfo{InfoBytes: testx.Must(metainfo.Encode(info))(t)}
		require.NoError(t, os.WriteFile(p+tracking.TorrentSuffix, testx.Must(metainfo.Encode(md))(t), 0600))

		actual, err := tracking.InfoFromPath(p)
		require.NoError(t, err)
		require.Equal(t, info.Name, actual.Name)
		require.Equal(t, info.TotalLength(), actual.TotalLength())
		require.Len(t, actual.Files, len(info.Files))
		for i, fi := range info.Files {
			require.Equal(t, fi.Path, actual.Files[i].Path)
			require.Equal(t, fi.Length, actual.Files[i].Length)
		}
	})

	t.Run("round trips a single file torrent", func(t *testing.T) {
		seedir := t.TempDir()
		info, _, err := torrenttest.Random(seedir, 64*bytesx.KiB)
		require.NoError(t, err)

		p := filepath.Join(t.TempDir(), info.Name)
		md := metainfo.MetaInfo{InfoBytes: testx.Must(metainfo.Encode(info))(t)}
		require.NoError(t, os.WriteFile(p+tracking.TorrentSuffix, testx.Must(metainfo.Encode(md))(t), 0600))

		actual, err := tracking.InfoFromPath(p)
		require.NoError(t, err)
		require.False(t, actual.IsDir())
		require.Equal(t, info.Name, actual.Name)
		require.Equal(t, info.TotalLength(), actual.TotalLength())
	})

	t.Run("missing torrent file returns an error", func(t *testing.T) {
		_, err := tracking.InfoFromPath(filepath.Join(t.TempDir(), "does-not-exist"))
		require.Error(t, err)
		require.True(t, errors.Is(err, fs.ErrNotExist))
	})

	t.Run("malformed torrent file returns a decode error", func(t *testing.T) {
		p := filepath.Join(t.TempDir(), "garbage")
		require.NoError(t, os.WriteFile(p+tracking.TorrentSuffix, []byte("not bencode"), 0600))

		_, err := tracking.InfoFromPath(p)
		require.Error(t, err)
	})
}

func TestFileInfoFromOffset(t *testing.T) {
	t.Run("returns the first file at offset zero", func(t *testing.T) {
		dir := t.TempDir()
		info := testx.Must(torrenttest.Tree(
			dir, rand.Reader, 16*bytesx.KiB, 64*bytesx.KiB,
			[]string{"file1.mkv", "file2.mkv", "file3.mkv"},
		))(t)

		p := filepath.Join(t.TempDir(), info.Name)
		md := metainfo.MetaInfo{InfoBytes: testx.Must(metainfo.Encode(info))(t)}
		require.NoError(t, os.WriteFile(p+tracking.TorrentSuffix, testx.Must(metainfo.Encode(md))(t), 0600))

		var files []metainfo.File
		for f := range metainfo.Files(info) {
			files = append(files, f)
		}
		require.Len(t, files, 3)

		actual, err := tracking.FileInfoFromOffset(p, 0)
		require.NoError(t, err)
		require.Equal(t, files[0], actual)
	})

	t.Run("returns a later file by its real offset", func(t *testing.T) {
		dir := t.TempDir()
		info := testx.Must(torrenttest.Tree(
			dir, rand.Reader, 16*bytesx.KiB, 64*bytesx.KiB,
			[]string{"file1.mkv", "file2.mkv", "file3.mkv"},
		))(t)

		p := filepath.Join(t.TempDir(), info.Name)
		md := metainfo.MetaInfo{InfoBytes: testx.Must(metainfo.Encode(info))(t)}
		require.NoError(t, os.WriteFile(p+tracking.TorrentSuffix, testx.Must(metainfo.Encode(md))(t), 0600))

		var files []metainfo.File
		for f := range metainfo.Files(info) {
			files = append(files, f)
		}
		require.Len(t, files, 3)
		want := files[2]
		require.Greater(t, want.Offset, uint64(0))

		actual, err := tracking.FileInfoFromOffset(p, want.Offset)
		require.NoError(t, err)
		require.Equal(t, want, actual)
	})

	t.Run("offset not aligned to a file boundary returns not exist", func(t *testing.T) {
		dir := t.TempDir()
		info := testx.Must(torrenttest.Tree(
			dir, rand.Reader, 16*bytesx.KiB, 64*bytesx.KiB,
			[]string{"file1.mkv", "file2.mkv", "file3.mkv"},
		))(t)

		p := filepath.Join(t.TempDir(), info.Name)
		md := metainfo.MetaInfo{InfoBytes: testx.Must(metainfo.Encode(info))(t)}
		require.NoError(t, os.WriteFile(p+tracking.TorrentSuffix, testx.Must(metainfo.Encode(md))(t), 0600))

		_, err := tracking.FileInfoFromOffset(p, uint64(info.TotalLength())+1)
		require.Error(t, err)
		require.True(t, errors.Is(err, fs.ErrNotExist))
	})

	t.Run("single file torrent resolves via info name at offset zero", func(t *testing.T) {
		seedir := t.TempDir()
		info, _, err := torrenttest.Random(seedir, 64*bytesx.KiB)
		require.NoError(t, err)

		p := filepath.Join(t.TempDir(), info.Name)
		md := metainfo.MetaInfo{InfoBytes: testx.Must(metainfo.Encode(info))(t)}
		require.NoError(t, os.WriteFile(p+tracking.TorrentSuffix, testx.Must(metainfo.Encode(md))(t), 0600))

		actual, err := tracking.FileInfoFromOffset(p, 0)
		require.NoError(t, err)
		require.Equal(t, info.Name, actual.Path)
		require.EqualValues(t, info.TotalLength(), actual.Length)
	})
}
