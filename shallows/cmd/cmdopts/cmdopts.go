package cmdopts

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"runtime/debug"
	"sync"

	"github.com/james-lawrence/torrent/dht/krpc"
	"github.com/retrovibed/retrovibed/internal/env"
	"github.com/retrovibed/retrovibed/internal/fsx"
	"github.com/retrovibed/retrovibed/internal/sshx"
	"github.com/retrovibed/retrovibed/internal/userx"
	"golang.org/x/crypto/ssh"
)

// Global command fields.
type Global struct {
	Context  context.Context    `kong:"-"`
	Shutdown context.CancelFunc `kong:"-"`
	Cleanup  *sync.WaitGroup    `kong:"-"`
}

type Version struct{}

func (t Version) Run(ctx *Global) (err error) {
	var (
		ok   bool
		info *debug.BuildInfo
	)

	if info, ok = debug.ReadBuildInfo(); ok {
		log.Println(info.Main.Path, info.Main.Version)
		return nil
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
	ssh.Signer
}

func (t *SSHID) AfterApply() (err error) {
	if t.Signer, err = sshx.AutoCached(sshx.NewKeyGen(), env.PrivateKeyPath()); err != nil {
		return err
	}

	return nil
}
