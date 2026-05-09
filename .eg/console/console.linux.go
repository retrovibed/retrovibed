package console

import (
	"context"
	"eg/compute/flatpakmods"
	"eg/compute/tarballs"
	"fmt"
	"time"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/egflatpak"
	"github.com/egdaemon/eg/runtime/x/wasi/eggithub"
	"github.com/egdaemon/eg/runtime/x/wasi/egtarball"
)

func BuildLinux(ctx context.Context, _ eg.Op) error {
	runtime := flutterRuntimev2(shell.Runtime())
	return shell.Run(
		ctx,
		runtime.New("rm -rf build/linux/x64/debug").Lenient(true),
		runtime.New("mkdir -p build/native_assets/linux"),
		runtime.New("flutter build linux --release lib/main.dart"),
	)
}

func BuildAndroidAPK(runtime shell.Command) eg.OpFn {
	return func(ctx context.Context, _ eg.Op) error {
		runtime = flutterRuntimev2(runtime)

		commit := eggit.EnvCommit()
		return shell.Run(
			ctx,
			runtime.New(
				fmt.Sprintf("flutter build apk --build-name=%s --build-number=%s --release lib/main.dart", tarballs.Version(), commit.StringReplace("%git.commit.unix%")),
			).Timeout(20*time.Minute),
			runtime.New("mv app-release.apk retrovibed.apk").Timeout(20*time.Minute).Directory(egenv.WorkingDirectory("console/build/app/outputs/apk/release")),
		)
	}
}

func BuildAndroidBundle(runtime shell.Command) eg.OpFn {
	return func(ctx context.Context, _ eg.Op) error {
		runtime = flutterRuntimev2(runtime)
		commit := eggit.EnvCommit()

		return shell.Run(
			ctx,
			runtime.New(
				fmt.Sprintf("flutter build appbundle --build-name=%s --build-number=%s --release lib/main.dart", tarballs.Version(), commit.StringReplace("%git.commit.unix%")),
			).Timeout(20*time.Minute),
			runtime.New("mv app-release.aab retrovibed.aab").Timeout(20*time.Minute).Directory(egenv.WorkingDirectory("console/build/app/outputs/bundle/release")),
		)
	}
}

func flatpak(final egflatpak.Module) *egflatpak.Builder {
	return egflatpak.New(
		"space.retrovibe.Console", "retrovibe",
		egflatpak.Option().SDK("org.gnome.Sdk", "50").Runtime("org.gnome.Platform", "50").
			Modules(
				flatpakmods.Libduckdb(),
				flatpakmods.Libass(),
				flatpakmods.Libbs2b(),
				flatpakmods.Libplacebo(),
				flatpakmods.Libx264(),
				flatpakmods.Libx265(),
				flatpakmods.Libffmpeg(),
				flatpakmods.Libmpv(),
				final,
			).
			AllowWayland().
			AllowDRI().
			AllowNetwork().
			AllowDownload().
			AllowMusic().
			AllowVideos().Allow(
			// we specify environment variables here so they show up in flatseal for easy adjustments.
			"--socket=pulseaudio",                           // for mpv
			"--filesystem=xdg-run/pipewire-0:ro",            // for mpv
			"--filesystem=~/.duckdb:create",                 // for duckdb
			"--socket=fallback-x11",                         // to appease the flatpak linter for flathub.
			"--share=ipc",                                   // enable standard desktop functionality.
			"--filesystem=xdg-run/gvfsd",                    // enable standard desktop functionality. (probably unnnecessary)
			"--env=LC_NUMERIC=C",                            // for mpv
			"--env=TMPDIR=/var/tmp/",                        // enaure golang sets its os.TempDir() to a working value.
			"--env=RETROVIBED_MDNS_DISABLED=true",           // disable MDNS when running in flatpak since it doesn't work.
			"--env=RETROVIBED_AUTO_IDENTIFY_MEDIA=false",    // automatically identify metadata for content. experimental.
			"--env=RETROVIBED_MEDIA_AUTO_ARCHIVE=true",      // enable content marked for archiving being uploaded into storage.
			"--env=RETROVIBED_MEDIA_AUTO_RECLAIM=false",     // allow automatically reclaiming disk space by removing archived data.
			"--env=RETROVIBED_TORRENT_AUTO_DISCOVERY=false", // peer scanning is an experimental feature.
			"--env=RETROVIBED_TORRENT_AUTO_BOOTSTRAP=true",  // auto bootstrap the dht from a global endpoint.
			"--env=RETROVIBED_TORRENT_ALLOW_SEEDING=true",   // allow uploading to peers
			"--env=RETROVIBED_TORRENT_PUBLIC_IP4=\"\"",      // manually set the public ipv4 address.
			"--env=RETROVIBED_TORRENT_PUBLIC_IP6=\"\"",      // manually set the public ipv6 address.
			"--env=RETROVIBED_JWT_SECRET=",                  // specify the jwt secret to use for signing tokens. generally this should not be necessary.
			"--env=RETROVIBED_SELF_SIGNED_HOSTS=127.0.0.1",  // TLS hosts to include in the self signed certificate.
		)...)
}

// // build ensures that the flatpak has all the necessary componentry for the generated manifest.
// func FlatpakBuild(ctx context.Context, op eg.Op) error {
// 	return egflatpak.Build(ctx, shell.Runtime().Timeout(30*time.Minute), flatpak(
// 		egflatpak.ModuleTarball(
// 			eggithub.DownloadURL(tarballs.Retrovibed()),
// 			egtarball.SHA256(tarballs.Retrovibed()),
// 		),
// 	))
// }

// Manifest generates the manifest for distribution.
func FlatpakManifest(b *tarballs.Build) eg.OpFn {
	return func(ctx context.Context, o eg.Op) error {
		// needs to be in an op so that the sha256 and download url are defered until this command actually runs.
		return eg.Sequential(
			egflatpak.ManifestOp(
				egenv.CacheDirectory(tarballs.Flatpak(b)),
				flatpak(
					moduleTarball(eggithub.DownloadURL(tarballs.Retrovibed(b)), egtarball.SHA256(tarballs.Retrovibed(b))),
				),
			),
		)(ctx, o)
	}
}

func moduleTarball(url, sha256d string) egflatpak.Module {
	return egflatpak.NewModule("retrovibed", "simple", egflatpak.ModuleOptions().Commands(
		"mv usr/share/applications/space.retrovibe.Console.desktop /app/share/applications/space.retrovibe.Console.desktop",
		"mv usr/share/icons/hicolor/scalable/apps/space.retrovibe.Console.svg /app/share/icons/hicolor/scalable/apps/space.retrovibe.Console.svg",
		"mv usr/share/metainfo/space.retrovibe.Console.metainfo.xml /app/share/metainfo/space.retrovibe.Console.metainfo.xml",
		"mv usr/share/licenses/space.retrovibe.Console /app/share/licenses/space.retrovibe.Console",
		"mv usr/lib/retrovibed/lib /app/lib",
		"mv usr/lib/retrovibed/* /app/bin",
		"rm -rf usr",
		"rm -rf etc",
	).Sources(egflatpak.SourceTarball(url, sha256d))...)
}
