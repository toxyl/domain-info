package main

import (
	"container/ring"
	"context"
	"fmt"
	"net"
	"time"
)

type DNSResolver struct {
	Host     string
	Port     uint
	resolver *net.Resolver
}

func (dr *DNSResolver) ResolveHost(host string) ([]string, error) {
	return dr.resolver.LookupHost(context.Background(), host)
}

func (dr *DNSResolver) ResolveCNAME(host string) (string, error) {
	return dr.resolver.LookupCNAME(context.Background(), host)
}

func (dr *DNSResolver) ResolveAddr(addr string) ([]string, error) {
	return dr.resolver.LookupAddr(context.Background(), addr)
}

func (dr *DNSResolver) ResolveIPAddr(host string) ([]net.IPAddr, error) {
	return dr.resolver.LookupIPAddr(context.Background(), host)
}

func (dr *DNSResolver) ResolveIP(host string) ([]net.IP, error) {
	return dr.resolver.LookupIP(context.Background(), "ip", host)
}

func (dr *DNSResolver) ResolveIPv4(host string) ([]net.IP, error) {
	return dr.resolver.LookupIP(context.Background(), "ip4", host)
}

func (dr *DNSResolver) ResolveIPv6(host string) ([]net.IP, error) {
	return dr.resolver.LookupIP(context.Background(), "ip6", host)
}

func (dr *DNSResolver) ResolveMX(name string) ([]*net.MX, error) {
	return dr.resolver.LookupMX(context.Background(), name)
}

func (dr *DNSResolver) ResolveNS(name string) ([]*net.NS, error) {
	return dr.resolver.LookupNS(context.Background(), name)
}

func (dr *DNSResolver) ResolveTXT(name string) ([]string, error) {
	return dr.resolver.LookupTXT(context.Background(), name)
}

func (dr *DNSResolver) ResolveSRV(service, proto, name string) (string, []*net.SRV, error) {
	return dr.resolver.LookupSRV(context.Background(), service, proto, name)
}

func (dr *DNSResolver) ResolvePort(network, service string) (int, error) {
	return dr.resolver.LookupPort(context.Background(), network, service)
}

func NewDNSResolver(host string, port uint, timeout uint) *DNSResolver {
	return &DNSResolver{
		Host: host,
		Port: port,
		resolver: &net.Resolver{
			PreferGo:     true,
			StrictErrors: false,
			Dial: func(ctx context.Context, network string, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: time.Millisecond * time.Duration(timeout),
				}
				srv := fmt.Sprintf("%s:%d", host, port)
				return d.DialContext(ctx, network, srv)
			},
		},
	}
}

type DNSResolverRing struct {
	ring *ring.Ring
}

func (rs *DNSResolverRing) Init(resolvers ...*DNSResolver) {
	for _, s := range resolvers {
		rs.ring.Value = s
		rs.ring = rs.ring.Next()
	}
}

func (rs *DNSResolverRing) Next() *DNSResolver {
	rs.ring = rs.ring.Next()
	return rs.ring.Value.(*DNSResolver)
}

func (rs *DNSResolverRing) Print() {
	rs.ring.Do(func(x interface{}) {
		fmt.Println(x)
	})
}

func NewDNSResolverRing(resolvers ...*DNSResolver) *DNSResolverRing {
	rs := &DNSResolverRing{}
	rs.ring = ring.New(len(resolvers))
	rs.Init(resolvers...)
	return rs
}
