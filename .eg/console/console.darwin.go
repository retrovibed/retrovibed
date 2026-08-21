package console

import (
	"context"
	"eg/compute/tarballs"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
	"github.com/egdaemon/eg/runtime/wasi/shell"
)

func CompileDarwinBinding(ctx context.Context, o eg.Op) error {
	runtime := flutterRuntimev2(shell.Runtime())
	libsdir := egenv.CacheDirectory("dev.native.libs")

	neuralsflags := "-L" + libsdir + " -lpredicttext"
	duckdbldflags := "-Wl,-force_load," + libsdir + "/libduckdb.a " + "-lc++"

	return shell.Run(
		ctx,
		runtime.New("mkdir -p ${RETROVIBED_SHARED_NATIVE_LIBS_DIRECTORY}"),
		runtime.New("go -C retrovibedbind build -buildmode=c-shared -buildvcs=true --tags duckdb_use_static_lib,retrovibed,neural -o ${RETROVIBED_SHARED_NATIVE_LIBS_DIRECTORY}/libretrovibed.dylib ./...").
			Timeout(egenv.TTL()).
			Environ("CGO_LDFLAGS", neuralsflags+" "+duckdbldflags),
	)
}

func BuildDarwin(ctx context.Context, _ eg.Op) error {
	commit := eggit.EnvCommit()
	runtime := flutterRuntimev2(shell.Runtime()).
		Environ("BUILD_NAME", tarballs.Version()).
		Environ("BUILD_NUMBER", commit.StringReplace("%git.commit.unix%")).
		Environ("NIX_RETROVIBED_SHARED_NATIVE_LIBS", egenv.CacheDirectory("dev.native.libs", "example.dylib")).
		// nothing in this pipeline builds x86_64 native libs; exclude it so flutter
		// doesn't attempt a universal binary and hand the code_assets hook's single
		// arm64 dylib to lipo twice. FLUTTER_XCODE_* env vars are forwarded to
		// xcodebuild as build settings (see flutter_tools/lib/src/ios/xcodeproj.dart).
		Environ("FLUTTER_XCODE_EXCLUDED_ARCHS", "x86_64")
	return shell.Run(
		ctx,
		runtime.New("rm -rf build/macos/{x64,arm64}/debug").Lenient(true),
		runtime.New("flutter build macos -v --build-name='${BUILD_NAME}' --build-number='${BUILD_NUMBER}' --release lib/main.dart"),
	)
}
