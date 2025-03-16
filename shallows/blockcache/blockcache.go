package blockcache

import (
	"os"

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

func (t DirCache) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return nil, nil
}
