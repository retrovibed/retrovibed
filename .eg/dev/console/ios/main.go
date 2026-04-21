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
func EnsureDuckDB(runtime shell.Command, duckdbsrc, duckdbinet, duckdbbuild, arch string) eg.OpFn {
	return func(ctx context.Context, op eg.Op) error {
		if err := console.EnsureDuckDBSource(duckdbsrc, duckdbinet)(ctx, op); err != nil {
			return err
		}

		libpath := filepath.Join(duckdbbuild, "src", "libduckdb_static.a")
		return egbug.DebugFailure(
			shell.Op(
				shell.Newf(
					"test -f %s || cmake -B %s -S %s -DCMAKE_SYSTEM_NAME=iOS -DCMAKE_OSX_ARCHITECTURES=%s -DCMAKE_OSX_SYSROOT=$(xcrun --sdk iphonesimulator --show-sdk-path) -DCMAKE_OSX_DEPLOYMENT_TARGET=16.0 -DCMAKE_BUILD_TYPE=Release -DEXTENSION_STATIC_BUILD=1 -DBUILD_SHELL=0 -DBUILD_UNITTESTS=0 -DDUCKDB_EXPLICIT_PLATFORM=ios_%s \"-DBUILD_EXTENSIONS=icu;json\" -GNinja",
					libpath, duckdbbuild, duckdbsrc, arch, arch,
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

func GoBuild(flutter shell.Command, duckdbsrc, duckdbbuild, goarch, target, outdir string) shell.Command {
	return flutter.Newf(
		"CC=\"$(xcrun --sdk iphonesimulator --find clang)\" "+
			"CXX=\"$(xcrun --sdk iphonesimulator --find clang++)\" "+
			"CGO_CFLAGS=\"-target %[3]s -isysroot $(xcrun --sdk iphonesimulator --show-sdk-path) -DDUCKDB_STATIC_BUILD -I%[1]s/src/include\" "+
			"CGO_CXXFLAGS=\"-target %[3]s -isysroot $(xcrun --sdk iphonesimulator --show-sdk-path)\" "+
			"CGO_LDFLAGS=\"-target %[3]s -isysroot $(xcrun --sdk iphonesimulator --show-sdk-path) -L%[2]s -lduckdb_static -lc++\" "+
			"GOOS=ios GOARCH=%[4]s CGO_ENABLED=1 "+
			"go -C retrovibedbind build -trimpath -buildmode=c-archive --tags duckdb_use_static_lib,localdev -o ../build/nativelib/%[5]s/libretrovibed.a ./...",
		duckdbsrc, duckdbbuild, target, goarch, outdir,
	).Timeout(5 * time.Minute)
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
	duckdbArm64 := egenv.CacheDirectory("duckdb", "v1.4.3", "ios-sim-arm64")
	duckdbX86 := egenv.CacheDirectory("duckdb", "v1.4.3", "ios-sim-x86_64")

	err := eg.Perform(
		ctx,
		eg.Sequential(
			shell.Op(
				shell.New("brew install go duckdb gpgme flutter ffmpeg@7 cocoapods cmake ninja"),
			),
			EnsureDuckDB(runtime, duckdbsrc, duckdbinet, duckdbArm64, "arm64"),
			EnsureDuckDB(runtime, duckdbsrc, duckdbinet, duckdbX86, "x86_64"),
			egbug.DebugFailure(
				shell.Op(
					flutter.New("mkdir -p build/nativelib/ios-sim-arm64 build/nativelib/ios-sim-x86_64"),
					GoBuild(flutter, duckdbsrc, duckdbArm64, "arm64", "arm64-apple-ios16.0-simulator", "ios-sim-arm64"),
					GoBuild(flutter, duckdbsrc, duckdbX86, "amd64", "x86_64-apple-ios16.0-simulator", "ios-sim-x86_64"),
				),
				eg.Sequential(
					egbug.Log("failed to build native library for iOS simulator"),
					Debug(runtime),
				),
			),
			egbug.DebugFailure(
				shell.Op(
					flutter.New("cp build/nativelib/ios-sim-arm64/libretrovibed.h build/nativelib/libretrovibed.h"),
					flutter.New("dart run ffigen --config ffigen.yaml --compiler-opts \"-I$(clang --print-resource-dir)/include\""),
				),
				eg.Sequential(
					egbug.Log("failed to generate ffi bindings"),
					Debug(runtime),
				),
			),
			shell.Op(
				flutter.New("lipo -create build/nativelib/ios-sim-arm64/libretrovibed.a build/nativelib/ios-sim-x86_64/libretrovibed.a -output ios/libretrovibed.a"),
				flutter.Newf("libtool -static -o ios/libduckdb_static.a $(find %s %s -name '*.a' ! -path '*/test/*')", duckdbArm64, duckdbX86),
				flutter.New("rm -rf ios/RetrovivedBind.framework ios/RetrovivedBind.xcframework && cp -r ios/RetrovivedBind ios/RetrovivedBind.framework"),
				flutter.New("bash -c 'cd ios && xcrun clang -arch arm64 -arch x86_64 -isysroot \"$(xcrun --sdk iphonesimulator --show-sdk-path)\" -mios-simulator-version-min=16.0 -shared -o RetrovivedBind.framework/RetrovivedBind $(for f in *.a; do printf -- \"-Wl,-force_load,%s \" \"$f\"; done) -lc++ -lresolv -framework CoreFoundation -framework Security -Wl,-install_name,@rpath/RetrovivedBind.framework/RetrovivedBind'"),
				flutter.New("bash -c 'cd ios && xcodebuild -create-xcframework -framework RetrovivedBind.framework -output RetrovivedBind.xcframework'"),
				flutter.New("flutter create --org space.retrovibe --platforms=ios ."),
				flutter.New("flutter pub get"),
				flutter.New("cd ios && pod install"),
				runtime.New("open -a Simulator"),
				runtime.New("xcrun simctl boot 'iPhone 16'").Lenient(true),
				runtime.New("xcrun simctl list devices booted").Attempts(15),
			),
			console.RunDev("flutter run -d iPhone"),
		),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
