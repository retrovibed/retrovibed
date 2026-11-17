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
	"github.com/egdaemon/eg/runtime/x/wasi/egarch"
	"github.com/egdaemon/eg/runtime/x/wasi/egbug"
	"github.com/egdaemon/eg/runtime/x/wasi/egtarball"
)

// diskdebug prints the state needed to diagnose "no space left on device" failures: the sizes of the mounts (tmpfs vs
// overlay), inodes, the temporary directories in effect, and the memory limit tmpfs is counted against.
func diskdebug(ctx context.Context, op eg.Op) error {
	debug := shell.Runtime().Lenient(true)
	return shell.Run(
		ctx,
		debug.New("du -xsh /tmp/* 2>/dev/null"),
		debug.New("df -h"),
		debug.New("df -ih /tmp"),
		debug.New("env | grep -i tmp"),
		debug.New("ls -lah ${GOTMPDIR} ${TMPDIR} /tmp"),
		debug.New("cat /sys/fs/cgroup/memory.max; free -m"),
	)
}

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
				shallows.NeuralsBuild(),
				eg.Sequential(
					egbug.DebugFailure(console.GenerateFlutter, diskdebug),
					egbug.DebugFailure(console.CompileBinding, diskdebug),
					egbug.DebugFailure(console.BuildLinux, diskdebug),
				),
				egbug.DebugFailure(shallows.Compile(), diskdebug),
				eg.Parallel(
					egbug.DebugFailure(build(), diskdebug),
					tarballs.Tarchive,
				),
			),
		),

		release.Release(),
	)

	if err != nil {
		log.Fatalln(err)
	}
}

func build() eg.OpFn {
	archive := tarballs.Retrovibed()
	return eg.Sequential(
		shell.Op(
			shell.Newf("mkdir -p %s", egtarball.Path(archive)),
		),
		eg.Sequential(
			console.Install(),
			shallows.Install(),
			shell.Op(
				shell.Newf("cp --verbose -R .dist/linux/* %s", egtarball.Path(archive)),
				shell.Newf(
					"cat .dist/linux/usr/share/applications/space.retrovibe.Console.desktop | envsubst > %s/usr/share/applications/space.retrovibe.Console.desktop",
					egtarball.Path(archive),
				).
					Environ("VERSION", tarballs.Version()).
					Environ("ARCH", egarch.Host()),
				shell.Newf(
					"cat usr/share/applications/space.retrovibe.Console.desktop",
				).Directory(egtarball.Path(archive)),
			),
			flathub.Metainfo(),
		),
		shell.Op(
			shell.Newf("echo 'tarballing %s -> %s'", egtarball.Path(archive), egtarball.Archive(archive)),
		),
		eg.Sequential(
			release.Tarball(),
			eg.Parallel(
				release.AppImageBuild(),
				console.FlatpakManifest(),
			),
		),
	)
}
