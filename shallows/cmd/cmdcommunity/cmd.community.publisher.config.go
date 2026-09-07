package cmdcommunity

import (
	"log"
	"path/filepath"

	"github.com/retrovibed/retrovibed/retroapi/publishplugin"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
)

// cmdPublisherConfig merges -e configuration into an already-installed
// plugin's .env file, without touching its compiled .wasm.
type cmdPublisherConfig struct {
	Name string   `arg:"" name:"name" help:"installed plugin name (its .wasm filename, with or without the extension)"`
	Env  []string `flag:"" name:"env" short:"e" type:"envvar" help:"KEY=VALUE pair or file://path merged into the plugin's .env config"`
}

func (t cmdPublisherConfig) Run(gctx *cmdopts.Global) error {
	plugindir, name, err := publisherPaths(t.Name)
	if err != nil {
		return err
	}

	dst := publishplugin.EnvPath(filepath.Join(plugindir, name+".wasm"))
	if err := mergeEnv(dst, t.Env); err != nil {
		return err
	}

	log.Println("publish plugin configuration updated", dst)

	return nil
}
