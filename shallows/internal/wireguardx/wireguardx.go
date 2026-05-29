package wireguardx

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"net/netip"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/retrovibed/retrovibed/shallows/internal/backoffx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/userx"
	"golang.org/x/text/encoding/unicode"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun/netstack"
)

const (
	DefaultMTU = 1420
	Current    = "_current"
)

type ParseError struct {
	why      string
	offender string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("%s: %q", e.why, e.offender)
}

func parseIPCidr(s string) (netip.Prefix, error) {
	ipcidr, err := netip.ParsePrefix(s)
	if err == nil {
		return ipcidr, nil
	}
	addr, err := netip.ParseAddr(s)
	if err != nil {
		return netip.Prefix{}, &ParseError{"Invalid IP address: ", s}
	}
	return netip.PrefixFrom(addr, addr.BitLen()), nil
}

func parseEndpoint(s string) (*Endpoint, error) {
	i := strings.LastIndexByte(s, ':')
	if i < 0 {
		return nil, &ParseError{"Missing port from endpoint", s}
	}
	host, portStr := s[:i], s[i+1:]
	if len(host) < 1 {
		return nil, &ParseError{"Invalid endpoint host", host}
	}
	port, err := parsePort(portStr)
	if err != nil {
		return nil, err
	}
	hostColon := strings.IndexByte(host, ':')
	if host[0] == '[' || host[len(host)-1] == ']' || hostColon > 0 {
		err := &ParseError{"Brackets must contain an IPv6 address", host}
		if len(host) > 3 && host[0] == '[' && host[len(host)-1] == ']' && hostColon > 0 {
			end := len(host) - 1
			if i := strings.LastIndexByte(host, '%'); i > 1 {
				end = i
			}
			maybeV6, err2 := netip.ParseAddr(host[1:end])
			if err2 != nil || !maybeV6.Is6() {
				return nil, err
			}
		} else {
			return nil, err
		}
		host = host[1 : len(host)-1]
	}
	return &Endpoint{host, port}, nil
}

func parseMTU(s string) (uint16, error) {
	m, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if m < 576 || m > 65535 {
		return 0, &ParseError{"Invalid MTU", s}
	}
	return uint16(m), nil
}

func parsePort(s string) (uint16, error) {
	m, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if m < 0 || m > 65535 {
		return 0, &ParseError{"Invalid port", s}
	}
	return uint16(m), nil
}

func parsePersistentKeepalive(s string) (uint16, error) {
	if s == "off" {
		return 0, nil
	}
	m, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if m < 0 || m > 65535 {
		return 0, &ParseError{"Invalid persistent keepalive", s}
	}
	return uint16(m), nil
}

func parseTableOff(s string) (bool, error) {
	switch s {
	case "off":
		return true, nil
	case "auto", "main":
		return false, nil
	}
	_, err := strconv.ParseUint(s, 10, 32)
	return false, err
}

func parseKeyBase64(s string) (*Key, error) {
	k, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, &ParseError{fmt.Sprintf("Invalid key: %v", err), s}
	}
	if len(k) != KeyLength {
		return nil, &ParseError{"Keys must decode to exactly 32 bytes", s}
	}
	var key Key
	copy(key[:], k)
	return &key, nil
}

func splitList(s string) ([]string, error) {
	var out []string
	for _, split := range strings.Split(s, ",") {
		trim := strings.TrimSpace(split)
		if len(trim) == 0 {
			return nil, &ParseError{"Two commas in a row", s}
		}
		out = append(out, trim)
	}
	return out, nil
}

type parserState int

const (
	inInterfaceSection parserState = iota
	inPeerSection
	notInASection
)

func (c *Config) maybeAddPeer(p *Peer) {
	if p != nil {
		c.Peers = append(c.Peers, *p)
	}
}

// Statistics is a point-in-time snapshot of WireGuard device state.
type Statistics struct {
	Timestamp         time.Time
	PeerKey           string
	KeepaliveInterval uint64
	TXBytes           uint64
	RXBytes           uint64
	LastHandshakeSec  int64
}

// Snapshot reads the current device state via UAPI.
func Snapshot(dev *device.Device) (Statistics, error) {
	uapiData, err := dev.IpcGet()
	if err != nil {
		return Statistics{}, fmt.Errorf("failed to read internal device memory: %w", err)
	}

	s := Statistics{Timestamp: time.Now()}
	scanner := bufio.NewScanner(strings.NewReader(uapiData))
	for scanner.Scan() {
		key, value, ok := strings.Cut(scanner.Text(), "=")
		if !ok {
			continue
		}
		switch key {
		case "public_key":
			s.PeerKey = value
		case "tx_bytes":
			s.TXBytes = errorsx.Zero(strconv.ParseUint(value, 10, 64))
		case "rx_bytes":
			s.RXBytes = errorsx.Zero(strconv.ParseUint(value, 10, 64))
		case "last_handshake_time_sec":
			s.LastHandshakeSec = errorsx.Zero(strconv.ParseInt(value, 10, 64))
		case "persistent_keepalive_interval":
			s.KeepaliveInterval = errorsx.Zero(strconv.ParseUint(value, 10, 64))
		}
	}

	return s, nil
}

// Diagnostic prints the current device state to w.
func Diagnostic(w io.Writer, dev *device.Device) error {
	s, err := Snapshot(dev)
	if err != nil {
		return err
	}

	errorsx.Zero(fmt.Fprintln(w, "=== WireGuard Engine Internal Diagnostics ==="))
	if s.PeerKey == "" {
		errorsx.Zero(fmt.Fprintln(w, "Status: INACTIVE - No peer configurations loaded in memory."))
		return nil
	}

	errorsx.Zero(fmt.Fprintf(w, "Peer Node:  %s\n", s.PeerKey))
	errorsx.Zero(fmt.Fprintf(w, "Keepalive:  %ds\n", s.KeepaliveInterval))
	errorsx.Zero(fmt.Fprintf(w, "Data Sent:  %d bytes\n", s.TXBytes))
	errorsx.Zero(fmt.Fprintf(w, "Data Recv:  %d bytes\n", s.RXBytes))

	if s.LastHandshakeSec == 0 {
		errorsx.Zero(fmt.Fprintln(w, "Handshake:  NEVER COMPLETED (Tunnel initializing or completely blocked)"))
		errorsx.Zero(fmt.Fprintln(w, "Conclusion: Server rate-limiting or firewall block is active."))
		return nil
	}

	elapsed := time.Now().Unix() - s.LastHandshakeSec
	errorsx.Zero(fmt.Fprintf(w, "Last Sync:  %d seconds ago\n", elapsed))

	errorsx.Zero(fmt.Fprintln(w, "--- Analysis ---"))
	if s.TXBytes > 0 && s.RXBytes == 0 {
		errorsx.Zero(fmt.Fprintln(w, "Alert:      Unbalanced Pipe (Data leaving host, but server dropping response)"))
		errorsx.Zero(fmt.Fprintln(w, "Conclusion: VPN active block / rate limit signature matched."))
	} else if elapsed > 180 {
		errorsx.Zero(fmt.Fprintln(w, "Alert:      Stale Handshake (Exceeded 3-minute protocol window)"))
		errorsx.Zero(fmt.Fprintln(w, "Conclusion: Connection dropped by remote endpoint."))
	} else {
		errorsx.Zero(fmt.Fprintln(w, "Status:     Healthy (Bidirectional data flow and valid handshakes)"))
	}

	return nil
}

type autohealState struct {
	prev            Statistics
	attempt         int
	lastRecovery    time.Time
	handshakeExpiry time.Duration
}

func newAutohealState() autohealState {
	return autohealState{
		lastRecovery:    time.Now(),
		handshakeExpiry: 185 * time.Second,
	}
}

// reset records the recovery timestamp, clears the stateful baseline, and
// increments attempt so the backoff grows before the next check.
func (a *autohealState) reset() {
	a.lastRecovery = time.Now()
	a.prev = Statistics{}
	a.attempt++
}

func (a *autohealState) needsRecovery(curr Statistics) bool {
	if curr.PeerKey == "" || curr.LastHandshakeSec == 0 {
		return false
	}
	// suppress until a handshake newer than last recovery/start is seen
	if curr.LastHandshakeSec <= a.lastRecovery.Unix() {
		return false
	}
	// keepalive configured but handshake stale → tunnel dropped
	if curr.KeepaliveInterval > 0 &&
		time.Duration(time.Now().Unix()-curr.LastHandshakeSec)*time.Second > a.handshakeExpiry {
		return true
	}
	// tx growing, rx frozen across two snapshots → server dropping responses
	if curr.TXBytes > a.prev.TXBytes && curr.RXBytes == a.prev.RXBytes && a.prev.RXBytes > 0 {
		return true
	}
	return false
}

// Autoheal monitors the WireGuard device and attempts Down/Up recovery when
// the tunnel appears dead. The backoff strategy controls polling frequency;
// attempt resets to 0 on a healthy check and grows on recovery triggers.
func Autoheal(ctx context.Context, dev *device.Device, b backoffx.Strategy) {
	state := newAutohealState()
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(b.Backoff(state.attempt)):
		}

		curr, err := Snapshot(dev)
		if err != nil {
			log.Println("wireguard autoheal: snapshot failed:", err)
			state.attempt++
			continue
		}

		errorsx.Log(Diagnostic(os.Stderr, dev))

		switch {
		case state.needsRecovery(curr):
			log.Println("wireguard autoheal: tunnel dead, initiating recovery")

			if err := errorsx.Compact(dev.Down(), dev.Up()); err != nil {
				log.Println("wireguard autoheal: recovery failed:", err)
				continue
			}

			state.reset()
		case curr.LastHandshakeSec > state.lastRecovery.Unix():
			state.attempt = 0
			// else: post-recovery settle window, leave attempt unchanged
		}

		state.prev = curr
	}
}

func FormatIPCSet(wcfg *Config) (ipcsets []string) {
	for _, peer := range wcfg.Peers {
		// ensure host is converted to an ip address.
		peer.Endpoint.Host = langx.FirstNonZero(errorsx.Zero(net.ResolveIPAddr("ip", peer.Endpoint.Host)).String(), peer.Endpoint.Host)
		ipcset := fmt.Sprintf("private_key=%x\n", wcfg.Interface.PrivateKey)
		ipcset += fmt.Sprintf("public_key=%x\n", peer.PublicKey)
		ipcset += fmt.Sprintf("preshared_key=%x\n", peer.PresharedKey)
		ipcset += fmt.Sprintf("endpoint=%s\n", peer.Endpoint.String())
		ipcset += fmt.Sprintf("persistent_keepalive_interval=%d\n", peer.PersistentKeepalive)

		for _, ip := range peer.AllowedIPs {
			ipcset += fmt.Sprintf("allowed_ip=%s\n", ip.String())
		}

		ipcsets = append(ipcsets, ipcset)
	}

	return ipcsets
}

func ConfigDirectory(rels ...string) string {
	return userx.DefaultConfigDir(userx.DefaultRelRoot(), "wireguard.d", filepath.Join(rels...))
}

func Latest() string {
	return ConfigDirectory(Current)
}

func Parse(path string) (*Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return FromWgQuick(string(raw), "retrovibed")
}

func FromWgQuick(s, name string) (_ *Config, err error) {
	// if !TunnelNameIsValid(name) {
	// 	return nil, &ParseError{fmt.Sprintf("Tunnel name is not valid"), name}
	// }
	lines := strings.Split(s, "\n")
	parserState := notInASection
	conf := Config{Name: name}
	sawPrivateKey := false
	var peer *Peer
	for _, line := range lines {
		line, _, _ = strings.Cut(line, "#")
		line = strings.TrimSpace(line)
		lineLower := strings.ToLower(line)
		if len(line) == 0 {
			continue
		}
		if lineLower == "[interface]" {
			conf.maybeAddPeer(peer)
			parserState = inInterfaceSection
			continue
		}
		if lineLower == "[peer]" {
			conf.maybeAddPeer(peer)
			peer = &Peer{}
			parserState = inPeerSection
			continue
		}
		if parserState == notInASection {
			return nil, &ParseError{why: "Line must occur in a section", offender: line}
		}
		equals := strings.IndexByte(line, '=')
		if equals < 0 {
			return nil, &ParseError{why: "Config key is missing an equals separator", offender: line}
		}
		key, val := strings.TrimSpace(lineLower[:equals]), strings.TrimSpace(line[equals+1:])
		if len(val) == 0 {
			return nil, &ParseError{why: "Key must have a value", offender: line}
		}
		switch parserState {
		case inInterfaceSection:
			switch key {
			case "privatekey":
				k, err := parseKeyBase64(val)
				if err != nil {
					return nil, err
				}
				conf.Interface.PrivateKey = *k
				sawPrivateKey = true
			case "listenport":
				p, err := parsePort(val)
				if err != nil {
					return nil, err
				}
				conf.Interface.ListenPort = p
			case "mtu":
				m, err := parseMTU(val)
				if err != nil {
					return nil, err
				}
				conf.Interface.MTU = m
			case "address":
				addresses, err := splitList(val)
				if err != nil {
					return nil, err
				}
				for _, address := range addresses {
					a, err := parseIPCidr(address)
					if err != nil {
						return nil, err
					}
					conf.Interface.Addresses = append(conf.Interface.Addresses, a)
				}
			case "dns":
				addresses, err := splitList(val)
				if err != nil {
					return nil, err
				}
				for _, address := range addresses {
					a, err := netip.ParseAddr(address)
					if err != nil {
						conf.Interface.DNSSearch = append(conf.Interface.DNSSearch, address)
					} else {
						conf.Interface.DNS = append(conf.Interface.DNS, a)
					}
				}
			case "preup":
				conf.Interface.PreUp = val
			case "postup":
				conf.Interface.PostUp = val
			case "predown":
				conf.Interface.PreDown = val
			case "postdown":
				conf.Interface.PostDown = val
			case "table":
				tableOff, err := parseTableOff(val)
				if err != nil {
					return nil, err
				}
				conf.Interface.TableOff = tableOff
			default:
				return nil, &ParseError{why: "Invalid key for [Interface] section", offender: key}
			}
		case inPeerSection:
			switch key {
			case "publickey":
				k, err := parseKeyBase64(val)
				if err != nil {
					return nil, err
				}
				peer.PublicKey = *k
			case "presharedkey":
				k, err := parseKeyBase64(val)
				if err != nil {
					return nil, err
				}
				peer.PresharedKey = *k
			case "allowedips":
				addresses, err := splitList(val)
				if err != nil {
					return nil, err
				}
				for _, address := range addresses {
					a, err := parseIPCidr(address)
					if err != nil {
						return nil, err
					}
					peer.AllowedIPs = append(peer.AllowedIPs, a)
				}
			case "persistentkeepalive":
				p, err := parsePersistentKeepalive(val)
				if err != nil {
					return nil, err
				}
				peer.PersistentKeepalive = p
			case "endpoint":
				e, err := parseEndpoint(val)
				if err != nil {
					return nil, err
				}
				peer.Endpoint = *e
			default:
				return nil, &ParseError{why: "Invalid key for [Peer] section", offender: key}
			}
		}
	}
	conf.maybeAddPeer(peer)

	if !sawPrivateKey {
		return nil, &ParseError{why: "An interface must have a private key", offender: "[none specified]"}
	}
	for _, p := range conf.Peers {
		if p.PublicKey.IsZero() {
			return nil, &ParseError{why: "All peers must have public keys", offender: "[none specified]"}
		}
	}

	return &conf, nil
}

func FromWgQuickWithUnknownEncoding(s, name string) (*Config, error) {
	c, firstErr := FromWgQuick(s, name)
	if firstErr == nil {
		return c, nil
	}
	for _, encoding := range unicode.All {
		decoded, err := encoding.NewDecoder().String(s)
		if err == nil {
			c, err := FromWgQuick(decoded, name)
			if err == nil {
				return c, nil
			}
		}
	}
	return nil, firstErr
}

func HostLookupAdapter(wnet *netstack.Net) dnsresolver {
	return dnsresolver{Net: wnet}
}

type dnsresolver struct {
	*netstack.Net
}

func (t dnsresolver) LookupHost(ctx context.Context, host string) (addrs []string, err error) {
	return t.LookupContextHost(ctx, host)
}
