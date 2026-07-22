package cmdmeta

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
)

type U12TLs struct {
	Pending  bool   `flag:"" name:"pending"  help:"include pending profiles"`
	Enabled  bool   `flag:"" name:"enabled"  help:"include enabled profiles"`
	Disabled bool   `flag:"" name:"disabled" help:"include disabled profiles"`
	Query    string `flag:"" name:"query"    help:"lucene text search (default field: display)" default:""`
}

func (t U12TLs) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID, daemon *cmdopts.Endpoint) (err error) {
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

func (t U12TLs) statuses() []uint32 {
	var s []uint32
	if t.Disabled {
		s = append(s, uint32(metaapi.ProfileStatus_DISABLED))
	}
	if t.Pending {
		s = append(s, uint32(metaapi.ProfileStatus_PENDING))
	}
	if t.Enabled {
		s = append(s, uint32(metaapi.ProfileStatus_ENABLED))
	}
	if len(s) == 0 {
		s = []uint32{uint32(metaapi.ProfileStatus_NONE)}
	}
	return s
}

func (t U12TLs) run(ctx context.Context, endpoint string, c *http.Client) (err error) {
	for _, status := range t.statuses() {
		if err = t.search(ctx, endpoint, c, status); err != nil {
			return err
		}
	}
	return nil
}

func (t U12TLs) search(ctx context.Context, endpoint string, c *http.Client, status uint32) (err error) {
	req := metaapi.ProfileSearchRequest{
		Query:  t.Query,
		Status: status,
		Limit:  100,
	}

	encoded, err := formx.NewEncoder().Encode(&req)
	if err != nil {
		return errorsx.Wrap(err, "unable to encode request")
	}

	hreq, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/meta/u12t/?"+encoded.Encode(), endpoint), nil)
	if err != nil {
		return errorsx.Wrap(err, "unable to create request")
	}

	resp, err := httpx.AsError(c.Do(hreq))
	if err != nil {
		return errorsx.Wrap(err, "request failed")
	}

	var result metaapi.ProfileSearchResponse
	if err = httpx.DecodeJSON(resp, &result); err != nil {
		return err
	}

	for _, p := range result.Items {
		line := fmt.Sprintf("id=%s display=%s", p.GetId(), p.GetDisplay())
		if label := u12tStatusLabel(status); label != "" {
			line += " status=" + label
		}
		if _, err = fmt.Println(line); err != nil {
			return err
		}
	}

	return nil
}

func u12tStatusLabel(status uint32) string {
	switch status {
	case uint32(metaapi.ProfileStatus_DISABLED):
		return "disabled"
	case uint32(metaapi.ProfileStatus_PENDING):
		return "pending"
	case uint32(metaapi.ProfileStatus_ENABLED):
		return "enabled"
	default:
		return ""
	}
}
