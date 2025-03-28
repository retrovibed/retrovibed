package torrent

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/james-lawrence/torrent/dht/krpc"

	"github.com/james-lawrence/torrent/tracker"
)

// Announces a torrent to a tracker at regular intervals, when peers are
// required.
type trackerScraper struct {
	u            url.URL
	t            *torrent
	lastAnnounce trackerAnnounceResult
}

func (ts *trackerScraper) statusLine() string {
	var w bytes.Buffer
	fmt.Fprintf(&w, "%q\t%s\t%s",
		ts.u.String(),
		func() string {
			na := time.Until(ts.lastAnnounce.Completed.Add(ts.lastAnnounce.Interval))
			if na > 0 {
				na /= time.Second
				na *= time.Second
				return na.String()
			} else {
				return "anytime"
			}
		}(),
		func() string {
			if ts.lastAnnounce.Err != nil {
				return ts.lastAnnounce.Err.Error()
			}
			if ts.lastAnnounce.Completed.IsZero() {
				return "never"
			}
			return fmt.Sprintf("%d peers", ts.lastAnnounce.NumPeers)
		}(),
	)
	return w.String()
}

type trackerAnnounceResult struct {
	Err       error
	NumPeers  int
	Interval  time.Duration
	Completed time.Time
}

func (me *trackerScraper) getIp() (ip net.IP, err error) {
	ips, err := net.LookupIP(me.u.Hostname())
	if err != nil {
		return
	}
	if len(ips) == 0 {
		err = errors.New("no ips")
		return
	}
	for _, ip = range ips {
		switch me.u.Scheme {
		case "udp4":
			if ip.To4() == nil {
				continue
			}
		case "udp6":
			if ip.To4() != nil {
				continue
			}
		}
		return
	}
	err = errors.New("no acceptable ips")
	return
}

func (me *trackerScraper) trackerUrl(ip net.IP) string {
	u := me.u
	if u.Port() != "" {
		u.Host = net.JoinHostPort(ip.String(), u.Port())
	}
	return u.String()
}

// Return how long to wait before trying again. For most errors, we return 5
// minutes, a relatively quick turn around for DNS changes.
func (me *trackerScraper) announce(event tracker.AnnounceEvent) (ret trackerAnnounceResult) {
	defer func() {
		ret.Completed = time.Now()
	}()
	ret.Interval = time.Minute
	ip, err := me.getIp()
	if err != nil {
		ret.Err = fmt.Errorf("error getting ip: %s", err)
		return
	}
	me.t.lock()
	req := me.t.announceRequest(event)
	me.t.unlock()
	//log.Printf("announcing %s %s to %q", me.t, req.Event, me.u.String())
	res, err := tracker.Announce{
		HTTPProxy:  me.t.config.HTTPProxy,
		UserAgent:  me.t.config.HTTPUserAgent,
		TrackerUrl: me.trackerUrl(ip),
		Request:    req,
		HostHeader: me.u.Host,
		ServerName: me.u.Hostname(),
		UdpNetwork: me.u.Scheme,
		ClientIp4:  krpc.NewNodeAddrFromIPPort(me.t.config.PublicIP4, 0),
		ClientIp6:  krpc.NewNodeAddrFromIPPort(me.t.config.PublicIP6, 0),
	}.Do()
	if err != nil {
		ret.Err = fmt.Errorf("error announcing: %s", err)
		return
	}
	me.t.AddPeers(Peers(nil).AppendFromTracker(res.Peers))
	ret.NumPeers = len(res.Peers)
	ret.Interval = time.Duration(res.Interval) * time.Second
	return
}

func (me *trackerScraper) Run() {
	defer me.announceStopped()
	// make sure first announce is a "started"
	e := tracker.Started
	for {
		ar := me.announce(e)
		// after first announce, get back to regular "none"
		e = tracker.None
		me.t.lock()
		me.lastAnnounce = ar
		me.t.unlock()

	wait:
		interval := ar.Interval
		if interval < time.Minute {
			interval = time.Minute
		}
		wantPeers := me.t.wantPeersEvent.LockedChan(me.t.locker())
		select {
		case <-wantPeers:
			if interval > time.Minute {
				interval = time.Minute
			}
			wantPeers = nil
		default:
		}

		select {
		case <-me.t.closed.LockedChan(me.t.locker()):
			return
		case <-wantPeers:
			goto wait
		case <-time.After(time.Until(ar.Completed.Add(interval))):
		}
	}
}

func (me *trackerScraper) announceStopped() {
	me.announce(tracker.Stopped)
}
