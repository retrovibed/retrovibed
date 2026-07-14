package cmdopts

import (
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

// ParseIP addresses
func ParseIP(ctx *kong.DecodeContext, target reflect.Value) (err error) {
	target.Set(reflect.ValueOf(net.ParseIP(ctx.Scan.Pop().String())))
	return nil
}

// ParseDurationInf parses a time.Duration flag value. In addition to normal
// time.ParseDuration syntax (e.g. "1h30m"), it accepts "infinity"
// (case-insensitive) as an alias for the maximum representable duration,
// used to signal "no timeout".
func ParseDurationInf(ctx *kong.DecodeContext, target reflect.Value) (err error) {
	t, err := ctx.Scan.PopValue("duration")
	if err != nil {
		return err
	}

	var d time.Duration
	switch v := t.Value.(type) {
	case string:
		switch {
		case strings.EqualFold(v, "infinity"):
			d = time.Duration(math.MaxInt)
		default:
			if d, err = time.ParseDuration(v); err != nil {
				return errorsx.Wrapf(err, "expected duration but got %q", v)
			}
		}
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		d = reflect.ValueOf(v).Convert(reflect.TypeOf(time.Duration(0))).Interface().(time.Duration)
	default:
		return fmt.Errorf("expected duration but got %q", v)
	}

	target.Set(reflect.ValueOf(d))
	return nil
}

func ParseTCPAddr(ctx *kong.DecodeContext, target reflect.Value) (err error) {
	if ctx.Scan.Len() == 0 {
		return nil
	}

	var (
		saddr = ctx.Scan.Pop().String()
	)

	var (
		addr *net.TCPAddr
	)

	if addr, err = net.ResolveTCPAddr("tcp", saddr); err != nil {
		return errorsx.Wrapf(err, "unable to resolve tcp address %s - %+v", saddr, ctx)
	}

	target.Set(reflect.ValueOf(addr))

	return nil
}

func ParseTCPAddrArray(ctx *kong.DecodeContext, target reflect.Value) (err error) {
	if ctx.Scan.Len() == 0 {
		return nil
	}

	var (
		results []*net.TCPAddr
		token   = ctx.Scan.Pop().String()
	)

	token = strings.ReplaceAll(token, "\n", " ")
	token = strings.ReplaceAll(token, ",", " ")
	for _, saddr := range strings.Split(token, " ") {
		var (
			addr *net.TCPAddr
		)

		if addr, err = net.ResolveTCPAddr("tcp", saddr); err != nil {
			return errorsx.Wrapf(err, "unable to resolve tcp address %s : %s", saddr, token)
		}

		results = append(results, addr)
	}

	target.Set(reflect.ValueOf(results))
	return nil
}

type Listener struct {
	uri *url.URL
	s   net.Listener
}

func (t Listener) MarshalText() (text []byte, err error) {
	return []byte(t.s.Addr().String()), nil
}

func (t *Listener) UnmarshalText(text []byte) (err error) {
	uri, err := url.Parse(string(text))
	if err != nil {
		return err
	}

	switch uri.Scheme {
	case "unix", "tcp", "tcp4", "tcp6", "udp":
		t.uri = uri
		return nil
	default:
		return errorsx.Wrapf(errors.ErrUnsupported, "network: %s", uri.String())
	}
}

func (t Listener) Socket() (net.Listener, error) {
	return net.Listen(t.uri.Scheme, t.uri.Host)
}

// IOOut is a flag type that opens a file for writing, or uses stdout when the path is "-".
type IOOut struct {
	path string
}

func (t *IOOut) UnmarshalText(text []byte) error {
	t.path = string(text)
	return nil
}

func (t IOOut) MarshalText() ([]byte, error) {
	return []byte(t.path), nil
}

// Open returns a WriteCloser for the output. When path is "-", writes go to
// fallback and Close is a no-op. Otherwise a new file is created (or truncated).
func (t IOOut) Open(fallback io.Writer) (io.WriteCloser, error) {
	if t.path == "-" {
		return nopWriteCloser{fallback}, nil
	}
	return os.OpenFile(t.path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
}

type nopWriteCloser struct{ io.Writer }

func (nopWriteCloser) Close() error { return nil }

type FileContents string

func (t FileContents) MarshalText() (text []byte, err error) {
	return []byte(t), nil
}

func (s *FileContents) UnmarshalText(text []byte) error {
	// Check for the @ prefix
	if path, ok := strings.CutPrefix(string(text), "@"); ok {
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %q: %w", path, err)
		}
		*s = FileContents(content)
	} else {
		*s = FileContents(text)
	}

	return nil
}
