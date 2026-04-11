package ddiscapi

import (
	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/internal/grpcx"
	"github.com/retrovibed/retrovibed/internal/langx"
	"github.com/retrovibed/retrovibed/internal/timex"
	"github.com/retrovibed/retrovibed/tracking"
)

func NewPeerFromTrackingPeer(mp tracking.Peer) (_ *Peer, err error) {
	var p Peer
	mp = langx.Clone(mp, timex.JSONSafeEncodeOption, timex.UTCEncodeOption)

	if err = grpcx.JSONDecode(mp, &p); err != nil {
		return nil, err
	}

	p.Infohash = mp.Peer
	p.Partition = mp.DdiscPartition
	p.Syncoffset = mp.DdiscSyncoffset

	return &p, nil
}

func NewTrackingPeerFromPeer(p *Peer, options ...tracking.PeerOption) (mp tracking.Peer, err error) {
	if err = grpcx.JSONEncode(p, &mp); err != nil {
		return mp, err
	}

	mp = tracking.NewPeer(
		int160.FromBytesOrZero(p.Infohash),
		tracking.PeerOptionDdisc(uuid.FromStringOrNil(p.Partition), uuid.FromStringOrNil(p.Syncoffset)),
		tracking.PeerOptionTimestampClone(mp),
		tracking.PeerOptionDescription(mp.Description),
		timex.JSONSafeDecodeOption,
		timex.UTCEncodeOption,
		langx.Compose(options...),
	)

	return mp, nil
}
