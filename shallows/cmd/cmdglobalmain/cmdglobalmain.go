package cmdglobalmain

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	"github.com/egdaemon/gdx"
	"github.com/egdaemon/gdx/konggdx"
	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/userx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdcommunity"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdddisc"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdetl"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdlibrary"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdmedia"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdmeta"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdtorrent"
	"github.com/retrovibed/retrovibed/shallows/cmd/retrovibe/daemons"
	"github.com/retrovibed/retrovibed/shallows/internal/debugx"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/willabides/kongplete"
	"golang.zx2c4.com/wireguard/device"

	_ "github.com/benbjohnson/immutable"
)

func Hostname() string {
	return stringsx.FirstNonBlank(errorsx.Zero(os.Hostname()), "localhost")
}

func Username() string {
	u := userx.CurrentUserOrDefault(userx.Zero())
	return stringsx.FirstNonBlank(u.Name, u.Username, u.Uid)
}

func Main(args ...string) {
	var shellCli struct {
		cmdopts.Global
		cmdopts.TLSConfig
		cmdopts.PeerID
		cmdopts.SSHID
		cmdopts.Endpoint
		Version   cmdopts.Version        `cmd:"" help:"display versioning information"`
		Identity  cmdmeta.Identity       `cmd:"" help:"identity management commands"`
		Cloud     cmdmeta.Cloud          `cmd:"" name:"cloud" help:"retrovibe.space cloud service commands"`
		U12t      cmdmeta.Usermanagement `cmd:"" name:"u12t" help:"user management commands"`
		Media     cmdmedia.Commands      `cmd:"" help:"media metadata management (import/export)"`
		Library   cmdlibrary.Commands    `cmd:"" help:"manage your media library"`
		Torrent   cmdtorrent.Commands    `cmd:"" help:"torrent commands"`
		Community cmdcommunity.Commands  `cmd:"" help:"community commands"`
		Ddisc     cmdddisc.Commands      `cmd:"" help:"media discovery commands, used to manage discovery of media"`
		ETL       cmdetl.Commands        `cmd:"" help:"etl commands for processing jsonl through llm endpoints"`
		Gdx       konggdx.Commands       `cmd:"" help:"pull profiles/traces from a running eg debug socket"`
		Daemon    daemons.Command        `cmd:"" help:"run the backend daemon" default:"true"`
		Console   cmdopts.CmdExec        `cmd:"" hidden:"" help:"open the retrovibe console (ui)"`
	}

	var (
		err error
		ctx *kong.Context
	)

	shellCli.Context, shellCli.Shutdown = context.WithCancel(context.Background())
	shellCli.Cleanup = &sync.WaitGroup{}
	shellCli.Console = cmdopts.NewCmdExec("retrovibe")

	log.SetFlags(log.Lshortfile | log.LUTC | log.Ltime)
	log.SetPrefix(fmt.Sprintf("%d ", os.Getpid()))

	log.Println("wireguard preallocated buffers per pool", device.PreallocatedBuffersPerPool)
	debugx.Println("jwt checksum", md5x.FormatUUID(md5x.Digest(env.JWTSecret())))

	go cmdopts.Cleanup(shellCli.Context, shellCli.Shutdown, shellCli.Cleanup, os.Kill, os.Interrupt, syscall.SIGTERM)(func() {
		log.Println("waiting for systems to shutdown")
	})

	go func() {
		path := userx.DefaultRuntimeDirectory("gdx.socket")
		os.Remove(path)

		l, err := net.Listen("unix", path)
		if err != nil {
			log.Println("unable to bind gdx debug socket", err)
			return
		}
		defer l.Close()

		go func() {
			<-shellCli.Context.Done()
			l.Close()
		}()

		if err := http.Serve(l, gdx.NewHTTPFn(gdx.Options().FromEnv())); err != nil && shellCli.Context.Err() == nil {
			log.Println("gdx debug server stopped", err)
		}
	}()

	tsstarted := time.Now().UTC()
	parser := kong.Must(
		&shellCli,
		kong.Name(userx.DefaultRelRoot()),
		kong.Description("daemon"),
		kong.Vars{
			"vars_private_key":                  env.PrivateKeyPath(),
			"vars_user_configuration_directory": userx.DefaultConfigDir(userx.DefaultRelRoot()),
			"vars_user_cache_directory":         userx.DefaultCacheDirectory(userx.DefaultRelRoot()),
			"vars_date_started":                 tsstarted.UTC().Format(time.DateOnly),
			"vars_timestamp_started":            tsstarted.UTC().Format(time.RFC3339),
			"vars_random_seed":                  uuid.Must(uuid.NewV4()).String(),
			"vars_cores":                        strconv.Itoa(runtime.GOMAXPROCS(0)),
			"vars_downloads_directory":          userx.DefaultDownloadDirectory(),
			"vars_media_directory":              env.MediaDir(),
			"env_http_endpoint":                 env.Endpoint,
			"env_mdns_advertise":                env.MDNSAdvertise,
			"env_mdns_discovery":                env.MDNSDiscovery,
			"env_auto_bootstrap":                env.AutoBootstrap,
			"env_auto_discovery":                env.AutoDiscovery,
			"env_auto_peertube":                 env.AutoPeerTube,
			"env_peertube_domain":               env.PeerTubeDomain,
			"env_self_signed_hosts":             env.SelfSignedHosts,
			"env_daemon_socket":                 env.DaemonSocket,
			"env_torrent_port":                  env.TorrentPort,
			"env_torrent_private":               env.TorrentPrivateNetwork,
			"env_torrent_ipv4":                  env.TorrentPublicIP4,
			"env_torrent_ipv6":                  env.TorrentPublicIP6,
			"env_torrent_directory_watch":       env.TorrentDirectoryWatch,
			"env_torrent_logging":               env.TorrentLogging,
			"env_torrent_debug":                 env.TorrentDebug,
			"env_discovery_index_ratio":         env.DDiscIndexRatio,
			"env_discovery_p2p_locate":          env.DDiscP2PLocate,
			"env_auto_identify_media":           env.AutoIdentifyMedia,
			"env_auto_locate_media":             env.AutoIdentifyMedia,
			"env_auto_archive":                  env.AutoArchive,
			"env_auto_reclaim":                  env.AutoReclaim,
			"vars_gdx_socket":                   userx.DefaultRuntimeDirectory(gdx.DefaultSocket),
		},
		kong.UsageOnError(),
		kong.Bind(
			&shellCli.TLSConfig,
			&shellCli.Global,
			&shellCli.PeerID,
			&shellCli.SSHID,
			&shellCli.Endpoint,
			shellCli.Context,
		),
		kong.BindTo(cmdopts.DeeppoolClientDefault{SSHID: &shellCli.SSHID}, (*cmdopts.DeeppoolClient)(nil)),
		kong.TypeMapper(reflect.TypeOf(&net.IP{}), kong.MapperFunc(cmdopts.ParseIP)),
		kong.TypeMapper(reflect.TypeOf(&net.TCPAddr{}), kong.MapperFunc(cmdopts.ParseTCPAddr)),
		kong.TypeMapper(reflect.TypeOf([]*net.TCPAddr(nil)), kong.MapperFunc(cmdopts.ParseTCPAddrArray)),
		kong.NamedMapper("durationinf", kong.MapperFunc(cmdopts.ParseDurationInf)),
		kong.NamedMapper("envvar", kong.MapperFunc(cmdopts.ParseEnviron)),
	)

	// Run kongplete.Complete to handle completion requests
	kongplete.Complete(parser)

	if ctx, err = parser.Parse(args); err != nil {
		log.Fatalln(err)
		return
	}

	if err = errorsx.LogErr(ctx.Run()); err != nil {
		shellCli.Shutdown()
	}

	shellCli.Cleanup.Wait()
	debugx.Println("system cleanup completed")
	ctx.FatalIfErrorf(err)
}
