package cmdopts

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/james-lawrence/torrent/dht/krpc"
	"github.com/pkg/errors"
	"github.com/retrovibed/retrovibed/retroapi/userx"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/envx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/sshx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"golang.org/x/crypto/ssh"
)

// Global command fields.
type Global struct {
	Verbosity int                `help:"increase verbosity of logging" short:"v" type:"counter" default:"0" env:"RETROVIBED_LOGGING_VERBOSITY"`
	Context   context.Context    `kong:"-"`
	Shutdown  context.CancelFunc `kong:"-"`
	Cleanup   *sync.WaitGroup    `kong:"-"`
}

func (t Global) AfterApply() (err error) {
	v := envx.Int(t.Verbosity, env.LoggingVerbosity)
	log.SetFlags(log.Flags() | log.Lshortfile)
	switch v {
	case 4: // NETWORK
		err = langx.FirstNonZero(
			err,
			os.Setenv(env.TorrentDebug, "true"),
			os.Setenv(env.DHTDebug, "true"),
		)
		fallthrough
	case 3: // TRACE
		err = langx.FirstNonZero(
			err,
			os.Setenv(env.TorrentLogging, "true"),
		)
		fallthrough
	case 2: // DEBUG
		fallthrough
	case 1: // INFO
		fallthrough
	default: // ERROR - minimal
	}

	return langx.FirstNonZero(
		err,
		os.Setenv(env.LoggingVerbosity, strconv.Itoa(v)),
	)
}

type Version struct{}

func (t Version) Run(ctx *Global) (err error) {
	if version, err := BuildVersion(); stringsx.Present(version) {
		log.Println(version)
		return nil
	} else {
		log.Println("failed to detect build version", err)
	}

	log.Println("unknown version")
	return nil
}

type PeerID krpc.ID

func (t *PeerID) AfterApply() error {
	rid, err := fsx.AutoCached(userx.DefaultConfigDir(userx.DefaultRelRoot(), "torrent.id"), func() ([]byte, error) {
		var id krpc.ID
		if _, err := rand.Read(id[:]); err != nil {
			return nil, err
		}
		return id[:], nil
	})
	if err != nil {
		return err
	}

	if n := copy(t[:], rid); n != len(t[:]) {
		return fmt.Errorf("invalid length %d vs %d", n, len(t[:]))
	}

	return nil
}

type SSHID struct {
	KeyPath string `flag:"" name:"private-key-path" default:"${vars_private_key}"`
}

func (t *SSHID) Signer() (signer ssh.Signer, err error) {
	signer, err = sshx.Load(t.KeyPath)
	return signer, errors.Wrapf(err, "failed to generate signer: %s", t.KeyPath)
}

type Endpoint struct {
	Endpoint string `flag:"" name:"endpoint" help:"http address for the retrovibed daemon" default:"https://localhost:9998" env:"${env_http_endpoint}"`
}
