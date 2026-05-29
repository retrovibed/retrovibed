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
	"path/filepath"
	"time"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/egbug"
	"github.com/egdaemon/eg/runtime/x/wasi/eggithub"
	"github.com/egdaemon/eg/runtime/x/wasi/eggolang"
	"github.com/egdaemon/eg/runtime/x/wasi/egtarball"
)

func tarinfo() *tarballs.Build {
	return &tarballs.Build{
		OS:   egenv.String("darin", "EG_COMPUTE_HOST_OS"),
		Arch: egenv.String("arm64", "EG_COMPUTE_HOST_ARCH"),
	}
}

// Baremetal command for darwin due to macosx nonsense for no cloud vms.
func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	ctx, done := context.WithTimeout(context.Background(), egenv.TTL())
	defer done()

	runtime := shell.Runtime().
		EnvironFrom(eggolang.Env()...).
		Environ("PUB_CACHE", egenv.CacheDirectory(".eg", "dart"))

	tarballapp := filepath.Join(egtarball.Path(tarballs.Retrovibed(tarinfo())), "retrovibed.app")
	appstoreapp := egenv.CacheDirectory("retrovibed.appstore.app")
	dmgpath := egenv.CacheDirectory("retrovibed.darwin.arm64.dmg")
	pkgpath := egenv.CacheDirectory("retrovibed.darwin.arm64.pkg")
	entitlements := egenv.WorkingDirectory("console", "macos", "Runner", "Release.entitlements")
	keychainPath := egenv.WorkspaceDirectory("apple.signing.keychain")
	flutter := runtime.Directory(egenv.WorkingDirectory("console")).Debug()
	shallows := runtime.Directory(egenv.WorkingDirectory("shallows"))
	commit := eggit.EnvCommit()
	duckdblibs := egenv.CacheDirectory("duckdb", ".darwin-arm64")
	neuralsdir := egenv.CacheDirectory("neurals")

	duckdbldflags := "-L" + duckdblibs + " " +
		"-Wl,-force_load," + duckdblibs + "/libduckdb.a " +
		"-lc++"

	apikey := egenv.String("", "RETROVIBED_APPLE_API_KEY")
	issuerid := egenv.String("", "RETROVIBED_APPLE_ISSUER_ID")

	err := eg.Perform(
		ctx,
		eg.Sequential(
			eg.Parallel(
				duckdb.MaybeBuild(
					filepath.Join(duckdblibs, "libduckdb.a"),
					duckdb.CompileDarwinRuntime("osx_arm64", "arm64"),
					duckdb.CompileDarwin,
					duckdb.CloneStaticBuild,
				),
				neurals.CompileDarwin(neuralsdir),
			),
			console.GenerateFlutter,
			egbug.DebugFailure(
				shell.Op(
					flutter.New("rm -rf build/macos/{x64,arm64}/debug").Lenient(true),
					flutter.New(fmt.Sprintf("flutter build macos --build-name=%s --build-number=%s --release lib/main.dart", tarballs.Version(), commit.StringReplace("%git.commit.unix%"))).Timeout(10*time.Minute),
				),
				shell.Op(shell.New("flutter failed to build app")),
			),
			shell.Op(
				flutter.Newf("CGO_LDFLAGS=\"%s\" go -C retrovibedbind build --tags duckdb_use_static_lib -buildmode=c-shared -o ../build/macos/Build/Products/Release/retrovibed.app/Contents/Frameworks/retrovibed.dylib ./...", duckdbldflags),
				flutter.New("tree build/macos/Build/Products/Release/retrovibed.app"),
				shell.Newf("mkdir -p %s", tarballapp),
				flutter.Newf("cp -R build/macos/Build/Products/Release/retrovibed.app/ %s/", tarballapp),
			),
			shell.Op(
				shallows.Newf(
					"CGO_LDFLAGS=\"%s -L%s -lpredicttext -Wl,-rpath,@executable_path/../Frameworks\" go install --tags duckdb_use_static_lib,retrovibed,neural ./cmd/...",
					duckdbldflags, neuralsdir,
				).Environ("GOBIN", filepath.Join(tarballapp, "Contents", "Helpers")),
				shell.Newf("cp %s/libpredicttext.dylib %s/Contents/Frameworks/", neuralsdir, tarballapp),
			),
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
