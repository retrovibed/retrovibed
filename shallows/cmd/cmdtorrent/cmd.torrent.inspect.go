package cmdtorrent

import (
	"github.com/alecthomas/kong"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/metainfox"
)

type cmdInspect struct {
	Files []string `arg:"" name:"files" help:"set of torrent files to read and print details for" required:"true"`
}

func (t cmdInspect) Run(kctx *kong.Context, gctx *cmdopts.Global) error {
	for _, path := range t.Files {
		if err := gctx.Context.Err(); err != nil {
			return err
		}

		md, err := metainfo.LoadFromFile(path)
		if err != nil {
			return errorsx.Wrapf(err, "failed to read file: %s", path)
		}

		if err := metainfox.NewPrinter(md).Print(kctx.Stdout); err != nil {
			return errorsx.Wrapf(err, "failed to format torrent: %s", path)
		}
	}

	return nil
}
