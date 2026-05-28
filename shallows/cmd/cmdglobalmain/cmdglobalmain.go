package cmdglobalmain

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	"github.com/gofrs/uuid/v5"
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
	"github.com/retrovibed/retrovibed/shallows/internal/envx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/internal/userx"
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
		Version   cmdopts.Version       `cmd:"" help:"display versioning information"`
		Identity  cmdmeta.Identity      `cmd:"" help:"identity management commands"`
		Media     cmdmedia.Commands     `cmd:"" help:"media metadata management (import/export)"`
		Library   cmdlibrary.Commands   `cmd:"" help:"manage your media library"`
		Torrent   cmdtorrent.Commands   `cmd:"" help:"torrent commands"`
		Community cmdcommunity.Commands `cmd:"" help:"community commands"`
		Discovery cmdddisc.Commands     `cmd:"" help:"media discovery commands, used to manage discovery of media"`
		ETL       cmdetl.Commands       `cmd:"" help:"etl commands for processing jsonl through llm endpoints"`
		Daemon    daemons.Command       `cmd:"" help:"run the backend daemon" default:"true"`
		Console   cmdopts.CmdExec       `cmd:"" hidden:"" help:"open the retrovibe console (ui)"`
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

	go debugx.DumpOnSignal(shellCli.Context, syscall.SIGUSR2)
	go cmdopts.Cleanup(shellCli.Context, shellCli.Shutdown, shellCli.Cleanup, os.Kill, os.Interrupt, syscall.SIGTERM)(func() {
		log.Println("waiting for systems to shutdown")
	})

	go debugx.OnSignal(shellCli.Context, func(ctx context.Context) error {
		type profilecfg struct {
			Mode     string        `json:"mode,omitempty"`
			Duration time.Duration `json:"duration,omitempty"`
		}

		var (
			cfg = profilecfg{
				Mode:     "cpu",
				Duration: time.Minute,
			}
		)

		path := userx.DefaultRuntimeDirectory("profile.cfg")
		if err := json.Unmarshal(errorsx.Zero(fsx.AutoCached(path, func() ([]byte, error) { return json.Marshal(cfg) })), &cfg); err != nil {
			log.Println("failed to load profiling configuration", err)
		}

		dctx, done := context.WithTimeout(ctx, cfg.Duration)
		defer done()
		log.Println("PROFILING INITIATED", cfg.Mode, cfg.Duration)
		defer log.Println("PROFILING COMPLETED", cfg.Mode, cfg.Duration)

		switch cfg.Mode {
		case "trace":
			return debugx.Trace(envx.String(os.TempDir(), userx.DefaultRuntimeDirectory()))(dctx)
		case "heap":
			return debugx.Heap(envx.String(os.TempDir(), userx.DefaultRuntimeDirectory()))(dctx)
		case "mem":
			return debugx.Memory(envx.String(os.TempDir(), userx.DefaultRuntimeDirectory()))(dctx)
		case "alloc":
			return debugx.Allocs(envx.String(os.TempDir(), userx.DefaultRuntimeDirectory()))(dctx)
		case "block":
			return debugx.Block(envx.String(os.TempDir(), userx.DefaultRuntimeDirectory()))(dctx)
		default:
			return debugx.CPU(envx.String(os.TempDir(), userx.DefaultRuntimeDirectory()))(dctx)
		}
	}, syscall.SIGUSR1)

	tsstarted := time.Now().UTC()
	parser := kong.Must(
		&shellCli,
		kong.Name(userx.DefaultRelRoot()),
		kong.Description("daemon"),
		kong.Vars{
			"vars_private_key":                  env.PrivateKeyPath(),
			"vars_user_configuration_directory": userx.DefaultConfigDir(userx.DefaultRelRoot()),
			"vars_date_started":                 tsstarted.UTC().Format(time.DateOnly),
			"vars_timestamp_started":            tsstarted.UTC().Format(time.RFC3339),
			"vars_random_seed":                  uuid.Must(uuid.NewV4()).String(),
			"vars_cores":                        strconv.Itoa(runtime.GOMAXPROCS(0)),
			"vars_downloads_directory":          userx.DefaultDownloadDirectory(),
			"env_auto_mdns":                     env.AutoMDNS,
			"env_auto_bootstrap":                env.AutoBootstrap,
			"env_auto_discovery":                env.AutoDiscovery,
			"env_self_signed_hosts":             env.SelfSignedHosts,
			"env_daemon_socket":                 env.DaemonSocket,
			"env_torrent_port":                  env.TorrentPort,
			"env_torrent_private":               env.TorrentPrivateNetwork,
			"env_torrent_ipv4":                  env.TorrentPublicIP4,
			"env_torrent_ipv6":                  env.TorrentPublicIP6,
			"env_torrent_directory_watch":       env.TorrentDirectoryWatch,
			"env_torrent_logging":               env.TorrentLogging,
			"env_torrent_debug":                 env.TorrentDebug,
			"env_auto_identify_media":           env.AutoIdentifyMedia,
			"env_auto_locate_media":             env.AutoIdentifyMedia,
			"env_auto_archive":                  env.AutoArchive,
			"env_auto_reclaim":                  env.AutoReclaim,
		},
		kong.UsageOnError(),
		kong.Bind(
			&shellCli.TLSConfig,
			&shellCli.Global,
			&shellCli.PeerID,
			&shellCli.SSHID,
		),
		kong.BindTo(cmdopts.DeeppoolClientDefault{}, (*cmdopts.DeeppoolClient)(nil)),
		kong.TypeMapper(reflect.TypeOf(&net.IP{}), kong.MapperFunc(cmdopts.ParseIP)),
		kong.TypeMapper(reflect.TypeOf(&net.TCPAddr{}), kong.MapperFunc(cmdopts.ParseTCPAddr)),
		kong.TypeMapper(reflect.TypeOf([]*net.TCPAddr(nil)), kong.MapperFunc(cmdopts.ParseTCPAddrArray)),
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
