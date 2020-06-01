package l2

import (
	"io"
	"net"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

// Conn provides the net.Conn interface for L2CAP connections
type Conn struct {
	s    *socket
	live bool
}

var (
	ErrClosedConn = errors.New("l2: read/write on closed conn")
)

func newConn(s *socket) (*Conn, error) {
	if err := s.SetSendTimeout(500 * time.Millisecond); err != nil {
		return nil, err
	}

	if err := s.SetRecvTimeout(500 * time.Millisecond); err != nil {
		return nil, err
	}

	return &Conn{s: s, live: true}, nil
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

	return newConn(s)
}

// Write reads the next packet from the connection.
func (c Conn) Read(p []byte) (int, error) {
	if !c.live {
		return 0, ErrClosedConn
	}

	n, err := c.s.Recv(p)
	if n == 0 || errors.Is(err, unix.EWOULDBLOCK) {
		if err := c.die(); err != nil {
			return 0, err
		}
		return 0, io.EOF
	}
	if err != nil {
		return 0, err
	}

	return n, nil
}

// Write sends the packet p through the connection.
func (c Conn) Write(p []byte) (int, error) {
	if !c.live {
		return 0, ErrClosedConn
	}
	return c.s.Send(p)
}

// Close closes the connection to further reads and writes.
func (c *Conn) Close() error {
	if c == nil {
		return nil
	}

	if !c.live {
		return nil
	}

	return c.die()
}

func (c Conn) die() error {
	c.live = false
	return c.s.Close()
}

// LocalAddr returns the address of local side of the connection.
func (c Conn) LocalAddr() net.Addr {
	return &Addr{MAC: c.s.name.Addr, PSM: c.s.name.PSM}
}

// RemoteAddr returns the address of remote endpoint of the connection.
func (c Conn) RemoteAddr() net.Addr {
	return &Addr{MAC: c.s.peer.Addr, PSM: c.s.peer.PSM}
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
