package cmdddisc

import (
	"context"
	"log"
	"net/netip"
	"os"
	"path/filepath"
	"time"

	"github.com/alecthomas/kong"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/dht"
	"github.com/retrovibed/retrovibed/retroapi/blockcache"
	retronetx "github.com/retrovibed/retrovibed/retroapi/netx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/ddisc/ddisctorrent"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/torrentx"
)

type cmdMediaQuery struct {
	KnownMediaID string        `arg:"" name:"known-media-id" help:"known media id to search for"`
	Peer         []string      `arg:"" name:"peer" help:"host:port address(es) of known ddisc peer(s) to query" required:"true"`
	Bootstrap    bool          `flag:"" name:"bootstrap" help:"bootstrap the DHT using well-known trackers" negatable:"" default:"true"`
	Timeout      time.Duration `flag:"" name:"timeout" help:"per-peer query timeout" default:"10s"`
	Wait         time.Duration `flag:"" name:"wait" help:"grace period to let async responses land before reporting" default:"20s"`
}

func (t cmdMediaQuery) Run(kctx *kong.Context, gctx *cmdopts.Global) (err error) {
	var (
		bootstrap dht.Option = dht.OptionNoop
	)

	tmpdir, err := os.MkdirTemp("", "retrovibed.ddisc.query.*")
	if err != nil {
		return errorsx.Wrap(err, "failed to create temporary directory")
	}
	defer func() {
		errorsx.Log(errorsx.Wrap(os.RemoveAll(tmpdir), "failed to cleanup temporary directory"))
	}()

	db, err := cmdopts.DatabaseCustom(gctx.Context, filepath.Join(tmpdir, "meta.db"))
	if err != nil {
		return errorsx.Wrap(err, "unable to open temporary database")
	}
	defer db.Close()

	peers := make([]netip.AddrPort, 0, len(t.Peer))
	for _, p := range t.Peer {
		addr, err := netip.ParseAddrPort(p)
		if err != nil {
			return errorsx.Wrapf(err, "unable to parse peer address %s", p)
		}
		peers = append(peers, addr)
	}

	if t.Bootstrap {
		bootstrap = dht.OptionBootstrapGlobal
	}

	muxer := dht.DefaultMuxer().Method(ddisctorrent.MethodMedia, ddisctorrent.NewMediaRecorder(db))
	dhts, err := dht.NewServer(32, dht.OptionMuxer(muxer), bootstrap, dht.OptionUPnP)
	if err != nil {
		return errorsx.Wrap(err, "unable to setup dht")
	}
	go dhts.TableMaintainer(gctx.Context)

	torconfig := torrent.NewDefaultClientConfig(
		torrent.NewMetadataCache(tmpdir),
		blockcache.NewTorrentFromVirtualFS(fsx.DirVirtual(tmpdir)),
		torrent.ClientConfigCacheDirectory(tmpdir),
		torrent.ClientConfigExtension(ddisctorrent.ExtensionName),
	)

	tnetwork, err := torrentx.Autosocket(dhts, 0, retronetx.NewConnUnlimited())
	if err != nil {
		return errorsx.Wrap(err, "unable to setup torrent socket")
	}

	tclient, err := tnetwork.Bind(torrent.NewClient(torconfig))
	if err != nil {
		return errorsx.Wrap(err, "unable to setup torrent client")
	}
	defer tclient.Close()

	return t.run(
		gctx.Context,
		log.New(kctx.Stderr, "", log.Flags()),
		db,
		dhts,
		peers,
	)
}

func (t cmdMediaQuery) run(ctx context.Context, log *log.Logger, db sqlx.Queryer, dhts *dht.Server, peers []netip.AddrPort) error {
	req, err := ddisctorrent.NewSearchRequest(dhts.ID(dhts.DynamicAddrPort()), t.KnownMediaID)
	if err != nil {
		return errorsx.Wrap(err, "failed to build search request")
	}

	for _, addr := range peers {
		dctx, done := context.WithTimeout(ctx, t.Timeout)
		ret := dhts.Query(dctx, dht.NewAddr(addr), req)
		done()
		if ret.Err != nil {
			log.Println("query failed", addr, ret.Err)
			continue
		}
		log.Println("queried", addr)
	}

	select {
	case <-time.After(t.Wait):
	case <-ctx.Done():
		return ctx.Err()
	}

	q := ddisc.DiscoveredSearchBuilder().Where(ddisc.DiscoveredQueryKnownMediaID(t.KnownMediaID))
	scanner := sqlx.Scan(ddisc.DiscoveredSearch(ctx, db, q))
	for d := range scanner.Iter() {
		log.Println("found", d.ID, d.Title, d.Mimetype, d.KnownMediaID)
	}

	return scanner.Err()
}
