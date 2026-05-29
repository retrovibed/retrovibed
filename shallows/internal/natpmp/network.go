package natpmp

import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"time"

	"github.com/retrovibed/retrovibed/shallows/internal/netx"
)

const nAT_PMP_PORT = 5351
const nAT_TRIES = 9
const nAT_INITIAL_MS = 250

type dialer interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

// A caller that implements the NAT-PMP RPC protocol.
type network struct {
	dialer
	gateway netip.Addr
}

func (n *network) call(ctx context.Context, msg []byte, timeout time.Duration) (result []byte, err error) {
	type udpconn interface {
		ReadFrom(b []byte) (n int, addr net.Addr, err error)
		SetDeadline(t time.Time) error
		Write(b []byte) (int, error)
	}
	// var server net.UDPAddr
	// server.IP = n.gateway
	// server.Port = nAT_PMP_PORT
	// conn, err := net.DialUDP("udp", nil, &server)
	_conn, err := n.DialContext(ctx, "udp", fmt.Sprintf("%s:%d", n.gateway.String(), nAT_PMP_PORT))
	if err != nil {
		return nil, err
	}

	defer _conn.Close()
	// conn := _conn.(*net.UDPConn)
	conn := _conn.(udpconn)

	// 16 bytes is the maximum result size.
	result = make([]byte, 16)

	var finalTimeout time.Time
	if timeout != 0 {
		finalTimeout = time.Now().Add(timeout)
	}

	needNewDeadline := true

	var tries uint
	for tries = 0; (tries < nAT_TRIES && finalTimeout.IsZero()) || time.Now().Before(finalTimeout); {
		if needNewDeadline {
			nextDeadline := time.Now().Add((nAT_INITIAL_MS << tries) * time.Millisecond)
			err = conn.SetDeadline(minTime(nextDeadline, finalTimeout))
			if err != nil {
				return
			}
			needNewDeadline = false
		}
		_, err = conn.Write(msg)
		if err != nil {
			return
		}
		var bytesRead int
		var _remoteAddr net.Addr
		bytesRead, _remoteAddr, err = conn.ReadFrom(result)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				tries++
				needNewDeadline = true
				continue
			}
			return
		}

		remoteaddr := netx.AddrPort(_remoteAddr)
		if remoteaddr.Addr().Compare(n.gateway) != 0 {
			// Ignore this packet.
			// Continue without increasing retransmission timeout or deadline.
			continue
		}
		// Trim result to actual number of bytes received
		if bytesRead < len(result) {
			result = result[:bytesRead]
		}
		return
	}

	return nil, fmt.Errorf("timed out trying to contact gateway")
}

func minTime(a, b time.Time) time.Time {
	if a.IsZero() {
		return b
	}
	if b.IsZero() {
		return a
	}
	if a.Before(b) {
		return a
	}
	return b
}
