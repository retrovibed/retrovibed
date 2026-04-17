package iox

import (
	"context"
	"errors"
	"io"
	"os"
	"sync/atomic"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/internal/errorsx"
)

// IgnoreEOF returns nil if err is io.EOF
func IgnoreEOF(err error) error {
	if !errors.Is(err, io.EOF) {
		return err
	}

	return nil
}

// Error return just the error from an IO call ignoring the number of bytes.
func Error(_ int64, err error) error {
	return err
}

type errReader struct {
	error
}

func (t errReader) Read([]byte) (int, error) {
	return 0, t
}

// ErrReader returns an io.Reader that returns the provided error.
func ErrReader(err error) io.Reader {
	return errReader{err}
}

// Rewind an io.Seeker
func Rewind(o io.Seeker) error {
	_, err := o.Seek(0, io.SeekStart)
	return err
}

type writeNopCloser struct {
	io.Writer
}

func (writeNopCloser) Close() error { return nil }

// WriteNopCloser returns a WriteCloser with a no-op Close method wrapping
// the provided Writer w.
func WriteNopCloser(w io.Writer) io.WriteCloser {
	return writeNopCloser{w}
}

// Copy a file to another path
func Copy(from, to string) error {
	in, err := os.Open(from)
	if err != nil {
		return err
	}
	defer in.Close()

	i, err := in.Stat()
	if err != nil {
		return err
	}

	out, err := os.OpenFile(to, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, i.Mode())
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return out.Close()
}

// String reads the entire string from a reader.
// if the reader is also a seeker it'll rewind before reading.
func String(in io.Reader) (s string, err error) {
	var (
		raw []byte
	)

	if seeker, ok := in.(io.Seeker); ok {
		if err = Rewind(seeker); err != nil {
			return "", err
		}
	}

	if raw, err = io.ReadAll(in); err != nil {
		return "", err
	}

	return string(raw), nil
}

type Copied struct {
	Result *uint64
}

func (t Copied) Write(b []byte) (n int, err error) {
	n = len(b)
	_ = atomic.AddUint64(t.Result, uint64(n))
	return n, nil
}

func Zero() io.Reader {
	return zero{}
}

type zero struct{}

func (t zero) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = 0
	}

	return len(p), nil
}

type SectionWriter struct {
	w        io.WriterAt
	off, len int64
}

func NewSectionWriter(w io.WriterAt, off, len int64) *SectionWriter {
	return &SectionWriter{w, off, len}
}

func (me *SectionWriter) WriteAt(b []byte, off int64) (n int, err error) {
	if off >= me.len {
		err = io.EOF
		return
	}
	if off+int64(len(b)) > me.len {
		b = b[:me.len-off]
	}
	return me.w.WriteAt(b, me.off+off)
}

func NewTimeoutWriter(timedout context.CancelFunc, d time.Duration, dst io.Writer) *TimeoutWriter {
	m := time.NewTimer(d)
	go func() {
		<-m.C
		timedout()
	}()
	return &TimeoutWriter{
		d:     d,
		m:     m,
		inner: dst,
	}
}

type TimeoutWriter struct {
	d     time.Duration
	m     *time.Timer
	inner io.Writer
}

func (t TimeoutWriter) Write(b []byte) (n int, err error) {
	t.m.Reset(t.d)
	return t.inner.Write(b)
}

type readCompositeCloser struct {
	io.Reader
	closefn []func() error
}

func (t readCompositeCloser) Close() (err error) {
	for _, fn := range t.closefn {
		err = errorsx.Compact(err, fn())
	}
	return err
}

// WriteNopCloser returns a WriteCloser with a no-op Close method wrapping
// the provided Writer w.
func ReaderCompositeCloser(w io.Reader, closers ...func() error) io.ReadCloser {
	return readCompositeCloser{Reader: w, closefn: closers}
}
