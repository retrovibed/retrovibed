package asynccompute

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/retrovibed/retrovibed/shallows/internal/contextx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
)

// pool of workers
type Pool[T any] struct {
	workers  int
	shutdown sync.WaitGroup // track active compute routines
	async    func(ctx context.Context, w T) error
	queued   chan pending
	failed   atomic.Pointer[error]
}

func (t *Pool[T]) Run(ctx context.Context, w T) error {
	select {
	case t.queued <- pending{workload: func() error { return t.async(ctx, w) }}:
		return nil
	case <-ctx.Done():
		return context.Cause(ctx)
	}
}

func (t *Pool[T]) Close() error {
	close(t.queued)
	t.shutdown.Wait()
	return langx.Autoderef(t.failed.Load())
}

func (t *Pool[T]) init() *Pool[T] {
	t.shutdown.Add(int(t.workers))
	for range t.workers {
		go func() {
			defer t.shutdown.Done()
			for pending := range t.queued {
				cause := errorsx.LogErr(pending.workload())
				t.failed.CompareAndSwap(nil, &cause)
			}
		}()
	}

	return t
}

type pending struct {
	workload func() error
}

type Option[T any] func(*Pool[T])

// Backlog sets the number of pending workloads allowed to queue up before
// Pool.Run blocks. n == 0 falls back to runtime.NumCPU(), since an unbuffered
// queue paired with the default worker count can deadlock.
func Backlog[T any](n uint16) Option[T] {
	return func(p *Pool[T]) {
		p.queued = make(chan pending, langx.FirstNonZero(n, uint16(runtime.NumCPU())))
	}
}

// Workers sets the number of worker goroutines draining the queue. n == 0
// falls back to runtime.NumCPU(), since zero workers would leave queued
// workloads unprocessed forever.
func Workers[T any](n uint16) Option[T] {
	return func(p *Pool[T]) {
		p.workers = int(langx.FirstNonZero(n, uint16(runtime.NumCPU())))
	}
}

func Compose[T any](options ...Option[T]) Option[T] {
	return func(p *Pool[T]) {
		for _, opt := range options {
			opt(p)
		}
	}
}

func New[T any](async func(ctx context.Context, w T) error, options ...Option[T]) *Pool[T] {
	return new(langx.Clone(Pool[T]{
		workers: runtime.NumCPU(),
		queued:  make(chan pending, runtime.NumCPU()),
		async:   async,
	}, options...)).init()
}

type closers interface {
	Close() error
}

// gracefully shutdown by invoking close and waiting until all workers
// complete or the context times out. Pools are shut down in the order
// given; every pool is closed even if an earlier one errors, so callers
// passing multiple pools that feed one another (e.g. pool -> insert)
// should list them in drain order.
func Shutdown(ctx context.Context, pools ...closers) error {
	var err error
	for _, p := range pools {
		dctx, cancelled := context.WithCancelCause(ctx)
		go func() {
			cancelled(p.Close())
		}()

		<-dctx.Done()
		err = errorsx.Compact(err, contextx.IgnoreCancelled(context.Cause(dctx)))
	}

	return err
}
