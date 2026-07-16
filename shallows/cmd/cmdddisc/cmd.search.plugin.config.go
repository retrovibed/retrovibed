package cmdddisc

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/retrovibed/retrovibed/retroapi/searchplugin"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/envx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

// cmdSearchPluginConfig merges -e configuration into an already-installed
// plugin's .env file, without touching its compiled .wasm.
type cmdSearchPluginConfig struct {
	Name string   `arg:"" name:"name" help:"installed plugin name (its .wasm filename, with or without the extension)"`
	Env  []string `flag:"" name:"env" short:"e" type:"envvar" help:"KEY=VALUE pair or file://path merged into the plugin's .env config"`
}

func (t cmdSearchPluginConfig) Run(gctx *cmdopts.Global) error {
	name := strings.TrimSuffix(t.Name, ".wasm")
	plugindir := searchplugin.SearchPluginDir()
	if err := os.MkdirAll(plugindir, 0o700); err != nil {
		return errorsx.Wrapf(err, "unable to create search plugin directory: %s", plugindir)
	}
	dst := filepath.Join(plugindir, name+".env")

	existing, err := envx.FromPath(dst)
	if err != nil {
		return errorsx.Wrapf(err, "unable to read existing plugin configuration: %s", dst)
	}
	if err := envx.WriteFile(dst, envx.Merge(existing, t.Env...)); err != nil {
		return errorsx.Wrapf(err, "unable to write plugin configuration: %s", dst)
	}

	log.Println("search plugin configuration updated", dst)
	return nil
}
