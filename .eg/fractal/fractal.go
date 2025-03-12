package fractal

import (
	"context"
	"eg/compute/tarballs"
	"os"
	"time"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/egflatpak"
	"github.com/egdaemon/eg/runtime/x/wasi/egfs"
	"github.com/egdaemon/eg/runtime/x/wasi/egtarball"
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
	dstdir := egtarball.Path(tarballs.Retrovibed())
	builddir := egenv.WorkingDirectory("fractal", "build", egfs.FindFirst(os.DirFS(egenv.WorkingDirectory("fractal", "build")), "bundle"))

	return shell.Run(
		ctx,
		runtime.Newf("mkdir -p %s", dstdir),
		runtime.Newf("cp -R %s/* %s", builddir, dstdir),
	)
}

func flatpak(final egflatpak.Module) *egflatpak.Builder {
	return egflatpak.New(
		"space.retrovibe.Client", "fractal",
		egflatpak.Option().SDK("org.gnome.Sdk", "47").Runtime("org.gnome.Platform", "47").
			Modules(
				libassmodule(),
				libbs2bmodule(),
				libplacebomodule(),
				libx264module(),
				libx265module(),
				libffmpegmodule(),
				mpvmodule(),
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
	// builddir := egenv.WorkingDirectory("fractal", "build", egfs.FindFirst(os.DirFS(egenv.WorkingDirectory("fractal", "build")), "bundle"))
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

// pulled from: https://github.com/flathub/io.mpv.Mpv/blob/d895bc41c09a17d0bdca40cd57f77340e44fdca5/io.mpv.Mpv.yml
func libx264module() egflatpak.Module {
	return egflatpak.NewModule("libx264", "autotools", egflatpak.ModuleOptions().ConfigOptions(
		"--disable-cli",
		"--enable-shared",
	).Cleanup(
		"/share/man",
	).Sources(
		egflatpak.SourceGit(
			"https://github.com/mirror/x264.git",
			"31e19f92f00c7003fa115047ce50978bc98c3a0d",
		),
	)...)
}

// pulled from: https://github.com/flathub/io.mpv.Mpv/blob/d895bc41c09a17d0bdca40cd57f77340e44fdca5/io.mpv.Mpv.yml
func libx265module() egflatpak.Module {
	// - name: x265
	return egflatpak.NewModule("libx265", "cmake", egflatpak.ModuleOptions().SubDirectory("source").ConfigOptions(
		"-DCMAKE_BUILD_TYPE=Release",
		"-DBUILD_STATIC=0",
	).Cleanup(
		"/share/man",
	).Sources(
		egflatpak.SourceGit(
			"https://bitbucket.org/multicoreware/x265_git.git",
			"6318f223684118a2c71f67f3f4633a9e35046b00",
			egflatpak.SourceOptions().Tag("4.0")...,
		),
	)...)
}

// pulled from: https://github.com/flathub/io.mpv.Mpv/blob/d895bc41c09a17d0bdca40cd57f77340e44fdca5/io.mpv.Mpv.yml
func libassmodule() egflatpak.Module {
	return egflatpak.NewModule("libass", "autotools", egflatpak.ModuleOptions().ConfigOptions(
		"--disable-static",
		"--enable-asm",
		"--enable-harfbuzz",
		"--enable-fontconfig",
	).Sources(
		egflatpak.SourceGit(
			"https://github.com/libass/libass.git",
			"e46aedea0a0d17da4c4ef49d84b94a7994664ab5",
			egflatpak.SourceOptions().Tag(
				"0.17.3",
			)...,
		),
	)...)
}

// pulled from: https://github.com/flathub/io.mpv.Mpv/blob/d895bc41c09a17d0bdca40cd57f77340e44fdca5/io.mpv.Mpv.yml
func libbs2bmodule() egflatpak.Module {
	return egflatpak.NewModule("libbs2b", "autotools", egflatpak.ModuleOptions().ConfigOptions(
		"--disable-static",
	).Sources(
		egflatpak.SourceTarball(
			"https://downloads.sourceforge.net/sourceforge/bs2b/libbs2b-3.1.0.tar.gz",
			"6aaafd81aae3898ee40148dd1349aab348db9bfae9767d0e66e0b07ddd4b2528",
			egflatpak.SourceOptions().Tag(
				"v1.3.3",
			)...,
		),
		egflatpak.SourceShell(
			egflatpak.SourceOptions().Commands(
				"sed -i -e 's/lzma/xz/g' configure.ac",
				"autoreconf -vif",
			)...,
		),
	)...)
}

// pulled from: https://github.com/flathub/io.mpv.Mpv/blob/d895bc41c09a17d0bdca40cd57f77340e44fdca5/io.mpv.Mpv.yml
func libffmpegmodule() egflatpak.Module {
	return egflatpak.NewModule("libffmpeg", "autotools", egflatpak.ModuleOptions().ConfigOptions(
		"--disable-debug",
		"--disable-doc",
		"--disable-static",
		"--enable-encoder=png",
		"--enable-gnutls",
		"--enable-gpl",
		"--enable-shared",
		"--enable-version3",
		"--enable-libaom",
		"--enable-libbs2b",
		"--enable-libdav1d",
		"--enable-libfreetype",
		"--enable-libmp3lame",
		"--enable-libopus",
		"--enable-libjxl",
		"--enable-libtheora",
		"--enable-libv4l2",
		"--enable-libvorbis",
		"--enable-libvpx",
		"--enable-vulkan",
		"--enable-libass",
		"--enable-libx264",
		"--enable-libx265",
		"--enable-libwebp",
		"--enable-libxml2",
		// "--enable-libmysofa",
	).Cleanup(
		"/share/ffmpeg/examples",
	).Sources(
		egflatpak.SourceGit(
			"https://github.com/FFmpeg/FFmpeg.git",
			"db69d06eeeab4f46da15030a80d539efb4503ca8",
			egflatpak.SourceOptions().Tag(
				"n7.1.1",
			)...,
		),
	)...)
}

// pulled from: https://github.com/flathub/io.mpv.Mpv/blob/d895bc41c09a17d0bdca40cd57f77340e44fdca5/io.mpv.Mpv.yml
func libplacebomodule() egflatpak.Module {
	return egflatpak.NewModule("libplacebo", "meson", egflatpak.ModuleOptions().ConfigOptions(
		"-Dvulkan=enabled",
		"-Dshaderc=enabled",
	).Sources(
		egflatpak.SourceGit(
			"https://github.com/haasn/libplacebo.git",
			"1fd3c7bde7b943fe8985c893310b5269a09b46c5",
			egflatpak.SourceOptions().Tag(
				"v7.349.0",
			)...,
		),
	)...)
}

// pulled from: https://github.com/flathub/io.mpv.Mpv/blob/d895bc41c09a17d0bdca40cd57f77340e44fdca5/io.mpv.Mpv.yml
func mpvmodule() egflatpak.Module {
	// https://github.com/mpv-player/mpv/releases/tag/v0.39.0
	return egflatpak.NewModule("mpv", "meson", egflatpak.ModuleOptions().ConfigOptions(
		"-Dlibmpv=true",
		// "-Dcdda=enabled",
		// "-Ddvbin=enabled",
		// "-Ddvdnav=enabled",
		// "-Dlibarchive=enabled",
		// "-Dsdl2=enabled",
		"-Dvulkan=enabled",
		"-Dmanpage-build=disabled",
		"-Dbuild-date=false",
	).PostInstall(
		"pwd",
		"ls -lha .",
	).Sources(egflatpak.SourceTarball("https://github.com/mpv-player/mpv/archive/refs/tags/v0.39.0.tar.gz", "2ca92437affb62c2b559b4419ea4785c70d023590500e8a52e95ea3ab4554683"))...)
}
