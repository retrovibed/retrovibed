package ddisctorrent

import (
	"context"

	"github.com/james-lawrence/torrent/bencode"
	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/dht/krpc"
	"github.com/retrovibed/retrovibed/ddisc"
	"github.com/retrovibed/retrovibed/internal/sqlx"
)

type SearchRequest struct {
	Q string       `bencode:"q,omitempty"` // ddisc_search_query
	A SearchParams `bencode:"a,omitempty"`
	T string       `bencode:"t"` // required: transaction ID
	Y string       `bencode:"y"` // required: type of the message: q for QUERY, r for RESPONSE, e for ERROR
}

type SearchParams struct {
	Peer         krpc.ID `bencode:"pid"`   // ID of the querying Node
	KnownMediaID string  `bencode:"kid"`   // required: known media id
	Token        string  `bencode:"token"` // required: authorization token
}

func NewSearchRequest(from int160.T, id string) (qi dht.QueryInput, err error) {
	req := SearchRequest{
		Y: "q",
		T: krpc.TimestampTransactionID(),
		Q: MethodSearch,
		A: SearchParams{
			Peer:         from.AsByteArray(),
			KnownMediaID: id,
		},
	}

	encoded, err := bencode.Marshal(req)
	return dht.NewEncodedRequest(req.Q, req.T, encoded), err
}

func NewSearch(q sqlx.Queryer) Search {
	return Search{q: q}
}

type Search struct {
	q sqlx.Queryer
}

func (t Search) Handle(ctx context.Context, source dht.Addr, s *dht.Server, raw []byte, _ *krpc.Msg) error {
	var (
		m SearchRequest
	)

	if err := bencode.Unmarshal(raw, &m); err != nil {
		return err
	}

	if _, err := s.SendMessageToNode(ctx, krpc.NewEmptyResponse(m.T), source, 1); err != nil {
		return err
	}

	seq := ddisc.FindMedia(t.q, m.A.KnownMediaID)
	for v := range seq.Each(ctx) {
		msg := Media{
			Q: MethodMedia,
			Y: "q",
			T: krpc.TimestampTransactionID(),
			A: mediaFromDiscovered(m.A.Token, s.ID(), &v),
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
