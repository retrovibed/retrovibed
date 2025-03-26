package console

import (
	"context"
	"eg/compute/errorsx"
	"eg/compute/flatpakmods"
	"eg/compute/tarballs"
	"os"
	"path/filepath"
	"time"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/egflatpak"
	"github.com/egdaemon/eg/runtime/x/wasi/egfs"
	"github.com/egdaemon/eg/runtime/x/wasi/eggolang"
	"github.com/egdaemon/eg/runtime/x/wasi/egtarball"
)

func flutterRuntime() shell.Command {
	return shell.Runtime().Directory(egenv.WorkingDirectory("console")).EnvironFrom(errorsx.Must(eggolang.Env())...).Environ("PUB_CACHE", egenv.CacheDirectory(".eg", "dart"))
}

func Build(ctx context.Context, _ eg.Op) error {
	runtime := flutterRuntime()
	return shell.Run(
		ctx,
		runtime.New("flutter create --platforms=linux ."),
		runtime.Newf("flutter build linux --release"),
		runtime.New("go -C retrovibedbind build -buildmode=c-shared --tags no_duckdb_arrow -o ../build/nativelib/retrovibed.so ./..."),
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
	return shell.Run(
		ctx,
		shell.New("PATH=\"${PATH}:${HOME}/.pub-cache/bin\" protoc --dart_out=grpc:console/lib/media -I.proto .proto/media.proto"),
		shell.New("PATH=\"${PATH}:${HOME}/.pub-cache/bin\" protoc --dart_out=grpc:console/lib/rss -I.proto .proto/rss.proto"),
		shell.New("dart run ffigen --config ffigen.yaml").Directory(egenv.WorkingDirectory("console")),
	)
}

func Install(ctx context.Context, op eg.Op) error {
	runtime := shell.Runtime()
	dstdir := egtarball.Path(tarballs.Retrovibed())
	builddir := egenv.WorkingDirectory("console", "build")
	bundledir := filepath.Join(builddir, egfs.FindFirst(os.DirFS(builddir), "bundle"))
	libdir := filepath.Join(builddir, "nativelib")
	return shell.Run(
		ctx,
		runtime.Newf("mkdir -p %s", dstdir),
		runtime.Newf("tree %s", dstdir),
		runtime.Newf("cp -R %s/* %s", bundledir, dstdir),
		runtime.Newf("cp -R %s/* %s/lib", libdir, dstdir),
	)
}

func flatpak(final egflatpak.Module) *egflatpak.Builder {
	return egflatpak.New(
		"space.retrovibe.Console", "console",
		egflatpak.Option().SDK("org.gnome.Sdk", "47").Runtime("org.gnome.Platform", "47").
			Modules(
				flatpakmods.Libass(),
				flatpakmods.Libbs2b(),
				flatpakmods.Libplacebo(),
				flatpakmods.Libx264(),
				flatpakmods.Libx265(),
				flatpakmods.Libffmpeg(),
				flatpakmods.Libmpv(),
				final,
			).
			AllowWayland().
			AllowDRI().
			AllowNetwork().
			AllowDownload().
			AllowMusic().
			AllowVideos().Allow(
			"--filesystem=host:ro",               // for mpv
			"--socket=pulseaudio",                // for mpv
			"--env=LC_NUMERIC=C",                 // for mpv
			"--filesystem=xdg-run/pipewire-0:ro", // for mpv
		)...)
}

// build ensures that the flatpak has all the necessary componentry for the generated manifest.
func FlatpakBuild(ctx context.Context, op eg.Op) error {
	// builddir := egenv.WorkingDirectory("console", "build", egfs.FindFirst(os.DirFS(egenv.WorkingDirectory("console", "build")), "bundle"))
	return egflatpak.Build(ctx, shell.Runtime().Timeout(30*time.Minute), flatpak(
		egflatpak.ModuleTarball(egtarball.GithubDownloadURL(tarballs.Retrovibed()), egtarball.SHA256(tarballs.Retrovibed())),
	// egflatpak.ModuleCopy(builddir),
	))
}

// Manifest generates the manifest for distribution.
func FlatpakManifest(ctx context.Context, o eg.Op) error {
	return egflatpak.ManifestOp(egenv.CacheDirectory("flatpak.client.yml"), flatpak(
		egflatpak.ModuleTarball(egtarball.GithubDownloadURL(tarballs.Retrovibed()), egtarball.SHA256(tarballs.Retrovibed())),
	))(ctx, o)
}
