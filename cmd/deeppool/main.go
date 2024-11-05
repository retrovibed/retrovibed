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
	"github.com/james-lawrence/deeppool/cmd/cmdopts"
	"github.com/james-lawrence/deeppool/internal/x/debugx"
	"github.com/willabides/kongplete"
)

func main() {
	var shellCli struct {
		cmdopts.Global
		Version cmdopts.Version `cmd:"" help:"display versioning information"`
		Daemon  cmdDaemon       `cmd:"" help:"run the backend daemon"`
		Torrent cmdTorrent      `cmd:"" help:"torrent related sub commands"`
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

	parser := kong.Must(
		&shellCli,
		kong.Name("dpool"),
		kong.Description("daemon"),
		kong.Vars{
			"vars_timestamp_started": time.Now().UTC().Format(time.RFC3339),
		},
		kong.UsageOnError(),
		kong.Bind(
			&shellCli.Global,
		),
		kong.TypeMapper(reflect.TypeOf(&net.IP{}), kong.MapperFunc(cmdopts.ParseIP)),
		kong.TypeMapper(reflect.TypeOf(&net.TCPAddr{}), kong.MapperFunc(cmdopts.ParseTCPAddr)),
		kong.TypeMapper(reflect.TypeOf([]*net.TCPAddr(nil)), kong.MapperFunc(cmdopts.ParseTCPAddrArray)),
	)

	// Run kongplete.Complete to handle completion requests
	kongplete.Complete(parser)

	if ctx, err = parser.Parse(os.Args[1:]); cmdopts.ReportError(err) != nil {
		ctx.FatalIfErrorf(err)
	}

	if err = cmdopts.ReportError(ctx.Run()); err != nil {
		shellCli.Shutdown()
	}

	shellCli.Cleanup.Wait()
	ctx.FatalIfErrorf(err)
}
