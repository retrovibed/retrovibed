package ddiscapi

import (
	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/grpcx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

func NewDiscoveryFromTrackingUnknownHash(mu tracking.UnknownHash) *Discovery {
	return &Discovery{
		Id:           mu.ID,
		Infohash:     mu.Infohash,
		Attempts:     uint32(mu.Attempts),
		NextCheck:    grpcx.EncodeTime(timex.RFC3339NanoEncode(mu.NextCheck)),
		CreatedAt:    grpcx.EncodeTime(timex.RFC3339NanoEncode(mu.CreatedAt)),
		UpdatedAt:    grpcx.EncodeTime(timex.RFC3339NanoEncode(mu.UpdatedAt)),
		KnownMediaId: uuid.Nil.String(),
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
		Id:           d.ID,
		Source:       d.Source,
		Infohash:     d.Infohash,
		Attempts:     d.Attempts,
		NextCheck:    grpcx.EncodeTime(timex.RFC3339NanoEncode(d.NextCheckAt)),
		CreatedAt:    grpcx.EncodeTime(timex.RFC3339NanoEncode(d.CreatedAt)),
		UpdatedAt:    grpcx.EncodeTime(timex.RFC3339NanoEncode(d.UpdatedAt)),
		Title:        d.Title,
		Description:  d.Description,
		Health:       d.Health,
		Bytes:        d.Bytes,
		PolicyRank:   uint32(d.PolicyRank),
		KnownMediaId: uuid.FromStringOrNil(d.KnownMediaID).String(),
	}
}
