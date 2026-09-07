package cmdcommunity

import (
	"os"
	"path/filepath"

	"github.com/retrovibed/retrovibed/retroapi/publishplugin"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

// cmdPublisherEnv prints the variables an installed plugin declares - the
// plugin's own answer to "what can I be configured with", which is also
// what the console renders its settings form from. Piping it into a plugin's
// .env is a reasonable way to start configuring one:
//
//	retrovibe community publisher env lemmy-movies > ~/.config/retrovibed/publish.d/lemmy-movies.env
type cmdPublisherEnv struct {
	Name string `arg:"" name:"name" help:"installed plugin name (its .wasm filename, with or without the extension)"`
}

func (t cmdPublisherEnv) Run(gctx *cmdopts.Global) error {
	plugindir, name, err := publisherPaths(t.Name)
	if err != nil {
		return err
	}

	// NewRegistry loads everything already sitting in publish.d before it
	// returns, so the plugin being asked about is resident by the time
	// Environment runs - no explicit Load needed. Same construction the
	// discovery locate command uses to run search plugins without a
	// daemon.
	reg, err := publishplugin.NewRegistry(gctx.Context)
	if err != nil {
		return errorsx.Wrap(err, "unable to start publish plugin registry")
	}

	declared, err := reg.Environment(gctx.Context, filepath.Join(plugindir, name+".wasm"))
	if err != nil {
		return err
	}

	_, err = os.Stdout.Write(declared)

	return err
}
