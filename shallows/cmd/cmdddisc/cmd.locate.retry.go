package cmdddisc

import (
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

type cmdLocateRetry struct {
	ID            string `flag:"" name:"id" help:"id of the locate request to retry" required:"true"`
	ResetAttempts bool   `flag:"" name:"reset-attempts" help:"also reset the attempt counter, restarting the cooldown backoff from its base" default:"false"`
}

func (t cmdLocateRetry) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID, daemon *cmdopts.Endpoint) (err error) {
	signer, err := id.Signer()
	if err != nil {
		return errorsx.Wrap(err, "failed to create signer")
	}

	c := authn.AutoOauth2Client(gctx.Context, tls.Config(), authn.EndpointSSHAuth(daemon.Endpoint), authn.SSHTokenSourceOptionSigner(signer))
	cc := authn.AuthzClientLibrary(tls.Config(), c, daemon.Endpoint)

	resp, err := ddiscapi.LocateRetry(gctx.Context, cc, daemon.Endpoint, t.ID, &ddiscapi.LocateRetryRequest{
		ResetAttempts: t.ResetAttempts,
	})
	if err != nil {
		return err
	}

	log.Println("locate request queued for retry", spew.Sdump(resp.Locate))

	return nil
}
