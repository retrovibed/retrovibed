package torrentx

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/james-lawrence/torrent"
	"github.com/retrovibed/retrovibed/internal/errorsx"

	"github.com/anacrolix/utp"
	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/dht/krpc"
	"github.com/james-lawrence/torrent/metainfo"
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

// read the info option from a on disk file
func OptionInfoFromFile(path string) torrent.Option {
	if minfo, err := metainfo.LoadFromFile(path); err == nil {
		return torrent.OptionInfo(minfo.InfoBytes)
		// if infob, err := os.ReadFile(path); err == nil {
		// return torrent.OptionInfo(infob)
	} else if !errors.Is(err, os.ErrNotExist) {
		log.Println("unable to load torrent info, will attempt to locate it from peers", err)
	}

	return torrent.OptionNoop
}

func RecordInfo(infopath string, dl torrent.Metadata) {
	if info := dl.InfoBytes; info != nil {
		errorsx.Log(errorsx.Wrap(os.WriteFile(infopath, info, 0600), "unable to record info file"))
	}
}
