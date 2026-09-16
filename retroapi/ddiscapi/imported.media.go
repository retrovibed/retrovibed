package ddiscapi

import (
	"fmt"
	"hash/fnv"
	"strconv"

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

// Well-known media-kind namespaces, passed as ImportedMediaUintID's
// namespaces argument to keep otherwise-colliding external ids from the
// same source apart (e.g. a TMDB movie and a TMDB tv show, or a TMDB
// movie and a TMDB tv episode, can have the same numeric id).
const (
	KindMovie   = "movie"
	KindSeries  = "series"
	KindEpisode = "episode"
	KindRelease = "release"
)

// importprefix is a type constraint for import source prefixes.
type importprefix interface {
	~string
}

func _uint[P importprefix, X ~uint64 | ~uint32](id X) P {
	return P(strconv.FormatUint(uint64(id), 10))
}

func _prefix[P importprefix](prefix P) []byte {
	return fnv.New32().Sum([]byte(prefix))[:4]
}

func _namespace[P importprefix, X ~uint64 | ~uint32](prefix P, id X, namespaces ...P) uint32 {
	d := fnv.New32()

	_, _ = d.Write([]byte(prefix))
	_, _ = d.Write([]byte(_uint[P](id)))

	for _, n := range namespaces {
		_, _ = d.Write([]byte(n))
	}

	return d.Sum32()
}

// ImportedMediaUintID creates a unique import id from a uint sequence.
func ImportedMediaUintID[P importprefix](prefix P, id uint64, namespaces ...P) string {
	l := id & 0x0000FFFFFFFFFFFF
	h := id & 0xFFFF000000000000 >> 56
	ns := _namespace(prefix, id, namespaces...)
	return fmt.Sprintf(
		"%x-%04x-%04x-%04x-%012x",
		_prefix(prefix),
		uint16(ns>>16), uint16(ns),
		h, l,
	)
}

// ImportedMediaUUID creates a unique import id from a uuid by mutating its first 4 bytes with the prefix checksum.
func ImportedMediaUUID[P importprefix](prefix P, id uuid.UUID) uuid.UUID {
	copy(id[:4], _prefix(prefix))
	return id
}
