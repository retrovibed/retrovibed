// Package pqueuetestx provides reusable test doubles and assertion helpers
// for code depending on pqueue.Queue (github.com/linxGnu/pqueue) and
// pqueuex.Handler (shallows/internal/pqueuex). Modeled on the assertion
// style of deeppool's nsqxtest.
package pqueuetestx

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"sync"
	"testing"

	"github.com/linxGnu/pqueue"
	"github.com/linxGnu/pqueue/entry"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/stretchr/testify/require"
)

var _ pqueue.Queue = (*Queue)(nil)

// Queue is an in-memory pqueue.Queue for tests. The real pqueue.New is
// file-backed and has no in-memory constructor, but callers only ever see
// the pqueue.Queue interface, so a hand-rolled fake is sufficient.
type Queue struct {
	mu    *sync.Mutex
	items [][]byte
}

// NewQueue returns an empty, ready-to-use in-memory Queue.
func NewQueue() *Queue {
	return &Queue{mu: &sync.Mutex{}}
}

func (q *Queue) Close() error { return nil }

func (q *Queue) Enqueue(e entry.Entry) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = append(q.items, append([]byte(nil), e...))
	return nil
}

// EnqueueBatch is a documented no-op: entry.Batch (vendor/github.com/linxGnu/pqueue/entry)
// stores its entries in an unexported field with no exported enumeration
// method, so a faithful in-memory fake cannot recover the batch's payloads
// without reimplementing entry.Batch's wire format. Nothing in this repo
// currently calls EnqueueBatch; if a caller starts to depend on it, this
// must be revisited.
func (q *Queue) EnqueueBatch(_ entry.Batch) error { return nil }

func (q *Queue) Dequeue(e *entry.Entry) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return false
	}
	*e = entry.Entry(q.items[0])
	q.items = q.items[1:]
	return true
}

func (q *Queue) Peek(e *entry.Entry) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return false
	}
	*e = entry.Entry(q.items[0])
	return true
}

// Len returns the number of entries currently backlogged (not yet dequeued).
func (q *Queue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}

// Snapshot returns a defensive, ordered copy of the entries currently
// backlogged (FIFO order: index 0 is next to be dequeued).
func (q *Queue) Snapshot() [][]byte {
	q.mu.Lock()
	defer q.mu.Unlock()
	out := make([][]byte, len(q.items))
	copy(out, q.items)
	return out
}

// HaveLength builds a LengthMatcher against the queue's current backlog.
func (q *Queue) HaveLength(n int) LengthMatcher {
	return LengthMatcher{messages: q.Snapshot(), Expected: n}
}

// HaveEnqueued builds a ContentMatcher checking whether b is present
// anywhere in the current backlog (order-independent "contains").
func (q *Queue) HaveEnqueued(b []byte) ContentMatcher {
	return ContentMatcher{messages: q.Snapshot(), Expected: b}
}

// Handler records every payload passed to Message and can be configured to
// fail on specific (0-indexed) calls via HandlerOptionErrorAt. It structurally
// satisfies pqueuex.Handler (Message(ctx context.Context, m []byte) error)
// without importing pqueuex, since pqueuex's own internal tests need to
// import this package, and pqueuex importing pqueuetestx importing pqueuex
// back would be a real import cycle.
type Handler struct {
	mu       *sync.Mutex
	received [][]byte
	errAt    map[int]error
}

// HandlerOption configures a Handler at construction time.
type HandlerOption func(*Handler)

// HandlerOptionErrorAt causes the idx'th (0-indexed) call to Message to
// return err instead of nil.
func HandlerOptionErrorAt(idx int, err error) HandlerOption {
	return func(h *Handler) { h.errAt[idx] = err }
}

// NewHandler returns a ready-to-use Handler, applying options via the same
// langx.Clone pattern pqueuex.NewWorker uses.
func NewHandler(options ...HandlerOption) *Handler {
	h := langx.Clone(Handler{mu: &sync.Mutex{}, errAt: map[int]error{}}, options...)
	return &h
}

func (h *Handler) Message(_ context.Context, m []byte) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	idx := len(h.received)
	h.received = append(h.received, append([]byte(nil), m...))
	return h.errAt[idx]
}

// Count returns the number of Message calls observed so far.
func (h *Handler) Count() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.received)
}

// Snapshot returns a defensive, ordered copy of every payload passed to
// Message so far, in call order.
func (h *Handler) Snapshot() [][]byte {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make([][]byte, len(h.received))
	copy(out, h.received)
	return out
}

// HaveLength builds a LengthMatcher against the handler's received count.
func (h *Handler) HaveLength(n int) LengthMatcher {
	return LengthMatcher{messages: h.Snapshot(), Expected: n}
}

// HaveReceived builds a ContentMatcher checking whether b was received at
// any point (order-independent "contains").
func (h *Handler) HaveReceived(b []byte) ContentMatcher {
	return ContentMatcher{messages: h.Snapshot(), Expected: b}
}

// HaveReceivedInOrder builds an OrderedMatcher checking that the exact
// sequence of received payloads equals want, in order. Unlike nsqx (a
// multi-producer, per-topic "contains" model), pqueue is a single FIFO/
// priority queue drained by one worker loop, so exact-sequence equality is
// the more natural default assertion here.
func (h *Handler) HaveReceivedInOrder(want ...[]byte) OrderedMatcher {
	return OrderedMatcher{messages: h.Snapshot(), Expected: want}
}

// LengthMatcher asserts a snapshot has an expected number of entries.
type LengthMatcher struct {
	messages [][]byte
	Expected int
}

// Match reports whether the snapshot's length equals Expected.
func (m LengthMatcher) Match() (bool, error) { return len(m.messages) == m.Expected, nil }

// FailureMessage describes why Match returned false.
func (m LengthMatcher) FailureMessage() string {
	return fmt.Sprintf("expected %d messages, have %d", m.Expected, len(m.messages))
}

// NegatedFailureMessage describes why a negated Match returned true.
func (m LengthMatcher) NegatedFailureMessage() string {
	return fmt.Sprintf("expected not to have %d messages", m.Expected)
}

// ContentMatcher asserts a snapshot contains a given payload anywhere,
// order-independent.
type ContentMatcher struct {
	messages [][]byte
	Expected []byte
}

// Match reports whether Expected is present anywhere in the snapshot.
func (m ContentMatcher) Match() (bool, error) {
	for _, msg := range m.messages {
		if bytes.Equal(msg, m.Expected) {
			return true, nil
		}
	}
	return false, nil
}

// FailureMessage describes why Match returned false.
func (m ContentMatcher) FailureMessage() string {
	return fmt.Sprintf("expected messages to contain:\n%v\ninstead contains:\n%v", m.Expected, m.messages)
}

// NegatedFailureMessage describes why a negated Match returned true.
func (m ContentMatcher) NegatedFailureMessage() string {
	return fmt.Sprintf("expected messages to not contain %v", m.Expected)
}

// OrderedMatcher asserts a snapshot equals an exact expected sequence.
type OrderedMatcher struct {
	messages [][]byte
	Expected [][]byte
}

// Match reports whether the snapshot equals Expected exactly, in order.
func (m OrderedMatcher) Match() (bool, error) { return reflect.DeepEqual(m.messages, m.Expected), nil }

// FailureMessage describes why Match returned false.
func (m OrderedMatcher) FailureMessage() string {
	return fmt.Sprintf("expected messages in order:\n%v\ninstead got:\n%v", m.Expected, m.messages)
}

// NegatedFailureMessage describes why a negated Match returned true.
func (m OrderedMatcher) NegatedFailureMessage() string {
	return fmt.Sprintf("expected messages to not equal, in order, %v", m.Expected)
}

// NewDisk returns a real, disk-backed pqueue.Queue rooted in a fresh
// t.TempDir(), failing the test via require.NoError if construction fails,
// and closing the queue automatically via t.Cleanup. This duplicates
// pqueuex.New's two lines rather than calling it directly, since pqueuex
// cannot be imported here (see the Handler doc comment above).
func NewDisk(t testing.TB) pqueue.Queue {
	t.Helper()
	dir := t.TempDir()
	require.NoError(t, fsx.MkDirs(0700, dir))
	q, err := pqueue.New(dir, 32)
	require.NoError(t, err)
	t.Cleanup(func() { _ = q.Close() })
	return q
}
