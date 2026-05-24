package neurals

import (
	"context"
	"path/filepath"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/egcargo"
	"github.com/egdaemon/eg/runtime/x/wasi/egfs"
)

// runtime returns a shell.Command with CARGO_HOME and CARGO_TARGET_DIR set to
// cache locations outside the source tree, rooted in the neurals source dir.
func runtime() shell.Command {
	return egcargo.Runtime().
		Directory(egenv.WorkingDirectory("neurals")).
		Environ("CARGO_TARGET_DIR", egcargo.CacheDirectory("neurals", "target"))
}

// Compile runs cargo build --release then rsyncs libpredicttext.so into dir.
func Compile(dir string) eg.OpFn {
	return func(ctx context.Context, op eg.Op) error {
		sruntime := runtime()
		return shell.Run(
			ctx,
			sruntime.New("cargo build --release").Timeout(egenv.TTL()),
			sruntime.Newf("mkdir -p %s && rsync -a ${CARGO_TARGET_DIR}/release/libpredicttext.so %s/", dir, dir),
		)
	}
}

// Clone copies libpredicttext.so from the cargo target dir into dir.
func Clone(dir string) eg.OpFn {
	return func(ctx context.Context, op eg.Op) error {
		sruntime := runtime()
		return shell.Run(
			ctx,
			sruntime.Newf("mkdir -p %s && rsync -a ${CARGO_TARGET_DIR}/release/libpredicttext.so %s/", dir, dir),
		)
	}
}

// MaybeBuild skips compilation if sopath already exists.
func MaybeBuild(sopath string, clone func(dir string) eg.OpFn) eg.OpFn {
	return eg.WhenFn(func(ctx context.Context) bool {
		return !egfs.FileExists(egenv.WorkingDirectory(sopath))
	}, eg.Sequential(
		Compile(egenv.WorkingDirectory(filepath.Dir(sopath))),
		clone(egenv.WorkingDirectory(filepath.Dir(sopath))),
	))
}
