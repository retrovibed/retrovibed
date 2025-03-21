package torrentx

import (
	"fmt"
	"net"

	"github.com/james-lawrence/torrent"
	"github.com/retrovibed/retrovibed/internal/errorsx"

	"github.com/anacrolix/utp"
	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/dht/krpc"
	"github.com/james-lawrence/torrent/sockets"
)

func Autosocket(p int) (_ torrent.Binder, err error) {
	var (
		s1, s2  sockets.Socket
		tsocket *utp.Socket
	)

	tsocket, err = utp.NewSocket("udp", fmt.Sprintf(":%d", p))
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to open utp socket")
	}

	s1 = sockets.New(tsocket, tsocket)
	if addr, ok := tsocket.Addr().(*net.UDPAddr); ok {
		s, err := net.Listen("tcp", fmt.Sprintf(":%d", addr.Port))
		if err != nil {
			return nil, errorsx.Wrap(err, "unable to open tcp socket")
		}
		s2 = sockets.New(s, &net.Dialer{})
	}

	return torrent.NewSocketsBind(s1, s2), nil
}

func NodesFromReply(ret dht.QueryResult) (retni []krpc.NodeInfo) {
	if err := ret.ToError(); err != nil {
		return nil
	}

	ret.Reply.R.ForAllNodes(func(ni krpc.NodeInfo) {
		retni = append(retni, ni)
	})
	return retni
}
