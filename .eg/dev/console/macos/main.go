package main

import (
	"context"
	"eg/compute/console"
	"log"
	"path/filepath"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/egbug"
)

func Debug(runtime shell.Command) eg.OpFn {
	return shell.Op(
		runtime.New("which go && go version"),
		runtime.New("which flutter && flutter --version"),
		runtime.New("which dart && dart --version"),
		runtime.New("ls -la /Library/Developer/CommandLineTools/usr/lib/clang/").Lenient(true),
		runtime.New("xcode-select -p"),
	)
}

// Baremetal command for macOS local development.
// Sets up the native library and ffigen bindings so flutter can run locally.
func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	ctx, done := context.WithTimeout(context.Background(), egenv.TTL())
	defer done()

	runtime := shell.Runtime()

	err := eg.Perform(
		ctx,
		eg.Sequential(
			shell.Op(
				shell.New("brew install go duckdb gpgme flutter ffmpeg@7 cocoapods"),
			),
			egbug.DebugFailure(
				console.GenerateDevDarwinBinding(runtime, console.DarwinLibDir()),
				eg.Sequential(
					egbug.Log("failed to build native library / generate ffi bindings"),
					Debug(runtime),
				),
			),
			console.BuildDevDarwin,
			// flutter run re-runs hook/build.dart, and RunDev builds its runtime from a
			// bare shell.Runtime(), so without this it inherits flutterRuntimev2's
			// dev.native.libs/example.so defaults instead of the darwin ones.
			console.RunDev(
				"flutter run -d macos --use-application-binary=build/macos/Build/Products/Debug/retrovibed.app",
				func(c shell.Command) shell.Command {
					return c.
						Environ("RETROVIBED_SHARED_NATIVE_LIBS_DIRECTORY", console.DarwinLibDir()).
						Environ("NIX_RETROVIBED_SHARED_NATIVE_LIBS", filepath.Join(console.DarwinLibDir(), "example.dylib"))
				},
			),
		),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
