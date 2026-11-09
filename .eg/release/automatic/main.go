package main

import (
	"context"
	"eg/compute/egx"
	"log"
	"time"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/eggithub"
	"github.com/egdaemon/eg/runtime/x/wasi/egpool"
)

func main() {
	ctx, done := context.WithTimeout(context.Background(), egenv.TTL())
	defer done()

	uploads := egpool.New(4, func(cmd string) eg.OpFn {
		return shell.Op(shell.Env().New(cmd).Attempts(3))
	})
	err := eg.Perform(
		ctx,
		eggit.AutoClone,
		shell.Op(
			shell.Env().New("eg login --seed=\"${RETROVIBED_EG_LOGIN_SEED}\""),
			// TODO: Revisit once darwin is working again.
			// shell.Env().New("gh workflow run release.ios.yml --ref main").Attempts(3),
			// shell.Env().New("gh workflow run release.macosx.yml --ref main").Attempts(3),
		),
		uploads.Run(egpool.Slice(
			"eg compute upload --no-follow --ttl=1h --arch=amd64 --cores=8 --memory=12g release/linux",
			// "eg compute upload --no-follow --ttl=3h --arch=arm64 --cores=3 --memory=3g release/linux",
			// "eg compute upload --no-follow --arch=amd64 --cores=8 --memory=6g release/archlinux",
			// "eg compute upload --no-follow --arch=amd64 --cores=8 --memory=6g release/retrokiosk",
			// "eg compute upload --no-follow --arch=amd64 --cores=8 --memory=6g release/debian",
			// shell.Env().New("eg compute upload release/android"), // TODO: android requires gcp credentials.
		)),
		eg.Sequential(
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
