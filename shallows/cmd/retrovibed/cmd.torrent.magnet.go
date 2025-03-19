package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/james-lawrence/deeppool/cmd/cmdopts"
	"github.com/james-lawrence/deeppool/internal/x/errorsx"
	"github.com/james-lawrence/deeppool/internal/x/userx"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/bencode"
)

type cmdTorrent struct {
	Magnet cmdTorrentMagnet `cmd:"" help:"insert magnet links for download"`
}

type cmdTorrentMagnet struct {
	Magnets []url.URL `arg:"" name:"magnet" help:"magnet uri to download" required:"true"`
}

func (t cmdTorrentMagnet) Run(ctx *cmdopts.Global) (err error) {
	for _, uri := range t.Magnets {
		m, cause := torrent.NewFromMagnet(uri.String())
		if cause != nil {
			err = errors.Join(err, errorsx.Wrap(cause, "unable to prepare magnet"))
			continue
		}

		encoded, err := bencode.Marshal(m.Metainfo())
		if cause != nil {
			err = errors.Join(err, errorsx.Wrap(cause, "unable to encode to torrent file"))
			continue
		}

		path := userx.DefaultDownloadDirectory(fmt.Sprintf("%s.torrent", m.InfoHash.HexString()))
		if cause := os.WriteFile(path, encoded, 0600); cause != nil {
			err = errors.Join(err, errorsx.Wrap(cause, "unable to write torrent file"))
			continue
		}

		if cause := os.Chmod(path, 0600); cause != nil {
			err = errors.Join(err, errorsx.Wrap(cause, "unable to touch torrent file"))
			continue
		}

		log.Println("NOOP MAGNET - Not implemented", userx.DefaultDownloadDirectory(fmt.Sprintf("%s.torrent", m.InfoHash.HexString())), spew.Sdump(m))
	}

	return err
}
