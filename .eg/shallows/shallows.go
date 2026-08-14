package shallows

import (
	"context"
	"eg/compute/neurals"
	"eg/compute/tarballs"
	"path/filepath"
	"strings"
	"time"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/eggolang"
	"github.com/egdaemon/eg/runtime/x/wasi/egtarball"
)

func testBuildTags() []string {
	return []string{"duckdb_use_lib"}
}

var buildTags = []string{"duckdb_use_lib", "retrovibed", "neural"}

func rootdir() string {
	return egenv.WorkingDirectory("shallows")
}

func shellruntime() shell.Command {
	return eggolang.Runtime().Directory(rootdir()).Environ(
		"CACHE_DIRECTORY", egenv.CacheDirectory(),
	).Environ(
		"XDG_CACHE_HOME", egenv.CacheDirectory("xdg"), // temporary until eg catches up.
	).Environ(
		"GOLANGCI_LINT_CACHE", egenv.CacheDirectory("golang-lint"),
	).Environ(
		"CGO_LDFLAGS", "-L"+egenv.CacheDirectory("neurals"),
	).Environ(
		"LD_LIBRARY_PATH", egenv.CacheDirectory("neurals"),
	)
}

func NeuralsBuild() eg.OpFn {
	return neurals.MaybeBuild(
		"neurals/libpredicttext.so",
		neurals.Compile,
		neurals.Clone,
	)
}

func Generate(ctx context.Context, op eg.Op) error {
	return eg.Sequential(
		GenerateProtocol,
		GenerateGogen,
	)(ctx, op)
}

func GenerateGogen(ctx context.Context, _ eg.Op) error {
	gruntime := shellruntime()
	return shell.Run(
		ctx,
		gruntime.New("go generate ./... && go fmt ./...").Timeout(30*time.Minute),
	)
}

func GenerateProtocol(ctx context.Context, op eg.Op) error {
	gruntime := shellruntime()
	return shell.Run(
		ctx,
		// Every invocation below uses paths=import with an explicit module= (rather than
		// paths=source_relative): output location is then computed purely from the Go
		// package's import path (the --go_opt=M value), completely independent of the
		// .proto file's own directory. That's what lets .proto/ be organized into
		// per-package subdirectories (needed so protoc-gen-dart can resolve cross-package
		// Dart imports, which has no M-option equivalent) without perturbing where any
		// generated .pb.go file lands, and without the nested-vs-flat descriptor-name
		// mismatches that paths=source_relative + a bucket-local proto_path produced for
		// files that share a Go package with a sibling that imports them (protoc-gen-go
		// chains sibling file-init calls by descriptor name, so those need to agree).
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mmeta/meta.search.proto=github.com/retrovibed/retrovibed/shallows/meta --go_opt=paths=import --go_opt=module=github.com/retrovibed/retrovibed/shallows --go_out=. meta/meta.search.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mmeta/meta.account.proto=github.com/retrovibed/retrovibed/retroapi/authn --go_opt=paths=import --go_opt=module=github.com/retrovibed/retrovibed/retroapi --go_out=../retroapi meta/meta.account.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mmeta/meta.profile.proto=github.com/retrovibed/retrovibed/retroapi/authn --go_opt=paths=import --go_opt=module=github.com/retrovibed/retrovibed/retroapi --go_out=../retroapi meta/meta.profile.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mmeta/meta.authz.proto=github.com/retrovibed/retrovibed/retroapi/authn --go_opt=paths=import --go_opt=module=github.com/retrovibed/retrovibed/retroapi --go_out=../retroapi meta/meta.authz.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mmeta/meta.daemon.proto=github.com/retrovibed/retrovibed/shallows/metaapi --go_opt=paths=import --go_opt=module=github.com/retrovibed/retrovibed/shallows --go_out=. meta/meta.daemon.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mwireguard/meta.wireguard.proto=github.com/retrovibed/retrovibed/shallows/metaapi --go_opt=paths=import --go_opt=module=github.com/retrovibed/retrovibed/shallows --go_out=. wireguard/meta.wireguard.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Maudio/meta.audio.proto=github.com/retrovibed/retrovibed/shallows/metaapi --go_opt=paths=import --go_opt=module=github.com/retrovibed/retrovibed/shallows --go_out=. audio/meta.audio.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mmeta/meta.network.proto=github.com/retrovibed/retrovibed/shallows/metaapi --go_opt=paths=import --go_opt=module=github.com/retrovibed/retrovibed/shallows --go_out=. meta/meta.network.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mmeta/meta.torrent.proto=github.com/retrovibed/retrovibed/shallows/metaapi --go_opt=paths=import --go_opt=module=github.com/retrovibed/retrovibed/shallows --go_out=. meta/meta.torrent.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mmeta/meta.dht.proto=github.com/retrovibed/retrovibed/shallows/metaapi --go_opt=paths=import --go_opt=module=github.com/retrovibed/retrovibed/shallows --go_out=. meta/meta.dht.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mmeta/meta.discovery.proto=github.com/retrovibed/retrovibed/shallows/metaapi --go_opt=paths=import --go_opt=module=github.com/retrovibed/retrovibed/shallows --go_out=. meta/meta.discovery.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mmeta/meta.authn.proto=github.com/retrovibed/retrovibed/retroapi/authn --go_opt=Mmeta/meta.account.proto=github.com/retrovibed/retrovibed/retroapi/authn --go_opt=Mmeta/meta.profile.proto=github.com/retrovibed/retrovibed/retroapi/authn --go_opt=paths=import --go_opt=module=github.com/retrovibed/retrovibed/retroapi --go_out=../retroapi meta/meta.authn.proto"),
		// media
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mmedia/media.proto=github.com/retrovibed/retrovibed/shallows/media --go_opt=paths=import --go_opt=module=github.com/retrovibed/retrovibed/shallows --go_out=. media/media.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mmedia/media.known.proto=github.com/retrovibed/retrovibed/shallows/media --go_opt=Mmeta/meta.search.proto=github.com/retrovibed/retrovibed/shallows/meta --go_opt=paths=import --go_opt=module=github.com/retrovibed/retrovibed/shallows --go_out=. media/media.known.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mmedia/media.recent.proto=github.com/retrovibed/retrovibed/shallows/media --go_opt=Mmeta/meta.search.proto=github.com/retrovibed/retrovibed/shallows/meta --go_opt=Mmedia/media.proto=github.com/retrovibed/retrovibed/shallows/media --go_opt=paths=import --go_opt=module=github.com/retrovibed/retrovibed/shallows --go_out=. media/media.recent.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mrss/rss.proto=github.com/retrovibed/retrovibed/shallows/rss --go_opt=paths=import --go_opt=module=github.com/retrovibed/retrovibed/shallows --go_out=. rss/rss.proto"),
		// remote control
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mmedia/media.remote.control.proto=github.com/retrovibed/retrovibed/shallows/mediaapi --go_opt=Mmedia/media.proto=github.com/retrovibed/retrovibed/shallows/media --go_opt=Mmeta/meta.daemon.proto=github.com/retrovibed/retrovibed/shallows/metaapi --go_opt=paths=import --go_opt=module=github.com/retrovibed/retrovibed/shallows --go_out=. media/media.remote.control.proto"),
		// block cache
		// content.addressable.storage.proto's M value below doesn't actually match the
		// retroapi module's real path (missing a "retrovibed" segment) - a pre-existing
		// inconsistency, left untouched (paths=import requires the M value to start with
		// module=, so this one keeps the old bucket-local-proto_path/source_relative form).
		gruntime.New("protoc --proto_path=../.proto --proto_path=../.proto/media --go_opt=Mcontent.addressable.storage.proto=github.com/retrovibed/retroapi/deeppool --go_opt=paths=source_relative --go_out=../retroapi/deeppool content.addressable.storage.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mmedia/media.id.proto=github.com/retrovibed/retrovibed/retroapi/deeppool --go_opt=paths=import --go_opt=module=github.com/retrovibed/retrovibed/retroapi --go_out=../retroapi media/media.id.proto"),
		// community
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mcommunity/community.proto=github.com/retrovibed/retrovibed/shallows/communityapi --go_opt=paths=import --go_opt=module=github.com/retrovibed/retrovibed/shallows --go_out=. community/community.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mcommunity/community.metrics.proto=github.com/retrovibed/retrovibed/shallows/communityapi --go_opt=paths=import --go_opt=module=github.com/retrovibed/retrovibed/shallows --go_out=. community/community.metrics.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mcommunity/community.publish.proto=github.com/retrovibed/retrovibed/shallows/communityapi --go_opt=Mcommunity/community.proto=github.com/retrovibed/retrovibed/shallows/communityapi --go_opt=paths=import --go_opt=module=github.com/retrovibed/retrovibed/shallows --go_out=. community/community.publish.proto"),
		// ddisc
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mddisc/ddisc.peers.proto=github.com/retrovibed/retrovibed/shallows/ddiscapi --go_opt=paths=import --go_opt=module=github.com/retrovibed/retrovibed/shallows --go_out=. ddisc/ddisc.peers.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mddisc/ddisc.discovery.proto=github.com/retrovibed/retrovibed/shallows/ddiscapi --go_opt=Mmeta/meta.search.proto=github.com/retrovibed/retrovibed/shallows/meta --go_opt=paths=import --go_opt=module=github.com/retrovibed/retrovibed/shallows --go_out=. ddisc/ddisc.discovery.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mddisc/ddisc.media.proto=github.com/retrovibed/retrovibed/shallows/ddiscapi --go_opt=Mmeta/meta.search.proto=github.com/retrovibed/retrovibed/shallows/meta --go_opt=paths=import --go_opt=module=github.com/retrovibed/retrovibed/shallows --go_out=. ddisc/ddisc.media.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mmedia/ddisc.locate.proto=github.com/retrovibed/retrovibed/shallows/ddiscapi --go_opt=paths=import --go_opt=module=github.com/retrovibed/retrovibed/shallows --go_out=. media/ddisc.locate.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mddisc/plugin/searchplugin.proto=github.com/retrovibed/retrovibed/shallows/ddiscapi --go_opt=paths=import --go_opt=module=github.com/retrovibed/retrovibed/shallows --go_out=. ddisc/plugin/searchplugin.proto"),
		// ddisc.import.proto is generated only into retroapi/ddiscapi (not shallows/ddiscapi) —
		// shallows code imports retroapi/ddiscapi directly for the Import type instead of
		// duplicating its own copy, since two packages registering the same proto file into
		// the global protobuf registry in the same binary panics at init.
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mddisc/ddisc.import.proto=github.com/retrovibed/retrovibed/retroapi/ddiscapi --go_opt=paths=import --go_opt=module=github.com/retrovibed/retrovibed/retroapi --go_out=../retroapi ddisc/ddisc.import.proto"),
		// settings
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mstorage/storage.proto=github.com/retrovibed/retrovibed/shallows/cmd/retrovibe/daemons --go_opt=paths=import --go_opt=module=github.com/retrovibed/retrovibed/shallows --go_out=. storage/storage.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mtorrents/torrent.proto=github.com/retrovibed/retrovibed/shallows/cmd/retrovibe/daemons --go_opt=paths=import --go_opt=module=github.com/retrovibed/retrovibed/shallows --go_out=. torrents/torrent.proto"),
	)
}

func Install(b *tarballs.Build) eg.OpFn {
	dstdir := filepath.Join(egtarball.Path(tarballs.Retrovibed(b)), "usr", "lib", "retrovibed")
	gruntime := shellruntime()

	return shell.Op(
		gruntime.Newf("mkdir -p ~/.duckdb %s && ln -sfn %s ~/.duckdb", egenv.CacheDirectory("duckdb"), egenv.CacheDirectory("duckdb")),
		gruntime.Newf(
			"go build -buildvcs=true -tags %s -o %s ./cmd/...",
			strings.Join(buildTags, ","),
			dstdir,
		),
		gruntime.Newf("cp libpredicttext.so %s/", dstdir).Directory(egenv.CacheDirectory("neurals")),
	)
}

func Compile() eg.OpFn {
	return eggolang.AutoCompile(
		eggolang.CompileOption.BuildOptions(
			eggolang.Build(
				eggolang.BuildOption.Tags(buildTags...),
				eggolang.BuildOption.WorkingDirectory(rootdir()),
				eggolang.BuildOption.Environ(
					"CGO_LDFLAGS=-L"+egenv.CacheDirectory("neurals"),
					"LD_LIBRARY_PATH="+egenv.CacheDirectory("neurals"),
				),
			),
		),
	)
}

func Test() eg.OpFn {
	return eg.Sequential(
		eggolang.AutoTest(
			eggolang.TestOption.BuildOptions(
				eggolang.Build(
					eggolang.BuildOption.Tags(testBuildTags()...),
					eggolang.BuildOption.WorkingDirectory(rootdir()),
				),
			),
		),
		eggolang.RecordCoverage,
	)
}

func Linting(ctx context.Context, _ eg.Op) error {
	gruntime := shellruntime()
	return shell.Run(
		ctx,
		gruntime.New("golangci-lint run ./..."),
	)
}

// func FlatpakManifest(ctx context.Context, o eg.Op) error {
// 	b := egflatpak.New(
// 		"space.retrovibe.Daemon", "retrovibed",
// 		egflatpak.Option().SDK("org.gnome.Sdk", "47").Runtime("org.gnome.Platform", "47").
// 			Modules(
// 				flatpakmods.Libduckdb(),
// 				egflatpak.NewModule("retrovibed", "simple", egflatpak.ModuleOptions().Commands(
// 					"cp -r . /app/bin",
// 				).Sources(
// 					egflatpak.SourceTarball(
// 						eggithub.DownloadURL(tarballs.Retrovibed()), egtarball.SHA256(tarballs.Retrovibed()),
// 						egflatpak.SourceOptions().Destination("retrovibed.tar.xz")...,
// 					),
// 				)...),
// 			).
// 			AllowWayland().
// 			AllowDRI().
// 			AllowNetwork().
// 			AllowDownload().
// 			AllowMusic().
// 			AllowVideos().Allow(
// 			"--filesystem=~/Downloads:ro",  // bug in flatpak doesn't properly grant access to xdg-download
// 			"--filesystem=~/Videos:create", // bug in flatpak doesn't properly grant full access to videos directory
// 			"--filesystem=~/Music:create",  // bug in flatpak doesn't properly grant full access to music directory
// 		)...)

// 	return egflatpak.ManifestOp(egenv.CacheDirectory("flatpak.daemon.yml"), b)(ctx, o)
// }
