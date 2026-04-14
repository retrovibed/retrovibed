package library

import (
	"context"
	"io/fs"
	"net/http"
	"time"

	"github.com/retrovibed/retrovibed/shallows/blockcache"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
)

func calculateBlockRange(blocklength uint64, offset uint64, length uint64) (doffset, dlength uint64) {
	doffset = (offset / blocklength) * blocklength
	end := min(doffset+blocklength, length)
	dlength = end - doffset
	return doffset, dlength
}

func New(c *http.Client, v fsx.Virtual, lookup func(ctx context.Context, s string) (md *Metadata, err error)) fs.FS {
	return vstoragefs{Virtual: v, metadata: lookup, c: c}
}

type vstoragefs struct {
	c *http.Client
	fsx.Virtual
	metadata func(ctx context.Context, s string) (md *Metadata, err error)
}

func (t vstoragefs) Open(name string) (fs.File, error) {
	ctx, done := context.WithTimeout(context.Background(), time.Second)
	defer done()

	md, err := t.metadata(ctx, name)
	if err != nil {
		return nil, err
	}

	dcache, err := blockcache.NewDirectoryCache(t.Path(md.ID))
	if err != nil {
		return nil, err
	}

	return blockcache.NewFile(NewDeeppoolReaderAt(t.c, *md, dcache), md.CreatedAt, md.ID, md.Bytes, 0600), nil
}
