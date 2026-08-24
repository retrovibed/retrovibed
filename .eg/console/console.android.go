package console

import (
	"context"
	"eg/compute/android"
	"eg/compute/tarballs"
	"fmt"
	"path/filepath"
	"time"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
	"github.com/egdaemon/eg/runtime/wasi/shell"
)

// androidRuntime applies flutterRuntimev2's defaults and then re-asserts the
// native-libs directory on top: flutterRuntimev2 unconditionally points
// RETROVIBED_SHARED_NATIVE_LIBS_DIRECTORY / NIX_RETROVIBED_SHARED_NATIVE_LIBS
// at its own dev.native.libs default, so applying it after the caller's
// runtime (as BuildAndroidAPK/Bundle do) would otherwise silently clobber
// whatever value the caller already set.
func androidRuntime(runtime shell.Command) shell.Command {
	runtime = flutterRuntimev2(runtime)
	dir := android.JNILibRoot()
	return runtime.
		Environ("RETROVIBED_SHARED_NATIVE_LIBS_DIRECTORY", dir).
		Environ("NIX_RETROVIBED_SHARED_NATIVE_LIBS", filepath.Join(dir, "example.so"))
}

func BuildAndroidAPK(runtime shell.Command) eg.OpFn {
	return func(ctx context.Context, _ eg.Op) error {
		runtime = androidRuntime(runtime)

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
		runtime = androidRuntime(runtime)
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
