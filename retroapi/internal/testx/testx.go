package testx

import (
	"context"
	"crypto/md5"
	"io"
	"iter"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/internal/errorsx"
	"github.com/stretchr/testify/require"
)

func Fixture(path ...string) string {
	return filepath.Join(".fixtures", filepath.Join(path...))
}

// Read a file at the given path.
func Read(path ...string) *os.File {
	return errorsx.Must(os.Open(filepath.Join(path...)))
}

// ReadString from the given file.
func ReadString(path ...string) string {
	return string(errorsx.Must(os.ReadFile(filepath.Join(path...))))
}

func ReadMD5(path ...string) string {
	d := md5.New()
	_ = errorsx.Must(d.Write(errorsx.Must(os.ReadFile(filepath.Join(path...)))))
	return uuid.FromBytesOrNil(d.Sum(nil)).String()
}

// Must is a small language extension for panicing on the common
// value, error return pattern. only used in tests.
func Must[T any](v T, err error) func(t testing.TB) T {
	return func(t testing.TB) T {
		require.NoError(t, err)
		return v
	}
}

// Tempenvvar temporarily set the environment variable.
func Tempenvvar(k, v string, do func()) {
	o := os.Getenv(k)
	defer os.Setenv(k, o) //nolint: errcheck
	if err := os.Setenv(k, v); err != nil {
		panic(err)
	}
	do()
}

func IOMD5(in io.Reader) string {
	digester := md5.New()
	errorsx.Must(io.Copy(digester, in))
	return uuid.FromBytesOrNil(digester.Sum(nil)).String()
}

func IOString(in io.Reader) string {
	return string(errorsx.Must(io.ReadAll(in)))
}

func IOBytes(in io.Reader) []byte {
	return errorsx.Must(io.ReadAll(in))
}

func Context(t testing.TB) (context.Context, context.CancelFunc) {
	return context.WithCancel(t.Context())
}

func SeqCount[T any](s iter.Seq[T]) (i uint64) {
	for range s {
		i++
	}
	return i
}

func Seq2Count[K any, T any](s iter.Seq2[K, T]) (i uint64) {
	for range s {
		i++
	}
	return i
}
