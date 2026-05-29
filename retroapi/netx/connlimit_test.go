package netx_test

import (
	"bytes"
	"context"
	"net"
	"strings"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/netx"
	"github.com/stretchr/testify/require"
)

type pipeListener struct {
	ch     chan net.Conn
	closed chan struct{}
	addr   net.Addr
}

func newPipeListener() *pipeListener {
	return &pipeListener{
		ch:     make(chan net.Conn, 16),
		closed: make(chan struct{}),
		addr:   &net.TCPAddr{},
	}
}

func (l *pipeListener) dial() net.Conn {
	server, client := net.Pipe()
	l.ch <- server
	return client
}

func (l *pipeListener) Accept() (net.Conn, error) {
	select {
	case conn := <-l.ch:
		return conn, nil
	case <-l.closed:
		return nil, net.ErrClosed
	}
}

func (l *pipeListener) Close() error {
	select {
	case <-l.closed:
	default:
		close(l.closed)
	}
	return nil
}

func (l *pipeListener) Addr() net.Addr { return l.addr }

type pipeDialer struct{}

func (d pipeDialer) DialContext(_ context.Context, _, _ string) (net.Conn, error) {
	_, client := net.Pipe()
	return client, nil
}

func TestConnLimitedListener(t *testing.T) {
	t.Run("allows connections up to max", func(t *testing.T) {
		cl := netx.NewConnLimited(2)
		raw := newPipeListener()
		l := cl.Listener(raw)

		_ = raw.dial()
		c1, err := l.Accept()
		require.NoError(t, err)
		require.NotNil(t, c1)

		_ = raw.dial()
		c2, err := l.Accept()
		require.NoError(t, err)
		require.NotNil(t, c2)

		c1.Close()
		c2.Close()
	})

	t.Run("rejects when at max", func(t *testing.T) {
		cl := netx.NewConnLimited(1)
		raw := newPipeListener()
		l := cl.Listener(raw)

		_ = raw.dial()
		c1, err := l.Accept()
		require.NoError(t, err)

		_ = raw.dial()
		c2, err := l.Accept()
		require.Error(t, err)
		require.Nil(t, c2)

		c1.Close()
	})

	t.Run("accepts after slot freed by close", func(t *testing.T) {
		cl := netx.NewConnLimited(1)
		raw := newPipeListener()
		l := cl.Listener(raw)

		_ = raw.dial()
		c1, err := l.Accept()
		require.NoError(t, err)
		c1.Close()

		_ = raw.dial()
		c2, err := l.Accept()
		require.NoError(t, err)
		require.NotNil(t, c2)
		c2.Close()
	})
}

func TestConnLimitedDialer(t *testing.T) {
	t.Run("allows dials up to max", func(t *testing.T) {
		cl := netx.NewConnLimited(2)
		d := cl.Dialer(pipeDialer{})

		c1, err := d.DialContext(context.Background(), "tcp", "127.0.0.1:0")
		require.NoError(t, err)

		c2, err := d.DialContext(context.Background(), "tcp", "127.0.0.1:0")
		require.NoError(t, err)

		c1.Close()
		c2.Close()
	})

	t.Run("rejects when at max", func(t *testing.T) {
		cl := netx.NewConnLimited(1)
		d := cl.Dialer(pipeDialer{})

		c1, err := d.DialContext(context.Background(), "tcp", "127.0.0.1:0")
		require.NoError(t, err)

		c2, err := d.DialContext(context.Background(), "tcp", "127.0.0.1:0")
		require.Error(t, err)
		require.Nil(t, c2)

		c1.Close()
	})

	t.Run("accepts after slot freed by close", func(t *testing.T) {
		cl := netx.NewConnLimited(1)
		d := cl.Dialer(pipeDialer{})

		c1, err := d.DialContext(context.Background(), "tcp", "127.0.0.1:0")
		require.NoError(t, err)
		c1.Close()

		c2, err := d.DialContext(context.Background(), "tcp", "127.0.0.1:0")
		require.NoError(t, err)
		require.NotNil(t, c2)
		c2.Close()
	})
}

func TestConnLimitMixed(t *testing.T) {
	t.Run("inbound and outbound share total limit", func(t *testing.T) {
		cl := netx.NewConnLimited(2)
		raw := newPipeListener()
		l := cl.Listener(raw)
		d := cl.Dialer(pipeDialer{})

		_ = raw.dial()
		c1, err := l.Accept()
		require.NoError(t, err)

		c2, err := d.DialContext(context.Background(), "tcp", "127.0.0.1:0")
		require.NoError(t, err)

		_ = raw.dial()
		c3, err := l.Accept()
		require.Error(t, err)
		require.Nil(t, c3)

		c4, err := d.DialContext(context.Background(), "tcp", "127.0.0.1:0")
		require.Error(t, err)
		require.Nil(t, c4)

		c1.Close()
		c2.Close()
	})
}

func TestConnLimitStatistics(t *testing.T) {
	cl := netx.NewConnLimited(10)
	d := cl.Dialer(pipeDialer{})

	c1, _ := d.DialContext(context.Background(), "tcp", "127.0.0.1:0")
	c2, _ := d.DialContext(context.Background(), "tcp", "127.0.0.1:0")
	c1.Close()

	var buf bytes.Buffer
	netx.ConnLimitStatistics(&buf, cl)

	out := buf.String()
	require.True(t, strings.Contains(out, "total=1"), "expected total=1 in: %s", out)
	require.True(t, strings.Contains(out, "outbound=1"), "expected outbound=1 in: %s", out)
	require.True(t, strings.Contains(out, "max=10"), "expected max=10 in: %s", out)

	_ = c2
}

func TestConnUnlimited(t *testing.T) {
	cl := netx.NewConnUnlimited()
	d := cl.Dialer(pipeDialer{})

	conns := make([]net.Conn, 100)
	for i := range conns {
		c, err := d.DialContext(context.Background(), "tcp", "127.0.0.1:0")
		require.NoError(t, err)
		conns[i] = c
	}
	for _, c := range conns {
		c.Close()
	}
}
