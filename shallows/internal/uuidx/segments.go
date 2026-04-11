package uuidx

import (
	"bytes"

	"github.com/gofrs/uuid/v5"
)

type Range struct {
	Begin uuid.UUID
	End   uuid.UUID
}

// RangeFromStrings panics if the uuids are malformed
func RangeFromStrings(begin, end string) Range {
	return Range{
		Begin: uuid.Must(uuid.FromString(begin)),
		End:   uuid.Must(uuid.FromString(end)),
	}
}

// RangeEverything covers the entire uuid range.
// from uuid.Nil to uuid.Max
func RangeEverything() Range {
	return Range{
		Begin: uuid.Nil,
		End:   uuid.Max,
	}
}

func NewSegmenter() *Segmenter {
	return &Segmenter{
		current: uuid.Nil,
		next:    uuid.Nil,
	}
}

type Segmenter struct {
	current uuid.UUID
	next    uuid.UUID
}

func (t *Segmenter) Next() bool {
	const (
		// we want to generate 1024 segments,
		// calculated increment using ceil(65535 / 1024)
		by = 64
	)

	t.current = t.next
	t.next = Advance16(t.next, by)
	future := Advance16(t.next, by)
	max := uuid.Max
	done := !bytes.Equal(t.next[:], max[:])

	if bytes.Equal(t.next[:], future[:]) {
		t.next = max
	}

	return done
}

func (t *Segmenter) Range() Range {
	return Range{
		Begin: t.current,
		End:   t.next,
	}
}
