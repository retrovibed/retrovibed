//go:build android

package netmonx

import (
	"bufio"
	"context"
	"encoding/hex"
	"errors"
	"net/netip"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"github.com/retrovibed/retrovibed/retroapi/internal/langx"
)

// Android's SELinux policy denies untrusted apps the netlink_route_socket
// permission, so a raw AF_NETLINK watcher (as used on plain Linux) isn't
// available here; New() falls back to polling instead.
func startWatcher(_ context.Context, _ chan<- struct{}) error {
	return errors.ErrUnsupported
}

// /sys/class/net/<iface>/uevent is still readable by apps since it only
// describes the device's own interfaces, not other processes' sockets.
func platformMetered(name string) bool {
	f, err := os.Open("/sys/class/net/" + name + "/uevent")
	if err != nil {
		return false
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if scanner.Text() == "DEVTYPE=wwan" {
			return true
		}
	}
	return false
}

// defaultRouteInterface connects a UDP socket to a public address without
// sending any packets to determine which interface the kernel selects for
// the default route. This avoids /proc/net/route which Android SELinux blocks.
func defaultRouteInterface() string {
	sock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		return ""
	}
	defer syscall.Close(sock)

	// UDP connect is purely local — no packets are sent.
	remote := syscall.SockaddrInet4{Port: 53, Addr: [4]byte{8, 8, 8, 8}}
	if err := syscall.Connect(sock, &remote); err != nil {
		return ""
	}

	sa, err := syscall.Getsockname(sock)
	if err != nil {
		return ""
	}
	local, ok := sa.(*syscall.SockaddrInet4)
	if !ok {
		return ""
	}

	return ifaceByIPv4(local.Addr)
}

// getFallbackState enumerates network state using only AF_INET socket ioctls,
// avoiding AF_NETLINK, /proc/net, and /sys/class/net which Android SELinux
// blocks for untrusted apps.
func getFallbackState() *State {
	sock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		return &State{}
	}
	defer syscall.Close(sock)

	buf := make([]byte, 64*ifreqSize)
	n, err := ioctlIfconf(sock, buf)
	if err != nil {
		return &State{}
	}

	s := &State{DefaultRouteInterface: defaultRouteInterface()}
	v6 := ifInet6Addrs()

	for i := 0; i+ifreqSize <= n; i += ifreqSize {
		entry := buf[i : i+ifreqSize]
		name := strings.TrimRight(string(entry[:16]), "\x00")

		var ip4 [4]byte
		copy(ip4[:], entry[20:24])
		addr := netip.AddrFrom4(ip4)

		if addr.IsLoopback() || addr.IsLinkLocalUnicast() || addr.IsUnspecified() {
			continue
		}

		bits := ifNetmaskBits(sock, name)
		metered := langx.FirstNonZero(platformMetered(name), isMeteredInterface(name))
		s.HaveV4 = true
		s.Networks = append(s.Networks, NetworkDetails{Name: name, IP: netip.PrefixFrom(addr, bits), Metered: metered})
	}

	for name, prefixes := range v6 {
		for _, prefix := range prefixes {
			a := prefix.Addr()
			if a.IsLoopback() || a.IsLinkLocalUnicast() {
				continue
			}
			metered := langx.FirstNonZero(platformMetered(name), isMeteredInterface(name))
			s.HaveV6 = true
			s.Networks = append(s.Networks, NetworkDetails{Name: name, IP: prefix, Metered: metered})
		}
	}

	return s
}

// ifaceByIPv4 returns the interface name that owns the given IPv4 address.
func ifaceByIPv4(target [4]byte) string {
	sock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		return ""
	}
	defer syscall.Close(sock)

	buf := make([]byte, 64*ifreqSize)
	n, err := ioctlIfconf(sock, buf)
	if err != nil {
		return ""
	}

	for i := 0; i+ifreqSize <= n; i += ifreqSize {
		entry := buf[i : i+ifreqSize]
		var ip4 [4]byte
		copy(ip4[:], entry[20:24])
		if ip4 == target {
			return strings.TrimRight(string(entry[:16]), "\x00")
		}
	}
	return ""
}

// ifreqSize is sizeof(struct ifreq) on 64-bit Linux/Android:
// 16-byte name + 24-byte union = 40 bytes.
const ifreqSize = 40

// ifreq is the kernel struct for socket-level network ioctls.
type ifreq [ifreqSize]byte

const (
	siocgifconf    uintptr = 0x8912
	siocgifnetmask uintptr = 0x891b
)

// ioctlIfconf calls SIOCGIFCONF to enumerate all IPv4 interfaces into buf.
// Returns the number of bytes written.
//
// struct ifconf on 64-bit Linux: int32 ifc_len (4B) + pad (4B) + *char ifc_buf (8B).
func ioctlIfconf(sock int, buf []byte) (int, error) {
	var ifc [16]byte
	*(*int32)(unsafe.Pointer(&ifc[0])) = int32(len(buf))
	*(*uintptr)(unsafe.Pointer(&ifc[8])) = uintptr(unsafe.Pointer(&buf[0]))
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(sock), siocgifconf, uintptr(unsafe.Pointer(&ifc[0])))
	runtime.KeepAlive(buf)
	if errno != 0 {
		return 0, errno
	}
	return int(*(*int32)(unsafe.Pointer(&ifc[0]))), nil
}

func ioctlIfreq(sock int, req uintptr, ifr *ifreq) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(sock), req, uintptr(unsafe.Pointer(ifr)))
	if errno != 0 {
		return errno
	}
	return nil
}

// ifNetmaskBits returns the prefix length for an interface's IPv4 netmask.
func ifNetmaskBits(sock int, name string) int {
	var req ifreq
	copy(req[:], name)
	if err := ioctlIfreq(sock, siocgifnetmask, &req); err != nil {
		return 32
	}
	var n int
	for _, b := range req[20:24] {
		for b != 0 {
			n += int(b & 1)
			b >>= 1
		}
	}
	return n
}

// ifInet6Addrs parses /proc/net/if_inet6 for per-interface IPv6 prefixes.
// This file remains readable by untrusted apps under Android's procfs policy.
func ifInet6Addrs() map[string][]netip.Prefix {
	result := make(map[string][]netip.Prefix)
	f, err := os.Open("/proc/net/if_inet6")
	if err != nil {
		return result
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 6 {
			continue
		}
		b, err := hex.DecodeString(fields[0])
		if err != nil || len(b) != 16 {
			continue
		}
		addr, ok := netip.AddrFromSlice(b)
		if !ok {
			continue
		}
		plen, err := strconv.ParseUint(fields[2], 16, 8)
		if err != nil {
			continue
		}
		prefix := netip.PrefixFrom(addr.Unmap(), int(plen))
		if prefix.IsValid() {
			result[fields[5]] = append(result[fields[5]], prefix)
		}
	}
	return result
}
