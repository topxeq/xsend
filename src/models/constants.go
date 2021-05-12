package models

import (
	"context"
	"fmt"
	"net"
	"time"
)

// TCP_BUFFER_SIZE is the maximum packet size
const TCP_BUFFER_SIZE = 1024 * 64

// DEFAULT_RELAY is the default relay used (can be set using --relay)
var (
	DEFAULT_RELAY      = "croc.schollz.com"
	DEFAULT_RELAY6     = "croc6.schollz.com"
	DEFAULT_PORT       = "9009"
	DEFAULT_PASSPHRASE = "pass123"
)

// lookupTimeout for DNS requests
const lookupTimeout = time.Second

// publicDns are servers to be queried if a local lookup fails
var publicDns = []string{
	"1.0.0.1",                // Cloudflare
	"1.1.1.1",                // Cloudflare
	"8.8.4.4",                // Google
	"8.8.8.8",                // Google
	"8.26.56.26",             // Comodo
	"208.67.220.220",         // Cisco OpenDNS
	"208.67.222.222",         // Cisco OpenDNS
	"[2001:4860:4860::8844]", // Google
	"[2001:4860:4860::8888]", // Google
}

func init() {
	var err error
	DEFAULT_RELAY, err = lookup(DEFAULT_RELAY)
	if err == nil {
		DEFAULT_RELAY += ":" + DEFAULT_PORT
	} else {
		DEFAULT_RELAY = ""
	}
	DEFAULT_RELAY6, err = lookup(DEFAULT_RELAY6)
	if err == nil {
		DEFAULT_RELAY6 = "[" + DEFAULT_RELAY6 + "]:" + DEFAULT_PORT
	} else {
		DEFAULT_RELAY6 = ""
	}
}

// lookup an IP address.
//
// Priority is given to local queries, and the system falls back to a list of
// public DNS servers.
func lookup(address string) (ipaddress string, err error) {
	ipaddress, err = localLookupIP(address)
	if err == nil {
		return
	}
	err = nil

	result := make(chan string, len(publicDns))
	for _, dns := range publicDns {
		go func(dns string) {
			s, _ := remoteLookupIP(address, dns)
			result <- s
		}(dns)
	}

	for i := 0; i < len(publicDns); i++ {
		ipaddress = <-result
		if ipaddress != "" {
			return
		}
	}

	err = fmt.Errorf("failed to lookup %s at any DNS server", address)
	return
}

// localLookupIP returns a host's IP address based on the local resolver.
func localLookupIP(address string) (ipaddress string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), lookupTimeout)
	defer cancel()

	ip, err := net.DefaultResolver.LookupHost(ctx, address)
	if err != nil {
		return
	}
	ipaddress = ip[0]
	return
}

// remoteLookupIP returns a host's IP address based on a remote DNS server.
func remoteLookupIP(address, dns string) (ipaddress string, err error) {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: lookupTimeout,
			}
			return d.DialContext(ctx, "udp", dns+":53")
		},
	}
	ip, err := r.LookupHost(context.Background(), address)
	if err != nil {
		return
	}
	ipaddress = ip[0]
	return
}
