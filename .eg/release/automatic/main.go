package main

import (
	"context"
	"eg/compute/debian"
	"eg/compute/egx"
	"eg/compute/maintainer"
	"log"
	"time"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/eggithub"
)

func main() {
	ctx, done := context.WithTimeout(context.Background(), egenv.TTL())
	defer done()

	deb := eg.Container(maintainer.Container)
	err := eg.Perform(
		ctx,
		eggit.AutoClone,
		// local builds run first to prevent enqueuing to infrastructure repeatedly
		// due to local failures.
		eg.Parallel(
			eg.Build(deb.BuildFromFile(".eg/Containerfile")),
			debian.Release,
		),
		shell.Op(
			// shell.New("gh workflow run release.ios.yml --ref main").Attempts(3),
			// shell.New("gh workflow run release.macosx.yml --ref main").Attempts(3),
			shell.Env().New("eg login --seed=\"${RETROVIBED_EG_LOGIN_SEED}\""),
			shell.Env().New("eg compute upload --arch=amd64 release/linux").Attempts(3),
			shell.Env().New("eg compute upload --arch=arm64 --cores=3 --memory=2g release/linux").Attempts(3),
			shell.Env().New("eg compute upload release/retrokiosk").Attempts(3),
			shell.Env().New("eg compute upload -e EG_SSH_KEY_SEED=${EG_SSH_KEY_SEED} release/archlinux").Attempts(3),
			// shell.Env().New("eg compute upload release/android"), // TODO: android requires gcp credentials.
		),
		eg.Module(
			ctx,
			deb,
			egx.RetryUntilSuccess(
				30*time.Second,
				eggithub.Promote(
					"PKGBUILD",
					"retrovibed.darwin.arm64.dmg",
					"retrovibed.linux.amd64.AppImage",
					"retrovibed.linux.amd64.AppImage.zsync",
					"retrovibed.linux.amd64.tar.xz",
					"retrovibed.linux.arm64.AppImage",
					"retrovibed.linux.arm64.AppImage.zsync",
					"retrovibed.linux.arm64.tar.xz",
					"retrovibed.source.tar.gz",
				),
			),
		),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
