package release

import (
	"context"
	"eg/compute/tarball"

	"github.com/egdaemon/eg/runtime/wasi/eg"
)

func Tarball(ctx context.Context, op eg.Op) error {
	return eg.Perform(
		ctx,
		tarball.Archive(tarball.Directory()),
	)
}
