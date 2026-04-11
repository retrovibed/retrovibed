package cmdtorrent

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/bencode"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/james-lawrence/torrent/storage"
	"github.com/retrovibed/retrovibed/blockcache"
	"github.com/retrovibed/retrovibed/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/internal/asynccompute"
	"github.com/retrovibed/retrovibed/internal/env"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/fsx"
	"github.com/retrovibed/retrovibed/internal/langx"
)

type filesystemVerify struct{}

func (t filesystemVerify) Run(gctx *cmdopts.Global) (err error) {
	type workload struct {
		id   int160.T
		path string
	}
	const (
		suffix = ".torrent"
	)

	tvfs := fsx.DirVirtual(env.TorrentDir())

	torrentstore := fsx.DirVirtual(env.TorrentDir())

	mdcache := torrent.NewMetadataCache(tvfs.Path())
	bmcache := torrent.NewBitmapCache(tvfs.Path())

	var tstore storage.ClientImpl = blockcache.NewTorrentFromVirtualFS(torrentstore)

	root, err := os.OpenRoot(tvfs.Path())
	if err != nil {
		return err
	}

	dir := fsx.Walk(root.FS())

	perform := func(_ context.Context, w workload) error {
		var (
			id   = w.id
			path = w.path
			info metainfo.Info
		)

		md, err := mdcache.Read(id)
		if err != nil {
			return errorsx.Wrapf(err, "unable to read medata of %s", id)
		}

		if err = bencode.Unmarshal(md.InfoBytes, &info); err != nil {
			return errorsx.Wrapf(err, "unable to decode metainfo of %s", id)
		}

		disk, err := tstore.OpenTorrent(&info, id)
		if err != nil {
			return errorsx.Wrapf(err, "unable to open storage of %s", id)
		}

		missing, verified, err := torrent.VerifyStored(gctx.Context, langx.Autoptr(md.Metainfo()), disk)
		if err != nil {
			return errorsx.Wrapf(err, "unable to verify data of %s", id)
		}

		log.Printf("completed %s - %s - %s - missing(%d) - verified(%d)\n", path, id, md.DisplayName, missing.GetCardinality(), verified.GetCardinality())

		if err = bmcache.Write(id, verified); err != nil {
			return errorsx.Wrapf(err, "unable to persist bitmap of %s", id)
		}

		return nil
	}

	arena := asynccompute.New(func(ctx context.Context, w workload) error {
		if err := perform(ctx, w); err != nil {
			log.Println(err)
			return err
		}

		return nil
	})

	for path := range dir.Walk() {
		if !strings.HasSuffix(path, suffix) {
			continue
		}

		id, err := int160.FromHexEncodedString(strings.TrimSuffix(path, suffix))
		if err != nil {
			log.Printf("unable to read id from %s - %v\n", path, err)
			continue
		}

		if err = arena.Run(gctx.Context, workload{id: id, path: path}); err != nil {
			log.Printf("unable to perform verification from %s - %v\n", path, err)
			continue
		}
	}

	return errorsx.Compact(asynccompute.Shutdown(gctx.Context, arena), dir.Err())
}
