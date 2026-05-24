package main

import (
	"context"
	"eg/compute/android"
	"eg/compute/console"
	"eg/compute/debuild/duckdb"
	"eg/compute/maintainer"
	"eg/compute/neurals"
	"log"
	"time"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/egbug"
	"github.com/egdaemon/eg/runtime/x/wasi/egfs"
	"github.com/egdaemon/eg/runtime/x/wasi/eggithub"
	"github.com/egdaemon/eg/runtime/x/wasi/egsecrets"
)

func androidruntime() shell.Command {
	ctx, done := context.WithTimeout(context.Background(), time.Minute)
	defer done()
	return shell.Env().
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
		eg.Sequential(
			eg.Build(eg.Container(maintainer.Container).BuildFromFile(".eg/Containerfile")),
			eg.Build(deb.BuildFromFile(".eg/release/android/Containerfile")),
			eg.Module(
				ctx,
				deb,
				eg.Sequential(
					// eggcp.CredentialsHack,
					CredentialsHack,
					// eg.WhenFn(egfs.FileNotExistsFn(egenv.CacheDirectory("android", "keystore")), android.SigningKey(egenv.CacheDirectory("android", "keystore"), "upload")),
					eg.WhenFn(egfs.FileNotExistsFn(egenv.CacheDirectory("android", "keystore")), egsecrets.CopyIntoFileOp(egenv.CacheDirectory("android", "keystore"), "gcpsm://retrovibed-prod/android-keystore/latest")),
					console.Generate,
					egbug.Log("generated console bindings"),
					eg.Parallel(
						console.GenerateStaticBinding(
							egenv.WorkingDirectory("console/android/app/src/main/jniLibs/x86_64"),
							console.AndroidRuntime("x86_64-none-linux-android31").
								Environ("GOARCH", "amd64"),
						),
						console.GenerateStaticBinding(
							egenv.WorkingDirectory("console/android/app/src/main/jniLibs/arm64-v8a"),
							console.AndroidRuntime("aarch64-none-linux-android31").
								Environ("GOARCH", "arm64"),
						),
					),
					egbug.Log("generated static libraries for android"),
					eg.Parallel(
						duckdb.MaybeBuild(
							"console/android/app/src/main/jniLibs/x86_64/libduckdb_static.a",
							duckdb.CompileAndroid("android_x86_64", "x86_64"),
							duckdb.CloneAndroid,
						),
						duckdb.MaybeBuild(
							"console/android/app/src/main/jniLibs/arm64-v8a/libduckdb_static.a",
							duckdb.CompileAndroid("android_arm64", "arm64-v8a"),
							duckdb.CloneAndroid,
						),
						neurals.CompileAndroid("x86_64-linux-android", egenv.WorkingDirectory("console/android/app/src/main/jniLibs/x86_64")),
						neurals.CompileAndroid("aarch64-linux-android", egenv.WorkingDirectory("console/android/app/src/main/jniLibs/arm64-v8a")),
					),
					egbug.Log("generated static libraries for duckdb"),
					console.BuildAndroidAPK(androidruntime()),
					console.BuildAndroidBundle(androidruntime()),
					eggithub.Release(
						egenv.WorkingDirectory("console/build/app/outputs/apk/release/retrovibed.apk"),
						egenv.WorkingDirectory("console/build/app/outputs/bundle/release/retrovibed.aab"),
					),
					android.Draft(android.ID, egenv.WorkingDirectory("console/build/app/outputs/bundle/release/retrovibed.aab"), "internal"),
				),
			),
		),
	)

	if err != nil {
		log.Fatalln(err)
	}
}

// temporary hack until remap directory functionality is working.
func CredentialsHack(ctx context.Context, o eg.Op) error {
	// encoded := os.Getenv(_eg.EnvUnsafeGcloudADCB64)
	// if encoded == "" {
	// 	return nil
	// }

	// decoded, err := base64.URLEncoding.DecodeString(encoded)
	// if err != nil {
	// 	return err
	// }

	// if err = os.WriteFile(egenv.WorkloadDirectory("gcloud", "application_default_credentials.json"), decoded, 0600); err != nil {
	// 	return err
	// }

	return shell.Run(
		ctx,
		// shell.Newf("echo %s | tr -- '-_' '+/' | base64 -d -i | install -D -m 600 /dev/stdin ~/.config/gcloud/application_default_credentials.json", encoded),
		// shell.Newf("echo %s | tr -- '-_' '+/' | base64 -d -i | install -D -m 600 /dev/stdin %s", encoded, egenv.WorkloadDirectory("gcloud", "application_default_credentials.json")),
		shell.Newf("which gcloud && gcloud config set auth/credential_file_override %s", egenv.WorkloadDirectory("gcloud", "application_default_credentials.json")),
	)
}
