package main

import (
	"context"
	"eg/compute/console"
	"log"
	"path/filepath"
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
		runtime.New("which cmake && cmake --version"),
		runtime.New("xcrun --sdk iphonesimulator --show-sdk-path"),
		runtime.New("xcode-select -p"),
	)
}

// EnsureDuckDB clones sources and builds the static library for iOS simulator.
func EnsureDuckDB(runtime shell.Command, duckdbsrc, duckdbinet, duckdbbuild string) eg.OpFn {
	return func(ctx context.Context, op eg.Op) error {
		if err := console.EnsureDuckDBSource(duckdbsrc, duckdbinet)(ctx, op); err != nil {
			return err
		}

		libpath := filepath.Join(duckdbbuild, "src", "libduckdb_static.a")
		return egbug.DebugFailure(
			shell.Op(
				shell.Newf(
					"test -f %s || cmake -B %s -S %s -DCMAKE_SYSTEM_NAME=iOS -DCMAKE_OSX_ARCHITECTURES=x86_64 -DCMAKE_OSX_SYSROOT=$(xcrun --sdk iphonesimulator --show-sdk-path) -DCMAKE_OSX_DEPLOYMENT_TARGET=16.0 -DCMAKE_BUILD_TYPE=Release -DEXTENSION_STATIC_BUILD=1 -DBUILD_SHELL=0 -DBUILD_UNITTESTS=0 -DDUCKDB_EXPLICIT_PLATFORM=ios_x86_64 \"-DBUILD_EXTENSIONS=icu;json\" -GNinja",
					libpath, duckdbbuild, duckdbsrc,
				),
				shell.Newf("test -f %s || cmake --build %s", libpath, duckdbbuild).Timeout(10*time.Minute),
			),
			eg.Sequential(
				egbug.Log("failed to build duckdb for iOS simulator"),
				Debug(runtime),
			),
		)(ctx, op)
	}
}

// Baremetal command for iOS simulator local development.
func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	ctx, done := context.WithTimeout(context.Background(), egenv.TTL())
	defer done()

	runtime := shell.Runtime().
		EnvironFrom(eggolang.Env()...).
		Environ("PUB_CACHE", egenv.CacheDirectory(".eg", "dart"))

	flutter := runtime.Directory(egenv.WorkingDirectory("console"))
	duckdbsrc := egenv.CacheDirectory("duckdb", "v1.4.3", "src")
	duckdbinet := egenv.CacheDirectory("duckdb", "v1.4.3", "inet")
	duckdbbuild := egenv.CacheDirectory("duckdb", "v1.4.3", "ios-sim-x86_64")

	err := eg.Perform(
		ctx,
		eg.Sequential(
			shell.Op(
				shell.New("brew install go duckdb gpgme flutter ffmpeg@7 cocoapods cmake ninja"),
			),
			EnsureDuckDB(runtime, duckdbsrc, duckdbinet, duckdbbuild),
			egbug.DebugFailure(
				shell.Op(
					flutter.New("mkdir -p build/nativelib/ios-sim-x86_64"),
					flutter.Newf(
						"GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 CC=\"$(xcrun --sdk iphonesimulator --find clang) -target x86_64-apple-ios16.0-simulator -isysroot $(xcrun --sdk iphonesimulator --show-sdk-path)\" CGO_CFLAGS=\"-DDUCKDB_STATIC_BUILD -I%s/src/include\" CGO_LDFLAGS=\"-L%s -lduckdb_static -lc++\" go -C retrovibedbind build -trimpath -buildmode=c-archive --tags duckdb_use_static_lib,localdev -o ../build/nativelib/ios-sim-x86_64/libretrovibed.a ./...",
						duckdbsrc, duckdbbuild,
					).Timeout(5*time.Minute),
				),
				eg.Sequential(
					egbug.Log("failed to build native library for iOS simulator"),
					Debug(runtime),
				),
			),
			egbug.DebugFailure(
				shell.Op(
					flutter.New("cp build/nativelib/ios-sim-x86_64/libretrovibed.h build/nativelib/libretrovibed.h"),
					flutter.New("dart run ffigen --config ffigen.yaml --compiler-opts \"-I$(clang --print-resource-dir)/include\""),
				),
				eg.Sequential(
					egbug.Log("failed to generate ffi bindings"),
					Debug(runtime),
				),
			),
			console.RetagSimulator(
				egenv.WorkingDirectory("console", "build", "nativelib", "ios-sim-x86_64", "libretrovibed.a"),
				egenv.WorkingDirectory("console", "ios", "libretrovibed.a"),
			),
			shell.Op(
				flutter.New("cp build/nativelib/ios-sim-x86_64/libretrovibed.h ios/Classes/libretrovibed.h"),
				flutter.Newf("libtool -static -o ios/libduckdb_static.a $(find %s -name '*.a' ! -path '*/test/*')", duckdbbuild),
				flutter.New("flutter create --org space.retrovibe --platforms=ios ."),
				flutter.New("flutter pub get"),
				flutter.New("cd ios && pod install"),
				runtime.New("open -a Simulator"),
				runtime.New("xcrun simctl list devices 'iOS 26.4' | grep -q 'Retrovibed Review' || xcrun simctl create 'Retrovibed Review' com.apple.CoreSimulator.SimDeviceType.iPad-Air-5th-generation com.apple.CoreSimulator.SimRuntime.iOS-26-4"),
				runtime.New("xcrun simctl boot 'Retrovibed Review'").Lenient(true),
				runtime.New("xcrun simctl list devices booted").Attempts(15),
			),
			console.RunDev("flutter run -d 'Retrovibed Review'"),
		),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
