package ddisctorrent

import (
	"context"
	"log"

	"github.com/james-lawrence/torrent/bencode"
	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/dht/krpc"
	"github.com/retrovibed/retrovibed/ddisc"
	"github.com/retrovibed/retrovibed/internal/langx"
	"github.com/retrovibed/retrovibed/internal/sqlx"
)

type SyncRequest struct {
	Q string     `bencode:"q,omitempty"` // ddisc_sync
	A SyncParams `bencode:"a,omitempty"`
	T string     `bencode:"t"` // required: transaction ID
	Y string     `bencode:"y"` // required: type of the message: q for QUERY, r for RESPONSE, e for ERROR
}

type SyncParams struct {
	Peer      krpc.ID `bencode:"pid"`       // ID of the querying Node
	Token     string  `bencode:"token"`     // required: authorization token for responses.
	Offset    string  `bencode:"offset"`    // required: sync offset.
	Partition string  `bencode:"partition"` // partition we want.
}

func NewSyncRequest(from int160.T, partition string, offset string) (qi dht.QueryInput, err error) {
	req := SyncRequest{
		Y: krpc.YQuery,
		T: krpc.TimestampTransactionID(),
		Q: MethodSync,
		A: SyncParams{
			Peer:      from.AsByteArray(),
			Partition: partition,
			Offset:    offset,
		},
	}

	encoded, err := bencode.Marshal(req)
	return dht.NewEncodedRequest(req.Q, req.T, encoded), err
}

func NewSync(q sqlx.Queryer) Sync {
	return Sync{
		q: q,
	}
}

type Sync struct {
	q sqlx.Queryer
}

func (t Sync) Handle(ctx context.Context, source dht.Addr, s *dht.Server, raw []byte, _ *krpc.Msg) error {
	var (
		m SyncRequest
	)

	log.Println("----------------------------------- sync request initiated -----------------------------------")
	defer log.Println("----------------------------------- sync request completed -----------------------------------")

	if err := bencode.Unmarshal(raw, &m); err != nil {
		return err
	}

	if _, err := s.SendMessageToNode(ctx, krpc.NewEmptyResponse(m.T), source, 1); err != nil {
		return err
	}
	seq := ddisc.SyncPartition(t.q, m.A.Partition, m.A.Offset)
	for v := range seq.Each(ctx) {
		msg := Media{
			Q: MethodMedia,
			Y: krpc.YQuery,
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

type SyncRecorder struct {
	q sqlx.Queryer
}

func (t SyncRecorder) Handle(ctx context.Context, source dht.Addr, s *dht.Server, raw []byte, _ *krpc.Msg) error {
	var (
		m Media
		d ddisc.Discovered
	)

	if err := bencode.Unmarshal(raw, &m); err != nil {
		return err
	}

	// TODO: validate token.

	// ack the message.
	if _, err := s.SendMessageToNode(ctx, krpc.NewEmptyResponse(m.T), source, 1); err != nil {
		return err
	}

	err := ddisc.DiscoveredInsertWithDefaults(ctx, t.q, ddisc.NewDiscovered(
		langx.Autoptr(int160.FromBytes(m.A.Infohash)),
		ddisc.DiscoveredOptionMimetype(m.A.Mimetype),
		ddisc.DiscoveredOptionKnownMedia(m.A.KnownMediaID),
		mediaToDiscovered(m.A),
	)).Scan(&d)

	if err != nil {
		return err
	}

	return nil
}
