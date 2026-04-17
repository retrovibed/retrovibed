package library

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"math/rand/v2"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/james-lawrence/torrent/storage"
	"github.com/retrovibed/retrovibed/retroapi/blockcache"
	"github.com/retrovibed/retrovibed/retroapi/deeppool"
	"github.com/retrovibed/retrovibed/shallows/internal/cryptox"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/iox"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
)

const (
	ErrCacheUnrecoverable = errorsx.String("cache layer blocked due to previous error")
)

func NewTorrentStorageFromHTTP(c *http.Client, q sqlx.Queryer, d storage.ClientImpl) *TorrentCacheStorage {
	return NewTorrentStorage(deeppool.NewRanger(c), q, d, blockcache.DefaultBlockLength)
}

func NewTorrentStorage(dl downloader, q sqlx.Queryer, d storage.ClientImpl, bsize uint64) *TorrentCacheStorage {
	if (bsize % 64) != 0 {
		// need to ensure out block size is a multiple of 64 for chacha encryption random access to work correctly.
		panic("attempted to create storage with a block size that is not a multiple of 64")
	}

	return &TorrentCacheStorage{
		d:           d,
		q:           q,
		dl:          dl,
		m:           &sync.Mutex{},
		blocklength: bsize,
	}
}

type TorrentCacheStorage struct {
	d           storage.ClientImpl
	q           sqlx.Queryer
	dl          downloader
	m           *sync.Mutex
	blocklength uint64
}

func (t *TorrentCacheStorage) OpenTorrent(info *metainfo.Info, id int160.T) (storage.TorrentImpl, error) {
	impl, err := t.d.OpenTorrent(info, id)
	if err != nil {
		return nil, err
	}

	return &fallback{
		id:          id,
		TorrentImpl: impl,
		q:           t.q,
		dl:          t.dl,
		m:           t.m,
		blocklength: t.blocklength,
	}, nil
}

func (t *TorrentCacheStorage) Close() error {
	return t.d.Close()
}

type fallback struct {
	dl downloader
	storage.TorrentImpl
	id           int160.T
	q            sqlx.Queryer
	m            *sync.Mutex
	cacheblocked atomic.Bool
	blocklength  uint64
}

func (t *fallback) downloadChunk(prng *rand.ChaCha8, id string, offset uint64, length uint64) error {
	// download a block length at a time.
	doffset, dlength := calculateBlockRange(t.blocklength, offset, length)

	// log.Println("------------------------------ download initiated", id, t.blocklength, offset, length, offset%64, offset%t.blocklength, "->", doffset, doffset+dlength)
	// defer log.Println("------------------------------ download completed", id, doffset, doffset+dlength)

	w, err := cryptox.NewOffsetWriterChaCha20(prng, io.NewOffsetWriter(t.TorrentImpl, int64(doffset)), uint32(doffset))
	if err != nil {
		return err
	}
	ctx, done := context.WithCancel(context.Background())
	defer done()

	return t.dl.Download(ctx, id, doffset, doffset+dlength, iox.NewTimeoutWriter(done, 3*time.Second, w))
}

func (t *fallback) ReadAt(p []byte, off int64) (n int, err error) {
	// defer log.Printf("ReadAt completed %3d %3d", off, off+int64(len(p)))
	// defer log.Printf("ReadAt completed %3d %3d %v", off, off+int64(len(p)), p)

	if n, err = t.TorrentImpl.ReadAt(p, off); err == nil {
		return n, nil
	} else if errorsx.Ignore(err, fs.ErrNotExist) != nil {
		return n, err
	}

	if t.cacheblocked.Load() {
		return n, errorsx.NewUnrecoverable(ErrCacheUnrecoverable)
	}

	t.m.Lock()
	defer t.m.Unlock()

	if n, err = t.TorrentImpl.ReadAt(p, off); err == nil {
		return n, err
	}

	ctx, done := context.WithTimeout(context.Background(), 3*time.Second)
	defer done()

	q := sqlx.Scan(MetadataForTorrentArchiveRetrieval(ctx, t.q, t.id.Bytes(), uint64(off), uint64(len(p))))
	for v := range q.Iter() {
		switch v.ArchiveID {
		case uuid.Nil.String():
			fallthrough
		case uuid.Max.String():
			continue
		default:
			// non min/max uuid.
		}

		if cerr := t.downloadChunk(MetadataChaCha8(v), v.ArchiveID, uint64(off), v.Bytes); cerr != nil {
			var (
				unrecoverable errorsx.Unrecoverable
			)

			if errors.As(cerr, &unrecoverable) {
				t.cacheblocked.Store(true)
			}

			return n, cerr
		}
	}

	if qerr := q.Err(); qerr != nil {
		return n, err
	}

	// re-read from disk. at this point we either succeeded or failed at resync from the archive.
	return t.TorrentImpl.ReadAt(p, off)
}

func (t *fallback) Close() error {
	return t.TorrentImpl.Close()
}
