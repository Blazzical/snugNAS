// Package mdns advertises snugNAS on the LAN via multicast DNS so that
// phones and other devices can resolve `<hostname>.local` to this host's IP
// without any router/DHCP/DNS configuration.
//
// On Windows the responder runs in our own process — Windows has no built-in
// mDNS responder, so the binary must be running for `snugnas.local` to
// resolve from other devices.
package mdns

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/grandcat/zeroconf"
)

// Publisher holds an active mDNS registration.
type Publisher struct {
	Hostname string
	IPs      []string
	server   *zeroconf.Server
}

// Publish starts an mDNS responder that announces `<hostname>.local` on this
// host's LAN IPv4 addresses, plus an `_http._tcp` service record on
// dashboardPort so Bonjour browsers can discover it. Caller must Shutdown.
//
// hostname may include a trailing ".local" or "." — both are stripped before
// publishing because zeroconf appends the domain itself.
func Publish(hostname string, dashboardPort int) (*Publisher, error) {
	if hostname == "" {
		return nil, errors.New("hostname is required")
	}
	name := canonicalHostname(hostname)
	if name == "" {
		return nil, fmt.Errorf("invalid hostname %q", hostname)
	}

	ips, err := lanIPv4()
	if err != nil {
		return nil, err
	}
	if len(ips) == 0 {
		return nil, errors.New("no non-loopback IPv4 interfaces — cannot publish mDNS")
	}

	ipStrs := make([]string, 0, len(ips))
	for _, ip := range ips {
		ipStrs = append(ipStrs, ip.String())
	}

	srv, err := zeroconf.RegisterProxy(
		name,
		"_http._tcp",
		"local.",
		dashboardPort,
		name,
		ipStrs,
		[]string{"path=/"},
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("mdns register: %w", err)
	}

	return &Publisher{
		Hostname: name + ".local",
		IPs:      ipStrs,
		server:   srv,
	}, nil
}

// Shutdown unregisters the mDNS service. Safe to call on a nil Publisher.
func (p *Publisher) Shutdown() {
	if p == nil || p.server == nil {
		return
	}
	p.server.Shutdown()
}

// canonicalHostname strips a trailing "." and ".local" so callers can pass
// "snugnas", "snugnas.local", or "snugnas.local." interchangeably.
func canonicalHostname(s string) string {
	return strings.TrimSuffix(strings.TrimSuffix(s, "."), ".local")
}

// lanIPv4 returns the IPv4 addresses of every non-loopback up interface.
func lanIPv4() ([]net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var out []net.IP
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			v4 := ipnet.IP.To4()
			if v4 == nil || v4.IsLinkLocalUnicast() {
				continue
			}
			out = append(out, v4)
		}
	}
	return out, nil
}
