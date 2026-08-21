package main

import (
	"context"
	"eg/compute/egx"
	"eg/compute/maintainer"
	"log"
	"time"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/egbug"
)

func main() {
	ctx, done := context.WithTimeout(context.Background(), egenv.TTL())
	defer done()

	deb := eg.Container(maintainer.Container)
	err := eg.Perform(
		ctx,
		eggit.AutoClone,
		eg.Parallel(
			eg.Build(deb.BuildFromFile(".eg/Containerfile")),
		),
		eg.Parallel(
			shell.Op(
				// environment propagation messed up -e {NAME} should be sufficient to copy values in.
				// TODO: Revisit once darwin is working again.
				shell.Env().New("gh workflow run release.ios.yml --ref james.lawrence/darwin-ios-fixes").Attempts(3),
				// shell.Env().New("gh workflow run release.macosx.yml --ref james.lawrence/darwin-ios-fixes").Attempts(3),
				// shell.Env().New("eg login --seed=\"${RETROVIBED_EG_LOGIN_SEED}\""),
				// shell.Env().New("eg compute upload --ttl=1h --arch=amd64 release/linux").Attempts(3),
				// shell.Env().New("eg compute upload --ttl=3h --arch=arm64 --cores=3 --memory=2g release/linux").Attempts(3),
				// shell.Env().New("eg compute upload -e EG_SSH_KEY_SEED=${EG_SSH_KEY_SEED} release/archlinux").Attempts(3),
				// shell.Env().New("eg compute upload -e EG_GPG_KEYRING_NAME=\"${EG_GPG_KEYRING_NAME}\" -e EG_GPG_KEYRING_EMAIL=\"${EG_GPG_KEYRING_EMAIL}\" -e EG_GPG_KEYRING_SEED=\"${EG_GPG_KEYRING_SEED}\" release/retrokiosk").Attempts(3).Debug(),
				// shell.Env().New("eg compute upload release/android"), // TODO: android requires gcp credentials.
			),
		),
		eg.Module(
			ctx,
			deb,
			eg.Sequential(
				egx.RetryUntilSuccess(
					30*time.Second,
					egbug.Log("derp done"),
					// eggithub.Promote(
					// 	"PKGBUILD",
					// 	"retrovibed.darwin.arm64.dmg",
					// 	"retrovibed.linux.amd64.AppImage",
					// 	"retrovibed.linux.amd64.AppImage.zsync",
					// 	"retrovibed.linux.amd64.tar.xz",
					// 	"retrovibed.linux.arm64.AppImage",
					// 	"retrovibed.linux.arm64.AppImage.zsync",
					// 	"retrovibed.linux.arm64.tar.xz",
					// 	"retrovibed.source.tar.gz",
					// ),
				),
			),
		),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
