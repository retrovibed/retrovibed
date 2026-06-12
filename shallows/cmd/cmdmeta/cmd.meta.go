package cmdmeta

import (
	"context"
	"strings"

	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/meta"
)

type Cloud struct {
	Register CloudRegister `cmd:"" help:"register the current identity with the cloud service"`
}

func Hostnames(ctx context.Context, q sqlx.Queryer) (res []string, _ error) {
	s := sqlx.Scan(meta.DaemonSearch(ctx, q, meta.DaemonSearchBuilder().Limit(128)))
	for d := range s.Iter() {
		before, _, _ := strings.Cut(d.Hostname, ":")
		res = append(res, before)
	}

	return res, s.Err()
}
