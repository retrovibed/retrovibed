package ddiscapi

import (
	"fmt"
	"hash/fnv"

	"github.com/gofrs/uuid/v5"
)

// Well-known known-media source ids, matching what shallows/cmd/cmdmedia's
// importers write into library_known_media.source.
const (
	SourceTMDB        = "tmdb"
	SourceTVDB        = "tvdb"
	SourceMusicbrainz = "musicbrainz"
	SourceDeeppool    = "deeppool"
	SourceUnspecified = "unspecified"
)

// importprefix is a type constraint for import source prefixes.
type importprefix interface {
	~string
}

// ImportedMediaUintID creates a unique import id from a uint sequence.
func ImportedMediaUintID[P importprefix](prefix P, id uint64) string {
	l := id & 0x0000FFFFFFFFFFFF
	h := id & 0xFFFF000000000000 >> 56
	return fmt.Sprintf("%x-0000-0000-%04x-%012x", fnv.New32().Sum([]byte(prefix))[:4], h, l)
}

// ImportedMediaUUID creates a unique import id from a uuid by mutating its first 4 bytes with the prefix checksum.
func ImportedMediaUUID[P importprefix](prefix P, id uuid.UUID) uuid.UUID {
	copy(id[:4], fnv.New32().Sum([]byte(prefix))[:4])
	return id
}
