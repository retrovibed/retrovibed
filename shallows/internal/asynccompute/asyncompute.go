package asynccompute

import (
	"context"
	"runtime"
	"sync"

	"github.com/retrovibed/retrovibed/internal/x/contextx"
	"github.com/retrovibed/retrovibed/internal/x/langx"
)

// pool of workers
type Pool[T any] struct {
	workers  int
	shutdown sync.WaitGroup // track active compute routines
	async    func(ctx context.Context, w T) error
	queued   chan pending[T]
}

func (t *Pool[T]) Run(ctx context.Context, w T) (context.Context, error) {
	ctx, cancelled := context.WithCancelCause(ctx)
	select {
	case t.queued <- pending[T]{ctx: ctx, workload: w, completed: cancelled}:
		return ctx, nil
	case <-ctx.Done():
		return ctx, context.Cause(ctx)
	}
}

func (t *Pool[T]) Close() {
	close(t.queued)
	t.shutdown.Wait()
}

func (t *Pool[T]) init() *Pool[T] {
	t.shutdown.Add(int(t.workers))
	for i := 0; i < t.workers; i++ {
		go func() {
			defer t.shutdown.Done()
			for pending := range t.queued {
				pending.completed(t.async(pending.ctx, pending.workload))
			}
		}()
	}

	return t
}

type pending[T any] struct {
	ctx       context.Context
	completed context.CancelCauseFunc
	workload  T
}

type option[T any] func(*Pool[T])

func Backlog[T any](n uint16) option[T] {
	return func(p *Pool[T]) {
		p.queued = make(chan pending[T], n)
	}
}

func Workers[T any](n uint16) option[T] {
	return func(p *Pool[T]) {
		p.workers = int(n)
	}
}

func New[T any](async func(ctx context.Context, w T) error, options ...option[T]) *Pool[T] {
	return langx.Autoptr(langx.Clone(Pool[T]{
		workers: runtime.NumCPU(),
		queued:  make(chan pending[T], runtime.NumCPU()),
		async:   async,
	}, options...)).init()
}

// gracefully shutdown by invoking close and waiting until all workers
// complete or the context times out.
func Shutdown[T any](ctx context.Context, p *Pool[T]) error {
	dctx, cancelled := context.WithCancelCause(ctx)
	go func() {
		p.Close()
		cancelled(nil)
	}()

	<-dctx.Done()
	return contextx.IgnoreCancelled(context.Cause(dctx))
}
