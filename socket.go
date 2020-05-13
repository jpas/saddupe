package main

import (
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
	"regexp"
	"strconv"
)

type BtAddr [6]byte

var BtAddrAny = BtAddr([6]byte{})

var btAddrPattern = regexp.MustCompile(`([[:xdigit:]]{2}):([[:xdigit:]]{2}):([[:xdigit:]]{2}):([[:xdigit:]]{2}):([[:xdigit:]]{2}):([[:xdigit:]]{2})`)

func NewBtAddr(addr string) (*BtAddr, error) {
	var b [6]byte

	match := btAddrPattern.FindStringSubmatch(addr)
	if len(match) == 0 {
		return nil, errors.Errorf("addr malformed %s", addr)
	}

	for i, m := range match[1:] {
		o, err := strconv.ParseUint(m[:2], 16, 8)
		if err != nil {
			return nil, errors.Wrap(err, "octect parse failed")
		}
		b[i] = byte(o)
	}

	a := BtAddr(b)
	return &a, nil
}

func (a BtAddr) Network() string {
	return "bt"
}

func (a BtAddr) String() string {
	b := [6]byte(a)
	return fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X", b[0], b[1], b[2], b[3], b[4], b[5])
}

func (a BtAddr) Bytes() [6]byte {
	return [6]byte(a)
}

type L2Addr struct {
	Addr BtAddr
	PSM  uint16
}

func NewL2Addr(addr string, psm uint16) (*L2Addr, error) {
	b, err := NewBtAddr(addr)
	if err != nil {
		return nil, errors.Wrap(err, "bad addr")
	}
	return &L2Addr{*b, psm}, nil
}

func (a L2Addr) Network() string {
	return "l2"
}

func (a L2Addr) String() string {
	return fmt.Sprintf("[%s]:%d", a.Addr, a.PSM)
}

func (a L2Addr) Sockaddr() unix.Sockaddr {
	sa := &unix.SockaddrL2{
		Addr: a.Addr,
		PSM:  a.PSM,
	}
	return unix.Sockaddr(sa)
}

type L2Socket struct {
	fd int
}

func NewL2Socket() (*L2Socket, error) {
	fd, err := unix.Socket(unix.AF_BLUETOOTH, unix.SOCK_SEQPACKET, unix.BTPROTO_L2CAP)
	if err != nil {
		return nil, err
	}
	return &L2Socket{fd}, nil
}

func (s *L2Socket) Read(b []byte) (int, error) {
	return unix.Read(s.fd, b)
}

func (s *L2Socket) Write(b []byte) (int, error) {
	return unix.Write(s.fd, b)
}

func (s *L2Socket) Close() error {
	return unix.Close(s.fd)
}

func (s *L2Socket) Bind(addr *L2Addr) error {
	if err := unix.SetsockoptInt(s.fd, unix.SOL_SOCKET, unix.SO_REUSEADDR, 1); err != nil {
		return errors.Wrap(err, "failed to set ReuseAddr")
	}
	return unix.Bind(s.fd, addr.Sockaddr())
}

func (s *L2Socket) Listen(n int) error {
	return unix.Listen(s.fd, n)
}

func (s *L2Socket) Accept() (*L2Socket, error) {
	fd, _, err := unix.Accept(s.fd)
	if err != nil {
		return nil, err
	}
	return &L2Socket{fd}, nil
}

func (s *L2Socket) Connect(addr *L2Addr) error {
	return unix.Connect(s.fd, addr.Sockaddr())
}

func (s *L2Socket) LocalAddr() *L2Addr {
	sn, err := unix.Getsockname(s.fd)
	if err != nil {
		panic(err)
	}
	sa := *sn.(*unix.SockaddrL2)

	la := &L2Addr{PSM: sa.PSM}
	// upstream forgets to reverse this on the way back up
	for i, b := range(sa.Addr) {
		la.Addr[5-i] = b
	}
	return la
}

func (s *L2Socket) RemoteAddr() *L2Addr {
	sn, err := unix.Getpeername(s.fd)
	if err != nil {
		panic(err)
	}
	sa := *sn.(*unix.SockaddrL2)

	la := &L2Addr{PSM: sa.PSM}
	// upstream forgets to reverse this on the way back up
	for i, b := range(sa.Addr) {
		la.Addr[5-i] = b
	}
	return la
}
