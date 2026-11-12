package asyncx

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/backoffx"
	"github.com/retrovibed/retrovibed/retroapi/errorsx"
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
	C       <-chan error
	cleanup func()
}

func (t *Wakeup) Close() error {
	t.cleanup()
	return nil
}

func NewWakeup(ctx context.Context) *Wakeup {
	q := make(chan error, 1)
	closing := make(chan error, 1)
	wakeup := sync.NewCond(&sync.Mutex{})
	closed := func() error {
		select {
		case err := <-closing:
			return err
		default:
			return nil
		}
	}
	a := &Wakeup{
		Cond: wakeup,
		C:    q,
		cleanup: sync.OnceFunc(func() {
			log.Println("async wakeup shutdown initiated")
			select {
			case closing <- ErrWakeupClosed:
				log.Println("async wakeup sent closed error")
				wakeup.L.Lock()
				wakeup.Broadcast()
				wakeup.L.Unlock()
			case <-ctx.Done():
				log.Println("async wakeup parent context was cancelled")
			}

			close(closing)
		}),
	}

	go func() {
		for {
			a.L.Lock()
			// closing is filled before the broadcast, so checking it under the lock
			// means the broadcast cannot be lost.
			if len(closing) == 0 {
				a.Wait()
			}
			a.L.Unlock()

			err := closed()
			select {
			case <-ctx.Done():
				return
			case q <- err:
				if errors.Is(err, ErrWakeupClosed) {
					return
				}
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
