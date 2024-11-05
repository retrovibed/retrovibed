package main

// import (
// 	"log"
// 	"os"

// 	"github.com/james-lawrence/torrent"
// 	"github.com/james-lawrence/torrent/metainfo"
// )

// func main() {
// 	src, err := os.Open("hello.world.txt")
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	minfo, err := metainfo.NewFromReader(src, metainfo.OptionDisplayName("hello.world.txt"))
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	md1, err := torrent.NewFromInfo(*minfo)
// 	// md1, err := torrent.NewFromReader(src, torrent.OptionDisplayName("hello.world.txt"))
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	log.Println("magnet uri:", torrent.NewMagnet(md1).String())
// 	md2, err := torrent.NewFromFile("hello.world.txt")
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	log.Println("magnet uri:", torrent.NewMagnet(md2).String())
// }

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
