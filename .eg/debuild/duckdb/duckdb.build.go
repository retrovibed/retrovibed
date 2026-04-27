package duckdb

import (
	"context"
	"eg/compute/android"
	"path/filepath"
	"time"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/egccache"
	"github.com/egdaemon/eg/runtime/x/wasi/egfs"
)

// download the version of duckdb we're using
func Download(ctx context.Context, op eg.Op) error {
	sruntime := shell.Runtime().Directory(egenv.CacheDirectory())
	return shell.Run(
		ctx,
		sruntime.Newf("test -d duckdb || git clone -b v%s --depth 1 https://github.com/duckdb/duckdb.git duckdb", version),
		sruntime.New("md5sum duckdb/src/include/duckdb.h"),
		sruntime.New("echo \"2a20d340931922b25919dd8a870365a9  duckdb/src/include/duckdb.h\" > duckdb.md5"),
		sruntime.New("md5sum -c duckdb.md5"),
	)
}

// compile and put the results into the specified directory.
func Compile(runtime shell.Command, build string) eg.OpFn {
	return func(ctx context.Context, op eg.Op) error {
		return compile(shell.Runtime(), "cmake -G \"Ninja\" -DEXTENSION_STATIC_BUILD=1 -DBUILD_EXTENSIONS=${DUCKDB_EXTENSIONS} -DENABLE_EXTENSION_AUTOLOADING=1 -DENABLE_EXTENSION_AUTOINSTALL=1 -DCMAKE_VERBOSE_MAKEFILE=on -DBUILD_UNITTESTS=0 -DBUILD_SHELL=0 -DCMAKE_BUILD_TYPE=Release .")(ctx, op)
	}
}

func compile(runtime shell.Command, build string) eg.OpFn {
	return func(ctx context.Context, o eg.Op) error {
		sruntime := runtime.
			EnvironFrom(egccache.Env()...).
			Environ("DUCKDB_EXTENSIONS", "autocomplete;json;parquet;icu;inet;fts").
			Directory(egenv.CacheDirectory("duckdb"))
		return shell.Run(
			ctx,
			sruntime.New(build).Timeout(egenv.TTL()),
			sruntime.New("cmake --build build --config Release").Timeout(30*time.Minute),
			sruntime.Newf("DESTDIR=\"%s\" cmake --install build --prefix=\"/\"", egenv.EphemeralDirectory("duckdb")),
		)
	}
}

// clone the compiled results from the ephemeral directory into the specified directory.
func Clone(dir string) eg.OpFn {
	return func(ctx context.Context, op eg.Op) error {
		sruntime := egccache.Runtime().Directory(egenv.CacheDirectory("duckdb"))
		return shell.Run(
			ctx,
			sruntime.Newf("rsync -avm --include='*/' --include='*.so' --include='*.a' --include='*.h' --exclude='*' %s/* %s", egenv.EphemeralDirectory("duckdb"), dir),
		)
	}
}

func CompileAndroid(platform, arch string) eg.OpFn {
	sruntime := egccache.Runtime().
		Environ("ANDROID_NDK", android.NDKPath).
		Environ("ANDROID_PLATFORM", android.Platform). // Now explicitly using the platform variable
		Environ("ANDROID_ABI", arch).
		Environ("PLATFORM_NAME", platform)

	return compile(sruntime, "cmake -G \"Ninja\" -DEXTENSION_STATIC_BUILD=1 -DDUCKDB_EXTRA_LINK_FLAGS=\"-llog -Wl,-z,max-page-size=16384\" -DCMAKE_SHARED_LINKER_FLAGS=\"-Wl,-z,max-page-size=16384\" -DCMAKE_EXE_LINKER_FLAGS=\"-Wl,-z,max-page-size=16384\" -DDUCKDB_EXTRA_LINK_FLAGS=\"-llog\" -DBUILD_EXTENSIONS=${DUCKDB_EXTENSIONS} -DENABLE_EXTENSION_AUTOLOADING=1 -DENABLE_EXTENSION_AUTOINSTALL=1 -DCMAKE_VERBOSE_MAKEFILE=on -DANDROID_PLATFORM=${ANDROID_PLATFORM} -DDUCKDB_EXPLICIT_PLATFORM=${PLATFORM_NAME} -DBUILD_UNITTESTS=0 -DBUILD_SHELL=1 -DANDROID_ABI=${ANDROID_ABI} -DCMAKE_TOOLCHAIN_FILE=${ANDROID_NDK}/build/cmake/android.toolchain.cmake -DCMAKE_BUILD_TYPE=Release -S . -B build")
}

// clone only the necessary files for android.
func CloneAndroid(dir string) eg.OpFn {
	return func(ctx context.Context, op eg.Op) error {
		sruntime := egccache.Runtime().Directory(egenv.CacheDirectory("duckdb"))
		return shell.Run(
			ctx,
			sruntime.Newf("rsync -avm --include='*/' --include='*.so' --include='*.a' --exclude='*' %s/* %s", egenv.EphemeralDirectory("duckdb", "lib"), dir),
		)
	}
}

func CompileDarwin(platform, arch string) eg.OpFn {
	sruntime := egccache.Runtime().
		Environ("PLATFORM_NAME", platform).
		Environ("ARCH", arch)

	return compile(sruntime, "cmake -G \"Ninja\" "+
		"-DCMAKE_OSX_ARCHITECTURES=${ARCH} "+
		"-DDUCKDB_EXTENSION_STATIC_BUILD=1 "+
		"-DEXTENSION_STATIC_BUILD=1 "+
		"-DBUILD_EXTENSIONS=${DUCKDB_EXTENSIONS} "+
		"-DENABLE_EXTENSION_AUTOLOADING=1 "+
		"-DENABLE_EXTENSION_AUTOINSTALL=1 "+
		"-DDUCKDB_EXPLICIT_PLATFORM=${PLATFORM_NAME} "+
		"-DBUILD_UNITTESTS=0 "+
		"-DBUILD_SHELL=0 "+
		"-DCMAKE_BUILD_TYPE=Release "+
		"-S . -B build")
}

// CloneDarwin copies static libraries and headers for macOS builds.
func CloneDarwin(dir string) eg.OpFn {
	return func(ctx context.Context, op eg.Op) error {
		sruntime := egccache.Runtime().Directory(egenv.CacheDirectory("duckdb"))
		return shell.Run(
			ctx,
			sruntime.Newf("mkdir -p %s", dir),
			sruntime.Newf("rsync -avm --include='*/' --include='*.h' --include='*.a' --exclude='*' %s/* %s",
				egenv.EphemeralDirectory("duckdb", "lib"),
				dir,
			),
		)
	}
}

func CompileIOS(platform, arch string) eg.OpFn {
	sruntime := egccache.Runtime().
		Environ("PLATFORM_NAME", platform).
		Environ("ARCH", arch)

	return compile(sruntime, "cmake -G \"Ninja\" "+
		"-DCMAKE_SYSTEM_NAME=iOS "+
		"-DCMAKE_OSX_ARCHITECTURES=${ARCH} "+
		"-DCMAKE_OSX_SYSROOT=$(xcrun --sdk iphoneos --show-sdk-path) "+
		"-DCMAKE_OSX_DEPLOYMENT_TARGET=16.0 "+
		// "-DAMALGAMATION_BUILD=1 "+
		"-DDUCKDB_EXTENSION_STATIC_BUILD=1 "+
		"-DEXTENSION_STATIC_BUILD=1 "+
		"-DBUILD_EXTENSIONS=${DUCKDB_EXTENSIONS} "+
		"-DENABLE_EXTENSION_AUTOLOADING=1 "+
		"-DENABLE_EXTENSION_AUTOINSTALL=1 "+
		"-DDUCKDB_EXPLICIT_PLATFORM=${PLATFORM_NAME} "+
		"-DBUILD_UNITTESTS=0 "+
		"-DBUILD_SHELL=0 "+
		"-DCMAKE_BUILD_TYPE=Release "+
		"-S . -B build")
}

// CloneIOS mirrors CloneAndroid but targets the ios directory
func CloneIOS(dir string) eg.OpFn {
	return func(ctx context.Context, op eg.Op) error {
		sruntime := egccache.Runtime().Directory(egenv.CacheDirectory("duckdb"))
		return shell.Run(
			ctx,
			sruntime.Newf("mkdir -p %s", dir),
			sruntime.Newf("rsync -avm --include='*/' --include='*.h' --include='*.a' --exclude='*' %s/* %s",
				egenv.EphemeralDirectory("duckdb", "lib"),
				dir,
			),
		)
	}
}

func MaybeBuild(sopath string, bop eg.OpFn, clone func(dir string) eg.OpFn) eg.OpFn {
	return eg.WhenFn(func(ctx context.Context) bool {
		return !egfs.FileExists(egenv.WorkingDirectory(sopath))
	}, eg.Sequential(
		Download,
		bop,
		clone(egenv.WorkingDirectory(filepath.Dir(sopath))),
	),
	)
}
