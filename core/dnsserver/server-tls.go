// Package dnsserver implements all the interfaces from Caddy, so that CoreDNS can be a servertype plugin.
package dnsserver

import (
	"fmt"
	"net"

	"github.com/miekg/dns"
)

// TLSServer represents an instance of a TLS-over-DNS-server.
type TLSServer struct {
	*Server
}

// NewTLSServer returns a new CoreDNS TLS server and compiles all middleware in to it.
func NewTLSServer(addr string, group []*Config) (*TLSServer, error) {

	s, err := NewServer(addr, group)
	if err != nil {
		return nil, err
	}

	return &TLSServer{Server: s}, nil
}

// Serve implements caddy.TCPServer interface.
func (s *TLSServer) Serve(l net.Listener) error {
	s.m.Lock()

	// Only fill out the TCP server for this one.
	s.server[tcp] = &dns.Server{Listener: l, Net: "tcp-tls", Handler: s.mux}
	s.m.Unlock()

	return s.server[tcp].ActivateAndServe()
}

// ServePacket This implements caddy.UDPServer interface.
func (s *TLSServer) ServePacket(p net.PacketConn) error { return nil }

// Listen implements caddy.TCPServer interface.
func (s *TLSServer) Listen() (net.Listener, error) {

	// Remove, but show our 'tls' directive has been picked up.
	for _, conf := range s.zones {
		fmt.Printf("%q\n", conf.TLSConfig)
	}

	l, err := net.Listen("tcp", s.Addr[len(ProtoTLS+"://"):])
	if err != nil {
		return nil, err
	}
	return l, nil
}

// ListenPacket implements caddy.UDPServer interface.
func (s *TLSServer) ListenPacket() (net.PacketConn, error) { return nil, nil }

// OnStartupComplete lists the sites served by this server
// and any relevant information, assuming Quiet is false.
func (s *TLSServer) OnStartupComplete() {
	if Quiet {
		return
	}

	for zone, config := range s.zones {
		fmt.Println(ProtoTLS + "://" + zone + ":" + config.Port)
	}
}
