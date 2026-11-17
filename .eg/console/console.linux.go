package console

import (
	"context"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/shell"
)

func CompileBinding(ctx context.Context, o eg.Op) error {
	runtime := flutterRuntimev2(shell.Runtime())

	return shell.Run(
		ctx,
		runtime.New("mkdir -p ${RETROVIBED_SHARED_NATIVE_LIBS_DIRECTORY} ${GOTMPDIR}"),
		runtime.New("go -C retrovibedbind build -buildmode=c-shared -buildvcs=true --tags duckdb_use_lib,retrovibed,neural -o ${RETROVIBED_SHARED_NATIVE_LIBS_DIRECTORY}/libretrovibed.so .").
			Timeout(egenv.TTL()).
			Environ("CGO_LDFLAGS", "-L"+egenv.DirectoryCache("neurals")),
	)
}

func BuildLinux(ctx context.Context, _ eg.Op) error {
	runtime := flutterRuntimev2(shell.Runtime())
	return shell.Run(
		ctx,
		runtime.New("flutter build linux --release lib/main.dart").Timeout(egenv.TTL()),
	)
}
