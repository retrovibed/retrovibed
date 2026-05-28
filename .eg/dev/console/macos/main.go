package main

import (
	"context"
	"eg/compute/console"
	"log"
	"time"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/egbug"
	"github.com/egdaemon/eg/runtime/x/wasi/eggolang"
	"github.com/egdaemon/eg/runtime/x/wasi/egrunnermacos"
)

// Image is the Tart-format OCI image the runner pulls. Cirrus Labs's
// macos-*-base ships with SSH (admin/admin), Xcode CLI tools, and Homebrew
// but no dev-language toolchain — the first BuildAndRun step installs the
// pieces this workload needs via brew.
const Image = "ghcr.io/cirruslabs/macos-sequoia-xcode:latest"

// Builds the native library and Flutter macOS app inside a pulled macOS guest
// VM so the host's Homebrew/Xcode state is irrelevant.
func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	ctx, done := context.WithTimeout(context.Background(), egenv.TTL())
	defer done()

	vm := egrunnermacos.New("retrovibed-macos").PullFrom(Image)

	err := eg.Perform(
		ctx,
		eg.Build(vm),
		eg.Module(ctx, vm.ToModuleRunner(), BuildAndRun),
	)

	if err != nil {
		log.Fatalln(err)
	}
}

// BuildAndRun executes inside the guest. The transpiler lifts this into the
// generated sub-module's main; only top-level identifiers reachable from this
// file survive, so the runtime setup lives here rather than in main().
func BuildAndRun(ctx context.Context, _ eg.Op) error {
	runtime := shell.Runtime().
		EnvironFrom(eggolang.Env()...).
		Environ("PUB_CACHE", egenv.CacheDirectory(".eg", "dart"))
	flutter := runtime.Directory(egenv.WorkingDirectory("console"))

	return eg.Perform(
		ctx,
		eg.Sequential(
			shell.Op(
				runtime.New("brew install go duckdb gpgme llvm flutter ffmpeg@7 cocoapods").Timeout(20*time.Minute),
			),
			egbug.DebugFailure(
				shell.Op(
					flutter.New("mkdir -p build/nativelib"),
					flutter.New("go -C retrovibedbind build -buildmode=c-shared --tags localdev -o ../build/nativelib/libretrovibed.dylib ./...").Timeout(5*time.Minute),
				),
				eg.Sequential(
					egbug.Log("failed to build native library"),
					Debug(runtime),
				),
			),
			egbug.DebugFailure(
				shell.Op(
					flutter.New(retry(`dart run ffigen --config ffigen.yaml --compiler-opts "-I$(clang --print-resource-dir)/include"`)),
				),
				eg.Sequential(
					egbug.Log("failed to generate ffi bindings"),
					Debug(runtime),
				),
			),
			shell.Op(
				flutter.New(retry(`flutter pub get`)),
				flutter.New(retry(`flutter build macos --debug --dart-define=EG_VM=true`)),
				flutter.New("cp build/nativelib/libretrovibed.dylib build/macos/Build/Products/Debug/retrovibed.app/Contents/MacOS/libretrovibed.dylib"),
				flutter.New("codesign --force --sign - build/macos/Build/Products/Debug/retrovibed.app"),
			),
			console.RunDevAt(
				"flutter run -d macos --use-application-binary=build/macos/Build/Products/Debug/retrovibed.app",
				"$(route -n get default | awk '/gateway:/{print $2}')",
			),
		),
	)
}

// retry wraps a shell command in a bounded retry loop. Tart's virtio-fs share
// occasionally returns ENOENT/EIO under load (pubspec.yaml briefly missing,
// hooks_runner stdout EIO); a fresh syscall after a short sleep usually
// succeeds.
func retry(cmd string) string {
	return `for i in 1 2 3; do ` + cmd + ` && exit 0; sleep 2; done; exit 1`
}

func Debug(runtime shell.Command) eg.OpFn {
	return shell.Op(
		runtime.New("which go && go version"),
		runtime.New("which flutter && flutter --version"),
		runtime.New("which dart && dart --version"),
		runtime.New("ls -la /Library/Developer/CommandLineTools/usr/lib/clang/").Lenient(true),
		runtime.New("xcode-select -p"),
	)
}
