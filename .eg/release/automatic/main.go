package main

import (
	"context"
	"log"
	"time"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
	"github.com/egdaemon/eg/runtime/wasi/egunsafe/ffigit"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/eggithub"
	"github.com/egdaemon/eg/runtime/x/wasi/egpool"
	"github.com/egdaemon/eg/runtime/x/wasi/egx"
)

func main() {
	ctx, done := context.WithTimeout(context.Background(), egenv.TTL())
	defer done()

	gh := shell.Env().Environ("GH_TOKEN", ffigit.Bearer())
	uploads := egpool.New(1, func(cmd string) eg.OpFn {
		return shell.Op(shell.Env().New(cmd).Attempts(3))
	})
	err := eg.Perform(
		ctx,
		eggit.AutoClone,
		eggithub.DraftClean(0),
		shell.Op(
			shell.Env().New("eg login --seed=\"${RETROVIBED_EG_LOGIN_SEED}\""),
			// TODO: Revisit once darwin is working again.
			gh.New("gh workflow run release.ios.yml --ref main").Attempts(3),
			gh.New("gh workflow run release.macosx.yml --ref main").Attempts(3),
		),
		uploads.Run(egpool.Slice(
			"eg compute upload --no-follow --arch=amd64 --cores=8 --memory=6g release/archlinux",
			"eg compute upload --no-follow --ttl=1h --arch=amd64 --cores=8 --memory=12g release/linux",
			"eg compute upload --no-follow --ttl=3h --arch=arm64 --cores=3 --memory=3g release/linux",
			"eg compute upload --no-follow --ttl=1h --arch=amd64 --cores=8 --memory=12g -e EG_GPG_KEYRING_NAME -e EG_GPG_KEYRING_EMAIL -e EG_GPG_KEYRING_SEED release/debian",
			// "eg compute upload --no-follow --arch=amd64 --cores=8 --memory=6g release/retrokiosk",
			// "eg compute upload --ttl=1h --arch=amd64 --cores=8 --memory=12g --gcp-auto release/android", // TODO: android requires gcp credentials.
		)),
		eg.Sequential(
			egx.Retry(0).Delay(30*time.Second).Run(
				eggithub.Promote(
					"retrovibed.darwin.arm64.dmg",
					"PKGBUILD",
					"retrovibed_*.deb",
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
