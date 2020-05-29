package l2

import (
	"fmt"

	"github.com/jpas/saddupe/hw"
	"golang.org/x/sys/unix"
)

// Addr implements net.Addr for L2CAP Addresses
type Addr struct {
	MAC hw.MAC
	PSM uint16
}

// NewAddr returns an Addr for L2CAP
func NewAddr(mac string, psm uint16) (*Addr, error) {
	m, err := hw.ParseMAC(mac)
	if err != nil {
		return nil, err
	}
	return &Addr{*m, psm}, nil
}

func (a *Addr) String() string {
	return fmt.Sprintf("[%02x]:%d", a.MAC, a.PSM)
}

// Network returns the "network" for L2CAP connections
func (a *Addr) Network() string {
	return "l2cap"
}

func (a *Addr) sockaddrL2() *unix.SockaddrL2 {
	return &unix.SockaddrL2{Addr: a.MAC, PSM: a.PSM}
}
