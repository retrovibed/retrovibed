package cmdtorrent

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	utp "github.com/anacrolix/go-libutp"
	"github.com/davecgh/go-spew/spew"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/internal/dhtx"
	"github.com/retrovibed/retrovibed/internal/env"
	"github.com/retrovibed/retrovibed/internal/envx"
	"github.com/retrovibed/retrovibed/internal/errorsx"
)

type cmdDHTTest struct {
	Debug           bool     `flag:"" name:"debug" help:"enable debug logging" default:"false" negatable:""`
	Port            uint16   `flag:"" name:"port" help:"port to listen on" default:"9654"`
	AnnounceTargets []string `arg:"" name:"target" help:"infohashes to announce"`
}

func (t cmdDHTTest) Run(gctx *cmdopts.Global) (err error) {
	var (
		dhtdebug = dht.OptionLogger(torrent.LogDiscard())
	)

	if envx.Boolean(t.Debug, env.DHTDebug) {
		dhtdebug = dht.OptionLogger(log.New(os.Stderr, "[dht] ", log.Flags()))
	}
	peerid := int160.New(errorsx.Zero(os.Hostname()))

	dhts, err := dht.NewServer(
		32,
		dht.OptionNodeID(peerid),
		dht.OptionMuxer(dht.DefaultMuxer()),
		dht.OptionBootstrapGlobal,
		dht.OptionUPnP,
		dhtdebug,
	)
	if err != nil {
		return errorsx.Wrap(err, "unable to setup dht")
	}

	s, err := utp.NewSocket("udp", fmt.Sprintf(":%d", t.Port))
	if err != nil {
		return errorsx.Wrap(err, "unable to setup socket")
	}

	if err := dhts.Serve(gctx.Context, s); err != nil {
		return errorsx.Wrap(err, "unable to serve dht")
	}

	go dhtx.Statistics(gctx.Context, 5*time.Second, dhts)

	bctx, done := context.WithTimeout(gctx.Context, 30*time.Second)
	defer done()
	bstat, err := dhts.Bootstrap(bctx)
	if errorsx.Ignore(err, context.DeadlineExceeded) != nil {
		return err
	}
	go dhts.TableMaintainer()

	log.Println("bootstrap stats", spew.Sdump(bstat))
	announce := func(_ctx context.Context, target int160.T) {
		ctx, done := context.WithTimeout(_ctx, time.Minute)
		defer done()
		log.Println("announce initiated", target)
		defer log.Println("announce completed", target)

		results, err := dhts.AnnounceTraversal(ctx, target, dht.AnnouncePeer(dhts, false))
		if err != nil {
			log.Println("announce failed", target, err)
			return
		}
		defer results.Close()

		for {
			select {
			case batch := <-results.Peers:
				for _, p := range batch.Peers {
					log.Println("located", target, p, p.AddrPort.Addr().Is4(), p.AddrPort.Addr().Is4In6(), p.AddrPort.Addr().Is6())
				}
			case <-results.Finished():
				return
			case <-gctx.Context.Done():
				return
			}
		}
	}

	for _, v := range t.AnnounceTargets {
		target, err := int160.FromHexEncodedString(v)
		if err != nil {
			log.Println("failed to decode target from", v, err)
			continue
		}

		go func() {
			for {
				announce(gctx.Context, target)
				time.Sleep(10 * time.Second)
			}
		}()
	}

	<-gctx.Context.Done()

	return nil
}
