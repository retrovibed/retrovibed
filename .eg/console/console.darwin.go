package console

import (
	"context"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/shell"
)

func BuildDarwin(ctx context.Context, _ eg.Op) error {
	runtime := flutterRuntimev2(shell.Runtime())
	return shell.Run(
		ctx,
		runtime.New("rm -rf build/macos/{x64,arm64}/debug").Lenient(true),
		runtime.New("flutter build macos --release lib/main.dart"),
	)
}
