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
			// shell.New("gh workflow run release.android.yml --ref main"),
			// shell.New("gh workflow run release.linux.yml --ref main"),
			// TODO: 2026-03-12 needs gpg credential generation enabled - shell.New("gh workflow run release.debian.yml --ref main"),
		),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
