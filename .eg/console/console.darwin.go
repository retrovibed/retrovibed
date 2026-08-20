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
	neuralsdir := egenv.CacheDirectory("neurals")
	neuralsflags := "-L" + neuralsdir + " -lpredicttext"

	duckdblibs := egenv.CacheDirectory("duckdb", ".darwin-arm64")
	duckdbldflags := "-L" + duckdblibs + " " +
		"-Wl,-force_load," + duckdblibs + "/libduckdb.a " +
		"-lc++"

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
		Environ("BUILD_NUMBER", commit.StringReplace("%git.commit.unix%"))
	return shell.Run(
		ctx,
		runtime.New("rm -rf build/macos/{x64,arm64}/debug").Lenient(true),
		runtime.New("flutter build macos --build-name='${BUILD_NAME}' --build-number='${BUILD_NUMBER}' --release lib/main.dart"),
	)
}
