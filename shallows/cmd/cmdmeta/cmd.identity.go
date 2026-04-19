package cmdmeta

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/sshx"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"golang.org/x/crypto/ssh"
)

func printIdentity(w io.Writer, s ssh.Signer, session *metaapi.Session) error {
	var (
		err3 error
	)
	_, err1 := fmt.Fprintln(w, "fingerprint", ssh.FingerprintSHA256(s.PublicKey()))
	_, err2 := fmt.Fprintln(w, "identity   ", md5x.String(ssh.FingerprintSHA256(s.PublicKey())))
	if session != nil {
		_, err3 = fmt.Fprintln(w, "account    ", session.Account.Id)
	}
	_, err4 := fmt.Fprintln(w, "public     ", strings.TrimSpace(string(ssh.MarshalAuthorizedKey(s.PublicKey()))))
	_, err5 := fmt.Fprintln(w, "base64     ", base64.URLEncoding.EncodeToString(s.PublicKey().Marshal()))
	return errorsx.Compact(err1, err2, err3, err4, err5)
}

type Identity struct {
	Generate  GenerateID  `cmd:"" help:"bootstrap the identity of the device itself, allows you to provide a seed for consistent generation"`
	Bootstrap Bootstrap   `cmd:"" help:"bootstrap authorized users into the system, used to initially provision the system"`
	Add       IdenAdd     `cmd:"" help:"add a user by public key via the api"`
	Register  Register    `cmd:"" help:"register the current identity with the cloud service"`
	Show      IdenDisplay `cmd:"" help:"display current identity"`
}

type IdenDisplay struct{}

func (t IdenDisplay) Run(gctx *cmdopts.Global, id *cmdopts.SSHID) (err error) {
	signer, err := id.Signer()
	if err != nil {
		return err
	}
	return printIdentity(os.Stdout, signer, nil)
}

type GenerateID struct {
	Force bool   `flag:"" name:"force" help:"force creation if a identity already exists, this is to prevent you from accidently destroying your identity"`
	Seed  string `arg:"" name:"seed" help:"used to seed the key generation, this command is used for when you want to maintain a persistent account identity easily" required:"true"`
}

func (t GenerateID) Run(gctx *cmdopts.Global) (err error) {
	id, err := sshx.Seeded(gctx.Context, t.Seed, t.Force, env.PrivateKeyPath())
	if err != nil {
		return err
	}

	return printIdentity(os.Stdout, id, nil)
}

type Register struct {
	Seed string `arg:"" name:"seed" help:"used to seed the key generation, this command is used for when you want to maintain a persistent account identity easily" `
}

func (t Register) Run(gctx *cmdopts.Global) (err error) {
	var (
		session *metaapi.Session
	)

	id, err := sshx.AutoCached(sshx.NewKeyGenSeeded(t.Seed), env.PrivateKeyPath())
	if err != nil {
		return err
	}

	if session, err = authn.Register(gctx.Context); err != nil {
		return errorsx.Wrap(err, "unable to register")
	}

	if err := printIdentity(os.Stdout, id, session); err != nil {
		return err
	}

	return nil
}
