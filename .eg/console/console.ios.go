package console

import (
	"context"
	"eg/compute/tarballs"
	"os"
	"path/filepath"
	"time"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
	"github.com/egdaemon/eg/runtime/wasi/shell"
)

// iosDeploymentTarget is the minimum iOS version we build/link against — also
// matches CMAKE_OSX_DEPLOYMENT_TARGET used for duckdb (.eg/debuild/duckdb/duckdb.build.go).
const iosDeploymentTarget = "16.0"

// iosTarget is the clang target triple for device builds. Shared between
// GenIOSCompilerEnv (writes it to ios.compile.env for shell commands to read
// via ${IOS_TARGET}) and CompileIOSBinding (needs the literal value: CGO_CFLAGS/
// CGO_LDFLAGS are env-var values handed directly to clang, never shell-expanded,
// so a ${IOS_TARGET} reference there is never substituted).
const iosTarget = "arm64-apple-ios" + iosDeploymentTarget

const extensionsCMakeConfig = `duckdb_extension_load(inet
    SOURCE_DIR ${CMAKE_CURRENT_LIST_DIR}/../../inet
    LOAD_TESTS
)
duckdb_extension_load(vss
    SOURCE_DIR ${CMAKE_CURRENT_LIST_DIR}/../../vss
)
`

// EnsureDuckDBSource clones DuckDB, inet, and vss extension sources, writes the
// extension_config_local.cmake so inet and vss are built as static extensions.
func EnsureDuckDBSource(duckdbsrc, duckdbinet, duckdbvss string) eg.OpFn {
	return func(ctx context.Context, op eg.Op) error {
		if err := shell.Op(
			shell.Newf("git clone --depth 1 --branch v1.5.3 https://github.com/duckdb/duckdb.git %s", duckdbsrc).Lenient(true),
			shell.Newf("git clone --depth 1 --branch v1.4-andium https://github.com/duckdb/duckdb_inet.git %s", duckdbinet).Lenient(true),
			shell.Newf("git clone --depth 1 --branch v1.5-variegata https://github.com/duckdb/duckdb-vss.git %s", duckdbvss).Lenient(true),
		)(ctx, op); err != nil {
			return err
		}

		return os.WriteFile(filepath.Join(duckdbsrc, "extension", "extension_config_local.cmake"), []byte(extensionsCMakeConfig), 0644)
	}
}

// baseIOSRuntime is flutterRuntimev2 plus the extras the iOS toolchain needs, but
// without the compiler env (CC/CXX/SDKROOT/IOS_TARGET/IPHONEOS_DEPLOYMENT_TARGET) loaded yet — used by
// GenIOSCompilerEnv itself, which is what generates that env file in the first place.
func baseIOSRuntime() shell.Command {
	return flutterRuntimev2(shell.Runtime()).
		Debug().
		Environ("IOSENV", egenv.WorkloadDirectory("ios.compile.env")).
		Environ("LANG", "en_US.UTF-8")
}

// GenIOSCompilerEnv discovers the iOS toolchain (CC/CXX/SDKROOT/IOS_TARGET/IPHONEOS_DEPLOYMENT_TARGET) via
// xcrun and writes it to a file for iosRuntime to read back. Must run before any
// step that calls iosRuntime.
func GenIOSCompilerEnv(ctx context.Context, op eg.Op) error {
	env := baseIOSRuntime().Directory(egenv.WorkloadDirectory())
	return shell.Op(
		env.New("echo \"CC=$(xcrun --sdk iphoneos --find clang)\" | tee ${IOSENV}"),
		env.New("echo \"CXX=$(xcrun --sdk iphoneos --find clang++)\" | tee -a ${IOSENV}"),
		env.New("echo \"SDKROOT=$(xcrun --sdk iphoneos --show-sdk-path)\" | tee -a ${IOSENV}"),
		env.New("echo \"IOS_TARGET="+iosTarget+"\" | tee -a ${IOSENV}"),
		env.New("echo \"IPHONEOS_DEPLOYMENT_TARGET="+iosDeploymentTarget+"\" | tee -a ${IOSENV}"),
	)(ctx, op)
}

func iosCompilerEnv() []string {
	return egenv.MustEnviron(egenv.Build().FromPath(egenv.WorkloadDirectory("ios.compile.env")))
}

func iosRuntime() shell.Command {
	return baseIOSRuntime().EnvironFrom(iosCompilerEnv()...)
}

// CompileIOSBinding cross-compiles the Go binding as a static archive
// (libretrovibed.a) for iOS/arm64, statically linking duckdb and predicttext
// from dev.native.libs.
func CompileIOSBinding(ctx context.Context, _ eg.Op) error {
	runtime := iosRuntime()
	libsdir := egenv.CacheDirectory("dev.native.libs")

	cflags := "-arch arm64"
	ldflags := "-arch arm64 -lc++ -framework CoreFoundation -framework Security " +
		"-L" + libsdir + " -lpredicttext " +
		"-Wl,-force_load," + filepath.Join(libsdir, "libduckdb.a") + " " +
		"-Wl,-force_load," + filepath.Join(libsdir, "libpredicttext.a")

	return shell.Run(
		ctx,
		runtime.New("mkdir -p ${RETROVIBED_SHARED_NATIVE_LIBS_DIRECTORY}"),
		runtime.New("go -C retrovibedbind build -trimpath -buildmode=c-archive --tags duckdb_use_static_lib,retrovibed,neural -o ${RETROVIBED_SHARED_NATIVE_LIBS_DIRECTORY}/libretrovibed.a ./...").
			Timeout(egenv.TTL()).
			Environ("GOOS", "ios").
			Environ("GOARCH", "arm64").
			Environ("CGO_ENABLED", "1").
			Environ("CGO_CFLAGS", cflags).
			Environ("CGO_LDFLAGS", ldflags),
	)
}

// BuildIOS assembles RetrovivedBind.xcframework from the static archives built by
// CompileIOSBinding (statically linked via -force_load, since code_assets'
// StaticLinking LinkMode is not yet supported by the Dart/Flutter SDK), then
// builds the signed-less release IPA.
func BuildIOS(ctx context.Context, _ eg.Op) error {
	runtime := iosRuntime()
	commit := eggit.EnvCommit()
	libsdir := egenv.CacheDirectory("dev.native.libs")

	return shell.Run(
		ctx,
		runtime.New("rm -rf ios/RetrovivedBind.framework ios/RetrovivedBind.xcframework && cp -r ios/RetrovivedBind ios/RetrovivedBind.framework"),
		runtime.Newf(
			"xcrun clang -target ${IOS_TARGET} -isysroot ${SDKROOT} -shared -o RetrovivedBind.framework/RetrovivedBind "+
				"-Wl,-force_load,%s -Wl,-force_load,%s -Wl,-force_load,%s "+
				"-lc++ -lresolv -framework CoreFoundation -framework Security -Wl,-install_name,@rpath/RetrovivedBind.framework/RetrovivedBind",
			filepath.Join(libsdir, "libretrovibed.a"),
			filepath.Join(libsdir, "libduckdb.a"),
			filepath.Join(libsdir, "libpredicttext.a"),
		).Directory(egenv.WorkingDirectory("console", "ios")),
		runtime.New("xcodebuild -create-xcframework -framework RetrovivedBind.framework -output RetrovivedBind.xcframework").Directory(egenv.WorkingDirectory("console", "ios")),
		runtime.New("flutter pub get"),
		runtime.New("pod install").Directory(egenv.WorkingDirectory("console", "ios")),
		runtime.Newf("flutter build ipa --build-name=%s --build-number=%s --no-codesign --release", tarballs.Version(), commit.StringReplace("%git.commit.unix%")).
			Timeout(15*time.Minute),
	)
}

func BuildIOSSimulator(ctx context.Context, _ eg.Op) error {
	runtime := flutterRuntimev2(shell.Runtime())
	return shell.Run(
		ctx,
		runtime.New("flutter build ios --debug --simulator"),
	)
}
