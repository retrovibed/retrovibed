package shallows

import (
	"context"
	"eg/compute/tarballs"
	"strings"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/egflatpak"
	"github.com/egdaemon/eg/runtime/x/wasi/eggolang"
	"github.com/egdaemon/eg/runtime/x/wasi/egtarball"
)

var buildTags = []string{"no_duckdb_arrow"}

func shellruntime() shell.Command {
	return eggolang.Runtime().Directory(egenv.WorkingDirectory("shallows"))
}

func Generate(ctx context.Context, _ eg.Op) error {
	gruntime := shellruntime()
	return shell.Run(
		ctx,
		gruntime.New("go generate ./... && go fmt ./..."),
	)
}

func Install(ctx context.Context, _ eg.Op) error {
	// go install -ldflags=\"-extldflags=-static\" -tags no_duckdb_arrow ./cmd/shallows/...
	dstdir := egtarball.Path(tarballs.Retrovibed())
	gruntime := shellruntime()
	return shell.Run(
		ctx,
		// gruntime.New("ldconfig -p | grep duckdb"),
		// gruntime.New("ld --verbose | grep SEARCH_DIR | tr -s ' ;'"),
		// gruntime.New("go env"),
		gruntime.Newf("go install -tags %s ./cmd/...", strings.Join(buildTags, ",")).Environ("GOBIN", dstdir),
	)
}

func Compile() eg.OpFn {
	return eggolang.AutoCompile(
		eggolang.CompileOption.BuildOptions(
			eggolang.Build(
				eggolang.BuildOption.Tags(buildTags...),
			),
		),
	)
}

func Test() eg.OpFn {
	return eg.Sequential(eggolang.AutoTest(
		eggolang.TestOption.BuildOptions(
			eggolang.Build(
				eggolang.BuildOption.Tags(buildTags...),
			),
		),
	),
		eggolang.RecordCoverage,
	)
}

func FlatpakManifest(ctx context.Context, o eg.Op) error {
	b := egflatpak.New(
		"space.retrovibe.Daemon", "retrovibed",
		egflatpak.Option().SDK("org.gnome.Sdk", "47").Runtime("org.gnome.Platform", "47").
			Modules(
				egflatpak.NewModule("duckdb", "simple", egflatpak.ModuleOptions().Commands(
					"cp -r . /app/lib",
				).Sources(
					egflatpak.SourceTarball(
						"https://github.com/duckdb/duckdb/releases/download/v1.1.3/libduckdb-linux-amd64.zip",
						"81199bf01b6d49941a38f426cad60e73c1ccd43f1f769a65ed8097d53fc7e40b",
						egflatpak.SourceOptions().Destination("duckdb.zip")...,
					),
				)...),
				egflatpak.NewModule("retrovibed", "simple", egflatpak.ModuleOptions().Commands(
					"cp -r . /app/bin",
				).Sources(
					egflatpak.SourceTarball(
						egtarball.GithubDownloadURL(tarballs.Retrovibed()), egtarball.SHA256(tarballs.Retrovibed()),
						egflatpak.SourceOptions().Destination("retrovibed.tar.xz")...,
					),
				)...),
			).
			AllowWayland().
			AllowDRI().
			AllowNetwork().
			AllowDownload().
			AllowMusic().
			AllowVideos().Allow(
			"--filesystem=/home/retrovibed/Downloads:ro", // bug in flatpak doesn't properly grant access to xdg-download
		)...)

	return egflatpak.ManifestOp(egenv.CacheDirectory("flatpak.daemon.yml"), b)(ctx, o)
}

// func Flatpak(ctx context.Context, op eg.Op) error {
// 	runtime := shell.Runtime()
// 	builddir := egenv.WorkingDirectory("console", "build", egfs.FindFirst(os.DirFS(egenv.WorkingDirectory("console", "build")), "bundle"))
// 	b := egflatpak.New(
// 		"space.retrovibe.Daemon", "console",
// 		egflatpak.Option().SDK("org.gnome.Sdk", "47").Runtime("org.gnome.Platform", "47").
// 			Modules(
// 				egflatpak.ModuleCopy(builddir),
// 			).
// 			AllowWayland().
// 			AllowDRI().
// 			AllowNetwork().
// 			AllowDownload().
// 			AllowMusic().
// 			AllowVideos()...)

// 	if err := egflatpak.Build(ctx, runtime, b); err != nil {
// 		return err
// 	}

// 	return nil
// }
