package l2

import (
	"net"
	"time"

	"github.com/pkg/errors"
)

// Conn provides the net.Conn interface for L2CAP connections
type Conn struct {
	s          *socket
	localAddr  net.Addr
	remoteAddr net.Addr
}

// NewConn returns a L2CAP connection
func NewConn(mac string, psm uint16) (net.Conn, error) {
	remoteAddr, err := NewAddr(mac, psm)
	if err != nil {
		return nil, errors.Wrap(err, "bad address")
	}

	s, err := newSocket()
	if err != nil {
		return nil, errors.Wrap(err, "unable to open socket")
	}

	if err := s.Connect(remoteAddr); err != nil {
		s.Close()
		return nil, errors.Wrap(err, "unable to connect")
	}

	localAddr, err := s.Getsockname()
	if err != nil {
		return nil, err
	}

	return &Conn{s, localAddr, remoteAddr}, nil
}

// Write reads the next packet from the connection.
func (c Conn) Read(p []byte) (int, error) {
	return c.s.Recv(p)
}

// Write sends the packet p through the connection.
func (c Conn) Write(p []byte) (int, error) {
	return c.s.Send(p)
}

// Close closes the connection to further reads and writes.
func (c *Conn) Close() error {
	if c == nil {
		return nil
	}
	return c.s.Close()
}

// LocalAddr returns the address of local side of the connection.
func (c Conn) LocalAddr() net.Addr {
	return c.localAddr
}

// RemoteAddr returns the address of remote endpoint of the connection.
func (c Conn) RemoteAddr() net.Addr {
	return c.remoteAddr
}

// SetDeadline sets the deadline for Read and Write.
func (c Conn) SetDeadline(t time.Time) error {
	return errors.New("not supported")
}

// SetReadDeadline sets the deadline for Read.
func (c Conn) SetReadDeadline(t time.Time) error {
	return errors.New("not supported")
}

// SetWriteDeadline sets the deadline for Write.
func (c Conn) SetWriteDeadline(t time.Time) error {
	return errors.New("not supported")
}
