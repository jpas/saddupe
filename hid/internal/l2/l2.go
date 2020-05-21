package l2

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
)

// Addr implements net.Addr for L2CAP Addresses
type Addr struct {
	MAC MAC
	PSM uint16
}

// NewAddr returns an Addr for L2CAP
func NewAddr(mac string, psm uint16) (*Addr, error) {
	m, err := ParseMAC(mac)
	if err != nil {
		return nil, err
	}
	return &Addr{*m, psm}, nil
}

func (a *Addr) String() string {
	return fmt.Sprintf("[%s]:%d", a.MAC, a.PSM)
}

// Network returns the "network" for L2CAP connections
func (a *Addr) Network() string {
	return "l2cap"
}

type MAC [6]byte

func ParseMAC(mac string) (*MAC, error) {
	hw, err := net.ParseMAC(mac)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse MAC address")
	}
	if len(hw) != 6 {
		return nil, errors.Errorf("invalid MAC address: %s", mac)
	}
	var m MAC
	for i := 0; i < 6; i++ {
		m[i] = hw[i]
	}

	return &m, nil
}

func (m MAC) Bytes() []byte {
	return m[:]
}

func (m MAC) String() string {
	return net.HardwareAddr(m.Bytes()).String()
}
