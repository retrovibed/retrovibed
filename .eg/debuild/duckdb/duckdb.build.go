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
			sruntime.New("cmake --build ${BUILD_DIRECTORY_REL} --config Release --parallel").Timeout(30*time.Minute),
			sruntime.New("DESTDIR=${BUILD_DIRECTORY} cmake --install ${BUILD_DIRECTORY_REL} --prefix=\"/\""),
		)
	}
}

// bundle deduplicates .o files by basename (first-seen wins) before archiving.
// Required when multiple extension archives share third-party object files (e.g. zstd, fastpfor).
func bundle(sruntime shell.Command) eg.OpFn {
	return func(ctx context.Context, op eg.Op) error {
		return shell.Run(
			ctx,
			sruntime.New("rm -rf bundle && mkdir -p bundle"),
			sruntime.New("rsync -av lib/*.a bundle/"),
			sruntime.New("find bundle -name '*.a' -exec mkdir -p {}.objects \\; -exec mv {} {}.objects \\;"),
			sruntime.New("find bundle -name '*.a' -execdir ar -x {} \\;"),
			sruntime.New("mkdir -p bundle/merged && find bundle -name '*.o' -not -path 'bundle/merged/*' -exec cp --update=none {} bundle/merged/ \\;"),
			// sruntime.New("tree -L 2 bundle/*"),
			sruntime.New("find bundle/merged -name '*.o' | xargs ar cr ${BUILD_DIRECTORY}/libduckdb.a"),
		)
	}
}

// bundlelibtool deduplicates .o files by basename (first-seen wins) before archiving.
// Uses xcrun ar and cp -n (macOS equivalents of ar and cp --update=none).
func bundlelibtool(sruntime shell.Command) eg.OpFn {
	return func(ctx context.Context, op eg.Op) error {
		return shell.Run(
			ctx,
			sruntime.New("rm -rf bundle && mkdir -p bundle"),
			sruntime.New("rsync -av lib/*.a bundle/"),
			sruntime.New("find bundle -name '*.a' -exec mkdir -p {}.objects \\; -exec mv {} {}.objects \\;"),
			sruntime.New("find bundle -name '*.a' -execdir xcrun ar -x {} \\;"),
			sruntime.New("mkdir -p bundle/merged && find bundle -name '*.o' -not -path 'bundle/merged/*' -exec cp -n {} bundle/merged/ \\;"),
			sruntime.New("find bundle/merged -name '*.o' | xargs xcrun libtool -static -o ${BUILD_DIRECTORY}/libduckdb.a"),
		)
	}
}

// compile and put the results into the specified directory.
func CompileDevRuntime() shell.Command {
	builddir := "build/dev"
	absbuilddir := egenv.EphemeralDirectory("duckdb", builddir)
	return shell.Runtime().
		Debug().
		Directory(absbuilddir).
		Environ("BUILD_DIRECTORY", absbuilddir).
		Environ("BUILD_DIRECTORY_REL", builddir)
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
	// absbuilddir := egenv.CacheDirectory("duckdb.bin", builddir)

	return egccache.Runtime().
		Debug().
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
			sruntime.New("rsync -avm --include='*/' --include='libduckdb.a' --exclude='*' ${BUILD_DIRECTORY}/* ${CLONE_DIRECTORY}/"),
		)
	}
}

// clone both the static and shared libraries from the build directory, flattened into CLONE_DIRECTORY.
func CloneBuild(sruntime shell.Command) eg.OpFn {
	return func(ctx context.Context, op eg.Op) error {
		return shell.Run(
			ctx,
			sruntime.New("cp ${BUILD_DIRECTORY}/libduckdb.a ${CLONE_DIRECTORY}/libduckdb.a"),
			sruntime.New("cp ${BUILD_DIRECTORY}/lib/libduckdb.so ${CLONE_DIRECTORY}/libduckdb.so"),
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
		Debug().
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
	fmt.Fprintf(&cmakevars, " -DBUILD_SHELL=OFF")

	builddir := fmt.Sprintf("build/%s-%s", platform, arch)
	absbuilddir := egenv.EphemeralDirectory("duckdb", builddir)

	return egccache.Runtime().
		Debug().
		Directory(absbuilddir).
		Environ("BUILD_DIRECTORY", absbuilddir).
		Environ("BUILD_DIRECTORY_REL", builddir).
		Environ("EXTRA_CMAKE_VARIABLES", cmakevars.String())
}

func CompileIOS(sruntime shell.Command) eg.OpFn {
	return eg.Sequential(
		compile(sruntime, egenv.WorkingDirectory(".dist", "duckdb_ios_config.cmake")),
		bundlelibtool(sruntime),
	)
}

func MaybeBuild(sopath string, runtime shell.Command, bop func(runtime shell.Command) eg.OpFn, clone func(runtime shell.Command) eg.OpFn) eg.OpFn {
	return eg.WhenFn(func(ctx context.Context) bool {
		return !egfs.FileExists(sopath)
	}, eg.Sequential(
		Download,
		bop(runtime),
		clone(runtime.Environ("CLONE_DIRECTORY", filepath.Dir(sopath))),
	),
	)
}
