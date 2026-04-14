package blockcache

import (
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
)

type cache interface {
	io.ReaderAt
	io.WriterAt
}

const DefaultBlockLength = 32 * bytesx.MiB

type OptionDirCache func(*DirCache)

func OptionDirCacheBlockLength(n int64) OptionDirCache {
	return func(dc *DirCache) {
		dc.BlockLength = n
	}
}

func NewDirectoryCache(dir string, options ...OptionDirCache) (*DirCache, error) {
	if !fsx.Exists(dir) {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return nil, errorsx.Wrapf(err, "unable to ensure directory exists: %s", dir)
		}
	}

	c := langx.Clone(DirCache{
		root:        dir,
		BlockLength: DefaultBlockLength,
	}, options...)
	return &c, nil
}

type DirCache struct {
	root        string
	BlockLength int64
}

func (t DirCache) path(off int64) string {
	block := off / int64(t.BlockLength)
	path := filepath.Join(t.root, strconv.FormatInt(block, 10))

	return path
}

func (t DirCache) offset(off int64) (int64, int64) {
	pos := off % t.BlockLength
	max := t.BlockLength - pos

	return pos, max
}

func (t DirCache) ReadAt(p []byte, off int64) (n int, err error) {
	readchunk := func(p []byte, off int64) (n int, err error) {
		dst, err := os.Open(t.path(off))
		if err != nil {
			return 0, errorsx.Wrapf(err, "unable to open for reading: offset(%d) n(%d)", off, len(p))
		}
		defer dst.Close() //nolint:errcheck

		start, max := t.offset(off)
		if int64(len(p)) > max {
			p = p[:max]
		}

		return dst.ReadAt(p, start)
	}

	for n < len(p) {
		_n, _err := readchunk(p[n:], off+int64(n))
		if _err != nil {
			return n + _n, _err
		}

		n += _n
	}

	return n, nil
}

func (t *DirCache) WriteAt(p []byte, off int64) (n int, err error) {
	writechunk := func(p []byte, ioff int64) (n int, err error) {
		path := t.path(ioff)
		dst, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return 0, errorsx.Wrapf(err, "unable to open for writing offset(%d) n(%d)", ioff, len(p))
		}
		defer dst.Close()

		start, max := t.offset(ioff)
		if int64(len(p)) > max {
			p = p[:max]
		}

		return dst.WriteAt(p, start)
	}

	for n < len(p) {
		_n, _err := writechunk(p[n:], off+int64(n))
		if _err != nil {
			return n + _n, _err
		}
		n += _n
	}

	return n, nil
}
