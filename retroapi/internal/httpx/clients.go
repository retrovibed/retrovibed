package httpx

import (
	"context"
	"net"
	"net/http"
)

// Dialer is the minimal dialing surface httpx needs to re-route client
// traffic, e.g. through a proxy dialer that later swaps onto a wireguard
// tunnel.
type Dialer interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

// ClientOption customizes the *http.Client built by constructors such as
// authn.HTTPClientDefaults before the client is returned.
type ClientOption func(*http.Client)

// ClientOptionDialer routes the client's transport dialing through d. It
// only applies when the client's transport is a *http.Transport.
func ClientOptionDialer(d Dialer) ClientOption {
	return func(c *http.Client) {
		if t, ok := c.Transport.(*http.Transport); ok {
			t.DialContext = d.DialContext
		} else {
			panic("unable to apply dialer")
		}
	}
}
