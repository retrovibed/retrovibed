package blockcache

import (
	"sync/atomic"

	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/james-lawrence/torrent/storage"

	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
)

type OptionTorrentCacheStorage func(*TorrentCacheStorage)

func OptionTorrentCacheStorageBlockLength(n uint64) OptionTorrentCacheStorage {
	return func(tcs *TorrentCacheStorage) {
		tcs.blength = n
	}
}

func NewTorrentFromVirtualFS(v fsx.Virtual, options ...OptionTorrentCacheStorage) *TorrentCacheStorage {
	return langx.Autoptr(langx.Clone(TorrentCacheStorage{v: v, blength: DefaultBlockLength}, options...))
}

var _ storage.ClientImpl = &TorrentCacheStorage{}

type TorrentCacheStorage struct {
	v       fsx.Virtual
	blength uint64
}

func (t *TorrentCacheStorage) OpenTorrent(info *metainfo.Info, infoHash int160.T) (storage.TorrentImpl, error) {
	dir := storage.InfoHashPathMaker(t.v.Path(), infoHash, info, nil)

	cache, err := NewDirectoryCache(dir, OptionDirCacheBlockLength(int64(t.blength)))
	if err != nil {
		return nil, err
	}

	return &fileTorrentImpl{
		cache: cache,
	}, nil
}

func (t *TorrentCacheStorage) Close() error {
	return nil
}

type fileTorrentImpl struct {
	closed atomic.Bool
	cache  *DirCache
}

func (t *fileTorrentImpl) ReadAt(p []byte, off int64) (n int, err error) {
	if t.closed.Load() {
		return 0, storage.ErrClosed()
	}
	return t.cache.ReadAt(p, off)
}

func (t *fileTorrentImpl) WriteAt(p []byte, off int64) (n int, err error) {
	if t.closed.Load() {
		return 0, storage.ErrClosed()
	}

	return t.cache.WriteAt(p, off)
}

func (t *fileTorrentImpl) Close() error {
	t.closed.Store(true)
	return nil
}
