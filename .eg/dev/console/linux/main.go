package main

import (
	"context"
	"eg/compute/console"
	"eg/compute/debuild/duckdb"
	"eg/compute/maintainer"
	"eg/compute/neurals"
	"log"
	"path/filepath"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
	"github.com/egdaemon/eg/runtime/wasi/shell"
)

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	ctx, done := context.WithTimeout(context.Background(), egenv.TTL())
	defer done()

	deb := eg.Container(maintainer.Container)
	err := eg.Perform(
		ctx,
		eggit.AutoClone,
		eg.Build(deb.BuildFromFile(".eg/Containerfile")),
		eg.Module(
			ctx,
			deb,
			eg.Sequential(
				console.Generate,
				eg.Parallel(
					duckdb.MaybeBuild(
						filepath.Join("dev.native.libs", "libduckdb.a"),
						duckdb.CompileDevRuntime(),
						duckdb.Compile,
						duckdb.CloneStaticBuild,
					),
					neurals.MaybeBuild(filepath.Join("dev.native.libs", "libpredicttext.a"), neurals.Compile, neurals.Clone),
				),
				console.GenerateDevStaticBinding(shell.Runtime(), egenv.WorkingDirectory("console", "build", "nativelib"), egenv.CacheDirectory("dev.native.libs")),
				console.BuildLinux,
			),
		),
		console.RunDev("flutter run -d linux"),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
