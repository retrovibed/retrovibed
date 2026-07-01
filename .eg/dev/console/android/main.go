package main

import (
	"context"
	"eg/compute/android"
	"eg/compute/console"
	"eg/compute/debuild/duckdb"
	"eg/compute/egx"
	"eg/compute/maintainer"
	"fmt"
	"log"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
	"github.com/egdaemon/eg/runtime/wasi/shell"
)

func Device(ctx context.Context, o eg.Op) error {
	return eg.Sequential(
		shell.Op(shell.New("adb -d get-serialno")),
		shell.Op(shell.New("adb -d wait-for-device shell 'while [ \"$(getprop sys.boot_completed)\" != \"1\" ]; do sleep 1; done'")),
		console.RunDev("flutter pub get && flutter run -d $(adb -d get-serialno)"),
	)(ctx, o)
}

func Emulator(port int) eg.OpFn {
	return eg.Sequential(
		shell.Op(
			shell.New("sdkmanager 'system-images;android-34;google_apis;x86_64'"),
			shell.New("echo 'no' | avdmanager create avd --force -n retrovibed -k 'system-images;android-34;google_apis;x86_64'"),
			shell.Newf("systemctl --user is-active --quiet retrovibed-android || systemd-run --user --unit=retrovibed-android --setenv=QT_QPA_PLATFORM=xcb --collect /opt/android-sdk/emulator/emulator -avd retrovibed -port %d", port),
			shell.Newf("adb -s emulator-%d wait-for-device shell 'while [ \"$(getprop sys.boot_completed)\" != \"1\" ]; do sleep 1; done'", port),
		),
		console.RunDev(fmt.Sprintf("flutter pub get && flutter run -d emulator-%d", port)),
	)
}

func Debug(runtime shell.Command) eg.OpFn {
	return shell.Op(
		runtime.New("which go && go version"),
		runtime.New("which flutter && flutter --version"),
		runtime.New("which dart && dart --version"),
	)
}

const (
	port = 5554
)

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
				duckdb.Download,
				eg.Parallel(
					duckdb.MaybeBuild(
						egenv.WorkingDirectory("console/android/app/src/main/jniLibs/x86_64/libduckdb.a"),
						duckdb.CompileAndroidRuntime("android_x86_64", "x86_64"),
						duckdb.CompileAndroid,
						duckdb.CloneBuild,
					),
					duckdb.MaybeBuild(
						egenv.WorkingDirectory("console/android/app/src/main/jniLibs/arm64-v8a/libduckdb.a"),
						duckdb.CompileAndroidRuntime("android_arm64", "arm64-v8a"),
						duckdb.CompileAndroid,
						duckdb.CloneBuild,
					),
				),
				console.Generate,
				eg.Parallel(
					console.GenerateDevStaticBinding(
						console.AndroidRuntime("x86_64-none-linux-android31").Environ("GOARCH", "amd64"),
						egenv.WorkingDirectory("console/android/app/src/main/jniLibs/x86_64"),
						egenv.WorkingDirectory("console/android/app/src/main/jniLibs/x86_64"),
					),
					console.GenerateDevStaticBinding(
						console.AndroidRuntime("aarch64-none-linux-android31").Environ("GOARCH", "arm64"),
						egenv.WorkingDirectory("console/android/app/src/main/jniLibs/arm64-v8a"),
						egenv.WorkingDirectory("console/android/app/src/main/jniLibs/arm64-v8a"),
					),
				),
				console.BuildLinux,
			),
		),
		egx.Fallback(
			Device,
			Emulator(port),
		),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
