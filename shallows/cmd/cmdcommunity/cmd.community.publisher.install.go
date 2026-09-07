package cmdcommunity

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/retrovibed/retrovibed/retroapi/publishplugin"
	"github.com/retrovibed/retrovibed/retroapi/userx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/envx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/pluginx"
)

// cmdPublisherInstall compiles a publishplugin repository to wasip1/wasm
// and installs it into publishplugin.PublishPluginDir(userx.DefaultConfigDir(...)),
// where the registry's fsnotify watch picks it up automatically. The
// search-plugin equivalent is cmdddisc's cmdSearchPluginInstall; the two
// differ only in which directory they target.
type cmdPublisherInstall struct {
	Repository string   `arg:"" name:"repository" help:"local directory or git URL of a go module implementing the publishplugin protocol"`
	Branch     string   `flag:"" name:"branch" help:"branch to clone when repository is a git URL" default:"main"`
	Package    string   `flag:"" name:"package" help:"import path (relative to the module root) of the plugin's main package" default:"."`
	Name       string   `flag:"" name:"name" optional:"" help:"filename (without .wasm) to install as; defaults to the repository's directory/base name"`
	Env        []string `flag:"" name:"env" short:"e" type:"envvar" optional:"" help:"KEY=VALUE pair or file://path merged into the installed plugin's .env config"`
	Bake       []string `flag:"" name:"bake" short:"b" optional:"" help:"main.KEY=VALUE pair baked into the binary at compile time via -ldflags -X; repeatable"`
}

func (t cmdPublisherInstall) Run(gctx *cmdopts.Global) (err error) {
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

	plugindir, name, err := publisherPaths(name)
	if err != nil {
		return err
	}

	dst := filepath.Join(plugindir, name+".wasm")
	tmp := dst + ".tmp"
	defer os.Remove(tmp)

	if err = pluginx.Compile(gctx.Context, dir, t.Package, tmp, t.Bake); err != nil {
		return errorsx.Wrapf(err, "unable to compile publish plugin: %s", t.Repository)
	}

	if err = publishplugin.VerifyWasmMagicPath(tmp); err != nil {
		return err
	}

	if err = os.Rename(tmp, dst); err != nil {
		return errorsx.Wrapf(err, "unable to install publish plugin: %s", dst)
	}

	if err = mergeEnv(publishplugin.EnvPath(dst), t.Env); err != nil {
		return err
	}

	log.Println("publish plugin installed", dst)
	log.Println("restart the daemon to make it selectable for a community")

	return nil
}

// publisherPaths resolves name to the publish.d directory it belongs in and
// its sanitized form, creating the directory when it doesn't exist yet.
// Sanitizing through publishplugin.SanitizeName is what keeps a name from
// escaping publish.d, regardless of what was typed.
func publisherPaths(name string) (plugindir string, sanitized string, err error) {
	sanitized = publishplugin.SanitizeName(strings.TrimSuffix(name, ".wasm"))
	if sanitized == "" {
		return "", "", errorsx.Errorf("unusable plugin name: %s", name)
	}

	plugindir = publishplugin.PublishPluginDir(userx.DefaultConfigDir(userx.DefaultRelRoot()))
	if err = os.MkdirAll(plugindir, 0o700); err != nil {
		return "", "", errorsx.Wrapf(err, "unable to create publish plugin directory: %s", plugindir)
	}

	return plugindir, sanitized, nil
}

// mergeEnv folds updates into the .env at path, leaving anything already
// configured there that wasn't mentioned alone.
func mergeEnv(path string, updates []string) error {
	existing, err := envx.FromPath(path)
	if err != nil {
		return errorsx.Wrapf(err, "unable to read existing plugin configuration: %s", path)
	}

	return errorsx.Wrapf(envx.WriteFile(path, envx.Merge(existing, updates...)), "unable to write plugin configuration: %s", path)
}
