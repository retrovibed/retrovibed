package main

import (
	"context"
	"eg/compute/console"
	"log"
	"time"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/egbug"
	"github.com/egdaemon/eg/runtime/x/wasi/eggolang"
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

	runtime := shell.Runtime().
		EnvironFrom(eggolang.Env()...).
		Environ("PUB_CACHE", egenv.CacheDirectory(".eg", "dart"))

	flutter := runtime.Directory(egenv.WorkingDirectory("console"))

	err := eg.Perform(
		ctx,
		eg.Sequential(
			shell.Op(
				shell.New("brew install go duckdb gpgme flutter ffmpeg@7 cocoapods"),
			),
			egbug.DebugFailure(
				shell.Op(
					flutter.New("mkdir -p build/nativelib"),
					flutter.New("go -C retrovibedbind build -buildmode=c-shared --tags localdev -o ../build/nativelib/retrovibed.dylib ./...").Timeout(5*time.Minute),
				),
				eg.Sequential(
					egbug.Log("failed to build native library"),
					Debug(runtime),
				),
			),
			egbug.DebugFailure(
				shell.Op(
					flutter.New("cp build/nativelib/retrovibed.h build/nativelib/libretrovibed.h"),
					flutter.New("dart run ffigen --config ffigen.yaml --compiler-opts \"-I$(clang --print-resource-dir)/include\""),
				),
				eg.Sequential(
					egbug.Log("failed to generate ffi bindings"),
					Debug(runtime),
				),
			),
			shell.Op(
				flutter.New("flutter pub get"),
				flutter.New("flutter build macos --debug"),
				flutter.New("cp build/nativelib/retrovibed.dylib build/macos/Build/Products/Debug/retrovibed.app/Contents/MacOS/retrovibed.dylib"),
				flutter.New("codesign --force --sign - build/macos/Build/Products/Debug/retrovibed.app"),
			),
			console.RunDev("flutter run -d macos --use-application-binary=build/macos/Build/Products/Debug/retrovibed.app"),
		),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
