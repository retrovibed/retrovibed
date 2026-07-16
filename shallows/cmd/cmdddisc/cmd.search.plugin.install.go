package cmdddisc

import (
	"bytes"
	"context"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/retrovibed/retrovibed/retroapi/searchplugin"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/envx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

// wasmMagic is the 4-byte binary header every valid wasm module begins with
// ("\0asm"), used to verify a freshly compiled plugin before it's installed
// into search.d.
var wasmMagic = []byte{0x00, 0x61, 0x73, 0x6D}

// cmdSearchPluginInstall compiles a searchplugin repository to wasip1/wasm
// and installs it into searchplugin.SearchPluginDir(), where the
// registry's fsnotify watch picks it up automatically.
type cmdSearchPluginInstall struct {
	Repository string   `arg:"" name:"repository" help:"local directory or git URL of a go module implementing the searchplugin protocol"`
	Branch     string   `flag:"" name:"branch" help:"branch to clone when repository is a git URL" default:"main"`
	Package    string   `flag:"" name:"package" help:"import path (relative to the module root) of the plugin's main package" default:"."`
	Name       string   `flag:"" name:"name" optional:"" help:"filename (without .wasm) to install as; defaults to the repository's directory/base name"`
	Env        []string `flag:"" name:"env" short:"e" type:"envvar" optional:"" help:"KEY=VALUE pair or file://path merged into the installed plugin's .env config"`
}

func (t cmdSearchPluginInstall) Run(gctx *cmdopts.Global) (err error) {
	dir := t.Repository
	if info, serr := os.Stat(dir); serr != nil || !info.IsDir() {
		if dir, err = cloneRepository(gctx.Context, t.Repository, t.Branch); err != nil {
			return err
		}
		defer os.RemoveAll(dir)
	}

	name := t.Name
	if name == "" {
		name = strings.TrimSuffix(filepath.Base(strings.TrimSuffix(t.Repository, "/")), ".git")
	}
	name = strings.TrimSuffix(name, ".wasm")

	plugindir := searchplugin.SearchPluginDir()
	if err = os.MkdirAll(plugindir, 0o700); err != nil {
		return errorsx.Wrapf(err, "unable to create search plugin directory: %s", plugindir)
	}

	dst := filepath.Join(plugindir, name+".wasm")
	tmp := dst + ".tmp"
	defer os.Remove(tmp)

	if err = compileWasm(gctx.Context, dir, t.Package, tmp); err != nil {
		return errorsx.Wrapf(err, "unable to compile search plugin: %s", t.Repository)
	}

	if err = verifyWasmMagic(tmp); err != nil {
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

// cloneRepository clones uri (a git URL) into a new temp directory and
// returns its path; the caller is responsible for removing it.
func cloneRepository(ctx context.Context, uri, branch string) (dir string, err error) {
	dir, err = os.MkdirTemp("", "retrovibed-search-install-*")
	if err != nil {
		return "", errorsx.Wrap(err, "unable to create temp directory")
	}

	cmd := exec.CommandContext(ctx, "git", "clone", "--branch", branch, "--single-branch", "--depth", "1", uri, dir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", errorsx.Wrapf(err, "unable to clone repository: %s", uri)
	}

	return dir, nil
}

// compileWasm cross-compiles pkg (relative to dir, a go module root) to a
// wasip1/wasm binary at output.
func compileWasm(ctx context.Context, dir, pkg, output string) error {
	cmd := exec.CommandContext(ctx, "go", "build", "-trimpath", "-o", output, pkg)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOOS=wasip1", "GOARCH=wasm")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// verifyWasmMagic confirms path begins with the wasm binary header, so an
// unexpected build output never gets installed into search.d.
func verifyWasmMagic(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return errorsx.Wrapf(err, "unable to open compiled plugin: %s", path)
	}
	defer f.Close()

	magic := make([]byte, len(wasmMagic))
	if _, err := f.Read(magic); err != nil {
		return errorsx.Wrapf(err, "unable to read compiled plugin: %s", path)
	}

	if !bytes.Equal(magic, wasmMagic) {
		return errorsx.Errorf("compiled output is not a valid wasm module: %s", path)
	}

	return nil
}
