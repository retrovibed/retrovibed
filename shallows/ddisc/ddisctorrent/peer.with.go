package ddisctorrent

import (
	"context"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/bencode"
	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/dht/krpc"
	"github.com/retrovibed/retrovibed/internal/errorsx"
)

func PeerWith(ctx context.Context, s *dht.Server, partition uuid.UUID, n krpc.NodeInfo) (_zero PeerInfo, err error) {
	var (
		resp MetaResponse
	)

	qi, err := NewMetaRequest(s.ID(), partition.String())
	if err != nil {
		return _zero, errorsx.Wrapf(err, "unable to generate search request: %s", n.ID)
	}

	ret := s.Query(ctx, dht.NewAddr(n.Addr.AddrPort), qi)
	if err := ret.Err; err != nil {
		return _zero, errorsx.Wrap(err, "meta query failed")
	}

	if err := bencode.Unmarshal(ret.Raw, &resp); err != nil {
		if _, ok := err.(bencode.ErrUnusedTrailingBytes); !ok {
			return _zero, errorsx.Wrapf(err, "unable to deserialize sample response: %T %s", err, n.ID)
		}
	}

	return resp.R, nil
}
