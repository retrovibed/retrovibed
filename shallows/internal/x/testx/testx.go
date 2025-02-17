package testx

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/james-lawrence/deeppool/internal/x/errorsx"
)

type deadline interface {
	Deadline() (deadline time.Time, ok bool)
}

func TempDir(t testing.TB) string {
	return t.TempDir()
}

func WithDeadline(t deadline) (context.Context, context.CancelFunc) {
	if d, ok := t.Deadline(); ok {
		return context.WithDeadline(context.Background(), d)
	}

	return context.WithCancel(context.Background())
}

func Fixture(path ...string) string {
	return filepath.Join(".fixtures", filepath.Join(path...))
}

// Read a file at the given path.
func Read(path ...string) io.Reader {
	return bytes.NewReader(errorsx.Must(os.ReadFile(filepath.Join(path...))))
}

// ReadString from the given file.
func ReadString(path ...string) string {
	return string(errorsx.Must(os.ReadFile(filepath.Join(path...))))
}

// Must is a small language extension for panicing on the common
// value, error return pattern. only used in tests.
func Must[T any](v T, err error) func(testing.TB) T {
	return func(t testing.TB) T {
		if err != nil {
			t.Fatal(err)
		}
		return v
	}
}

// Tempenvvar temporarily set the environment variable.
func Tempenvvar(k, v string, do func()) {
	o := os.Getenv(k)
	defer os.Setenv(k, o)
	if err := os.Setenv(k, v); err != nil {
		panic(err)
	}
	do()
}

func IOString(in io.Reader) string {
	return string(errorsx.Must(io.ReadAll(in)))
}

func IOBytes(in io.Reader) []byte {
	return errorsx.Must(io.ReadAll(in))
}
