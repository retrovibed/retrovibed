package main

import (
	"context"
	"eg/compute/console"
	"eg/compute/debuild/duckdb"
	"eg/compute/neurals"
	"eg/compute/release"
	"eg/compute/tarballs"
	"log"
	"path/filepath"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/egbug"
	"github.com/egdaemon/eg/runtime/x/wasi/eggithub"
	"github.com/egdaemon/eg/runtime/x/wasi/egtarball"
)

func tarinfo() *tarballs.Build {
	return &tarballs.Build{
		OS:   egenv.String("darwin", "EG_COMPUTE_HOST_OS"),
		Arch: egenv.String("arm64", "EG_COMPUTE_HOST_ARCH"),
	}
}

// Baremetal command for darwin due to macosx nonsense for no cloud vms.
func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	ctx, done := context.WithTimeout(context.Background(), egenv.TTL())
	defer done()

	tarballapp := filepath.Join(egtarball.Path(tarballs.Retrovibed(tarinfo())), "retrovibed.app")
	appstoreapp := egenv.CacheDirectory("retrovibed.appstore.app")
	dmgpath := egenv.CacheDirectory("retrovibed.darwin.arm64.dmg")
	pkgpath := egenv.CacheDirectory("retrovibed.darwin.arm64.pkg")
	entitlements := egenv.WorkingDirectory("console", "macos", "Runner", "Release.entitlements")
	keychainPath := egenv.WorkspaceDirectory("apple.signing.keychain")

	apikey := egenv.String("", "RETROVIBED_APPLE_API_KEY")
	issuerid := egenv.String("", "RETROVIBED_APPLE_ISSUER_ID")

	err := eg.Perform(
		ctx,
		eg.Sequential(
			// eg.Parallel(
			duckdb.MaybeBuild(
				egenv.CacheDirectory("dev.native.libs", "libduckdb.a"),
				duckdb.CompileDarwinRuntime("osx_arm64", "arm64"),
				duckdb.CompileDarwin,
				duckdb.CloneBuild,
			),
			neurals.CompileDarwin(egenv.CacheDirectory("dev.native.libs")),
			// ),
			egbug.DirectoryTree(egenv.CacheDirectory("dev.native.libs")),
			console.GenerateFlutter,
			egbug.DebugFailure(
				console.CompileDarwinBinding,
				egbug.Log("flutter failed to build binding"),
			),
			egbug.DebugFailure(
				console.BuildDarwin,
				eg.Sequential(
					egbug.Log("flutter failed to build app debug initiated"),
					egbug.DirectoryTree(egenv.CacheDirectory("dev.native.libs")),
					// tarballapp doesn't exist yet at this point (created later by the
					// mkdir+cp step below); the flutter build's own output directory is
					// what's actually populated here.
					egbug.DirectoryTree(egenv.WorkingDirectory("console", "build", "macos", "Build", "Products", "Release")),
					// architectures of every native lib in play, to catch lipo/universal-binary
					// mismatches (e.g. a single-arch dylib handed to lipo for multiple slots).
					shell.Op(
						shell.Newf(
							"find %s %s -name '*.dylib' -o -name '*.a' | xargs -I{} sh -c 'echo {}: $(lipo -info {} 2>&1)'",
							egenv.CacheDirectory("dev.native.libs"),
							egenv.WorkingDirectory("console", "build", "macos"),
						).Lenient(true),
						shell.Newf(
							"xcodebuild -workspace %s -scheme Runner -configuration Release -showBuildSettings | grep -E 'ARCHS|EXCLUDED_ARCHS'",
							egenv.WorkingDirectory("console", "macos", "Runner.xcworkspace"),
						).Lenient(true),
					),
					egbug.DirectoryTree(tarballapp),
					egbug.Log("flutter failed to build app debug completed"),
				),
			),
			egbug.DirectoryTree(egenv.WorkingDirectory("console", "build", "macos", "Build", "Products", "Release", "retrovibed.app")),
			// Copy the built .app bundle to the staging workspace directory
			shell.Op(
				shell.Newf("mkdir -p %s", tarballapp),
				shell.Newf(
					"cp -R %s/ %s/",
					egenv.WorkingDirectory("console", "build", "macos", "Build", "Products", "Release", "retrovibed.app"),
					tarballapp,
				),
			),
			egbug.DirectoryTree(tarballapp),
			release.KeychainPEM(
				egenv.String("", "APPLE_SIGNING_KEY"),
				egenv.String("", "APPLE_SIGNING_CER"),
			),
			release.KeychainAppendPEM(
				"installer",
				egenv.String("", "APPLE_MACOS_INSTALLER_KEY"),
				egenv.String("", "APPLE_MACOS_INSTALLER_CERT"),
			),
			release.KeychainAppendPEM(
				"appstore",
				egenv.String("", "APPLE_MACOS_APPSTORE_KEY"),
				egenv.String("", "APPLE_MACOS_APPSTORE_CERT"),
			),
			release.AuthKey(
				apikey,
				egenv.String("", "RETROVIBED_APPLE_AUTH_KEY"),
			),
			shell.Op(
				shell.Newf("security unlock-keychain -p %s %s", egenv.RunID(), keychainPath),
				shell.Newf("codesign --deep --force --options runtime --sign \"Developer ID Application\" --keychain %s %s", keychainPath, tarballapp),
			),
			release.DarwinDmg(tarinfo()),
			shell.Op(
				shell.Newf("security unlock-keychain -p %s %s", egenv.RunID(), keychainPath),
				shell.Newf("codesign --force --sign \"Developer ID Application\" --keychain %s %s", keychainPath, dmgpath),
				shell.Newf("xcrun notarytool submit %s --key ~/.private_keys/AuthKey_${APPLE_API_KEY}.p8 --key-id ${APPLE_API_KEY} --issuer ${APPLE_ISSUER_ID} --wait", dmgpath).
					Environ("APPLE_API_KEY", apikey).
					Environ("APPLE_ISSUER_ID", issuerid),
				shell.Newf("xcrun stapler staple %s", dmgpath),
			),
			shell.Op(
				shell.Newf("rm -rf %s && cp -R %s %s", appstoreapp, tarballapp, appstoreapp),
			),
			release.EmbedProvisioningProfile(
				egenv.String("", "APPLE_MACOS_APPSTORE_PROFILE"),
				filepath.Join(appstoreapp, "Contents", "embedded.provisionprofile"),
			),
			shell.Op(
				shell.Newf("rm -rf %s/Contents/Helpers %s/Contents/Frameworks/retrovibed.h", appstoreapp, appstoreapp),
				shell.Newf("chmod -R a+rX %s", appstoreapp),
				shell.Newf("security unlock-keychain -p %s %s", egenv.RunID(), keychainPath),
				shell.Newf("find %s/Contents/Frameworks -depth -name '*.framework' -o -name '*.dylib' | xargs -I{} codesign --force --options runtime --sign \"Apple Distribution\" --keychain %s {}", appstoreapp, keychainPath),
				shell.Newf("codesign --force --options runtime --sign \"Apple Distribution\" --entitlements %s --keychain %s %s", entitlements, keychainPath, appstoreapp),
				shell.Newf("productbuild --component %s /Applications --sign \"3rd Party Mac Developer Installer\" --keychain %s %s", appstoreapp, keychainPath, pkgpath),
				shell.Newf("xcrun altool --upload-app --type macos -f %s --apiKey ${APPLE_API_KEY} --apiIssuer ${APPLE_ISSUER_ID} 2>&1 | tee /tmp/altool.log && ! grep -q 'UPLOAD FAILED' /tmp/altool.log", pkgpath).
					Environ("APPLE_API_KEY", apikey).
					Environ("APPLE_ISSUER_ID", issuerid),
			),
			eggithub.Release(dmgpath),
		),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
