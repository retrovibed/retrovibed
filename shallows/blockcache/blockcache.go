package blockcache

import (
	"io/fs"

	"github.com/james-lawrence/deeppool/internal/x/bytesx"
)

// Open(name string) (File, error)
func NewDirectoryCache(dir string) DirCache {
	return DirCache{
		BlockLength: 32 * bytesx.MiB,
	}
}

type DirCache struct {
	BlockLength int64
}

func (t DirCache) Open(name string) (_ fs.File, err error) {
	return nil, nil
}
