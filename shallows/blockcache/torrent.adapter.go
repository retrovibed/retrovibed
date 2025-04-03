package blockcache

// revisit once torrent lib is improved.
// import (
// 	"io"
// 	"os"
// 	"path/filepath"

// 	"github.com/james-lawrence/torrent/metainfo"
// 	"github.com/james-lawrence/torrent/storage"

// 	"github.com/retrovibed/retrovibed/internal/fsx"
// 	"github.com/retrovibed/retrovibed/internal/iox"
// )

// func NewTorrentFromVirtualFS(v fsx.Virtual) *TorrentCacheStorage {
// 	return &TorrentCacheStorage{v: v}
// }

// var _ storage.ClientImpl = &TorrentCacheStorage{}

// type TorrentCacheStorage struct {
// 	v fsx.Virtual
// }

// func (t *TorrentCacheStorage) OpenTorrent(info *metainfo.Info, infoHash metainfo.Hash) (storage.TorrentImpl, error) {
// 	path := filepath.Join("tmp", infoHash.HexString())
// 	if err := t.v.MkDirAll(filepath.Dir(path), 0700); err != nil {
// 		return nil, err
// 	}

// 	osio, err := t.v.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0600)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// initialize the file to all zeros
// 	n, err := io.Copy(osio, iox.ZeroReader{})
// 	if err != nil {
// 		return nil, err
// 	}

// 	if n < info.TotalLength() {
// 		return nil, io.ErrShortWrite
// 	}

// 	return &fileTorrentImpl{
// 		dir:      path,
// 		info:     info,
// 		infoHash: infoHash,
// 	}, nil
// }

// func (t *TorrentCacheStorage) Close() error {
// 	return nil
// }

// type fileTorrentImpl struct {
// 	dst      *os.File
// 	dir      string
// 	info     *metainfo.Info
// 	infoHash metainfo.Hash
// }

// func (t *fileTorrentImpl) Close() error {
// 	return t.dst.Close()
// }
