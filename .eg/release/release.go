package release

import (
	"context"
	"eg/compute/tarball"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
)

func Tarball(ctx context.Context, op eg.Op) error {
	archive := tarball.GitPattern("retrovibed")
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
			tarball.Archive(tarball.GitPattern("retrovibed")),
			egenv.CacheDirectory("flatpak.client.yml"),
			egenv.CacheDirectory("flatpak.daemon.yml"),
		),
	)
}
