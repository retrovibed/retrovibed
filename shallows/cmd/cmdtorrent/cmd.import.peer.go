package cmdtorrent

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/bep0051"
	"github.com/james-lawrence/torrent/connections"
	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/retroapi/blockcache"
	retronetx "github.com/retrovibed/retrovibed/retroapi/netx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/cmd/retrovibe/daemons"
	"github.com/retrovibed/retrovibed/shallows/internal/asynccompute"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/envx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/internal/torrentx"
	"github.com/retrovibed/retrovibed/shallows/internal/userx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

type logging interface {
	Println(v ...any)
	Printf(format string, v ...any)
	Print(v ...any)
}

type importPeer struct {
	Peer           []string `flag:"" name:"peer" help:"peer(s) to connect to and download the provided torrents from" default:"localhost:10000"`
	Directory      string   `flag:"" name:"directory" help:"specify the directory to download torrents into" default:""`
	Archive        bool     `flag:"" name:"archive" help:"mark imported media for archival" default:"false"`
	Reclaim        bool     `flag:"" name:"reclaim" help:"mark imported media for disk space reclaimation" default:"false"`
	Magnets        string   `arg:"" name:"magnets" help:"file containing magnet links to download, defaults to stdin" default:""`
	TorrentPrivate bool     `flag:"" name:"torrent-private" help:"restrict torrent connections to private networks" env:"${env_torrent_private}" default:"false"`
	Concurrency    uint16   `flag:"" name:"concurrency" help:"specify the number of workers" default:"${vars_cores}"`
}

func (t importPeer) torrents(tstore fsx.Virtual) iter.Seq2[torrentRecord, torrent.Metadata] {
	return func(yield func(torrentRecord, torrent.Metadata) bool) {
		src := os.Stdin

		if stringsx.Present(t.Magnets) {
			var (
				err error
			)

			if src, err = os.Open(t.Magnets); err != nil {
				log.Fatalln("failed to open magnets source", err)
			}
		}

		dec := json.NewDecoder(src)
		for dec.More() {
			var rec torrentRecord
			if err := dec.Decode(&rec); err != nil {
				log.Println("unable to decode torrent record", err)
				return
			}

			magnet, err := torrent.NewFromMagnet(rec.Magnet)
			if err != nil {
				log.Println("unable to parse magnet link", err)
				return
			}

			infopath := tstore.Path(fmt.Sprintf("%s.torrent", magnet.ID))
			if _, err := os.Stat(infopath); err == nil {
				if magnet.InfoBytes, err = os.ReadFile(infopath); err != nil {
					log.Println("failed to read torrent info", err)
				}
			}

			if !yield(rec, magnet) {
				return
			}
		}
	}
}

func (t importPeer) Run(gctx *cmdopts.Global, sshid *cmdopts.SSHID) (err error) {
	type workload struct {
		record torrentRecord
		meta   torrent.Metadata
	}

	var (
		db        *sql.DB
		tnetwork  torrent.Binder
		bootstrap torrent.ClientConfigOption = torrent.ClientConfigNoop
		firewall  torrent.ClientConfigOption = torrent.ClientConfigNoop
	)

	gctx.Cleanup.Add(1)
	defer gctx.Cleanup.Done()
	if db, err = cmdopts.DatabaseMeta(gctx.Context); err != nil {
		return err
	}
	defer db.Close()

	rootstore := fsx.DirVirtual(userx.DefaultDataDirectory(userx.DefaultRelRoot()))
	torrentstore := fsx.DirVirtual(stringsx.FirstNonBlank(t.Directory, env.TorrentDir()))
	mediastore := fsx.DirVirtual(env.MediaDir())

	if err := fsx.MkDirs(0700, torrentstore.Path(), mediastore.Path()); err != nil {
		return err
	}

	async := asyncx.NewWakeup(gctx.Context)
	errorsx.Log(daemons.AutoArchival(gctx.Context, db, mediastore, async, t.Reclaim))

	peers := make([]torrent.Peer, 0, 128)
	for _, p := range t.Peer {
		host, port, err := net.SplitHostPort(p)
		if err != nil {
			return errorsx.Wrap(err, "unable to setup torrent client")
		}

		addrs, err := net.DefaultResolver.LookupIP(gctx.Context, "ip4", host)
		if err != nil {
			return errorsx.Wrap(err, "unable to resolve host")
		}

		peers = append(peers, torrent.NewPeerDeprecated(
			int160.Zero(),
			addrs[0],
			uint16(errorsx.Must(strconv.Atoi(port))), torrent.PeerOptionTrusted(true),
		))
	}

	if t.TorrentPrivate {
		log.Println("disabling public networks for torrent")
		firewall = torrent.ClientConfigFirewall(connections.NewFirewall(
			connections.Private{},
			connections.BanInvalidPort{},
			connections.NewBloomBanIP(10*time.Minute),
		))
	}

	var torrentlogging logging = torrent.LogDiscard()
	if envx.Boolean(false, env.TorrentLogging) {
		torrentlogging = log.New(os.Stderr, "[torrent] ", log.Flags())
	}
	tm := dht.DefaultMuxer().
		Method(bep0051.Query, bep0051.NewEndpoint(bep0051.EmptySampler{}))
	_dht, err := dht.NewServer(
		32,
		dht.OptionMuxer(tm),
	)
	if err != nil {
		return errorsx.Wrap(err, "failed to initialize dht")
	}

	tstore := blockcache.NewTorrentFromVirtualFS(torrentstore)

	torconfig := torrent.NewDefaultClientConfig(
		torrent.NewMetadataCache(torrentstore.Path()),
		tstore,
		torrent.ClientConfigCacheDirectory(torrentstore.Path()),
		torrent.ClientConfigPEX(false),
		torrent.ClientConfigSeed(false),
		torrent.ClientConfigInfoLogger(torrentlogging),
		torrent.ClientConfigDebugLogger(torrentlogging),
		torrent.ClientConfigMaxOutstandingRequests(2048),
		torrent.ClientConfigHTTPUserAgent("retrovibed/0.0"),
		torrent.ClientConfigConnectionClosed(func(ih int160.T, stats torrent.ConnStats, remaining int) {
			if stats.BytesWrittenData.Uint64() == 0 {
				return
			}

			var md tracking.Metadata
			ictx, done := context.WithTimeout(gctx.Context, 3*time.Second)
			defer done()
			if err := tracking.MetadataUploadedByID(ictx, db, ih.Bytes(), stats.BytesWrittenData.Uint64()).Scan(&md); err != nil {
				log.Println(errorsx.Wrapf(err, "%s: unable to record uploaded metrics", ih.String()))
				return
			}

			if remaining == 0 {
				time.AfterFunc(time.Minute, func() {
					log.Println("connection closed, and no remaining connections, TODO gracefully remove")
				})
			}
		}),
		bootstrap,
		firewall,
	)

	if tnetwork, err = torrentx.Autosocket(_dht, 0, retronetx.NewConnUnlimited()); err != nil {
		return errorsx.Wrap(err, "unable to setup torrent socket")
	}

	tclient, err := tnetwork.Bind(torrent.NewClient(torconfig))
	if err != nil {
		return errorsx.Wrap(err, "unable to bind torrent to socket")
	}
	defer tclient.Close()

	importfn := func(ctx context.Context, w workload) (err error) {
		var (
			info *metainfo.Info
		)

		log.Println("import initiated", w.meta.ID)
		defer log.Println("import completed", w.meta.ID)
		defer func() {
			if err == nil {
				return
			}
		}()

		if len(w.meta.InfoBytes) == 0 {
			var (
				cause error
			)
			log.Printf("awaiting torrent info %s\n", w.meta.ID)

			info, cause = tclient.Info(ctx, w.meta, torrent.TunePeers(peers...))
			if cause != nil {
				return errorsx.Wrapf(cause, "failed to retrieve torrent info %s", w.meta.ID)
			}

			log.Println("torrent info received", w.meta.ID.String(), spew.Sdump(w.meta.Metainfo()))
			if w.meta, cause = torrent.NewFromInfo(info); err != nil {
				return errorsx.Wrapf(cause, "failed to retrieve torrent info %s", w.meta.ID)
			}

			if err = os.WriteFile(torrentstore.Path(fmt.Sprintf("%s.torrent", w.meta.ID)), errorsx.Must(metainfo.Encode(w.meta.Metainfo())), 0600); err != nil {
				return errorsx.Wrapf(cause, "failed to record torrent %s %v", w.meta.ID, cause)
			}
		} else {
			if _info, cause := w.meta.Metainfo().UnmarshalInfo(); cause != nil {
				return errorsx.Wrapf(cause, "failed to resume torrent %s", w.meta.ID)
			} else {
				info = &_info
			}
		}

		lmd := tracking.NewMetadata(
			langx.Autoptr(w.meta.ID),
			tracking.MetadataOptionFromInfo(info),
			tracking.MetadataOptionTrackers(w.meta.Trackers...),
			tracking.MetadataOptionAutoArchive(t.Archive),
			tracking.MetadataOptionEncryptionSeed(w.record.EncryptionSeed),
			tracking.MetadataOptionInitiate,
			tracking.MetadataOptionAutoDescription,
			tracking.MetadataOptionAutoHidden,
		)

		if cause := tracking.MetadataInsertWithDefaults(ctx, db, lmd).Scan(&lmd); cause != nil {
			return errorsx.Wrapf(cause, "failed to record metadata %s", w.meta.ID.String())
		}

		log.Println("---------------------------------- starting", w.meta.ID, lmd.ID, lmd.Downloaded, lmd.Bytes)
		dl, _, cause := tclient.Start(w.meta, torrent.TuneDisableTrackers, torrent.TunePeers(peers...), torrent.TuneVerifyFull)
		if cause != nil {
			return errorsx.Wrapf(cause, "failed to start magnet %s - %T: %+v\n", w.meta.ID.String(), cause, cause)
		}
		defer func() {
			errorsx.Log(errorsx.Wrapf(tclient.Stop(w.meta), "failed to shutdown torrent %s %v", w.meta.ID.String(), cause))
		}()

		log.Println("---------------------------------- downloading", w.meta.ID)
		if cause := tracking.DownloadInto(ctx, db, rootstore, library.NewQueryerCleanerAuto(), &lmd, dl, io.Discard); cause != nil {
			return errorsx.LogErr(errorsx.Wrapf(cause, "failed to download %s %v", w.meta.ID.String(), cause))
		}

		if cause := tclient.Stop(w.meta); cause != nil {
			return errorsx.LogErr(errorsx.Wrapf(cause, "failed to shutdown torrent %s %v", w.meta.ID.String(), cause))
		}

		log.Println("---------------------------------- COMPLETED", w.meta.ID, lmd.Downloaded, lmd.Bytes)
		if t.Archive {
			async.Broadcast()
		}

		return nil
	}

	arena := asynccompute.New(func(ctx context.Context, w workload) error {
		if err := importfn(ctx, w); err != nil {
			fmt.Println(w.record.Magnet)
			return err
		}

		return nil
	}, asynccompute.Workers[workload](t.Concurrency))

	for rec, v := range t.torrents(torrentstore) {
		if err := arena.Run(gctx.Context, workload{record: rec, meta: v}); err != nil {
			return errorsx.Compact(err, asynccompute.Shutdown(gctx.Context, arena))
		}
	}

	if err = asynccompute.Shutdown(gctx.Context, arena); err != nil {
		return err
	}

	return async.Close()
}
