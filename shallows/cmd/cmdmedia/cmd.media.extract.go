package cmdmedia

import (
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

type mediaInspectFile struct {
	Path string `arg:"" name:"path" help:"file to look at" required:"true"`
}

func (t mediaInspectFile) Run(gctx *cmdopts.Global) (err error) {
	src, err := os.Open(t.Path)
	if err != nil {
		return errorsx.Wrap(err, "unable to open file")
	}
	defer src.Close()

	e, err := ddisc.Extract(src)
	if err != nil {
		return errorsx.Wrap(err, "unable to determine mimetype")
	}

	log.Println(spew.Sdump(e))

	return nil
}
