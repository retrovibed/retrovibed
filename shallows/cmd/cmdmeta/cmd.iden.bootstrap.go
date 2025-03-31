package cmdmeta

import (
	"context"
	"database/sql"
	"os"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/retrovibed/retrovibed/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/langx"
	"github.com/retrovibed/retrovibed/internal/sshx"
	"github.com/retrovibed/retrovibed/meta"
	"github.com/retrovibed/retrovibed/meta/identityssh"
)

type Bootstrap struct {
	PublicKey      BootstrapPublicKey  `cmd:"" help:"a single ssh public key into the system"`
	AuthorizedFile BootstrapAuthorized `cmd:"" help:"authorize public keys from a authorized keys file"`
}

type BootstrapPublicKey struct {
	PublicKey string `arg:"" name:"pubkey" help:"public key to add" required:"true"`
}

func (t BootstrapPublicKey) Run(gctx *cmdopts.Global) (err error) {
	var (
		db     *sql.DB
		parsed sshx.Parsed
	)

	if db, err = Database(gctx.Context); err != nil {
		return err
	}
	defer db.Close()

	ctx, done := context.WithTimeout(gctx.Context, 10*time.Second)
	defer done()

	if parsed.PublicKey, parsed.Comment, parsed.Options, _, err = ssh.ParseAuthorizedKey([]byte(t.PublicKey)); err != nil {
		return errorsx.Wrap(err, "unable to parse public key")
	}

	return identityssh.ImportParsed(ctx, db, parsed)
}

type BootstrapAuthorized struct {
	Path string `arg:"" name:"authorized_keys" help:"path to authorized key file to import" required:"true"`
}

func (t BootstrapAuthorized) Run(gctx *cmdopts.Global) (err error) {
	var (
		db *sql.DB
	)

	if db, err = Database(gctx.Context); err != nil {
		return err
	}
	defer db.Close()

	ctx, done := context.WithTimeout(gctx.Context, 10*time.Second)
	defer done()

	encoded, err := os.ReadFile(t.Path)
	if err != nil {
		return errorsx.Wrapf(err, "unable to read authorized keys from %s", t.Path)
	}

	for parsed := range sshx.ParseAuthorizedKeys(encoded) {
		p := meta.Profile{
			Description: parsed.Comment,
		}
		if err = meta.ProfileInsertWithDefaults(ctx, db, p).Scan(&p); err != nil {
			return errorsx.Wrap(err, "unable to create profile")
		}

		authz := langx.Clone(meta.Authz{}, meta.AuthzOptionAdmin)
		if err = meta.AuthzInsertWithDefaults(ctx, db, authz).Scan(&authz); err != nil {
			return errorsx.Wrap(err, "unable to setup authorizations")
		}

		iden := identityssh.Identity{
			ID:        sshx.FingerprintMD5(parsed.PublicKey),
			PublicKey: sshx.EncodeBase64PublicKey(parsed.PublicKey),
			ProfileID: p.ID,
		}

		if err = identityssh.IdentityInsertWithDefaults(ctx, db, iden).Scan(&iden); err != nil {
			return errorsx.Wrap(err, "unable create ssh identity")
		}
	}

	return nil
}
