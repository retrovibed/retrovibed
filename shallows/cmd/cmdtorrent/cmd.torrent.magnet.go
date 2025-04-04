package cmdtorrent

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/bencode"
	"github.com/retrovibed/retrovibed/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/userx"
)

type cmdMagnet struct {
	Magnets []url.URL `arg:"" name:"magnet" help:"magnet uri to download" required:"true"`
}

func (t cmdMagnet) Run(ctx *cmdopts.Global) (err error) {
	for _, uri := range t.Magnets {
		m, cause := torrent.NewFromMagnet(uri.String())
		if cause != nil {
			err = errors.Join(err, errorsx.Wrap(cause, "unable to prepare magnet"))
			continue
		}

		encoded, cause := bencode.Marshal(m.Metainfo())
		if cause != nil {
			err = errors.Join(err, errorsx.Wrap(cause, "unable to encode to torrent file"))
			continue
		}

		path := userx.DefaultDownloadDirectory(fmt.Sprintf("%s.torrent", m.ID.HexString()))
		if cause := os.WriteFile(path, encoded, 0600); cause != nil {
			err = errors.Join(err, errorsx.Wrap(cause, "unable to write torrent file"))
			continue
		}

		if cause := os.Chmod(path, 0600); cause != nil {
			err = errors.Join(err, errorsx.Wrap(cause, "unable to touch torrent file"))
			continue
		}

		log.Println("NOOP MAGNET - Not implemented", userx.DefaultDownloadDirectory(fmt.Sprintf("%s.torrent", m.ID.HexString())), spew.Sdump(m))
	}

	return err
}
