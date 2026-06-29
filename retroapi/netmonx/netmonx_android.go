//go:build android

package netmonx

import (
	"bufio"
	"context"
	"errors"
	"math"
	"os"
	"strconv"
	"strings"
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

// /proc/net/route exposes only the device's own routing table, which
// remains readable by apps under Android's procfs restrictions.
func defaultRouteInterface() string {
	f, err := os.Open("/proc/net/route")
	if err != nil {
		return ""
	}
	defer f.Close()

	best := ""
	bestMetric := math.MaxInt32

	scanner := bufio.NewScanner(f)
	// Skip header line.
	scanner.Scan()
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 8 {
			continue
		}
		// Destination is hex little-endian; 00000000 = default route.
		if fields[1] != "00000000" {
			continue
		}
		metric, err := strconv.Atoi(fields[6])
		if err != nil {
			continue
		}
		if best == "" || metric < bestMetric {
			best = fields[0]
			bestMetric = metric
		}
	}

	return best
}
