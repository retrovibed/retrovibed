package pqueuex

import (
	"context"
	"log"
	"time"

	"github.com/linxGnu/pqueue"
	"github.com/linxGnu/pqueue/entry"
	"github.com/retrovibed/retrovibed/retroapi/backoffx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
)

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
