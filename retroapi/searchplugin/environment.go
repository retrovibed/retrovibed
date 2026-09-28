package searchplugin

import (
	"bytes"
	"context"
	"errors"

	"github.com/retrovibed/retrovibed/retroapi/errorsx"
	"github.com/tetratelabs/wazero/sys"
)

// Environment invokes the plugin installed at path as:
//
//	<path> env
//
// returning, verbatim, the .env document it writes to stdout - the
// variables that plugin understands, with a comment per variable describing
// it (see retroapi/envfile for the exact convention). Same contract as
// publishplugin.Registry.Environment: the plugin itself is the schema the
// console renders its settings form from.
//
// Unlike Search this bypasses the shared worker pool; the call is short,
// rare, and interactive. A plugin predating the env command exits non-zero
// (unknown subcommand), which is returned as an error.
func (r *Registry) Environment(ctx context.Context, path string) ([]byte, error) {
	compiled, ok := r.lookup(path)
	if !ok {
		return nil, errorsx.Wrapf(ErrNotLoaded, "path: %s", path)
	}

	cfg, err := r.sandbox(path, path, "env")
	if err != nil {
		return nil, err
	}

	// deliberately anonymous, unlike runSearchJob: an environment read can
	// land while a search against the same plugin is mid-flight, and wazero
	// refuses to instantiate a second module under a name that is already live.
	var stdout bytes.Buffer
	cfg = cfg.WithName("env").WithStdout(&stdout)

	mod, runErr := r.runtime.InstantiateModule(ctx, compiled, cfg)
	if mod != nil {
		errorsx.Log(mod.Close(ctx))
	}

	if exit, ok := errors.AsType[*sys.ExitError](runErr); ok {
		if exit.ExitCode() != 0 {
			return nil, errorsx.Wrapf(exit, "search plugin exited non-zero: %s", path)
		}
	} else if runErr != nil {
		return nil, errorsx.Wrapf(runErr, "unable to run search plugin: %s", path)
	}

	return stdout.Bytes(), nil
}
