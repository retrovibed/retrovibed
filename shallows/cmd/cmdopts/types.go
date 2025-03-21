package cmdopts

import (
	"errors"
	"net"
	"net/url"
	"reflect"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/retrovibed/retrovibed/internal/errorsx"
)

// ParseIP addresses
func ParseIP(ctx *kong.DecodeContext, target reflect.Value) (err error) {
	target.Set(reflect.ValueOf(net.ParseIP(ctx.Scan.Pop().String())))
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
