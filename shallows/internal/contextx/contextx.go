package contextx

import (
	"context"
	"errors"
	"sync"
)

type keys int

const (
	contextKeyWaitgroup keys = iota
)

func WithWaitGroup(ctx context.Context, wg *sync.WaitGroup) context.Context {
	return context.WithValue(ctx, contextKeyWaitgroup, wg)
}

// NewWaitGroup - adds a waitgroup to the context.
func NewWaitGroup(ctx context.Context) context.Context {
	return WithWaitGroup(ctx, &sync.WaitGroup{})
}

// WaitGroup - retrieve the waitgroup from the context.
func WaitGroup(ctx context.Context) (*sync.WaitGroup, bool) {
	wg, ok := ctx.Value(contextKeyWaitgroup).(*sync.WaitGroup)
	return wg, ok
}

// WaitGroupAdd - increment the waitgroup by delta.
func WaitGroupAdd(ctx context.Context, delta int) {
	if wg, ok := WaitGroup(ctx); ok {
		wg.Add(delta)
	}
}

// WaitGroupDone - decrement the waitgroup
func WaitGroupDone(ctx context.Context) {
	if wg, ok := WaitGroup(ctx); ok {
		wg.Done()
	}
}

func IgnoreDeadlineExceeded(err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		return nil
	}
	return err
}

func IsDeadlineExceeded(err error) bool {
	return errors.Is(err, context.DeadlineExceeded)
}

func IgnoreCancelled(err error) error {
	if errors.Is(err, context.Canceled) {
		return nil
	}
	return err
}

func IsCancelled(err error) bool {
	return errors.Is(err, context.Canceled)
}
