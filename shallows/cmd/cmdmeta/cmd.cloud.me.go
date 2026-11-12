package cmdmeta

import (
	"context"
	"net/http"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/retroapi/errorsx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"golang.org/x/crypto/ssh"
)

// CloudMe manages your own profile, and the other profiles associated with your ssh key.
type CloudMe struct {
	Ls      CloudMeLs      `cmd:"" name:"ls" help:"show the current profile and every profile associated with your ssh key"`
	Edit    CloudMeEdit    `cmd:"" name:"edit" help:"edit your profile"`
	Enable  CloudMeEnable  `cmd:"" name:"enable" help:"enable a profile associated with your ssh key"`
	Disable CloudMeDisable `cmd:"" name:"disable" help:"disable a profile associated with your ssh key, the current profile cannot be disabled"`
}

// cloudClient resolves the profile to act as and returns a client authorized for it. the profile is selected
// by --account-id, or interactively when the ssh key is associated with multiple.
func cloudClient(ctx context.Context, gctx *cmdopts.Global, signer ssh.Signer) (*http.Client, *authn.Session, error) {
	session, err := authn.AutoRegistrationOnce(ctx, signer, gctx.AccountID, false, cmdopts.ChooseProfile)
	if err != nil {
		return nil, nil, errorsx.Wrap(err, "unable to determine the current profile")
	}

	c, err := authn.AutoJWTClient(ctx, signer, session.Account.GetId())
	if err != nil {
		return nil, nil, errorsx.Wrap(err, "unable to authenticate")
	}

	return c, session, nil
}
