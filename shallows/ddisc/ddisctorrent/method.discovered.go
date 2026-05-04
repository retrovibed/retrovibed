package ddisctorrent

import (
	"context"

	"github.com/james-lawrence/torrent/bencode"
	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/dht/krpc"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/cryptox"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
)

type DiscoveredRequest struct {
	Q string           `bencode:"q,omitempty"` // ddisc_discovered
	A DiscoveredParams `bencode:"a,omitempty"`
	T string           `bencode:"t"` // required: transaction ID
	Y string           `bencode:"y"` // required: type of the message: q for QUERY, r for RESPONSE, e for ERROR
}

type DiscoveredParams struct {
	Peer   krpc.ID `bencode:"pid"`    // ID of the querying Node
	Token  string  `bencode:"token"`  // required: authorization token for responses.
	Offset string  `bencode:"offset"` // required: sync offset.
	Ratio  uint8   `bencode:"ratio"`  // ratio
}

func NewDiscoveredRequest(from krpc.ID, ratio uint8, offset string) (qi dht.QueryInput, err error) {
	req := DiscoveredRequest{
		Y: "q",
		T: krpc.TimestampTransactionID(),
		Q: MethodDisc,
		A: DiscoveredParams{
			Peer:   from,
			Ratio:  ratio,
			Offset: offset,
		},
	}

	encoded, err := bencode.Marshal(req)
	return dht.NewEncodedRequest(req.Q, req.T, encoded), err
}

func NewDiscovered(q sqlx.Queryer, p *ddisc.Partition) Discovered {
	return Discovered{
		q:          q,
		partitions: p,
	}
}

type Discovered struct {
	q          sqlx.Queryer
	partitions *ddisc.Partition
}

func (t Discovered) Handle(ctx context.Context, source dht.Addr, s *dht.Server, b dht.Binding, raw []byte, _ *krpc.Msg) error {
	var (
		m DiscoveredRequest
	)

	if err := bencode.Unmarshal(raw, &m); err != nil {
		return err
	}

	if _, err := s.SendMessageToNode(ctx, krpc.NewEmptyResponse(m.T), source, 1); err != nil {
		return err
	}

	n := t.partitions.Max(m.A.Peer[:])
	block := ddisc.FilterRatio(cryptox.NewChaCha8(n[:]), m.A.Ratio)

	seq := ddisc.SyncDiscovered(t.q, block, m.A.Offset)
	for v := range seq.Each(ctx) {
		msg := Media{
			Q: MethodMedia,
			Y: "q",
			T: krpc.TimestampTransactionID(),
			A: mediaFromDiscovered(m.A.Token, b.ID(), &v),
		}

		if _, err := s.SendMessageToNode(ctx, msg, source, 1); err != nil {
			return err
		}
	}

	if err := seq.Err(); err != nil {
		return err
	}

	return nil
}
