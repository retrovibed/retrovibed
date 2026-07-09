package ddisc

import (
	"context"
	"errors"
	"iter"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"golang.org/x/time/rate"
)

// UnknownEnqueueOption configures EnqueueUnknown.
type UnknownEnqueueOption func(*UnknownEnqueue)

// UnknownEnqueueOptionLimiter rate limits how often a miss is allowed to
// enqueue a known_media_id. Share one limiter across many Wrap calls (e.g.
// one held by a DHT search handler across all its requests) to actually
// bound enqueue frequency — a limiter constructed fresh per call has no
// effect. Without this option, every miss enqueues (or idempotently touches
// an already-queued row).
func UnknownEnqueueOptionLimiter(l *rate.Limiter) UnknownEnqueueOption {
	return func(t *UnknownEnqueue) {
		t.limiter = l
	}
}

// EnqueueUnknown builds a reusable policy for enqueuing known-media-ids into
// ddisc_search_queue for background discovery via external search plugins
// whenever a Wrapped sequence turns up nothing — this is meant to apply
// even to lookups triggered by a remote peer's DHT query, so misses become a
// signal worth investigating regardless of who's asking.
func EnqueueUnknown(options ...UnknownEnqueueOption) UnknownEnqueue {
	t := UnknownEnqueue{
		limiter: rate.NewLimiter(rate.Every(time.Minute), 1),
	}
	for _, opt := range options {
		opt(&t)
	}
	return t
}

type UnknownEnqueue struct {
	limiter *rate.Limiter
}

// Wrap s so that, if it yields nothing, kid is enqueued for search plugin
// discovery via q. Callers see the same empty-result, nil-error contract as
// s itself; only a genuine underlying error propagates.
func (t UnknownEnqueue) Wrap(q sqlx.Queryer, kid string, s iterx.Seq[Discovered]) iterx.Seq[Discovered] {
	return &unknownEnqueueSeq{
		inner: iterx.NotFound(s),
		q:     q,
		kid:   kid,
		cfg:   t,
	}
}

type unknownEnqueueSeq struct {
	inner iterx.Seq[Discovered]
	q     sqlx.Queryer
	kid   string
	cfg   UnknownEnqueue
	err   error
}

func (t *unknownEnqueueSeq) Each(ctx context.Context) iter.Seq[Discovered] {
	return func(yield func(Discovered) bool) {
		for v := range t.inner.Each(ctx) {
			if !yield(v) {
				return
			}
		}

		if !errors.Is(t.inner.Err(), iterx.ErrNotFound) {
			t.err = t.inner.Err()
			return
		}

		if !t.cfg.limiter.Allow() {
			return
		}

		var sq SearchQueue
		errorsx.Log(errorsx.Wrap(SearchQueueEnqueue(ctx, t.q, SearchQueue{KnownMediaID: t.kid}).Scan(&sq), "unable to enqueue known media for search"))
	}
}

func (t *unknownEnqueueSeq) Err() error {
	return t.err
}
