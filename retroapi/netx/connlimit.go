package netx

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"sync"
	"sync/atomic"
)

// Dialer missing interface from the net package.
type Dialer interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

type ConnLimit struct {
	total    atomic.Uint64
	inbound  atomic.Uint64
	outbound atomic.Uint64
	rejected atomic.Uint64
	max      uint64
}

func NewConnUnlimited() *ConnLimit {
	return NewConnLimited(math.MaxUint64)
}

func NewConnLimited(n uint64) *ConnLimit {
	return &ConnLimit{max: n}
}

func (c *ConnLimit) Listener(l net.Listener) net.Listener {
	if pc, ok := l.(net.PacketConn); ok {
		return &limitedPacketConnListener{PacketConn: pc, listener: l, cl: c}
	}
	return &limitedListener{Listener: l, cl: c}
}

func (c *ConnLimit) Dialer(d Dialer) Dialer {
	return &limitedDialer{Dialer: d, cl: c}
}

func ConnLimitStatistics(dst io.Writer, c *ConnLimit) {
	fmt.Fprintf(dst, "connections total=%d inbound=%d outbound=%d max=%d rejected=%d\n",
		c.total.Load(), c.inbound.Load(), c.outbound.Load(), c.max, c.rejected.Load())
}

type limitedListener struct {
	net.Listener
	cl *ConnLimit
}

func (l *limitedListener) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}

	for {
		cur := l.cl.total.Load()
		if cur >= l.cl.max {
			conn.Close()
			l.cl.rejected.Add(1)
			return nil, errors.New("connection limit reached")
		}
		if l.cl.total.CompareAndSwap(cur, cur+1) {
			l.cl.inbound.Add(1)
			break
		}
	}

	return &limitedConn{Conn: conn, cl: l.cl, inbound: true}, nil
}

// limitedPacketConnListener wraps a net.Listener that also implements net.PacketConn
// (e.g. a UTP socket), preserving the PacketConn interface so callers like BindDHT
// can detect and use it for packet-based protocols.
type limitedPacketConnListener struct {
	net.PacketConn
	listener net.Listener
	cl       *ConnLimit
}

func (l *limitedPacketConnListener) Accept() (net.Conn, error) {
	conn, err := l.listener.Accept()
	if err != nil {
		return nil, err
	}

	for {
		cur := l.cl.total.Load()
		if cur >= l.cl.max {
			conn.Close()
			l.cl.rejected.Add(1)
			return nil, errors.New("connection limit reached")
		}
		if l.cl.total.CompareAndSwap(cur, cur+1) {
			l.cl.inbound.Add(1)
			break
		}
	}

	return &limitedConn{Conn: conn, cl: l.cl, inbound: true}, nil
}

func (l *limitedPacketConnListener) Addr() net.Addr  { return l.listener.Addr() }
func (l *limitedPacketConnListener) Close() error    { return l.listener.Close() }

type limitedDialer struct {
	Dialer
	cl *ConnLimit
}

func (d *limitedDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	for {
		cur := d.cl.total.Load()
		if cur >= d.cl.max {
			d.cl.rejected.Add(1)
			return nil, errors.New("connection limit reached")
		}
		if d.cl.total.CompareAndSwap(cur, cur+1) {
			d.cl.outbound.Add(1)
			break
		}
	}

	conn, err := d.Dialer.DialContext(ctx, network, address)
	if err != nil {
		d.cl.total.Add(^uint64(0))
		d.cl.outbound.Add(^uint64(0))
		return nil, err
	}

	return &limitedConn{Conn: conn, cl: d.cl, inbound: false}, nil
}

type limitedConn struct {
	net.Conn
	once    sync.Once
	cl      *ConnLimit
	inbound bool
}

func (c *limitedConn) Close() error {
	c.once.Do(func() {
		c.cl.total.Add(^uint64(0))
		if c.inbound {
			c.cl.inbound.Add(^uint64(0))
		} else {
			c.cl.outbound.Add(^uint64(0))
		}
	})
	return c.Conn.Close()
}
