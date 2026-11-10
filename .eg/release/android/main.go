package main

import (
	"context"
	"eg/compute/android"
	"eg/compute/console"
	"eg/compute/debuild/duckdb"
	"eg/compute/eggradlex"
	"eg/compute/maintainer"
	"eg/compute/neurals"
	"eg/compute/release"
	"log"
	"time"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/egbug"
	"github.com/egdaemon/eg/runtime/x/wasi/egfs"
	"github.com/egdaemon/eg/runtime/x/wasi/eggithub"
	"github.com/egdaemon/eg/runtime/x/wasi/egsecrets"
)

// androidruntime does not set the native-libs directory itself —
// console.BuildAndroidAPK/BuildAndroidBundle apply that after
// flutterRuntimev2, since flutterRuntimev2 unconditionally overrides it
// with its own dev.native.libs default and would clobber it if set here
// first.
func androidruntime() shell.Command {
	ctx, done := context.WithTimeout(context.Background(), time.Minute)
	defer done()

	return shell.Env().
		EnvironFrom(eggradlex.Env()...).
		Environ("ANDROID_HOME", "/opt/android-sdk").
		EnvironFrom(egsecrets.Env(ctx, "gcpsm://retrovibed-prod/android-keystore-prod-env")...).
		Environ("RETROVIBED_ANDROID_KEY_STORE_PATH", egenv.CacheDirectory("android", "keystore"))
}

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	ctx, done := context.WithTimeout(context.Background(), egenv.TTL())
	defer done()
	deb := eg.Container(android.Container)
	err := eg.Perform(
		ctx,
		eggit.AutoClone,
		eg.Build(eg.Container(maintainer.Container).BuildFromFile(".eg/Containerfile")),
		eg.Build(deb.BuildFromFile(".eg/release/android/Containerfile")),
		eg.Module(
			ctx,
			deb,
			eg.Sequential(
				eg.WhenFn(egfs.FileNotExistsFn(egenv.CacheDirectory("android", "keystore")), egsecrets.CopyIntoFileOp(egenv.CacheDirectory("android", "keystore"), "gcpsm://retrovibed-prod/android-keystore/latest")),
				console.Generate,
				egbug.Log("generated console bindings"),
				eg.Parallel(
					duckdb.MaybeBuild(
						android.JNIStagingDir("x64", "libduckdb.a"),
						duckdb.CompileAndroidRuntime("android_x86_64", "x86_64"),
						duckdb.CompileAndroid,
						duckdb.CloneStaticBuild,
					),
					duckdb.MaybeBuild(
						android.JNIStagingDir("arm64", "libduckdb.a"),
						duckdb.CompileAndroidRuntime("android_arm64", "arm64-v8a"),
						duckdb.CompileAndroid,
						duckdb.CloneStaticBuild,
					),
					neurals.CompileAndroid("x86_64", android.JNIStagingDir("x64")),
					neurals.CompileAndroid("arm64-v8a", android.JNIStagingDir("arm64")),
				),
				egbug.Log("generated static libraries for android"),
				eg.Parallel(
					console.GenerateStaticBinding(
						android.JNILibDir("x64"),
						console.AndroidRuntime("x86_64-none-linux-android31").
							Environ("GOARCH", "amd64"),
						android.JNIStagingDir("x64"),
					),
					console.GenerateStaticBinding(
						android.JNILibDir("arm64"),
						console.AndroidRuntime("aarch64-none-linux-android31").
							Environ("GOARCH", "arm64"),
						android.JNIStagingDir("arm64"),
					),
				),
				egbug.Log("generated static libraries for duckdb"),
				console.BuildAndroidAPK(androidruntime()),
				console.BuildAndroidBundle(androidruntime()),
				release.Dist(
					egenv.WorkingDirectory("console/build/app/outputs/apk/release/retrovibed.apk"),
					egenv.WorkingDirectory("console/build/app/outputs/bundle/release/retrovibed.aab"),
				),
				eggithub.Draft(
					egenv.WorkingDirectory("console/build/app/outputs/apk/release/retrovibed.apk"),
					egenv.WorkingDirectory("console/build/app/outputs/bundle/release/retrovibed.aab"),
				),
				android.Draft(android.ID, egenv.WorkingDirectory("console/build/app/outputs/bundle/release/retrovibed.aab"), "internal"),
			),
		),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
