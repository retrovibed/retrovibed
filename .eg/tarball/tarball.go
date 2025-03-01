package tarball

import (
	"context"
	"path/filepath"
	"time"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
	"github.com/egdaemon/eg/runtime/wasi/shell"
)

// Directory that will be converted into a tarball
func Directory(paths ...string) string {
	return egenv.RuntimeDirectory("tarball", filepath.Join(paths...))
}

func Archive(dir string) eg.OpFn {
	return func(ctx context.Context, o eg.Op) error {
		c := eggit.EnvCommit()

		ts := time.Now()
		return shell.Run(
			ctx,
			shell.Newf("tar -C %s -Jcvf archive.tar.xz .", dir),
			shell.New("ls -lha ."),
			// shell.Newf("tree %s", dir),
			// shell.Newf("echo gh release create --target %s v%d.%d.%d archive.tar.xz", c.Hash.String(), ts.Year(), ts.Month(), ts.UnixNano()),
			shell.New("gh auth login"),
			shell.Newf("gh release create --target %s v%d.%d.%d archive.tar.xz", c.Hash.String(), ts.Year(), ts.Month(), ts.UnixNano()).Environ(
				"GH_TOKEN", egenv.String("", "EG_GIT_AUTH_ACCESS_TOKEN"),
			),
		)
	}
}
