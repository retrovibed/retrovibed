package cmdcommunity

import (
	"log"
	"os"
	"path/filepath"

	"github.com/retrovibed/retrovibed/retroapi/publishplugin"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
)

// cmdPublisherLink installs a second copy of an already-installed plugin
// under a new name, as a symlink rather than a duplicated binary.
//
// A publish plugin's entire identity - its .env, its config.d and cache.d
// directories, its row in the catalog - is derived from the filename it was
// installed as, so two links to one module are two independently configured
// publishers. That is how a single lemmy plugin posts to several different
// communities: one link per target, each with its own credentials.
type cmdPublisherLink struct {
	Source string   `arg:"" name:"source" help:"name of the already installed plugin to link to"`
	Name   string   `arg:"" name:"name" help:"name to install the new copy as"`
	Env    []string `flag:"" name:"env" short:"e" type:"envvar" optional:"" help:"KEY=VALUE pair or file://path written into the new copy's .env config"`
}

func (t cmdPublisherLink) Run(gctx *cmdopts.Global) (err error) {
	plugindir, name, err := publisherPaths(t.Name)
	if err != nil {
		return err
	}

	_, source, err := publisherPaths(t.Source)
	if err != nil {
		return err
	}

	if source == name {
		return errorsx.Errorf("a plugin cannot be linked to itself: %s", name)
	}

	src := filepath.Join(plugindir, source+".wasm")
	if _, err = os.Stat(src); err != nil {
		return errorsx.Wrapf(err, "unable to find the plugin being linked to: %s", src)
	}

	dst := filepath.Join(plugindir, name+".wasm")

	// replacing an existing link is the normal way to repoint one, but
	// refuse to clobber a real binary - that would silently discard an
	// install.
	if info, serr := os.Lstat(dst); serr == nil {
		if info.Mode()&os.ModeSymlink == 0 {
			return errorsx.Errorf("refusing to replace an installed plugin with a link: %s", dst)
		}

		if err = os.Remove(dst); err != nil {
			return errorsx.Wrapf(err, "unable to replace the existing link: %s", dst)
		}
	}

	// deliberately relative: publish.d is copied and relocated between
	// hosts (and into containers), and an absolute link would dangle the
	// moment it moved.
	if err = os.Symlink(source+".wasm", dst); err != nil {
		return errorsx.Wrapf(err, "unable to link publish plugin: %s", dst)
	}

	if err = mergeEnv(publishplugin.EnvPath(dst), t.Env); err != nil {
		errorsx.Log(errorsx.Wrap(fsx.IgnoreIsNotExist(os.Remove(dst)), "unable to remove the incomplete link"))
		return err
	}

	log.Println("publish plugin linked", dst, "->", src)
	log.Println("restart the daemon to make it selectable for a community")

	return nil
}
