package main

import (
	"context"
	"eg/compute/console"
	"eg/compute/debuild/duckdb"
	"eg/compute/release"
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
		EnvironFrom(eggolang.Env()...).
		Environ("LANG", "en_US.UTF-8").
		Environ("PUB_CACHE", egenv.CacheDirectory(".eg", "dart"))

	return runtime.Directory(egenv.WorkingDirectory("console"))
}

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	ctx, done := context.WithTimeout(context.Background(), egenv.TTL())
	defer done()

	flutter := flutterRuntime()

	err := eg.Perform(
		ctx,
		eg.Sequential(
			duckdb.MaybeBuild(".eg.cache/duckdb/.arm64/libduckdb_static.a", duckdb.CompileIOS("ios_arm64", "arm64"), duckdb.CloneIOS),
			// Copy all static libs to the ios folder so Xcode/Pods can find them directly
			shell.Op(shell.Newf("cp .eg.cache/duckdb/.arm64/*.a console/ios/")),
			console.GenerateFlutter,
			egbug.DebugFailure(
				shell.Op(
					flutter.Newf(
						"CC=\"$(xcrun --sdk iphoneos --find clang)\" "+
							"CXX=\"$(xcrun --sdk iphoneos --find clang++)\" "+
							"SDKROOT=\"$(xcrun --sdk iphoneos --show-sdk-path)\" "+
							"CGO_CFLAGS=\"-target arm64-apple-ios16.0 -isysroot $(xcrun --sdk iphoneos --show-sdk-path) -I%[1]s -DDUCKDB_STATIC_BUILD\" "+
							"CGO_LDFLAGS=\"-target arm64-apple-ios16.0 -isysroot $(xcrun --sdk iphoneos --show-sdk-path) -lc++ -framework CoreFoundation -framework Security\" "+
							"go -C retrovibedbind build -trimpath -buildmode=c-archive --tags duckdb_use_static_lib -o ../ios/libretrovibed.a ./...",
						egenv.CacheDirectory("duckdb", ".arm64"),
					).
						Environ("GOOS", "ios").
						Environ("GOARCH", "arm64").
						Environ("CGO_ENABLED", "1").
						Timeout(5*time.Minute),
				),
				shell.Op(shell.New("echo 'go failed to build for iOS'")),
			),
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
			iosbuild,
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

func iosbuild(ctx context.Context, op eg.Op) error {
	flutter := flutterRuntime()
	commit := eggit.EnvCommit()

	return eg.Sequential(
		egbug.DebugFailure(
			shell.Op(
				flutter.New("cp ios/libretrovibed.h ios/Classes/libretrovibed.h"),
				flutter.New("flutter pub get"),
				flutter.New("pod install").Directory(egenv.WorkingDirectory("console", "ios")),
				flutter.New(commit.StringReplace("flutter build ipa --build-name=%git.commit.year%.%git.commit.month%.%git.commit.day% --build-number=%git.commit.unix% --no-codesign --release --build-number=%git.commit.unix%")).
					Timeout(15*time.Minute),
			),
			shell.Op(
				shell.New("echo flutter failed to build iOS app"),
			),
		),
	)(ctx, op)
}
