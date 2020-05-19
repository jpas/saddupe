package l2

import (
	"golang.org/x/sys/unix"
)

type socket struct {
	fd int
}

func newSocket() (*socket, error) {
	// TODO(jpas) Is unix.SOCK_CLOEXEC required?
	fd, err := unix.Socket(
		unix.AF_BLUETOOTH,
		unix.SOCK_SEQPACKET|unix.SOCK_CLOEXEC,
		unix.BTPROTO_L2CAP,
	)
	if err != nil {
		return nil, err
	}
	return &socket{fd}, nil
}

func (s socket) Send(p []byte) (int, error) {
	err := unix.Sendto(s.fd, p, 0, nil)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (s socket) Recv(p []byte) (int, error) {
	n, _, err := unix.Recvfrom(s.fd, p, 0)
	return n, err
}

func (s socket) Connect(addr *Addr) error {
	sa := sockaddrL2FromAddr(addr)
	return unix.Connect(s.fd, sa)
}

func (s socket) Listen(n int) error {
	return unix.Listen(s.fd, n)
}

func (s socket) Bind(addr *Addr) error {
	sa := sockaddrL2FromAddr(addr)
	return unix.Bind(s.fd, sa)
}

func (s socket) Accept() (*socket, *Addr, error) {
	fd, sa, err := unix.Accept(s.fd)
	if err != nil {
		return nil, nil, err
	}
	return &socket{fd}, l2AddrFromSockaddr(sa), nil
}

func (s *socket) Close() error {
	if s == nil {
		return nil
	}
	return unix.Close(s.fd)
}

func (s socket) Getsockname() (*Addr, error) {
	sa, err := unix.Getsockname(s.fd)
	if err != nil {
		return nil, err
	}
	return l2AddrFromSockaddr(sa), nil
}

func (s socket) Getpeername() (*Addr, error) {
	sa, err := unix.Getpeername(s.fd)
	if err != nil {
		return nil, err
	}
	return l2AddrFromSockaddr(sa), nil
}

// sockaddrL2FromAddr returns a unix.SockaddrL2 corresponding to addr, but with the Addr field in big endian.
//
// When SockaddrL2 is converted to RawSockaddrL2 it is converted back to little endian.
// This conversion is _not_ done when beign returned from Accept, Getpeername, or Getsockname.
func sockaddrL2FromAddr(addr *Addr) *unix.SockaddrL2 {
	sa := &unix.SockaddrL2{PSM: addr.PSM}
	for i := 0; i < 6; i++ {
		sa.Addr[i] = addr.MAC[i]
	}
	return sa
}

func l2AddrFromSockaddr(sa unix.Sockaddr) *Addr {
	l2sa := sa.(*unix.SockaddrL2)
	addr := &Addr{PSM: l2sa.PSM}
	for i := 0; i < 6; i++ {
		addr.MAC[i] = l2sa.Addr[5-i]
	}
	return addr
}
