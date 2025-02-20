package shallows

import (
	"context"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/eggolang"
)

func runtime() shell.Command {
	return eggolang.Runtime().Directory(egenv.WorkingDirectory("shallows"))
}

func Generate(ctx context.Context, _ eg.Op) error {
	gruntime := runtime()
	return shell.Run(
		ctx,
		gruntime.New("go generate ./... && go fmt ./..."),
	)
}
