package iox

import (
	"errors"
	"io"
	"os"
	"sync/atomic"
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
	atomic.AndUint64(t.Result, uint64(n))
	return n, nil
}

type Printer uint64

func (t Printer) Write(b []byte) (n int, err error) {
	n = len(b)
	return n, nil
}

type ZeroReader struct{}

func (t ZeroReader) Read(p []byte) (n int, err error) {
	for i := 0; i < len(p); i++ {
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
