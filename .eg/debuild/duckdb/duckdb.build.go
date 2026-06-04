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

type cmakeVars struct{ strings.Builder }

func cmake() *cmakeVars { return &cmakeVars{} }

func (c *cmakeVars) set(key, value string) *cmakeVars {
	if c.Len() > 0 {
		c.WriteByte(' ')
	}
	fmt.Fprintf(&c.Builder, `-D%s=%s`, key, value)
	return c
}

// quoted quotes the value. Use only when the flags are embedded directly in a shell
// command string — the shell strips the quotes. Do NOT use for values placed in an
// env var expanded later via ${VAR}: the shell does not re-process quotes during
// variable expansion, so cmake would receive the literal quote characters.
func (c *cmakeVars) quoted(key, value string) *cmakeVars {
	if c.Len() > 0 {
		c.WriteByte(' ')
	}
	fmt.Fprintf(&c.Builder, `-D%s="%s"`, key, value)
	return c
}

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
		cfg := cmake().
			set("EXTENSION_STATIC_BUILD", "1"). // this needs to happen *before* our configurations are loaded. jfc.
			quoted("DUCKDB_EXTENSION_CONFIGS", strings.Join(cmakeconfigs, ";")).
			quoted("BUILD_EXTENSIONS", "inet;autocomplete;json;parquet;icu;vss")
		return shell.Run(
			ctx,
			sruntime.Newf("cmake -G \"Ninja\" -S . -B ${BUILD_DIRECTORY_REL} %s ${EXTRA_CMAKE_VARIABLES}", cfg).
				Timeout(egenv.TTL()),
			sruntime.New("cmake --build ${BUILD_DIRECTORY_REL} --config Release --parallel --verbose").Timeout(30*time.Minute),
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
			sruntime.New("rsync -av --exclude='libdummy_static_extension_loader.a' lib/*.a bundle/"),
			sruntime.New("find bundle -name '*.a' -exec mkdir -p {}.objects \\; -exec mv {} {}.objects \\;"),
			sruntime.New("find bundle -name '*.a' -execdir ar -x {} \\;"),
			sruntime.New("mkdir -p bundle/merged && find bundle -name '*.o' -not -path 'bundle/merged/*' -exec sh -c 'cp --update=none \"$1\" \"bundle/merged/$(md5sum \"$1\" | cut -c1-32).o\"' _ {} \\;"),
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
			sruntime.New("rsync -av --exclude='libdummy_static_extension_loader.a' lib/*.a bundle/"),
			sruntime.New("find bundle -name '*.a' -exec mkdir -p {}.objects \\; -exec mv {} {}.objects \\;"),
			sruntime.New("find bundle -name '*.a' -execdir xcrun ar -x {} \\;"),
			sruntime.New("mkdir -p bundle/merged && find bundle -name '*.o' -not -path 'bundle/merged/*' -exec sh -c 'cp -n \"$1\" \"bundle/merged/$(md5 -q \"$1\").o\"' _ {} \\;"),
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
	cmakevars := cmake().
		set("CMAKE_TOOLCHAIN_FILE", fmt.Sprintf("%s/build/cmake/android.toolchain.cmake", android.NDKPath)).
		set("ANDROID_PLATFORM", android.Platform).
		set("ANDROID_ABI", arch).
		set("DUCKDB_EXPLICIT_PLATFORM", platform).
		set("BUILD_UNITTESTS", "OFF")

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

func CompileDev(sruntime shell.Command) eg.OpFn {
	return eg.Sequential(
		compile(sruntime, egenv.WorkingDirectory(".dist", "duckdb_config.cmake")),
		bundle(sruntime),
	)
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
			sruntime.New("mkdir -p ${CLONE_DIRECTORY}"),
			sruntime.New("rsync -avm --include='*/' --include='libduckdb.a' --exclude='*' ${BUILD_DIRECTORY}/* ${CLONE_DIRECTORY}/"),
		)
	}
}

// clone both the static and shared libraries from the build directory, flattened into CLONE_DIRECTORY.
func CloneBuild(sruntime shell.Command) eg.OpFn {
	return func(ctx context.Context, op eg.Op) error {
		return shell.Run(
			ctx,
			sruntime.New("mkdir -p ${CLONE_DIRECTORY}"),
			sruntime.New("cp ${BUILD_DIRECTORY}/libduckdb.a ${CLONE_DIRECTORY}/libduckdb.a"),
			sruntime.New("cp ${BUILD_DIRECTORY}/lib/libduckdb.* ${CLONE_DIRECTORY}/"),
		)
	}
}

func CompileDarwinRuntime(platform, arch string) shell.Command {
	cmakevars := cmake().
		set("CMAKE_OSX_ARCHITECTURES", arch).
		set("DUCKDB_EXPLICIT_PLATFORM", platform).
		set("BUILD_UNITTESTS", "OFF").
		set("BUILD_SHELL", "OFF")

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
		bundlelibtool(sruntime),
	)
}

func CompileIOSRuntime(platform, arch string) shell.Command {
	cmakevars := cmake().
		set("CMAKE_SYSTEM_NAME", "iOS").
		set("CMAKE_OSX_ARCHITECTURES", arch).
		set("CMAKE_OSX_DEPLOYMENT_TARGET", "16.0").
		set("DUCKDB_EXPLICIT_PLATFORM", platform).
		set("BUILD_UNITTESTS", "OFF").
		set("BUILD_SHELL", "OFF")

	builddir := fmt.Sprintf("build/%s-%s", platform, arch)
	absbuilddir := egenv.EphemeralDirectory("duckdb", builddir)
	// absbuilddir := egenv.CacheDirectory("duckdb.bin", builddir)

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
