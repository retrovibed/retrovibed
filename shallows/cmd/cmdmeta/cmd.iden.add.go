package cmdmeta

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"golang.org/x/crypto/ssh"
)

type IdenAdd struct {
	Endpoint  string `flag:"" name:"endpoint" help:"http address of the retrovibed instance" default:"localhost:9998"`
	PublicKey string `arg:"" name:"pubkey" help:"public key to add" required:"true"`
}

func (t IdenAdd) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID) (err error) {
	ctx, done := context.WithTimeout(gctx.Context, 10*time.Second)
	defer done()

	signer, err := id.Signer()
	if err != nil {
		return errorsx.Wrap(err, "failed to create signer")
	}

	c := authn.AutoOauth2Client(gctx.Context, tls.Config(), authn.EndpointSSHAuth(fmt.Sprintf("https://%s", t.Endpoint)), authn.SSHTokenSourceOptionSigner(signer))
	cc := metaapi.AuthzClientLibrary(tls.Config(), c, t.Endpoint)

	return t.run(ctx, cc)
}

func (t IdenAdd) run(ctx context.Context, c *http.Client) (err error) {
	pubkey, comment, _, _, err := ssh.ParseAuthorizedKey([]byte(t.PublicKey))
	if err != nil {
		return errorsx.Wrap(err, "invalid public key")
	}

	var result metaapi.ProfileCreateResponse

	encoded, err := json.Marshal(&metaapi.ProfileCreateRequest{
		Profile:   &metaapi.Profile{Display: comment},
		PublicKey: strings.TrimSpace(string(ssh.MarshalAuthorizedKey(pubkey))),
	})
	if err != nil {
		return errorsx.Wrap(err, "unable to encode request")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("https://%s/meta/u12t/", t.Endpoint), bytes.NewReader(encoded))
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

	encoded, err = json.Marshal(&metaapi.AuthzGrantRequest{
		Token: &metaapi.Token{
			LibraryRead: true,
		},
	})
	if err != nil {
		return errorsx.Wrap(err, "unable to encode authz request")
	}

	req, err = http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("https://%s/meta/authz/%s", t.Endpoint, result.Profile.GetId()), bytes.NewReader(encoded))
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
