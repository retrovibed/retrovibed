package main

import (
	"context"
	"eg/compute/console"
	"eg/compute/debuild/duckdb"
	"eg/compute/egapplex"
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
			eg.Parallel(
				duckdb.MaybeBuild(
					egenv.CacheDirectory("dev.native.libs", "libduckdb.a"),
					duckdb.CompileDarwinRuntime("osx_arm64", "arm64"),
					duckdb.CompileDarwin,
					duckdb.CloneBuild,
				),
				neurals.CompileDarwin(egenv.CacheDirectory("dev.native.libs")),
			),
			console.GenerateFlutter,
			console.CompileDarwinBinding,
			console.BuildDarwin,
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
			egapplex.AuthKey(
				apikey,
				egenv.Base64(nil, "RETROVIBED_APPLE_AUTH_KEY"),
			),
			egapplex.Sign("Developer ID Application", keychainPath, tarballapp, egapplex.SignDeep(), egapplex.SignRuntime()),
			release.DarwinDmg(tarinfo()),
			egapplex.Sign("Developer ID Application", keychainPath, dmgpath),
			egapplex.Notarize(apikey, issuerid, dmgpath),
			shell.Op(
				shell.Newf("rm -rf %s && cp -R %s %s", appstoreapp, tarballapp, appstoreapp),
			),
			egapplex.Provision(
				egenv.Base64(nil, "APPLE_MACOS_APPSTORE_PROFILE"),
				filepath.Join(appstoreapp, "Contents", "embedded.provisionprofile"),
			),
			shell.Op(
				shell.Newf("rm -rf %s/Contents/Helpers %s/Contents/Frameworks/retrovibed.h", appstoreapp, appstoreapp),
				shell.Newf("chmod -R a+rX %s", appstoreapp),
				shell.Newf("security unlock-keychain -p %s %s", egenv.RunID(), keychainPath),
				shell.Newf("find %s/Contents/Frameworks -depth -name '*.framework' -o -name '*.dylib' | xargs -I{} codesign --force --options runtime --sign \"Apple Distribution\" --keychain %s {}", appstoreapp, keychainPath),
			),
			egapplex.Sign("Apple Distribution", keychainPath, appstoreapp, egapplex.SignRuntime(), egapplex.SignEntitlements(entitlements)),
			shell.Op(
				shell.Newf("productbuild --component %s /Applications --sign \"3rd Party Mac Developer Installer\" --keychain %s %s", appstoreapp, keychainPath, pkgpath),
			),
			egapplex.Upload(apikey, issuerid, pkgpath, "macos"),
			eggithub.Release(dmgpath),
		),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
