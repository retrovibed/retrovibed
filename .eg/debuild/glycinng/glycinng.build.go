package glycinng

import (
	"context"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/shell"
)

const (
	version = "0.3.1"
	// commit pinned by the "0.3.1" tag; guards against the tag being moved upstream.
	commit = "70772888bd6a8847b9fe2a382d6c8c6f50edb35c"
)

// Download clones the pinned glycin-ng release and vendors its cargo
// dependencies into the checkout, since Launchpad's builders have no network
// access during the actual package build (see debian/rules). Requires cargo,
// so this must run as a module step inside Runner()'s container, not via a
// bare top-level eg.Perform (the default eg container has no Rust toolchain).
func Download(ctx context.Context, op eg.Op) error {
	sruntime := shell.Runtime().Directory(egenv.CacheDirectory())
	return shell.Run(
		ctx,
		// 3 attempts to deal with racey behavior around cloning the repo multiple times in parallel.
		sruntime.Newf("test -d glycinng || git clone -b %s --depth 1 https://github.com/QaidVoid/glycin-ng.git glycinng", version).Attempts(3),
		sruntime.New("git -C glycinng rev-parse HEAD > glycinng.actual.commit"),
		sruntime.Newf("echo \"%s\" > glycinng.expected.commit", commit),
		sruntime.New("diff glycinng.actual.commit glycinng.expected.commit"),
		sruntime.New("cd glycinng && cargo vendor vendor > /dev/null"),
	)
}
