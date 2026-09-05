package cmdddisc

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/retrovibed/retrovibed/retroapi/searchplugin"
	"github.com/retrovibed/retrovibed/retroapi/userx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/envx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/pluginx"
)

// cmdSearchPluginInstall compiles a searchplugin repository to wasip1/wasm
// and installs it into searchplugin.SearchPluginDir(userx.DefaultConfigDir(...)),
// where the registry's fsnotify watch picks it up automatically.
type cmdSearchPluginInstall struct {
	Repository string   `arg:"" name:"repository" help:"local directory or git URL of a go module implementing the searchplugin protocol"`
	Branch     string   `flag:"" name:"branch" help:"branch to clone when repository is a git URL" default:"main"`
	Package    string   `flag:"" name:"package" help:"import path (relative to the module root) of the plugin's main package" default:"."`
	Name       string   `flag:"" name:"name" optional:"" help:"filename (without .wasm) to install as; defaults to the repository's directory/base name"`
	Env        []string `flag:"" name:"env" short:"e" type:"envvar" optional:"" help:"KEY=VALUE pair or file://path merged into the installed plugin's .env config"`
	Bake       []string `flag:"" name:"bake" short:"b" optional:"" help:"main.KEY=VALUE pair baked into the binary at compile time via -ldflags -X; repeatable"`
}

func (t cmdSearchPluginInstall) Run(gctx *cmdopts.Global) (err error) {
	dir := t.Repository
	if info, serr := os.Stat(dir); serr != nil || !info.IsDir() {
		if dir, err = pluginx.Clone(gctx.Context, t.Repository, t.Branch); err != nil {
			return err
		}
		defer os.RemoveAll(dir)
	}

	name := t.Name
	if name == "" {
		name = strings.TrimSuffix(filepath.Base(strings.TrimSuffix(t.Repository, "/")), ".git")
	}
	name = strings.TrimSuffix(name, ".wasm")

	plugindir := searchplugin.SearchPluginDir(userx.DefaultConfigDir(userx.DefaultRelRoot()))
	if err = os.MkdirAll(plugindir, 0o700); err != nil {
		return errorsx.Wrapf(err, "unable to create search plugin directory: %s", plugindir)
	}

	dst := filepath.Join(plugindir, name+".wasm")
	tmp := dst + ".tmp"
	defer os.Remove(tmp)

	if err = pluginx.Compile(gctx.Context, dir, t.Package, tmp, t.Bake); err != nil {
		return errorsx.Wrapf(err, "unable to compile search plugin: %s", t.Repository)
	}

	if err = searchplugin.VerifyWasmMagicPath(tmp); err != nil {
		return err
	}

	if err = os.Rename(tmp, dst); err != nil {
		return errorsx.Wrapf(err, "unable to install search plugin: %s", dst)
	}

	envdst := filepath.Join(plugindir, name+".env")
	existing, err := envx.FromPath(envdst)
	if err != nil {
		return errorsx.Wrapf(err, "unable to read existing plugin configuration: %s", envdst)
	}
	if err = envx.WriteFile(envdst, envx.Merge(existing, t.Env...)); err != nil {
		return errorsx.Wrapf(err, "unable to write plugin configuration: %s", envdst)
	}

	log.Println("search plugin installed", dst)
	return nil
}
