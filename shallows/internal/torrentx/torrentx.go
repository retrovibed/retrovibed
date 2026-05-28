package torrentx

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"iter"
	"log"
	"net"
	"net/netip"
	"os"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/RoaringBitmap/roaring/v2"
	"github.com/james-lawrence/torrent"
	retronetx "github.com/retrovibed/retrovibed/retroapi/netx"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/backoffx"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/debugx"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/envx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/natpmp"
	"github.com/retrovibed/retrovibed/shallows/internal/netx"
	"github.com/retrovibed/retrovibed/shallows/internal/slicesx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/internal/wireguardx"
	"golang.org/x/time/rate"

	"github.com/anacrolix/missinggo/pubsub"
	"github.com/anacrolix/utp"
	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/dht/krpc"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/james-lawrence/torrent/sockets"
	"github.com/james-lawrence/torrent/tracker"

	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun/netstack"
)

func HashUID(md *int160.T) string {
	return md5x.FormatUUID(md5x.Digest(md.Bytes()))
}

func ClearIdleTorrents(ctx context.Context, d time.Duration, c *torrent.Client) {
	timex.Every(d, func() {
		errorsx.Log(c.Tune(torrent.ClientOperationClearIdleTorrents(func(stats torrent.Stats) bool {
			dead := stats.LastConnection.Add(d).Before(time.Now())
			if stats.Seeders == 0 && dead {
				return true
			}

			if max(stats.ActivePeers, stats.PendingPeers, stats.TotalPeers) > 0 {
				return false
			}

			return dead
		})))
	})
}

func AnnouncerFromClient(c *torrent.Client) tracker.Announce {
	return c.Config().AnnounceRequest(c)
}

func Peers(ctx context.Context, s *dht.Server, id int160.T) ([]dht.Peer, error) {
	seq, err := torrent.DHTAnnounce(ctx, s, id)
	if err != nil {
		return nil, errorsx.Wrapf(err, "failed to announce partition %s", id)
	}

	for pv := range seq.Each(ctx) {
		return pv.Peers, nil
	}

	return nil, errorsx.Wrapf(seq.Err(), "failed to announce partition %s", id)
}

func localsocket(network string, port uint16) (s0 *utp.Socket, s1 net.Listener, err error) {
	if s0, err = utp.NewSocket(network, fmt.Sprintf(":%d", port)); err != nil {
		return nil, nil, errorsx.Wrap(err, "unable to open utp socket")
	}

	if addr, ok := s0.Addr().(*net.UDPAddr); ok {
		if s1, err = net.Listen("tcp", fmt.Sprintf(":%d", addr.Port)); err != nil {
			return nil, nil, errorsx.Wrap(langx.FirstNonZero(err, s0.Close()), "unable to open tcp socket")
		}
	}

	return s0, s1, nil
}

func FileLargestRange(info *metainfo.Info) (m metainfo.FileInfo, offset, length int64) {
	m = FileLargest(info)
	return m, m.Offset(info), m.Length
}

func FileLargest(info *metainfo.Info) (m metainfo.FileInfo) {
	if len(info.Files) == 0 {
		return metainfo.FileInfo{Path: nil, Length: info.TotalLength()}
	}

	for _, c := range info.Files {
		if m.Length <= c.Length {
			m = c
		}
	}

	return m
}

func FileFirst(info *metainfo.Info, only func(metainfo.FileInfo) bool) (m metainfo.FileInfo) {
	if len(info.Files) == 0 {
		return metainfo.FileInfo{Path: nil, Length: info.TotalLength()}
	}

	for _, c := range info.Files {
		if only(c) {
			return c
		}
	}

	return m
}

func InfoFilePrint(info *metainfo.Info) {
	if len(info.Files) == 0 {
		log.Println(info.Name)
		return
	}

	for idx, c := range info.Files {
		log.Println(idx, c.Length, c.DisplayPath(info))
	}
}

func FilePrint(info *metainfo.Info, files ...metainfo.FileInfo) {
	for idx, c := range files {
		log.Println(idx, c.Length, c.DisplayPath(info))
	}
}

func FileBitmap(info *metainfo.Info, c metainfo.FileInfo) (m roaring.Bitmap) {
	return m
}

func Autosocket(_dht *dht.Server, p uint16, cl *retronetx.ConnLimit) (_ torrent.Binder, err error) {
	var (
		s1 *utp.Socket
		s2 net.Listener
	)

	if s1, s2, err = localsocket("udp", p); err != nil {
		return nil, err
	}

	return torrent.NewSocketsBind(
		sockets.New(cl.Listener(s1), cl.Dialer(s1)),
		sockets.New(cl.Listener(s2), cl.Dialer(&net.Dialer{})),
	).Options(torrent.BinderOptionDHT(_dht)), nil
}

func WireguardSocket(ctx context.Context, wcfg *wireguardx.Config) (_ *netstack.Net, err error) {
	var (
		// logger    = device.NewLogger(device.LogLevelError, "")
		logger = device.NewLogger(device.LogLevelVerbose, "")
	)

	tun, tnet, err := netstack.CreateNetTUN(
		slicesx.MapTransform(func(n netip.Prefix) netip.Addr { return n.Addr() }, wcfg.Interface.Addresses...),
		wcfg.Interface.DNS,
		langx.FirstNonZero(int(wcfg.Interface.MTU), wireguardx.DefaultMTU),
	)
	if err != nil {
		return nil, errorsx.Wrap(err, "failed to create network tun device")
	}

	dev := device.NewDevice(tun, conn.NewDefaultBind(), logger)

	w := asyncx.NewWakeup(ctx)

	diagnostic := func(ctx context.Context) error {
		return wireguardx.Diagnostic(os.Stderr, dev)
	}
	go debugx.OnSignal(ctx, diagnostic, syscall.SIGUSR1)
	go asyncx.Periodic(ctx, w, backoffx.Constant(5*time.Second), "wireguard statistics")
	asyncx.Background(ctx, w, diagnostic)

	for _, ipcset := range wireguardx.FormatIPCSet(wcfg) {
		if err = dev.IpcSet(ipcset); err != nil {
			return nil, errorsx.Wrap(err, "invalid ipcset for peer")
		}
	}

	if err = dev.Up(); err != nil {
		return nil, errorsx.Wrap(err, "network device failed to come up")
	}

	return tnet, nil
}

func SetupTorrentBinder(tnet *netstack.Net, port uint16, cl *retronetx.ConnLimit, opts ...torrent.BinderOption) (_ torrent.Binder, err error) {
	var (
		s0, s1    sockets.Socket
		utpsocket *utp.Socket
	)

	conn, err := tnet.ListenUDP(&net.UDPAddr{Port: int(port)})
	if err != nil {
		return nil, errorsx.Wrap(err, "failed to listen on port")
	}

	// log.Println("remote addr", errorsx.Zero(tun.Name()), conn.RemoteAddr())
	// log.Println("WAAAAT", errorsx.Zero(dev.IpcGet()))

	if utpsocket, err = utp.NewSocketFromPacketConn(conn); err != nil {
		return nil, errorsx.Wrap(err, "failed to create utp socket")
	}

	s0 = sockets.New(cl.Listener(utpsocket), cl.Dialer(utpsocket))
	if addr, ok := utpsocket.Addr().(*net.UDPAddr); ok {
		s, err := tnet.ListenTCP(&net.TCPAddr{Port: addr.Port})
		if err != nil {
			return nil, errorsx.Wrap(err, "unable to open tcp socket")
		}
		s1 = sockets.New(cl.Listener(s), cl.Dialer(tnet))
	}

	return torrent.NewSocketsBind(s0, s1).Options(opts...), nil
}

func ExternalPort(wcfg *wireguardx.Config, d netx.Dialer, port uint16) (_zero netip.AddrPort, _ time.Duration, err error) {
	dnsgateway := slicesx.FirstOrZero(wcfg.Interface.DNS...)

	client := natpmp.NewClient(dnsgateway, natpmp.OptionTimeout(15*time.Second), natpmp.OptionDialer(d))

	ex, err := client.GetExternalAddress()
	if err != nil {
		return _zero, 0, errorsx.Wrapf(err, "unable to determine external ip: %s", dnsgateway)
	}

	result, err := client.AddPortMapping("tcp", int(port), int(port), int(time.Hour/time.Second))
	if err != nil {
		return _zero, 0, errorsx.Wrapf(err, "unable to map port: %s", dnsgateway)
	}

	return netip.AddrPortFrom(netip.AddrFrom4(ex.ExternalIPAddress), result.MappedExternalPort), time.Duration(result.PortMappingLifetimeInSeconds) * time.Second, nil
}

func AutomaticIP(wcfg *wireguardx.Config, d netx.Dialer, port uint16) dht.Option {
	if port > 0 {
		return dht.OptionStaticPort(port)
	}

	return dht.OptionDynamicPort(func(ctx context.Context, sc *dht.Server, q dht.Binding, id int160.T, bestaddr netip.AddrPort, local net.PacketConn) (iter.Seq[netip.AddrPort], error) {
		return func(yield func(netip.AddrPort) bool) {
			for {
				addr, d, err := ExternalPort(wcfg, d, port)
				if err != nil {
					log.Println("failed to map ports", err)
					time.Sleep(time.Minute)
					continue
				}

				if !yield(addr) {
					return
				}

				time.Sleep(d)
			}
		}, nil
	})
}

// d.AnnounceTraversal(ctx, id, dht.AnnouncePeer(d, false))
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
		if minfo.ID().Cmp(int160.New([]byte(minfo.InfoBytes))) == 0 {
			return torrent.OptionInfo(minfo.InfoBytes)
		}

		panic("tisk tisk")
	} else if !errors.Is(err, os.ErrNotExist) {
		log.Println("unable to load torrent info, will attempt to locate it from peers", err)
	}

	return torrent.OptionNoop
}

func OptionTracker(tracker string) torrent.Option {
	if stringsx.Blank(tracker) {
		return torrent.OptionNoop
	}

	return torrent.OptionTrackers(tracker)
}

func Info(dl torrent.Torrent) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		stats := dl.Stats()
		info := dl.Info()
		md := dl.Metadata()

		var infoStatus string
		if info != nil {
			infoStatus = "downloaded"
		} else {
			infoStatus = "missing"
		}

		var peers []torrent.Peer
		if err := dl.Tune(torrent.TuneReadPeers(&peers)); err != nil {
			log.Println("unable to read peers", err)
			peers = nil
		}

		var b = &bytes.Buffer{}
		tw := tabwriter.NewWriter(b, 1, 0, 2, ' ', 0)

		fmt.Fprintf(tw, "%s - %s: info(%s) seeding(%t) (last connection: %s)\n",
			md.ID, md.DisplayName, infoStatus, stats.Seeding, timex.Human(stats.LastConnection))

		_, _ = fmt.Fprintf(tw, "  Maximum Allowed Peers %d\n", stats.MaximumAllowedPeers)
		_, _ = fmt.Fprintf(tw, "  Active Peers          %d\n", stats.ActivePeers)
		_, _ = fmt.Fprintf(tw, "  Half-Open Peers       %d\n", stats.HalfOpenPeers)
		_, _ = fmt.Fprintf(tw, "  Pending Peers         %d\n", stats.PendingPeers)
		_, _ = fmt.Fprintf(tw, "  Total Peers           %d\n", stats.TotalPeers)
		_, _ = fmt.Fprintf(tw, "  Seeders               %d\n", stats.Seeders)

		_, _ = fmt.Fprintf(tw, "  Pieces Missing        %d\n", stats.Missing)
		_, _ = fmt.Fprintf(tw, "  Pieces Outstanding    %d\n", stats.Outstanding)
		_, _ = fmt.Fprintf(tw, "  Pieces Unverified     %d\n", stats.Unverified)
		_, _ = fmt.Fprintf(tw, "  Pieces Completed      %d\n", stats.Completed)
		_, _ = fmt.Fprintf(tw, "  Pieces Failed         %d\n", stats.Failed)

		_, _ = fmt.Fprintf(tw, "  Bytes Written         %s\n", bytesx.Unit(stats.BytesWritten.Int64()))
		_, _ = fmt.Fprintf(tw, "  Bytes Read            %s\n", bytesx.Unit(stats.BytesRead.Int64()))
		_, _ = fmt.Fprintf(tw, "  Chunks Written        %d\n", stats.ChunksWritten.Int64())
		_, _ = fmt.Fprintf(tw, "  Chunks Read           %d\n", stats.ChunksRead.Int64())
		_, _ = fmt.Fprintf(tw, "  DHT Announce          %d\n", stats.DHTAnnounce.Int64())

		if len(peers) > 0 {
			active, pending := partitionPeers(peers)

			if len(active) > 0 {
				_, _ = fmt.Fprintln(tw, "  Active Connections")
				_, _ = fmt.Fprintf(tw, "    %-40s  %-28s    %s\n",
					"Peer ID", "Address", "Encrypted")
				for _, p := range active {
					_, _ = fmt.Fprintf(tw, "    %040s  %-28s  %t\n",
						p.ID.String(), p.String(), p.SupportsEncryption,
					)
				}
			}

			if len(pending) > 0 {
				_, _ = fmt.Fprintln(tw, "  Pending Peers")
				_, _ = fmt.Fprintf(tw, "    %-40s  %-28s  %-14s  %s\n",
					"Peer ID", "Address", "Last Attempt", "Attempts")
				for _, p := range pending {
					_, _ = fmt.Fprintf(tw, "    %040s  %-28s  %-14s  %d\n",
						p.ID.String(), p.String(), timex.Human(p.LastAttempt), p.Attempts,
					)
				}
			}
		}

		_ = tw.Flush()

		log.Println(b.String())

		if err := dl.Tune(torrent.TuneNewConns); err != nil {
			log.Println("unable to request new connections", err)
		}

		return nil
	}
}

func partitionPeers(peers []torrent.Peer) (active, pending []torrent.Peer) {
	for _, p := range peers {
		if p.LastAttempt.IsZero() || p.Attempts == 0 {
			active = append(active, p)
		} else {
			pending = append(pending, p)
		}
	}
	return active, pending
}

func DownloadProgress(ctx context.Context, dl torrent.Torrent) {
	var (
		statsfreq = envx.Duration(1*time.Minute, env.TorrentDownloadStats)
		sub       pubsub.Subscription
	)

	log.Println("monitoring download progress initiated", dl.Metadata().ID, statsfreq)
	defer log.Println("monitoring download progress completed", dl.Metadata().ID)

	// Revisit once resume is working.
	if err := dl.Tune(torrent.TuneSubscribe(&sub)); err != nil {
		log.Println("unable to subscribe", err)
		return
	}
	defer sub.Close()

	statst := time.NewTicker(statsfreq)
	l := rate.NewLimiter(rate.Every(time.Second), 1)

	for {
		select {
		case <-statst.C:
			stats := dl.Stats()
			info := dl.Info()
			md := dl.Metadata()

			log.Printf(
				"%s - %s: info(%t) %s\n", md.ID, md.DisplayName, info != nil, stats,
			)

			if err := dl.Tune(torrent.TuneNewConns); err != nil {
				log.Println("unable to request new connections", err)
				continue
			}
		case <-sub.Values:
			if !l.Allow() {
				continue
			}

			statst.Reset(statsfreq)
			stats := dl.Stats()
			md := dl.Metadata()

			log.Printf(
				"%s - %s: info(%t) %s\n", md.ID, md.DisplayName, true, stats,
			)
		case <-ctx.Done():
			return
		}
	}
}
