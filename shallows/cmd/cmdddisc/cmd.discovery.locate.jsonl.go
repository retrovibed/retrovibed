package cmdddisc

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/asynccompute"
	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/internal/backoffx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/library"
)

type cmdMediaLocateJSONL struct {
	Endpoint string `flag:"" name:"library" help:"http address for the library you want to connect to" default:"localhost:9998"`
	Backlog  uint16 `flag:"" name:"backlog" help:"number of batches to allowed to queue up" default:"128"`
	Workers  uint16 `flag:"" name:"workers" help:"number of async database workers to run" default:"1"`
}

func (t cmdMediaLocateJSONL) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID) (err error) {
	signer, err := id.Signer()
	if err != nil {
		return errorsx.Wrap(err, "failed to create signer")
	}

	c := authn.AutoOauth2Client(gctx.Context, tls.Config(), authn.EndpointSSHAuth(fmt.Sprintf("https://%s", t.Endpoint)), authn.SSHTokenSourceOptionSigner(signer))
	cc := authn.AuthzClientLibrary(tls.Config(), c, t.Endpoint)

	return t.run(gctx.Context, cc, t.Endpoint, os.Stdin)
}

func (t cmdMediaLocateJSONL) run(ctx context.Context, c *http.Client, endpoint string, r io.Reader) error {
	inserts := asynccompute.New(func(ctx context.Context, v library.Known) (err error) {
		if stringsx.Blank(v.Mimetype) {
			log.Println("skipping locate request, missing mimetype", v.UID, v.Title)
			return nil
		}

		bs := backoffx.New(
			backoffx.Constant(time.Second),
			backoffx.JitterRandom(200*time.Millisecond),
		)
		_, err = backoffx.AttemptV(ctx, bs, func(ctx context.Context, attempt uint) (struct{}, error) {
			if attempt >= 5 {
				return struct{}{}, backoffx.ErrStopAttempts
			}

			rsp, err := ddiscapi.LocateCreate(ctx, c, endpoint, &ddiscapi.LocateCreateRequest{
				Locate: &ddiscapi.Locate{
					Query:        v.Title,
					Mimetype:     v.Mimetype,
					KnownMediaId: v.UID,
				},
			})
			if err != nil {
				return struct{}{}, errorsx.Wrap(err, "failed to submit locate request")
			}

			log.Println("locate request submitted", rsp.Locate.Id, rsp.Locate.Query)
			return struct{}{}, nil
		})
		return err
	}, asynccompute.Backlog[library.Known](t.Backlog), asynccompute.Workers[library.Known](t.Workers))

	d := jsonl.Iter[library.Known](jsonl.NewDecoder(r))

	for v := range d.Each(ctx) {
		if err := inserts.Run(ctx, v); err != nil {
			return errorsx.Compact(err, asynccompute.Shutdown(ctx, inserts))
		}
	}

	return errorsx.Compact(d.Err(), asynccompute.Shutdown(ctx, inserts))
}
