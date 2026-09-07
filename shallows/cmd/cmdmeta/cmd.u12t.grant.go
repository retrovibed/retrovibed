package cmdmeta

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/grpcx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
)

type U12TGrant struct {
	LibraryRead     bool   `flag:"" name:"library-read"     help:"grant library read access"               negatable:"" default:"true"`
	LibraryModify   bool   `flag:"" name:"library-modify"   help:"grant library modify access"             negatable:"" default:"false"`
	RemoteControl   bool   `flag:"" name:"remote-control"   help:"grant remote control access"             negatable:"" default:"false"`
	BillingRead     bool   `flag:"" name:"billing-read"     help:"grant billing read access"               negatable:"" default:"false"`
	BillingModify   bool   `flag:"" name:"billing-modify"   help:"grant billing modify access"             negatable:"" default:"false"`
	CommunityModify bool   `flag:"" name:"community-modify" help:"grant community modify access"           negatable:"" default:"false"`
	Usermanagement  bool   `flag:"" name:"usermanagement"   help:"grant usermanagement access"             negatable:"" default:"false"`
	ProfileID       string `arg:"" name:"profile.id"        help:"profile id to grant access"             required:"true"`
}

func (t U12TGrant) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID, daemon *cmdopts.Endpoint) (err error) {
	ctx, done := context.WithTimeout(gctx.Context, 10*time.Second)
	defer done()

	signer, err := id.Signer()
	if err != nil {
		return errorsx.Wrap(err, "failed to create signer")
	}

	c := authn.AutoOauth2Client(gctx.Context, tls.Config(), authn.EndpointSSHAuth(daemon.Endpoint), authn.SSHTokenSourceOptionSigner(signer))
	cc := authn.AuthzClientLibrary(tls.Config(), c, daemon.Endpoint)

	return t.run(ctx, daemon.Endpoint, cc)
}

func (t U12TGrant) run(ctx context.Context, endpoint string, c *http.Client) (err error) {
	// fetch the current profile state so we can preserve existing fields in the update
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/meta/u12t/%s", endpoint, t.ProfileID), nil)
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
	encoded, err := jsonx.Marshal(&metaapi.ProfileUpdateRequest{Profile: lookup.Profile})
	if err != nil {
		return errorsx.Wrap(err, "unable to encode profile update")
	}

	req, err = http.NewRequestWithContext(ctx, http.MethodPatch, fmt.Sprintf("%s/meta/u12t/%s", endpoint, t.ProfileID), bytes.NewReader(encoded))
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

	encoded, err = jsonx.Marshal(&metaapi.AuthzGrantRequest{
		Token: &metaapi.Token{
			LibraryRead:     t.LibraryRead,
			LibraryModify:   t.LibraryModify,
			RemoteControl:   t.RemoteControl,
			BillingRead:     t.BillingRead,
			BillingModify:   t.BillingModify,
			CommunityModify: t.CommunityModify,
			Usermanagement:  t.Usermanagement,
		},
	})

	if err != nil {
		return errorsx.Wrap(err, "unable to encode authz request")
	}

	req, err = http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/meta/authz/%s", endpoint, t.ProfileID), bytes.NewReader(encoded))
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
