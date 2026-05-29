package duckdb

import (
	"context"
	"eg/compute/android"
	"fmt"
	"path/filepath"
	"strings"

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
		// 3 attempts to deal with racey behavior around cloning the repo multiple times in parallel.
		sruntime.Newf("test -d duckdb || git clone -b v%s --depth 1 https://github.com/duckdb/duckdb.git duckdb", version).Attempts(3),
		sruntime.New("md5sum duckdb/src/include/duckdb.h"),
		sruntime.New("echo \"fcdba922a5ef1ac7373134cb915d204b  duckdb/src/include/duckdb.h\" > duckdb.md5"),
		sruntime.New("md5sum -c duckdb.md5"),
	)
}

// compile and put the results into the specified directory.
func Compile(runtime shell.Command) eg.OpFn {
	return func(ctx context.Context, op eg.Op) error {
		return compile(shell.Runtime(), egenv.WorkingDirectory(".dist", "duckdb_config.cmake"))(ctx, op)
	}
}

func compile(runtime shell.Command, cmakeconfigs ...string) eg.OpFn {
	return func(ctx context.Context, o eg.Op) error {
		sruntime := runtime.
			EnvironFrom(egccache.Env()...).
			Directory(egenv.CacheDirectory("duckdb"))
		return shell.Run(
			ctx,
			sruntime.New("GEN=ninja make bundle-library").
				Environ("DUCKDB_EXTENSIONS", "inet").
				Environ(
					"EXTENSION_CONFIGS", strings.Join(cmakeconfigs, ";"),
				).Timeout(egenv.TTL()),
			sruntime.Newf("mkdir -p %s", egenv.EphemeralDirectory("duckdb")),
			sruntime.Newf("cp build/release/libduckdb_bundle.a %s/", egenv.EphemeralDirectory("duckdb")),
			sruntime.Newf("cp build/release/src/libduckdb.so %s/", egenv.EphemeralDirectory("duckdb")),
		)
	}
}

// clone the compiled results from the ephemeral directory into the specified directory.
func Clone(dir string) eg.OpFn {
	return func(ctx context.Context, op eg.Op) error {
		sruntime := egccache.Runtime().Directory(egenv.CacheDirectory("duckdb"))
		return shell.Run(
			ctx,
			sruntime.Newf("ls -lha %s/*.a", egenv.EphemeralDirectory("duckdb")),
			sruntime.Newf("mkdir -p %s", dir),
			sruntime.Newf("rsync -avm --include='*/' --include='*.so' --include='*.a' --include='*.h' --exclude='*' %s/* %s", egenv.EphemeralDirectory("duckdb"), dir),
		)
	}
}

func CompileAndroid(platform, arch string) eg.OpFn {
	var cmakevars strings.Builder
	fmt.Fprintf(&cmakevars, "-DCMAKE_TOOLCHAIN_FILE=%s/build/cmake/android.toolchain.cmake", android.NDKPath)

	sruntime := egccache.Runtime().
		Environ("ANDROID_NDK", android.NDKPath).
		Environ("ANDROID_PLATFORM", android.Platform).
		Environ("ANDROID_ABI", arch).
		Environ("DUCKDB_PLATFORM", platform).
		Environ("EXTRA_CMAKE_VARIABLES", cmakevars.String())

	return compile(sruntime, egenv.WorkingDirectory(".dist", "duckdb_android_config.cmake"))
}

// clone only the necessary files for android.
func CloneAndroid(dir string) eg.OpFn {
	return func(ctx context.Context, op eg.Op) error {
		sruntime := egccache.Runtime().Directory(egenv.CacheDirectory("duckdb"))
		return shell.Run(
			ctx,
			sruntime.Newf("ls -lha %s/*.a", egenv.EphemeralDirectory("duckdb", "lib")),
			sruntime.Newf("rsync -avm --include='*/' --include='*.so' --include='*.a' --exclude='*' %s/* %s", egenv.EphemeralDirectory("duckdb", "lib"), dir),
		)
	}
}

func CompileDarwin(platform, arch string) eg.OpFn {
	var cmakevars strings.Builder
	fmt.Fprintf(&cmakevars, "-DCMAKE_OSX_ARCHITECTURES=%s", arch)

	sruntime := egccache.Runtime().
		Environ("DUCKDB_PLATFORM", platform).
		Environ("EXTRA_CMAKE_VARIABLES", cmakevars.String())

	return compile(sruntime, egenv.WorkingDirectory(".dist", "duckdb_darwin_config.cmake"))
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
	var cmakevars strings.Builder
	fmt.Fprintf(&cmakevars, "-DCMAKE_SYSTEM_NAME=iOS")
	fmt.Fprintf(&cmakevars, " -DCMAKE_OSX_ARCHITECTURES=%s", arch)
	fmt.Fprintf(&cmakevars, " -DCMAKE_OSX_DEPLOYMENT_TARGET=16.0")

	sruntime := egccache.Runtime().
		Environ("DUCKDB_PLATFORM", platform).
		Environ("EXTRA_CMAKE_VARIABLES", cmakevars.String())

	return compile(sruntime, egenv.WorkingDirectory(".dist", "duckdb_ios_config.cmake"))
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
		return !egfs.FileExists(egenv.CacheDirectory(sopath))
	}, eg.Sequential(
		Download,
		bop,
		clone(egenv.CacheDirectory(filepath.Dir(sopath))),
	),
	)
}
