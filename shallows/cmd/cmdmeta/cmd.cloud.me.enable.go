package cmdmeta

import (
	"log"

	"github.com/retrovibed/retrovibed/retroapi/errorsx"
	"github.com/retrovibed/retrovibed/retroapi/me"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
)

type CloudMeEnable struct {
	Profile string `arg:"" name:"profile-id" help:"id of the profile to enable, see cloud me ls"`
}

func (t CloudMeEnable) Run(gctx *cmdopts.Global, id *cmdopts.SSHID) (err error) {
	signer, err := id.Signer()
	if err != nil {
		return err
	}

	c, _, err := cloudClient(gctx.Context, gctx, signer)
	if err != nil {
		return err
	}

	p, err := me.EnableProfile(gctx.Context, c, t.Profile)
	if err != nil {
		if me.Forbidden(err) {
			return errorsx.Wrap(err, "the current profile cannot be changed, select a different one with --account-id")
		}

		if me.NotFound(err) {
			return errorsx.Wrap(err, "no profile with that id is associated with your ssh key")
		}

		return err
	}

	log.Println("enabled", p.Id, p.Display)

	return nil
}
