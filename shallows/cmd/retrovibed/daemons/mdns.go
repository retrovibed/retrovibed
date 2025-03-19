package daemons

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/hashicorp/mdns"
	"github.com/retrovibed/retrovibed/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/internal/x/errorsx"
	"github.com/retrovibed/retrovibed/internal/x/netx"
)

func MulticastService(ctx context.Context, addr net.Listener) error {
	ipaddr := netx.AddrPort(addr.Addr())
	if ipaddr == nil {
		return fmt.Errorf("unable to get ip address from listener: %s", addr.Addr().String())
	}

	hostname := fmt.Sprintf("%s.", errorsx.Zero(os.Hostname()))
	info := []string{"shallows"}
	service, err := mdns.NewMDNSService(cmdopts.MachineID(), "_shallows._udp", "local.", hostname, int(ipaddr.Port()), nil, info)
	if err != nil {
		return err
	}

	server, err := mdns.NewServer(&mdns.Config{Zone: service})
	if err != nil {
		return err
	}
	go func() {
		defer server.Shutdown()
		<-ctx.Done()
	}()

	return nil
}
