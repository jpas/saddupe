package l2

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
)

// Addr implements net.Addr for L2CAP Addresses
type Addr struct {
	MAC net.HardwareAddr
	PSM uint16
}

// NewAddr returns an Addr for L2CAP
func NewAddr(mac string, psm uint16) (*Addr, error) {
	if len(mac) != 17 {
		return nil, errors.Errorf("invalid MAC address: %s", mac)
	}
	hw, err := net.ParseMAC(mac)
	if err != nil {
		return nil, err
	}
	return &Addr{hw, psm}, nil
}

func (a *Addr) String() string {
	hw := net.HardwareAddr(a.MAC[:]) // does propper string formatting for us
	return fmt.Sprintf("[%s]:%d", hw, a.PSM)
}

// Network returns the "network" for L2CAP connections
func (a *Addr) Network() string {
	return "l2cap"
}
