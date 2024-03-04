package manager

import (
	"net"
	"net/url"
	"slices"
	"strings"

	"github.com/seal-io/utils/netx"
	"k8s.io/client-go/rest"

	"github.com/seal-io/walrus/pkg/system"
)

func isLoopbackClusterNearby(restCfg *rest.Config) bool {
	// Extract host from rest config.
	var host string
	if strings.Contains(restCfg.Host, "://") {
		u, _ := url.Parse(restCfg.Host)
		host = u.Host
	} else {
		host = restCfg.Host
	}
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	} else if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}

	// Detect host in a fast pass way.
	knownLoopbackHosts := []string{
		"kubernetes.docker.internal",
		"host.docker.internal",
		"localhost",
		"127.0.0.1",
		"[::1]",
		"[::1%lo0]",
	}
	if slices.Contains(knownLoopbackHosts, host) {
		return true
	}

	// Detect host in a slow pass way.
	subnets := make([]netx.IPv4, 0, system.Subnets.Get().Len())
	for _, v := range system.Subnets.Get().List() {
		sn := netx.MustIPv4FromCIDR(v)
		subnets = append(subnets, sn)
	}

	// IP detect.
	if ip := net.ParseIP(host); ip != nil {
		for j := range subnets {
			if subnets[j].Contains(ip) {
				return true
			}
		}

		return false
	}

	// Or DNS lookup.
	ips, err := net.LookupIP(host)
	if err != nil {
		return false
	}

	for i := range ips {
		if ips[i].IsLoopback() {
			return true
		}
		for j := range subnets {
			if subnets[j].Contains(ips[i]) {
				return true
			}
		}
	}

	return false
}
