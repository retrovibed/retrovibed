package ddisc

import (
	"context"
	"encoding/hex"
	"log"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/james-lawrence/torrent/storage"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/torrentx"
)

// IdentifyOne locates peers for disc's infohash, downloads a sample of the torrent,
// and detects its mimetype. peerTimeout bounds the DHT peer lookup; infoTimeout bounds
// the wait for torrent metadata/info once peers are being searched for. It performs no
// database writes - callers are responsible for persisting the returned Discovered value.
func IdentifyOne(ctx context.Context, dhts *dht.Server, tclient *torrent.Client, ttstore storage.ClientImpl, peerTimeout, infoTimeout time.Duration, disc Discovered, tuners ...torrent.Tuner) (_ Discovered, err error) {
	available := func(ctx context.Context, disc Discovered) ([]dht.Peer, error) {
		dctx, done := context.WithTimeout(ctx, peerTimeout)
		defer done()
		return torrentx.Peers(dctx, dhts, int160.FromBytes(disc.Infohash))
	}

	peers, err := available(ctx, disc)
	if err != nil {
		return disc, errorsx.Wrapf(err, "no peers available %s", disc.ID)
	}
	log.Println("located", len(peers), "initial peers")

	dctx, done := context.WithTimeout(ctx, infoTimeout)
	defer done()

	log.Println("identify initiated", disc.ID, hex.EncodeToString(disc.Infohash))
	defer log.Println("identify completed", disc.ID, hex.EncodeToString(disc.Infohash))

	id := int160.FromBytes(disc.Infohash)
	metadata, err := torrent.New(metainfo.Hash(disc.Infohash), torrent.OptionStorage(ttstore))
	if err != nil {
		return disc, errorsx.Wrapf(err, "unable to create metadata from infohash %s", disc.ID)
	}
	defer tclient.Stop(metadata)

	info, err := tclient.Info(
		dctx,
		metadata,
		append([]torrent.Tuner{torrent.TuneAnnounceUntilComplete, torrent.TuneNewConns}, tuners...)...,
	)
	if err != nil {
		return disc, errorsx.Wrapf(err, "unable to initialize torrent for %s", disc.ID)
	}

	torrentx.FilePrint(
		info,
		torrentx.FileFirst(info, func(fi metainfo.FileInfo) bool { return strings.HasSuffix(fi.DisplayPath(info), ".nfo") }),
		torrentx.FileLargest(info),
	)

	_, off, length := torrentx.FileLargestRange(info)
	metadata, err = torrent.NewFromInfo(info, torrent.OptionStorage(ttstore))
	if err != nil {
		return disc, errorsx.Wrapf(err, "unable to create metadata from info %s", id)
	}

	tt, _, err := tclient.Start(
		metadata,
		append([]torrent.Tuner{
			torrent.TuneResetBitmaps,
			torrent.TuneVerifyRange(off, min(length, 32*bytesx.KiB)),
			torrent.TuneAnnounceUntilComplete,
			torrent.TuneNewConns,
		}, tuners...)...,
	)
	if err != nil {
		return disc, errorsx.Wrapf(err, "unable to initialize torrent for %s", disc.ID)
	}

	log.Println("attempting to download", id, bytesx.Unit(length))
	r := torrent.DownloadRange(dctx, tt, off, min(length, 32*bytesx.KiB))
	go func() {
		torrentx.DownloadProgress(dctx, tt)
		errorsx.Log(errorsx.Wrapf(r.Close(), "%s failed to close reader", id))
		errorsx.Log(errorsx.Wrapf(tclient.Stop(metadata), "%s failed to stop torrent", id))
	}()

	mime, err := Mimetype(r)
	if err != nil {
		return disc, errorsx.Wrapf(err, "unable to initialize torrent for %s", disc.ID)
	}

	disc = langx.Clone(
		disc,
		DiscoveredOptionMimetype(mime.String()),
		DiscoveredOptionKnownMedia(uuid.Max.String()),
	)

	return disc, nil
}
