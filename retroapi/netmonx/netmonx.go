package netmonx

import (
	"context"
	"iter"
	"sync"
	"sync/atomic"

	"github.com/retrovibed/retrovibed/retroapi/atomicx"
	"tailscale.com/net/netmon"
	"tailscale.com/types/logger"
	"tailscale.com/util/eventbus"
)

var (
	globalMon  atomic.Pointer[Monitor]
	globalOnce sync.Once
)

func Global() *Monitor {
	globalOnce.Do(func() {
		m, err := New()
		if err != nil {
			return
		}
		globalMon.Store(m)
	})
	return globalMon.Load()
}

func SetMetered(b bool) {
	m := Global()
	if m == nil {
		return
	}
	m.SetMetered(b)
}

func (m *Monitor) SetMetered(b bool) {
	s := m.current.Load()
	s.IsExpensive = b
	m.current.Store(s)
	m.mon.InjectEvent()
}

func Metered() bool {
	m := Global()
	if m == nil {
		return false
	}

	return m.Metered()
}

type Monitor struct {
	bus     *eventbus.Bus
	mon     *netmon.Monitor
	current *atomic.Pointer[netmon.State]
	ch      chan netmon.ChangeDelta
	done    chan struct{}
	once    sync.Once
	err     error
}

func New() (*Monitor, error) {
	m := &Monitor{
		bus:     eventbus.New(),
		ch:      make(chan netmon.ChangeDelta, 16),
		done:    make(chan struct{}),
		current: atomicx.Pointer(netmon.State{}),
	}

	mon, err := netmon.New(m.bus, logger.Discard)
	if err != nil {
		m.bus.Close()
		return nil, err
	}

	mon.RegisterChangeCallback(func(delta *netmon.ChangeDelta) {
		s := delta.CurrentState()
		s.IsExpensive = m.current.Load().IsExpensive
		m.current.Store(s)

		select {
		case m.ch <- *delta:
		case <-m.done:
		}
	})

	m.mon = mon
	mon.Start()

	if st := mon.InterfaceState(); st != nil {
		m.current.Store(st)
	}

	return m, nil
}

func (m *Monitor) Each(ctx context.Context) iter.Seq[netmon.ChangeDelta] {
	return func(yield func(netmon.ChangeDelta) bool) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-m.done:
				return
			case v := <-m.ch:
				if !yield(v) {
					return
				}
			}
		}
	}
}

func (m *Monitor) Metered() bool {
	return m.current.Load().IsExpensive
}

func (m *Monitor) Err() error {
	return m.err
}

func (m *Monitor) Close() error {
	var err error
	m.once.Do(func() {
		close(m.done)
		err = m.mon.Close()
		m.bus.Close()
		m.err = err
	})
	return err
}
