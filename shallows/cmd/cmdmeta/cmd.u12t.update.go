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
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
)

type U12TUpdate struct {
	Display   *string `flag:"" name:"display" help:"set the profile's display name"`
	ProfileID string  `arg:"" name:"profile.id" help:"profile id to update" required:"true"`
}

func (t U12TUpdate) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID, daemon *cmdopts.Endpoint) (err error) {
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

func (t U12TUpdate) run(ctx context.Context, endpoint string, c *http.Client) (err error) {
	// fetch the current profile state so we can preserve fields the caller didn't override
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

	lookup.Profile.Display = *langx.FirstNonZero(t.Display, &lookup.Profile.Display)

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
		return errorsx.Wrap(err, "profile update failed")
	}

	if err = httpx.DecodeJSON(resp, &updated); err != nil {
		return err
	}

	_, err = fmt.Printf("id=%s display=%s\n", updated.Profile.Id, updated.Profile.Display)
	return err
}
