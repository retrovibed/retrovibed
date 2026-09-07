package main

import (
	"context"
	"log"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
	"github.com/egdaemon/eg/runtime/wasi/shell"
)

// plugin is one publishplugin.Registry program: the package to build and
// the filename to install it as. The two differ here, and the filename is
// the one that matters - a plugin's .env, config and cache directories are
// all keyed off it (publishplugin.EnvPath, PluginConfigDir, PluginCacheDir),
// and each package's own documentation names the binary it expects.
type plugin struct {
	pkg  string
	name string
}

// plugins are retrovibed's publishplugin programs; each is its own main
// package built independently to wasip1/wasm.
var plugins = []plugin{
	{pkg: "retroapi/publishplugin-activitypub", name: "lemmy"},
	{pkg: "retroapi/examples/publishplugin-noop", name: "noop"},
}

func runtime() shell.Command {
	return shell.Runtime().Directory(egenv.WorkingDirectory())
}

func wasmdir() string {
	return egenv.EphemeralDirectory("retrovibed", "publishers", "wasm")
}

func main() {
	log.SetFlags(log.Lshortfile | log.LUTC | log.Ltime)
	ctx, done := context.WithTimeout(context.Background(), egenv.TTL())
	defer done()

	log.Println("TTL", egenv.TTL())

	c1 := eg.Container("retrovibed.publishers")

	err := eg.Perform(
		ctx,
		eggit.AutoClone,
		eg.Build(
			c1.BuildFromFile(".eg/release/publishers/Containerfile"),
		),
		eg.Module(
			ctx,
			c1,
			eg.Sequential(
				Setup,
				Build,
				Publish,
			),
		),
	)

	if err != nil {
		log.Fatalln(err)
	}
}

func Setup(ctx context.Context, _ eg.Op) error {
	return shell.Run(
		ctx,
		shell.Newf("mkdir -p %s", wasmdir()),
		runtime().New("retrovibe version"),
		runtime().New("retrovibe identity generate ${RETROVIBED_IDENTITY_SEED}"),
		runtime().New("retrovibe identity show"),
	)
}

// Build cross compiles each plugin from the repository root. Unlike
// retrodscrape's equivalent workload this deliberately leaves GOWORK alone:
// retrovibed vendors at the workspace root (vendor/modules.txt is a
// "## workspace" vendor tree, and go.work is committed), so the workspace is
// what lets these packages - which live in the retroapi module - resolve
// their dependencies offline.
func Build(ctx context.Context, _ eg.Op) error {
	rt := runtime().
		Environ("GOOS", "wasip1").
		Environ("GOARCH", "wasm")

	cmds := make([]shell.Command, 0, len(plugins))
	for _, p := range plugins {
		cmds = append(cmds, rt.Newf("go build -trimpath -o %s/%s.wasm ./%s", wasmdir(), p.name, p.pkg))
	}

	return shell.Run(ctx, cmds...)
}

// Publish tags the built wasm binaries with
// application/vnd.retrovibed.publish.module (matches
// retroapi/mimex.RetrovibedPublishModule) and publishes them into the
// retropublish community via the retrovibed daemon at
// host.containers.internal:9998 (data-plane); community create/info instead
// talk to the deeppool cloud service directly, keyed off the identity
// generated in Setup. See .README.md for what's expected to already be
// running, and for why nothing installs these automatically yet.
func Publish(ctx context.Context, _ eg.Op) error {
	rt := runtime()

	return shell.Run(
		ctx,
		rt.Newf(
			"retrovibe library import directory --no-directory-prefix --mimetype application/vnd.retrovibed.publish.module --insecure --endpoint=\"https://host.containers.internal:9998\" %s"+
				" | retrovibe community info retropublish"+
				" | retrovibe library publish --no-dry-run --endpoint=\"https://host.containers.internal:9998\"",
			wasmdir(),
		).Timeout(egenv.TTL()),
	)
}
