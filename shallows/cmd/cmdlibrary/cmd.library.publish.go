package cmdlibrary

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/communityapi"
	"github.com/retrovibed/retrovibed/shallows/internal/debugx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
)

type cmdPublish struct {
	DryRun bool `flag:"" name:"dry-run" help:"print what would be published without actually publishing" negatable:"" default:"true"`
}

func (t cmdPublish) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, daemon *cmdopts.Endpoint) (err error) {
	httpc := authn.AutoOauth2Client(gctx.Context, tls.Config(), authn.EndpointSSHAuth(daemon.Endpoint))
	return t.run(gctx.Context, daemon.Endpoint, jsonl.NewEncoder(os.Stdout), os.Stdin, httpc)
}

func (t cmdPublish) run(ctx context.Context, endpoint string, enc *jsonl.Encoder, r io.Reader, c *http.Client) error {
	var com communityapi.Community

	debugx.Println("reading community from stdin")
	d := jsonl.NewDecoder(r)
	if err := d.Decode(&com); err != nil {
		return errorsx.Wrap(err, "unable to decode community from stdin")
	}

	var derr error
	var lmd library.Metadata
	for derr = d.Decode(&lmd); derr == nil; derr = d.Decode(&lmd) {
		if t.DryRun {
			if err := enc.Encode(langx.Clone(lmd, timex.JSONSafeEncodeOption)); err != nil {
				return err
			}
			continue
		}

		resp, err := t.publishItem(ctx, endpoint, c, &com, lmd.ID)
		if err != nil {
			return errorsx.Wrapf(err, "failed to publish library content: %s", lmd.ID)
		}

		if err = enc.Encode(resp.PublishedContent); err != nil {
			return errorsx.Wrap(err, "unable to write published content")
		}
	}

	return errorsx.Ignore(derr, io.EOF)
}

func (t cmdPublish) publishItem(ctx context.Context, endpoint string, c *http.Client, com *communityapi.Community, libraryID string) (*communityapi.PublishContentResponse, error) {
	var (
		err  error
		req  *http.Request
		resp *http.Response
		msg  communityapi.PublishContentResponse
	)

	body, err := jsonx.Marshal(&communityapi.PublishContentRequest{
		PublishedContent: &communityapi.PublishedContent{
			LibraryId:   libraryID,
			CommunityId: com.Id,
		},
		PublishMode: com.DefaultPublishMode,
	})

	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("%s/c/p/%s", endpoint, com.Id)
	if req, err = http.NewRequestWithContext(ctx, http.MethodPost, uri, bytes.NewReader(body)); err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", mimex.JSON)

	if resp, err = httpx.AsError(c.Do(req)); err != nil {
		return nil, errorsx.Wrap(err, "failed to publish library content")
	}
	defer resp.Body.Close()

	if err = jsonx.UnmarshalRead(resp.Body, &msg); err != nil {
		return nil, errorsx.Wrap(err, "failed to decode response")
	}

	return &msg, nil
}
