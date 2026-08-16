package retrozsync

import (
	"context"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/shell"
)

// download the version of github.com/cph6/zsync we're packaging and vendor its
// module dependencies (golang.org/x/crypto, golang.org/x/sync) so the actual
// package build has no network access requirement — launchpad builders block
// network access during the build step.
func Download(ctx context.Context, op eg.Op) error {
	sruntime := shell.Runtime().Directory(egenv.CacheDirectory())
	return shell.Run(
		ctx,
		// 3 attempts to deal with racey behavior around cloning the repo multiple times in parallel.
		sruntime.Newf("test -d retrozsync || git clone -b v%s --depth 1 https://github.com/cph6/zsync.git retrozsync", version).Attempts(3),
		sruntime.New("md5sum retrozsync/go.mod"),
		sruntime.New("echo \"cb929c8563349d0778576b4224fdb40b  retrozsync/go.mod\" > retrozsync.md5"),
		sruntime.New("md5sum -c retrozsync.md5"),
		shell.Runtime().Directory(egenv.CacheDirectory("retrozsync")).New("go mod vendor"),
	)
}
