package netmonx

import (
	"context"
	"iter"
	"log"
	"net"
	"net/netip"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/internal/langx"
)

var (
	globalMon  atomic.Pointer[Monitor]
	globalOnce sync.Once
)

func Global() *Monitor {
	globalOnce.Do(func() {
		m, err := New()
		if err != nil {
			log.Println("failed to initialize global monitor", err)
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

func Metered() bool {
	m := Global()
	if m == nil {
		return true // assume worst case.
	}
	return m.Metered()
}

// InterfaceDetails holds the addresses and metered status of a single interface.
type InterfaceDetails struct {
	Name    string
	IPs     []netip.Prefix
	Metered bool
}

// State represents a snapshot of network interface state.
type State struct {
	Interfaces            []InterfaceDetails
	HaveV4                bool
	HaveV6                bool
	DefaultRouteInterface string
}

// ChangeDelta is emitted each time the network state changes.
type ChangeDelta struct {
	Old                     *State // nil when IsInitialState is true
	New                     *State
	DefaultInterfaceChanged bool
	InterfaceIPsChanged     bool
	IsInitialState          bool
}

type Monitor struct {
	current atomic.Pointer[State]
	metered atomic.Pointer[bool]
	events  chan ChangeDelta
	notify  chan struct{}
	done    chan struct{}
	once    sync.Once
	err     error
}

func New() (*Monitor, error) {
	m := &Monitor{
		events: make(chan ChangeDelta, 16),
		notify: make(chan struct{}, 1),
		done:   make(chan struct{}),
	}

	s, err := getState()
	if err != nil {
		return nil, err
	}
	m.current.Store(s)

	ctx, cancel := context.WithCancel(context.Background())
	if err := startWatcher(ctx, m.notify); err != nil {
		log.Println("netmonx: native watcher unavailable, using polling:", err)
		startPoll(ctx, m.notify)
	}

	go m.run(cancel)
	return m, nil
}

func (m *Monitor) run(cancel context.CancelFunc) {
	defer cancel()
	defer close(m.events)

	// Emit initial state.
	initial := m.current.Load()
	select {
	case m.events <- ChangeDelta{New: initial, IsInitialState: true}:
	case <-m.done:
		return
	}

	for {
		select {
		case <-m.done:
			return
		case <-m.notify:
		}
		m.refresh()
	}
}

func startPoll(ctx context.Context, notify chan<- struct{}) {
	go func() {
		t := time.NewTicker(30 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				select {
				case notify <- struct{}{}:
				default:
				}
			}
		}
	}()
}

func (m *Monitor) refresh() {
	newState, err := getState()
	if err != nil {
		return
	}

	old := m.current.Load()
	if statesEqual(old, newState) {
		return
	}

	m.current.Store(newState)

	delta := ChangeDelta{
		Old:                     old,
		New:                     newState,
		DefaultInterfaceChanged: old.DefaultRouteInterface != newState.DefaultRouteInterface,
		InterfaceIPsChanged:     !interfacesEqual(old.Interfaces, newState.Interfaces),
	}

	select {
	case m.events <- delta:
	case <-m.done:
	}
}

func (m *Monitor) Each(ctx context.Context) iter.Seq[ChangeDelta] {
	return func(yield func(ChangeDelta) bool) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-m.done:
				return
			case v, ok := <-m.events:
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

func (m *Monitor) Metered() bool {
	if metered := langx.Zero(m.metered.Load()); metered {
		return true
	}

	s := m.current.Load()
	if s == nil {
		return true
	}

	for _, iface := range s.Interfaces {
		if iface.Name == s.DefaultRouteInterface {
			return iface.Metered
		}
	}

	return false
}

func (m *Monitor) SetMetered(b bool) {
	m.metered.Store(&b)

	// Signal background goroutine to emit a delta.
	select {
	case m.notify <- struct{}{}:
	default:
	}
}

func (m *Monitor) Err() error {
	return m.err
}

func (m *Monitor) Close() error {
	m.once.Do(func() {
		close(m.done)
	})
	return m.err
}

func getState() (*State, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	s := &State{}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		var prefixes []netip.Prefix
		for _, addr := range addrs {
			var prefix netip.Prefix
			switch v := addr.(type) {
			case *net.IPNet:
				ip, ok := netip.AddrFromSlice(v.IP)
				if !ok {
					continue
				}
				bits, _ := v.Mask.Size()
				prefix = netip.PrefixFrom(ip.Unmap(), bits)
			case *net.IPAddr:
				ip, ok := netip.AddrFromSlice(v.IP)
				if !ok {
					continue
				}
				prefix = netip.PrefixFrom(ip.Unmap(), ip.BitLen())
			}
			if !prefix.IsValid() {
				continue
			}
			a := prefix.Addr()
			if a.IsLoopback() || a.IsLinkLocalUnicast() {
				continue
			}
			prefixes = append(prefixes, prefix)
			if a.Is4() {
				s.HaveV4 = true
			} else if a.Is6() {
				s.HaveV6 = true
			}
		}

		if len(prefixes) > 0 {
			s.Interfaces = append(s.Interfaces, InterfaceDetails{
				Name:    iface.Name,
				IPs:     prefixes,
				Metered: langx.FirstNonZero(platformMetered(iface.Name), isMeteredInterface(iface.Name)),
			})
		}
	}

	s.DefaultRouteInterface = defaultRouteInterface()

	return s, nil
}

func isMeteredInterface(name string) bool {
	return strings.HasPrefix(name, "wwan") ||
		strings.HasPrefix(name, "rmnet") ||
		strings.HasPrefix(name, "pdp_ip") ||
		strings.HasPrefix(name, "ccmni")
}

func statesEqual(a, b *State) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.HaveV4 == b.HaveV4 &&
		a.HaveV6 == b.HaveV6 &&
		a.DefaultRouteInterface == b.DefaultRouteInterface &&
		interfacesEqual(a.Interfaces, b.Interfaces)
}

func interfacesEqual(a, b []InterfaceDetails) bool {
	if len(a) != len(b) {
		return false
	}
	bByName := make(map[string]InterfaceDetails, len(b))
	for _, d := range b {
		bByName[d.Name] = d
	}
	for _, da := range a {
		db, ok := bByName[da.Name]
		if !ok || da.Metered != db.Metered || !prefixSetsEqual(da.IPs, db.IPs) {
			return false
		}
	}
	return true
}

func prefixSetsEqual(a, b []netip.Prefix) bool {
	if len(a) != len(b) {
		return false
	}
	set := make(map[netip.Prefix]bool, len(a))
	for _, p := range a {
		set[p] = true
	}
	for _, p := range b {
		if !set[p] {
			return false
		}
	}
	return true
}
