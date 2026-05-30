// Package egrunnermacos runs eg modules inside a macOS guest VM
// driven by Apple's Virtualization.framework via the host-side macvm proxy.
//
// The user-facing surface mirrors runtime/wasi/eg.ContainerRunner: declare
// a Runner with FromIPSW (first-boot restore) or FromBundle (reuse cached
// bundle), then drive it through eg.Build / eg.Module just like a container.
package egrunnermacos

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/egdaemon/eg"
	"github.com/egdaemon/eg/internal/envx"
	"github.com/egdaemon/eg/internal/errorsx"
	"github.com/egdaemon/eg/runtime/wasi/egunsafe/ffiegmacvm"
)

type option []string

func (option) workdir(dir string) option {
	return []string{"-w", dir}
}

func (option) envvar(k, v string) option {
	if v == "" {
		return []string{"-e", k}
	}
	return []string{"-e", fmt.Sprintf("%s=%s", k, v)}
}

func (option) literal(args ...string) option {
	return args
}

// New returns a Runner that will boot a macOS VM identified by name.
func New(name string) Runner {
	return Runner{
		name:  name,
		built: &sync.Once{},
	}
}

// Runner declares a macOS VM and the command/module it will host.
type Runner struct {
	name    string
	image   string
	ipsw    string
	bundle  string
	cmd     []string
	options []option
	built   *sync.Once
}

// PullFrom pulls a Tart-format macOS image from an OCI registry (the parity
// path with ContainerRunner.PullFrom). Compatible images live at e.g.
// ghcr.io/cirruslabs/macos-sequoia-base:latest.
func (t Runner) PullFrom(image string) Runner {
	t.image = image
	return t
}

// FromIPSW restores macOS into a fresh bundle from the given IPSW. The
// resulting guest is unprovisioned (OOBE) — useful for building golden
// images, not for direct workload execution. Use PullFrom for parity with
// container workflows.
func (t Runner) FromIPSW(path string) Runner {
	t.ipsw = path
	return t
}

// FromBundle boots from an existing VM bundle directory and skips both
// IPSW restore and OCI pull.
func (t Runner) FromBundle(path string) Runner {
	t.bundle = path
	return t
}

// Command is the shell command to execute inside the guest.
func (t Runner) Command(s string) Runner {
	t.cmd = strings.Split(s, " ")
	return t
}

func (t Runner) OptionLiteral(args ...string) Runner {
	t.options = append(t.options, option(nil).literal(args...))
	return t
}

func (t Runner) OptionWorkingDirectory(dir string) Runner {
	t.options = append(t.options, option{}.workdir(dir))
	return t
}

func (t Runner) OptionEnvVar(k string) Runner {
	t.options = append(t.options, option{}.envvar(k, ""))
	return t
}

func (t Runner) OptionEnv(k, v string) Runner {
	t.options = append(t.options, option{}.envvar(k, v))
	return t
}

func (t Runner) flatOptions() []string {
	out := make([]string, 0, len(t.options))
	for _, o := range t.options {
		out = append(out, o...)
	}
	return out
}

// CompileWith provisions the bundle the first time it is called. The path is
// chosen by the most-specific option set on the Runner:
//
//  1. PullFrom — fetch a Tart-format image from an OCI registry (container parity).
//  2. FromIPSW — drive vz.MacOSInstaller against a local IPSW.
//  3. FromBundle — assume the bundle already exists; no-op.
//
// Idempotent across runs because the host's bundle cache is keyed by runner name.
func (t Runner) CompileWith(ctx context.Context) (err error) {
	t.built.Do(func() {
		if t.image == "" && t.ipsw == "" && t.bundle == "" {
			err = errorsx.Errorf("macvm runner %q requires PullFrom, FromIPSW, or FromBundle", t.name)
			return
		}
		if t.image != "" {
			err = errorsx.Wrapf(ffiegmacvm.Pull(ctx, t.name, t.image, t.bundle, t.flatOptions()), "unable to pull bundle: %s", t.name)
			return
		}
		if t.ipsw != "" {
			err = errorsx.Wrapf(ffiegmacvm.Build(ctx, t.name, t.ipsw, t.bundle, t.flatOptions()), "unable to restore bundle: %s", t.name)
		}
	})
	return err
}

// RunWith boots the VM and executes the configured Command via SSH.
func (t Runner) RunWith(ctx context.Context, mpath string) error {
	return errorsx.Wrapf(
		ffiegmacvm.Run(ctx, t.name, t.bundle, mpath, t.cmd, t.flatOptions()),
		"unable to run macvm: %s", t.name,
	)
}

// ToModuleRunner returns a Runner variant that dispatches Module rather than Run.
func (t Runner) ToModuleRunner() ModuleRunner {
	return ModuleRunner{Runner: t}
}

// ModuleRunner runs the workload as a nested eg module inside the guest.
type ModuleRunner struct {
	Runner
}

func (t ModuleRunner) RunWith(ctx context.Context, mpath string) error {
	opts := append(t.flatOptions(), "-e", fmt.Sprintf("%s=%d", eg.EnvComputeModuleNestedLevel, envx.Int(-1, eg.EnvComputeModuleNestedLevel)+1))
	return errorsx.Wrapf(
		ffiegmacvm.Module(ctx, t.name, t.bundle, mpath, opts),
		"unable to run macvm module: %s", t.name,
	)
}
