package pqueuex

import (
	"context"
	"errors"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/linxGnu/pqueue/entry"
	"github.com/retrovibed/retrovibed/retroapi/backoffx"
	"github.com/retrovibed/retrovibed/shallows/internal/pqueuetestx"
	"github.com/stretchr/testify/require"
)

func TestConsumeDequeuesAndDispatches(t *testing.T) {
	q := pqueuetestx.NewQueue()
	require.NoError(t, q.Enqueue(entry.Entry("hello")))
	h := pqueuetestx.NewHandler()
	w := worker{wq: q, s: backoffx.Constant(time.Millisecond), fn: h}

	ctx, cancel := context.WithCancel(t.Context())
	done := make(chan struct{})
	go func() {
		w.Consume(ctx)
		close(done)
	}()

	require.Eventually(t, func() bool {
		return h.Count() == 1
	}, time.Second, time.Millisecond)
	require.Equal(t, [][]byte{[]byte("hello")}, h.Snapshot())

	cancel()
	require.NoError(t, q.Enqueue(entry.Entry("unblock")))
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Consume did not return after context cancellation")
	}
}

func TestConsumeProcessesMultipleMessagesInOrder(t *testing.T) {
	q := pqueuetestx.NewQueue()
	require.NoError(t, q.Enqueue(entry.Entry("first")))
	require.NoError(t, q.Enqueue(entry.Entry("second")))
	require.NoError(t, q.Enqueue(entry.Entry("third")))
	h := pqueuetestx.NewHandler()
	w := worker{wq: q, s: backoffx.Constant(time.Millisecond), fn: h}

	ctx, cancel := context.WithCancel(t.Context())
	done := make(chan struct{})
	go func() {
		w.Consume(ctx)
		close(done)
	}()

	require.Eventually(t, func() bool {
		return h.Count() == 3
	}, time.Second, time.Millisecond)
	require.Equal(t, [][]byte{[]byte("first"), []byte("second"), []byte("third")}, h.Snapshot())

	cancel()
	require.NoError(t, q.Enqueue(entry.Entry("unblock")))
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Consume did not return after context cancellation")
	}
}

func TestConsumeContinuesAfterHandlerError(t *testing.T) {
	q := pqueuetestx.NewQueue()
	require.NoError(t, q.Enqueue(entry.Entry("bad")))
	require.NoError(t, q.Enqueue(entry.Entry("good")))
	h := pqueuetestx.NewHandler(pqueuetestx.HandlerOptionErrorAt(0, errors.New("boom")))
	w := worker{wq: q, s: backoffx.Constant(time.Millisecond), fn: h}

	ctx, cancel := context.WithCancel(t.Context())
	done := make(chan struct{})
	go func() {
		w.Consume(ctx)
		close(done)
	}()

	// Both messages are dispatched exactly once each even though the first
	// failed: a handler error is logged and the loop moves on, it does not
	// retry or re-enqueue the failed message.
	require.Eventually(t, func() bool {
		return h.Count() == 2
	}, time.Second, time.Millisecond)
	require.Equal(t, [][]byte{[]byte("bad"), []byte("good")}, h.Snapshot())

	cancel()
	require.NoError(t, q.Enqueue(entry.Entry("unblock")))
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Consume did not return after context cancellation")
	}
}

func TestConsumeRetriesUntilMessageAvailable(t *testing.T) {
	q := pqueuetestx.NewQueue()
	h := pqueuetestx.NewHandler()

	var (
		mu       sync.Mutex
		attempts []int
	)
	strategy := backoffx.StrategyFunc(func(attempt int) time.Duration {
		mu.Lock()
		attempts = append(attempts, attempt)
		mu.Unlock()
		return time.Millisecond
	})
	w := worker{wq: q, s: strategy, fn: h}

	ctx, cancel := context.WithCancel(t.Context())
	done := make(chan struct{})
	go func() {
		w.Consume(ctx)
		close(done)
	}()

	require.Eventually(t, func() bool {
		mu.Lock()
		defer mu.Unlock()
		return len(attempts) >= 3
	}, time.Second, time.Millisecond, "expected the backoff strategy to be consulted while the queue is empty")

	// Snapshot before enqueueing: once the message is available and processed,
	// the outer loop starts a fresh attempts iterator for its next (empty-queue)
	// wait, which would otherwise contaminate this assertion with a reset count.
	mu.Lock()
	seen := append([]int(nil), attempts...)
	mu.Unlock()
	require.True(t, sort.IntsAreSorted(seen), "expected non-decreasing attempt numbers: %v", seen)
	require.Zero(t, seen[0])

	require.NoError(t, q.Enqueue(entry.Entry("delayed")))
	require.Eventually(t, func() bool {
		return h.Count() == 1
	}, time.Second, time.Millisecond)
	require.Equal(t, [][]byte{[]byte("delayed")}, h.Snapshot())

	cancel()
	require.NoError(t, q.Enqueue(entry.Entry("unblock")))
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Consume did not return after context cancellation")
	}
}

// TestConsumeIgnoresContextCancellationWhileQueueEmpty characterizes a known
// limitation of the current implementation: worker.Consume's inner retry loop
// ranges over backoffx.Iter, which never checks ctx.Err() itself, so
// cancellation is only observed between successfully dequeued messages. This
// is not desired long-term behavior, just the current one.
func TestConsumeIgnoresContextCancellationWhileQueueEmpty(t *testing.T) {
	q := pqueuetestx.NewQueue()
	h := pqueuetestx.NewHandler()

	// Wait for the first Backoff call before cancelling, so cancellation
	// happens once Consume is provably inside the inner wait loop rather than
	// racing its very first ctx.Err() check at the top of the outer loop.
	started := make(chan struct{})
	var once sync.Once
	strategy := backoffx.StrategyFunc(func(attempt int) time.Duration {
		once.Do(func() { close(started) })
		return time.Millisecond
	})
	w := worker{wq: q, s: strategy, fn: h}

	ctx, cancel := context.WithCancel(t.Context())
	done := make(chan struct{})
	go func() {
		w.Consume(ctx)
		close(done)
	}()

	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("Consume never attempted to dequeue")
	}

	cancel()

	select {
	case <-done:
		t.Fatal("Consume returned despite an empty queue; ctx-cancellation-while-waiting behavior may have changed, update this test")
	case <-time.After(100 * time.Millisecond):
	}

	// Unblock the still-running goroutine now that the limitation has been observed.
	require.NoError(t, q.Enqueue(entry.Entry("unblock")))
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Consume did not return after being unblocked")
	}
}

func TestNewWorkerWiresDefaults(t *testing.T) {
	q := pqueuetestx.NewQueue()
	require.NoError(t, q.Enqueue(entry.Entry("hello")))
	h := pqueuetestx.NewHandler()
	w := NewWorker(q, h)

	ctx, cancel := context.WithCancel(t.Context())
	done := make(chan struct{})
	go func() {
		w.Consume(ctx)
		close(done)
	}()

	require.Eventually(t, func() bool {
		return h.Count() == 1
	}, time.Second, time.Millisecond)
	require.Equal(t, [][]byte{[]byte("hello")}, h.Snapshot())

	cancel()
	require.NoError(t, q.Enqueue(entry.Entry("unblock")))
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Consume did not return after context cancellation")
	}
}
