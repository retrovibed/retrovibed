package cmdddisc

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/dht"
	"github.com/retrovibed/retrovibed/retroapi/blockcache"
	retronetx "github.com/retrovibed/retrovibed/retroapi/netx"
	"github.com/retrovibed/retrovibed/retroapi/searchplugin"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/ddisc/ddisctorrent"
	"github.com/retrovibed/retrovibed/shallows/internal/cryptox"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/torrentx"
)

// cmdMediaDiscover is a manual/diagnostic entry point into the same
// local/partition/plugin discovery chain LocateMedia runs automatically.
type cmdMediaDiscover struct {
	KnownMediaID string `arg:"" name:"known-media-id" help:"known media id to search for"`
	Title        string `flag:"" name:"title" help:"title to search external search plugins for; the plugin strategy no-ops without one"`
	Category     string `flag:"" name:"category" help:"mimex category (video, audio, image, text, application) to search external search plugins for"`
	Bootstrap    bool   `flag:"" name:"bootstrap" help:"bootstrap the DHT using well-known trackers" negatable:"" default:"true"`
	Partitions   uint8  `flag:"" name:"partitions" help:"number of partitions to split the infohash space into, must match the daemon(s) you're querying" default:"128"`
	Seed         string `flag:"" name:"seed" help:"seed to generate partition spaces, must match the daemon(s) you're querying" default:"retrovibed-ddisc"`
}

func (t cmdMediaDiscover) Run(kctx *kong.Context, gctx *cmdopts.Global) (err error) {
	var (
		bootstrap dht.Option = dht.OptionNoop
	)

	tmpdir, err := os.MkdirTemp("", "retrovibed.ddisc.discover.*")
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

	plugins, err := searchplugin.NewRegistry(gctx.Context)
	if err != nil {
		return errorsx.Wrap(err, "unable to start search plugin registry")
	}

	partitions := ddisc.Partitions(uint16(t.Partitions), cryptox.NewChaCha8(t.Seed))

	return t.run(
		gctx.Context,
		log.New(kctx.Stderr, "", log.Flags()),
		db,
		dhts,
		partitions,
		plugins,
	)
}

func (t cmdMediaDiscover) run(ctx context.Context, log *log.Logger, db sqlx.Queryer, dhts *dht.Server, partitions *ddisc.Partition, plugins *searchplugin.Registry) error {
	req := ddisc.DiscoverRequest{
		KnownMediaID: t.KnownMediaID,
		Title:        t.Title,
		Category:     t.Category,
	}

	seq := ddisc.Discover(
		ctx,
		req,
		ddisc.LocalStrategy(db),
		ddisctorrent.NewPartitionStrategy(dhts, partitions),
		ddisc.PluginStrategy(db, plugins),
	)

	found := false
	for d := range seq.Each(ctx) {
		found = true
		log.Println("found", d.ID, d.Title, d.Mimetype, d.KnownMediaID)
	}

	if err := seq.Err(); err != nil {
		return errorsx.Wrap(err, "discover failed")
	}

	if !found {
		log.Println("no results found")
	}

	return nil
}
