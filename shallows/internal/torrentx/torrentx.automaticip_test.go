package torrentx

import (
	"context"
	"encoding/binary"
	"net"
	"net/netip"
	"testing"
	"time"

	"github.com/james-lawrence/torrent/dht"
	"github.com/retrovibed/retrovibed/shallows/internal/wireguardx"
	"github.com/stretchr/testify/require"
)

// natpmpdialer ignores the requested address and dials the fake gateway instead.
type natpmpdialer struct {
	gateway string
}

func (t natpmpdialer) DialContext(ctx context.Context, network, _ string) (net.Conn, error) {
	var d net.Dialer
	return d.DialContext(ctx, network, t.gateway)
}

func TestAutomaticIP(t *testing.T) {
	t.Run("server close returns while waiting for the port mapping to expire", func(t *testing.T) {
		gateway, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
		require.NoError(t, err)
		t.Cleanup(func() { _ = gateway.Close() })

		// minimal NAT-PMP responder: external address and port mapping requests, mapping lives for an hour.
		go func() {
			var req [16]byte
			for {
				n, from, err := gateway.ReadFromUDP(req[:])
				if err != nil {
					return
				}

				switch {
				case n == 2 && req[1] == 0:
					resp := make([]byte, 12)
					resp[1] = 0x80
					copy(resp[8:12], net.IPv4(203, 0, 113, 7).To4())
					_, _ = gateway.WriteToUDP(resp, from)
				case n == 12 && req[1] == 2:
					resp := make([]byte, 16)
					resp[1] = 0x82
					copy(resp[8:12], req[4:8]) // internal port, mapped external port.
					binary.BigEndian.PutUint32(resp[12:16], 3600)
					_, _ = gateway.WriteToUDP(resp, from)
				}
			}
		}()

		wcfg := &wireguardx.Config{}
		wcfg.Interface.DNS = []netip.Addr{netip.MustParseAddr("127.0.0.1")}

		s, err := dht.NewServer(20, AutomaticIP(wcfg, natpmpdialer{gateway: gateway.LocalAddr().String()}, 0))
		require.NoError(t, err)

		pc, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
		require.NoError(t, err)
		t.Cleanup(func() { _ = pc.Close() })

		_, err = s.ServeBinding(t.Context(), pc, pc.LocalAddr().(*net.UDPAddr).AddrPort())
		require.NoError(t, err)

		done := make(chan struct{})
		go func() {
			defer close(done)
			s.Close()
		}()

		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("Close blocked on the port mapping refresh goroutine")
		}
	})
}
