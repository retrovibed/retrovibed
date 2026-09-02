package pqueuex

import (
	"context"
	"log"
	"time"

	"github.com/linxGnu/pqueue"
	"github.com/linxGnu/pqueue/entry"
	"github.com/retrovibed/retrovibed/retroapi/backoffx"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
)

func New(dir string) (q pqueue.Queue, err error) {
	if err := fsx.MkDirs(0700, dir); err != nil {
		return nil, errorsx.Wrapf(err, "unable to create queues dir: %s", dir)
	}

	if q, err = pqueue.New(dir, 32); err != nil {
		return nil, errorsx.Wrapf(err, "unable to create queue: %s", dir)
	}

	return q, nil
}

func Enqueue(ctx context.Context, q pqueue.Queue, v any) error {
	encoded, err := jsonx.Marshal(v)
	if err != nil {
		return errorsx.Wrap(err, "failed to encode")
	}

	return errorsx.Wrap(q.Enqueue(entry.Entry(encoded)), "failed to enqueue")
}

type Handler interface {
	Message(ctx context.Context, m []byte) error
}

func NewWorker(wq pqueue.Queue, fn Handler, options ...func(*worker)) worker {
	s := backoffx.New(
		backoffx.Exponential(200*time.Millisecond),
		backoffx.Maximum(30*time.Second),
		backoffx.JitterRandom(50*time.Millisecond),
	)

	return langx.Clone(worker{wq: wq, s: s, fn: fn}, options...)
}

type worker struct {
	wq pqueue.Queue
	s  backoffx.Strategy
	fn Handler
}

// consume the queue until context is cancelled.
func (t worker) Consume(ctx context.Context) {
	for ctx.Err() == nil {
		var (
			m        entry.Entry
			attempts = backoffx.Iter(t.s)
		)

		for _, delay := range attempts {
			if t.wq.Dequeue(&m) {
				break
			}

			time.Sleep(delay)
		}

		if err := t.fn.Message(ctx, m); err != nil {
			log.Println("failed to process message", err)
			continue
		}
	}
}
