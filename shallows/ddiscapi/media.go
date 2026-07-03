package ddiscapi

import (
	"time"

	"github.com/retrovibed/retrovibed/shallows/ddisc"
)

// NewMediaFromDiscovered converts a ddisc.Discovered into its wire representation.
//
// this is a manual field mapping rather than the generic grpcx.JSONDecode round-trip
// used elsewhere (e.g. NewPeerFromTrackingPeer) because ddisc.Discovered.Bytes is
// tagged `json:"bytes,string"` while the generated Media.Bytes has no `,string`
// modifier, which breaks a generic JSON round trip.
func NewMediaFromDiscovered(d ddisc.Discovered) *Media {
	return &Media{
		Id:           d.ID,
		Infohash:     d.Infohash,
		Title:        d.Title,
		Description:  d.Description,
		Mimetype:     d.Mimetype,
		KnownMediaId: d.KnownMediaID,
		Partition:    d.Partition,
		Bytes:        d.Bytes,
		Attempts:     d.Attempts,
		NextCheckAt:  d.NextCheckAt.Format(time.RFC3339),
		CreatedAt:    d.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    d.UpdatedAt.Format(time.RFC3339),
	}
}
