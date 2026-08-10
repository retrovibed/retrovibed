package ddisc

import (
	"encoding/binary"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/uuidx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
)

// KnownMediaFromDiscovered maps a Discovered candidate (already carrying a
// resolved known-media id, e.g. after KnownMediaDetector has run) to a
// library.Known catalog entry keyed by kid - see KnownMediaDynamic, which is
// the only caller. Callers are responsible for only calling this with a
// resolved, non-sentinel kid.
func KnownMediaFromDiscovered(kid uuid.UUID, d Discovered) library.Known {
	// kid is folded into the hash (not just Title/Description) so the md5
	// column's UNIQUE NOT NULL constraint can't collide across two
	// different known-media ids that happen to share a blank/identical
	// title+description - a collision would fail the insert outright,
	// since ON CONFLICT (uid) only covers a uid conflict, not a md5 one.
	contentmd5 := md5x.Digest(kid.String(), d.Title, d.Description)
	uidmd5 := uuid.FromBytesOrNil(contentmd5.Sum(nil))

	return library.Known{
		UID:             kid.String(),
		Md5:             uidmd5.String(),
		Md5Lower:        binary.LittleEndian.Uint64(uuidx.LowN(uidmd5, 64)),
		Title:           d.Title,
		Overview:        d.Description,
		PosterPath:      d.PosterURI,
		Source:          d.Source,
		Mimetype:        Generalize(d.Mimetype),
		Released:        timex.Inf(), // no release date available from the wire format
		AutoDescription: d.Title,
	}
}
