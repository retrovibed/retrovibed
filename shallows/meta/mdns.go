package meta

import (
	"context"
	"iter"
	"time"

	"github.com/hashicorp/mdns"
	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
)

// MDNSService is the service name advertised/queried for retrovibed daemons on the LAN.
const MDNSService = "_retrovibed._udp"

// MDNSLookup performs an mdns query for the given service, returning the discovered entries.
// injectable so tests can supply canned entries instead of hitting the network.
type MDNSLookup func(ctx context.Context, service string, timeout time.Duration) ([]*mdns.ServiceEntry, error)

// DefaultMDNSLookup queries the LAN via multicast dns for the given service.
func DefaultMDNSLookup(ctx context.Context, service string, timeout time.Duration) ([]*mdns.ServiceEntry, error) {
	entries := make(chan *mdns.ServiceEntry, 32)
	params := mdns.DefaultParams(service)
	params.Entries = entries
	params.Timeout = timeout

	if err := mdns.QueryContext(ctx, params); err != nil {
		return nil, err
	}
	close(entries)

	results := make([]*mdns.ServiceEntry, 0, len(entries))
	for e := range entries {
		results = append(results, e)
	}

	return results, nil
}

// DiscoverOnce performs a single mdns scan for retrovibed peers on the LAN, upserting
// each discovered peer into meta_daemons keyed by hostname as it's found and yielding
// it immediately. CreatedAt/UpdatedAt on the yielded Daemon let callers distinguish a
// newly discovered peer (CreatedAt == UpdatedAt) from one that was already known (this
// scan only touched UpdatedAt) — DiscoverOnce itself makes no such distinction.
func DiscoverOnce(q sqlx.Queryer, lookup MDNSLookup) iterx.Seq[Daemon] {
	return &discoverSeq{q: q, lookup: lookup}
}

type discoverSeq struct {
	q      sqlx.Queryer
	lookup MDNSLookup
	err    error
}

func (t *discoverSeq) Each(ctx context.Context) iter.Seq[Daemon] {
	return func(yield func(Daemon) bool) {
		entries, err := t.lookup(ctx, MDNSService, 5*time.Second)
		if err != nil {
			t.err = err
			return
		}

		for _, e := range entries {
			d := Daemon{Hostname: e.Host}
			DaemonOptionMaybeID(&d)
			DaemonOptionEnsureDescription(&d)

			if err := DaemonInsertWithDefaults(ctx, t.q, d).Scan(&d); err != nil {
				t.err = err
				return
			}

			if !yield(d) {
				return
			}
		}
	}
}

func (t *discoverSeq) Err() error {
	return t.err
}
