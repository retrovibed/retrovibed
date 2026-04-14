package cmdtorrent

import (
	"crypto/md5"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/storage"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/ddisc/ddisctorrent"
	"github.com/retrovibed/retrovibed/shallows/internal/dhtx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/torrentx"
)

type cmdDownload struct {
	Magnet    url.URL   `arg:"" name:"magnet" help:"magnet uri to download" required:"true"`
	Peers     []url.URL `flag:"" name:"peer" help:"uri of a peer to download from"`
	Bootstrap bool      `flag:"" name:"bootstrap" help:"bootstrap the dht using well known trackers"`
}

func (t cmdDownload) Run(gctx *cmdopts.Global) error {
	var (
		tnetwork  torrent.Binder
		digest               = md5.New()
		bootstrap dht.Option = dht.OptionNoop
	)

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
		dht.OptionLogger(log.Default()),
	)
	if err != nil {
		log.Fatalln(err)
	}

	go dhtx.Statistics(gctx.Context, time.Minute, _dht)
	go _dht.TableMaintainer()

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
		torrent.ClientConfigDebugLogger(log.Default()),
		torrent.ClientConfigInfoLogger(log.Default()),
		torrent.ClientConfigMaxOutstandingRequests(2048),
		torrent.ClientConfigExtension(ddisctorrent.ExtensionName),
	)

	if tnetwork, err = torrentx.Autosocket(_dht, 0); err != nil {
		return errorsx.Wrap(err, "unable to setup torrent socket")
	}

	tclient, err := tnetwork.Bind(torrent.NewClient(torconfig))
	if err != nil {
		return errorsx.Wrap(err, "unable to setup torrent client")
	}

	dl, _, err := tclient.Start(md, torrent.TunePeers(peers...))
	if err != nil {
		return errorsx.Wrap(err, "unable to start torrent")
	}
	go torrentx.DownloadProgress(gctx.Context, dl)

	n, err := torrent.DownloadInto(gctx.Context, digest, dl)
	if err != nil {
		return errorsx.Wrap(err, "torrent download failed")
	}

	log.Println("download succeed", n, "bytes", md5x.FormatUUID(digest))
	return nil
}
