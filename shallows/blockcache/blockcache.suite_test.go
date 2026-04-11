package blockcache_test

import (
	"fmt"
	"io"
	"os"
	"sync"
	"testing"

	"github.com/retrovibed/retrovibed/internal/testx"
)

func TestMain(m *testing.M) {
	testx.Logging()
	os.Exit(m.Run())
}

type ByteArrayCache struct {
	data []byte
	mu   sync.RWMutex
}

func NewByteArrayCache(capacity uint16) *ByteArrayCache {
	return &ByteArrayCache{
		data: make([]byte, capacity),
	}
}

func (fbc *ByteArrayCache) ReadAt(p []byte, off int64) (n int, err error) {
	fbc.mu.RLock()
	defer fbc.mu.RUnlock()

	if off < 0 {
		return 0, fmt.Errorf("ReadAt: negative offset: %d", off)
	}

	start := int(off)
	end := start + len(p)

	if start >= len(fbc.data) {
		return 0, io.EOF
	}
	if end > len(fbc.data) {
		end = len(fbc.data)
		err = io.EOF
	}

	n = copy(p, fbc.data[start:end])
	return n, err
}

func (fbc *ByteArrayCache) WriteAt(p []byte, off int64) (n int, err error) {
	fbc.mu.Lock()
	defer fbc.mu.Unlock()

	if off < 0 {
		return 0, fmt.Errorf("WriteAt: negative offset: %d", off)
	}

	start := int(off)
	end := start + len(p)

	if end > len(fbc.data) {
		return 0, fmt.Errorf("WriteAt: write beyond cache capacity (offset %d, data length %d, capacity %d)", off, len(p), len(fbc.data))
	}

	n = copy(fbc.data[start:end], p)
	return n, nil
}
