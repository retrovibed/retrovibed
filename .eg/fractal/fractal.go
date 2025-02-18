package fractal

import (
	"context"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/shell"
)

func flutterRuntime() shell.Command {
	return shell.Runtime().Directory(egenv.WorkingDirectory("fractal")).Environ("PUB_CACHE", egenv.CacheDirectory(".eg", "dart"))
}

func Build(ctx context.Context, _ eg.Op) error {
	runtime := flutterRuntime()
	return shell.Run(
		ctx,
		runtime.New("flutter create --platforms=linux ."),
		runtime.Newf("flutter build linux"),
	)
}

func Tests(ctx context.Context, _ eg.Op) error {
	runtime := flutterRuntime()
	return shell.Run(
		ctx,
		runtime.New("flutter test"),
	)
}

func Linting(ctx context.Context, _ eg.Op) error {
	runtime := flutterRuntime()
	return shell.Run(
		ctx,
		runtime.New("flutter analyze"),
	)
}

func Generate(ctx context.Context, _ eg.Op) error {
	runtime := flutterRuntime()
	return shell.Run(
		ctx,
		runtime.New("echo generating"),
	)
}
