package shallows

import (
	"context"
	"eg/compute/tarball"
	"os"
	"strings"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/egflatpak"
	"github.com/egdaemon/eg/runtime/x/wasi/egfs"
	"github.com/egdaemon/eg/runtime/x/wasi/eggolang"
)

var buildTags = []string{"no_duckdb_arrow"}

func runtime() shell.Command {
	return eggolang.Runtime().Directory(egenv.WorkingDirectory("shallows"))
}

func Generate(ctx context.Context, _ eg.Op) error {
	gruntime := runtime()
	return shell.Run(
		ctx,
		gruntime.New("go generate ./... && go fmt ./..."),
	)
}

func Install(ctx context.Context, _ eg.Op) error {
	dstdir := tarball.Path(tarball.GitPattern("retrovibed"))
	gruntime := runtime()
	return shell.Run(
		ctx,
		gruntime.New("ldconfig -p | grep duckdb"),
		gruntime.New("ld --verbose | grep SEARCH_DIR | tr -s ' ;'"),
		gruntime.New("go env"),
		gruntime.Newf("go install -ldflags=\"-extldflags=-static\" -tags %s ./cmd/shallows/...", strings.Join(buildTags, ",")).Environ("GOBIN", dstdir),
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

func Flatpak(ctx context.Context, op eg.Op) error {
	runtime := shell.Runtime()
	builddir := egenv.WorkingDirectory("fractal", "build", egfs.FindFirst(os.DirFS(egenv.WorkingDirectory("fractal", "build")), "bundle"))

	b := egflatpak.New(
		"space.retrovibe.Daemon", "fractal",
		egflatpak.Option.SDK("org.gnome.Sdk", "47").Runtime("org.gnome.Platform", "47").
			CopyModule(builddir).
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
