package release

import (
	"context"
	"eg/compute/tarballs"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/x/wasi/egtarball"
)

func Tarball(ctx context.Context, op eg.Op) error {
	archive := tarballs.Retrovibed()
	return eg.Perform(
		ctx,
		egtarball.Pack(archive),
		egtarball.SHA256Op(archive),
	)
}

func Release(ctx context.Context, op eg.Op) error {
	return eg.Perform(
		ctx,
		egtarball.Github(
			egtarball.Archive(tarballs.Retrovibed()),
			egenv.CacheDirectory("flatpak.client.yml"),
			egenv.CacheDirectory("flatpak.daemon.yml"),
		),
	)
}
