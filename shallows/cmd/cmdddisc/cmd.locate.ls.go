package cmdddisc

import (
	"fmt"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

type cmdLocateList struct {
	ID          []string `flag:"" name:"id" help:"only show entries matching the given id(s)"`
	Query       string   `flag:"" name:"query" help:"only show entries whose query matches this text"`
	Pending     bool     `flag:"" name:"pending" help:"only show entries the daemon will attempt on its next pass"`
	Completed   bool     `flag:"" name:"completed" help:"only show entries that have been completed"`
	Offset      uint64   `flag:"" name:"offset" help:"page offset for pagination (multiplied by the result limit)"`
	MinAttempts uint64   `flag:"" name:"min-attempts" help:"only show entries with at least this many attempts"`
	MaxAttempts uint64   `flag:"" name:"max-attempts" help:"only show entries with at most this many attempts"`
}

func (t cmdLocateList) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID, daemon *cmdopts.Endpoint) (err error) {
	signer, err := id.Signer()
	if err != nil {
		return errorsx.Wrap(err, "failed to create signer")
	}

	c := authn.AutoOauth2Client(gctx.Context, tls.Config(), authn.EndpointSSHAuth(daemon.Endpoint), authn.SSHTokenSourceOptionSigner(signer))
	cc := authn.AuthzClientLibrary(tls.Config(), c, daemon.Endpoint)

	result, err := ddiscapi.LocateSearch(gctx.Context, cc, daemon.Endpoint, &ddiscapi.LocateSearchRequest{
		Id:          t.ID,
		Query:       t.Query,
		Pending:     t.Pending,
		Completed:   t.Completed,
		AttemptsMin: t.MinAttempts,
		AttemptsMax: t.MaxAttempts,
		Offset:      t.Offset,
		Limit:       100,
	})
	if err != nil {
		return err
	}

	for _, l := range result.Items {
		if _, err = fmt.Printf("id=%s query=%q mimetype=%s attempts=%d next_check=%s located=%s tombstoned=%s\n", l.Id, l.Query, l.Mimetype, l.Attempts, l.NextCheckAt, l.LocatedTorrentId, l.TombstonedAt); err != nil {
			return err
		}
	}

	return nil
}
