package daemons

import (
	"net"

	"github.com/egdaemon/wasinet/wasinet/wnetruntime"
	"github.com/retrovibed/retrovibed/shallows/internal/wireguardx"
	"golang.zx2c4.com/wireguard/tun/netstack"
)

// searchPluginSocket builds the wnetruntime.Socket search plugins get their
// network access through: routed through the wireguard tunnel when one is
// configured (wgnet != nil, and *netstack.Net already satisfies
// wnetruntime.Dialer directly), otherwise the host's own network stack.
// Either way only public addresses are reachable.
func searchPluginSocket(wgnet *netstack.Net) wnetruntime.Socket {
	if wgnet == nil {
		return wnetruntime.Virtual(&net.Dialer{}, &net.ListenConfig{}, net.DefaultResolver, wnetruntime.PublicFirewall())
	}
	return wnetruntime.Virtual(wgnet, wireguardx.PacketDialerAdapter(wgnet), wireguardx.HostLookupAdapter(wgnet), wnetruntime.PublicFirewall())
}
