package netx

import (
	"log"
	"net"
	"strconv"
)

func DefaultIfZero(fallback net.IP, v net.IP) net.IP {
	if v != nil {
		return v
	}

	return fallback
}

// HostIP ...
func HostIP(host string) net.IP {
	ip, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		log.Println("failed to resolve ip for", host, "falling back to 127.0.0.1:", err)
		return net.ParseIP("127.0.0.1")
	}

	return ip.IP
}

func Port(s string) (p uint16, err error) {
	var (
		sport string
		port  uint64
	)

	if _, sport, err = net.SplitHostPort(s); err != nil {
		return 0, err
	}

	if port, err = strconv.ParseUint(sport, 10, 16); err != nil {
		return 0, err
	}

	return uint16(port), nil
}

func IP(s string) net.IP {
	var (
		err  error
		host string
	)

	if host, _, err = net.SplitHostPort(s); err != nil {
		log.Println("unable to parse host", host)
		return nil
	}

	return HostIP(host)
}
