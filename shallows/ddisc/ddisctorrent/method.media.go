package ddisctorrent

import (
	"context"

	"github.com/james-lawrence/torrent/bencode"
	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/dht/krpc"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
)

type Media struct {
	Q string `bencode:"q,omitempty"` // ddisc_search_result
	A media  `bencode:"a,omitempty"`
	T string `bencode:"t"` // required: transaction ID
	Y string `bencode:"y"` // required: type of the message: q for QUERY, r for RESPONSE, e for ERROR
}

func NewMediaRecorder(q sqlx.Queryer) SearchRecorder {
	return SearchRecorder{q: q}
}

type SearchRecorder struct {
	q sqlx.Queryer
}

func (t SearchRecorder) Handle(ctx context.Context, source dht.Addr, s *dht.Server, b dht.Binding, raw []byte, _ *krpc.Msg) error {
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
		new(int160.FromBytes(m.A.Infohash)),
		metainfo.Magnet{InfoHash: metainfo.Hash(m.A.Infohash)}.String(),
		ddisc.DiscoveredOptionMimetype(m.A.Mimetype),
		ddisc.DiscoveredOptionKnownMedia(m.A.KnownMediaID),
		mediaToDiscovered(m.A),
	)).Scan(&d)

	if err != nil {
		return err
	}

	return nil
}
