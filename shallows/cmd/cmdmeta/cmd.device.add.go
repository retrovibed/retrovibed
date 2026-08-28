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
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"golang.org/x/crypto/ssh"
)

type DeviceAdd struct {
	Hostname    string `arg:"" name:"hostname" help:"hostname[:port] of the device to add" required:"true"`
	Description string `flag:"" name:"description" help:"human readable description"`
	Default     bool   `flag:"" name:"default" help:"mark this device as the default"`
	Downloads   bool   `flag:"" name:"download" help:"mark this device as the download target"`
	Force       bool   `flag:"" name:"force" help:"skip the reachability check"`
}

func (t DeviceAdd) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID, daemon *cmdopts.Endpoint) (err error) {
	ctx, done := context.WithTimeout(gctx.Context, 10*time.Second)
	defer done()

	signer, err := id.Signer()
	if err != nil {
		return errorsx.Wrap(err, "failed to create signer")
	}

	if !t.Force {
		if err = t.checkReachable(ctx, tls, signer); err != nil {
			return errorsx.Wrap(err, "device unreachable, use --force to add it anyway")
		}
	}

	c := authn.AutoOauth2Client(gctx.Context, tls.Config(), authn.EndpointSSHAuth(daemon.Endpoint), authn.SSHTokenSourceOptionSigner(signer))
	cc := authn.AuthzClientLibrary(tls.Config(), c, daemon.Endpoint)

	return t.run(ctx, daemon.Endpoint, cc)
}

// checkReachable mirrors console/lib/meta/api.dart's `connectable`: confirm
// the target is up, then confirm our identity is actually authorized there
// (not just that the host responds).
func (t DeviceAdd) checkReachable(ctx context.Context, tls *cmdopts.TLSConfig, signer ssh.Signer) (err error) {
	target := fmt.Sprintf("https://%s", t.Hostname)

	hc := &http.Client{Transport: &http.Transport{TLSClientConfig: tls.Config()}, Timeout: 5 * time.Second}
	hreq, err := http.NewRequestWithContext(ctx, http.MethodGet, target+"/healthz", nil)
	if err != nil {
		return errorsx.Wrap(err, "unable to create healthz request")
	}

	if _, err = httpx.AsError(hc.Do(hreq)); err != nil {
		return errorsx.Wrap(err, "healthz check failed")
	}

	ac := authn.AutoOauth2Client(ctx, tls.Config(), authn.EndpointSSHAuth(target), authn.SSHTokenSourceOptionSigner(signer))
	areq, err := http.NewRequestWithContext(ctx, http.MethodGet, target+"/meta/authz/", nil)
	if err != nil {
		return errorsx.Wrap(err, "unable to create authz request")
	}

	if _, err = httpx.AsError(ac.Do(areq)); err != nil {
		return errorsx.Wrap(err, "authz check failed")
	}

	return nil
}

func (t DeviceAdd) run(ctx context.Context, endpoint string, c *http.Client) (err error) {
	encoded, err := json.Marshal(&metaapi.DaemonCreateRequest{
		Daemon: &metaapi.Daemon{
			Hostname:    t.Hostname,
			Description: t.Description,
			Default:     t.Default,
			Downloads:   t.Downloads,
		},
	})
	if err != nil {
		return errorsx.Wrap(err, "unable to encode request")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/meta/d/", endpoint), bytes.NewReader(encoded))
	if err != nil {
		return errorsx.Wrap(err, "unable to create request")
	}

	resp, err := httpx.AsError(c.Do(req))
	if err != nil {
		return errorsx.Wrap(err, "add device failed")
	}

	var result metaapi.DaemonCreateResponse
	if err = httpx.DecodeJSON(resp, &result); err != nil {
		return err
	}

	_, err = fmt.Printf("id=%s hostname=%s description=%s\n", result.Daemon.Id, result.Daemon.Hostname, result.Daemon.Description)
	return err
}
