package cmdmeta

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/grpcx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
)

type Usermanagement struct {
	Ls     U12TLs     `cmd:"" help:"list profiles"`
	Grant  U12TGrant  `cmd:"" help:"grant access to a profile"`
	Revoke U12TRevoke `cmd:"" help:"revoke access from a profile"`
}

type U12TLs struct {
	Endpoint string `flag:"" name:"endpoint" help:"http address of the retrovibed instance" default:"localhost:9998"`
	Pending  bool   `flag:"" name:"pending" help:"filter to only show pending profiles"`
}

func (t U12TLs) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID) (err error) {
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

func (t U12TLs) run(ctx context.Context, c *http.Client) (err error) {
	var (
		result metaapi.ProfileSearchResponse
		req    metaapi.ProfileSearchRequest
	)

	if t.Pending {
		req.Status = uint32(metaapi.ProfileStatus_PENDING)
	} else {
		req.Status = uint32(metaapi.ProfileStatus_DISABLED)
	}
	req.Limit = 100

	encoded, err := formx.NewEncoder().Encode(&req)
	if err != nil {
		return errorsx.Wrap(err, "unable to encode request")
	}

	hreq, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://%s/meta/u12t/?"+encoded.Encode(), t.Endpoint), nil)
	if err != nil {
		return errorsx.Wrap(err, "unable to create request")
	}

	resp, err := httpx.AsError(c.Do(hreq))
	if err != nil {
		return errorsx.Wrap(err, "request failed")
	}

	if err = httpx.DecodeJSON(resp, &result); err != nil {
		return err
	}

	for _, p := range result.Items {
		_, err := fmt.Printf("id=%s display=%s\n", p.GetId(), p.GetDisplay())
		if err != nil {
			return err
		}
	}

	return nil
}

type U12TGrant struct {
	Endpoint  string `flag:"" name:"endpoint" help:"http address of the retrovibed instance" default:"localhost:9998"`
	ProfileID string `arg:"" name:"profile.id" help:"profile id to grant access" required:"true"`
}

func (t U12TGrant) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID) (err error) {
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

func (t U12TGrant) run(ctx context.Context, c *http.Client) (err error) {
	// fetch the current profile state so we can preserve existing fields in the update
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://%s/meta/u12t/%s", t.Endpoint, t.ProfileID), nil)
	if err != nil {
		return errorsx.Wrap(err, "unable to create profile request")
	}

	var lookup metaapi.ProfileLookupResponse
	resp, err := httpx.AsError(c.Do(req))
	if err != nil {
		return errorsx.Wrap(err, "profile lookup failed")
	}

	if err = httpx.DecodeJSON(resp, &lookup); err != nil {
		return err
	}

	// enable the profile by setting disabled_pending_approval_at to RFC3339 infinity
	lookup.Profile.DisabledPendingApprovalAt = grpcx.EncodeTime(timex.RFC3339NanoEncode(timex.Inf()))
	encoded, err := json.Marshal(&metaapi.ProfileUpdateRequest{Profile: lookup.Profile})
	if err != nil {
		return errorsx.Wrap(err, "unable to encode profile update")
	}

	req, err = http.NewRequestWithContext(ctx, http.MethodPatch, fmt.Sprintf("https://%s/meta/u12t/%s", t.Endpoint, t.ProfileID), bytes.NewReader(encoded))
	if err != nil {
		return errorsx.Wrap(err, "unable to create update request")
	}

	var updated metaapi.ProfileUpdateResponse
	resp, err = httpx.AsError(c.Do(req))
	if err != nil {
		return errorsx.Wrap(err, "profile enable failed")
	}

	if err = httpx.DecodeJSON(resp, &updated); err != nil {
		return err
	}

	// grant library read access
	encoded, err = json.Marshal(&metaapi.AuthzGrantRequest{
		Token: &metaapi.Token{
			LibraryRead: true,
		},
	})
	if err != nil {
		return errorsx.Wrap(err, "unable to encode authz request")
	}

	req, err = http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("https://%s/meta/authz/%s", t.Endpoint, t.ProfileID), bytes.NewReader(encoded))
	if err != nil {
		return errorsx.Wrap(err, "unable to create authz request")
	}

	var grant metaapi.AuthzGrantResponse
	resp, err = httpx.AsError(c.Do(req))
	if err != nil {
		return errorsx.Wrap(err, "authz grant failed")
	}

	return httpx.DecodeJSON(resp, &grant)
}

type U12TRevoke struct {
	Endpoint  string `flag:"" name:"endpoint" help:"http address of the retrovibed instance" default:"localhost:9998"`
	ProfileID string `arg:"" name:"profile.id" help:"profile id to revoke access" required:"true"`
}

func (t U12TRevoke) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID) (err error) {
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

func (t U12TRevoke) run(ctx context.Context, c *http.Client) (err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("https://%s/meta/authz/%s", t.Endpoint, t.ProfileID), nil)
	if err != nil {
		return errorsx.Wrap(err, "unable to create request")
	}

	var revoke metaapi.AuthzRevokeResponse
	resp, err := httpx.AsError(c.Do(req))
	if err != nil {
		return errorsx.Wrap(err, "authz revoke failed")
	}

	return httpx.DecodeJSON(resp, &revoke)
}
