package main

import (
	"context"
	"eg/compute/console"
	"eg/compute/debuild/duckdb"
	"eg/compute/egapplex"
	"eg/compute/neurals"
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

	apikey := egenv.String("", "RETROVIBED_APPLE_API_KEY")
	issuerid := egenv.String("", "RETROVIBED_APPLE_ISSUER_ID")

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
			egapplex.KeychainP12(
				egenv.Base64(nil, "RETROVIBED_APPLE_SIGNING_KEY"),
				egenv.String("", "RETROVIBED_APPLE_SIGNING_PASSWORD"),
			),
			egapplex.Provision(
				egenv.Base64(nil, "RETROVIBED_APPLE_PROFILE"),
			),
			egapplex.AuthKey(
				apikey,
				egenv.Base64(nil, "RETROVIBED_APPLE_AUTH_KEY"),
			),
			egapplex.UnlockKeychain(egenv.WorkspaceDirectory("apple.signing.keychain")),
			shell.Op(
				flutter.Newf(
					"xcodebuild -exportArchive -archivePath build/ios/archive/Runner.xcarchive -exportPath build/ios/ipa -exportOptionsPlist ios/ExportOptions.plist OTHER_CODE_SIGN_FLAGS=\"--keychain %s\"",
					egenv.WorkspaceDirectory("apple.signing.keychain"),
				).Timeout(10*time.Minute),
			),
			egapplex.Upload(apikey, issuerid, egenv.WorkingDirectory("console", "build/ios/ipa/*.ipa"), "ios"),
		),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
