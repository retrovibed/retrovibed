package cmdmeta

import (
	"log"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/backups"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sshx"
)

type CloudBackup struct{}

// takes an encrypted snapshot of the metadata database and uploads it now, rather than
// waiting for the daemon's hourly pass.
func (t CloudBackup) Run(gctx *cmdopts.Global) (err error) {
	id, err := sshx.Load(env.PrivateKeyPath())
	if err != nil {
		return err
	}

	c, err := authn.AutoJWTClient(gctx.Context, id)
	if err != nil {
		return errorsx.Wrap(err, "unable to authenticate")
	}

	key, err := backups.ResolveKey(gctx.Context, c)
	if err != nil {
		return err
	}

	db, err := cmdopts.DatabaseMeta(gctx.Context)
	if err != nil {
		return err
	}
	defer db.Close()

	m, err := backups.Run(gctx.Context, c, db, cmdopts.MachineID(), key)
	if err != nil {
		return err
	}

	log.Println("backup uploaded", m.Id, m.Bytes)
	return nil
}
