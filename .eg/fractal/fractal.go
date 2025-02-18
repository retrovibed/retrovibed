package fractal

import (
	"context"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/shell"
)

func Build(ctx context.Context, _ eg.Op) error {
	return shell.Run(
		ctx,
		shell.New("which flutter"),
	)
}

func Tests(ctx context.Context, _ eg.Op) error {
	return shell.Run(
		ctx,
		shell.New("true"),
	)
}
