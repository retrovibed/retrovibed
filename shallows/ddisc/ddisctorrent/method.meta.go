package ddisctorrent

import (
	"context"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/bencode"
	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/dht/krpc"
)

type MetaRequest struct {
	Q string   `bencode:"q,omitempty"` // ddisc_search_query
	A PeerInfo `bencode:"a,omitempty"`
	T string   `bencode:"t"` // required: transaction ID
	Y string   `bencode:"y"` // required: type of the message: q for QUERY, r for RESPONSE, e for ERROR
}

type PeerInfo struct {
	Peer      krpc.ID `bencode:"id"`        // ID of the querying Node
	Version   string  `bencode:"version"`   // version of the ddisc protocol this node runs.
	Partition string  `bencode:"partition"` // required: partition the requesting node tracks.
}

func NewMetaRequest(from int160.T, partition string) (qi dht.QueryInput, err error) {
	req := MetaRequest{
		Y: krpc.YQuery,
		T: krpc.TimestampTransactionID(),
		Q: MethodMeta,
		A: PeerInfo{
			Peer:      from.AsByteArray(),
			Partition: partition,
			Version:   ExtensionName,
		},
	}

	encoded, err := bencode.Marshal(req)
	return dht.NewEncodedRequest(req.Q, req.T, encoded), err
}

type MetaResponse struct {
	Y string   `bencode:"y"` // required: type of the message: q for QUERY, r for RESPONSE, e for ERROR
	T string   `bencode:"t"` // required: transaction ID
	R PeerInfo `bencode:"r,omitempty"`
}

func NewMetaResponse(from int160.T, partition, tid string) MetaResponse {
	req := MetaResponse{
		Y: krpc.YResponse,
		T: tid,
		R: PeerInfo{
			Peer:      from.AsByteArray(),
			Partition: partition,
			Version:   ExtensionName,
		},
	}

	return req
}

func NewMeta(partition uuid.UUID) Meta {
	return Meta{partition: partition}
}

type Meta struct {
	partition uuid.UUID
}

func (t Meta) Handle(ctx context.Context, source dht.Addr, s *dht.Server, raw []byte, _ *krpc.Msg) error {
	var (
		m MetaRequest
	)

	if err := bencode.Unmarshal(raw, &m); err != nil {
		return err
	}

	resp := NewMetaResponse(s.ID(), t.partition.String(), m.T)

	if _, err := s.SendMessageToNode(ctx, resp, source, 3); err != nil {
		return err
	}

	return nil
}
