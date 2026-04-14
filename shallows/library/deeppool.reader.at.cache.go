package library

import (
	"context"
	"io"
	"io/fs"
	"math/rand/v2"
	"net/http"
	"sync"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/shallows/blockcache"
	"github.com/retrovibed/retrovibed/shallows/deeppool"
	"github.com/retrovibed/retrovibed/shallows/internal/cryptox"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/iox"
	"github.com/retrovibed/retrovibed/shallows/internal/uuidx"
)

type localstorage interface {
	io.ReaderAt
	io.WriterAt
}

type downloader interface {
	Download(ctx context.Context, id string, offset uint64, length uint64, dst io.Writer) error
}

func NewDeeppoolReaderAt(c *http.Client, md Metadata, l localstorage) *DeeppoolReaderAtCache {
	return &DeeppoolReaderAtCache{md: md, localstorage: l, m: &sync.Mutex{}, d: deeppool.NewRanger(c), blocklength: blockcache.DefaultBlockLength}
}

type DeeppoolReaderAtCache struct {
	localstorage
	d           downloader
	md          Metadata
	m           *sync.Mutex
	blocklength uint64
}

func (t *DeeppoolReaderAtCache) downloadChunk(prng *rand.ChaCha8, id string, offset uint64, length uint64) error {
	doffset, dlength := calculateBlockRange(t.blocklength, offset, length)
	if length <= doffset {
		dlength = t.blocklength
	}

	w, err := cryptox.NewOffsetWriterChaCha20(prng, io.NewOffsetWriter(t.localstorage, int64(doffset)), uint32(doffset))
	if err != nil {
		return err
	}

	ctx, done := context.WithCancel(context.Background())
	defer done()

	return t.d.Download(ctx, id, doffset, doffset+dlength, iox.NewTimeoutWriter(done, 3*time.Second, w))
}

func (t *DeeppoolReaderAtCache) ReadAt(p []byte, off int64) (n int, err error) {
	if n, err = t.localstorage.ReadAt(p, off); err == nil {
		return n, nil
	} else if errorsx.Ignore(err, fs.ErrNotExist) != nil {
		return n, err
	}

	t.m.Lock()
	defer t.m.Unlock()

	if n, err = t.localstorage.ReadAt(p, off); err == nil {
		return n, nil
	}

	if err = t.downloadChunk(MetadataChaCha8(t.md), t.md.ArchiveID, t.md.DiskOffset+uint64(off), t.md.DiskOffset+t.md.Bytes); err != nil {
		return n, errorsx.Wrap(err, "failed to donload from archive")
	}

	// re-read from disk. at this point we either succeeded or failed at resync from the archive.
	return t.localstorage.ReadAt(p, off)
}

func MetadataChaCha8(md Metadata) *rand.ChaCha8 {
	return cryptox.NewChaCha8(uuidx.FirstNonNil(uuid.FromStringOrNil(md.EncryptionSeed), uuid.FromStringOrNil(md.ID)).Bytes())
}
