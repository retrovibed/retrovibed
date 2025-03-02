package release

import (
	"context"
	"eg/compute/tarball"
	"eg/compute/tarballs"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
)

func Tarball(ctx context.Context, op eg.Op) error {
	archive := tarballs.Retrovibed()
	return eg.Perform(
		ctx,
		tarball.Pack(archive),
		tarball.SHA256Op(archive),
	)
}

func Release(ctx context.Context, op eg.Op) error {
	return eg.Perform(
		ctx,
		tarball.Github(
			tarball.Archive(tarballs.Retrovibed()),
			egenv.CacheDirectory("flatpak.client.yml"),
			egenv.CacheDirectory("flatpak.daemon.yml"),
		),
	)
}
