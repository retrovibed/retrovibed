//go:build linux && !android

package netmonx

import (
	"bufio"
	"context"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/mdlayher/netlink"
	"golang.org/x/sys/unix"
)

func startWatcher(ctx context.Context, notify chan<- struct{}) error {
	conn, err := netlink.Dial(unix.NETLINK_ROUTE, &netlink.Config{
		Groups: unix.RTMGRP_IPV4_IFADDR | unix.RTMGRP_IPV6_IFADDR |
			unix.RTMGRP_IPV4_ROUTE | unix.RTMGRP_IPV6_ROUTE | unix.RTMGRP_LINK,
	})
	if err != nil {
		return err
	}

	// Close connection when context is cancelled to unblock Receive.
	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	go func() {
		for {
			_, err := conn.Receive()
			if err != nil {
				return
			}
			select {
			case notify <- struct{}{}:
			default:
			}
		}
	}()

	return nil
}

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

func getFallbackState() *State { return nil }

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
