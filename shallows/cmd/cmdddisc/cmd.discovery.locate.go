package cmdddisc

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/storage"
	"github.com/retrovibed/retrovibed/retroapi/ddiscapi"
	"github.com/retrovibed/retrovibed/retroapi/iterx"
	retronetx "github.com/retrovibed/retrovibed/retroapi/netx"
	"github.com/retrovibed/retrovibed/retroapi/searchplugin"
	"github.com/retrovibed/retrovibed/retroapi/userx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/cmd/retrovibe/daemons"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/ddisc/ddisctorrent"
	"github.com/retrovibed/retrovibed/shallows/internal/cryptox"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/netipx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/internal/torrentx"
)

// searchPlugins mirrors daemons.searchPlugins' method set structurally (that
// type is unexported, so we can't reference it by name) - the narrow
// interface daemons.LocateMedia actually needs from *searchplugin.Registry.
type searchPlugins interface {
	Search(ctx context.Context, category, query string) iterx.Seq[*ddiscapi.Import]
}

// errLocateFound stops the retry loop in cmdMediaLocate.Run once the queued
// locate request has been resolved to a downloaded torrent.
var errLocateFound = errorsx.String("locate: candidate found")

type cmdMediaLocate struct {
	Database     string        `flag:"" name:"database" help:"database to read/write" default:"${vars_user_configuration_directory}/meta.db"`
	Query        string        `arg:"" name:"query" help:"title or free-text search to locate media for"`
	Mimetype     string        `flag:"" name:"mimetype" help:"mimetype/category (video, audio, image, text, application) to search for" required:""`
	KnownMediaID string        `flag:"" name:"known-media-id" help:"known media id, if already resolved from a catalog search" optional:""`
	Partitions   uint8         `flag:"" name:"discovery-partition" help:"number of partitions to split the infohash space into, must match the swarm's configuration" default:"128"`
	Seed         string        `flag:"" name:"discovery-seed" help:"seed to generate partition spaces, must match the swarm's configuration" default:"retrovibed-ddisc"`
	Interval     time.Duration `flag:"" name:"interval" help:"how often to re-run the discover/rank/download pass while waiting on async DHT responses" default:"20s"`
	Timeout      time.Duration `flag:"" name:"timeout" help:"give up waiting for a candidate after this long; the locate request itself remains queued for a running daemon to pick up later" default:"2m"`
	Bootstrap    bool          `flag:"" name:"dht-bootstrap" help:"bootstrap the DHT using well-known trackers" negatable:"" default:"true"`
	DHTPeers     []string      `flag:"" name:"dht-peers" help:"use these dht peers as the sole bootstrap nodes instead of the public network" hidden:"true"`
}

func (t cmdMediaLocate) Run(gctx *cmdopts.Global) (err error) {
	db, err := cmdopts.DatabaseCustom(gctx.Context, t.Database)
	if err != nil {
		return err
	}
	defer db.Close()

	loc := ddisc.NewLocate(t.Query, t.Mimetype, ddisc.LocateOptionKnownMedia(langx.FirstNonZero(t.KnownMediaID, uuid.Nil.String())))
	if err := ddisc.LocateInsertWithDefaults(gctx.Context, db, loc).Scan(&loc); err != nil {
		return errorsx.Wrap(err, "failed to queue locate request")
	}

	log.Println("locate request queued", loc.ID, loc.Query)

	dhts, tclient, err := t.torrentClient(db)
	if err != nil {
		return errorsx.Wrap(err, "unable to setup torrent client")
	}
	defer tclient.Close()

	var plugins searchPlugins
	if reg, err := searchplugin.NewRegistry(gctx.Context); err != nil {
		log.Println("search plugins unavailable, continuing with DHT discovery only:", err)
	} else {
		plugins = reg
	}

	partitions := ddisc.Partitions(uint16(t.Partitions), cryptox.NewChaCha8(t.Seed))
	policy := ddisc.DefaultPolicy()
	disc := &daemons.DiscoverySettings{
		LocateP2P:  true,
		Partitions: uint32(t.Partitions),
		Seed:       t.Seed,
	}

	bctx, bcancel := context.WithTimeout(gctx.Context, t.Timeout)
	defer bcancel()

	err = timex.NowAndEvery(bctx, t.Interval, func(ctx context.Context) error {
		if err := daemons.LocateMedia(ctx, db, tclient, disc, dhts, partitions, plugins, policy); err != nil {
			return err
		}

		var found ddisc.Locate
		if err := ddisc.LocateFindByID(ctx, db, loc.ID).Scan(&found); err != nil {
			return err
		}

		if found.LocatedTorrentID != uuid.Max.String() {
			log.Println("media located", found.ID, found.LocatedTorrentID)
			return errLocateFound
		}

		return nil
	})

	if errors.Is(err, errLocateFound) || errors.Is(err, context.DeadlineExceeded) {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("timed out waiting for a candidate; the locate request remains queued for a running daemon to pick up later")
		}
		return nil
	}

	return err
}

func (t cmdMediaLocate) torrentClient(db sqlx.Queryer) (dhts *dht.Server, tclient *torrent.Client, err error) {
	var (
		globalbootstrap dht.Option = dht.OptionNoop
	)

	cachedir := userx.DefaultCacheDirectory("torrentddisc")

	if t.Bootstrap {
		globalbootstrap = dht.OptionBootstrapGlobal
	}

	muxer := dht.DefaultMuxer().Method(ddisctorrent.MethodMedia, ddisctorrent.NewMediaRecorder(db))

	dhtOptions := []dht.Option{
		dht.OptionMuxer(muxer),
		dht.OptionBootstrapAddrPort(netipx.AddrPortFromStrings(t.DHTPeers...)...),
		globalbootstrap,
		dht.OptionUPnP,
	}

	if dhts, err = dht.NewServer(32, dhtOptions...); err != nil {
		return nil, nil, errorsx.Wrap(err, "failed to initialize dht")
	}
	go dhts.TableMaintainer(context.Background())

	tnetwork, err := torrentx.Autosocket(dhts, 0, retronetx.NewConnUnlimited())
	if err != nil {
		return nil, nil, errorsx.Wrap(err, "unable to setup torrent socket")
	}

	ttstore := storage.NewFile(cachedir)

	torconfig := torrent.NewDefaultClientConfig(
		torrent.NewMetadataCache(cachedir),
		ttstore,
		torrent.ClientConfigCacheDirectory(cachedir),
		torrent.ClientConfigExtension(ddisctorrent.ExtensionName),
		torrent.ClientConfigSeed(false),
		torrent.ClientConfigPEX(false),
		torrent.ClientConfigHTTPUserAgent("retrovibed/0.0"),
	)

	if tclient, err = tnetwork.Bind(torrent.NewClient(torconfig)); err != nil {
		return nil, nil, errorsx.Wrap(err, "unable to bind torrent to socket")
	}

	return dhts, tclient, nil
}
