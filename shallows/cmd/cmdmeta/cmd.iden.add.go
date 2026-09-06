package cmdmeta

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"golang.org/x/crypto/ssh"
)

type IdenAdd struct {
	PublicKey string `arg:"" name:"pubkey" help:"public key to add" required:"true"`
	Username  string `field:"" name:"username" help:"username for the profile, defaults to the ssh comment if not provided"`
}

func (t IdenAdd) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID, daemon *cmdopts.Endpoint) (err error) {
	ctx, done := context.WithTimeout(gctx.Context, 10*time.Second)
	defer done()

	signer, err := id.Signer()
	if err != nil {
		return errorsx.Wrap(err, "failed to create signer")
	}

	c := authn.AutoOauth2Client(
		gctx.Context,
		tls.Config(),
		authn.EndpointSSHAuth(daemon.Endpoint),
		authn.SSHTokenSourceOptionSigner(signer),
	)
	cc := authn.AuthzClientLibrary(tls.Config(), c, daemon.Endpoint)

	return t.run(ctx, daemon.Endpoint, cc)
}

func (t IdenAdd) run(ctx context.Context, endpoint string, c *http.Client) (err error) {
	pubkey, comment, _, _, err := ssh.ParseAuthorizedKey([]byte(t.PublicKey))
	if err != nil {
		return errorsx.Wrap(err, "invalid public key")
	}

	var result metaapi.ProfileCreateResponse

	encoded, err := jsonx.Marshal(&metaapi.ProfileCreateRequest{
		Profile:   &metaapi.Profile{Display: langx.FirstNonZero(t.Username, comment)},
		PublicKey: strings.TrimSpace(string(ssh.MarshalAuthorizedKey(pubkey))),
	})

	if err != nil {
		return errorsx.Wrap(err, "unable to encode request")
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/meta/u12t/", endpoint),
		bytes.NewReader(encoded),
	)
	if err != nil {
		return errorsx.Wrap(err, "unable to create http request")
	}

	resp, err := httpx.AsError(c.Do(req))
	if err != nil {
		return errorsx.Wrap(err, "http request failed")
	}

	if err = httpx.DecodeJSON(resp, &result); err != nil {
		return err
	}

	encoded, err = jsonx.Marshal(&metaapi.AuthzGrantRequest{
		Token: &metaapi.Token{
			LibraryRead: true,
		},
	})

	if err != nil {
		return errorsx.Wrap(err, "unable to encode authz request")
	}

	req, err = http.NewRequestWithContext(ctx,
		http.MethodPost,
		fmt.Sprintf("%s/meta/authz/%s",
			endpoint,
			result.Profile.GetId(),
		),
		bytes.NewReader(encoded),
	)
	if err != nil {
		return errorsx.Wrap(err, "unable to create authz http request")
	}

	var grant metaapi.AuthzGrantResponse
	resp, err = httpx.AsError(c.Do(req))
	if err != nil {
		return errorsx.Wrap(err, "authz grant request failed")
	}

	return httpx.DecodeJSON(resp, &grant)
}
