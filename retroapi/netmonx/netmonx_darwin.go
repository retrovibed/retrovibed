//go:build darwin

package netmonx

import (
	"context"
	"net"
	"syscall"

	"golang.org/x/net/route"
	"golang.org/x/sys/unix"
)

func startWatcher(ctx context.Context, notify chan<- struct{}) error {
	fd, err := unix.Socket(unix.AF_ROUTE, unix.SOCK_RAW, 0)
	if err != nil {
		return err
	}
	unix.CloseOnExec(fd)

	// Close socket when context is cancelled to unblock Read.
	go func() {
		<-ctx.Done()
		unix.Close(fd)
	}()

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := unix.Read(fd, buf)
			if err != nil || n <= 0 {
				return
			}
			select {
			case notify <- struct{}{}:
			default:
			}
		}
	}()

	return nil
}

func platformMetered(_ string) bool { return false }

func getFallbackState() *State { return nil }

func defaultRouteInterface() string {
	rib, err := route.FetchRIB(syscall.AF_INET, route.RIBTypeRoute, 0)
	if err != nil {
		return ""
	}
	msgs, err := route.ParseRIB(route.RIBTypeRoute, rib)
	if err != nil {
		return ""
	}

	for _, msg := range msgs {
		rm, ok := msg.(*route.RouteMessage)
		if !ok || len(rm.Addrs) == 0 {
			continue
		}
		// Default route: destination address is 0.0.0.0
		dst, ok := rm.Addrs[0].(*route.Inet4Addr)
		if !ok {
			continue
		}
		if dst.IP != [4]byte{0, 0, 0, 0} {
			continue
		}
		// Interface index is stored in rm.Index.
		iface, err := net.InterfaceByIndex(rm.Index)
		if err != nil {
			continue
		}
		return iface.Name
	}
	return ""
}
