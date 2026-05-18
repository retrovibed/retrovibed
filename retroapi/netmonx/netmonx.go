package netmonx

import (
	"context"
	"iter"
	"sync"

	"tailscale.com/net/netmon"
)

type Monitor struct {
	mon  *netmon.Monitor
	ch   chan *netmon.Delta
	done chan struct{}
	once sync.Once
	err  error
}

func New() (*Monitor, error) {
	m := &Monitor{
		ch:   make(chan *netmon.Delta, 16),
		done: make(chan struct{}),
	}

	mon, err := netmon.New(func(n *netmon.Delta) {
		select {
		case m.ch <- n:
		case <-m.done:
		}
	})
	if err != nil {
		return nil, err
	}

	m.mon = mon
	mon.Start()
	return m, nil
}

func (m *Monitor) Each(ctx context.Context) iter.Seq[*netmon.Delta] {
	return func(yield func(*netmon.Delta) bool) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-m.done:
				return
			case v, ok := <-m.ch:
				if !ok {
					return
				}
				if !yield(v) {
					return
				}
			}
		}
	}
}

func (m *Monitor) Err() error {
	return m.err
}

func (m *Monitor) Close() error {
	m.once.Do(func() {
		close(m.done)
		close(m.ch)
	})
	return m.mon.Close()
}
