package ddiscapi

import (
	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/grpcx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"google.golang.org/protobuf/encoding/protojson"
)

func (t *Discovery) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(t)
}

func (t *Discovery) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, t)
}

func (t *DiscoveryDownloadRequest) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(t)
}

func (t *DiscoveryDownloadRequest) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, t)
}

func (t *DiscoveryDownloadResponse) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(t)
}

func (t *DiscoveryDownloadResponse) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, t)
}

func NewDiscoveryFromTrackingUnknownHash(mu tracking.UnknownHash) *Discovery {
	mg := metainfo.NewMagnetFromInfohash(mu.Infohash)

	return &Discovery{
		Id:           mu.ID,
		Infohash:     mu.Infohash,
		Attempts:     uint32(mu.Attempts),
		NextCheck:    grpcx.EncodeTime(timex.RFC3339NanoEncode(mu.NextCheck)),
		CreatedAt:    grpcx.EncodeTime(timex.RFC3339NanoEncode(mu.CreatedAt)),
		UpdatedAt:    grpcx.EncodeTime(timex.RFC3339NanoEncode(mu.UpdatedAt)),
		KnownMediaId: uuid.Nil.String(),
		Uri:          mg.String(),
	}
}

// NewDiscoveryFromDiscovered converts a ddisc.Discovered into its wire representation.
//
// this is a manual field mapping rather than the generic grpcx.JSONDecode round-trip
// used for NewDiscoveryFromTrackingUnknownHash above, because ddisc.Discovered.NextCheckAt
// doesn't share a json key with Discovery.NextCheck - same class of problem
// NewMediaFromDiscovered documents for this same source type.
func NewDiscoveryFromDiscovered(d ddisc.Discovered) *Discovery {
	return &Discovery{
		Id:               d.ID,
		Source:           d.Source,
		Infohash:         d.Infohash,
		Attempts:         d.Attempts,
		NextCheck:        grpcx.EncodeTime(timex.RFC3339NanoEncode(d.NextCheckAt)),
		CreatedAt:        grpcx.EncodeTime(timex.RFC3339NanoEncode(d.CreatedAt)),
		UpdatedAt:        grpcx.EncodeTime(timex.RFC3339NanoEncode(d.UpdatedAt)),
		Title:            d.Title,
		Description:      d.Description,
		Health:           d.Health,
		Bytes:            d.Bytes,
		PolicyRank:       uint32(d.PolicyRank),
		KnownMediaId:     uuid.FromStringOrNil(d.KnownMediaID).String(),
		Uri:              d.URI,
		AcquisitionState: AcquisitionState(d.AcquisitionState),
	}
}

// NewDiscoveredFromDiscovery converts a client-submitted Discovery back into
// a ddisc.Discovered - the inverse of NewDiscoveryFromDiscovered. Used only
// for candidates the caller streamed from an ephemeral (never persisted)
// strategy (Known/Plugin/PeerTube) and is now asking to download; d.Id is
// preserved as-is so the corrected-infohash insert in
// ddisc.DownloadDiscovered lands on the identity the client already has.
func NewDiscoveredFromDiscovery(d *Discovery) ddisc.Discovered {
	ih := int160.FromBytesOrZero(d.Infohash)
	m := ddisc.NewDiscovered(
		&ih,
		ddisc.DiscoveredOptionTitle(d.Title),
		ddisc.DiscoveredOptionDescription(d.Description),
		ddisc.DiscoveredOptionHealth(d.Health),
		ddisc.DiscoveredOptionURI(d.Uri),
		ddisc.DiscoveredOptionKnownMedia(uuid.FromStringOrNil(d.KnownMediaId).String()),
	)
	m.ID = d.Id
	m.Source = d.Source
	m.Bytes = d.Bytes
	m.Attempts = d.Attempts
	return m
}
