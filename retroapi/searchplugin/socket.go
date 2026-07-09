package searchplugin

import (
	"context"
	"fmt"
	"net/netip"

	"github.com/egdaemon/wasinet/wasinet/wnetruntime"
	"golang.org/x/sys/unix"
)

// publicOnlySocket wraps a wnetruntime.Socket and rejects outbound
// connections to non-public address ranges (loopback/private/link-local/
// multicast/unspecified), delegating everything else to the wrapped socket.
// wnetruntime's own allow-list (wnetruntime.OptionAllow) isn't actually
// enforced by wnetruntime itself, so this is the sole gate.
type publicOnlySocket struct {
	wnetruntime.Socket
}

func newPublicOnlySocket() publicOnlySocket {
	return publicOnlySocket{Socket: wnetruntime.Unrestricted()}
}

func (t publicOnlySocket) Connect(ctx context.Context, fd int, sa unix.Sockaddr) error {
	if addr, ok := sockaddrAddr(sa); ok && !isPublic(addr) {
		return fmt.Errorf("searchplugin: connection to non-public address %s blocked", addr)
	}
	return t.Socket.Connect(ctx, fd, sa)
}

func (t publicOnlySocket) SendTo(ctx context.Context, fd int, sa unix.Sockaddr, vecs [][]byte, oob []byte, flags int) (int, error) {
	if addr, ok := sockaddrAddr(sa); ok && !isPublic(addr) {
		return 0, fmt.Errorf("searchplugin: send to non-public address %s blocked", addr)
	}
	return t.Socket.SendTo(ctx, fd, sa, vecs, oob, flags)
}

func sockaddrAddr(sa unix.Sockaddr) (netip.Addr, bool) {
	switch a := sa.(type) {
	case *unix.SockaddrInet4:
		return netip.AddrFrom4(a.Addr), true
	case *unix.SockaddrInet6:
		return netip.AddrFrom16(a.Addr), true
	default:
		return netip.Addr{}, false
	}
}

func isPublic(addr netip.Addr) bool {
	return !(addr.IsPrivate() || addr.IsLoopback() || addr.IsLinkLocalUnicast() || addr.IsLinkLocalMulticast() || addr.IsMulticast() || addr.IsUnspecified())
}
