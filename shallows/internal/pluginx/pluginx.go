// Package pluginx holds the mechanics both plugin install commands share -
// search plugins (retroapi/searchplugin) and publish plugins
// (retroapi/publishplugin) are installed the same way: fetch a go module,
// cross compile its main package to wasip1/wasm, drop the result in the
// registry's watched directory. Only the directory and the sidecar
// bookkeeping differ, and those stay in the commands themselves.
package pluginx

import (
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

// Clone clones uri (a git URL) into a new temp directory and returns its
// path; the caller is responsible for removing it.
func Clone(ctx context.Context, uri, branch string) (dir string, err error) {
	dir, err = os.MkdirTemp("", "retrovibed-plugin-install-*")
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

// Compile cross-compiles pkg (relative to dir, a go module root) to a
// wasip1/wasm binary at output. Each bake entry ("main.KEY=VALUE") becomes
// a `-X` linker flag, letting an install bake install-specific defaults
// into the binary (see retroapi/examples/searchplugin-noop's source var for
// what a plugin does with a baked value).
func Compile(ctx context.Context, dir, pkg, output string, bake []string) error {
	xflags := make([]string, 0, len(bake)*2)
	for _, kv := range bake {
		xflags = append(xflags, "-X", kv)
	}

	cmd := exec.CommandContext(ctx, "go", "build", "-trimpath", "-ldflags", strings.Join(xflags, " "), "-o", output, pkg)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOOS=wasip1", "GOARCH=wasm", "GOWORK=off")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
