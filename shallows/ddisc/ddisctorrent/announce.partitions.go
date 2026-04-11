package ddisctorrent

import (
	"context"
	"log"
	"net/netip"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/internal/backoffx"
	"github.com/retrovibed/retrovibed/internal/env"
	"github.com/retrovibed/retrovibed/internal/envx"
	"github.com/retrovibed/retrovibed/internal/errorsx"
)

// announce the provided partitions, and identify any useful peers for syncing.
func Announce(ctx context.Context, cln *torrent.Client, s *dht.Server, partitions ...uuid.UUID) (err error) {
	log.Println("----------------------------------------- announce initiated", len(partitions))
	defer log.Println("----------------------------------------- announce completed", len(partitions))
	for {
		for _, p := range partitions {
			id := int160.New(p.Bytes())

			performsync := func(self, to netip.AddrPort) error {
				var (
					pid int160.T
				)

				if self.Compare(to) == 0 {
					log.Println("skipping self", self, "vs", to)
					return nil
				}

				// ping the node to determine its id, allowing us to check if we've talked to them before.
				if ret := dht.PingDuration(ctx, 2*time.Second, s, to, s.ID()); ret.Err != nil {
					return errorsx.Wrap(ret.Err, "ping failed")
				} else if pid = ret.Reply.SenderID().Int160(); pid == s.ID() {
					log.Println("skipping self", s.ID(), "vs", pid)
					return nil
				} else {
					log.Println("pinged", pid, to, "we are", ret.Reply.IP)
				}

				dctx, done := context.WithTimeout(ctx, time.Second)
				defer done()

				req, err := NewSyncRequest(s.ID(), p.String(), uuid.Nil.String())
				if err != nil {
					return errorsx.Wrap(err, "failed to create sync request")
				}

				ret := s.Query(dctx, dht.NewAddr(to), req)
				if ret.Err != nil {
					return ret.Err
				}

				return nil
			}

			seq, err := torrent.DHTAnnounce(ctx, s, id)
			if err != nil {
				log.Println("failed to announce partition", p, err)
				continue
			}

			for pv := range seq.Each(ctx) {
				addrport := s.AddrPort()
				for _, peer := range pv.Peers {
					log.Println("sync initiated", id, peer)
					if err := performsync(addrport, peer.AddrPort); err != nil {
						log.Println("failed to request sync with", id, peer, err)
						continue
					}
					log.Println("sync completed", id, peer)
				}
			}

			if err := seq.Err(); err != nil {
				log.Println("failed partition sync", err)
				continue
			}
		}

		select {
		case <-time.After(envx.Duration(time.Hour, env.DDiscFrequency) + backoffx.Random(10*time.Second)):
		case <-ctx.Done():
			return context.Cause(ctx)
		}
	}
}
