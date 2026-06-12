package cmdmeta

import (
	"os"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/sshx"
)

type Identity struct {
	Generate  GenerateID   `cmd:"" help:"bootstrap the identity of the device itself, allows you to provide a seed for consistent generation"`
	Bootstrap Bootstrap    `cmd:"" help:"bootstrap authorized users into the system, used to initially provision the system"`
	Add       IdenAdd      `cmd:"" help:"add a user by public key via the api"`
	Show      IdenDisplay  `cmd:"" help:"display current identity"`
	Register  IdenRegister `cmd:"" help:"register the current identity with the cloud service"`
}

type IdenDisplay struct{}

func (t IdenDisplay) Run(gctx *cmdopts.Global, id *cmdopts.SSHID) (err error) {
	signer, err := id.Signer()
	if err != nil {
		return err
	}
	return authn.PrintIdentity(os.Stdout, signer, nil)
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

	return authn.PrintIdentity(os.Stdout, id, nil)
}
