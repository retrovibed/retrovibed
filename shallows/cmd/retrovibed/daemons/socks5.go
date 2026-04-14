package daemons

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/retrovibed/retrovibed/shallows/dnscache"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/netx"
	"github.com/things-go/go-socks5"
)

type WireguardResolver struct {
	dnscache.Resolver
}

func (t WireguardResolver) Resolve(ctx context.Context, name string) (context.Context, net.IP, error) {
	addrs, err := t.LookupHost(ctx, name)
	return ctx, net.ParseIP(langx.FirstNonZero(addrs...)), err
}

func Socks5(ctx context.Context, d netx.Dialer, dns socks5.NameResolver, l net.Listener) error {
	server := socks5.NewServer(
		socks5.WithDial(d.DialContext),
		socks5.WithResolver(dns),
		socks5.WithLogger(socks5.NewLogger(log.New(os.Stderr, "socks5: ", log.Flags()))),
	)
	go func() {
		errorsx.Log(errorsx.Wrap(server.Serve(l), "socks5 failed"))
	}()

	log.Println("socks5 service enabled", l.Addr().String())
	return nil
}
