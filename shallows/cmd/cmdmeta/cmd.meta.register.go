package cmdmeta

import (
	"io"
	"os"

	"golang.org/x/crypto/ssh"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sshx"
)

type CloudRegister struct{}

func (t CloudRegister) Run(gctx *cmdopts.Global) (err error) {
	id, err := sshx.Load(env.PrivateKeyPath())
	if err != nil {
		return err
	}

	session, err := authn.Register(gctx.Context)
	if err != nil {
		return errorsx.Wrap(err, "unable to register")
	}

	return t.run(os.Stdout, id, session)
}

func (t CloudRegister) run(w io.Writer, id ssh.Signer, session *authn.Session) error {
	return authn.PrintIdentity(w, id, session)
}
