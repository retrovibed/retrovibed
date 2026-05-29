package duckdb

import (
	"context"
	"eg/compute/android"
	"fmt"
	"path/filepath"
	"strings"
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
		// 3 attempts to deal with racey behavior around cloning the repo multiple times in parallel.
		sruntime.Newf("test -d duckdb || git clone -b v%s --depth 1 https://github.com/duckdb/duckdb.git duckdb", version).Attempts(3),
		sruntime.New("md5sum duckdb/src/include/duckdb.h"),
		sruntime.New("echo \"fcdba922a5ef1ac7373134cb915d204b  duckdb/src/include/duckdb.h\" > duckdb.md5"),
		sruntime.New("md5sum -c duckdb.md5"),
	)
}

// compile and put the results into the specified directory.
func CompileDevRuntime() shell.Command {
	return shell.Runtime()
}

// compile and put the results into the specified directory.
func Compile(runtime shell.Command) eg.OpFn {
	return func(ctx context.Context, op eg.Op) error {
		return compile(runtime, egenv.WorkingDirectory(".dist", "duckdb_config.cmake"))(ctx, op)
	}
}

func compile(runtime shell.Command, cmakeconfigs ...string) eg.OpFn {
	return func(ctx context.Context, o eg.Op) error {
		sruntime := runtime.
			Directory(egenv.CacheDirectory("duckdb"))

		return shell.Run(
			ctx,
			sruntime.New("cmake -G \"Ninja\" -S . -B ${BUILD_DIRECTORY_REL} ${EXTRA_CMAKE_VARIABLES}").
				Environ("DUCKDB_EXTENSIONS", "inet").
				Environ(
					"EXTENSION_CONFIGS", strings.Join(cmakeconfigs, ";"),
				).Timeout(egenv.TTL()),
			sruntime.New("cmake --build ${BUILD_DIRECTORY_REL} --config Release").Timeout(30*time.Minute),
			sruntime.New("DESTDIR=${BUILD_DIRECTORY} cmake --install ${BUILD_DIRECTORY_REL} --prefix=\"/\""),
		)
	}
}

// bundle replicates `make bundle-library` for a specific cmake binary directory.
// the runtime must already be set to the proper directory before invoking.
func bundle(sruntime shell.Command) eg.OpFn {
	return func(ctx context.Context, op eg.Op) error {
		return shell.Run(
			ctx,
			sruntime.New("rm -rf bundle && mkdir -p bundle"),
			sruntime.New("cp lib/*.a bundle/."),
			sruntime.New("find bundle -name '*.a' -exec mkdir -p {}.objects \\; -exec mv {} {}.objects \\;"),
			sruntime.New("find bundle -name '*.a' -execdir ar -x {} \\;"),
			// bundle-library-o: archive all .o files into libduckdb.a
			sruntime.New("cd bundle && echo ./*/*.o | xargs ar cr ${BUILD_DIRECTORY}/libduckdb.a"),
		)
	}
}

func CompileAndroidRuntime(platform, arch string) shell.Command {
	var cmakevars strings.Builder
	fmt.Fprintf(&cmakevars, "-DCMAKE_TOOLCHAIN_FILE=%s/build/cmake/android.toolchain.cmake", android.NDKPath)
	fmt.Fprintf(&cmakevars, " -DANDROID_PLATFORM=%s", android.Platform)
	fmt.Fprintf(&cmakevars, " -DANDROID_ABI=%s", arch)
	fmt.Fprintf(&cmakevars, " -DDUCKDB_EXPLICIT_PLATFORM=%s", platform)
	fmt.Fprintf(&cmakevars, " -DBUILD_UNITTESTS=OFF")

	builddir := fmt.Sprintf("build/%s", platform)
	absbuilddir := egenv.EphemeralDirectory("duckdb", builddir)

	return egccache.Runtime().Debug().
		Directory(absbuilddir).
		Environ("BUILD_DIRECTORY", absbuilddir).
		Environ("BUILD_DIRECTORY_REL", builddir).
		Environ("EXTRA_CMAKE_VARIABLES", cmakevars.String())
}

func CompileAndroid(sruntime shell.Command) eg.OpFn {
	return eg.Sequential(
		compile(sruntime, egenv.WorkingDirectory(".dist", "duckdb_android_config.cmake")),
		bundle(sruntime),
	)
}

// clone only the necessary files for a static build
func CloneStaticBuild(sruntime shell.Command) eg.OpFn {
	return func(ctx context.Context, op eg.Op) error {
		return shell.Run(
			ctx,
			sruntime.New("rsync -avm --include='*/' --include='libduckdb.a' --exclude='*' ${BUILD_DIRECTORY}/* ${CLONE_DIRECTORY}"),
		)
	}
}

func CompileDarwinRuntime(platform, arch string) shell.Command {
	var cmakevars strings.Builder
	fmt.Fprintf(&cmakevars, "-DCMAKE_OSX_ARCHITECTURES=%s", arch)
	fmt.Fprintf(&cmakevars, " -DDUCKDB_EXPLICIT_PLATFORM=%s", platform)
	fmt.Fprintf(&cmakevars, " -DBUILD_UNITTESTS=OFF")

	builddir := fmt.Sprintf("build/%s-%s", platform, arch)
	absbuilddir := egenv.EphemeralDirectory("duckdb", builddir)

	return egccache.Runtime().
		Directory(absbuilddir).
		Environ("BUILD_DIRECTORY", absbuilddir).
		Environ("BUILD_DIRECTORY_REL", builddir).
		Environ("EXTRA_CMAKE_VARIABLES", cmakevars.String())
}

func CompileDarwin(sruntime shell.Command) eg.OpFn {
	return eg.Sequential(
		compile(sruntime, egenv.WorkingDirectory(".dist", "duckdb_darwin_config.cmake")),
		bundle(sruntime),
	)
}

func CompileIOSRuntime(platform, arch string) shell.Command {
	var cmakevars strings.Builder
	fmt.Fprintf(&cmakevars, "-DCMAKE_SYSTEM_NAME=iOS")
	fmt.Fprintf(&cmakevars, " -DCMAKE_OSX_ARCHITECTURES=%s", arch)
	fmt.Fprintf(&cmakevars, " -DCMAKE_OSX_DEPLOYMENT_TARGET=16.0")
	fmt.Fprintf(&cmakevars, " -DDUCKDB_EXPLICIT_PLATFORM=%s", platform)
	fmt.Fprintf(&cmakevars, " -DBUILD_UNITTESTS=OFF")

	builddir := fmt.Sprintf("build/%s-%s", platform, arch)
	absbuilddir := egenv.EphemeralDirectory("duckdb", builddir)

	return egccache.Runtime().
		Directory(absbuilddir).
		Environ("BUILD_DIRECTORY", absbuilddir).
		Environ("BUILD_DIRECTORY_REL", builddir).
		Environ("EXTRA_CMAKE_VARIABLES", cmakevars.String())
}

func CompileIOS(sruntime shell.Command) eg.OpFn {
	return eg.Sequential(
		compile(sruntime, egenv.WorkingDirectory(".dist", "duckdb_ios_config.cmake")),
		bundle(sruntime),
	)
}

func MaybeBuild(sopath string, runtime shell.Command, bop func(runtime shell.Command) eg.OpFn, clone func(runtime shell.Command) eg.OpFn) eg.OpFn {
	return eg.WhenFn(func(ctx context.Context) bool {
		return !egfs.FileExists(egenv.WorkingDirectory(sopath))
	}, eg.Sequential(
		Download,
		bop(runtime),
		clone(runtime.Environ("CLONE_DIRECTORY", egenv.WorkingDirectory(filepath.Dir(sopath)))),
	),
	)
}
