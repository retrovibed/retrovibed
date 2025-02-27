package release

import (
	"context"
	"os"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/egflatpak"
	"github.com/egdaemon/eg/runtime/x/wasi/egfs"
)

func Flatpak(ctx context.Context, op eg.Op) error {
	runtime := shell.Runtime()

	builddir := egenv.WorkingDirectory("fractal", "build", egfs.FindFirst(os.DirFS(egenv.WorkingDirectory("fractal", "build")), "bundle"))

	b := egflatpak.New("space.retrovibe.Daemon", egflatpak.Option.CopyModule(builddir)...)
	if err := egflatpak.Build(ctx, runtime, b); err != nil {
		return err
	}

	return nil
}
