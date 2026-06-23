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
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mmeta.search.proto=github.com/retrovibed/retrovibed/shallows/meta --go_opt=paths=source_relative --go_out=meta meta.search.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mmeta.account.proto=github.com/retrovibed/retrovibed/retroapi/authn --go_opt=paths=source_relative --go_out=../retroapi/authn meta.account.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mmeta.profile.proto=github.com/retrovibed/retrovibed/retroapi/authn --go_opt=paths=source_relative --go_out=../retroapi/authn meta.profile.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mmeta.authz.proto=github.com/retrovibed/retrovibed/retroapi/authn --go_opt=paths=source_relative --go_out=../retroapi/authn meta.authz.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mmeta.daemon.proto=github.com/retrovibed/retrovibed/shallows/metaapi --go_opt=paths=source_relative --go_out=metaapi meta.daemon.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mmeta.wireguard.proto=github.com/retrovibed/retrovibed/shallows/metaapi --go_opt=paths=source_relative --go_out=metaapi meta.wireguard.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mmeta.network.proto=github.com/retrovibed/retrovibed/shallows/metaapi --go_opt=paths=source_relative --go_out=metaapi meta.network.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mmeta.torrent.proto=github.com/retrovibed/retrovibed/shallows/metaapi --go_opt=paths=source_relative --go_out=metaapi meta.torrent.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mmeta.dht.proto=github.com/retrovibed/retrovibed/shallows/metaapi --go_opt=paths=source_relative --go_out=metaapi meta.dht.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mmeta.discovery.proto=github.com/retrovibed/retrovibed/shallows/metaapi --go_opt=paths=source_relative --go_out=metaapi meta.discovery.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mmeta.authn.proto=github.com/retrovibed/retrovibed/retroapi/authn --go_opt=Mmeta.account.proto=github.com/retrovibed/retrovibed/retroapi/authn --go_opt=Mmeta.profile.proto=github.com/retrovibed/retrovibed/retroapi/authn --go_opt=paths=source_relative --go_out=../retroapi/authn meta.authn.proto"),
		// media
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mmedia.proto=github.com/retrovibed/retrovibed/shallows/media --go_opt=paths=source_relative --go_out=media media.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mmedia.known.proto=github.com/retrovibed/retrovibed/shallows/media --go_opt=Mmeta.search.proto=github.com/retrovibed/retrovibed/shallows/meta --go_opt=paths=source_relative --go_out=media media.known.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mmedia.recent.proto=github.com/retrovibed/retrovibed/shallows/media --go_opt=Mmeta.search.proto=github.com/retrovibed/retrovibed/shallows/meta --go_opt=Mmedia.proto=github.com/retrovibed/retrovibed/shallows/media --go_opt=paths=source_relative --go_out=media media.recent.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mmedia.locate.proto=github.com/retrovibed/retrovibed/shallows/media --go_opt=paths=source_relative --go_out=media media.locate.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mrss.proto=github.com/retrovibed/retrovibed/shallows/rss --go_opt=paths=source_relative --go_out=rss rss.proto"),
		// block cache
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mcontent.addressable.storage.proto=github.com/retrovibed/retroapi/deeppool --go_opt=paths=source_relative --go_out=../retroapi/deeppool content.addressable.storage.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mmedia.id.proto=github.com/retrovibed/retrovibed/retroapi/deeppool --go_opt=paths=source_relative --go_out=../retroapi/deeppool media.id.proto"),
		// community
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mcommunity.proto=github.com/retrovibed/retrovibed/shallows/communityapi --go_opt=paths=source_relative --go_out=communityapi community.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mcommunity.metrics.proto=github.com/retrovibed/retrovibed/shallows/communityapi --go_opt=paths=source_relative --go_out=communityapi community.metrics.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mcommunity.publish.proto=github.com/retrovibed/retrovibed/shallows/communityapi --go_opt=Mcommunity.proto=github.com/retrovibed/retrovibed/shallows/communityapi --go_opt=paths=source_relative --go_out=communityapi community.publish.proto"),
		// ddisc
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mddisc.peers.proto=github.com/retrovibed/retrovibed/shallows/ddiscapi --go_opt=paths=source_relative --go_out=ddiscapi ddisc.peers.proto"),
		// settings
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mstorage.proto=github.com/retrovibed/retrovibed/shallows/cmd/retrovibe/daemons --go_opt=paths=source_relative --go_out=cmd/retrovibe/daemons storage.proto"),
		gruntime.New("protoc --proto_path=../.proto --go_opt=Mtorrent.proto=github.com/retrovibed/retrovibed/shallows/cmd/retrovibe/daemons --go_opt=paths=source_relative --go_out=cmd/retrovibe/daemons torrent.proto"),
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
