package asyncx

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/retrovibed/retrovibed/internal/backoffx"
	"github.com/retrovibed/retrovibed/internal/errorsx"
)

const (
	ErrWakeupClosed = errorsx.String("wakeup closed")
)

func Periodic(ctx context.Context, async *Wakeup, b backoffx.Strategy, msg string) {
	for _, delay := range backoffx.Iter(b) {
		async.Broadcast()
		log.Println(msg, delay)
		time.Sleep(delay)
	}
}

type Wakeup struct {
	*sync.Cond
	done    chan struct{}
	C       <-chan error
	cleanup func()
}

func (t *Wakeup) Close() error {
	t.cleanup()
	return nil
}

func NewWakeup(ctx context.Context) *Wakeup {
	ictx, done := context.WithCancel(ctx)
	q := make(chan error, 1)
	wakeup := sync.NewCond(&sync.Mutex{})
	donesig := make(chan struct{})
	a := &Wakeup{
		Cond: wakeup,
		done: donesig,
		C:    q,
		cleanup: sync.OnceFunc(func() {
			log.Println("async wakeup shutdown initiated")
			done()
			select {
			case q <- ErrWakeupClosed:
				log.Println("async wakeup sent closed error")
				wakeup.Broadcast()
			case <-ctx.Done():
				log.Println("async wakeup parent context was cancelled")
			}

			close(donesig)
		}),
	}

	go func() {
		for {
			a.L.Lock()
			a.Wait()
			a.L.Unlock()

			select {
			case <-ictx.Done():
				return
			case <-donesig:
				return
			case q <- nil:
			}
		}
	}()

	return a
}

func Run(ctx context.Context, w *Wakeup, do func(context.Context) error) error {
	for {
		if err := do(ctx); err != nil {
			return err
		}

		select {
		case err := <-w.C:
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return errorsx.Ignore(context.Cause(ctx), context.Canceled)
		}
	}
}

// Same as run but in a goroutine that logs any error that occurs.
func Background(ctx context.Context, w *Wakeup, do func(context.Context) error) {
	go func() {
		errorsx.Log(Run(ctx, w, do))
	}()
}
