package cmdtorrent

import (
	"crypto/md5"
	"log"
	"net/url"
	"os"
	"syscall"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/storage"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/ddisc/ddisctorrent"
	"github.com/retrovibed/retrovibed/shallows/internal/debugx"
	"github.com/retrovibed/retrovibed/shallows/internal/dhtx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/torrentx"
)

type cmdDownload struct {
	Magnet    url.URL   `arg:"" name:"magnet" help:"magnet uri to download (e.g. magnet:?xt=urn:btih:...)" required:"true"`
	Peers     []url.URL `flag:"" name:"peer" help:"peer uri to connect to (e.g. p0000000000000000000000000000000000000000://[::1]:3196) — may be repeated"`
	Bootstrap bool      `flag:"" name:"bootstrap" help:"bootstrap the DHT using well-known trackers (default: true; use --no-bootstrap to skip)" negatable:"" default:"true"`
}

func (t cmdDownload) Run(gctx *cmdopts.Global) error {
	var (
		tnetwork     torrent.Binder
		digest                                  = md5.New()
		bootstrap    dht.Option                 = dht.OptionNoop
		dhtdebug     dht.Option                 = dht.OptionNoop
		torrentinfo                             = torrent.ClientConfigInfoLogger(log.Default())
		torrentdebug torrent.ClientConfigOption = torrent.ClientConfigNoop
	)

	if gctx.Verbosity >= 2 {
		log.Println("-------------------------- DHT DEBUG LOGGING ENABLED --------------------------")
		dhtdebug = dht.OptionLogger(log.Default())
	}

	if gctx.Verbosity >= 1 {
		log.Println("-------------------------- TORRENT DEBUG LOGGING ENABLED --------------------------")
		torrentdebug = torrent.ClientConfigDebugLogger(log.Default())
	}

	log.Println("-------------------------- ", gctx.Verbosity, " --------------------------")

	if t.Bootstrap {
		bootstrap = dht.OptionBootstrapGlobal
	}

	md, err := torrent.NewFromMagnet(t.Magnet.String())
	if err != nil {
		return errorsx.Wrap(err, "unable to parse magnet uri")
	}

	peers := torrent.NewPeersFromURI(t.Peers...)

	tmpdir, err := os.MkdirTemp("", "retrovibed.torrent.test.*")
	if err != nil {
		return errorsx.Wrap(err, "failed to create temporary directory for torrent testing")
	}
	defer func() {
		errorsx.Log(errorsx.Wrap(os.RemoveAll(tmpdir), "failed to cleanup temp directory"))
	}()

	rootvfs := fsx.DirVirtual(tmpdir)
	tvfs := fsx.DirVirtual(rootvfs.Path("torrents"))
	tstore := storage.NewFile(tvfs.Path())
	_dht, err := dht.NewServer(
		32,
		bootstrap,
		dht.OptionUPnP,
		dhtdebug,
	)
	if err != nil {
		log.Fatalln(err)
	}

	go _dht.TableMaintainer(gctx.Context)

	torconfig := torrent.NewDefaultClientConfig(
		torrent.NewMetadataCache(tvfs.Path()),
		tstore,
		torrent.ClientConfigCacheDirectory(tvfs.Path()),
		torrent.ClientConfigPEX(true),
		torrent.ClientConfigSeed(true),
		torrent.ClientConfigDialTimeouts(4*time.Second, 3*time.Minute),
		torrent.ClientConfigHTTPUserAgent("retrovibed/0.0"),
		torrent.ClientConfigConnectionClosed(func(id int160.T, stats torrent.ConnStats, remaining int) {
			log.Println("connection closed", id, remaining, spew.Sdump(stats))
		}),
		torrent.ClientConfigMaxOutstandingRequests(2048),
		torrent.ClientConfigExtension(ddisctorrent.ExtensionName),
		torrentinfo,
		torrentdebug,
	)

	if tnetwork, err = torrentx.Autosocket(_dht, 0); err != nil {
		return errorsx.Wrap(err, "unable to setup torrent socket")
	}

	tclient, err := tnetwork.Bind(torrent.NewClient(torconfig))
	if err != nil {
		return errorsx.Wrap(err, "unable to setup torrent client")
	}

	dl, _, err := tclient.Start(md, torrent.TunePeers(peers...), torrent.TuneAnnounceUntilComplete)
	if err != nil {
		return errorsx.Wrap(err, "unable to start torrent")
	}

	go debugx.OnSignal(gctx.Context, dhtx.Statistics(_dht), syscall.SIGUSR1)
	go debugx.OnSignal(gctx.Context, torrentx.Info(dl), syscall.SIGUSR1)
	go torrentx.DownloadProgress(gctx.Context, dl)

	n, err := torrent.DownloadInto(gctx.Context, digest, dl)
	if err != nil {
		return errorsx.Wrap(err, "torrent download failed")
	}

	log.Println("download succeed", n, "bytes", md5x.FormatUUID(digest))
	return nil
}
