package l2

import (
	"errors"
	"time"

	"golang.org/x/sys/unix"
)

type socket struct {
	fd   int
	name *unix.SockaddrL2
	peer *unix.SockaddrL2
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
	return newSocketFromFd(fd), nil
}

func newSocketFromFd(fd int) *socket {
	s := &socket{fd: fd}

	name, err := s.getsockname()
	if err == nil {
		s.name = name
	}

	peer, err := s.getpeername()
	if err == nil {
		s.peer = peer
	}

	return s
}

func (s *socket) Send(p []byte) (int, error) {
	if s.peer == nil {
		return 0, unix.ENOTCONN
	}
	err := unix.Sendto(s.fd, p, 0, s.peer)
	for errors.Is(err, unix.EINTR) {
		err = unix.Sendto(s.fd, p, 0, s.peer)
	}
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (s *socket) Recv(p []byte) (int, error) {
	n, _, err := unix.Recvfrom(s.fd, p, 0)
	for errors.Is(err, unix.EINTR) {
		n, _, err = unix.Recvfrom(s.fd, p, 0)
	}
	if err != nil {
		return 0, err
	}
	return n, err
}

func (s *socket) timeout(opt int) (time.Duration, error) {
	tv, err := unix.GetsockoptTimeval(s.fd, unix.SOL_SOCKET, opt)
	if err != nil {
		return 0, err
	}
	return time.Duration(tv.Nano()), nil
}

func (s *socket) setTimeout(opt int, t time.Duration) error {
	tv := unix.NsecToTimeval(t.Nanoseconds())
	return unix.SetsockoptTimeval(s.fd, unix.SOL_SOCKET, opt, &tv)
}

func (s *socket) SendTimeout() (time.Duration, error) {
	return s.timeout(unix.SO_SNDTIMEO)
}

func (s *socket) SetSendTimeout(t time.Duration) error {
	return s.setTimeout(unix.SO_SNDTIMEO, t)
}

func (s *socket) RecvTimeout() (time.Duration, error) {
	return s.timeout(unix.SO_RCVTIMEO)
}

func (s *socket) SetRecvTimeout(t time.Duration) error {
	return s.setTimeout(unix.SO_RCVTIMEO, t)
}

func (s *socket) Connect(addr *Addr) error {
	sa := addr.sockaddrL2()

	// When SockaddrL2 is converted to RawSockaddrL2 it's byte order is swapped, but
	// we already have it in network order so we must swap it so that they can put it
	// back into network order
	for i := 0; i < 3; i++ {
		sa.Addr[i], sa.Addr[5-i] = sa.Addr[5-i], sa.Addr[i]
	}

	err := unix.Connect(s.fd, sa)
	for errors.Is(err, unix.EINTR) {
		err = unix.Connect(s.fd, sa)
	}

	name, err := s.getsockname()
	if err != nil {
		return err
	}
	s.name = name

	peer, err := s.getpeername()
	if err != nil {
		return err
	}
	s.peer = peer

	return err
}

func (s *socket) Listen(n int) error {
	return unix.Listen(s.fd, n)
}

func (s *socket) Bind(addr *Addr) error {
	sa := addr.sockaddrL2()

	// When SockaddrL2 is converted to RawSockaddrL2 it's byte order is swapped, but
	// we already have it in network order so we must swap it so that they can put it
	// back into network order
	for i := 0; i < 3; i++ {
		sa.Addr[i], sa.Addr[5-i] = sa.Addr[5-i], sa.Addr[i]
	}

	err := unix.Bind(s.fd, sa)
	if err != nil {
		return err
	}

	name, err := s.getsockname()
	if err != nil {
		return err
	}
	s.name = name

	return nil
}

func (s *socket) Accept() (*socket, error) {
	fd, _, err := unix.Accept(s.fd)
	for errors.Is(err, unix.EINTR) {
		fd, _, err = unix.Accept(s.fd)
	}
	if err != nil {
		return nil, err
	}
	return newSocketFromFd(fd), nil
}

func (s *socket) Close() error {
	if s == nil {
		return nil
	}
	return unix.Close(s.fd)
}

func newSockaddrL2(sa unix.Sockaddr) (*unix.SockaddrL2, error) {
	l2, ok := sa.(*unix.SockaddrL2)
	if !ok {
		return nil, errors.New("not a *unix.SockaddrL2")
	}
	return l2, nil
}

func (s *socket) getsockname() (*unix.SockaddrL2, error) {
	sa, err := unix.Getpeername(s.fd)
	if err != nil {
		return nil, err
	}
	return newSockaddrL2(sa)
}

func (s *socket) getpeername() (*unix.SockaddrL2, error) {
	sa, err := unix.Getpeername(s.fd)
	if err != nil {
		return nil, err
	}
	return newSockaddrL2(sa)
}
