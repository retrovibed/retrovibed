package timex

import (
	"context"
	"log"
	"math"
	"reflect"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/internal/errorsx"
)

// Inf - positive infinity no time can be larger.
// see https://stackoverflow.com/questions/25065055/what-is-the-maximum-time-time-in-go/32620397
func Inf() time.Time {
	return time.Unix(math.MaxInt64-62135596800, 999999999)
}

func NegInf() time.Time {
	return time.Unix(math.MinInt64, math.MinInt64)
}

// Max select the maximum timestamp from the set.
func Max(ds ...time.Time) (d time.Time) {
	d = NegInf()
	for _, c := range ds {
		if c.After(d) {
			d = c
		}
	}

	return d
}

// Min select the minimum timpstamp from the set.
func Min(ds ...time.Time) (d time.Time) {
	d = Inf()
	for _, c := range ds {
		if c.Before(d) {
			d = c
		}
	}

	return d
}

// Run the provided function after the duration.
func After(d time.Duration, do func()) {
	go func() {
		log.Println("sleepy")
		time.Sleep(d)
		log.Println("awake")
		do()
	}()
}

// Every executes the provided function every duration.
func Every(d time.Duration, do func()) {
	for range time.Tick(d) {
		do()
	}
}

// NowAndEvery executes the provided function immeditately and every duration.
func NowAndEvery(ctx context.Context, d time.Duration, do func(context.Context) error) error {
	if err := do(ctx); err != nil {
		return err
	}

	t := time.NewTicker(d)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			if err := do(ctx); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func NowAndEveryVoid(ctx context.Context, d time.Duration, do func(context.Context)) {
	errorsx.Log(NowAndEvery(ctx, d, func(ctx context.Context) error {
		do(ctx)
		return nil
	}))
}

// DurationOrDefault ...
func DurationOrDefault(a, b time.Duration) time.Duration {
	if a == 0 {
		return b
	}
	return a
}

// DurationMax select the maximum duration from the set.
func DurationMax(ds ...time.Duration) (d time.Duration) {
	for _, c := range ds {
		if c > d {
			d = c
		}
	}

	return d
}

// returns the first non-zero timestamp
func FirstNonZero(s ...time.Time) (_zero time.Time) {
	for _, v := range s {
		if v.IsZero() {
			continue
		}

		return v
	}

	return _zero
}

// DurationMin select the minimum duration from the set.
func DurationMin(ds ...time.Duration) (d time.Duration) {
	d = math.MaxInt64

	for _, c := range ds {
		if c < d {
			d = c
		}
	}

	return d
}

// SafeReset stops and drains the timer (if necessary) and then resets.
func SafeReset(t *time.Timer, d time.Duration) {
	if !t.Stop() {
		select {
		case <-t.C:
		default:
		}
	}
	t.Reset(d)
}

type Clock struct{}

func (t Clock) Now() time.Time {
	return time.Now()
}

// RFC3339NanoMax truncate to the maximum value for RFC3339.
func RFC3339NanoMax(t time.Time) time.Time {
	ts := RFC3339Inf()
	if t.Before(ts) {
		return t
	}

	return ts
}

// RFC3339NanoMin truncate to the minimum value for RFC3339.
func RFC3339NanoMin(t time.Time) time.Time {
	ts := RFC3339NegInf()
	if t.After(ts) && !t.Equal(NegInf()) {
		return t
	}

	return ts
}

// RFC3339NanoEncode truncate time to RFC3339NanoEncode
func RFC3339NanoEncode(t time.Time) time.Time {
	return RFC3339NanoMax(RFC3339NanoMin(t))
}

// RFC3339Nano truncate time to RFC3339Nano
func RFC3339NanoDecode(t time.Time) time.Time {
	return RFC3339NanoMaxDecode(RFC3339NanoMinDecode(t))
}

// RFC3339NanoMinDecode convert minimum value for RFC3339 to time.Time.
func RFC3339NanoMinDecode(t time.Time) time.Time {
	ts := RFC3339NegInf()
	if t.Before(ts) || t.Equal(ts) {
		return NegInf()
	}

	return t
}

// RFC3339NanoMaxDecode  convert maximum value for RFC3339 to time.Time.
func RFC3339NanoMaxDecode(t time.Time) time.Time {
	ts := RFC3339Inf()
	if t.After(ts) || t.Equal(ts) {
		return Inf()
	}

	return t
}

// RFC3339NegInf neg infinity representation
func RFC3339NegInf() time.Time {
	return time.Date(0000, 01, 1, 1, 1, 1, 0, time.UTC)
}

// RFC3339Inf  infinity representation
func RFC3339Inf() time.Time {
	return time.Date(9999, time.December, 31, 23, 59, 59, 999000000, time.UTC)
}

// StartOfDay returns the start of the day (00:00:00) for the given time in UTC.
func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

// EndOfDay returns the end of the day (23:59:59.999999999) for the given time in UTC.
func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, time.UTC)
}

// StartOfWeek returns the start of the week (Monday 00:00:00) for the given time in UTC.
func StartOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return StartOfDay(t.AddDate(0, 0, -weekday+1))
}

// StartOfMonth returns the start of the month for the given time in UTC.
func StartOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
}

func NewRangeEverything() Range {
	return Range{
		Start: NegInf(),
		End:   Inf(),
	}
}

func NewRangeDuration(d time.Duration) Range {
	ts := time.Now()
	return Range{
		Start: ts.Add(-1 * d),
		End:   ts,
	}
}

func NewRangeISO8601(b string, e string) (Range, error) {
	begin, err := time.Parse(time.RFC3339Nano, b)
	if err != nil {
		return Range{}, err
	}

	end, err := time.Parse(time.RFC3339Nano, e)
	if err != nil {
		return Range{}, err
	}

	return Range{Start: RFC3339NanoDecode(begin), End: RFC3339NanoDecode(end)}, nil
}

// Range represents a time range with start and end times.
type Range struct {
	Start time.Time
	End   time.Time
}

func UTCEncodeOption[T any](v *T) {
	_jsonsacodec(reflect.ValueOf(v), func(ts time.Time) time.Time {
		return ts.In(time.UTC)
	})
}

func JSONSafeEncodeOption[T any](v *T) {
	JSONSafeEncode(v)
}

func JSONSafeDecodeOption[T any](v *T) {
	JSONSafeDecode(v)
}

func JSONSafeDecode[T any](v T) T {
	metav := reflect.ValueOf(v)
	_jsonsacodec(metav, RFC3339NanoDecode)
	return v
}

func JSONSafeEncode[T any](v T) T {
	metav := reflect.ValueOf(v)
	_jsonsacodec(metav, RFC3339NanoEncode)
	return v
}

func _jsonsacodec(v reflect.Value, m func(time.Time) time.Time) {
	if v.CanInterface() {
		if ts, ok := v.Interface().(time.Time); ok {
			v.Set(reflect.ValueOf(m(ts)))
			return
		}
	}

	switch v.Kind() {
	case reflect.Struct:
		for _, nv := range reflect.VisibleFields(v.Type()) {
			_jsonsacodec(v.FieldByIndex(nv.Index), m)
		}
	// case reflect.Slice, reflect.Array:
	// case reflect.Interface:
	case reflect.Pointer:
		if v.IsNil() {
			v = reflect.Zero(v.Type().Elem())
		} else {
			v = v.Elem()
		}
		_jsonsacodec(v, m)
	default:
		// do nothing
	}
}
