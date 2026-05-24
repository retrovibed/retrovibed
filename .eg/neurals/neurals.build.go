package neurals

import (
	"context"
	"path/filepath"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/egfs"
)

// runtime returns a shell.Command with CARGO_HOME and CARGO_TARGET_DIR set to
// cache locations outside the source tree, rooted in the neurals source dir.
func runtime() shell.Command {
	return shell.Runtime().
		Directory(egenv.WorkingDirectory("neurals")).
		Environ("CARGO_TARGET_DIR", egenv.CacheDirectory("neurals", "target"))
}

// Compile runs cargo build --release then rsyncs libpredicttext.so into dir.
func Compile(dir string) eg.OpFn {
	return func(ctx context.Context, op eg.Op) error {
		sruntime := runtime()
		return shell.Run(
			ctx,
			sruntime.New("cargo build --release").Timeout(egenv.TTL()),
			sruntime.Newf("mkdir -p %s && rsync -a ${CARGO_TARGET_DIR}/release/libpredicttext.so ${CARGO_TARGET_DIR}/release/libpredicttext.a %s/", dir, dir),
		)
	}
}

// Clone copies libpredicttext.so from the cargo target dir into dir.
func Clone(dir string) eg.OpFn {
	return func(ctx context.Context, op eg.Op) error {
		sruntime := runtime()
		return shell.Run(
			ctx,
			sruntime.Newf("mkdir -p %s && rsync -a ${CARGO_TARGET_DIR}/release/libpredicttext.so ${CARGO_TARGET_DIR}/release/libpredicttext.a %s/", dir, dir),
		)
	}
}

// CompileDarwin builds libpredicttext.dylib for the host macOS architecture,
// rsyncs it into dir, and sets its install name so it can be embedded in an app bundle.
func CompileDarwin(dir string) eg.OpFn {
	return func(ctx context.Context, op eg.Op) error {
		sruntime := runtime()
		return shell.Run(
			ctx,
			sruntime.New("cargo build --release").Timeout(egenv.TTL()),
			sruntime.Newf("mkdir -p %s && rsync -a ${CARGO_TARGET_DIR}/release/libpredicttext.dylib %s/", dir, dir),
			shell.Newf("install_name_tool -id @rpath/libpredicttext.dylib %s/libpredicttext.dylib", dir),
		)
	}
}

// androidRustTarget maps an Android ABI name to its Rust target triple.
func androidRustTarget(abi string) string {
	switch abi {
	case "arm64-v8a":
		return "aarch64-linux-android"
	case "x86_64":
		return "x86_64-linux-android"
	default:
		return abi
	}
}

// CompileAndroid cross-compiles libpredicttext.so for the given Android ABI
// using cargo-ndk and rsyncs the result into dir.
func CompileAndroid(abi, dir string) eg.OpFn {
	return func(ctx context.Context, op eg.Op) error {
		rustTarget := androidRustTarget(abi)
		sruntime := runtime().Environ("ANDROID_NDK_HOME", "/opt/android-sdk/ndk/27.0.12077973")
		return shell.Run(
			ctx,
			sruntime.Newf("cargo ndk --target %s --platform 31 build --release", abi).Timeout(egenv.TTL()),
			sruntime.Newf("mkdir -p %s && rsync -a ${CARGO_TARGET_DIR}/%s/release/libpredicttext.so ${CARGO_TARGET_DIR}/%s/release/libpredicttext.a %s/", dir, rustTarget, rustTarget, dir),
		)
	}
}

// CompileIOS cross-compiles libpredicttext.a for arm64-apple-ios and rsyncs it into dir.
func CompileIOS(dir string) eg.OpFn {
	return func(ctx context.Context, op eg.Op) error {
		sruntime := runtime()
		return shell.Run(
			ctx,
			sruntime.New("cargo build --release --target aarch64-apple-ios").Timeout(egenv.TTL()),
			sruntime.Newf("mkdir -p %s && rsync -a ${CARGO_TARGET_DIR}/aarch64-apple-ios/release/libpredicttext.a %s/", dir, dir),
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
