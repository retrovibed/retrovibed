package main

import (
	"context"
	"log"
	"net"
	"os"
	"reflect"
	"sync"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	"github.com/gofrs/uuid"
	"github.com/retrovibed/retrovibed/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/internal/x/debugx"
	"github.com/retrovibed/retrovibed/internal/x/envx"
	"github.com/retrovibed/retrovibed/internal/x/userx"
	"github.com/willabides/kongplete"

	_ "github.com/benbjohnson/immutable"
)

func main() {
	var shellCli struct {
		cmdopts.Global
		cmdopts.PeerID
		cmdopts.SSHID
		Version  cmdopts.Version `cmd:"" help:"display versioning information"`
		Identity cmdMetaIdentity `cmd:"" help:"identity management commands"`
		Daemon   cmdDaemon       `cmd:"" help:"run the backend daemon" default:"true"`
		Torrent  cmdTorrent      `cmd:"" help:"torrent related sub commands"`
	}

	var (
		err error
		ctx *kong.Context
	)

	shellCli.Context, shellCli.Shutdown = context.WithCancel(context.Background())
	shellCli.Cleanup = &sync.WaitGroup{}

	log.SetFlags(log.Lshortfile | log.LUTC | log.Ltime)

	go debugx.DumpOnSignal(shellCli.Context, syscall.SIGUSR2)
	go cmdopts.Cleanup(shellCli.Context, shellCli.Shutdown, shellCli.Cleanup, os.Kill, os.Interrupt)(func() {
		log.Println("waiting for systems to shutdown")
	})

	go debugx.OnSignal(shellCli.Context, func(ctx context.Context) error {
		dctx, done := context.WithTimeout(ctx, envx.Duration(time.Second, "DEEPPOOL_PROFILING_DURATION"))
		defer done()

		log.Println("PROFILING INITIATED")
		defer log.Println("PROFILING COMPLETED")

		switch envx.String("cpu", "DEEPPOOL_PROFILING_STRATEGY") {
		case "heap":
			return debugx.Heap(envx.String(os.TempDir(), "CACHE_DIRECTORY"))(dctx)
		case "mem":
			return debugx.Memory(envx.String(os.TempDir(), "CACHE_DIRECTORY"))(dctx)
		default:
			return debugx.CPU(envx.String(os.TempDir(), "CACHE_DIRECTORY"))(dctx)
		}
	}, syscall.SIGUSR1)

	parser := kong.Must(
		&shellCli,
		kong.Name(userx.DefaultRelRoot()),
		kong.Description("daemon"),
		kong.Vars{
			"vars_timestamp_started": time.Now().UTC().Format(time.RFC3339),
			"vars_random_seed":       uuid.Must(uuid.NewV4()).String(),
		},
		kong.UsageOnError(),
		kong.Bind(
			&shellCli.Global,
		),
		kong.Bind(
			&shellCli.PeerID,
		),
		kong.Bind(
			&shellCli.SSHID,
		),
		kong.TypeMapper(reflect.TypeOf(&net.IP{}), kong.MapperFunc(cmdopts.ParseIP)),
		kong.TypeMapper(reflect.TypeOf(&net.TCPAddr{}), kong.MapperFunc(cmdopts.ParseTCPAddr)),
		kong.TypeMapper(reflect.TypeOf([]*net.TCPAddr(nil)), kong.MapperFunc(cmdopts.ParseTCPAddrArray)),
	)

	// Run kongplete.Complete to handle completion requests
	kongplete.Complete(parser)

	if ctx, err = parser.Parse(os.Args[1:]); cmdopts.ReportError(err) != nil {
		log.Fatalln(err)
	}

	if err = cmdopts.ReportError(ctx.Run()); err != nil {
		shellCli.Shutdown()
	}

	shellCli.Cleanup.Wait()
	ctx.FatalIfErrorf(err)
}
