package metaapi

import (
	"context"

	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
)

// QueryDiscoveryDiagnostics reports the state of the infohash identification pipeline.
func QueryDiscoveryDiagnostics(ctx context.Context, q sqlx.Queryer) (d *DiscoveryDiagnostics, err error) {
	var (
		queued, indexing, offload, indexed int
	)

	if queued, err = sqlx.Count(ctx, q, "SELECT COUNT (*) FROM torrents_unknown_infohashes WHERE next_check < NOW()"); err != nil {
		return d, err
	}

	if indexing, err = sqlx.Count(ctx, q, "SELECT COUNT (*) FROM ddisc_media WHERE known_media_id = 'ffffffff-ffff-ffff-ffff-ffffffffffff'"); err != nil {
		return d, err
	}

	if offload, err = sqlx.Count(ctx, q, "SELECT COUNT (*) FROM ddisc_media WHERE known_media_id = '00000000-0000-0000-0000-000000000000'"); err != nil {
		return d, err
	}

	if indexed, err = sqlx.Count(ctx, q, "SELECT COUNT (*) FROM ddisc_media WHERE known_media_id NOT IN ('00000000-0000-0000-0000-000000000000', 'ffffffff-ffff-ffff-ffff-ffffffffffff')"); err != nil {
		return d, err
	}

	return &DiscoveryDiagnostics{
		Queued:   uint64(queued),
		Indexing: uint64(indexing),
		Offload:  uint64(offload),
		Indexed:  uint64(indexed),
	}, nil
}
