package l2

import (
	"net"
	"os/exec"
)

// Listener implements net.Listener for Bluetooth L2CAP, can be used to accept pair requests from another Blueooth device.
type Listener struct {
	s    *socket
	addr *Addr
}

// BluezInputBindHack toggles the restarting of the bluetooth service so that we may bind
var BluezInputBindHack = true

// NewListener returns a net.Listener for L2CAP connections
func NewListener(mac string, psm uint16) (net.Listener, error) {
	var s *socket

	addr, err := NewAddr(mac, psm)
	if err != nil {
		return nil, err
	}

	s, err = newSocket()
	if err != nil {
		return nil, err
	}

	if BluezInputBindHack {
		// restart bluetooth.service
		if err := exec.Command("systemctl", "restart", "bluetooth").Run(); err != nil {
			s.Close()
			return nil, err
		}
	}

	if err := s.Bind(addr); err != nil {
		s.Close()
		return nil, err
	}

	if err := s.Listen(1); err != nil {
		s.Close()
		return nil, err
	}

	return &Listener{s, addr}, nil
}

// Accept returns a new connection
func (l Listener) Accept() (net.Conn, error) {
	s, remoteAddr, err := l.s.Accept()
	if err != nil {
		return nil, err
	}
	return &Conn{s, l.addr, remoteAddr}, nil
}

// Addr returns the address of the Listener
func (l Listener) Addr() net.Addr {
	return l.addr
}

// Close closes the listener for further operation
func (l *Listener) Close() error {
	if l == nil {
		return nil
	}
	return l.s.Close()
}
