package ffiegmacvm

import (
	"context"

	"github.com/egdaemon/eg"
	"github.com/egdaemon/eg/internal/errorsx"
	"github.com/egdaemon/eg/interp/macvm"
	"github.com/egdaemon/eg/runtime/wasi/egunsafe"
)

func Pull(ctx context.Context, name, image, bundle string, options []string) error {
	cc, err := egunsafe.DialControlSocket(ctx)
	if err != nil {
		return err
	}
	_, err = macvm.NewProxyClient(cc).Pull(ctx, &macvm.PullRequest{
		Name:    name,
		Image:   image,
		Bundle:  bundle,
		Options: options,
	})
	return errorsx.Wrap(err, "macvm pull failed")
}

func Build(ctx context.Context, name, ipsw, bundle string, options []string) error {
	cc, err := egunsafe.DialControlSocket(ctx)
	if err != nil {
		return err
	}
	_, err = macvm.NewProxyClient(cc).Build(ctx, &macvm.BuildRequest{
		Name:    name,
		Ipsw:    ipsw,
		Bundle:  bundle,
		Options: options,
	})
	return errorsx.Wrap(err, "macvm build failed")
}

// Run dispatches a single shell command into the VM identified by name.
// The bundle path and runner name match the values used by Build so the
// host's bundle cache is shared across Build/Run/Module calls.
func Run(ctx context.Context, name, bundle, _ string, cmd, options []string) error {
	cc, err := egunsafe.DialControlSocket(ctx)
	if err != nil {
		return err
	}
	_, err = macvm.NewProxyClient(cc).Run(ctx, &macvm.RunRequest{
		Name:    name,
		Bundle:  bundle,
		Command: cmd,
		Options: options,
	})
	return errorsx.Wrap(err, "macvm run failed")
}

// Module dispatches the nested eg interpreter against modulepath inside the VM.
func Module(ctx context.Context, name, bundle, modulepath string, options []string) error {
	cc, err := egunsafe.DialControlSocket(ctx)
	if err != nil {
		return err
	}
	_, err = macvm.NewProxyClient(cc).Module(ctx, &macvm.ModuleRequest{
		Name:    name,
		Bundle:  bundle,
		Module:  modulepath,
		Mdir:    eg.DefaultMountRoot(eg.RuntimeDirectory),
		Options: options,
	})
	return errorsx.Wrap(err, "macvm module failed")
}
