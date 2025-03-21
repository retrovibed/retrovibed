package shallows

import (
	"context"
	"eg/compute/flatpakmods"
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

func rootdir() string {
	return egenv.WorkingDirectory("shallows")
}

func shellruntime() shell.Command {
	return eggolang.Runtime().Directory(rootdir())
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
		gruntime.Newf("go install -tags %s ./cmd/...", strings.Join(buildTags, ",")).Environ("GOBIN", dstdir),
	)
}

func Compile() eg.OpFn {
	return eggolang.AutoCompile(
		eggolang.CompileOption.BuildOptions(
			eggolang.Build(
				eggolang.BuildOption.Tags(buildTags...),
				eggolang.BuildOption.WorkingDirectory(rootdir()),
			),
		),
	)
}

func Test() eg.OpFn {
	return eg.Sequential(eggolang.AutoTest(
		eggolang.TestOption.BuildOptions(
			eggolang.Build(
				eggolang.BuildOption.Tags(buildTags...),
				eggolang.BuildOption.WorkingDirectory(rootdir()),
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
				flatpakmods.Libduckdb(),
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
			"--filesystem=~/Downloads:ro",  // bug in flatpak doesn't properly grant access to xdg-download
			"--filesystem=~/Videos:create", // bug in flatpak doesn't properly grant full access to videos directory
			"--filesystem=~/Music:create",  // bug in flatpak doesn't properly grant full access to music directory
		)...)

	return egflatpak.ManifestOp(egenv.CacheDirectory("flatpak.daemon.yml"), b)(ctx, o)
}
