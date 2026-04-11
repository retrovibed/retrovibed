package dnscache

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"
)

// stubResolver is a Resolver whose behaviour is controlled per call.
type stubResolver struct {
	calls int
	addrs []string
	err   error
}

func (s *stubResolver) LookupHost(_ context.Context, _ string) ([]string, error) {
	s.calls++
	return s.addrs, s.err
}

// stubDialer records the address it was asked to dial.
type stubDialer struct {
	dialedAddr string
	conn       net.Conn
	err        error
}

func (s *stubDialer) DialContext(_ context.Context, _, address string) (net.Conn, error) {
	s.dialedAddr = address
	return s.conn, s.err
}

func TestCache(t *testing.T) {
	t.Run("cache miss calls through to resolver", func(t *testing.T) {
		r := &stubResolver{addrs: []string{"1.2.3.4"}}
		c := New(r)

		addrs, err := c.LookupHost(t.Context(), "example.com")
		require.NoError(t, err)
		require.Equal(t, []string{"1.2.3.4"}, addrs)
		require.Equal(t, 1, r.calls)
	})

	t.Run("cache hit does not call resolver again", func(t *testing.T) {
		r := &stubResolver{addrs: []string{"1.2.3.4"}}
		c := New(r)

		_, err := c.LookupHost(t.Context(), "example.com")
		require.NoError(t, err)

		addrs, err := c.LookupHost(t.Context(), "example.com")
		require.NoError(t, err)
		require.Equal(t, []string{"1.2.3.4"}, addrs)
		require.Equal(t, 1, r.calls, "resolver should only be called once")
	})

	t.Run("different hosts are cached independently", func(t *testing.T) {
		r := &stubResolver{addrs: []string{"1.2.3.4"}}
		c := New(r)

		_, err := c.LookupHost(t.Context(), "a.example.com")
		require.NoError(t, err)

		_, err = c.LookupHost(t.Context(), "b.example.com")
		require.NoError(t, err)

		require.Equal(t, 2, r.calls)
	})

	t.Run("expired entry calls resolver again", func(t *testing.T) {
		r := &stubResolver{addrs: []string{"1.2.3.4"}}
		c := New(r, CacheOptionTTL(-time.Second))

		_, err := c.LookupHost(t.Context(), "example.com")
		require.NoError(t, err)

		_, err = c.LookupHost(t.Context(), "example.com")
		require.NoError(t, err)

		require.Equal(t, 2, r.calls, "expired entry should trigger a new lookup")
	})

	t.Run("rate limiter blocks when context is cancelled", func(t *testing.T) {
		r := &stubResolver{addrs: []string{"1.2.3.4"}}
		c := New(r,
			CacheOptionTTL(-time.Second),
			CacheOptionLimiter(rate.NewLimiter(rate.Every(time.Hour), 1)),
		)

		// first lookup consumes the only token
		_, err := c.LookupHost(t.Context(), "example.com")
		require.NoError(t, err)

		// second lookup must wait for the refill; cancel ctx to unblock
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		_, err = c.LookupHost(ctx, "example.com")
		require.ErrorIs(t, err, context.Canceled)
		require.Equal(t, 1, r.calls, "resolver should not be reached when limiter blocks")
	})

	t.Run("resolver error is propagated and not cached", func(t *testing.T) {
		lookupErr := errors.New("dns failure")
		r := &stubResolver{err: lookupErr}
		c := New(r)

		_, err := c.LookupHost(t.Context(), "example.com")
		require.ErrorIs(t, err, lookupErr)

		// second call should also reach the resolver since nothing was cached
		r.err = nil
		r.addrs = []string{"1.2.3.4"}
		addrs, err := c.LookupHost(t.Context(), "example.com")
		require.NoError(t, err)
		require.Equal(t, []string{"1.2.3.4"}, addrs)
		require.Equal(t, 2, r.calls)
	})
}

func TestProxyPtr(t *testing.T) {
	t.Run("delegates LookupHost to stored resolver", func(t *testing.T) {
		r := &stubResolver{addrs: []string{"1.2.3.4"}}
		p := &ProxyPtr{_delegate: New(r)}

		addrs, err := p.LookupHost(t.Context(), "example.com")
		require.NoError(t, err)
		require.Equal(t, []string{"1.2.3.4"}, addrs)
	})

	t.Run("Store replaces the delegate", func(t *testing.T) {
		first := &stubResolver{addrs: []string{"1.1.1.1"}}
		second := &stubResolver{addrs: []string{"2.2.2.2"}}

		p := &ProxyPtr{_delegate: New(first)}

		addrs, err := p.LookupHost(t.Context(), "example.com")
		require.NoError(t, err)
		require.Equal(t, []string{"1.1.1.1"}, addrs)

		p.Store(New(second))

		addrs, err = p.LookupHost(t.Context(), "example.com")
		require.NoError(t, err)
		require.Equal(t, []string{"2.2.2.2"}, addrs)
	})

	t.Run("AutoProxyResolver initializes with default resolver", func(t *testing.T) {
		p := AutoProxyResolver()
		require.NotNil(t, p)
		require.NotNil(t, p._delegate)
	})
}

func TestDialer(t *testing.T) {
	t.Run("resolves host and dials resolved address", func(t *testing.T) {
		r := &stubResolver{addrs: []string{"1.2.3.4"}}
		d := &stubDialer{}

		dialer := NewDialer(New(r), d)
		_, _ = dialer.DialContext(t.Context(), "tcp", "example.com:8080")

		require.Equal(t, "1.2.3.4:8080", d.dialedAddr)
	})

	t.Run("propagates resolver error", func(t *testing.T) {
		lookupErr := errors.New("dns failure")
		r := &stubResolver{err: lookupErr}
		d := &stubDialer{}

		dialer := NewDialer(New(r), d)
		_, err := dialer.DialContext(t.Context(), "tcp", "example.com:8080")

		require.ErrorIs(t, err, lookupErr)
		require.Empty(t, d.dialedAddr)
	})

	t.Run("propagates dial error", func(t *testing.T) {
		dialErr := errors.New("connection refused")
		r := &stubResolver{addrs: []string{"1.2.3.4"}}
		d := &stubDialer{err: dialErr}

		dialer := NewDialer(New(r), d)
		_, err := dialer.DialContext(t.Context(), "tcp", "example.com:8080")

		require.ErrorIs(t, err, dialErr)
	})

	t.Run("returns error on malformed address", func(t *testing.T) {
		r := &stubResolver{addrs: []string{"1.2.3.4"}}
		d := &stubDialer{}

		dialer := NewDialer(New(r), d)
		_, err := dialer.DialContext(t.Context(), "tcp", "not-a-valid-address")

		require.Error(t, err)
	})
}
