package timex_test

import (
	"log"
	"testing"
	"time"

	"github.com/retrovibed/retrovibed/internal/timex"
	"github.com/stretchr/testify/require"
)

func TestJSONSafeDecodeNowShouldRemainUnchanged(t *testing.T) {
	type foo struct {
		Timestamp time.Time
		Bar       struct {
			Timestamp time.Time
		}
	}

	ts := time.Now()
	tmp := timex.JSONSafeDecode(&foo{Timestamp: ts, Bar: struct{ Timestamp time.Time }{Timestamp: ts}})
	require.Equal(t, tmp.Timestamp, ts)
	require.Equal(t, tmp.Bar.Timestamp, ts)
}

func TestJSONSafeDecodeInfShouldBeAdjusted(t *testing.T) {
	type foo struct {
		Timestamp time.Time
		Bar       struct {
			Timestamp time.Time
		}
	}

	tmp := timex.JSONSafeDecode(&foo{Timestamp: timex.RFC3339Inf(), Bar: struct{ Timestamp time.Time }{Timestamp: timex.RFC3339Inf()}})
	log.Println(tmp.Timestamp, timex.Inf())
	require.Equal(t, tmp.Timestamp, timex.RFC3339NanoDecode(timex.Inf()))
	require.NotEqual(t, tmp.Timestamp, timex.RFC3339Inf())

	require.Equal(t, tmp.Bar.Timestamp, timex.RFC3339NanoDecode(timex.Inf()))
	require.NotEqual(t, tmp.Bar.Timestamp, timex.RFC3339Inf())
}

func TestJSONSafeEncodeNowShouldRemainUnchanged(t *testing.T) {
	type foo struct {
		Timestamp time.Time
		Bar       struct {
			Timestamp time.Time
		}
	}

	ts := time.Now()
	tmp := timex.JSONSafeEncode(&foo{Timestamp: ts, Bar: struct{ Timestamp time.Time }{Timestamp: ts}})
	require.Equal(t, tmp.Timestamp, ts)
	require.Equal(t, tmp.Bar.Timestamp, ts)
}

func TestJSONSafeEncodeInfShouldBeAdjusted(t *testing.T) {
	type foo struct {
		Timestamp time.Time
		Bar       struct {
			Timestamp time.Time
		}
	}

	tmp := timex.JSONSafeEncode(&foo{Timestamp: timex.Inf(), Bar: struct{ Timestamp time.Time }{Timestamp: timex.Inf()}})
	log.Println(tmp.Timestamp, timex.Inf())
	require.Equal(t, tmp.Timestamp, timex.RFC3339Inf())
	require.NotEqual(t, tmp.Timestamp, timex.Inf())

	require.Equal(t, tmp.Bar.Timestamp, timex.RFC3339Inf())
	require.NotEqual(t, tmp.Bar.Timestamp, timex.Inf())
}
