// Package dnscache caches DNS lookups
package dnscache

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/retrovibed/retrovibed/shallows/internal/debugx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/netx"
	"golang.org/x/time/rate"
)

type Resolver interface {
	LookupHost(ctx context.Context, host string) (addrs []string, err error)
}

func AutoProxyResolver() *ProxyPtr {
	return &ProxyPtr{_delegate: New(net.DefaultResolver)}
}

type ProxyPtr struct {
	mu        sync.RWMutex
	_delegate Resolver
}

func (t *ProxyPtr) Store(c Resolver) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t._delegate = c
}

func (t *ProxyPtr) LookupHost(ctx context.Context, host string) (addrs []string, err error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t._delegate.LookupHost(ctx, host)
}

type record struct {
	expired   time.Time
	addresses []string
}

type Cache struct {
	Resolver
	ttl     time.Duration
	limiter *rate.Limiter
	lock    sync.RWMutex
	cache   map[string]record
}

func New(r Resolver, options ...func(*Cache)) *Cache {
	return langx.Autoptr(langx.Clone(Cache{
		Resolver: r,
		ttl:      10 * time.Minute,
		limiter:  rate.NewLimiter(rate.Inf, 0),
		cache:    make(map[string]record, 128),
	}, options...))
}

func CacheOptionTTL(d time.Duration) func(*Cache) {
	return func(c *Cache) {
		c.ttl = d
	}
}

func CacheOptionLimiter(l *rate.Limiter) func(*Cache) {
	return func(c *Cache) {
		c.limiter = l
	}
}

// CacheOptionRateLimit sets the limiter from a uint32 events-per-second value.
// A value of 0 is treated as infinite (no limiting).
func CacheOptionRateLimit(n uint32) func(*Cache) {
	return func(c *Cache) {
		if n == 0 {
			c.limiter = rate.NewLimiter(rate.Inf, 0)
		} else {
			c.limiter = rate.NewLimiter(rate.Limit(n), 1)
		}
	}
}

func (t *Cache) LookupHost(ctx context.Context, host string) (addrs []string, err error) {
	t.lock.RLock()
	ts := time.Now()
	cached, ok := t.cache[host]
	flush := len(t.cache) > 256
	t.lock.RUnlock()

	if ok && cached.expired.After(ts) {
		return cached.addresses, nil
	}

	if err = t.limiter.Wait(ctx); err != nil {
		return nil, err
	}

	t.lock.Lock()
	defer t.lock.Unlock()

	cached, ok = t.cache[host]
	if ok && cached.expired.After(ts) {
		return cached.addresses, nil
	}

	defer debugx.Println("------ cache miss", host, ts, ">", cached.expired, t.limiter.Limit(), ok, flush)

	ips, err := t.Resolver.LookupHost(ctx, host)
	if err != nil {
		return nil, err
	}

	if flush {
		t.cache = make(map[string]record, 128)
	}
	t.cache[host] = record{expired: time.Now().Add(t.ttl), addresses: ips}

	return ips, nil
}

func NewDialer(r Resolver, d netx.Dialer) Dialer {
	return Dialer{Resolver: r, Dialer: d}
}

type Dialer struct {
	Resolver
	netx.Dialer
}

func (t Dialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}

	addrs, err := t.LookupHost(ctx, host)
	if err != nil {
		return nil, err
	}

	return t.Dialer.DialContext(ctx, network, net.JoinHostPort(langx.FirstNonZero(addrs...), port))
}
