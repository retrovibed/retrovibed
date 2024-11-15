package main

import (
	"errors"
	"log"
	"net/url"

	"github.com/davecgh/go-spew/spew"
	"github.com/james-lawrence/deeppool/cmd/cmdopts"
	"github.com/james-lawrence/deeppool/internal/x/errorsx"
	"github.com/james-lawrence/torrent"
)

type cmdTorrent struct {
	Magnet cmdTorrentMagnet `cmd:"" help:"insert magnet links for download"`
}

type cmdTorrentMagnet struct {
	Magnets []url.URL `arg:"" name:"megnet" help:"magnet uri to download" required:"true"`
}

func (t cmdTorrentMagnet) Run(ctx *cmdopts.Global) (err error) {
	for _, uri := range t.Magnets {
		m, cause := torrent.NewFromMagnet(uri.String())
		if cause != nil {
			err = errors.Join(err, errorsx.Wrap(cause, "unable to prepare magnet"))
			continue
		}
		log.Println("NOOP MAGNET - Not implemented", spew.Sdump(m))
	}
	return err
}
