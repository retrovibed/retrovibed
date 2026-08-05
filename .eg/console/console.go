package console

import (
	"context"
	"eg/compute/tarballs"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/egfs"
	"github.com/egdaemon/eg/runtime/x/wasi/eggolang"
	"github.com/egdaemon/eg/runtime/x/wasi/egtarball"
)

func flutterRuntimev2(v shell.Command) shell.Command {
	return v.
		Directory(egenv.WorkingDirectory("console")).
		EnvironFrom(eggolang.Env()...).
		Environ("PUB_CACHE", egenv.CacheDirectory(".eg", "dart"))
}

func Tests(ctx context.Context, _ eg.Op) error {
	runtime := flutterRuntimev2(shell.Runtime())
	return shell.Run(
		ctx,
		runtime.New("flutter test"),
	)
}

func Linting(ctx context.Context, _ eg.Op) error {
	runtime := flutterRuntimev2(shell.Runtime())
	return shell.Run(
		ctx,
		runtime.New("flutter analyze"),
	)
}

func GenerateBinding(ctx context.Context, _ eg.Op) error {
	runtime := flutterRuntimev2(shell.Runtime())
	return shell.Run(
		ctx,
		runtime.New("go -C retrovibedbind build -buildmode=c-shared -buildvcs=true --tags duckdb_use_lib -o ../build/nativelib/libretrovibed.so ./..."),
		runtime.New("dart run ffigen --config ffigen.yaml"),
	)
}

func CompileBinding(b *tarballs.Build) eg.OpFn {
	return func(ctx context.Context, o eg.Op) error {
		runtime := flutterRuntimev2(shell.Runtime())
		return shell.Run(
			ctx,
			runtime.New("go -C retrovibedbind build -buildmode=c-shared -buildvcs=true --tags duckdb_use_lib,retrovibed,neural -o ../build/nativelib/libretrovibed.so .").
				Environ("CGO_LDFLAGS", "-L"+egenv.CacheDirectory("neurals")),
		)
	}
}

func GenerateDevBinding(runtime shell.Command, outdir string, staticdirs ...string) eg.OpFn {
	return func(ctx context.Context, _ eg.Op) error {
		var cgoFlags strings.Builder

		for _, dir := range staticdirs {
			fmt.Fprintf(&cgoFlags, " -L%s", dir)
		}

		runtime := flutterRuntimev2(runtime)
		return shell.Run(
			ctx,
			runtime.New("go -C retrovibedbind build -buildmode=c-shared -buildvcs=true --tags duckdb_use_lib,localdev,retrovibed,neural -o ${OUTPUT}/libretrovibed.so ./...").
				Environ("OUTPUT", outdir).
				Environ("CGO_LDFLAGS", strings.TrimSpace(cgoFlags.String())),
			runtime.New("dart run ffigen --config ffigen.yaml"),
		)
	}
}

func GenerateDevStaticBinding(rt shell.Command, outdir string, staticdirs ...string) eg.OpFn {
	return func(ctx context.Context, _ eg.Op) error {
		var cgoFlags strings.Builder

		for _, dir := range staticdirs {
			fmt.Fprintf(&cgoFlags, " -L%s", dir)
		}

		// 1. Link to the clean dynamic libduckdb.so to drop duplicate symbol issues
		// 2. Explicitly bind libpredicttext.a statically so its Rust code is baked inside
		// duckdb is an absolute disaster to statically link
		fmt.Fprintf(&cgoFlags, " -static-libstdc++ -Wl,-z,max-page-size=16384")

		runtime := flutterRuntimev2(rt)
		return shell.Run(
			ctx,
			runtime.Newf(
				"go -C retrovibedbind build -trimpath -buildmode=c-shared -buildvcs=true --tags duckdb_use_static_lib,localdev,retrovibed,neural -o %s/libretrovibed.so .",
				outdir,
			).
				Environ("CGO_LDFLAGS", strings.TrimSpace(cgoFlags.String())),
		)
	}
}

func GenerateStaticBinding(dir string, rt shell.Command) eg.OpFn {
	return func(ctx context.Context, _ eg.Op) error {
		var cgoFlags strings.Builder
		fmt.Fprintf(&cgoFlags, "-L%s", dir)
		fmt.Fprintf(&cgoFlags, " -static-libstdc++ -Wl,-z,max-page-size=16384")

		runtime := flutterRuntimev2(rt)
		return shell.Run(
			ctx,
			runtime.Newf(
				"go -C retrovibedbind build -trimpath -buildmode=c-shared -buildvcs=true --tags duckdb_use_static_lib,retrovibed,neural -o %s/libretrovibed.so .",
				dir,
			).Environ("CGO_LDFLAGS", strings.TrimSpace(cgoFlags.String())),
		)
	}
}

func AndroidRuntime(target string) shell.Command {
	return shell.Runtime().
		Environ("GOOS", "android").
		Environ("CGO_ENABLED", "1").
		Environ("GRADLE_USER_HOME", egenv.CacheDirectory("gradle")).
		Environ("CC", fmt.Sprintf("/opt/android-sdk/ndk/27.0.12077973/toolchains/llvm/prebuilt/linux-x86_64/bin/clang --sysroot=/opt/android-sdk/ndk/27.0.12077973/toolchains/llvm/prebuilt/linux-x86_64/sysroot -target %s", target)).
		Environ("CXX", fmt.Sprintf("/opt/android-sdk/ndk/27.0.12077973/toolchains/llvm/prebuilt/linux-x86_64/bin/clang++ --sysroot=/opt/android-sdk/ndk/27.0.12077973/toolchains/llvm/prebuilt/linux-x86_64/sysroot -target %s", target))
}

func Generate(ctx context.Context, op eg.Op) error {
	return eg.Sequential(
		GenerateFlutter,
		GenerateProtocol,
	)(ctx, op)
}

func MaskDartTool(ctx context.Context, _ eg.Op) error {
	// dart relies on .dart_tool cached data to work properly.
	// but it doesnt allow providing the directory location.
	// so to behavior properly between host/guest runtimes
	// we need to mask it.
	runtime := flutterRuntimev2(shell.Runtime())
	return shell.Run(
		ctx,
		runtime.New("rm -rf .dart_tool && mkdir .dart_tool && mount -t tmpfs tmpfs .dart_tool").Privileged(),
		runtime.New("tree -L 1 .dart_tool"),
	)
}

func GenerateFlutter(ctx context.Context, _ eg.Op) error {
	runtime := flutterRuntimev2(shell.Runtime())
	return shell.Run(
		ctx,
		// runtime.New("flutter clean"), // This should not be committed if uncommented. should only be needed rarely
		runtime.New("flutter create --org space.retrovibe --platforms=linux,macos,ios,android ."),
		runtime.New("flutter pub get"),
		runtime.New("dart run flutter_launcher_icons"),
	)
}

func GenerateProtocol(ctx context.Context, _ eg.Op) error {
	runtime := shell.Runtime()
	return shell.Run(
		ctx,
		runtime.New("PATH=\"${PATH}:${HOME}/.pub-cache/bin\" protoc --dart_out=grpc:console/lib/meta -I.proto .proto/meta.daemon.proto"),
		runtime.New("PATH=\"${PATH}:${HOME}/.pub-cache/bin\" protoc --dart_out=grpc:console/lib/meta -I.proto .proto/meta.profile.proto"),
		runtime.New("PATH=\"${PATH}:${HOME}/.pub-cache/bin\" protoc --dart_out=grpc:console/lib/meta -I.proto .proto/meta.authz.proto"),
		runtime.New("PATH=\"${PATH}:${HOME}/.pub-cache/bin\" protoc --dart_out=grpc:console/lib/meta -I.proto .proto/meta.account.proto"),
		runtime.New("PATH=\"${PATH}:${HOME}/.pub-cache/bin\" protoc --dart_out=grpc:console/lib/meta -I.proto .proto/meta.authn.proto"),
		runtime.New("PATH=\"${PATH}:${HOME}/.pub-cache/bin\" protoc --dart_out=grpc:console/lib/meta -I.proto .proto/meta.network.proto"),
		runtime.New("PATH=\"${PATH}:${HOME}/.pub-cache/bin\" protoc --dart_out=grpc:console/lib/meta -I.proto .proto/meta.torrent.proto"),
		runtime.New("PATH=\"${PATH}:${HOME}/.pub-cache/bin\" protoc --dart_out=grpc:console/lib/meta -I.proto .proto/meta.dht.proto"),
		runtime.New("PATH=\"${PATH}:${HOME}/.pub-cache/bin\" protoc --dart_out=grpc:console/lib/meta -I.proto .proto/meta.discovery.proto"),
		runtime.New("PATH=\"${PATH}:${HOME}/.pub-cache/bin\" protoc --dart_out=grpc:console/lib/billing -I.proto .proto/meta.billing.proto"),
		runtime.New("PATH=\"${PATH}:${HOME}/.pub-cache/bin\" protoc --dart_out=grpc:console/lib/quotas -I.proto .proto/meta.quota.proto"),
		runtime.New("PATH=\"${PATH}:${HOME}/.pub-cache/bin\" protoc --dart_out=grpc:console/lib/wireguard -I.proto .proto/meta.wireguard.proto"),
		runtime.New("PATH=\"${PATH}:${HOME}/.pub-cache/bin\" protoc --dart_out=grpc:console/lib/media -I.proto .proto/media.proto"),
		runtime.New("PATH=\"${PATH}:${HOME}/.pub-cache/bin\" protoc --dart_out=grpc:console/lib/media -I.proto .proto/meta.search.proto .proto/media.known.proto"),
		runtime.New("PATH=\"${PATH}:${HOME}/.pub-cache/bin\" protoc --dart_out=grpc:console/lib/media -I.proto .proto/media.recent.proto"),
		runtime.New("PATH=\"${PATH}:${HOME}/.pub-cache/bin\" protoc --dart_out=grpc:console/lib/media -I.proto .proto/ddisc.locate.proto"),
		runtime.New("PATH=\"${PATH}:${HOME}/.pub-cache/bin\" protoc --dart_out=grpc:console/lib/media -I.proto .proto/content.addressable.storage.proto"),
		runtime.New("PATH=\"${PATH}:${HOME}/.pub-cache/bin\" protoc --dart_out=grpc:console/lib/ddisc -I.proto .proto/meta.search.proto .proto/ddisc.discovery.proto"),
		runtime.New("PATH=\"${PATH}:${HOME}/.pub-cache/bin\" protoc --dart_out=grpc:console/lib/rss -I.proto .proto/rss.proto"),
		runtime.New("PATH=\"${PATH}:${HOME}/.pub-cache/bin\" protoc --dart_out=grpc:console/lib/community -I.proto .proto/community.proto"),
		runtime.New("PATH=\"${PATH}:${HOME}/.pub-cache/bin\" protoc --dart_out=grpc:console/lib/community -I.proto .proto/community.proto .proto/community.metrics.proto"),
		runtime.New("PATH=\"${PATH}:${HOME}/.pub-cache/bin\" protoc --dart_out=grpc:console/lib/community -I.proto .proto/community.proto .proto/community.publish.proto"),
		// frontend exclusive
		runtime.New("PATH=\"${PATH}:${HOME}/.pub-cache/bin\" protoc --dart_out=grpc:console/lib/storage -I.proto .proto/storage.proto"),
		runtime.New("PATH=\"${PATH}:${HOME}/.pub-cache/bin\" protoc --dart_out=grpc:console/lib/torrents -I.proto .proto/torrent.proto"),
		runtime.New("PATH=\"${PATH}:${HOME}/.pub-cache/bin\" protoc --dart_out=grpc:console/lib/ddisc/plugin -I.proto .proto/searchplugin.proto"),
	)
}

func Install(b *tarballs.Build) eg.OpFn {
	runtime := shell.Runtime()
	dstdir := filepath.Join(egtarball.Path(tarballs.Retrovibed(b)), "usr", "lib", "retrovibed")
	builddir := egenv.WorkingDirectory("console", "build")
	linuxdir := filepath.Join(builddir, "linux")
	libdir := filepath.Join(builddir, "nativelib")
	return eg.Sequential(
		CompileBinding(b),
		func(ctx context.Context, op eg.Op) error {
			// resolved at execution time since the bundle directory only
			// exists once console.BuildLinux has run.
			bundledir := filepath.Join(linuxdir, egfs.FindFirst(os.DirFS(linuxdir), "bundle"))
			return shell.Run(
				ctx,
				runtime.Newf("mkdir -p %s", dstdir),
				runtime.Newf("ls -lha  %s/retrovibed", bundledir).Lenient(true),
				runtime.Newf("cp -R %s/* %s", bundledir, dstdir),
				runtime.Newf("cp -R %s/* %s/lib", libdir, dstdir),
				// runtime.Newf("tree %s", dstdir),
			)
		},
	)
}

func runDev(cmd string, rt shell.Command, envopts ...func(shell.Command) shell.Command) eg.OpFn {
	return func(ctx context.Context, _ eg.Op) error {
		runtime := flutterRuntimev2(rt).
			// Environ("RETROVIBED_TORRENT_DEBUG", "true").
			Environ("RETROVIBED_TORRENT_AUTO_DISCOVERY", "true").
			Environ("RETROVIBED_METADATA_AUTODOWNLOAD", "false").
			Environ("RETROVIBED_DDISC_INDEX_RATIO", "99").
			Environ("RETROVIBED_DDISC_BACKGROUND_FREQUENCY", "20s").
			Environ("RETROVIBED_DDISC_BACKGROUND_WORKERS", "1")
		for _, envopt := range envopts {
			runtime = envopt(runtime)
		}
		return shell.Run(
			ctx,
			runtime.New(cmd).Timeout(egenv.TTL()),
		)
	}
}

func RunDev(cmd string, envopts ...func(shell.Command) shell.Command) eg.OpFn {
	rt := shell.Runtime().
		Environ("RETROVIBED_META_ENDPOINT", "localhost:8081").
		Environ("RETROVIBED_CONSOLE_ENDPOINT", "localhost:8080")
		// Environ("RETROVIBED_META_ENDPOINT", "api.retrovibe.space").
		// Environ("RETROVIBED_CONSOLE_ENDPOINT", "console.retrovibe.space")
	return runDev(cmd, rt, envopts...)
}
