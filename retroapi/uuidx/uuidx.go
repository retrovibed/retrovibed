package uuidx

import (
	"encoding/binary"
	"fmt"
	"math"
	"time"

	"github.com/gofrs/uuid/v5"
)

// WithSuffix makes an almost all balls UUID
func WithSuffix(end int) string {
	return fmt.Sprintf("00000000-0000-0000-0000-%012d", end)
}

// prepend timestamp to a uuid; this is useful for ensuring timebased ordering.
func WithTimestampPrefix(uid uuid.UUID, ts time.Time) uuid.UUID {
	buf := uid.Bytes()
	binary.BigEndian.PutUint64(buf[:8], uint64(ts.UnixMicro()))
	return uuid.FromBytesOrNil(buf)
}

func Advance16(u uuid.UUID, by uint16) uuid.UUID {
	dup := u

	if v := binary.BigEndian.Uint16(dup[:]); math.MaxUint16-v < by {
		binary.BigEndian.PutUint16(dup[:], math.MaxUint16)
	} else {
		binary.BigEndian.PutUint16(dup[:], v+by)
	}

	return dup
}

// returns the first non-zero uuid or zero if all are zero.
func FirstNonNil(uids ...uuid.UUID) uuid.UUID {
	for _, uid := range uids {
		if !uid.IsNil() {
			return uid
		}
	}

	return uuid.Nil
}

func HighN(id uuid.UUID, n int) []byte {
	ubytes := id.Bytes()

	if n <= 0 {
		return []byte(nil)
	}

	if n >= len(ubytes) {
		return ubytes
	}

	return ubytes[:n]
}

func LowN(id uuid.UUID, n int) []byte {
	ubytes := id.Bytes()

	if n <= 0 {
		return []byte(nil)
	}
	if n >= len(ubytes) {
		return append([]byte(nil), ubytes...)
	}

	return ubytes[:n]
}
