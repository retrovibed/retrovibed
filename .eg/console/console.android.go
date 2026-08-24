package console

import (
	"context"
	"eg/compute/tarballs"
	"fmt"
	"time"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
	"github.com/egdaemon/eg/runtime/wasi/shell"
)

func BuildAndroidAPK(runtime shell.Command) eg.OpFn {
	return func(ctx context.Context, _ eg.Op) error {
		runtime = flutterRuntimev2(runtime)

		commit := eggit.EnvCommit()
		return shell.Run(
			ctx,
			runtime.New(
				fmt.Sprintf("flutter build apk --build-name=%s --build-number=%s --release lib/main.dart", tarballs.Version(), commit.StringReplace("%git.commit.unix%")),
			).Timeout(20*time.Minute),
			runtime.New("mv app-release.apk retrovibed.apk").Timeout(20*time.Minute).Directory(egenv.WorkingDirectory("console/build/app/outputs/apk/release")),
		)
	}
}

func BuildAndroidBundle(runtime shell.Command) eg.OpFn {
	return func(ctx context.Context, _ eg.Op) error {
		runtime = flutterRuntimev2(runtime)
		commit := eggit.EnvCommit()

		return shell.Run(
			ctx,
			runtime.New(
				fmt.Sprintf("flutter build appbundle --build-name=%s --build-number=%s --release lib/main.dart", tarballs.Version(), commit.StringReplace("%git.commit.unix%")),
			).Timeout(20*time.Minute),
			runtime.New("mv app-release.aab retrovibed.aab").Timeout(20*time.Minute).Directory(egenv.WorkingDirectory("console/build/app/outputs/bundle/release")),
		)
	}
}
