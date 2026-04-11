package cmdtorrent

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/internal/env"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/fsx"
	"github.com/retrovibed/retrovibed/internal/stringsx"
)

type exportMagnets struct {
	Path   string `arg:"" name:"path" help:"file to write magnet urls out to, defaults to stdout" default:"" required:"false"`
	Legacy bool   `flag:"" name:"legacy" help:"enable legacy storage structure for migration" default:"false" hidden:"true"`
}

func (t exportMagnets) Run(gctx *cmdopts.Global, id *cmdopts.SSHID) (err error) {
	const (
		suffix = ".torrent"
	)

	dst := os.Stdout
	if stringsx.Present(t.Path) {
		if dst, err = os.OpenFile(t.Path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600); err != nil {
			return err
		}
	}

	tvfs := fsx.DirVirtual(env.TorrentDir())

	mdcache := torrent.NewMetadataCache(tvfs.Path())

	root, err := os.OpenRoot(tvfs.Path())
	if err != nil {
		return err
	}

	dir := fsx.Walk(root.FS())
	bmc := torrent.NewBitmapCache(tvfs.Path())

	for path := range dir.Walk() {
		if !strings.HasSuffix(path, suffix) {
			continue
		}

		id, err := int160.FromHexEncodedString(strings.TrimSuffix(path, suffix))
		if err != nil {
			return errorsx.Wrapf(err, "unable to read id from %s", path)
		}

		if !fsx.Exists(tvfs.Path(id.String())) {
			continue
		}

		bm, err := bmc.Read(id)
		if err != nil {
			log.Println(errorsx.Wrapf(err, "unable to open bitmap cache %s", id))
			continue
		}

		if bm.GetCardinality() == 0 {
			log.Printf("bitmap - skipping due to no data available %s\n", id)
			continue
		}

		md, err := mdcache.Read(id)
		if err != nil {
			return errorsx.Wrapf(err, "unable to read medata of %s", id)
		}

		mg := metainfo.Magnet{
			InfoHash:    metainfo.Hash(md.ID.Bytes()),
			Trackers:    md.Trackers,
			DisplayName: md.DisplayName,
		}

		if _, err = fmt.Fprintln(dst, mg.String()); err != nil {
			return err
		}
	}

	return nil
}
