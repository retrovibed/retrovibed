package ddisc

import (
	"encoding/binary"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/ddiscapi"
	"github.com/retrovibed/retrovibed/retroapi/uuidx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
)

// KnownMediaFromImport maps a search plugin's Import result to a
// library.Known catalog entry keyed by kid. Callers are responsible for
// only calling this with a resolved, non-sentinel kid (see pluginSeq.Each).
func KnownMediaFromImport(kid uuid.UUID, mimetype string, imp *ddiscapi.Import) library.Known {
	// kid is folded into the hash (not just Title/Overview) so the md5
	// column's UNIQUE NOT NULL constraint can't collide across two
	// different known-media ids that happen to share a blank/identical
	// title+overview (title is optional here) - a collision would fail
	// the insert outright, since ON CONFLICT (uid) only covers a uid
	// conflict, not a md5 one.
	contentmd5 := md5x.Digest(kid.String(), imp.Title, imp.Overview)
	uidmd5 := uuid.FromBytesOrNil(contentmd5.Sum(nil))

	return library.Known{
		UID:             kid.String(),
		Md5:             uidmd5.String(),
		Md5Lower:        binary.LittleEndian.Uint64(uuidx.LowN(uidmd5, 64)),
		Title:           imp.Title,
		Overview:        imp.Overview,
		Popularity:      imp.Popularity,
		PosterPath:      imp.PosterPath,
		Source:          imp.Source,
		Mimetype:        Generalize(mimetype),
		Released:        timex.Inf(), // no release date available from the wire format
		AutoDescription: imp.Title,
	}
}
