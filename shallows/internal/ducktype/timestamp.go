package ducktype

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Status byte

const (
	Undefined Status = iota
	Null
	Present
)

type InfinityModifier int8

const (
	Infinity         InfinityModifier = 1
	None             InfinityModifier = 0
	NegativeInfinity InfinityModifier = -Infinity
)

func (im InfinityModifier) String() string {
	switch im {
	case None:
		return "none"
	case Infinity:
		return "infinity"
	case NegativeInfinity:
		return "-infinity"
	default:
		return "invalid"
	}
}

func NewNullTime(ts time.Time) NullTime {
	// these two timestamps are what duckdb returns for pos/neg infinity.
	var (
		inf    = time.UnixMicro(9223372036854775807)
		neginf = time.UnixMicro(-9223372036854775807)
	)

	if ts.Equal(inf) || ts.After(inf) {
		return NullTime{Status: Present, InfinityModifier: Infinity, Time: ts}
	}

	if ts.Equal(neginf) || ts.Before(neginf) {
		return NullTime{Status: Present, InfinityModifier: NegativeInfinity, Time: ts}
	}

	return NullTime{
		Time:   ts,
		Status: Present,
	}
}

type NullTime struct {
	Time             time.Time // Time must always be in UTC.
	Status           Status
	InfinityModifier InfinityModifier
}

func (t *NullTime) Infinity() {
	t.Status = Present
	t.InfinityModifier = Infinity
}

func (t *NullTime) NegativeInfinity() {
	t.Status = Present
	t.InfinityModifier = NegativeInfinity
}

// Scan implements the database/sql Scanner interface.
func (t *NullTime) Scan(src interface{}) error {
	if src == nil {
		*t = NullTime{Status: Null}
		return nil
	}

	switch src := src.(type) {
	case time.Time:
		*t = NewNullTime(src)
		return nil
	}

	return fmt.Errorf("cannot scan %T", src)
}

// Value implements the database/sql/driver Valuer interface.
func (t NullTime) Value() (driver.Value, error) {
	switch t.Status {
	case Present:
		if t.InfinityModifier != None {
			return t.InfinityModifier.String(), nil
		}
		return t.Time, nil
	case Null:
		return nil, nil
	default:
		return nil, fmt.Errorf("undefined timestamp value")
	}
}
