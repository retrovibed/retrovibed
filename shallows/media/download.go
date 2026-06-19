package media

import (
	"time"

	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/shallows/internal/grpcx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"google.golang.org/protobuf/encoding/protojson"
)

func (t *DownloadUpdateRequest) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(t)
}

func (t *DownloadUpdateRequest) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, t)
}

type DownloadOption func(*Download)

func MetadataFromDownload(dl *Download) (md tracking.Metadata, err error) {
	md = tracking.NewMetadata(new(int160.Zero()))
	md.Bytes = dl.Bytes
	md.Downloaded = dl.Downloaded
	md.Peers = uint16(dl.Peers)
	md.Seeding = dl.Distributing

	if md.PausedAt, err = grpcx.DecodeTime(dl.PausedAt); err != nil {
		return md, err
	}
	if md.InitiatedAt, err = grpcx.DecodeTime(dl.InitiatedAt); err != nil {
		return md, err
	}
	if md.VerifyAt, err = grpcx.DecodeTime(dl.VerifyAt); err != nil {
		return md, err
	}
	if md.CompletedAt, err = grpcx.DecodeTime(dl.CompletedAt); err != nil {
		return md, err
	}

	if m := dl.Media; m != nil {
		md.Description = m.Description
		md.KnownMediaID = m.KnownMediaId
		md.Mimetype = m.Mimetype
		md.EncryptionSeed = m.EncryptionSeed
	}

	return md, nil
}

func DownloadOptionFromTorrentMetadata(cc tracking.Metadata) DownloadOption {
	return func(c *Download) {
		c.Media = new(langx.Clone(Media{}, MediaOptionFromTorrentMetadata(cc)))
		c.Bytes = cc.Bytes
		c.Downloaded = cc.Downloaded
		c.Peers = uint32(cc.Peers)
		c.PausedAt = grpcx.EncodeTime(cc.PausedAt)
		c.InitiatedAt = grpcx.EncodeTime(cc.InitiatedAt)
		c.VerifyAt = grpcx.EncodeTime(cc.VerifyAt)
		c.CompletedAt = grpcx.EncodeTime(cc.CompletedAt)
		c.Distributing = cc.Seeding
		c.Path = int160.FromBytes(cc.Infohash).String()
	}
}

func DownloadOptionFromTorrent(t torrent.Torrent) DownloadOption {
	return func(c *Download) {
		cc := t.Stats()

		info := langx.Zero(t.Info())
		c.Bytes = uint64(info.TotalLength())
		c.Downloaded = uint64(t.BytesCompleted())
		c.Peers = uint32(cc.ActivePeers)
		c.PeersHalf = uint32(cc.HalfOpenPeers)
		c.PeersAvailable = uint32(cc.TotalPeers)
		c.PeersSeeders = uint32(cc.Seeders)
		if langx.FirstNonZero(cc.Missing, cc.Unverified, cc.Outstanding) == 0 && cc.Completed > 0 {
			c.CompletedAt = grpcx.EncodeTime(time.Now())
		}
	}
}
