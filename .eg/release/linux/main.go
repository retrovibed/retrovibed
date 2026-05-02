package main

import (
	"context"
	"eg/compute/console"
	"eg/compute/flathub"
	"eg/compute/maintainer"
	"eg/compute/release"
	"eg/compute/shallows"
	"eg/compute/tarballs"
	"log"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/egtarball"
)

func main() {
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
				eg.Parallel(
					shallows.Generate,
					console.Generate,
				),
				eg.Parallel(
					console.BuildLinux,
					shallows.Compile(),
				),
				eg.Parallel(
					console.Tests,
					console.Linting,
					shallows.Test(),
				),
				eg.Parallel(
					build(),
					tarballs.Tarchive,
				),
			),
		),
		release.Release(tarinfo()),
	)

	if err != nil {
		log.Fatalln(err)
	}
}

func tarinfo() *tarballs.Build {
	return &tarballs.Build{
		OS:   egenv.String("linux", "EG_COMPUTE_HOST_OS"),
		Arch: egenv.String("amd64", "EG_COMPUTE_HOST_ARCH"),
	}
}

func build() eg.OpFn {
	b := tarinfo()
	archive := tarballs.Retrovibed(b)
	return eg.Sequential(
		shell.Op(
			shell.Newf("mkdir -p %s", egtarball.Path(archive)),
		),
		eg.Sequential(
			console.Install(b),
			shallows.Install(b),
			shell.Op(
				shell.Newf("cp --verbose -R .dist/linux/* %s", egtarball.Path(archive)),
				shell.Newf(
					"cat .dist/linux/usr/share/applications/retrovibed.desktop | envsubst > %s/usr/share/applications/retrovibed.desktop",
					egtarball.Path(archive),
				).
					Environ("VERSION", tarballs.Version()).
					Environ("ARCH", b.Arch),
				shell.Newf(
					"cat usr/share/applications/retrovibed.desktop",
				).Directory(egtarball.Path(archive)),
			),
			flathub.Metainfo(b),
		),
		shell.Op(
			shell.Newf("echo 'tarballing %s -> %s'", egtarball.Path(archive), egtarball.Archive(archive)),
		),
		eg.Parallel(
			release.Tarball(b),
			release.AppImageBuild(b),
		),
		console.FlatpakManifest(b),
	)
}
