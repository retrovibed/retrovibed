package fractal

import (
	"context"
	"eg/compute/tarball"
	"eg/compute/tarballs"
	"os"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/egflatpak"
	"github.com/egdaemon/eg/runtime/x/wasi/egfs"
)

func flutterRuntime() shell.Command {
	return shell.Runtime().Directory(egenv.WorkingDirectory("fractal")).Environ("PUB_CACHE", egenv.CacheDirectory(".eg", "dart"))
}

func Build(ctx context.Context, _ eg.Op) error {
	runtime := flutterRuntime()
	return shell.Run(
		ctx,
		runtime.New("flutter create --platforms=linux ."),
		runtime.Newf("flutter build bundle"),
		runtime.Newf("flutter build linux"),
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
		shell.New("PATH=\"${PATH}:${HOME}/.pub-cache/bin\" protoc --dart_out=grpc:fractal/lib/media -I.proto .proto/media.proto"),
		shell.New("PATH=\"${PATH}:${HOME}/.pub-cache/bin\" protoc --dart_out=grpc:fractal/lib/rss -I.proto .proto/rss.proto"),
	)
}

func Install(ctx context.Context, op eg.Op) error {
	runtime := shell.Runtime()
	dstdir := tarball.Path(tarballs.Retrovibed())
	builddir := egenv.WorkingDirectory("fractal", "build", egfs.FindFirst(os.DirFS(egenv.WorkingDirectory("fractal", "build")), "bundle"))

	return shell.Run(
		ctx,
		runtime.Newf("mkdir -p %s", dstdir),
		runtime.Newf("cp -R %s/* %s", builddir, dstdir),
	)
}

func FlatpakDir(ctx context.Context, op eg.Op) error {
	runtime := shell.Runtime()
	builddir := egenv.WorkingDirectory("fractal", "build", egfs.FindFirst(os.DirFS(egenv.WorkingDirectory("fractal", "build")), "bundle"))

	b := egflatpak.New(
		"space.retrovibe.Daemon", "fractal",
		egflatpak.Option().SDK("org.gnome.Sdk", "47").Runtime("org.gnome.Platform", "47").
			Modules(egflatpak.ModuleCopy(builddir)).
			AllowWayland().
			AllowDRI().
			AllowNetwork().
			AllowDownload().
			AllowMusic().
			AllowVideos()...)

	if err := egflatpak.Build(ctx, runtime, b); err != nil {
		return err
	}

	return nil
}

func FlatpakManifest(ctx context.Context, o eg.Op) error {
	b := egflatpak.New(
		"space.retrovibe.Client", "fractal",
		egflatpak.Option().SDK("org.gnome.Sdk", "47").Runtime("org.gnome.Platform", "47").
			Modules(
				egflatpak.ModuleTarball(tarball.GithubDownloadURL(tarballs.Retrovibed()), tarball.SHA256(tarballs.Retrovibed())),
			).
			AllowWayland().
			AllowDRI().
			AllowNetwork().
			AllowDownload().
			AllowMusic().
			AllowVideos()...)

	return egflatpak.ManifestOp(egenv.CacheDirectory("flatpak.client.yml"), b)(ctx, o)
}
