package main

import (
	"context"
	"log"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
	"github.com/egdaemon/eg/runtime/wasi/shell"
)

func main() {
	ctx, done := context.WithTimeout(context.Background(), egenv.TTL())
	defer done()

	err := eg.Perform(
		ctx,
		eggit.AutoClone,
		shell.Op(
			shell.New("gh workflow run release.ios.yml --ref main"),
			shell.New("gh workflow run release.macosx.yml --ref main"),
			shell.New("mkdir -p ~/.ssh && eg login --seed=\"${RETROVIBED_EG_LOGIN_SEED}\""),
			shell.New("eg compute upload release/linux"),
			shell.New("eg compute upload --secret ${RETROVIBED_RELEASE_AUTOMATIC_SECRET} release/debian"),
			shell.New("eg compute upload --secret ${RETROVIBED_RELEASE_AUTOMATIC_SECRET} release/archlinux"),
			// shell.New("eg compute upload release/android"), // TODO: android requires gcp credentials.
		),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
