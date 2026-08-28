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
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
)

type DeviceEdit struct {
	DeviceID    string  `arg:"" name:"device.id" help:"device id to edit" required:"true"`
	Default     bool    `flag:"" name:"default" help:"mark this device as the default"`
	Downloads   bool    `flag:"" name:"download" help:"mark this device as the download target"`
	Description *string `flag:"" name:"description" help:"update the description"`
}

func (t DeviceEdit) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID, daemon *cmdopts.Endpoint) (err error) {
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

func (t DeviceEdit) run(ctx context.Context, endpoint string, c *http.Client) (err error) {
	var current *metaapi.Daemon

	if t.Default {
		if current, err = t.touch(ctx, endpoint, c); err != nil {
			return err
		}
	}

	if t.Downloads || t.Description != nil {
		if current, err = t.patch(ctx, endpoint, c); err != nil {
			return err
		}
	}

	if current == nil {
		if current, err = t.lookup(ctx, endpoint, c); err != nil {
			return err
		}
	}

	_, err = fmt.Printf("id=%s hostname=%s description=%s default=%t downloads=%t\n", current.Id, current.Hostname, current.Description, current.Default, current.Downloads)
	return err
}

// touch marks this device as the default via PUT /meta/d/{id}.
func (t DeviceEdit) touch(ctx context.Context, endpoint string, c *http.Client) (*metaapi.Daemon, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, fmt.Sprintf("%s/meta/d/%s", endpoint, t.DeviceID), nil)
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to create request")
	}

	resp, err := httpx.AsError(c.Do(req))
	if err != nil {
		return nil, errorsx.Wrap(err, "mark default failed")
	}

	var result metaapi.DaemonLookupResponse
	if err = httpx.DecodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return result.Daemon, nil
}

// patch applies downloads/description changes via POST /meta/d/{id}, which
// replaces the full record — so the existing hostname/description/default
// are fetched first and carried forward for whichever fields weren't
// overridden (same fetch-then-patch idiom as cmd.u12t.update.go).
func (t DeviceEdit) patch(ctx context.Context, endpoint string, c *http.Client) (*metaapi.Daemon, error) {
	existing, err := t.lookup(ctx, endpoint, c)
	if err != nil {
		return nil, err
	}

	patched := &metaapi.Daemon{
		Hostname:    existing.Hostname,
		Description: *langx.FirstNonZero(t.Description, &existing.Description),
		Default:     existing.Default,
		Downloads:   existing.Downloads || t.Downloads,
	}

	encoded, err := json.Marshal(&metaapi.DaemonUpdateRequest{Daemon: patched})
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to encode request")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/meta/d/%s", endpoint, t.DeviceID), bytes.NewReader(encoded))
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to create request")
	}

	resp, err := httpx.AsError(c.Do(req))
	if err != nil {
		return nil, errorsx.Wrap(err, "edit device failed")
	}

	var result metaapi.DaemonUpdateResponse
	if err = httpx.DecodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return result.Daemon, nil
}

// lookup finds the device by id from the first page of known devices — there
// is no GET-by-id route for daemons (only /latest and the paginated /).
func (t DeviceEdit) lookup(ctx context.Context, endpoint string, c *http.Client) (*metaapi.Daemon, error) {
	encoded, err := formx.NewEncoder().Encode(&metaapi.DaemonSearchRequest{Limit: 128})
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to encode request")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/meta/d/?"+encoded.Encode(), endpoint), nil)
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to create request")
	}

	resp, err := httpx.AsError(c.Do(req))
	if err != nil {
		return nil, errorsx.Wrap(err, "lookup failed")
	}

	var result metaapi.DaemonSearchResponse
	if err = httpx.DecodeJSON(resp, &result); err != nil {
		return nil, err
	}

	for _, d := range result.Items {
		if d.Id == t.DeviceID {
			return d, nil
		}
	}

	return nil, errorsx.Errorf("device not found: %s", t.DeviceID)
}
