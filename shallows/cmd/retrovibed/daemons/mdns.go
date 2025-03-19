package daemons

import (
	"context"
	"fmt"
	"log"
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
	info := []string{"retrovibed"}
	service, err := mdns.NewMDNSService(cmdopts.MachineID(), "_retrovibed._udp", "local.", hostname, int(ipaddr.Port()), nil, info)
	if err != nil {
		return err
	}

	log.Println("mdns", service.Instance, service.Service, service.Domain, service.HostName, service.Port, service.TXT)
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
