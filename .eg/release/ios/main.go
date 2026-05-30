package main

import (
	"context"
	"eg/compute/console"
	"eg/compute/debuild/duckdb"
	"eg/compute/neurals"
	"eg/compute/release"
	"eg/compute/tarballs"
	"fmt"
	"log"
	"time"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/egbug"
	"github.com/egdaemon/eg/runtime/x/wasi/eggolang"
)

func flutterRuntime() shell.Command {
	runtime := shell.Runtime().
		Debug().
		EnvironFrom(eggolang.Env()...).
		Environ("IOSENV", egenv.WorkloadDirectory("ios.compile.env")).
		Environ("LANG", "en_US.UTF-8").
		Environ("PUB_CACHE", egenv.CacheDirectory(".eg", "dart"))

	return runtime.Directory(egenv.WorkingDirectory("console"))
}

func geniOSCompilerEnv(ctx context.Context, op eg.Op) error {
	env := flutterRuntime().
		Directory(egenv.WorkloadDirectory())
	return shell.Op(
		env.New("echo \"CC=$(xcrun --sdk iphoneos --find clang)\" | tee ${IOSENV}"),
		env.New("echo \"CXX=$(xcrun --sdk iphoneos --find clang++)\" | tee -a ${IOSENV}"),
		env.New("echo \"SDKROOT=$(xcrun --sdk iphoneos --show-sdk-path)\" | tee -a ${IOSENV}"),
		env.New("echo \"IOS_TARGET=arm64-apple-ios16.0\" | tee -a ${IOSENV}"),
	)(ctx, op)
}

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	ctx, done := context.WithTimeout(context.Background(), egenv.TTL())
	defer done()

	flutter := flutterRuntime()
	err := eg.Perform(
		ctx,
		eg.Sequential(
			geniOSCompilerEnv,
			eg.Parallel(
				duckdb.MaybeBuild(
					egenv.WorkingDirectory("console", "ios", "libduckdb.a"),
					duckdb.CompileIOSRuntime("ios_arm64", "arm64"),
					duckdb.CompileIOS,
					duckdb.CloneStaticBuild,
				),
				neurals.CompileIOS(egenv.WorkingDirectory("console", "ios")),
			),
			console.GenerateFlutter,
			gobuild,
			iosbuild,
			release.Keychain(
				egenv.String("", "RETROVIBED_APPLE_SIGNING_KEY"),
				egenv.String("", "RETROVIBED_APPLE_SIGNING_PASSWORD"),
			),
			release.ProvisioningProfile(
				egenv.String("", "RETROVIBED_APPLE_PROFILE"),
			),
			release.AuthKey(
				egenv.String("", "RETROVIBED_APPLE_API_KEY"),
				egenv.String("", "RETROVIBED_APPLE_AUTH_KEY"),
			),
			shell.Op(
				shell.Newf("security unlock-keychain -p %s %s", egenv.RunID(), egenv.WorkspaceDirectory("apple.signing.keychain")),
				flutter.Newf(
					"xcodebuild -exportArchive -archivePath build/ios/archive/Runner.xcarchive -exportPath build/ios/ipa -exportOptionsPlist ios/ExportOptions.plist OTHER_CODE_SIGN_FLAGS=\"--keychain %s\"",
					egenv.WorkspaceDirectory("apple.signing.keychain"),
				).Timeout(10*time.Minute),
				flutter.New("xcrun altool --upload-app --type ios -f build/ios/ipa/*.ipa --apiKey ${RETROVIBED_APPLE_API_KEY} --apiIssuer ${RETROVIBED_APPLE_ISSUER_ID}").
					Environ("RETROVIBED_APPLE_ISSUER_ID", egenv.String("", "RETROVIBED_APPLE_ISSUER_ID")).
					Environ("RETROVIBED_APPLE_API_KEY", egenv.String("", "RETROVIBED_APPLE_API_KEY")).
					Debug(),
			),
		),
	)

	if err != nil {
		log.Fatalln(err)
	}
}

func iOSCompilerEnv() []string {
	return egenv.MustEnviron(egenv.Build().FromPath(egenv.WorkloadDirectory("ios.compile.env")))
}

func gobuild(ctx context.Context, op eg.Op) error {
	flutter := flutterRuntime().
		EnvironFrom(iOSCompilerEnv()...)
	return egbug.DebugFailure(
		shell.Op(
			flutter.Newf(
				"CGO_CFLAGS=\"-target ${IOS_TARGET} -I%[1]s\" "+
					"CGO_LDFLAGS=\"-target ${IOS_TARGET} -lc++ -framework CoreFoundation -framework Security -L$(pwd)/console/ios -lpredicttext -Wl,-force_load,libduckdb.a -Wl,-force_load,libpredicttext.a\" "+
					"go -C retrovibedbind build -trimpath -buildmode=c-archive --tags duckdb_use_static_lib,retrovibed,neural -o ../ios/libretrovibed.a ./...",
				egenv.CacheDirectory("duckdb", ".arm64"),
			).
				Environ("GOOS", "ios").
				Environ("GOARCH", "arm64").
				Environ("CGO_ENABLED", "1").
				Timeout(30*time.Minute),
		),
		shell.Op(shell.New("echo 'go failed to build for iOS'")),
	)(ctx, op)
}

func iosbuild(ctx context.Context, op eg.Op) error {
	flutter := flutterRuntime().EnvironFrom(iOSCompilerEnv()...)
	commit := eggit.EnvCommit()

	return eg.Sequential(
		egbug.DebugFailure(
			shell.Op(
				flutter.New("rm -rf ios/RetrovivedBind.framework ios/RetrovivedBind.xcframework && cp -r ios/RetrovivedBind ios/RetrovivedBind.framework"),
				flutter.New("xcrun clang -target ${IOS_TARGET} -isysroot ${SDKROOT} -shared -o RetrovivedBind.framework/RetrovivedBind -Wl,-force_load,libretrovibed.a -Wl,-force_load,libduckdb.a -Wl,-force_load,libpredicttext.a -lc++ -lresolv -framework CoreFoundation -framework Security -Wl,-install_name,@rpath/RetrovivedBind.framework/RetrovivedBind").Directory(egenv.WorkingDirectory("console", "ios")),
				flutter.New("xcodebuild -create-xcframework -framework RetrovivedBind.framework -output RetrovivedBind.xcframework").Directory(egenv.WorkingDirectory("console", "ios")),
				flutter.New("flutter pub get"),
				flutter.New("pod install").Directory(egenv.WorkingDirectory("console", "ios")),
				flutter.New(fmt.Sprintf("flutter build ipa --build-name=%s --build-number=%s --no-codesign --release", tarballs.Version(), commit.StringReplace("%git.commit.unix%"))).
					Timeout(15*time.Minute),
			),
			shell.Op(
				shell.New("echo flutter failed to build iOS app"),
			),
		),
	)(ctx, op)
}
