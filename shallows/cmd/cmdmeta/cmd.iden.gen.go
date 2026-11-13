package cmdmeta

import (
	"context"
	"database/sql"
	"os"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/sshx"
	"github.com/retrovibed/retrovibed/shallows/meta/identityssh"
)

type GenerateID struct {
	Force     bool   `flag:"" name:"force" help:"force creation if a identity already exists, this is to prevent you from accidently destroying your identity"`
	Bootstrap bool   `flag:"" name:"bootstrap" help:"also bootstrap the generated identity as an authorized admin of the system"`
	Seed      string `arg:"" name:"seed" help:"used to seed the key generation, this command is used for when you want to maintain a persistent account identity easily" required:"true"`
}

func (t GenerateID) Run(gctx *cmdopts.Global, _id *cmdopts.SSHID) (err error) {
	id, err := sshx.Seeded(gctx.Context, t.Seed, t.Force, _id.KeyPath)
	if err != nil {
		return err
	}

	if t.Bootstrap {
		var db *sql.DB

		if db, err = cmdopts.DatabaseMeta(gctx.Context); err != nil {
			return err
		}
		defer db.Close()

		ctx, done := context.WithTimeout(gctx.Context, 10*time.Second)
		defer done()

		if err = identityssh.InitializeAdmin(ctx, db, sshx.Parsed{PublicKey: id.PublicKey()}); err != nil {
			return err
		}
	}

	return authn.PrintIdentity(os.Stdout, id, nil)
}
