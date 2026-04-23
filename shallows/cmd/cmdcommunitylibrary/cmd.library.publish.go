package cmdcommunitylibrary

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/debugx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/meta"
)

type cmdPublish struct {
	Endpoint string `flag:"" name:"peer" help:"http address for the retrovibed daemon" default:"http://localhost:9998"`
	DryRun   bool   `flag:"" name:"dry-run" help:"print what would be published without actually publishing" negatable:"" default:"true"`
}

func (t cmdPublish) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig) (err error) {
	httpc := authn.AutoOauth2Client(gctx.Context, tls.Config(), authn.EndpointSSHAuth(t.Endpoint))

	debugx.Println("reading community from stdin")
	var com meta.Community
	if err = json.NewDecoder(os.Stdin).Decode(&com); err != nil {
		return err
	}

	return t.run(gctx.Context, jsonl.NewEncoder(os.Stdout), os.Stdin, httpc, &com)
}

func (t cmdPublish) run(ctx context.Context, enc *jsonl.Encoder, r io.Reader, c *http.Client, com *meta.Community) error {
	d := jsonl.NewDecoder(r)

	var derr error
	var lmd library.Metadata
	for derr = d.Decode(&lmd); derr == nil; derr = d.Decode(&lmd) {
		if t.DryRun {
			if err := enc.Encode(langx.Clone(lmd, timex.JSONSafeEncodeOption)); err != nil {
				return err
			}
			continue
		}

		resp, err := t.publishItem(ctx, c, com, lmd.ID)
		if err != nil {
			return err
		}

		if err = enc.Encode(resp.PublishedContent); err != nil {
			return errorsx.Wrap(err, "unable to write published content")
		}
	}

	return errorsx.Ignore(derr, io.EOF)
}

func (t cmdPublish) publishItem(ctx context.Context, c *http.Client, com *meta.Community, libraryID string) (*meta.PublishContentResponse, error) {
	var (
		err  error
		req  *http.Request
		resp *http.Response
		msg  meta.PublishContentResponse
	)

	body, err := json.Marshal(&meta.PublishContentRequest{
		PublishedContent: &meta.PublishedContent{LibraryId: libraryID},
		PublishMode:      com.DefaultPublishMode,
	})
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("%s/c/%s/publish", t.Endpoint, com.Id)
	if req, err = http.NewRequestWithContext(ctx, http.MethodPost, uri, bytes.NewReader(body)); err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", mimex.JSON)

	if resp, err = httpx.AsError(c.Do(req)); err != nil {
		return nil, errorsx.Wrap(err, "failed to publish library content")
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		return nil, errorsx.Wrap(err, "failed to decode response")
	}

	return &msg, nil
}
