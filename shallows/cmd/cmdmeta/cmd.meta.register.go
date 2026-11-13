package cmdmeta

import (
	"io"
	"os"

	"golang.org/x/crypto/ssh"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/retroapi/errorsx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
)

type CloudRegister struct {
	Account string `name:"account" help:"optional id of the account to register with, an administrator of the account must approve the registration"`
}

func (t CloudRegister) Run(gctx *cmdopts.Global, _id *cmdopts.SSHID) (err error) {
	id, err := _id.Signer()
	if err != nil {
		return err
	}

	session, err := authn.AutoRegistrationOnce(gctx.Context, id, t.Account, true, cmdopts.ChooseProfile)
	if err != nil {
		return errorsx.Wrap(err, "unable to register")
	}

	return t.run(os.Stdout, id, session)
}

func (t CloudRegister) run(w io.Writer, id ssh.Signer, session *authn.Session) error {
	return authn.PrintIdentity(w, id, session)
}
