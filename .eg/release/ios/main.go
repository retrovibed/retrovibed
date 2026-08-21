package main

import (
	"context"
	"eg/compute/console"
	"eg/compute/debuild/duckdb"
	"eg/compute/neurals"
	"eg/compute/release"
	"log"
	"time"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/egbug"
)

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	ctx, done := context.WithTimeout(context.Background(), egenv.TTL())
	defer done()

	flutter := shell.Runtime().Directory(egenv.WorkingDirectory("console"))
	err := eg.Perform(
		ctx,
		eg.Sequential(
			console.GenIOSCompilerEnv,
			eg.Parallel(
				duckdb.MaybeBuild(
					egenv.CacheDirectory("dev.native.libs", "libduckdb.a"),
					duckdb.CompileIOSRuntime("ios_arm64", "arm64"),
					duckdb.CompileIOS,
					duckdb.CloneStaticBuild,
				),
				neurals.CompileIOS(egenv.CacheDirectory("dev.native.libs")),
			),
			console.GenerateFlutter,
			egbug.DebugFailure(
				console.CompileIOSBinding,
				egbug.Log("bindings failed to build for iOS"),
			),
			egbug.DebugFailure(
				console.BuildIOS,
				egbug.Log("flutter failed to build iOS app"),
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
			shell.Op(
				shell.Newf("security unlock-keychain -p %s %s", egenv.RunID(), egenv.WorkspaceDirectory("apple.signing.keychain")),
				flutter.Newf(
					"xcodebuild -exportArchive -archivePath build/ios/archive/Runner.xcarchive -exportPath build/ios/ipa -exportOptionsPlist ios/ExportOptions.plist OTHER_CODE_SIGN_FLAGS=\"--keychain %s\"",
					egenv.WorkspaceDirectory("apple.signing.keychain"),
				).Timeout(10*time.Minute),
				flutter.New("xcrun altool --upload-app --type ios -f build/ios/ipa/*.ipa --apiKey ${RETROVIBED_APPLE_API_KEY} --apiIssuer ${RETROVIBED_APPLE_ISSUER_ID}").
					Environ("RETROVIBED_APPLE_ISSUER_ID", egenv.String("", "RETROVIBED_APPLE_ISSUER_ID")).
					Environ("RETROVIBED_APPLE_API_KEY", egenv.String("", "RETROVIBED_APPLE_API_KEY")),
			),
		),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
