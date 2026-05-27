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
	"github.com/egdaemon/eg/runtime/x/wasi/egbug"
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
				egbug.Log("DERP DERP DERP 0"),
				eg.Parallel(
					duckdb.MaybeBuild(filepath.Join("dev.native.libs", "libduckdb_bundle.a"), duckdb.Compile(shell.Runtime()), duckdb.Clone),
					neurals.MaybeBuild(filepath.Join("dev.native.libs", "libpredicttext.a"), neurals.Compile, neurals.Clone),
				),
				egbug.Log("DERP DERP DERP 1"),
				console.GenerateDevStaticBinding(shell.Runtime(), egenv.WorkingDirectory("console", "build", "nativelib"), egenv.CacheDirectory("dev.native.libs")),
				egbug.Log("DERP DERP DERP 2"),
				console.Generate,
				egbug.Log("DERP DERP DERP 3"),
				console.BuildLinux,
				egbug.Log("DERP DERP DERP 4"),
			),
		),
		console.RunDev("flutter run -d linux"),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
