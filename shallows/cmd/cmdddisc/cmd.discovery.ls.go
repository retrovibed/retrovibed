package cmdddisc

import (
	"fmt"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/meta"
)

type cmdDiscoveryList struct {
	Endpoint    string        `flag:"" name:"library" help:"http address for the library you want to connect to" default:"localhost:9998"`
	NextCheck   time.Duration `flag:"" name:"next-check" help:"only show entries due within the the next window" default:"0m"`
	ID          []string      `flag:"" name:"id" help:"only show entries matching the given id(s)"`
	Offset      uint64        `flag:"" name:"offset" help:"page offset for pagination (multiplied by the result limit)"`
	MinAttempts uint64        `flag:"" name:"min-attempts" help:"only show entries with at least this many attempts"`
	MaxAttempts uint64        `flag:"" name:"max-attempts" help:"only show entries with at most this many attempts"`
}

func (t cmdDiscoveryList) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID) (err error) {
	signer, err := id.Signer()
	if err != nil {
		return errorsx.Wrap(err, "failed to create signer")
	}

	c := authn.AutoOauth2Client(gctx.Context, tls.Config(), authn.EndpointSSHAuth(fmt.Sprintf("https://%s", t.Endpoint)), authn.SSHTokenSourceOptionSigner(signer))
	cc := authn.AuthzClientLibrary(tls.Config(), c, t.Endpoint)

	result, err := ddiscapi.DiscoverySearch(gctx.Context, cc, t.Endpoint, &ddiscapi.DiscoverySearchRequest{
		NextCheck:   meta.NewDateRange(timex.NewRangeWithin(t.NextCheck)),
		Id:          t.ID,
		Offset:      t.Offset,
		AttemptsMin: t.MinAttempts,
		AttemptsMax: t.MaxAttempts,
		Limit:       100,
	})
	if err != nil {
		return err
	}

	for _, d := range result.Items {
		if _, err = fmt.Printf("id=%s infohash=%x attempts=%s next_check=%s\n", d.Id, d.Infohash, d.Attempts, d.NextCheck); err != nil {
			return err
		}
	}

	return nil
}
