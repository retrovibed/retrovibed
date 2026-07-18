package cmdddisc

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/dht/int160"
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
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/netipx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/internal/torrentx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

// searchPlugins mirrors daemons.searchPlugins' method set structurally (that
// type is unexported, so we can't reference it by name) - the narrow
// interface daemons.LocateMedia actually needs from *searchplugin.Registry.
type searchPlugins interface {
	Search(ctx context.Context, mimetypes []string, query string, adult bool) iterx.Seq[*ddiscapi.Import]
}

// errLocateFound stops the retry loop in cmdMediaLocate.Run once the queued
// locate request has had a candidate ranked and selected.
var errLocateFound = errorsx.String("locate: candidate found")

type cmdMediaLocate struct {
	Database     string        `flag:"" name:"database" help:"database to read/write" default:"${vars_user_configuration_directory}/meta.db"`
	Query        string        `arg:"" name:"query" help:"title or free-text search to locate media for"`
	Mimetype     string        `flag:"" name:"mimetype" help:"mimetype/category (video, audio, image, text, application) to search for" required:""`
	KnownMediaID string        `flag:"" name:"known-media-id" help:"known media id, if already resolved from a catalog search" optional:""`
	Partitions   uint8         `flag:"" name:"discovery-partition" help:"number of partitions to split the infohash space into, must match the swarm's configuration" default:"128"`
	Seed         string        `flag:"" name:"discovery-seed" help:"seed to generate partition spaces, must match the swarm's configuration" default:"retrovibed-ddisc"`
	Interval     time.Duration `flag:"" name:"interval" help:"how often to re-run the discover/rank pass while waiting on async DHT responses" default:"20s"`
	Timeout      time.Duration `flag:"" name:"timeout" type:"durationinf" help:"give up waiting for a candidate after this long, use 'infinity' to wait forever; the locate request itself remains queued for a running daemon to pick up later" default:"infinity"`
	Bootstrap    bool          `flag:"" name:"dht-bootstrap" help:"bootstrap the DHT using well-known trackers" negatable:"" default:"true"`
	Adult        bool          `flag:"" name:"adult" help:"allow adult content in search plugin results" negatable:"" default:"false"`
	DHTPeers     []string      `flag:"" name:"dht-peers" help:"use these dht peers as the sole bootstrap nodes instead of the public network" hidden:"true"`
}

func (t cmdMediaLocate) Run(gctx *cmdopts.Global) (err error) {
	db, err := cmdopts.DatabaseCustom(gctx.Context, t.Database)
	if err != nil {
		return err
	}
	defer db.Close()

	loc := ddisc.NewLocate(t.Query, t.Mimetype, ddisc.LocateOptionKnownMedia(langx.FirstNonZero(t.KnownMediaID, uuid.Nil.String())), ddisc.LocateOptionAdult(t.Adult))
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

	var found ddisc.Discovered
	err = timex.NowAndEvery(bctx, t.Interval, func(ctx context.Context) error {
		d, ferr := daemons.Locate(ctx, db, disc, dhts, partitions, plugins, policy, loc)
		if errors.Is(ferr, ddisc.ErrNoCandidate) {
			return nil
		} else if ferr != nil {
			return ferr
		}

		found = d
		return errLocateFound
	})

	if errors.Is(err, errLocateFound) {
		log.Println("candidate found", found.ID, found.Title, found.PolicyRank, int160.FromBytes(found.Infohash))
		importer := tracking.NewURIImport(db, http.DefaultClient, fsx.DirVirtual(userx.DefaultCacheDirectory("torrentddisc")))
		if err := daemons.DiscoveredDownload(bctx, db, importer, loc, found); err != nil {
			return errorsx.Wrap(err, "unable to initiate download")
		}
		log.Println("media located", loc.ID, found.ID)
		return nil
	} else if errors.Is(err, context.DeadlineExceeded) {
		log.Println("timed out waiting for a candidate; the locate request remains queued for a running daemon to pick up later")
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
