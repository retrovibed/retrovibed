package cmdmeta

import (
	"context"
	"database/sql"
	"os"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sshx"
	"github.com/retrovibed/retrovibed/shallows/meta/identityssh"
)

type Bootstrap struct {
	PublicKey      BootstrapPublicKey  `cmd:"" help:"a single ssh public key into the system"`
	AuthorizedFile BootstrapAuthorized `cmd:"" help:"authorize public keys from a authorized keys file"`
}

type BootstrapPublicKey struct {
	PublicKey string `arg:"" name:"pubkey" help:"public key to add" required:"true"`
}

func (t BootstrapPublicKey) Run(gctx *cmdopts.Global) (err error) {
	var db *sql.DB

	if db, err = cmdopts.DatabaseMeta(gctx.Context); err != nil {
		return err
	}
	defer db.Close()

	ctx, done := context.WithTimeout(gctx.Context, 10*time.Second)
	defer done()

	return t.run(ctx, db)
}

func (t BootstrapPublicKey) run(ctx context.Context, db *sql.DB) (err error) {
	var parsed sshx.Parsed

	if parsed.PublicKey, parsed.Comment, parsed.Options, _, err = ssh.ParseAuthorizedKey([]byte(t.PublicKey)); err != nil {
		return errorsx.Wrap(err, "unable to parse public key")
	}

	return identityssh.InitializeAdmin(ctx, db, parsed)
}

type BootstrapAuthorized struct {
	Path string `arg:"" name:"authorized_keys" help:"path to authorized key file to import" required:"true"`
}

func (t BootstrapAuthorized) Run(gctx *cmdopts.Global) (err error) {
	var db *sql.DB

	if db, err = cmdopts.DatabaseMeta(gctx.Context); err != nil {
		return err
	}
	defer db.Close()

	ctx, done := context.WithTimeout(gctx.Context, 10*time.Second)
	defer done()

	return t.run(ctx, db)
}

func (t BootstrapAuthorized) run(ctx context.Context, db *sql.DB) (err error) {
	encoded, err := os.ReadFile(t.Path)
	if err != nil {
		return errorsx.Wrapf(err, "unable to read authorized keys from %s", t.Path)
	}

	for parsed := range sshx.ParseAuthorizedKeys(encoded) {
		if err = identityssh.InitializeAdmin(ctx, db, parsed); err != nil {
			return err
		}
	}

	return nil
}
