package cmdddisc

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/retroapi/authn"
	ddiscapiimport "github.com/retrovibed/retrovibed/retroapi/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/internal/asynccompute"
	"github.com/retrovibed/retrovibed/shallows/internal/backoffx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
)

type cmdDiscoveryImportJSONL struct {
	Endpoint string `flag:"" name:"library" help:"http address for the library you want to connect to" default:"localhost:9998"`
	Backlog  uint16 `flag:"" name:"backlog" help:"number of batches to allowed to queue up" default:"128"`
	Workers  uint16 `flag:"" name:"workers" help:"number of async database workers to run" default:"1"`
}

func (t cmdDiscoveryImportJSONL) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID) (err error) {
	c := authn.AutoOauth2Client(gctx.Context, tls.Config(), authn.EndpointSSHAuth(fmt.Sprintf("https://%s", t.Endpoint)))
	return t.run(gctx.Context, c, t.Endpoint, os.Stdin)
}

func (t cmdDiscoveryImportJSONL) run(ctx context.Context, c *http.Client, endpoint string, r io.Reader) error {
	inserts := asynccompute.New(func(ctx context.Context, v *ddiscapiimport.Import) (err error) {
		bs := backoffx.New(
			backoffx.Constant(time.Second),
			backoffx.JitterRandom(200*time.Millisecond),
		)
		_, err = backoffx.AttemptV(ctx, bs, func(ctx context.Context, attempt uint) (struct{}, error) {
			if attempt >= 5 {
				return struct{}{}, backoffx.ErrStopAttempts
			}

			m, err := metainfo.ParseMagnetURI(v.Magnet)
			if err != nil {
				return struct{}{}, errorsx.Wrap(err, "failed to parse magnet uri")
			}

			rsp, err := ddiscapi.DiscoveryCreate(ctx, c, endpoint, &ddiscapi.DiscoveryCreateRequest{
				Discovery: &ddiscapi.Discovery{
					Infohash: m.InfoHash.Bytes(),
				},
			})
			if err != nil {
				return struct{}{}, errorsx.Wrap(err, "failed to import magnet uri")
			}

			log.Println("magnet uri import", spew.Sdump(rsp.Discovery))
			return struct{}{}, nil
		})
		return err
	}, asynccompute.Backlog[*ddiscapiimport.Import](t.Backlog), asynccompute.Workers[*ddiscapiimport.Import](t.Workers))

	d := jsonl.Iter[*ddiscapiimport.Import](jsonl.NewDecoder(r))

	for v := range d.Each(ctx) {
		if err := inserts.Run(ctx, v); err != nil {
			return errorsx.Compact(err, asynccompute.Shutdown(ctx, inserts))
		}
	}

	if err := errorsx.Compact(d.Err(), asynccompute.Shutdown(ctx, inserts)); err != nil {
		return err
	}

	return nil
}
