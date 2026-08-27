package tracking

import (
	"io/fs"
	"path/filepath"

	"github.com/RoaringBitmap/roaring/v2"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

const (
	BitmapSuffix  = ".bitmap"
	TorrentSuffix = ".torrent"
)

// InfoFromPath reads the torrent info artifact stored on disk at path.
func InfoFromPath(path string) (info *metainfo.Info, err error) {
	mi, err := metainfo.LoadFromFile(path + TorrentSuffix)
	if err != nil {
		return nil, err
	}

	i, err := mi.UnmarshalInfo()
	if err != nil {
		return nil, err
	}

	return &i, nil
}

// FileInfoFromOffset reads the torrent info artifact stored on disk at path and
// returns the file entry whose offset matches zerooffset.
func FileInfoFromOffset(path string, zerooffset uint64) (_z metainfo.File, err error) {
	info, err := InfoFromPath(path)
	if err != nil {
		return _z, err
	}

	for n := range metainfo.Files(info) {
		if n.Offset == zerooffset {
			return n, nil
		}
	}

	return _z, fs.ErrNotExist
}

// BitmapFromPath reads the bitmap artifact stored on disk at path.
func BitmapFromPath(path string) (_z *roaring.Bitmap, err error) {
	return torrent.NewBitmapCache(filepath.Dir(path)).Read(errorsx.Zero(int160.FromHexEncodedString(filepath.Base(path))))
}
